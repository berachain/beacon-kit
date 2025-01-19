// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package components

import (
	"cosmossdk.io/core/store"
	"path/filepath"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	depositstore "github.com/berachain/beacon-kit/storage/deposit"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
)

// DepositStoreInput is the input for the dep inject framework.
type DepositStoreInput[
LoggerT log.AdvancedLogger[LoggerT],
] struct {
	depinject.In
	Logger  LoggerT
	AppOpts config.AppOptions
}

type syncedPDB struct {
	pdb dbm.DB
}

func (s *syncedPDB) Get(key []byte) ([]byte, error) {
	return s.pdb.Get(key)
}

func (s *syncedPDB) Has(key []byte) (bool, error) {
	return s.pdb.Has(key)
}

func (s *syncedPDB) Set(key, value []byte) error {
	return s.pdb.SetSync(key, value)
}

func (s *syncedPDB) Delete(key []byte) error {
	return s.pdb.Delete(key)
}

func (s *syncedPDB) Iterator(start, end []byte) (store.Iterator, error) {
	return s.pdb.Iterator(start, end)
}

func (s *syncedPDB) ReverseIterator(start, end []byte) (store.Iterator, error) {
	return s.pdb.ReverseIterator(start, end)
}

func (s *syncedPDB) NewBatch() store.Batch {
	return s.pdb.NewBatch()
}

func (s *syncedPDB) NewBatchWithSize(i int) store.Batch {
	return s.pdb.NewBatchWithSize(i)
}

func (s *syncedPDB) Close() error {
	return s.pdb.Close()
}

var _ store.KVStoreWithBatch = &syncedPDB{}

// ProvideDepositStore is a function that provides the module to the
// application.
func ProvideDepositStore[
LoggerT log.AdvancedLogger[LoggerT],
](
	in DepositStoreInput[LoggerT],
) (*depositstore.KVStore, error) {
	var (
		rootDir = cast.ToString(in.AppOpts.Get(flags.FlagHome))
		dataDir = filepath.Join(rootDir, "data")
		name    = "deposits"
	)

	pdb, err := dbm.NewDB(name, dbm.PebbleDBBackend, dataDir)
	if err != nil {
		return nil, err
	}
	spdb := &syncedPDB{pdb}

	// pass a closure to close the db as its not supported by the KVStoreService interface
	closeFunc := func() error { return spdb.Close() }

	return depositstore.NewStore(
		storage.NewKVStoreProvider(spdb),
		closeFunc,
		in.Logger.With("service", "deposit-store"),
	), nil
}
