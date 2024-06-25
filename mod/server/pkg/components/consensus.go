package components

import (
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/consensus"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/comet"
)

type ConsensusVersion uint8

const (
	CometBFTConsensus ConsensusVersion = iota
	RollKitConsensus
)

type ConsensusInput[T transaction.Tx, ValidatorUpdateT any] struct {
	depinject.In

	Version ConsensusVersion
	TxCodec transaction.Codec[T]
}

func ProvideConsensus[T transaction.Tx, ValidatorUpdateT any](
	in ConsensusInput[T, ValidatorUpdateT],
) consensus.Consensus[T, ValidatorUpdateT] {
	switch in.Version {
	case CometBFTConsensus:
		return comet.NewConsensus[T](
			in.TxCodec,
		)
	case RollKitConsensus:
		panic("not implemented")
	default:
		panic("unknown consensus version")
	}
}
