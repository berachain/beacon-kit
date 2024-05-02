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
	"path/filepath"

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/cockroachdb/errors"
	"github.com/spf13/afero"
)

// DB represents a filesystem backed key-value store.
// It is useful for storing amounts of data that exceed what is
// performant to store in a traditional key-value database.
type DB struct {
	fs        afero.Fs
	logger    log.Logger[any]
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
		db.logger.Warn("overriding existing key", "key", key)
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
