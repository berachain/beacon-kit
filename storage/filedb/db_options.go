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

package filedb

import (
	"os"

	"github.com/berachain/beacon-kit/log"
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
func WithLogger(logger log.Logger) Option {
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
