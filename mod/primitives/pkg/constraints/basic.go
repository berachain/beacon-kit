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

import (
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

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

// EmptyWithForkVersion is a constraint that requires a type to have an Empty
// method.
type EmptyWithVersion[SelfT any] interface {
	Empty(uint32) SelfT
}

// Empty is a constraint that requires a type to have an Empty method.
type Empty[SelfT any] interface {
	Empty() SelfT
}

// IsNil is a constraint that requires a type to have an IsNil method.
type Nillable interface {
	IsNil() bool
}

// Versionable is a constraint that requires a type to have a Version method.
type Versionable interface {
	Version() uint32
}

// SSZField is an interface for types that can be used as SSZ fields.
// It requires methods for creating an empty instance, computing a hash tree
// root,
// and allows pointer access to the underlying type.
type SSZField[SelfPtrT, SelfT any] interface {
	// Empty returns a new empty instance of the type.
	Empty() SelfPtrT
	// HashTreeRootWith computes the hash tree root of the object using the
	// provided HashWalker.
	HashTreeRootWith(fastssz.HashWalker) error
	// Allows pointer access to the underlying type.
	*SelfT
}

// StaticSSZField is an interface for SSZ fields with a static size.
// It embeds the SSZField interface and the ssz.StaticObject interface.
type StaticSSZField[SelfPtrT, SelfT any] interface {
	ssz.StaticObject
	SSZField[SelfPtrT, SelfT]
}

// DynamicSSZField is an interface for SSZ fields with a dynamic size.
// It embeds the SSZField interface and the ssz.DynamicObject interface.
type DynamicSSZField[SelfPtrT, SelfT any] interface {
	ssz.DynamicObject
	SSZField[SelfPtrT, SelfT]
}
