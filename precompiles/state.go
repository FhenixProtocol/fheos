package precompiles

import (
	"errors"
	"github.com/fhenixprotocol/go-tfhe"
)

type FheosState struct {
	FheosVersion uint64
	Storage      Storage
}

const FheosVersion = uint64(1)

func (fs *FheosState) GetCiphertext(hash tfhe.Hash, isTx bool) (*tfhe.Ciphertext, error) {
	return fs.Storage.GetCt(hash, isTx)
}

func (fs *FheosState) SetCiphertext(ct *tfhe.Ciphertext, isTx bool) error {
	return fs.Storage.PutCt(ct.Hash(), ct, isTx)
}

func InitializeFheosState(burner GasBurner) (*FheosState, error) {
	storage := InitStorage(burner)

	if storage == nil {
		logger.Error("failed to open storage for fheos state")
		return nil, errors.New("failed to open storage for fheos state")
	}

	err := storage.PutVersion(FheosVersion)
	if err != nil {
		logger.Error("failed to write version into fheos db ", err)
		return nil, errors.New("failed to write version into fheos db ")
	}

	return &FheosState{
		FheosVersion,
		storage,
	}, nil

}

func OpenFheosState(burner GasBurner) (*FheosState, error) {
	storage := InitStorage(burner)
	version, err := storage.GetVersion()
	if err != nil {
		logger.Error("failed to read version from fheos db ", err)
		return nil, err
	}

	if version != FheosVersion {
		logger.Error("fheos version is corrupted")
		return nil, errors.New("fheos version is corrupted")
	}

	return &FheosState{
		version,
		storage,
	}, nil
}

// The following functions are useful for future implementation of storage based on geth
// The reason we are not using them now is because they store the entire map as bytes on every set
// which is not efficient for large maps - not even for a few ciphertexts.
//func InitializeFheosState(stateDB vm.StateDB, burner GasBurner) (*FheosState, error) {
//	storage := NewStorage(stateDB, burner)
//	fheosVersion, err := storage.GetUint64ByUint64(versionOffset)
//	if err != nil {
//		return nil, err
//	}
//	if fheosVersion != 0 {
//		return nil, errors.New("fheos state is already initialized")
//	}
//
//	_ = storage.SetUint64ByUint64(versionOffset, 1)
//
//	b, err := encodeMap(make(map[tfhe.Hash]tfhe.Ciphertext))
//	if err != nil {
//		return nil, err
//	}
//
//	ctStorage := storage.OpenBytesStorage([]byte{ctOffset})
//	_ = ctStorage.Set(b)
//
//	aState, err := OpenFheosState(stateDB, burner)
//	if err != nil {
//		return nil, err
//	}
//
//	return aState, nil
//}
//
//func OpenFheosState(stateDB vm.StateDB, burner GasBurner) (*FheosState, error) {
//	storage := NewStorage(stateDB, burner)
//	fheosVersion, err := storage.GetUint64ByUint64(versionOffset)
//	if err != nil {
//		return nil, err
//	}
//	if fheosVersion == 0 {
//		return nil, errors.New("fheos state is uninitialized")
//	}
//	return &FheosState{
//		fheosVersion,
//		storage.OpenBytesStorage([]byte{ctOffset}),
//		burner,
//	}, nil
//}
