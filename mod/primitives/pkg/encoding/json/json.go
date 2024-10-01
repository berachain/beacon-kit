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
//
//nolint:gochecknoglobals // intentional aliases.
package json

import (
	"encoding/json"
)

// Marshaler is the interface implemented by types that
// can marshal themselves into valid JSON.
type Marshaler = gojson.Marshaler

// Unmarshaler is the interface implemented by types
// that can unmarshal a JSON description of themselves.
type Unmarshaler = gojson.Unmarshaler

// Marshal is a wrapper for gojson.Marshal, which provides high-performance JSON
// encoding.
var Marshal = json.Marshal

// MarshalIndent is a wrapper for gojson.MarshalIndent, which provides
// high-performance JSON encoding with indentation.
var MarshalIndent = json.MarshalIndent

// Unmarshal is a wrapper for gojson.Unmarshal, which provides high-performance
// JSON decoding.
var Unmarshal = json.Unmarshal

// RawMessage is an alias for json.RawMessage, represensting a raw encoded JSON
// value. It implements Marshaler and Unmarshaler and can be used to delay JSON
// decoding or precompute a JSON encoding.
type RawMessage = json.RawMessage
