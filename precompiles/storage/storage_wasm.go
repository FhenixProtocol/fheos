package storage

import (
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/go-tfhe"
)

type WasmStorage struct {
}

type Storage interface {
	Put(t types.DataType, key []byte, val []byte) error
	Get(t types.DataType, key []byte) ([]byte, error)
	GetVersion() (uint64, error)
	PutVersion(v uint64) error
	PutCt(h tfhe.Hash, cipher *tfhe.Ciphertext) error
	GetCt(h tfhe.Hash) (*tfhe.Ciphertext, error)
}

func InitStorage(_ string) Storage {
	return WasmStorage{}
}

func (store WasmStorage) Put(_ types.DataType, _ []byte, _ []byte) error {
	return nil
}

func (store WasmStorage) Get(t types.DataType, key []byte) ([]byte, error) {
	return nil, nil
}

func (store WasmStorage) GetVersion() (uint64, error) {
	return 0, nil
}

func (store WasmStorage) PutVersion(v uint64) error {
	return nil
}

func (store WasmStorage) PutCt(h tfhe.Hash, cipher *tfhe.Ciphertext) error {
	return nil
}

func (store WasmStorage) GetCt(h tfhe.Hash) (*tfhe.Ciphertext, error) {
	return nil, nil
}

func (store WasmStorage) Size() uint64 { return 0 }
