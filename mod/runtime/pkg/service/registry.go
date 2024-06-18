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

package service

import (
	"context"
	"reflect"

	"github.com/berachain/beacon-kit/mod/log"
)

// Basic is the minimal interface for a service.
type Basic interface {
	// Start spawns any goroutines required by the service.
	Start(ctx context.Context) error
	// Name returns the name of the service.
	Name() string
}

// Registry provides a useful pattern for managing services.
// It allows for ease of dependency management and ensures services
// dependent on others use the same references in memory.
type Registry struct {
	// logger is the logger for the Registry.
	logger log.Logger[any]
	// services is a map of service type -> service instance.
	services map[string]Basic
	// serviceTypes is an ordered slice of registered service types.
	serviceTypes []string
}

// NewRegistry starts a registry instance for convenience.
func NewRegistry(opts ...RegistryOption) *Registry {
	r := &Registry{
		services: make(map[string]Basic),
	}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			panic(err)
		}
	}
	return r
}

// StartAll initialized each service in order of registration.
func (s *Registry) StartAll(ctx context.Context) error {
	s.logger.Info("starting services", "num", len(s.serviceTypes))
	for _, typeName := range s.serviceTypes {
		s.logger.Info("starting service", "type", typeName)
		svc := s.services[typeName]
		if svc == nil {
			s.logger.Error("service not found", "type", typeName)
			continue
		}

		if err := svc.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

// RegisterService appends a service constructor function to the service
// registry.
func (s *Registry) RegisterService(service Basic) error {
	typeName := service.Name()
	if _, exists := s.services[typeName]; exists {
		return errServiceAlreadyExists(typeName)
	}
	s.services[typeName] = service
	s.serviceTypes = append(s.serviceTypes, typeName)
	return nil
}

// FetchService takes in a struct pointer and sets the value of that pointer
// to a service currently stored in the service registry. This ensures the
// input argument is set to the right pointer that refers to the originally
// registered service.
func (s *Registry) FetchService(service interface{}) error {
	serviceType := reflect.TypeOf(service)
	if serviceType.Kind() != reflect.Ptr ||
		serviceType.Elem().Kind() != reflect.Ptr {
		return errInputIsNotPointer(serviceType)
	}

	element := reflect.ValueOf(service).Elem()

	typeName := ""
	for name, svc := range s.services {
		svcType := reflect.TypeOf(svc)
		if svcType.AssignableTo(serviceType.Elem()) {
			typeName = name
			break
		}
	}

	if typeName == "" {
		return errUnknownService(serviceType)
	}

	if running, ok := s.services[typeName]; ok {
		element.Set(reflect.ValueOf(running))
		return nil
	}
	return errUnknownService(serviceType)
}
