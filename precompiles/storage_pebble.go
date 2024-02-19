package precompiles

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/fhenixprotocol/go-tfhe"
	"log"
	"sync"

	"github.com/cockroachdb/pebble"
)

var (
	instance *PebbleStorage
	once     sync.Once
)

type PebbleStorage struct {
	db *pebble.DB
}

// NewPebbleStorage ensures a single PebbleStorage instance
func NewPebbleStorage() Storage {
	once.Do(func() {
		path := getDbPath()
		db, err := pebble.Open(path, &pebble.Options{})
		if err != nil {
			// Consider handling this error more gracefully instead of using log.Fatalf
			// For example, you could return an error to the caller.
			log.Fatalf("failed to open pebble db: %v", err)
		}
		instance = &PebbleStorage{db: db}
	})
	return instance
}
func (p *PebbleStorage) Put(t DataType, key []byte, val []byte) error {
	// Use DataType as part of the key to differentiate data types
	prefixedKey := append([]byte(fmt.Sprintf("%d_", t)), key...)
	return p.db.Set(prefixedKey, val, pebble.Sync)
}

func (p *PebbleStorage) Get(t DataType, key []byte) ([]byte, error) {
	prefixedKey := append([]byte(fmt.Sprintf("%d_", t)), key...)
	val, closer, err := p.db.Get(prefixedKey)
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	// Make a copy of the data since val becomes invalid after the closer is called
	valCopy := make([]byte, len(val))
	copy(valCopy, val)

	return valCopy, nil
}

func (p *PebbleStorage) GetVersion() (uint64, error) {
	key := []byte("version")
	val, closer, err := p.db.Get(key)
	if err != nil {
		return 0, err
	}
	defer closer.Close()

	// Assuming the version is stored as a uint64
	var version uint64
	buf := bytes.NewBuffer(val)
	err = gob.NewDecoder(buf).Decode(&version)
	if err != nil {
		return 0, err
	}

	return version, nil
}

func (p *PebbleStorage) PutVersion(v uint64) error {
	key := []byte("version")
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(v)
	if err != nil {
		return err
	}

	return p.db.Set(key, buf.Bytes(), pebble.Sync)
}

func (p *PebbleStorage) PutCt(h tfhe.Hash, cipher *tfhe.Ciphertext) error {
	// Serialize Ciphertext
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(cipher)
	if err != nil {
		return err
	}

	// Use hash as key
	return p.db.Set(h[:], buf.Bytes(), pebble.Sync)
}

func (p *PebbleStorage) GetCt(h tfhe.Hash) (*tfhe.Ciphertext, error) {
	val, closer, err := p.db.Get(h[:])
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	var cipher tfhe.Ciphertext
	err = gob.NewDecoder(bytes.NewBuffer(val)).Decode(&cipher)
	if err != nil {
		return nil, err
	}

	return &cipher, nil
}
