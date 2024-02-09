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

package service

import (
	"context"
	"reflect"

	"cosmossdk.io/log"
)

// Basic is the minimal interface for a service.
type Basic interface {
	// Start spawns any goroutines required by the service.
	Start(ctx context.Context)
	// Status returns error if the service is not considered healthy.
	Status() error
}

// Registry provides a useful pattern for managing services.
// It allows for ease of dependency management and ensures services
// dependent on others use the same references in memory.
type Registry struct {
	// logger is the logger for the Registry.
	logger log.Logger
	// services is a map of service type -> service instance.
	services map[reflect.Type]Basic
	// serviceTypes is an ordered slice of registered service types.
	serviceTypes []reflect.Type // keep an ordered slice of registered service types.
}

// NewRegistry starts a registry instance for convenience.
func NewRegistry(opts ...RegistryOption) *Registry {
	r := &Registry{
		services: make(map[reflect.Type]Basic),
	}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			panic(err)
		}
	}
	return r
}

// StartAll initialized each service in order of registration.
func (s *Registry) StartAll(ctx context.Context) {
	s.logger.Info("starting services", "num", len(s.serviceTypes))
	for _, kind := range s.serviceTypes {
		s.logger.Info("starting service", "type", kind)
		go s.services[kind].Start(ctx)
	}
}

// Statuses returns a map of Service type -> error. The map will be populated
// with the results of each service.Status() method call.
func (s *Registry) Statuses() map[reflect.Type]error {
	m := make(map[reflect.Type]error, len(s.serviceTypes))
	for _, kind := range s.serviceTypes {
		m[kind] = s.services[kind].Status() //#nosec:G703 // todo:test.
	}
	return m
}

// RegisterService appends a service constructor function to the service
// registry.
func (s *Registry) RegisterService(service Basic) error {
	kind := reflect.TypeOf(service)
	if _, exists := s.services[kind]; exists {
		return errServiceAlreadyExists(kind.String())
	}
	s.services[kind] = service
	s.serviceTypes = append(s.serviceTypes, kind)
	return nil
}

// FetchService takes in a struct pointer and sets the value of that pointer
// to a service currently stored in the service registry. This ensures the input argument is
// set to the right pointer that refers to the originally registered service.
func (s *Registry) FetchService(service interface{}) error {
	if reflect.TypeOf(service).Kind() != reflect.Ptr {
		return errInputIsNotPointer(reflect.TypeOf(service))
	}
	element := reflect.ValueOf(service).Elem()
	if running, ok := s.services[element.Type()]; ok {
		element.Set(reflect.ValueOf(running))
		return nil
	}
	return errUnknownService(reflect.TypeOf(service))
}
