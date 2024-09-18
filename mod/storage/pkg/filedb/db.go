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
	"path/filepath"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/spf13/afero"
)

// DB represents a filesystem backed key-value store.
// It is useful for storing amounts of data that exceed what is
// performant to store in a traditional key-value database.
type DB struct {
	fs        afero.Fs
	logger    log.Logger
	rootDir   string
	extension string
	dirPerms  os.FileMode
}

// NewDB creates a new instance of the DB.
func NewDB(opts ...Option) *DB {
	db := &DB{}
	for _, opt := range opts {
		if err := opt(db); err != nil {
			panic(errors.Wrap(err, "failed to apply option"))
		}
	}

	db.fs = afero.NewBasePathFs(afero.NewOsFs(), db.rootDir)
	return db
}

// Get retrieves the value for a key.
func (db *DB) Get(key []byte) ([]byte, error) {
	return afero.ReadFile(db.fs, db.pathForKey(key))
}

// Has returns true if the key exists in the database.
func (db *DB) Has(key []byte) (bool, error) {
	exists, err := afero.Exists(db.fs, db.pathForKey(key))
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Set stores the value for a key.
func (db *DB) Set(key []byte, value []byte) error {
	if exists, err := afero.Exists(db.fs, db.pathForKey(key)); err != nil {
		return err
	} else if exists {
		db.logger.Warn("Overriding existing key", "key", key)
	}

	if err := db.fs.MkdirAll(
		filepath.Dir(db.pathForKey(key)), db.dirPerms,
	); err != nil {
		return err
	}

	file, err := db.fs.Create(db.pathForKey(key))
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer file.Close()

	n, err := file.Write(value)
	if err != nil {
		return errors.Wrap(err, "failed to write to file")
	}
	db.logger.Debug("wrote %d bytes to %s", n, db.pathForKey(key))

	return nil
}

// Delete removes the value for a key.
func (db *DB) Delete(key []byte) error {
	return db.fs.RemoveAll(db.pathForKey(key))
}

// pathForKey returns the path for a key.
// TODO: for efficient storage we should expand this path
func (db *DB) pathForKey(key []byte) string {
	return string(key) + "." + db.extension
}
