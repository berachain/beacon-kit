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
// defined by:
// *v = UnusedType(buf[0])
func TestDecodeUnusedTypeEquality(t *testing.T) {
	t.Parallel()
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "decode-unused-type-empty", args: args{buf: []byte{0x00}}},
		{name: "decode-unused-type-one", args: args{buf: []byte{0x01}}},
		{name: "decode-unused-type-max", args: args{buf: []byte{0xff}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut1 := new(common.UnusedType)
			if err := ssz.DecodeFromBytes(tt.args.buf, ut1); err != nil {
				t.Errorf("DecodeFromBytes() error = %v", err)
			}
			ut2 := new(common.UnusedType)
			*ut2 = common.UnusedType(tt.args.buf[0])
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
