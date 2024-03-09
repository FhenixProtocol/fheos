//go:build amd64 || arm64

package storage

import (
	"github.com/fhenixprotocol/fheos/precompiles/storage/pebble"
	"github.com/fhenixprotocol/fheos/precompiles/types"
)

type Storage interface {
	Put(t types.DataType, key []byte, val []byte) error
	Get(t types.DataType, key []byte) ([]byte, error)
	GetVersion() (uint64, error)
	PutVersion(v uint64) error
	PutCt(h fhe.Hash, cipher *fhe.FheEncrypted) error
	GetCt(h fhe.Hash) (*fhe.FheEncrypted, error)
}

func InitStorage(path string) Storage {
	// in theory we could add a db selector here... in practice we really don't want to since one is really better than the
	// other
	//storage := leveldb.NewStorage(path)
	storage := pebble.NewStorage(path)
	return storage
}
