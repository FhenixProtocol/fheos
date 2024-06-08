package precompiles

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	"os"
	"time"
)

type FheosState struct {
	FheosVersion uint64
	Storage      storage2.FheosStorage
	//MaxUintValue *big.Int // This should contain the max value of the supported uint type
}

func (fs *FheosState) GetRandomForGasEstimation() []byte {
	return []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
	}
}

const FheosVersion = uint64(1001)

func getDbPath() string {
	dbPath := os.Getenv("FHEOS_DB_PATH")
	if dbPath == "" {
		return os.TempDir() + "/fheos"
	}

	return dbPath
}

var State *FheosState = nil

func (fs *FheosState) GetCiphertext(hash types.Hash) (*types.CipherTextRepresentation, error) {
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "db", "get")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	sharedCt, err := fs.Storage.GetCt(hash)
	if sharedCt == nil {
		return nil, err
	}

	return &sharedCt.Ciphertext, nil
}

func (fs *FheosState) DereferenceCiphertext(hash types.Hash) error {
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "db", "dereference")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}
	sharedCt, _ := fs.Storage.GetCt(hash)
	if sharedCt == nil {
		return nil
	}

	sharedCt.RefCount--
	if sharedCt.RefCount == 0 {
		logger.Info("Deleted ciphertext", "hash", hex.EncodeToString(hash[:]))
		return fs.Storage.DeleteCt(hash)
	}

	logger.Info("Decremented ciphertext ref count", "hash", hex.EncodeToString(hash[:]), "newRefCount", sharedCt.RefCount)
	return fs.Storage.PutCt(hash, sharedCt)
}

func (fs *FheosState) incrementRefCountHelper(ctHash types.Hash, sharedCt *types.SharedCiphertext) error {
	// Skip count for trivially encrypted ciphertexts
	if fhe.IsTriviallyEncryptedCtHash(ctHash) {
		return nil
	}

	sharedCt.RefCount++
	logger.Info("Incremented ciphertext ref count", "hash", hex.EncodeToString(ctHash[:]), "newRefCount", sharedCt.RefCount)
	return fs.Storage.PutCt(ctHash, sharedCt)
}

func (fs *FheosState) SetCiphertext(ctHash types.Hash, ct *types.CipherTextRepresentation) error {
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "db", "put")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	currentCt, err := fs.Storage.GetCt(ctHash)
	if err != nil {
		newCt := &types.SharedCiphertext{
			Ciphertext: *ct,
			RefCount:   1,
		}
		return fs.Storage.PutCt(ctHash, newCt)
	}

	return fs.incrementRefCountHelper(ctHash, currentCt)
}

func (fs *FheosState) ReferenceCiphertext(hash types.Hash, ct *types.SharedCiphertext) error {
	return fs.incrementRefCountHelper(hash, ct)
}

func createFheosState(storage storage2.FheosStorage, version uint64) {
	State = &FheosState{
		version,
		storage,
	}
}

func InitializeFheosState() error {
	store, err := storage2.InitStorage(getDbPath())

	if err != nil {
		logger.Error("failed to open storage for fheos State")
		return err
	}

	err = store.PutVersion(FheosVersion)
	if err != nil {
		logger.Error("failed to write version into fheos db", "err", err)
		return errors.New("failed to write version into fheos db")
	}

	createFheosState(*store, FheosVersion)

	return nil
}
