package precompiles

import (
	"github.com/fhenixprotocol/go-tfhe"
)

type WasmStorage struct {
}

func InitStorage(burner GasBurner) Storage {
	return WasmStorage{}
}

func (store WasmStorage) Put(_ DataType, _ []byte, _ []byte, _ bool) error {
	return nil
}

func (store WasmStorage) Get(t DataType, key []byte, _ bool) ([]byte, error) {
	return nil, nil
}

func (store WasmStorage) GetVersion() (uint64, error) {
	return 0, nil
}

func (store WasmStorage) PutVersion(v uint64) error {
	return nil
}

func (store WasmStorage) PutCt(h tfhe.Hash, cipher *tfhe.Ciphertext, _ bool) error {
	return nil
}

func (store WasmStorage) GetCt(h tfhe.Hash, _ bool) (*tfhe.Ciphertext, error) {
	return nil, nil
}
