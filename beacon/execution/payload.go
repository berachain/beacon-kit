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

package execution

import (
	"context"
	"crypto/rand"

	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

// TODO: this whole function is just a hacky patch, doesn't do anything significant.
func (s *Service) getPayloadAttributes(
	_ context.Context, slot primitives.Slot, timestamp uint64,
) (payloadattribute.Attributer, error) {
	// TODO: modularize andn make better.
	requestedVersion := s.beaconCfg.ActiveForkVersion(primitives.Epoch(slot))
	emptyAttri := payloadattribute.EmptyWithVersion(requestedVersion)

	// TODO: this is a hack to fill the PrevRandao field. It is not verifiable or safe
	// or anything, it's literally just to get basic functionality.
	var random [32]byte
	if _, err := rand.Read(random[:]); err != nil {
		return emptyAttri, err
	}

	// // TODO: need to support multiple fork versions.
	// if requestedVersion != 3 { //
	// 	s.logger.Error("Could not get payload attribute due to unknown state version")
	// 	return emptyAttri, nil
	// }

	return payloadattribute.New(&enginev1.PayloadAttributesV2{
		Timestamp:             timestamp,
		SuggestedFeeRecipient: s.etherbase.Bytes(),
		// TODO: support withdrawls of BGT here.
		Withdrawals: nil,
		// TODO: we need to implement this correctly.
		PrevRandao: append([]byte{}, random[:]...),
	})
	// var attr payloadattribute.Attributer
	// switch requestedVersion {
	// case version.Deneb:
	// 	withdrawals, err := st.ExpectedWithdrawals()
	// 	if err != nil {
	// 		// log.WithError(err).Error("Could not get expected withdrawals to get payload attribute")
	// 		return emptyAttri, err
	// 	}
	// 	attr, err = payloadattribute.New(&enginev1.PayloadAttributesV3{
	// 		Timestamp:             uint64(t.Unix()),
	// 		PrevRandao:            prevRando,
	// 		SuggestedFeeRecipient: val.FeeRecipient[:],
	// 		Withdrawals:           withdrawals,
	// 		ParentBeaconBlockRoot: headRoot,
	// 	})
	// 	if err != nil {
	// 		log.WithError(err).Error("Could not get payload attribute")
	// 		return false, emptyAttri
	// 	}
	// case version.Capella:
	// 	withdrawals, err := st.ExpectedWithdrawals()
	// 	if err != nil {
	// 		log.WithError(err).Error("Could not get expected withdrawals to get payload attribute")
	// 		return false, emptyAttri
	// 	}
	// 	attr, err = payloadattribute.New(&enginev1.PayloadAttributesV2{
	// 		Timestamp:             uint64(t.Unix()),
	// 		PrevRandao:            prevRando,
	// 		SuggestedFeeRecipient: val.FeeRecipient[:],
	// 		Withdrawals:           withdrawals,
	// 	})
	// 	if err != nil {
	// 		log.WithError(err).Error("Could not get payload attribute")
	// 		return false, emptyAttri
	// 	}
	// case version.Bellatrix:
	// 	attr, err = payloadattribute.New(&enginev1.PayloadAttributes{
	// 		Timestamp:             uint64(t.Unix()),
	// 		PrevRandao:            prevRando,
	// 		SuggestedFeeRecipient: val.FeeRecipient[:],
	// 	})
	// 	if err != nil {
	// 		log.WithError(err).Error("Could not get payload attribute")
	// 		return false, emptyAttri
	// 	}
	// default:
	// 	log.WithField("vers
	// on", version).Error("Could not get payload attribute due to unknown state version")
	// 	return false, emptyAttri
	// }
}
