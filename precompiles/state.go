package precompiles

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	"math/big"
	"os"
	"time"
)

type FheosState struct {
	FheosVersion uint64
	Storage      storage2.FheosStorage
	EZero        [][]byte // Preencrypted 0s for each uint type
	//MaxUintValue *big.Int // This should contain the max value of the supported uint type
}

func (fs *FheosState) GetZero(uintType fhe.EncryptionType) *fhe.FheEncrypted {
	// todo (eshel): allow for other security zones?
	ct, err := fhe.EncryptPlainText(*big.NewInt(0), uintType, 0)
	if err != nil {
		return nil
	}
	return ct
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
	return fs.Storage.GetCt(hash)
}

func (fs *FheosState) SetCiphertext(ct *types.CipherTextRepresentation) error {
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "db", "put")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	result := fs.Storage.PutCt(types.Hash((*fhe.FheEncrypted)(ct.Data).Hash()), ct)

	return result
}

func createFheosState(storage storage2.FheosStorage, version uint64) error {
	State = &FheosState{
		version,
		storage,
		nil,
	}

	tempTp := TxParams{
		Commit:        false,
		GasEstimation: false,
		EthCall:       true,
	}

	// todo: refactor this - currently it causes crashing and weirdness if you try to add a new type
	// also it's not very flexible
	zero := make([]byte, 32)
	var err error
	ezero := make([][]byte, 14)

	for i := 0; i < 6; i++ {
		// todo (eshel): allow for other security zones?
		ezero[i], _, err = TrivialEncrypt(zero, byte(i), 0, &tempTp)
		if err != nil {
			logger.Error("failed to encrypt 0 for ezero", "toType", i, "err", err)
			// don't error out - this should be handled dynamically later - otherwise it just requires the backend to do work
			// that might not be necessary right now, and makes it more annoying for unit tests that might not need to encrypt
		}
	}

	// todo (eshel): allow for other security zones?
	ezero[13], _, err = TrivialEncrypt(zero, byte(13), 0, &tempTp)
	if err != nil {
		logger.Error("failed to encrypt 0 for ezero", "toType", 13, "err", err)
	}

	State.EZero = ezero

	return nil
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

	err = createFheosState(*store, FheosVersion)

	if err != nil {
		logger.Error("failed to create fheos State", "err", err)
		return errors.New("failed to create fheos State")
	}

	return nil
}
