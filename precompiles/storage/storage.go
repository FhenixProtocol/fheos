package storage

import (
	"github.com/fhenixprotocol/fheos/precompiles/storage/leveldb"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/go-tfhe"
	"os"
)

type Storage interface {
	Put(t types.DataType, key []byte, val []byte) error
	Get(t types.DataType, key []byte) ([]byte, error)
	GetVersion() (uint64, error)
	PutVersion(v uint64) error
	PutCt(h tfhe.Hash, cipher *tfhe.Ciphertext) error
	GetCt(h tfhe.Hash) (*tfhe.Ciphertext, error)
	Size() uint64
}

const DBPath = "/home/user/fhenix/fheosdb"

func getDbPath() string {
	dbPath := os.Getenv("FHEOS_DB_PATH")
	if dbPath == "" {
		return DBPath
	}

	return dbPath
}

func InitStorage() Storage {
	storage := leveldb.LevelDbStorage{
		DbPath: getDbPath(),
	}

	return storage
}
