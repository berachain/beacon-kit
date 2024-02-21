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

package ethclient_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/engine/ethclient"
	"github.com/itsdevbear/bolaris/engine/ethclient/mocks"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPayloadV3_BasicSanityCheck(t *testing.T) {
	tests := []struct {
		name            string
		versionedHashes []common.Hash
		parentBlockRoot *common.Hash
		ret             error
		wantErr         bool
	}{
		{
			name:            "success",
			versionedHashes: []common.Hash{{}, {}},
			parentBlockRoot: &common.Hash{},
			ret:             nil,
			wantErr:         false,
		},
		{
			name:            "error",
			versionedHashes: []common.Hash{{}, {}},
			parentBlockRoot: &common.Hash{},
			ret:             errors.New("my rpc error"),
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRPCClient := new(mocks.GethRPCClient)
			client := &ethclient.Eth1Client{GethRPCClient: mockRPCClient}

			payload := &enginev1.ExecutionPayloadDeneb{}

			mockRPCClient.EXPECT().
				CallContext(mock.Anything, mock.Anything,
					ethclient.NewPayloadMethodV3, payload,
					tt.versionedHashes, tt.parentBlockRoot).
				Return(tt.ret).
				Once()

			_, err := client.NewPayloadV3(
				context.Background(), payload, tt.versionedHashes, tt.parentBlockRoot,
			)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewPayloadV2_BasicSanityCheck(t *testing.T) {
	tests := []struct {
		name    string
		ret     error
		wantErr bool
	}{
		{
			name:    "success",
			ret:     nil,
			wantErr: false,
		},
		{
			name:    "rpc error",
			ret:     errors.New("rpc call failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRPCClient := new(mocks.GethRPCClient)
			client := &ethclient.Eth1Client{GethRPCClient: mockRPCClient}

			payload := &enginev1.ExecutionPayloadCapella{}

			mockRPCClient.EXPECT().
				CallContext(mock.Anything, mock.Anything,
					ethclient.NewPayloadMethodV2, payload).
				Return(tt.ret).
				Once()

			_, err := client.NewPayloadV2(context.Background(), payload)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestForkchoiceUpdatedV3_BasicSanityCheck(t *testing.T) {
	tests := []struct {
		name    string
		state   *enginev1.ForkchoiceState
		attrs   *enginev1.PayloadAttributesV3
		ret     error
		wantErr bool
	}{
		{
			name:    "nil response should error",
			state:   &enginev1.ForkchoiceState{},
			attrs:   &enginev1.PayloadAttributesV3{},
			ret:     nil,
			wantErr: true,
		},
		{
			name:    "call context returns error",
			state:   &enginev1.ForkchoiceState{},
			attrs:   &enginev1.PayloadAttributesV3{},
			ret:     errors.New("my custom rpc error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRPCClient := new(mocks.GethRPCClient)
			client := &ethclient.Eth1Client{GethRPCClient: mockRPCClient}

			mockRPCClient.EXPECT().
				CallContext(mock.Anything, mock.Anything,
					ethclient.ForkchoiceUpdatedMethodV3, tt.state, tt.attrs).
				Return(tt.ret).
				Once()

			_, err := client.ForkchoiceUpdatedV3(context.Background(), tt.state, tt.attrs)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestForkchoiceUpdatedV2_BasicSanityCheck(t *testing.T) {
	tests := []struct {
		name    string
		state   *enginev1.ForkchoiceState
		attrs   *enginev1.PayloadAttributesV2
		ret     error
		wantErr bool
	}{
		{
			name:    "nil response should error",
			state:   &enginev1.ForkchoiceState{},
			attrs:   &enginev1.PayloadAttributesV2{},
			ret:     nil,
			wantErr: true,
		},
		{
			name:    "call context returns error",
			state:   &enginev1.ForkchoiceState{},
			attrs:   &enginev1.PayloadAttributesV2{},
			ret:     errors.New("custom rpc error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRPCClient := new(mocks.GethRPCClient)
			client := &ethclient.Eth1Client{GethRPCClient: mockRPCClient}

			mockRPCClient.EXPECT().
				CallContext(mock.Anything, mock.Anything,
					ethclient.ForkchoiceUpdatedMethodV2,
					tt.state, tt.attrs,
				).
				Return(tt.ret).
				Once()

			_, err := client.ForkchoiceUpdatedV2(context.Background(), tt.state, tt.attrs)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetPayloadV3_BasicSanityCheck(t *testing.T) {
	tests := []struct {
		name    string
		pid     enginev1.PayloadIDBytes
		ret     error
		wantErr bool
	}{
		{
			name:    "nil response is desired",
			pid:     enginev1.PayloadIDBytes{},
			ret:     nil,
			wantErr: false,
		},
		{
			name:    "call context returns error",
			pid:     enginev1.PayloadIDBytes{},
			ret:     errors.New("my custom rpc error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRPCClient := new(mocks.GethRPCClient)
			client := &ethclient.Eth1Client{GethRPCClient: mockRPCClient}

			mockRPCClient.EXPECT().
				CallContext(mock.Anything, mock.Anything,
					"engine_getPayloadV3", tt.pid,
				).
				Return(tt.ret).
				Once()

			_, err := client.GetPayloadV3(context.Background(), tt.pid)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetPayloadV2_BasicSanityCheck(t *testing.T) {
	tests := []struct {
		name    string
		pid     enginev1.PayloadIDBytes
		ret     error
		wantErr bool
	}{
		{
			name:    "nil response is desired",
			pid:     enginev1.PayloadIDBytes{},
			ret:     nil,
			wantErr: false,
		},
		{
			name:    "call context returns error",
			pid:     enginev1.PayloadIDBytes{},
			ret:     errors.New("my custom rpc error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRPCClient := new(mocks.GethRPCClient)
			client := &ethclient.Eth1Client{GethRPCClient: mockRPCClient}

			mockRPCClient.EXPECT().
				CallContext(mock.Anything, mock.Anything,
					"engine_getPayloadV2", tt.pid,
				).
				Return(tt.ret).
				Once()

			_, err := client.GetPayloadV2(context.Background(), tt.pid)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
