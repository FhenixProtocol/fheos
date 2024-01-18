package precompiles

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"sync"

	"github.com/ethereum/go-ethereum/params"
	"github.com/fhenixprotocol/go-tfhe"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var LevelDbLock sync.RWMutex

type LevelDbStorage struct {
	dbPath string
	burner GasBurner
}

const DBPath = "/home/user/fhenix/fheosdb"
const StorageReadCost = params.SloadGasEIP2200
const StorageWriteCost = params.SstoreSetGasEIP2200

func InitStorage(burner GasBurner) Storage {
	storage := LevelDbStorage{
		dbPath: DBPath,
		burner: burner,
	}

	return storage
}

func (store LevelDbStorage) OpenDB(readonly bool) *leveldb.DB {
	if readonly {
		LevelDbLock.RLock()
	} else {
		LevelDbLock.Lock()
	}

	db, err := leveldb.OpenFile(store.dbPath, &opt.Options{ReadOnly: readonly})
	if err != nil {
		logger.Error("failed to open fheos db ", err)
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
		logger.Error("failed to close fheos db ", err)
		panic(err)
	}

	logger.Debug("fheos db closed")
}
func (store LevelDbStorage) Put(t DataType, key []byte, val []byte) error {
	db := store.OpenDB(false)
	defer closeDB(db, false)

	tb := make([]byte, 8)
	binary.BigEndian.PutUint64(tb, uint64(t))
	extendedKey := append(tb, key...)

	err := db.Put(extendedKey, val, nil)
	if err != nil {
		logger.Error("failed to write into fheos db ", err)
		return err
	}

	return nil
}

func (store LevelDbStorage) Get(t DataType, key []byte) ([]byte, error) {
	db := store.OpenDB(true)
	defer closeDB(db, true)

	tb := make([]byte, 8)
	binary.BigEndian.PutUint64(tb, uint64(t))
	extendedKey := append(tb, key...)

	val, err := db.Get(extendedKey, nil)
	if err != nil {
		logger.Error("failed to read from fheos db ", err)
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
		logger.Error("failed to encode ciphertext ", err)
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
		logger.Error("failed to decode ciphertext ", err)
		return nil, err
	}

	return &cipher, nil
}
