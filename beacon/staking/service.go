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

package staking

import (
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	stakingtypes "github.com/berachain/beacon-kit/beacon/staking/types"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/lib/skiplist"
	"github.com/berachain/beacon-kit/runtime/service"
)

// Service represents the staking service.
type Service struct {
	service.BaseService

	// TODO: make these persist to disk for restarts.
	depositQueue  *skiplist.Skiplist[enginetypes.Withdrawal]
	redirectQueue *skiplist.Skiplist[stakingtypes.Redirect]
	withdrawQueue *skiplist.Skiplist[beacontypes.Deposit]
	// vcp is responsible for applying validator set changes.
	vcp ValsetChangeProvider
}
