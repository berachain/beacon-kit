package blocks

import enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

type ExecutionPayloadCapella struct{}

type FunBlock struct {
	Henlo            []byte `ssz-size:"32"`
	ExecutionPayload *enginev1.ExecutionPayloadCapella
}
