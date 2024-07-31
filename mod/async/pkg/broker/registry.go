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

package broker

import (
	"context"
	"reflect"

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
)

// grand hack central for now.
type Registry struct {
	brokers map[string]types.Broker
}

func NewRegistry() *Registry {
	return &Registry{
		brokers: make(map[string]types.Broker),
	}
}

func (r *Registry) RegisterBroker(broker types.Broker) {
	r.brokers[broker.Name()] = broker
}

func (r *Registry) FetchBroker(broker types.Broker) error {
	brokerType := reflect.TypeOf(broker)
	if brokerType.Kind() != reflect.Ptr {
		return types.ErrInputIsNotPointer(brokerType)
	}

	elem := reflect.ValueOf(broker).Elem()

	typeName := ""
	for name, b := range r.brokers {
		bType := reflect.TypeOf(b)
		if bType.AssignableTo(bType.Elem()) {
			typeName = name
			break
		}
	}

	if typeName == "" {
		return types.ErrUnknownService(brokerType)
	}

	if registeredBroker, ok := r.brokers[typeName]; ok {
		elem.Set(reflect.ValueOf(registeredBroker))
	}

	return nil
}

func (s *Registry) StartAll(ctx context.Context) {
	for _, broker := range s.brokers {
		broker.Start(ctx)
	}
}
