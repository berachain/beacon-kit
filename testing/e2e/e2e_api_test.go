// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package e2e_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	beaconapi "github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/berachain/beacon-kit/config/spec"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/ethereum/go-ethereum/common"
)

const localHost = "localhost"

// BeaconHTTPClient wraps http.Client with baseURL.
// This is needed to change the default baseURL of the http.Client.
type BeaconHTTPClient struct {
	*http.Client
	baseURL string
}

// Get overrides http.Client.Get to use baseURL
func (c *BeaconHTTPClient) Get(path string) (*http.Response, error) {
	return c.Client.Get(c.baseURL + path)
}

// initHTTPBeaconTest initializes the http client for the beacon node api.
// It gets the public ports from the consensus client and creates a http client with the baseURL.
// This is needed where we want to test the node-api directly as
// some of the methods are not present in go-eth2-client library.
func (s *BeaconKitE2ESuite) initHTTPBeaconTest() *BeaconHTTPClient {
	// Initialize consensus client to get public ports.
	ports := s.initBeaconTest().GetPublicPorts()

	// As kurtosis assigns public port randomly, we need to get the port from the consensus client.
	// Get node-api port and format it.
	portStr := ports["node-api"].String()
	portNum := strings.Replace(portStr, "/0", "", 1)

	hostPort := net.JoinHostPort(localHost, portNum)
	// Create client with baseURL
	return &BeaconHTTPClient{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: "http://" + hostPort,
	}
}

// initBeaconTest initializes the any tests for the beacon node api.
func (s *BeaconKitE2ESuite) initBeaconTest() *types.ConsensusClient {
	// Wait for execution block 5.
	err := s.WaitForFinalizedBlockNumber(5)
	s.Require().NoError(err)

	// Get the consensus client.
	client := s.ConsensusClients()[config.ClientValidator0]
	s.Require().NotNil(client)

	return client
}

// TestBeaconStateRoot tests the beacon node api for beacon state root.
func (s *BeaconKitE2ESuite) TestBeaconStateRoot() {
	client := s.initBeaconTest()

	// Ensure the state root is not nil.
	stateRootResp, err := client.BeaconStateRoot(
		s.Ctx(),
		&beaconapi.BeaconStateRootOpts{
			State: utils.StateIDHead,
		},
	)
	s.Require().NoError(err)
	s.Require().NotEmpty(stateRootResp)
	s.Require().False(stateRootResp.Data.IsZero())
}

// TestBeaconValidatorsWithIndices tests the beacon node api for beacon validators with indices.
func (s *BeaconKitE2ESuite) TestBeaconValidatorsWithIndices() {
	client := s.initBeaconTest()

	indices := []phase0.ValidatorIndex{0}
	validatorsResp, err := client.Validators(
		s.Ctx(),
		&beaconapi.ValidatorsOpts{
			State:   utils.StateIDHead,
			Indices: indices,
		},
	)
	s.Require().NoError(err)
	s.Require().NotNil(validatorsResp)

	validatorData := validatorsResp.Data
	s.Require().NotNil(validatorData, "Validator data should not be nil")
	s.Require().
		Len(validatorData, len(indices), "Number of validator responses should match number of requested indices")

	validator := validatorData[0]
	s.Require().NotNil(validator, "Validator should not be nil")
	s.Require().Equal(phase0.ValidatorIndex(0), validator.Index, "Should be validator index 0")

	s.Require().
		NotEmpty(validator.Validator.PublicKey, "Validator public key should not be empty")
	s.Require().
		Len(validator.Validator.PublicKey, 48, "Validator public key should be 48 bytes long")

	s.Require().
		NotEmpty(validator.Validator.WithdrawalCredentials, "Withdrawal credentials should not be empty")
	s.Require().
		Len(validator.Validator.WithdrawalCredentials, 32, "Withdrawal credentials should be 32 bytes long")

	s.Require().
		True(validator.Validator.EffectiveBalance > 0, "Effective balance should be positive")
	s.Require().
		True(validator.Validator.EffectiveBalance <= 32e9, "Effective balance should not exceed 32 ETH")

	s.Require().
		False(validator.Validator.Slashed, "Slashed status should not be true")

	s.Require().
		True(validator.Validator.ActivationEpoch >= validator.Validator.ActivationEligibilityEpoch,
			"Activation epoch should be greater than or equal to activation eligibility epoch")

	s.Require().
		True(validator.Validator.WithdrawableEpoch >= validator.Validator.ExitEpoch,
			"Withdrawable epoch should be greater than or equal to exit epoch")

	s.Require().
		NotEmpty(validator.Status, "Validator status should not be empty")

	s.Require().
		True(validator.Balance > 0, "Validator balance should be positive")
	s.Require().
		True(validator.Balance <= 32e9, "Validator balance should not exceed 32 ETH")
}

// TestValidatorsEmptyIndicesAndStatuses tests that querying validators with empty indices and empty statuses returns all validators.
// Empty indices and statuses is same as not populating the indices and statuses. Basically, querying by State.
func (s *BeaconKitE2ESuite) TestValidatorsEmptyIndicesAndStatuses() {
	client := s.initBeaconTest()

	// Query validators with empty indices and empty statuses
	emptyIndices := []phase0.ValidatorIndex{}
	emptyStatuses := []apiv1.ValidatorState{}

	validatorsResp, err := client.Validators(
		s.Ctx(),
		&beaconapi.ValidatorsOpts{
			State:           utils.StateIDHead,
			Indices:         emptyIndices,
			ValidatorStates: emptyStatuses,
		},
	)

	s.Require().NoError(err)
	s.Require().NotNil(validatorsResp)

	// Verify we got all validators
	validatorData := validatorsResp.Data
	s.Require().NotNil(validatorData, "Validator data should not be nil")
	s.Require().Equal(config.NumValidators, len(validatorData),
		"Should return all validators when using empty statuses and empty indices")

	// Verify each validator has required fields
	for _, validator := range validatorData {
		s.Require().NotNil(validator, "Validator should not be nil")
		s.Require().NotEmpty(validator.Validator.PublicKey, "Validator public key should not be empty")
		s.Require().Len(validator.Validator.PublicKey, 48, "Validator public key should be 48 bytes long")
		s.Require().NotEmpty(validator.Validator.WithdrawalCredentials,
			"Withdrawal credentials should not be empty")
		s.Require().Len(validator.Validator.WithdrawalCredentials, 32,
			"Withdrawal credentials should be 32 bytes long")
		s.Require().True(validator.Validator.EffectiveBalance > 0,
			"Effective balance should be positive")
	}
}

// TestValidatorsWithMultipleIndices tests querying multiple specific validator indices.
func (s *BeaconKitE2ESuite) TestValidatorsWithMultipleIndices() {
	client := s.initBeaconTest()
	indices := []phase0.ValidatorIndex{0, 1, 2}

	validatorsResp, err := client.Validators(s.Ctx(), &beaconapi.ValidatorsOpts{
		State:   utils.StateIDHead,
		Indices: indices,
	})
	s.Require().NoError(err)
	s.Require().NotNil(validatorsResp)
	s.Require().Len(validatorsResp.Data, len(indices))
}

// TestValidatorsWithInvalidIndex tests querying a non-existent validator index
// This should return an empty list of validators.
func (s *BeaconKitE2ESuite) TestValidatorsWithInvalidIndex() {
	client := s.initBeaconTest()
	indices := []phase0.ValidatorIndex{999999} // Invalid index

	validatorsResp, err := client.Validators(s.Ctx(), &beaconapi.ValidatorsOpts{
		State:   utils.StateIDHead,
		Indices: indices,
	})
	s.Require().NoError(err)
	s.Require().NotNil(validatorsResp)

	// No validators returned
	s.Require().Len(validatorsResp.Data, 0)
}

// TestValidatorsWithSpecificStatus tests filtering validators by status.
func (s *BeaconKitE2ESuite) TestValidatorsWithSpecificStatus() {
	client := s.initBeaconTest()

	validatorsResp, err := client.Validators(s.Ctx(), &beaconapi.ValidatorsOpts{
		State:           utils.StateIDHead,
		ValidatorStates: []apiv1.ValidatorState{apiv1.ValidatorStateActiveOngoing},
	})
	s.Require().NoError(err)
	s.Require().NotNil(validatorsResp)

	// Verify all returned validators have the requested status
	for _, validator := range validatorsResp.Data {
		s.Require().Equal(apiv1.ValidatorStateActiveOngoing, validator.Status)
	}
}

// TestValidatorBalances tests querying validator balances.
func (s *BeaconKitE2ESuite) TestValidatorBalances() {
	client := s.initBeaconTest()

	balancesResp, err := client.ValidatorBalances(s.Ctx(), &beaconapi.ValidatorBalancesOpts{
		State: utils.StateIDHead,
	})
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp)

	// Verify the response is not empty
	s.Require().NotNil(balancesResp.Data)
	s.Require().NotEmpty(balancesResp.Data)

	balanceMap := balancesResp.Data
	for _, balance := range balanceMap {
		s.Require().True(balance > 0, "Validator balance should be positive")
		// 4e12 Gwei = 4 * 10^12 Gwei = 4,000,000,000,000 Gwei = 4000 BERA
		s.Require().True(balance <= 4e12, "Validator balance should not exceed 4e12 gwei (4000 BERA)")
	}
}

// TestValidatorBalancesWithSpecificIndices tests querying validator balances with specific indices.
func (s *BeaconKitE2ESuite) TestValidatorBalancesWithSpecificIndices() {
	client := s.initBeaconTest()

	indices := []phase0.ValidatorIndex{0}

	balancesResp, err := client.ValidatorBalances(s.Ctx(), &beaconapi.ValidatorBalancesOpts{
		State:   utils.StateIDHead,
		Indices: indices,
	})
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp)

	// Verify the response is not empty
	s.Require().NotNil(balancesResp.Data)
	s.Require().Len(balancesResp.Data, len(indices))

	// Verify balance data
	for index, balance := range balancesResp.Data {
		s.Require().NotNil(balance)
		s.Require().Contains(indices, index)
		s.Require().True(balance > 0, "Validator balance should be positive")
		// 4e12 Gwei = 4 * 10^12 Gwei = 4,000,000,000,000 Gwei = 4000 BERA
		s.Require().True(balance <= 4e12, "Validator balance should not exceed 4000 BERA")
	}
}

// TestValidatorBalancesMultipleIndices tests querying balances for multiple validator indices.
func (s *BeaconKitE2ESuite) TestValidatorBalancesMultipleIndices() {
	client := s.initBeaconTest()
	indices := []phase0.ValidatorIndex{0, 1, 2}

	balancesResp, err := client.ValidatorBalances(
		s.Ctx(),
		&beaconapi.ValidatorBalancesOpts{
			State:   utils.StateIDHead,
			Indices: indices,
		},
	)

	s.Require().NoError(err)
	s.Require().NotNil(balancesResp)
	s.Require().Len(balancesResp.Data, len(indices))

	// Verify all requested indices are present
	returnedIndices := make(map[phase0.ValidatorIndex]struct{})
	for index, balance := range balancesResp.Data {
		returnedIndices[index] = struct{}{}
		s.Require().True(balance > 0)
		// 4e12 Gwei = 4 * 10^12 Gwei = 4,000,000,000,000 Gwei = 4000 BERA
		s.Require().True(balance <= 4e12, "Validator balance should not exceed 4000 BERA")
	}
	for _, idx := range indices {
		_, exists := returnedIndices[idx]
		s.Require().True(exists, "Expected validator index not found in response")
	}
}

// TestValidatorBalancesWithInvalidIndex tests querying validator balances with an invalid index.
func (s *BeaconKitE2ESuite) TestValidatorBalancesWithInvalidIndex() {
	client := s.initBeaconTest()

	indices := []phase0.ValidatorIndex{999999} // Invalid index

	balancesResp, err := client.ValidatorBalances(s.Ctx(), &beaconapi.ValidatorBalancesOpts{
		State:   utils.StateIDHead,
		Indices: indices,
	})
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp)
	// Should return an empty list of balances
	s.Require().Len(balancesResp.Data, 0)
}

// TestValidatorBalancesWithPubkey tests querying validator balances using a public key.
func (s *BeaconKitE2ESuite) TestValidatorBalancesWithPubkey() {
	client := s.initBeaconTest()

	// First call validators to get the validator public key
	validatorsResp, err := client.Validators(s.Ctx(), &beaconapi.ValidatorsOpts{
		State: utils.StateIDHead,
	})
	s.Require().NoError(err)
	s.Require().NotNil(validatorsResp)

	validator := validatorsResp.Data[0]
	s.Require().NotNil(validator)
	pubkey := validator.Validator.PublicKey

	balancesResp, err := client.ValidatorBalances(s.Ctx(), &beaconapi.ValidatorBalancesOpts{
		State:   utils.StateIDHead,
		Indices: []phase0.ValidatorIndex{},
		PubKeys: []phase0.BLSPubKey{pubkey},
	})
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp)
	s.Require().Len(balancesResp.Data, 1)

	// Verify balance data
	for _, balance := range balancesResp.Data {
		s.Require().NotNil(balance)
		s.Require().True(balance > 0, "Validator balance should be positive")
		// 4e12 Gwei = 4 * 10^12 Gwei = 4,000,000,000,000 Gwei = 4000 BERA
		s.Require().True(balance <= 4e12, "Validator balance should not exceed 4000 BERA")
	}
}

// TestValidatorBalancesWithInvalidPubkey tests querying validator balances using a public key.
func (s *BeaconKitE2ESuite) TestValidatorBalancesWithInvalidPubkey() {
	client := s.initBeaconTest()

	// Example validator pubkey (48 bytes with 0x prefix)
	notFoundPubkey := "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a"

	balancesResp, err := client.ValidatorBalances(s.Ctx(), &beaconapi.ValidatorBalancesOpts{
		State:   utils.StateIDHead,
		Indices: []phase0.ValidatorIndex{},
		PubKeys: []phase0.BLSPubKey{phase0.BLSPubKey(common.FromHex(notFoundPubkey))},
	})
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp)
	// Should return an empty list of balances
	s.Require().Len(balancesResp.Data, 0)
}

// Helper functions
// decodeResponse is a generic function that decodes an HTTP response into the specified type T.
func decodeResponse[T any](resp *http.Response) (T, error) {
	var result T

	if resp == nil {
		return result, errors.New("nil response")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	// First decode into GenericResponse.
	var genericResp beacontypes.GenericResponse
	if err = json.Unmarshal(bodyBytes, &genericResp); err != nil {
		return result, err
	}

	// Convert the data field to JSON.
	dataBytes, err := json.Marshal(genericResp.Data)
	if err != nil {
		return result, err
	}

	// Unmarshal into the target type.
	if err = json.Unmarshal(dataBytes, &result); err != nil {
		return result, err
	}

	return result, nil
}

// decodeValidatorResponse decodes a single validator response.
func (s *BeaconKitE2ESuite) decodeValidatorResponse(resp *http.Response) (*beacontypes.ValidatorData, error) {
	validator, err := decodeResponse[beacontypes.ValidatorData](resp)
	if err != nil {
		return nil, err
	}
	return &validator, nil
}

// getStateValidator gets the state validator by index or pubkey.
func (s *BeaconKitE2ESuite) getStateValidator(stateID, validatorID string) (*http.Response, error) {
	client := s.initHTTPBeaconTest()

	resp, err := client.Get(fmt.Sprintf("/eth/v1/beacon/states/%s/validators/%s", stateID, validatorID))
	if err != nil {
		return nil, fmt.Errorf("failed to get validator: %w", err)
	}
	if resp == nil {
		return nil, errors.New("received nil response")
	}

	return resp, nil
}

// TestGetStateValidatorByIndex tests getting the state validator by index.
func (s *BeaconKitE2ESuite) TestGetStateValidatorByIndex() {
	resp, err := s.getStateValidator(utils.StateIDHead, "0")
	s.Require().NoError(err)
	s.Require().NotNil(resp, "response should not be nil")
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	validatorResp, err := s.decodeValidatorResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(validatorResp, "validator response should not be nil")

	s.Require().Equal(uint64(0), validatorResp.Index)
	s.Require().Equal("active_ongoing", validatorResp.Status)
}

// TestGetStateValidatorBySlotAndIndex tests getting the state validator by slot and index.
func (s *BeaconKitE2ESuite) TestGetStateValidatorBySlotAndIndex() {
	// Here the stateID is a slot number which is 10.
	resp, err := s.getStateValidator("10", "0")
	s.Require().NoError(err)
	s.Require().NotNil(resp, "response should not be nil")
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	validatorResp, err := s.decodeValidatorResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(validatorResp, "validator response should not be nil")

	s.Require().Equal(uint64(0), validatorResp.Index)
	s.Require().Equal("active_ongoing", validatorResp.Status)
}

// TestGetStateValidatorByPubkey tests getting the state validator by pubkey.
func (s *BeaconKitE2ESuite) TestGetStateValidatorByPubkey() {
	// First call validators to get the validator public key
	resp, err := s.getStateValidator(utils.StateIDHead, "0")
	s.Require().NoError(err)
	s.Require().NotNil(resp, "response should not be nil")
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	validatorResp, err := s.decodeValidatorResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(validatorResp, "validator response should not be nil")

	// Retrieve the public key.
	pubkey := validatorResp.Validator.PublicKey

	// Actual test starts here.
	resp, err = s.getStateValidator(utils.StateIDHead, pubkey)
	s.Require().NoError(err)
	s.Require().NotNil(resp, "response should not be nil")
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	validatorResp, errInDecode := s.decodeValidatorResponse(resp)
	s.Require().NoError(errInDecode)
	s.Require().NotNil(validatorResp, "validator response should not be nil")

	s.Require().Equal("active_ongoing", validatorResp.Status)
}

// TestGetStateValidatorInvalidID tests getting the state validator with an invalid id.
func (s *BeaconKitE2ESuite) TestGetStateValidatorInvalidID() {
	resp, err := s.getStateValidator(utils.StateIDHead, "invalid_id")
	s.Require().NoError(err)
	s.Require().NotNil(resp, "response should not be nil")
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()
}

// getValidatorBalances gets the validator balances for the given stateID and optional ids.
func (s *BeaconKitE2ESuite) getValidatorBalances(stateID string, ids ...string) (*http.Response, error) {
	client := s.initHTTPBeaconTest()

	url := fmt.Sprintf("/eth/v1/beacon/states/%s/validator_balances", stateID)

	// Add ID parameters if provided
	if len(ids) > 0 {
		queryParams := make([]string, 0, len(ids))
		for _, id := range ids {
			queryParams = append(queryParams, "id="+id)
		}
		url = url + "?" + strings.Join(queryParams, "&")
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get validator balances: %w", err)
	}
	if resp == nil {
		return nil, errors.New("received nil response")
	}

	return resp, nil
}

// decodeValidatorBalancesResponse decodes a response containing validator balances.
func (s *BeaconKitE2ESuite) decodeValidatorBalancesResponse(resp *http.Response) (*[]beacontypes.ValidatorBalanceData, error) {
	balances, err := decodeResponse[[]beacontypes.ValidatorBalanceData](resp)
	if err != nil {
		return nil, err
	}
	return &balances, nil
}

// TestGetValidatorBalances tests querying validator balances for state head.
func (s *BeaconKitE2ESuite) TestGetValidatorBalances() {
	resp, err := s.getValidatorBalances(utils.StateIDHead)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	balancesResp, err := s.decodeValidatorBalancesResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp, "balances response should not be nil")
	s.Require().NotEmpty(balancesResp, "balances response should not be empty")

	for _, balance := range *balancesResp {
		s.Require().True(balance.Balance > 0, "Validator balance should be positive")
		// 4e12 Gwei = 4 * 10^12 Gwei = 4,000,000,000,000 Gwei = 4000 BERA
		s.Require().True(balance.Balance <= 4e12, "Validator balance should not exceed 4000 BERA")
	}
}

// TestGetValidatorBalancesWithSpecificID tests querying validator balances with specific ID.
func (s *BeaconKitE2ESuite) TestGetValidatorBalancesWithSpecificID() {
	resp, err := s.getValidatorBalances(utils.StateIDHead, "0")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	balancesResp, err := s.decodeValidatorBalancesResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp, "balances response should not be nil")

	s.Require().Len(*balancesResp, 1)
	s.Require().True((*balancesResp)[0].Balance > 0, "Validator balance should be positive")
	// 4e12 Gwei = 4 * 10^12 Gwei = 4,000,000,000,000 Gwei = 4000 BERA
	s.Require().True((*balancesResp)[0].Balance <= 4e12, "Validator balance should not exceed 4000 BERA")
}

// TestGetValidatorBalancesWithMultipleIDs tests querying validator balances with multiple IDs.
func (s *BeaconKitE2ESuite) TestGetValidatorBalancesWithMultipleIDs() {
	resp, err := s.getValidatorBalances(utils.StateIDHead, "0", "1")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	balancesResp, err := s.decodeValidatorBalancesResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp, "balances response should not be nil")
	s.Require().NotEmpty(balancesResp, "balances response should not be empty")

	// The response should contain 2 validator balances.
	s.Require().Len(*balancesResp, 2)

	for _, balance := range *balancesResp {
		s.Require().True(balance.Balance > 0, "Validator balance should be positive")
		// 4e12 Gwei = 4 * 10^12 Gwei = 4,000,000,000,000 Gwei = 4000 BERA
		s.Require().True(balance.Balance <= 4e12, "Validator balance should not exceed 4000 BERA")
	}
}

// TestGetValidatorBalancesWithInvalidID tests querying validator balances with invalid ID.
func (s *BeaconKitE2ESuite) TestGetValidatorBalancesWithInvalidID() {
	resp, err := s.getValidatorBalances(utils.StateIDHead, "invalid_id")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	balancesResp, err := s.decodeValidatorBalancesResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp, "balances response should not be nil")
	s.Require().Len(*balancesResp, 0)
}

// TestGetValidatorBalancesWithNonExistentIndex tests querying validator balances with non-existent index.
func (s *BeaconKitE2ESuite) TestGetValidatorBalancesWithNonExistentIndex() {
	resp, err := s.getValidatorBalances(utils.StateIDHead, "99999")
	s.Require().NoError(err)
	// If an index does not match any validator, no balance will be returned but this will not cause an error.
	// The response should be 200 OK with empty data.
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	balancesResp, err := s.decodeValidatorBalancesResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp, "balances response should not be nil")
	s.Require().Len(*balancesResp, 0)
}

// TestGetValidatorBalancesWithPublicKey tests querying validator balances with public key.
func (s *BeaconKitE2ESuite) TestGetValidatorBalancesWithPublicKey() {
	client := s.initBeaconTest()

	// First call validators to get the validator public key
	validatorsResp, err := client.Validators(s.Ctx(), &beaconapi.ValidatorsOpts{
		State: utils.StateIDHead,
	})
	s.Require().NoError(err)
	s.Require().NotNil(validatorsResp)

	validator := validatorsResp.Data[0]
	s.Require().NotNil(validator)
	pubkey := validator.Validator.PublicKey

	resp, err := s.getValidatorBalances(utils.StateIDHead, pubkey.String())
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	balancesResp, err := s.decodeValidatorBalancesResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp, "balances response should not be nil")
	s.Require().NotEmpty(balancesResp, "balances response should not be empty")

	s.Require().Len(*balancesResp, 1)
	s.Require().True((*balancesResp)[0].Balance > 0, "Validator balance should be positive")
	// 4e12 Gwei = 4 * 10^12 Gwei = 4,000,000,000,000 Gwei = 4000 BERA
	s.Require().True((*balancesResp)[0].Balance <= 4e12, "Validator balance should not exceed 4000 BERA")
}

// TestGetValidatorBalancesWithInvalidPublicKey tests querying validator balances with invalid public key.
func (s *BeaconKitE2ESuite) TestGetValidatorBalancesWithInvalidPublicKey() {
	// Example validator pubkey (48 bytes with 0x prefix)
	notFoundPubkey := "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a"

	resp, err := s.getValidatorBalances(utils.StateIDHead, notFoundPubkey)
	s.Require().NoError(err)
	// If public key does not match any validator, no balance will be returned but this will not cause an error.
	// The response should be 200 OK with empty data.
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	balancesResp, err := s.decodeValidatorBalancesResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp, "balances response should not be nil")
	s.Require().Len(*balancesResp, 0)
}

// TestGetValidatorBalancesForGenesis tests querying validator balances for state genesis.
func (s *BeaconKitE2ESuite) TestGetValidatorBalancesForGenesis() {
	resp, err := s.getValidatorBalances(utils.StateIDGenesis)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	balancesResp, err := s.decodeValidatorBalancesResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(balancesResp, "balances response should not be nil")
	s.Require().NotEmpty(balancesResp, "balances response should not be empty")

	for _, balance := range *balancesResp {
		s.Require().True(balance.Balance > 0, "Validator balance should be positive")
		// At genesis, the validator balance is 32 BERA.
		// 32e9 Gwei = 32 * 10^9 Gwei = 32,000,000,000 Gwei = 32 BERA
		s.Require().True(balance.Balance <= 32e9, "Validator balance should not exceed 32 BERA")
	}
}

// getValidator gets the validator with the given stateID and optional options.
func (s *BeaconKitE2ESuite) getValidator(stateID string, options ...map[string]string) (*http.Response, error) {
	client := s.initHTTPBeaconTest()

	url := fmt.Sprintf("/eth/v1/beacon/states/%s/validators", stateID)

	// Process options if provided
	if len(options) > 0 {
		queryParams := []string{}

		// Extract options from the map
		for _, opts := range options {
			for key, value := range opts {
				queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
			}
		}

		// Add query parameters to URL if any exist
		if len(queryParams) > 0 {
			url = url + "?" + strings.Join(queryParams, "&")
		}
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get validator: %w", err)
	}
	if resp == nil {
		return nil, errors.New("received nil response")
	}
	return resp, nil
}

// decodeValidatorsResponse decodes a response containing multiple validators.
func (s *BeaconKitE2ESuite) decodeValidatorsResponse(resp *http.Response) (*[]beacontypes.ValidatorData, error) {
	validators, err := decodeResponse[[]beacontypes.ValidatorData](resp)
	if err != nil {
		return nil, err
	}
	return &validators, nil
}

// TestGetValidatorsWithStateHead tests querying validators with state head.
func (s *BeaconKitE2ESuite) TestGetValidatorsWithStateHead() {
	resp, err := s.getValidator(utils.StateIDHead)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	validators, err := s.decodeValidatorsResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(validators)

	for _, validator := range *validators {
		s.Require().True(validator.Status == "active_ongoing")
		s.Require().True(validator.Balance > 0, "Validator balance should be positive")
		// 4e12 Gwei = 4 * 10^12 Gwei = 4,000,000,000,000 Gwei = 4000 BERA
		s.Require().True(validator.Balance <= 4e12, "Validator balance should not exceed 4000 BERA")
	}
}

// TestGetValidatorsWithID tests querying validators with ID parameter.
func (s *BeaconKitE2ESuite) TestGetValidatorsWithID() {
	resp, err := s.getValidator(utils.StateIDHead, map[string]string{
		"id": "0",
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	validators, err := s.decodeValidatorsResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(validators)
	s.Require().Len(*validators, 1)
	s.Require().Equal(uint64(0), (*validators)[0].Index, "Validator index should be 0")
	s.Require().Equal("active_ongoing", (*validators)[0].Status, "Validator status should be active_ongoing")
	s.Require().True((*validators)[0].Balance > 0, "Validator balance should be positive")
	// 4e12 Gwei = 4 * 10^12 Gwei = 4,000,000,000,000 Gwei = 4000 BERA
	s.Require().True((*validators)[0].Balance <= 4e12, "Validator balance should not exceed 4000 BERA")
}

// TestGetValidatorsWithStatus tests querying validators with status parameter.
func (s *BeaconKitE2ESuite) TestGetValidatorsWithStatus() {
	resp, err := s.getValidator(utils.StateIDHead, map[string]string{
		"status": "active_ongoing",
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	validators, err := s.decodeValidatorsResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(validators)

	for _, validator := range *validators {
		s.Require().Equal("active_ongoing", validator.Status, "Validator status should be active_ongoing")
		s.Require().True(validator.Balance > 0, "Validator balance should be positive")
		// 4e12 Gwei = 4 * 10^12 Gwei = 4,000,000,000,000 Gwei = 4000 BERA
		s.Require().True(validator.Balance <= 4e12, "Validator balance should not exceed 4000 BERA")
	}
}

// TestGetValidatorsWithIDAndStatus tests querying validators with both ID and status parameters.
func (s *BeaconKitE2ESuite) TestGetValidatorsWithIDAndStatus() {
	// Call with both ID and status
	resp, err := s.getValidator(utils.StateIDHead, map[string]string{
		"id":     "0",
		"status": "active_ongoing",
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	// Decode and verify response
	validators, err := s.decodeValidatorsResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(validators)

	s.Require().Len(*validators, 1)
	// Verify the validator has the expected ID and status
	s.Require().Equal(uint64(0), (*validators)[0].Index, "Validator index should be 0")
	s.Require().Equal("active_ongoing", (*validators)[0].Status, "Validator status should be active_ongoing")
	s.Require().True((*validators)[0].Balance > 0, "Validator balance should be positive")
	// 4e12 Gwei = 4 * 10^12 Gwei = 4,000,000,000,000 Gwei = 4000 BERA
	s.Require().True((*validators)[0].Balance <= 4e12, "Validator balance should not exceed 4000 BERA")
}

// TestGetValidatorsWithStateGenesis tests querying validators with state genesis.
func (s *BeaconKitE2ESuite) TestGetValidatorsWithStateGenesis() {
	resp, err := s.getValidator(utils.StateIDGenesis)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	validators, err := s.decodeValidatorsResponse(resp)
	s.Require().NoError(err)
	s.Require().NotNil(validators)

	for _, validator := range *validators {
		s.Require().True(validator.Status == "active_ongoing")
		s.Require().True(validator.Balance > 0, "Validator balance should be positive")
		// 32e9 Gwei = 32 * 10^9 Gwei = 32,000,000,000 Gwei = 32 BERA
		s.Require().True(validator.Balance <= 32e9, "Validator balance should not exceed 32 BERA")
	}
}

// TestGenesis tests querying the genesis of the beacon node.
func (s *BeaconKitE2ESuite) TestGenesis() {
	client := s.initBeaconTest()

	genesisResp, err := client.Genesis(s.Ctx(), &beaconapi.GenesisOpts{})
	s.Require().NoError(err)
	s.Require().NotNil(genesisResp)
	s.Require().NotEmpty(genesisResp.Data)

	genesis := genesisResp.Data

	s.Require().NotZero(genesis.GenesisTime, "Genesis time should not be zero")
	s.Require().True(genesis.GenesisTime.Unix() < time.Now().Unix(), "Genesis time should be in the past")

	s.Require().NotEmpty(genesis.GenesisValidatorsRoot, "Genesis validators root should not be empty")
	s.Require().NotEmpty(genesis.GenesisForkVersion, "Genesis fork version should not be empty")

	expectedVersion := version.Electra() // TODO: change this back to Deneb once devnet spec is updated.
	s.Require().Equal(
		expectedVersion[:],
		genesis.GenesisForkVersion[:],
		"Genesis fork version does not match expected value",
	)
}

// TestConfigSpec tests querying the config spec.
func (s *BeaconKitE2ESuite) TestConfigSpec() {
	client := s.initBeaconTest()

	resp, err := client.Spec(s.Ctx(), &beaconapi.SpecOpts{})
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	specData := resp.Data
	s.Require().NotNil(specData)

	// TODO: make test use configurable chain spec.
	chainspec, err := spec.DevnetChainSpec()
	s.Require().NoError(err)

	// Verify essential config parameters exist and have expected types
	s.Require().Contains(specData, "DEPOSIT_CONTRACT_ADDRESS")
	addressBytes, ok := specData["DEPOSIT_CONTRACT_ADDRESS"].([]byte)
	s.Require().True(ok, "DEPOSIT_CONTRACT_ADDRESS should be a byte array")
	var executionAddress common.Address
	copy(executionAddress[:], addressBytes)

	// Convert chainspec's ExecutionAddress to common.Address
	depositAddr := common.Address(chainspec.DepositContractAddress())
	s.Require().Equal(depositAddr, executionAddress)

	s.Require().Contains(specData, "DEPOSIT_NETWORK_ID")
	s.Require().Equal(chainspec.DepositEth1ChainID(), specData["DEPOSIT_NETWORK_ID"])

	s.Require().Contains(specData, "DOMAIN_AGGREGATE_AND_PROOF")
	s.Require().EqualValues(chainspec.DomainTypeAggregateAndProof(), specData["DOMAIN_AGGREGATE_AND_PROOF"])

	// Check penalty quotients
	s.Require().Contains(specData, "INACTIVITY_PENALTY_QUOTIENT")
	s.Require().Zero(specData["INACTIVITY_PENALTY_QUOTIENT"])

	s.Require().Contains(specData, "INACTIVITY_PENALTY_QUOTIENT_ALTAIR")
	s.Require().Zero(specData["INACTIVITY_PENALTY_QUOTIENT_ALTAIR"])
}

// TestBeaconBlockHeaderByHead tests querying the beacon block header by head.
func (s *BeaconKitE2ESuite) TestBeaconBlockHeaderByHead() {
	client := s.initBeaconTest()

	resp, err := client.BeaconBlockHeader(s.Ctx(), &beaconapi.BeaconBlockHeaderOpts{
		Block: utils.StateIDHead,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	header := resp.Data
	s.Require().NotNil(header)

	s.Require().NotNil(header.Root, "Block root should not be nil")
	s.Require().False(header.Root.IsZero(), "Block root should not be zero")

	s.Require().NotNil(header.Header, "Header should not be nil")
	s.Require().NotNil(header.Header.Message, "Header message should not be nil")

	message := header.Header.Message
	s.Require().True(message.Slot > 0, "Slot should be greater than 0 for head")
	s.Require().NotNil(message.ParentRoot, "ParentRoot should not be nil")
	s.Require().False(message.ParentRoot.IsZero(), "ParentRoot should not be zero")
	s.Require().NotNil(message.StateRoot, "StateRoot should not be nil")
	s.Require().False(message.StateRoot.IsZero(), "StateRoot should not be zero")
	s.Require().NotNil(message.BodyRoot, "BodyRoot should not be nil")
	s.Require().False(message.BodyRoot.IsZero(), "BodyRoot should not be zero")
}

// TestBeaconBlockHeaderByGenesis tests querying the beacon block header by genesis.
func (s *BeaconKitE2ESuite) TestBeaconBlockHeaderByGenesis() {
	client := s.initBeaconTest()

	resp, err := client.BeaconBlockHeader(s.Ctx(), &beaconapi.BeaconBlockHeaderOpts{
		Block: utils.StateIDGenesis,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	header := resp.Data
	s.Require().NotNil(header)

	// Verify genesis header fields
	s.Require().NotNil(header.Root, "Genesis block root should not be nil")
	s.Require().False(header.Root.IsZero(), "Genesis block root should not be zero")

	s.Require().NotNil(header.Header, "Genesis header should not be nil")
	s.Require().NotNil(header.Header.Message, "Genesis header message should not be nil")

	message := header.Header.Message
	// Currently querying for genesis will return slot 1.
	s.Require().Equal(phase0.Slot(1), message.Slot, "Genesis slot should be 1")
	s.Require().NotNil(message.ProposerIndex, "Genesis proposer index should not be nil")
	s.Require().GreaterOrEqual(message.ProposerIndex, phase0.ValidatorIndex(0), "Genesis proposer index should be greater than 0")
	s.Require().NotNil(message.ParentRoot, "Genesis parent root should not be nil")
	s.Require().NotNil(message.StateRoot, "Genesis state root should not be nil")
	s.Require().NotNil(message.BodyRoot, "Genesis body root should not be nil")
}

// TestBeaconBlockHeaderBySlot tests querying the beacon block header by a specific slot.
func (s *BeaconKitE2ESuite) TestBeaconBlockHeaderBySlot() {
	client := s.initBeaconTest()

	resp, err := client.BeaconBlockHeader(s.Ctx(), &beaconapi.BeaconBlockHeaderOpts{
		Block: fmt.Sprintf("%d", phase0.Slot(10)),
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	header := resp.Data
	s.Require().NotNil(header)

	s.Require().NotNil(header.Root, "Block root should not be nil")
	s.Require().False(header.Root.IsZero(), "Block root should not be zero")

	s.Require().NotNil(header.Header, "Header should not be nil")
	s.Require().NotNil(header.Header.Message, "Header message should not be nil")

	message := header.Header.Message
	s.Require().Equal(phase0.Slot(10), message.Slot, "Slot should be 10")
	s.Require().NotNil(message.ParentRoot, "ParentRoot should not be nil")
	s.Require().False(message.ParentRoot.IsZero(), "ParentRoot should not be zero")
	s.Require().NotNil(message.StateRoot, "StateRoot should not be nil")
	s.Require().False(message.StateRoot.IsZero(), "StateRoot should not be zero")
	s.Require().NotNil(message.BodyRoot, "BodyRoot should not be nil")
	s.Require().False(message.BodyRoot.IsZero(), "BodyRoot should not be zero")
}
