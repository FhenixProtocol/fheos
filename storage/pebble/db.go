//go:build amd64 || arm64

package pebble

import (
	"bytes"
	"encoding/gob"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/fhenixprotocol/fheos/precompiles/types"
)

var (
	instance *EthDbWrapper
	once     sync.Once
)

type EthDbWrapper struct {
	types.Storage
	db ethdb.Database
}

// NewStorage ensures a single EthDbWrapper instance
func NewStorage(path string) (*EthDbWrapper, error) {
	once.Do(func() {
		db, err := rawdb.NewPebbleDBDatabase(path, 128, 128, "fheos", false, false, nil)
		if err != nil {
			log.Fatalf("Error creating PebbleDBDatabase: %v", err)
		}
		instance = &EthDbWrapper{db: db}
	})
	return instance, nil
}
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

func (p *EthDbWrapper) PutCt(h types.Hash, cipher *types.FheEncrypted) error {
	// Serialize Ciphertext
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(cipher)
	if err != nil {
		return err
	}

	// Use hash as key
	return p.db.Put(h[:], buf.Bytes())
}

func (p *EthDbWrapper) HasCt(h types.Hash) bool {
	isPresent, err := p.db.Has(h[:])
	if err != nil {
		return false
	}

	return isPresent
}

func (p *EthDbWrapper) GetCt(h types.Hash) (*types.FheEncrypted, error) {
	val, err := p.db.Get(h[:])
	if err != nil {
		return nil, err
	}

	var cipher types.FheEncrypted
	err = gob.NewDecoder(bytes.NewBuffer(val)).Decode(&cipher)
	if err != nil {
		return nil, err
	}

	return &cipher, nil
}

func (p *EthDbWrapper) DeleteCt(h types.Hash) error {
	return p.db.Delete(h[:])
}
