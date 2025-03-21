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

package common_test

import (
	"reflect"
	"testing"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/karalabe/ssz"
)

// Verify that DecodeFromBytes produces the same UnusedType obj as the previous implementation
// defined by *v = UnusedType(buf[0])
func TestDecodeUnusedTypeEquality(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		buf     []byte
		wantErr bool
	}{
		{name: "decode-unused-type-empty", buf: []byte{0x00}, wantErr: false},
		{name: "decode-unused-type-one", buf: []byte{0x01}, wantErr: false},
		{name: "decode-unused-type-max", buf: []byte{0xff}, wantErr: false},
		{name: "decode-unused-type-too-long", buf: []byte{0xff, 0xff}, wantErr: true},
		{name: "decode-unused-type-too-short", buf: []byte{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(common.UnusedType)
			if err := ssz.DecodeFromBytes(tt.buf, got); err != nil {
				if tt.wantErr {
					return
				}
				t.Errorf("DecodeFromBytes() error = %v", err)
			}
			want := common.UnusedType(tt.buf[0])
			if !reflect.DeepEqual(got, &want) {
				t.Errorf("MarshalSSZ() got = %v, want %v", got, want)
			}
		})
	}
}

// Verify that MarshalSSZ produces the same bytes as the previous implementation
// defined by:
// []byte{uint8(*ut)}
func TestEncodeUnusedTypeEquality(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		ut   common.UnusedType
	}{
		{name: "encode-unused-type-empty", ut: common.UnusedType(0)},
		{name: "encode-unused-type-one", ut: common.UnusedType(1)},
		{name: "encode-unused-type-max", ut: ^common.UnusedType(0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ut.MarshalSSZ()
			if err != nil {
				t.Errorf("MarshalSSZ() error = %v", err)
				return
			}
			want := []byte{uint8(tt.ut)}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("MarshalSSZ() got = %v, want %v", got, want)
			}
		})
	}
}
