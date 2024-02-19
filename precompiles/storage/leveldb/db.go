//go:build amd64 || arm64

package leveldb

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/ethereum/go-ethereum/log"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/go-tfhe"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"sync"
)

var logger log.Logger

const (
	version types.DataType = iota
	ct
)

var LevelDbLock sync.RWMutex

type LevelDbStorage struct {
	DbPath string
}

func (store LevelDbStorage) OpenDB(readonly bool) *leveldb.DB {
	if readonly {
		LevelDbLock.RLock()
	} else {
		LevelDbLock.Lock()
	}

	db, err := leveldb.OpenFile(store.DbPath, &opt.Options{ReadOnly: readonly})
	if err != nil {
		logger.Error("failed to open fheos db", "err", err)
		panic(err)
	}

	return db
}

func closeDB(db *leveldb.DB, readonly bool) {
	if readonly {
		defer LevelDbLock.RUnlock()
	} else {
		defer LevelDbLock.Unlock()
	}

	err := db.Close()
	if err != nil {
		logger.Error("failed to close fheos db", "err", err)
		panic(err)
	}

	logger.Debug("fheos db closed")
}
func (store LevelDbStorage) Put(t types.DataType, key []byte, val []byte) error {
	db := store.OpenDB(false)
	defer closeDB(db, false)

	tb := make([]byte, 8)
	binary.BigEndian.PutUint64(tb, uint64(t))
	extendedKey := append(tb, key...)

	err := db.Put(extendedKey, val, nil)
	if err != nil {
		logger.Error("failed to write into fheos db", "err", err)
		return err
	}

	return nil
}

func (store LevelDbStorage) Size() uint64 {
	db := store.OpenDB(true)
	defer closeDB(db, true)

	dbStats := leveldb.DBStats{}
	err := db.Stats(&dbStats)
	if err != nil {
		logger.Error("failed to get stats from db")
		return 0
	}

	size := dbStats.LevelSizes.Sum()
	return uint64(size)
}

func (store LevelDbStorage) Get(t types.DataType, key []byte) ([]byte, error) {
	db := store.OpenDB(true)
	defer closeDB(db, true)

	tb := make([]byte, 8)
	binary.BigEndian.PutUint64(tb, uint64(t))
	extendedKey := append(tb, key...)

	val, err := db.Get(extendedKey, nil)
	if err != nil {
		logger.Error("failed to read from fheos db", "err", err)
		return nil, err
	}

	return val, nil
}

func (store LevelDbStorage) GetVersion() (uint64, error) {
	v, err := store.Get(version, []byte{})
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(v), nil
}

func (store LevelDbStorage) PutVersion(v uint64) error {
	vb := make([]byte, 8)
	binary.BigEndian.PutUint64(vb, v)

	return store.Put(version, []byte{}, vb)
}

func (store LevelDbStorage) PutCt(h tfhe.Hash, cipher *tfhe.Ciphertext) error {
	var cipherBuffer bytes.Buffer
	enc := gob.NewEncoder(&cipherBuffer)
	err := enc.Encode(*cipher)
	if err != nil {
		logger.Error("failed to encode ciphertext", "err", err)
		return err
	}

	return store.Put(ct, h[:], cipherBuffer.Bytes())
}

func (store LevelDbStorage) GetCt(h tfhe.Hash) (*tfhe.Ciphertext, error) {
	v, err := store.Get(ct, h[:])
	if err != nil {
		return nil, err
	}

	var cipher tfhe.Ciphertext
	dec := gob.NewDecoder(bytes.NewReader(v))
	err = dec.Decode(&cipher)
	if err != nil {
		logger.Error("failed to decode ciphertext", "err", err)
		return nil, err
	}

	return &cipher, nil
}
