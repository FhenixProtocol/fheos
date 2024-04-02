package ephemeraldb

// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package memorydb implements the key-value database layer based on memory maps.

import (
	"errors"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"sync"
)

var (
	// errMemorydbClosed is returned if a memory database was already closed at the
	// invocation of a data access operation.
	errMemorydbClosed = errors.New("database closed")

	// errMemorydbNotFound is returned if a key is requested that is not found in
	// the provided memory database.
	ErrMemorydbNotFound = errors.New("not found")
)

// ClearAll removes all entries from the database.
func (db *Database) ClearAll() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		return errMemorydbClosed
	}

	// Reinitialize the map to clear all entries
	db.db = make(map[string][]byte)
	return nil
}

type Database struct {
	db          map[string][]byte
	version     uint64
	encryptedDb map[types.Hash]*types.FheEncrypted
	lock        sync.RWMutex
}

func New() *Database {
	return &Database{
		db:          make(map[string][]byte),
		encryptedDb: make(map[types.Hash]*types.FheEncrypted),
	}
}

func (db *Database) GetVersion() (uint64, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return 0, errMemorydbClosed
	}
	return db.version, nil
}

func (db *Database) PutVersion(v uint64) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		return errMemorydbClosed
	}
	db.version = v
	return nil
}

func (db *Database) PutCt(h types.Hash, cipher *types.FheEncrypted) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		return errMemorydbClosed
	}

	// Store the encrypted data in a separate map.
	db.encryptedDb[h] = cipher
	return nil
}

func (db *Database) GetCt(h types.Hash) (*types.FheEncrypted, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return nil, errMemorydbClosed
	}

	cipher, ok := db.encryptedDb[h]
	if !ok {
		return nil, ErrMemorydbNotFound
	}
	// Assuming a deep copy is not required for *types.FheEncrypted, return directly.
	return cipher, nil
}
