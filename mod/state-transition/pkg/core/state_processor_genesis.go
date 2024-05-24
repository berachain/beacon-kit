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

package core

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// TODO: properly create a genesis type and use the FromEth1 ideaology.
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, ContextT,
]) ProcessGenesisState(
	_ ContextT,
	st BeaconStateT,
	data *deneb.BeaconState,
) ([]*transition.ValidatorUpdate, error) {
	// TODO: this should not just dumb a beaconstate object.
	if err := st.WriteGenesisStateDeneb(data); err != nil {
		return nil, err
	}

	updates, err := sp.processSyncCommitteeUpdates(st)
	if err != nil {
		return nil, err
	}

	st.Save()
	return updates, nil
}
