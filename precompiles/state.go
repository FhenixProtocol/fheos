package precompiles

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	"io"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/metrics"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
)

type FheosState struct {
	FheosVersion   uint64
	Storage        storage2.FheosStorage
	DecryptResults *types.DecryptionResults
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

func createFheosState(storage storage2.FheosStorage, version uint64) {
	State = &FheosState{
		version,
		storage,
		types.NewDecryptionResultsMap(),
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

func GetSerializedDecryptionResult(key types.PendingDecryption) ([]byte, error) {
	if State == nil {
		return nil, errors.New("fheos state is not initialized")
	}

	if State.DecryptResults == nil {
		return nil, errors.New("DecryptionResults is not initialized in fheos state")
	}

	return State.DecryptResults.GetSerializedDecryptionResult(key)
}

func LoadMultipleResolvedDecryptions(reader io.Reader) error {
	// parse the number of resolved decryptions
	var numDecryptions int32
	err := binary.Read(reader, binary.LittleEndian, &numDecryptions)
	if err != nil {
		return err
	}

	for i := int32(0); i < numDecryptions; i++ {
		err = LoadResolvedDecryption(reader)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadResolvedDecryption(reader io.Reader) error {
	if State == nil {
		return errors.New("fheos state is not initialized")
	}

	if State.DecryptResults == nil {
		return errors.New("fheos state is not initialized")
	}

	return State.DecryptResults.LoadResolvedDecryption(reader)
}
