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

package constraints

// ForkTyped is a constraint that requires a type to have an Empty method.
type ForkTyped[SelfT any] interface {
	EmptyWithVersion[SelfT]
	Versionable
	Nillable
}

// EngineType represents the constraints for a type that is
// used within the context of sending over the EngineAPI.
type EngineType[SelfT any] interface {
	EmptyWithVersion[SelfT]
	Versionable
	Nillable
	JSONMarshallable
}

// EmptyWithVersion is a constraint that requires a type to have an Empty
// method.
type EmptyWithVersion[SelfT any] interface {
	Empty(uint32) SelfT
}

// Empty is a constraint that requires a type to have an Empty method.
type Empty[SelfT any] interface {
	Empty() SelfT
}

// Nillable is a constraint that requires a type to have an IsNil method.
type Nillable interface {
	IsNil() bool
}

// Versionable is a constraint that requires a type to have a Version method.
type Versionable interface {
	Version() uint32
}
