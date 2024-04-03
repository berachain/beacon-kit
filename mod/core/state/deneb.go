// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package state

import (
	"errors"
	"sort"

	"github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// DefaultBeaconStateDeneb returns a default BeaconStateDeneb.
//
// TODO: take in BeaconConfig params to determine the
// default length of the arrays, which we are currently
// and INCORRECTLY setting to 0.
func DefaultBeaconStateDeneb() *BeaconStateDeneb {
	//nolint:gomnd // default allocs.
	return &BeaconStateDeneb{
		GenesisValidatorsRoot: primitives.Root{},

		Slot: 0,
		LatestBlockHeader: &types.BeaconBlockHeader{
			Slot:          0,
			ProposerIndex: 0,
			ParentRoot:    primitives.Root{},
			StateRoot:     primitives.Root{},
			BodyRoot:      primitives.Root{},
		},
		BlockRoots: make([][32]byte, 1),
		StateRoots: make([][32]byte, 1),
		Eth1BlockHash: common.HexToHash(
			"0xa63c365d92faa4de2a64a80ed4759c3e9dfa939065c10af08d2d8d017a29f5f4",
		),
		Eth1DepositIndex: 0,
		Validators:       make([]*types.Validator, 0),
		Balances:         make([]uint64, 0),
		RandaoMixes:      make([][32]byte, 8),
		Slashings:        make([]uint64, 1),
		TotalSlashing:    0,
	}
}

// TODO: should we replace ? in ssz-size with values to ensure we are hash tree
// rooting correctly?
//
//go:generate go run github.com/fjl/gencodec -type BeaconStateDeneb -field-override beaconStateDenebJSONMarshaling -out deneb.json.go
//nolint:lll // various json tags.
type BeaconStateDeneb struct {
	// Versioning
	//
	//nolint:lll
	GenesisValidatorsRoot primitives.Root `json:"genesisValidatorsRoot" ssz-size:"32"`
	Slot                  primitives.Slot `json:"slot"`

	// History
	LatestBlockHeader *types.BeaconBlockHeader `json:"latestBlockHeader"`
	BlockRoots        [][32]byte               `json:"blockRoots"        ssz-size:"?,32" ssz-max:"8192"`
	StateRoots        [][32]byte               `json:"stateRoots"        ssz-size:"?,32" ssz-max:"8192"`

	// Eth1
	Eth1BlockHash    primitives.ExecutionHash `json:"eth1BlockHash"    ssz-size:"32"`
	Eth1DepositIndex uint64                   `json:"eth1DepositIndex"`

	// Registry
	Validators []*types.Validator `json:"validators" ssz-max:"1099511627776"`
	Balances   []uint64           `json:"balances"   ssz-max:"1099511627776"`

	// Randomness
	RandaoMixes [][32]byte `json:"randaoMixes" ssz-size:"?,32" ssz-max:"65536"`

	// Withdrawals
	NextWithdrawalIndex          uint64 `json:"nextWithdrawalIndex"`
	NextWithdrawalValidatorIndex uint64 `json:"nextWithdrawalValidatorIndex"`

	// Slashing
	Slashings     []uint64 `json:"slashings"     ssz-max:"1099511627776"`
	TotalSlashing uint64   `json:"totalSlashing"`
}

// Copy returns a deep copy of BeaconStateDeneb.
func (s *BeaconStateDeneb) Copy() BeaconState {
	return &BeaconStateDeneb{
		GenesisValidatorsRoot: s.GenesisValidatorsRoot,
		Slot:                  s.Slot,
		LatestBlockHeader:     s.LatestBlockHeader.Copy(),
		BlockRoots:            append([][32]byte{}, s.BlockRoots...),
		StateRoots:            append([][32]byte{}, s.StateRoots...),
		Eth1BlockHash:         s.Eth1BlockHash,
		Eth1DepositIndex:      s.Eth1DepositIndex,
		Validators:            append([]*types.Validator{}, s.Validators...),
		Balances:              append([]uint64{}, s.Balances...),
		RandaoMixes:           append([][32]byte{}, s.RandaoMixes...),
		Slashings:             append([]uint64{}, s.Slashings...),
		TotalSlashing:         s.TotalSlashing,
	}
}

func (s *BeaconStateDeneb) Save() {
}

// String returns a string representation of BeaconStateDeneb.
func (s *BeaconStateDeneb) String() string {
	return "TODO: BeaconStateDeneb"
}

// IncreaseBalance increases the balance of a validator.
func (s *BeaconStateDeneb) IncreaseBalance(
	idx primitives.ValidatorIndex,
	delta primitives.Gwei,
) error {
	s.Balances[idx] += uint64(delta)
	return nil
}

// DecreaseBalance decreases the balance of a validator.
func (s *BeaconStateDeneb) DecreaseBalance(
	idx primitives.ValidatorIndex,
	delta primitives.Gwei,
) error {
	s.Balances[idx] -= uint64(delta)
	return nil
}

func (s *BeaconStateDeneb) GetTotalActiveBalances(
	slotsPerEpoch uint64,
) (primitives.Gwei, error) {
	epoch, err := s.GetCurrentEpoch(slotsPerEpoch)
	if err != nil {
		return 0, err
	}

	totalActiveBalances := primitives.Gwei(0)
	for _, v := range s.Validators {
		if v.IsActive(epoch) {
			totalActiveBalances += v.EffectiveBalance
		}
	}
	return totalActiveBalances, nil
}

func (s *BeaconStateDeneb) GetCurrentEpoch(
	slotsPerEpoch uint64,
) (primitives.Epoch, error) {
	return primitives.Epoch(uint64(s.Slot) / slotsPerEpoch), nil
}

// UpdateBlockRootAtIndex sets a block root in the BeaconStore.
func (s *BeaconStateDeneb) UpdateBlockRootAtIndex(
	index uint64,
	root primitives.Root,
) error {
	s.BlockRoots[index] = root
	return nil
}

// GetBlockRoot retrieves the block root from the BeaconStore.
func (s *BeaconStateDeneb) GetBlockRootAtIndex(
	index uint64,
) (primitives.Root, error) {
	return s.BlockRoots[index], nil
}

// SetLatestBlockHeader sets the latest block header in the BeaconStore.
func (s *BeaconStateDeneb) SetLatestBlockHeader(
	header *types.BeaconBlockHeader,
) error {
	s.LatestBlockHeader = header
	return nil
}

// GetLatestBlockHeader retrieves the latest block header from the BeaconStore.
func (s *BeaconStateDeneb) GetLatestBlockHeader() (*types.BeaconBlockHeader, error) {
	return s.LatestBlockHeader, nil
}

// UpdateEth1BlockHash sets the Eth1 hash in the BeaconStore.
func (s *BeaconStateDeneb) UpdateEth1BlockHash(
	hash primitives.ExecutionHash,
) error {
	s.Eth1BlockHash = hash
	return nil
}

// GetEth1BlockHash retrieves the Eth1 hash from the BeaconStore.
func (s *BeaconStateDeneb) GetEth1BlockHash() (primitives.ExecutionHash, error) {
	return s.Eth1BlockHash, nil
}

func (s *BeaconStateDeneb) UpdateRandaoMixAtIndex(
	index uint64,
	mix primitives.Bytes32,
) error {
	s.RandaoMixes[index] = mix
	return nil
}

func (s *BeaconStateDeneb) GetRandaoMixAtIndex(
	index uint64,
) (primitives.Bytes32, error) {
	return s.RandaoMixes[index], nil
}

// UpdateStateRootAtIndex sets the state root at the given slot.
func (s *BeaconStateDeneb) UpdateStateRootAtIndex(
	idx uint64,
	stateRoot primitives.Root,
) error {
	s.StateRoots[idx] = stateRoot
	return nil
}

// StateRootAtIndex returns the state root at the given slot.
func (s *BeaconStateDeneb) StateRootAtIndex(
	idx uint64,
) (primitives.Root, error) {
	return s.StateRoots[idx], nil
}

func (s *BeaconStateDeneb) UpdateSlashingAtIndex(
	index uint64,
	amount primitives.Gwei,
) error {
	total := s.TotalSlashing
	oldValue := s.Slashings[index]
	s.Slashings[index] += total - (oldValue) + uint64(amount)
	return nil
}

func (s *BeaconStateDeneb) GetSlashingAtIndex(
	index uint64,
) (primitives.Gwei, error) {
	return primitives.Gwei(s.Slashings[index]), nil
}

func (s *BeaconStateDeneb) GetTotalSlashing() (primitives.Gwei, error) {
	return primitives.Gwei(s.TotalSlashing), nil
}

func (s *BeaconStateDeneb) SetTotalSlashing(amount primitives.Gwei) error {
	s.TotalSlashing = uint64(amount)
	return nil
}

func (s *BeaconStateDeneb) GetSlot() (primitives.Slot, error) {
	return s.Slot, nil
}

func (s *BeaconStateDeneb) SetSlot(slot primitives.Slot) error {
	s.Slot = slot
	return nil
}

func (s *BeaconStateDeneb) SetGenesisValidatorsRoot(
	root primitives.Root,
) error {
	s.GenesisValidatorsRoot = root
	return nil
}

func (s *BeaconStateDeneb) GetGenesisValidatorsRoot() (primitives.Root, error) {
	return s.GenesisValidatorsRoot, nil
}

func (s *BeaconStateDeneb) AddValidator(val *types.Validator) error {
	// Ensure the validator does not already exist
	for _, v := range s.Validators {
		if v.Pubkey == val.Pubkey {
			return errors.New("validator already exists")
		}
	}

	// Append the new validator to the list of validators
	s.Validators = append(s.Validators, val)

	// Update the total balance with the validator's balance
	s.Balances = append(s.Balances, uint64(val.EffectiveBalance))
	return nil
}

func (s *BeaconStateDeneb) ValidatorByIndex(
	index primitives.ValidatorIndex,
) (*types.Validator, error) {
	return s.Validators[index], nil
}

func (s *BeaconStateDeneb) UpdateValidatorAtIndex(
	index primitives.ValidatorIndex,
	val *types.Validator,
) error {
	s.Validators[index] = val
	return nil
}

func (s *BeaconStateDeneb) RemoveValidatorAtIndex(
	idx primitives.ValidatorIndex,
) error {
	return nil
}

func (s *BeaconStateDeneb) ValidatorIndexByPubkey(
	pubkey []byte,
) (primitives.ValidatorIndex, error) {
	for i, v := range s.Validators {
		if v.Pubkey == [48]byte(pubkey) {
			return primitives.ValidatorIndex(i), nil
		}
	}
	return 0, errors.New("validator not found")
}

func (s *BeaconStateDeneb) GetValidators() ([]*types.Validator, error) {
	return s.Validators, nil
}

func (s *BeaconStateDeneb) GetValidatorsByEffectiveBalance() ([]*types.Validator, error) {
	validatorsCopy := make([]*types.Validator, len(s.Validators))
	copy(validatorsCopy, s.Validators)
	sort.Slice(validatorsCopy, func(i, j int) bool {
		return validatorsCopy[i].EffectiveBalance > validatorsCopy[j].EffectiveBalance
	})
	return validatorsCopy, nil
}

// beaconStateDenebJSONMarshaling is a type used to marshal/unmarshal
// BeaconStateDeneb.
type beaconStateDenebJSONMarshaling struct {
	GenesisValidatorsRoot hexutil.Bytes
	BlockRoots            []primitives.Root
	StateRoots            []primitives.Root
	RandaoMixes           []primitives.Root
}
