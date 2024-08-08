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

package services

import (
	"context"

	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// ReflectionService implements the cosmos.reflection.v1 service.
type ReflectionService struct {
	reflectionv1.UnimplementedReflectionServiceServer
	files *descriptorpb.FileDescriptorSet
}

func NewReflectionService() (*ReflectionService, error) {
	fds, err := proto.MergedGlobalFileDescriptors()
	if err != nil {
		return nil, err
	}

	return &ReflectionService{files: fds}, nil
}

func (r ReflectionService) FileDescriptors(
	_ context.Context,
	_ *reflectionv1.FileDescriptorsRequest,
) (*reflectionv1.FileDescriptorsResponse, error) {
	return &reflectionv1.FileDescriptorsResponse{
		Files: r.files.GetFile(),
	}, nil
}

var _ reflectionv1.ReflectionServiceServer = &ReflectionService{}
