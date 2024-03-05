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

package enginev1

import (
	"github.com/itsdevbear/bolaris/config/version"
	"github.com/itsdevbear/bolaris/crypto/sha256"
	byteslib "github.com/itsdevbear/bolaris/lib/bytes"
	"github.com/itsdevbear/bolaris/math"
	fssz "github.com/prysmaticlabs/fastssz"
	"google.golang.org/protobuf/proto"
)

// Version returns the version identifier for the ExecutionPayloadDeneb.
func (p *ExecutionPayloadEnvelope) Version() int {
	switch p.GetPayload().(type) {
	case *ExecutionPayloadEnvelope_Deneb:
		return version.Deneb
	case *ExecutionPayloadEnvelope_DenebHeader:
		return version.Deneb
	default:
		return 0
	}
}

func (p *ExecutionPayloadEnvelope) IsEmpty() bool {
	if p.GetPayload() == nil {
		return true
	}

	return p.ToProto() == nil
}

// IsBlinded indicates whether the payload is blinded. For
// ExecutionPayloadDeneb,
// this is always false.
func (p *ExecutionPayloadEnvelope) IsBlinded() bool {
	switch p.GetPayload().(type) {
	case *ExecutionPayloadEnvelope_Deneb:
		return false
	case *ExecutionPayloadEnvelope_DenebHeader:
		return true
	}
	return false
}

// ToProto returns the ExecutionPayloadDeneb as a proto.Message.
func (p *ExecutionPayloadEnvelope) ToProto() proto.Message {
	return p.getPayload()
}

// GetValue returns the value of the payload.
func (p *ExecutionPayloadEnvelope) GetValue() math.Wei {
	if p.PayloadValue == nil {
		return math.Wei{}
	}

	// We have to convert big endian to little endian because the value is
	// coming
	// from the execution layer.
	return math.WeiFromBytes(
		byteslib.CopyAndReverseEndianess(p.GetPayloadValue()),
	)
}

// GetBlockHash retrieves the block hash from the payload.
func (p *ExecutionPayloadEnvelope) GetBlockHash() []byte {
	return p.getPayload().(interface{ GetBlockHash() []byte }).GetBlockHash()
}

// GetParentHash retrieves the parent hash from the payload.
func (p *ExecutionPayloadEnvelope) GetParentHash() []byte {
	return p.getPayload().(interface{ GetParentHash() []byte }).GetParentHash()
}

// GetTransactions retrieves the transactions from the payload.
func (p *ExecutionPayloadEnvelope) GetTransactions() [][]byte {
	if p.IsBlinded() {
		return [][]byte{}
	}
	return p.getPayload().(interface{ GetTransactions() [][]byte }).
		GetTransactions()
}

// GetTransactionsRoot retrieves the transactions root from the payload.
func (p *ExecutionPayloadEnvelope) GetTransactionsRoot() []byte {
	payload := p.getPayload()
	if p.IsBlinded() {
		return payload.(interface{ GetTransactionsRoot() []byte }).
			GetTransactionsRoot()
	}
	return sha256.HashRootAndMixinLengthAsBzSlice(
		payload.(interface{ GetTransactions() [][]byte }).
			GetTransactions())
}

// GetWithdrawals retrieves the withdrawals from the payload.
func (p *ExecutionPayloadEnvelope) GetWithdrawals() []*Withdrawal {
	if p.IsBlinded() {
		return []*Withdrawal{}
	}
	return p.getPayload().(interface{ GetWithdrawals() []*Withdrawal }).
		GetWithdrawals()
}

// GetWithdrawalsRoot retrieves the withdrawals root from the payload.
func (p *ExecutionPayloadEnvelope) GetWithdrawalsRoot() []byte {
	payload := p.getPayload()
	if p.IsBlinded() {
		return payload.(interface{ GetWithdrawalsRoot() []byte }).
			GetWithdrawalsRoot()
	}
	return sha256.HashRootAndMixinLengthAsSlice[*Withdrawal](
		payload.(interface{ GetWithdrawals() []*Withdrawal }).GetWithdrawals())
}

// HashTreeRoot calculates the hash tree root of the payload.
func (p *ExecutionPayloadEnvelope) HashTreeRoot() ([32]byte, error) {
	return p.getPayload().(interface{ HashTreeRoot() ([32]byte, error) }).
		HashTreeRoot()
}

// HashTreeRootWith calculates the hash tree root of the payload using a
// provided hasher.
func (p *ExecutionPayloadEnvelope) HashTreeRootWith(h *fssz.Hasher) error {
	return p.getPayload().(interface {
		HashTreeRootWith(h *fssz.Hasher) error
	}).HashTreeRootWith(h)
}

// MarshalSSZ marshals the payload into SSZ-encoded bytes.
func (p *ExecutionPayloadEnvelope) MarshalSSZ() ([]byte, error) {
	return p.getPayload().(interface{ MarshalSSZ() ([]byte, error) }).MarshalSSZ()
}

// UnmarshalSSZ unmarshals SSZ-encoded bytes into the payload.
func (p *ExecutionPayloadEnvelope) UnmarshalSSZ(bz []byte) error {
	return p.getPayload().(interface{ UnmarshalSSZ([]byte) error }).UnmarshalSSZ(
		bz,
	)
}

// MarshalSSZTo marshals the payload into SSZ-encoded bytes and appends it to
// the
// provided byte slice.
func (p *ExecutionPayloadEnvelope) MarshalSSZTo(bz []byte) ([]byte, error) {
	return p.getPayload().(interface{ MarshalSSZTo([]byte) ([]byte, error) }).
		MarshalSSZTo(bz)
}

// SizeSSZ returns the size of the SSZ-encoded payload in bytes.
func (p *ExecutionPayloadEnvelope) SizeSSZ() int {
	return p.getPayload().(interface{ SizeSSZ() int }).SizeSSZ()
}

// GetPayload returns the payload of the ExecutionPayloadEnvelope as an
// interface{}.
// The actual type of the returned value depends on which payload is set.
// The caller will need to type assert the returned value to use it.
func (p *ExecutionPayloadEnvelope) getPayload() proto.Message {
	switch {
	case p.GetDeneb() != nil:
		return p.GetDeneb()
	case p.GetDenebHeader() != nil:
		return p.GetDenebHeader()
	default:
		return nil
	}
}
