package precompiles

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/fhenixprotocol/go-tfhe"
	"math/big"
	"time"
)

type FheosState struct {
	FheosVersion uint64
	Storage      Storage
	EZero        [][]byte // Preencrypted 0s for each uint type
	MaxUintValue *big.Int // This should contain the max value of the supported uint type
}

const FheosVersion = uint64(1)

var state *FheosState = nil

func (fs *FheosState) GetStorageSize() uint64 {
	storageSize := fs.Storage.Size()
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "db", "size")
		metrics.GetOrRegisterGauge(h, nil).Update(int64(storageSize))
	}
	return storageSize
}

func (fs *FheosState) GetCiphertext(hash tfhe.Hash) (*tfhe.Ciphertext, error) {
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "db", "get")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}
	return fs.Storage.GetCt(hash)
}

func (fs *FheosState) SetCiphertext(ct *tfhe.Ciphertext) error {
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "db", "put")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	result := fs.Storage.PutCt(ct.Hash(), ct)

	// This checks the size of the db and logs it - we don't really need the result
	if metrics.Enabled {
		_ = fs.Storage.Size()
	}

	return result
}

func createFheosState(storage *Storage, version uint64) error {
	state = &FheosState{
		version,
		*storage,
		nil,
		new(big.Int).SetUint64(^uint64(0)),
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
			logger.Error("failed to encrypt 0 for ezero", "toType", i, "err", err)
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
		logger.Error("failed to write version into fheos db", "err", err)
		return errors.New("failed to write version into fheos db")
	}

	err = createFheosState(&storage, FheosVersion)

	if err != nil {
		logger.Error("failed to create fheos state", "err", err)
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
