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
	RandomCounter  uint64
	DecryptResults *types.DecryptionResults
	//MaxUintValue *big.Int // This should contain the max value of the supported uint type
}

func (fs *FheosState) GetEmptyKeyForGasEstimation() []byte {
	return types.SerializeCiphertextKey(types.GetEmptyCiphertextKey())
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

func (fs *FheosState) GetCiphertext(hash types.Hash) (*types.FheEncrypted, error) {
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

func (fs *FheosState) SetCiphertext(ct *types.FheEncrypted) error {
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "db", "put")
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	result := fs.Storage.PutCt(types.Hash((*fhe.FheEncrypted)(ct).GetHash()), ct)

	return result
}

func (fs *FheosState) GetRandomCounter() uint64 {
	return fs.RandomCounter
}

func (fs *FheosState) IncRandomCounter() uint64 {
	fs.RandomCounter += 1
	return fs.RandomCounter
}

func createFheosState(storage storage2.FheosStorage, version uint64) {
	State = &FheosState{
		version,
		storage,
		0,
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

	logger.Debug("Loading resolved decryptions", "numDecryptions", numDecryptions)

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
