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

// TODO: Deposit needs manual fastssz migration to handle dual interface compatibility (like BeaconBlockBody)
// go:generate sszgen -path . -objs Deposit -output deposit_sszgen.go -include ../../primitives/common,../../primitives/crypto,../../primitives/math,../../primitives/bytes

package types

import (
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// depositSize is the size of the SSZ encoding of a Deposit.
const depositSize = 192 // 48 + 32 + 8 + 96 + 8

// Compile-time assertions to ensure Deposit implements necessary interfaces.
var (
	_ ssz.StaticObject                    = (*Deposit)(nil)
	_ constraints.SSZMarshallableRootable = (*Deposit)(nil)
)

// Deposit into the consensus layer from the deposit contract in the execution
// layer.
type Deposit struct {
	// Public key of the validator specified in the deposit.
	Pubkey crypto.BLSPubkey `json:"pubkey"`
	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials WithdrawalCredentials `json:"credentials"`
	// Deposit amount in gwei.
	Amount math.Gwei `json:"amount"`
	// Signature of the deposit data.
	Signature crypto.BLSSignature `json:"signature"`
	// Index of the deposit in the deposit contract.
	Index uint64 `json:"index"`
}

func NewEmptyDeposit() *Deposit {
	return &Deposit{}
}

// Equals returns true if the Deposit is equal to the other.
func (d *Deposit) Equals(o *Deposit) bool {
	return d.Pubkey == o.Pubkey &&
		d.Credentials == o.Credentials &&
		d.Amount == o.Amount &&
		d.Signature == o.Signature &&
		d.Index == o.Index
}

// VerifySignature verifies the deposit data and signature.
func (d *Deposit) VerifySignature(
	forkData *ForkData,
	domainType common.DomainType,
	signatureVerificationFn func(
		pubkey crypto.BLSPubkey, message []byte, signature crypto.BLSSignature,
	) error,
) error {
	return (&DepositMessage{
		Pubkey:      d.Pubkey,
		Credentials: d.Credentials,
		Amount:      d.Amount,
	}).VerifyCreateValidator(
		forkData, d.Signature,
		domainType, signatureVerificationFn,
	)
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// DefineSSZ defines the SSZ encoding for the Deposit object.
func (d *Deposit) DefineSSZ(c *ssz.Codec) {
	ssz.DefineStaticBytes(c, &d.Pubkey)
	ssz.DefineStaticBytes(c, &d.Credentials)
	ssz.DefineUint64(c, &d.Amount)
	ssz.DefineStaticBytes(c, &d.Signature)
	ssz.DefineUint64(c, &d.Index)
}

// MarshalSSZ marshals the Deposit object to SSZ format.
func (d *Deposit) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(d))
	return buf, ssz.EncodeToBytes(buf, d)
}

func (*Deposit) ValidateAfterDecodingSSZ() error { return nil }

// SizeSSZ returns the SSZ encoded size of the Deposit object.
func (d *Deposit) SizeSSZ(*ssz.Sizer) uint32 {
	return depositSize
}

// HashTreeRoot computes the Merkleization of the Deposit object.
func (d *Deposit) HashTreeRoot() common.Root {
	return ssz.HashSequential(d)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith ssz hashes the Deposit object with a hasher.
func (d *Deposit) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Pubkey'
	hh.PutBytes(d.Pubkey[:])

	// Field (1) 'Credentials'
	hh.PutBytes(d.Credentials[:])

	// Field (2) 'Amount'
	hh.PutUint64(uint64(d.Amount))

	// Field (3) 'Signature'
	hh.PutBytes(d.Signature[:])

	// Field (4) 'Index'
	hh.PutUint64(d.Index)

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the Deposit object.
func (d *Deposit) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(d)
}

// MarshalSSZTo marshals the Deposit object into a pre-allocated byte slice.
func (d *Deposit) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := d.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the Deposit object.
func (d *Deposit) UnmarshalSSZ(buf []byte) error {
	// For now, delegate to karalabe/ssz for unmarshaling
	// This preserves compatibility during migration
	return ssz.DecodeFromBytes(buf, d)
}

// SizeSSZFastSSZ returns the ssz encoded size in bytes for the Deposit (fastssz).
// TODO: Rename to SizeSSZ() once karalabe/ssz is fully removed.
func (d *Deposit) SizeSSZFastSSZ() (size int) {
	return depositSize
}

/* -------------------------------------------------------------------------- */
/*                             Getters and Setters                            */
/* -------------------------------------------------------------------------- */

// GetAmount returns the deposit amount in gwei.
func (d *Deposit) GetAmount() math.Gwei {
	return d.Amount
}

// GetPubkey returns the public key of the validator specified in the deposit.
func (d *Deposit) GetPubkey() crypto.BLSPubkey {
	return d.Pubkey
}

// GetIndex returns the index of the deposit in the deposit contract.
func (d *Deposit) GetIndex() math.U64 {
	return math.U64(d.Index)
}

// GetSignature returns the signature of the deposit data.
func (d *Deposit) GetSignature() crypto.BLSSignature {
	return d.Signature
}

// GetWithdrawalCredentials returns the staking credentials of the deposit.
func (d *Deposit) GetWithdrawalCredentials() WithdrawalCredentials {
	return d.Credentials
}

// HasEth1WithdrawalCredentials as defined in the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/validator.md#eth1_address_withdrawal_prefix
func (d *Deposit) HasEth1WithdrawalCredentials() bool {
	return d.Credentials.IsValidEth1WithdrawalCredentials()
}
