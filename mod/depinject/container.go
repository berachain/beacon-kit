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

package depinject

import (
	"reflect"

	"go.uber.org/dig"
)

// Container is a wrapper around dig.Container that provides syntactic sugar
// and adds convenience methods for building and injecting dependencies.
type Container struct {
	*dig.Container
}

// NewContainer creates a new Container.
func NewContainer() *Container {
	return &Container{
		Container: dig.New(),
	}
}

// Provide is a helper function that provides multiple constructors to the
// container. It takes an arbitrary number of constructor functions, and adds
// them to the container.
func (c *Container) Provide(constructors ...any) error {
	for _, constructor := range constructors {
		if err := c.Container.Provide(constructor); err != nil {
			return ProvideError(err, constructor)
		}
	}
	return nil
}

// Supply is a helper function that supplies multiple values to the container.
// It takes an arbitrary number of values of possibly different types,
// and adds them to the container.
func (c *Container) Supply(values ...any) error {
	for _, value := range values {
		valueType := reflect.TypeOf(value)
		provideFunc := reflect.MakeFunc(
			reflect.FuncOf(nil, []reflect.Type{valueType}, false),
			func(args []reflect.Value) (results []reflect.Value) {
				return []reflect.Value{reflect.ValueOf(value)}
			},
		)
		if err := c.Provide(provideFunc.Interface()); err != nil {
			return err
		}
	}
	return nil
}

// Inject is a helper function that retrieves multiple dependencies from
// the container. It takes an arbitrary number of pointers to objects of
// possibly different types, invokes the container for each type, and assigns
// the values to the provided pointers.
func (c *Container) Inject(targets ...interface{}) error {
	for _, target := range targets {
		targetVal := reflect.ValueOf(target)
		if targetVal.Kind() != reflect.Ptr || targetVal.IsNil() {
			return ErrTargetMustBePointer
		}
		targetType := targetVal.Elem().Type()

		// This function infers the type of the target from the provided pointer
		// and returns a function that sets the target to the provided argument.
		fn := reflect.MakeFunc(
			reflect.FuncOf([]reflect.Type{targetType}, []reflect.Type{}, false),
			func(args []reflect.Value) (results []reflect.Value) {
				targetVal.Elem().Set(args[0])
				return []reflect.Value{}
			},
		)

		if err := c.Invoke(fn.Interface()); err != nil {
			return err
		}
	}

	return nil
}

func (c *Container) Invoke(fn any) error {
	if err := c.Container.Invoke(fn); err != nil {
		return err
	}
	return nil
}
