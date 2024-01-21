package precompiles

import (
	"errors"
	"github.com/fhenixprotocol/go-tfhe"
)

type FheosState struct {
	FheosVersion uint64
	Storage      Storage
	EZero        [][]byte // Preencrypted 0s for each uint type
}

const FheosVersion = uint64(1)

var state *FheosState = nil

func (fs *FheosState) GetCiphertext(hash tfhe.Hash) (*tfhe.Ciphertext, error) {
	return fs.Storage.GetCt(hash)
}

func (fs *FheosState) SetCiphertext(ct *tfhe.Ciphertext) error {
	return fs.Storage.PutCt(ct.Hash(), ct)
}

func createFheosState(storage *Storage, version uint64) error {
	state = &FheosState{
		version,
		*storage,
		nil,
	}

	tempTp := TxParams{
		Commit:        false,
		GasEstimation: false,
		EthCall:       true,
	}

	zero := make([]byte, 32)
	var err error
	ezero := make([][]byte, 3)

	for i := 0; i < 3; i++ {
		ezero[i], _, err = TrivialEncrypt(zero, byte(i), &tempTp)
		if err != nil {
			logger.Error("failed to encrypt 0 for ezero ", i, err)
			return err
		}
	}

	state.EZero = ezero

	return nil
}

func InitializeFheosState() error {
	storage := InitStorage()

	if storage == nil {
		logger.Error("failed to open storage for fheos state")
		return errors.New("failed to open storage for fheos state")
	}

	err := storage.PutVersion(FheosVersion)
	if err != nil {
		logger.Error("failed to write version into fheos db ", err)
		return errors.New("failed to write version into fheos db")
	}

	err = createFheosState(&storage, FheosVersion)

	if err != nil {
		logger.Error("failed to create fheos state ", err)
		return errors.New("failed to create fheos state")
	}

	return nil
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
