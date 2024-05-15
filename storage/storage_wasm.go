package storage

import (
	"github.com/fhenixprotocol/fheos/precompiles/types"
)

type FheosStorage struct {
}

type IFheosStorage interface {
	//Put(t types.DataType, key []byte, val []byte) error
	//Get(t types.DataType, key []byte) ([]byte, error)
	GetVersion() (uint64, error)
	PutVersion(v uint64) error
	PutCt(h types.Hash, cipher *types.FheEncrypted) error
	GetCt(h types.Hash) (*types.FheEncrypted, error)
}

func InitStorage(_ string) (*FheosStorage, error) {
	return &FheosStorage{}, nil
}

//func (store WasmStorage) Put(_ types.DataType, _ []byte, _ []byte) error {
//	return nil
//}
//
//func (store WasmStorage) Get(t types.DataType, key []byte) ([]byte, error) {
//	return nil, nil
//}

func (store FheosStorage) GetVersion() (uint64, error) {
	return 0, nil
}

func (store FheosStorage) PutVersion(v uint64) error {
	return nil
}

func (store FheosStorage) PutCt(h types.Hash, cipher *types.CipherTextRepresentation) error {
	return nil
}

func (store FheosStorage) GetCt(h types.Hash) (*types.CipherTextRepresentation, error) {
	return nil, nil
}

func (store FheosStorage) Size() uint64 { return 0 }
