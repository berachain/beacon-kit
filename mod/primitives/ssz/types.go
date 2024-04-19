// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package ssz

// Basic defines an interface for SSZ basic types which includes methods for
// determining the size of the SSZ encoding and computing the hash tree root.
type Basic[SpecT any, RootT ~[32]byte] interface {
	// SizeSSZ returns the size in bytes of the SSZ-encoded data.
	SizeSSZ() int
	// HashTreeRoot computes and returns the hash tree root of the data as RootT
	// and an error if the computation fails.
	HashTreeRoot( /*...args*/ ) (RootT, error)
}

// Composite is an interface that embeds the Basic interface. It is used for
// types that are composed of other SSZ encodable values.
type Composite[SpecT any, RootT ~[32]byte] interface {
	Basic[SpecT, RootT]
}

// Container is an interface for SSZ container types that can be marshaled and
// unmarshaled.
type Container[SpecT any, RootT ~[32]byte] interface {
	Composite[SpecT, RootT]
}
