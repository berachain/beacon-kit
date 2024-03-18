package e2e_test

import (
	stakingabi "github.com/berachain/beacon-kit/contracts/abi"
	byteslib "github.com/berachain/beacon-kit/lib/bytes"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TestForgeScriptExecution tests the execution of a forge script
// against the beacon-kit network.
func (s *BeaconKitE2ESuite) TestDepositContract() {
	client := s.ConsensusClients()["cl-validator-beaconkit-0"]

	pubkey, err := client.GetPubKey(s.Ctx())
	s.Require().NoError(err)

	_, err = client.GetConsensusPower(s.Ctx())
	s.Require().NoError(err)

	dc, err := stakingabi.NewBeaconDepositContract(
		common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"),
		s.JSONRPCBalancer(),
	)
	s.Require().NoError(err)

	bz := byteslib.ToBytes32(s.GenesisAccount().Address().Bytes())
	bz[0] = 0x01

	tx, err := dc.Deposit(&bind.TransactOpts{
		From: s.GenesisAccount().Address(),
	}, pubkey, bz[:], 32e09, nil)
	s.Require().NoError(err)

	chainID, err := s.JSONRPCBalancer().ChainID(s.Ctx())
	s.Require().NoError(err)
	tx, err = s.GenesisAccount().SignTx(chainID, tx)
	s.Require().NoError(err)

	err = s.JSONRPCBalancer().SendTransaction(s.Ctx(), tx)
	s.Require().NoError(err)

	var receipt *types.Receipt
	receipt, err = bind.WaitMined(s.Ctx(), s.JSONRPCBalancer(), tx)
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), receipt.Status)
}
