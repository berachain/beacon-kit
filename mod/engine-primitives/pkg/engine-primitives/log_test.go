// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package engineprimitives_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//nolint:lll // test input data.
var unmarshalLogTests = map[string]struct {
	input     string
	want      *engineprimitives.Log
	wantError error
}{
	"ok": {
		input: `{"address":"0xecf8f87f810ecf450940c9f60066b4a7a501d6a7","blockHash":"0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056","blockNumber":"0x1ecfa4","data":"0x000000000000000000000000000000000000000000000001a055690d9db80000","logIndex":"0x2","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615"],"transactionHash":"0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e","transactionIndex":"0x3"}`,
		want: &engineprimitives.Log{
			Address: common.NewExecutionAddressFromHex(
				"0xecf8f87f810ecf450940c9f60066b4a7a501d6a7",
			),
			BlockHash: common.NewExecutionHashFromHex(
				"0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056",
			),
			BlockNumber: 2019236,
			Data: hexutil.MustDecode(
				"0x000000000000000000000000000000000000000000000001a055690d9db80000",
			),
			Index:   2,
			TxIndex: 3,
			TxHash: common.NewExecutionHashFromHex(
				"0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e",
			),
			Topics: []common.ExecutionHash{
				common.NewExecutionHashFromHex(
					"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
				),
				common.NewExecutionHashFromHex(
					"0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615",
				),
			},
		},
	},
	"empty data": {
		input: `{"address":"0xecf8f87f810ecf450940c9f60066b4a7a501d6a7","blockHash":"0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056","blockNumber":"0x1ecfa4","data":"0x","logIndex":"0x2","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615"],"transactionHash":"0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e","transactionIndex":"0x3"}`,
		want: &engineprimitives.Log{
			Address: common.NewExecutionAddressFromHex(
				"0xecf8f87f810ecf450940c9f60066b4a7a501d6a7",
			),
			BlockHash: common.NewExecutionHashFromHex(
				"0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056",
			),
			BlockNumber: 2019236,
			Data:        []byte{},
			Index:       2,
			TxIndex:     3,
			TxHash: common.NewExecutionHashFromHex(
				"0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e",
			),
			Topics: []common.ExecutionHash{
				common.NewExecutionHashFromHex(
					"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
				),
				common.NewExecutionHashFromHex(
					"0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615",
				),
			},
		},
	},
	"missing block fields (pending logs)": {
		input: `{"address":"0xecf8f87f810ecf450940c9f60066b4a7a501d6a7","data":"0x","logIndex":"0x0","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"],"transactionHash":"0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e","transactionIndex":"0x3"}`,
		want: &engineprimitives.Log{
			Address: common.NewExecutionAddressFromHex(
				"0xecf8f87f810ecf450940c9f60066b4a7a501d6a7",
			),
			BlockHash:   common.ExecutionHash{},
			BlockNumber: 0,
			Data:        []byte{},
			Index:       0,
			TxIndex:     3,
			TxHash: common.NewExecutionHashFromHex(
				"0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e",
			),
			Topics: []common.ExecutionHash{
				common.NewExecutionHashFromHex(
					"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
				),
			},
		},
	},
	"Removed: true": {
		input: `{"address":"0xecf8f87f810ecf450940c9f60066b4a7a501d6a7","blockHash":"0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056","blockNumber":"0x1ecfa4","data":"0x","logIndex":"0x2","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"],"transactionHash":"0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e","transactionIndex":"0x3","removed":true}`,
		want: &engineprimitives.Log{
			Address: common.NewExecutionAddressFromHex(
				"0xecf8f87f810ecf450940c9f60066b4a7a501d6a7",
			),
			BlockHash: common.NewExecutionHashFromHex(
				"0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056",
			),
			BlockNumber: 2019236,
			Data:        []byte{},
			Index:       2,
			TxIndex:     3,
			TxHash: common.NewExecutionHashFromHex(
				"0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e",
			),
			Topics: []common.ExecutionHash{
				common.NewExecutionHashFromHex(
					"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
				),
			},
			Removed: true,
		},
	},
	"missing data": {
		input:     `{"address":"0xecf8f87f810ecf450940c9f60066b4a7a501d6a7","blockHash":"0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056","blockNumber":"0x1ecfa4","logIndex":"0x2","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615","0x000000000000000000000000f9dff387dcb5cc4cca5b91adb07a95f54e9f1bb6"],"transactionHash":"0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e","transactionIndex":"0x3"}`,
		wantError: errors.New("missing required field 'data' for Log"),
	},
}

func TestUnmarshalLog(t *testing.T) {
	dumper := spew.ConfigState{DisableMethods: true, Indent: "    "}
	for name, test := range unmarshalLogTests {
		var log *engineprimitives.Log
		err := json.Unmarshal([]byte(test.input), &log)
		checkError(t, name, err, test.wantError)
		if test.wantError == nil && err == nil {
			if !reflect.DeepEqual(log, test.want) {
				t.Errorf(
					"test %q:\nGOT %sWANT %s", name,
					dumper.Sdump(log), dumper.Sdump(test.want),
				)
			}
		}
	}
}

func checkError(t *testing.T, testname string, got, want error) bool {
	t.Helper()
	if got == nil {
		if want != nil {
			t.Errorf("test %q: got no error, want %q", testname, want)
			return false
		}
		return true
	}
	if want == nil {
		t.Errorf("test %q: unexpected error %q", testname, got)
	} else if got.Error() != want.Error() {
		t.Errorf("test %q: got error %q, want %q", testname, got, want)
	}
	return false
}
