// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package store

import (
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"google.golang.org/protobuf/proto"

	"cosmossdk.io/store"
)

type Forkchoice struct {
	store store.KVStore
}

func NewForkchoice(store store.KVStore) *Forkchoice {
	return &Forkchoice{
		store: store,
	}
}

func (f *Forkchoice) Store(fcs *enginev1.ForkchoiceState) error {
	bz, err := proto.Marshal(fcs)
	if err != nil {
		return err
	}
	f.store.Set([]byte("forkchoice"), bz)
	return nil
}

func (f *Forkchoice) Retrieve() (*enginev1.ForkchoiceState, error) {
	bz := f.store.Get([]byte("forkchoice"))
	fcs := &enginev1.ForkchoiceState{}
	if err := proto.Unmarshal(bz, fcs); err != nil {
		return nil, err
	}
	return fcs, nil
}
