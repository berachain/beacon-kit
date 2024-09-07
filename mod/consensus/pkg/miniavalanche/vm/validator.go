package vm

import (
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/encoding"
)

type Validator struct {
	NodeID ids.NodeID
	Weight uint64

	// note no nonce, not validatorTxID to distinguish restaked validators
	// TODO: mind this once dynamic validators set is introduced

	bytes []byte
	id    ids.ID
}

func NewValidator(nodeID ids.NodeID, weight uint64) (*Validator, error) {
	val := &Validator{
		NodeID: nodeID,
		Weight: weight,
	}
	return val, val.initValID()
}

func ParseValidator(valBytes []byte) (*Validator, error) {
	val := &Validator{}
	if err := encoding.Decode(valBytes, &val); err != nil {
		return nil, fmt.Errorf("unable to parse validator: %w", err)
	}

	return val, val.initValID()
}

func (v *Validator) initValID() error {
	bytes, err := encoding.Encode(v)
	if err != nil {
		return fmt.Errorf("failed encoding validator %v: %w", v, err)
	}
	v.bytes = bytes
	v.id = hashing.ComputeHash256Array(v.bytes)
	return nil
}
