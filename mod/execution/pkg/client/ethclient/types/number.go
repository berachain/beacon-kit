package types

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	SafeBlockNumbeStr      = "safe"
	FinalizedBlockNumStr   = "finalized"
	LatestBlockNumStr      = "latest"
	PendingBlockNumStr     = "pending"
	EarliestBlockNumberStr = "earliest"
)

const (
	SafeBlockNumber      = BlockNumber(-4)
	FinalizedBlockNumber = BlockNumber(-3)
	LatestBlockNumber    = BlockNumber(-2)
	PendingBlockNumber   = BlockNumber(-1)
	EarliestBlockNumber  = BlockNumber(0)
)

type BlockNumber int64

// UnmarshalJSON parses the given JSON fragment into a BlockNumber. It supports:
// - "safe", "finalized", "latest", "earliest" or "pending" as string arguments
// - the block number
// Returned errors:
// - an invalid block number error when the given argument isn't a known strings
// - an out of range error when the given block number is either too little or too large
func (bn *BlockNumber) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}

	switch input {
	case EarliestBlockNumberStr:
		*bn = EarliestBlockNumber
		return nil
	case LatestBlockNumStr:
		*bn = LatestBlockNumber
		return nil
	case PendingBlockNumStr:
		*bn = PendingBlockNumber
		return nil
	case FinalizedBlockNumStr:
		*bn = FinalizedBlockNumber
		return nil
	case SafeBlockNumbeStr:
		*bn = SafeBlockNumber
		return nil
	}

	blckNum, err := hexutil.DecodeUint64(input)
	if err != nil {
		return err
	}
	if blckNum > math.MaxInt64 {
		return errors.New("block number larger than int64")
	}
	*bn = BlockNumber(blckNum)
	return nil
}

// MarshalText implements encoding.TextMarshaler. It marshals:
// - "safe", "finalized", "latest", "earliest" or "pending" as strings
// - other numbers as hex
func (bn BlockNumber) MarshalText() ([]byte, error) {
	return []byte(bn.String()), nil
}

func (bn BlockNumber) String() string {
	switch bn {
	case EarliestBlockNumber:
		return EarliestBlockNumberStr
	case LatestBlockNumber:
		return LatestBlockNumStr
	case PendingBlockNumber:
		return PendingBlockNumStr
	case FinalizedBlockNumber:
		return FinalizedBlockNumStr
	case SafeBlockNumber:
		return SafeBlockNumbeStr
	default:
		if bn < 0 {
			return fmt.Sprintf("<invalid %d>", bn)
		}
		return hexutil.Uint64(bn).String()
	}
}
