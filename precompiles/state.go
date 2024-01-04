package precompiles

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/fhenixprotocol/go-tfhe"
)

const (
	versionOffset = iota + 30
	ctOffset
)

type FheosState struct {
	FheosVersion uint64
	CtStorage    BytesStorage
	Burner       GasBurner
}

func encodeMap(m map[tfhe.Hash]tfhe.Ciphertext) ([]byte, error) {
	b := new(bytes.Buffer)
	encoder := gob.NewEncoder(b)

	err := encoder.Encode(m)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (fs *FheosState) GetCiphertextMap() (map[tfhe.Hash]tfhe.Ciphertext, error) {
	ctMap := make(map[tfhe.Hash]tfhe.Ciphertext)
	serialized, err := fs.CtStorage.Get()

	if err != nil {
		return ctMap, err
	}

	decoder := gob.NewDecoder(bytes.NewReader(serialized))
	err = decoder.Decode(&ctMap)
	if err != nil {
		return ctMap, err
	}

	return ctMap, nil
}

func (fs *FheosState) SetCiphertextMap(ctMap *map[tfhe.Hash]tfhe.Ciphertext) error {
	encoded, err := encodeMap(*ctMap)
	if err != nil {
		return err
	}

	return fs.CtStorage.Set(encoded)
}

func (fs *FheosState) GetCiphertext(hash tfhe.Hash) (*tfhe.Ciphertext, error) {
	ctMap, err := fs.GetCiphertextMap()
	if err != nil {
		return nil, err
	}

	ct, ok := ctMap[hash]
	if !ok {
		return nil, nil
	}

	return &ct, nil
}

func (fs *FheosState) SetCiphertext(ct *tfhe.Ciphertext) error {
	ctMap, err := fs.GetCiphertextMap()
	if err != nil {
		return err
	}

	_, ok := ctMap[ct.Hash()]
	if ok {
		return nil
	}

	ctMap[ct.Hash()] = *ct
	return fs.SetCiphertextMap(&ctMap)
}

func InitializeFheosState(stateDB vm.StateDB, burner GasBurner) (*FheosState, error) {
	storage := NewStorage(stateDB, burner)
	fheosVersion, err := storage.GetUint64ByUint64(versionOffset)
	if err != nil {
		return nil, err
	}
	if fheosVersion != 0 {
		return nil, errors.New("fheos state is already initialized")
	}

	_ = storage.SetUint64ByUint64(versionOffset, 1)

	b, err := encodeMap(make(map[tfhe.Hash]tfhe.Ciphertext))
	if err != nil {
		return nil, err
	}

	ctStorage := storage.OpenBytesStorage([]byte{ctOffset})
	_ = ctStorage.Set(b)

	aState, err := OpenFheosState(stateDB, burner)
	if err != nil {
		return nil, err
	}

	return aState, nil
}

func OpenFheosState(stateDB vm.StateDB, burner GasBurner) (*FheosState, error) {
	storage := NewStorage(stateDB, burner)
	fheosVersion, err := storage.GetUint64ByUint64(versionOffset)
	if err != nil {
		return nil, err
	}
	if fheosVersion == 0 {
		return nil, errors.New("fheos state is uninitialized")
	}
	return &FheosState{
		fheosVersion,
		storage.OpenBytesStorage([]byte{ctOffset}),
		burner,
	}, nil
}
