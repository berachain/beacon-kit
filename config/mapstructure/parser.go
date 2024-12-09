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

package mapstructure

import (
	"net/url"
	"reflect"

	"github.com/berachain/beacon-kit/primitives/common"
	beaconurl "github.com/berachain/beacon-kit/primitives/net/url"
	"github.com/mitchellh/mapstructure"
)

// StringToExecutionAddressFunc returns a DecodeHookFunc that converts
// string to a `primitives.ExecutionAddresses` by parsing the string.
func StringToExecutionAddressFunc() mapstructure.DecodeHookFunc {
	return StringTo(
		func(s string) (common.ExecutionAddress, error) {
			return common.NewExecutionAddressFromHex(s), nil
		},
	)
}

// StringToDialURLFunc returns a DecodeHookFunc that converts
// string to *url.URL by parsing the string.
func StringToDialURLFunc() mapstructure.DecodeHookFunc {
	return StringTo(
		func(s string) (*url.URL, error) {
			url, err := url.Parse(s)
			if err != nil {
				return nil, err
			}
			return url, nil
		},
	)
}

// StringToConnectionURLFunc returns a DecodeHookFunc that converts
// string to *beaconurl.ConnectionURL by parsing the string.
func StringToConnectionURLFunc() mapstructure.DecodeHookFunc {
	return StringTo(beaconurl.NewFromRaw)
}

// StringTo is a helper function for creating DecodeHookFuncs that convert
// string to a specific type by parsing the string.
func StringTo[T any](
	constructor func(string) (T, error),
) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		var retType T
		if t != reflect.TypeOf(retType) {
			return data, nil
		}

		// Convert it by parsing
		//nolint:errcheck // should be safe
		return constructor(data.(string))
	}
}
