package precompiles

import (
	"github.com/fhenixprotocol/go-tfhe"
)

type WasmStorage struct {
}

func InitStorage(burner GasBurner) Storage {
	return WasmStorage{}
}

func (store WasmStorage) Put(_ DataType, _ []byte, _ []byte) error {
	return nil
}

func (store WasmStorage) Get(t DataType, key []byte) ([]byte, error) {
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
