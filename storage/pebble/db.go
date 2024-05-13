//go:build amd64 || arm64

package pebble

import (
	"bytes"
	"encoding/gob"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"log"
	"sync"
)

var (
	instance *EthDbWrapper
	once     sync.Once
)

type EthDbWrapper struct {
	types.Storage
	db ethdb.Database
}

// NewPebbleStorage ensures a single EthDbWrapper instance
func NewStorage(path string) (*EthDbWrapper, error) {
	once.Do(func() {
		db, err := rawdb.NewPebbleDBDatabase(path, 128, 128, "fheos", false, false)
		if err != nil {
			log.Fatalf("Error creating PebbleDBDatabase: %v", err)
		}
		instance = &EthDbWrapper{db: db}
	})
	return instance, nil
}

//func (p *Storage) Put(t types.DataType, key []byte, val []byte) error {
//	// Use DataType as part of the key to differentiate data types
//	prefixedKey := append([]byte(fmt.Sprintf("%d_", t)), key...)
//	return p.db.Set(prefixedKey, val, pebble.Sync)
//}
//
//func (p *Storage) Get(t types.DataType, key []byte) ([]byte, error) {
//	prefixedKey := append([]byte(fmt.Sprintf("%d_", t)), key...)
//	val, closer, err := p.db.Get(prefixedKey)
//	if err != nil {
//		return nil, err
//	}
//	defer closer.Close()
//
//	// Make a copy of the data since val becomes invalid after the closer is called
//	valCopy := make([]byte, len(val))
//	copy(valCopy, val)
//
//	return valCopy, nil
//}

func (p *EthDbWrapper) GetVersion() (uint64, error) {
	key := []byte("version")
	val, err := p.db.Get(key)
	if err != nil {
		return 0, err
	}

	// Assuming the version is stored as a uint64
	var version uint64
	buf := bytes.NewBuffer(val)
	err = gob.NewDecoder(buf).Decode(&version)
	if err != nil {
		return 0, err
	}

	return version, nil
}

func (p *EthDbWrapper) PutVersion(v uint64) error {
	key := []byte("version")
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(v)
	if err != nil {
		return err
	}

	return p.db.Put(key, buf.Bytes())
}

func (p *EthDbWrapper) PutCt(h types.Hash, cipher *types.CipherTextRepresentation) error {
	// Serialize Ciphertext
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(cipher)
	if err != nil {
		return err
	}

	// Use hash as key
	return p.db.Put(h[:], buf.Bytes())
}

func (p *EthDbWrapper) GetCt(h types.Hash) (*types.CipherTextRepresentation, error) {
	val, err := p.db.Get(h[:])
	if err != nil {
		return nil, err
	}

	var cipher types.CipherTextRepresentation
	err = gob.NewDecoder(bytes.NewBuffer(val)).Decode(&cipher)
	if err != nil {
		return nil, err
	}

	return &cipher, nil
}
