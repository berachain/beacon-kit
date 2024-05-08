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

package service

import (
	"context"
	"reflect"

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/sourcegraph/conc"
)

// Basic is the minimal interface for a service.
type Basic interface {
	// Start spawns any goroutines required by the service.
	Start(ctx context.Context)
	// Name returns the name of the service.
	Name() string
	// Status returns error if the service is not considered healthy.
	Status() error
	// WaitForHealthy blocks until the service is healthy.
	WaitForHealthy(ctx context.Context)
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
func (s *Registry) StartAll(ctx context.Context) {
	s.logger.Info("starting services", "num", len(s.serviceTypes))
	for _, typeName := range s.serviceTypes {
		s.logger.Info("starting service", "type", typeName)
		svc := s.services[typeName]
		if svc == nil {
			s.logger.Error("service not found", "type", typeName)
			continue
		}
		svc.Start(ctx)
	}
}

// Statuses returns a map of Service type -> error. The map will be populated
// with the results of each service.Status() method call.
func (s *Registry) Statuses(services ...string) map[string]error {
	if len(services) == 0 {
		services = s.serviceTypes
	}

	m := make(map[string]error)
	for _, typeName := range services {
		if typeName == "" {
			s.logger.Error("empty service type")
		}
		svc := s.services[typeName]
		if svc == nil {
			s.logger.Error("service not found", "type", typeName)
			continue
		}
		//#nosec:G703 // false positive?
		m[typeName] = svc.Status()
	}
	return m
}

// Statuses returns a map of Service type -> error. The map will be populated
// with the results of each service.Status() method call.
func (s *Registry) WaitForHealthy(ctx context.Context, services ...string) {
	wg := conc.NewWaitGroup()
	for _, typeName := range services {
		wg.Go(func() { s.services[typeName].WaitForHealthy(ctx) })
	}
	wg.Wait()
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
