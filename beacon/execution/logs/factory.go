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

package logs

import (
	"errors"
	"reflect"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type Factory struct {
	addressToAbi map[common.Address]*ethabi.ABI
	sigToName    map[common.Hash]string
	sigToType    map[common.Hash]reflect.Type
}

func NewFactory() *Factory {
	return &Factory{
		addressToAbi: make(map[common.Address]*ethabi.ABI),
		sigToName:    make(map[common.Hash]string),
		sigToType:    make(map[common.Hash]reflect.Type),
	}
}

func (f *Factory) RegisterLog(
	contractAddress common.Address,
	contractAbi *ethabi.ABI,
	eventName string,
	eventType reflect.Type,
) {
	eventID := contractAbi.Events[eventName].ID
	f.addressToAbi[contractAddress] = contractAbi
	f.sigToName[eventID] = eventName
	f.sigToType[eventID] = eventType
}

func (f *Factory) UnmarshalEthLog(log *ethtypes.Log) (any, error) {
	var (
		contractAbi *ethabi.ABI
		eventType   reflect.Type
		eventName   string
		ok          bool
	)

	sig := log.Topics[0]
	if contractAbi, ok = f.addressToAbi[log.Address]; !ok {
		return nil, errors.New("abi not found for log address")
	}
	if eventType, ok = f.sigToType[sig]; !ok {
		return nil, errors.New("type not found for log signature")
	}
	if eventName, ok = f.sigToName[sig]; !ok {
		return nil, errors.New("name not found for log signature")
	}

	into := reflect.New(eventType).Interface()
	if err := contractAbi.UnpackIntoInterface(into, eventName, log.Data); err != nil {
		return nil, err
	}
	return into, nil
}
