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

package types

type Type uint8

const (
	Basic Type = iota
	Vector
	List
	Container
)

// IsBasic returns true if the type is a basic type.
func (t Type) IsBasic() bool {
	return t == Basic
}

// IsElements returns true if the type is an enumerable type.
func (t Type) IsElements() bool {
	return t == Vector || t == List
}

// IsComposite returns true if the type is a composite type.
func (t Type) IsComposite() bool {
	return t == Vector || t == List || t == Container
}

// IsEnumerable returns true if the type is an enumerable type.
func (t Type) IsEnumerable() bool {
	return t == Vector || t == List
}

// IsContainer returns true if the type is a container type.
func (t Type) IsContainer() bool {
	return t == Container
}
