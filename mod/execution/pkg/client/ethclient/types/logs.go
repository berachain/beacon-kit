package types

import (
	"errors"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// FilterArgs is a map of string to any, used to configure a filter query.
type FilterArgs map[string]any

// New creates a new FilterArgs with the given addresses and topics.
func (f FilterArgs) New(
	addrs []common.ExecutionAddress, topics []common.ExecutionHash,
) FilterArgs {
	return FilterArgs{
		"address": addrs,
		"topics":  topics,
	}
}

// SetFromBlock sets the "fromBlock" field in the FilterArgs, ensuring that
// BlockHash is not also set.
func (f FilterArgs) SetFromBlock(block BlockNumber) error {
	var err error
	if f["fromBlock"], err = block.MarshalText(); err != nil {
		return err
	}
	if f["blockHash"] != nil {
		return errors.New("cannot specify both BlockHash and FromBlock/ToBlock")
	}
	return nil
}

// SetToBlock sets the "toBlock" field in the FilterArgs, ensuring that
// BlockHash is not also set.
func (f FilterArgs) SetToBlock(block BlockNumber) error {
	var err error
	if f["toBlock"], err = block.MarshalText(); err != nil {
		return err
	}
	if f["blockHash"] != nil {
		return errors.New("cannot specify both BlockHash and FromBlock/ToBlock")
	}
	return nil
}

// SetBlockHash sets the "blockHash" field in the FilterArgs, ensuring that
// FromBlock and ToBlock are not also set.
func (f FilterArgs) SetBlockHash(hash common.ExecutionHash) error {
	if f["fromBlock"] != nil || f["toBlock"] != nil {
		return errors.New("cannot specify both BlockHash and FromBlock/ToBlock")
	}
	f["blockHash"] = hash.Hex()
	return nil
}
