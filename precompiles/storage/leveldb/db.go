//go:build amd64 || arm64

package leveldb

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
	"sync"
)

const (
	version types.DataType = iota
	ct
)

var Lock sync.RWMutex

type Storage struct {
	dbPath string
}

func NewStorage(path string) Storage {
	storage := Storage{
		dbPath: path,
	}

	return storage
}

func (store Storage) OpenDB(readonly bool) *leveldb.DB {
	if readonly {
		Lock.RLock()
	} else {
		Lock.Lock()
	}

	db, err := leveldb.OpenFile(store.dbPath, &opt.Options{ReadOnly: readonly})
	if err != nil {
		log.Fatalf("failed to open fheos db. %s: %s", "err", err)
	}

	return db
}

func closeDB(db *leveldb.DB, readonly bool) {
	if readonly {
		defer Lock.RUnlock()
	} else {
		defer Lock.Unlock()
	}

	err := db.Close()
	if err != nil {
		log.Fatalf("failed to close fheos db. %s: %s", "err", err)
	}
}

func (store Storage) Put(t types.DataType, key []byte, val []byte) error {
	db := store.OpenDB(false)
	defer closeDB(db, false)

	tb := make([]byte, 8)
	binary.BigEndian.PutUint64(tb, uint64(t))
	extendedKey := append(tb, key...)

	err := db.Put(extendedKey, val, nil)
	if err != nil {
		return err
	}

	return nil
}

func (store Storage) Get(t types.DataType, key []byte) ([]byte, error) {
	db := store.OpenDB(true)
	defer closeDB(db, true)

	tb := make([]byte, 8)
	binary.BigEndian.PutUint64(tb, uint64(t))
	extendedKey := append(tb, key...)

	val, err := db.Get(extendedKey, nil)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (store Storage) GetVersion() (uint64, error) {
	v, err := store.Get(version, []byte{})
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(v), nil
}

func (store Storage) PutVersion(v uint64) error {
	vb := make([]byte, 8)
	binary.BigEndian.PutUint64(vb, v)

	return store.Put(version, []byte{}, vb)
}

func (store Storage) PutCt(h fhe.Hash, cipher *fhe.FheEncrypted) error {
	var cipherBuffer bytes.Buffer
	enc := gob.NewEncoder(&cipherBuffer)
	err := enc.Encode(*cipher)
	if err != nil {
		return err
	}

	return store.Put(ct, h[:], cipherBuffer.Bytes())
}

func (store Storage) GetCt(h fhe.Hash) (*fhe.FheEncrypted, error) {
	v, err := store.Get(ct, h[:])
	if err != nil {
		return nil, err
	}

	var cipher fhe.FheEncrypted
	dec := gob.NewDecoder(bytes.NewReader(v))
	err = dec.Decode(&cipher)
	if err != nil {
		return nil, err
	}

	return &cipher, nil
}
