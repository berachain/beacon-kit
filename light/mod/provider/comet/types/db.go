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

package types

import (
	"strings"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/light/store"
	dbs "github.com/cometbft/cometbft/light/store/db"
)

const (
	primaryKey   = "primary"
	witnessesKey = "witnesses"
)

type DB struct {
	dbm.DB
	Store store.Store
}

func NewDB(dir string, chainID string) (*DB, error) {
	// db to store existing light client info
	// this includes any existing primary and witness addresses
	lightDB, err := dbm.NewGoLevelDB("light-client-db", dir)
	if err != nil {
		return nil, err
	}
	// create a db with chain id as prefix
	db := dbm.NewPrefixDB(lightDB, []byte(chainID))

	store := dbs.New(lightDB, chainID)
	return &DB{
		db,
		store,
	}, nil
}

// checkForExistingProviders checks the db for existing primary and witness
// providers
// Returns primary, witnesses, error.
func (db *DB) CheckForExistingProviders() (string, []string, error) {
	primaryBytes, err := db.Get([]byte(primaryKey))
	if err != nil {
		return "", []string{""}, err
	}
	witnessesBytes, err := db.Get([]byte(witnessesKey))
	if err != nil {
		return "", []string{""}, err
	}
	witnessesAddrs := strings.Split(string(witnessesBytes), ",")
	return string(primaryBytes), witnessesAddrs, nil
}

// stores the primary and witness providers in the db
// call this on the first run of the light client.
func (db *DB) SaveProviders(primaryAddr, witnessesAddrs string) error {
	err := db.Set([]byte(primaryKey), []byte(primaryAddr))
	if err != nil {
		return err
	}
	err = db.Set([]byte(witnessesKey), []byte(witnessesAddrs))
	if err != nil {
		return err
	}
	return nil
}
