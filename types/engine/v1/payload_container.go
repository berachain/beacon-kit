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
	"github.com/holiman/uint256"
	"github.com/itsdevbear/bolaris/crypto/sha256"
	byteslib "github.com/itsdevbear/bolaris/lib/bytes"
	"github.com/itsdevbear/bolaris/math"
	"github.com/itsdevbear/bolaris/types/consensus/version"
	fssz "github.com/prysmaticlabs/fastssz"
	"google.golang.org/protobuf/proto"
)

// Version returns the version identifier for the ExecutionPayloadDeneb.
func (p *ExecutionPayloadContainer) Version() int {
	switch p.GetPayload().(type) {
	case *ExecutionPayloadContainer_Deneb:
		return version.Deneb
	case *ExecutionPayloadContainer_DenebHeader:
		return version.Deneb
	case *ExecutionPayloadContainer_Capella:
		return version.Capella
	case *ExecutionPayloadContainer_CapellaHeader:
		return version.Capella
	default:
		return 0
	}
}

// IsBlinded indicates whether the payload is blinded. For ExecutionPayloadDeneb,
// this is always false.
func (p *ExecutionPayloadContainer) IsBlinded() bool {
	switch p.GetPayload().(type) {
	case *ExecutionPayloadContainer_Deneb:
		return false
	case *ExecutionPayloadContainer_DenebHeader:
		return true
	case *ExecutionPayloadContainer_Capella:
		return false
	case *ExecutionPayloadContainer_CapellaHeader:
		return true
	}
	return false
}

// ToProto returns the ExecutionPayloadDeneb as a proto.Message.
func (p *ExecutionPayloadContainer) ToProto() proto.Message {
	return p.getPayload()
}

// GetValue returns the value of the payload.
func (p *ExecutionPayloadContainer) GetValue() math.Wei {
	if p.PayloadValue == nil {
		return math.ZeroWei()
	}

	return uint256.NewInt(0).SetBytes(
		byteslib.CopyAndReverseEndianess(p.GetPayloadValue()))
}

// GetBlockHash retrieves the block hash from the payload.
func (p *ExecutionPayloadContainer) GetBlockHash() []byte {
	return p.getPayload().(interface{ GetBlockHash() []byte }).GetBlockHash()
}

// GetParentHash retrieves the parent hash from the payload.
func (p *ExecutionPayloadContainer) GetParentHash() []byte {
	return p.getPayload().(interface{ GetParentHash() []byte }).GetParentHash()
}

// GetTransactions retrieves the transactions from the payload.
func (p *ExecutionPayloadContainer) GetTransactions() [][]byte {
	if p.IsBlinded() {
		return [][]byte{}
	}
	return p.getPayload().(interface{ GetTransactions() [][]byte }).GetTransactions()
}

// GetTransactionsRoot retrieves the transactions root from the payload.
func (p *ExecutionPayloadContainer) GetTransactionsRoot() []byte {
	payload := p.getPayload()
	if p.IsBlinded() {
		return payload.(interface{ GetTransactionsRoot() []byte }).GetTransactionsRoot()
	}
	return sha256.HashRootAndMixinLengthAsBzSlice(
		payload.(interface{ GetTransactions() [][]byte }).GetTransactions())
}

// GetWithdrawals retrieves the withdrawals from the payload.
func (p *ExecutionPayloadContainer) GetWithdrawals() []*Withdrawal {
	if p.IsBlinded() {
		return []*Withdrawal{}
	}
	return p.getPayload().(interface{ GetWithdrawals() []*Withdrawal }).GetWithdrawals()
}

// GetWithdrawalsRoot retrieves the withdrawals root from the payload.
func (p *ExecutionPayloadContainer) GetWithdrawalsRoot() []byte {
	payload := p.getPayload()
	if p.IsBlinded() {
		return payload.(interface{ GetWithdrawalsRoot() []byte }).GetWithdrawalsRoot()
	}
	return sha256.HashRootAndMixinLengthAsSlice[*Withdrawal](
		payload.(interface{ GetWithdrawals() []*Withdrawal }).GetWithdrawals())
}

// HashTreeRoot calculates the hash tree root of the payload.
func (p *ExecutionPayloadContainer) HashTreeRoot() ([32]byte, error) {
	return p.getPayload().(interface{ HashTreeRoot() ([32]byte, error) }).HashTreeRoot()
}

// HashTreeRootWith calculates the hash tree root of the payload using a provided hasher.
func (p *ExecutionPayloadContainer) HashTreeRootWith(h *fssz.Hasher) error {
	return p.getPayload().(interface {
		HashTreeRootWith(h *fssz.Hasher) error
	}).HashTreeRootWith(h)
}

// MarshalSSZ marshals the payload into SSZ-encoded bytes.
func (p *ExecutionPayloadContainer) MarshalSSZ() ([]byte, error) {
	return p.getPayload().(interface{ MarshalSSZ() ([]byte, error) }).MarshalSSZ()
}

// UnmarshalSSZ unmarshals SSZ-encoded bytes into the payload.
func (p *ExecutionPayloadContainer) UnmarshalSSZ(bz []byte) error {
	return p.getPayload().(interface{ UnmarshalSSZ([]byte) error }).UnmarshalSSZ(bz)
}

// MarshalSSZTo marshals the payload into SSZ-encoded bytes and appends it to the
// provided byte slice.
func (p *ExecutionPayloadContainer) MarshalSSZTo(bz []byte) ([]byte, error) {
	return p.getPayload().(interface{ MarshalSSZTo([]byte) ([]byte, error) }).MarshalSSZTo(bz)
}

// SizeSSZ returns the size of the SSZ-encoded payload in bytes.
func (p *ExecutionPayloadContainer) SizeSSZ() int {
	return p.getPayload().(interface{ SizeSSZ() int }).SizeSSZ()
}

// GetPayload returns the payload of the ExecutionPayloadContainer as an interface{}.
// The actual type of the returned value depends on which payload is set.
// The caller will need to type assert the returned value to use it.
func (p *ExecutionPayloadContainer) getPayload() proto.Message {
	switch {
	case p.GetCapella() != nil:
		return p.GetCapella()
	case p.GetCapellaHeader() != nil:
		return p.GetCapellaHeader()
	case p.GetDeneb() != nil:
		return p.GetDeneb()
	case p.GetDenebHeader() != nil:
		return p.GetDenebHeader()
	default:
		return nil
	}
}
