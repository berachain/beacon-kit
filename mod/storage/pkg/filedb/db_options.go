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

package filedb

import (
	"os"

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/spf13/afero"
)

type Option func(*DB) error

// WithAferoFS sets the filesystem for the database.
// NOTE: Should only be used for testing.
func WithAferoFS(fs afero.Fs) Option {
	return func(db *DB) error {
		db.fs = fs
		return nil
	}
}

// WithDirectoryPermissions sets the permissions for the directory.
func WithDirectoryPermissions(permissions os.FileMode) Option {
	return func(db *DB) error {
		db.dirPerms = permissions
		return nil
	}
}

// WithFileExtension sets the file extension for the database.
func WithFileExtension(extension string) Option {
	return func(db *DB) error {
		db.extension = extension
		return nil
	}
}

// WithLogger sets the logger for the database.
func WithLogger(logger log.Logger[any]) Option {
	return func(db *DB) error {
		db.logger = logger
		return nil
	}
}

// WithRootDirectory sets the root directory for the database.
func WithRootDirectory(rootDir string) Option {
	return func(db *DB) error {
		db.rootDir = rootDir
		return nil
	}
}
