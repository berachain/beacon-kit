// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package bls

import "github.com/prysmaticlabs/prysm/v4/crypto/bls/blst"

func (privKey PrivKey) Sign(digestBz []byte) ([]byte, error) {
	secretKey, err := blst.SecretKeyFromBytes(privKey.Key)
	if err != nil {
		return nil, err
	}

	sig := secretKey.Sign(digestBz)
	return sig.Marshal(), nil
}

func (pubKey PubKey) VerifySignature(msg, sig []byte) bool {
	if len(sig) != 96 { //nolint:gomnd // bls this bls that...
		return false
	}
	if len(msg) != 32 { //nolint:gomnd // bls this bls that...
		return false
	}

	pubK, _ := blst.PublicKeyFromBytes(pubKey.Key)
	ok, err := blst.VerifySignature(sig, [32]byte(msg[:32]), pubK)
	if err != nil {
		return false
	}
	return ok
}
