package types

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

type PendingDecryption struct {
	Hash fhe.Hash
	Type PrecompileName
}

type DecryptionRecord struct {
	Value     any
	Timestamp time.Time
}

type DecryptionResults struct {
	data map[PendingDecryption]DecryptionRecord
	mu   sync.RWMutex
}

func NewDecryptionResultsMap() *DecryptionResults {
	return &DecryptionResults{
		data: make(map[PendingDecryption]DecryptionRecord),
	}
}

// CreateEmptyRecord creates a new empty record for the given PendingDecryption key
// if it doesn't already exist in the DecryptionResults map.
// This function is intended to be used once a parallel decryption is initiated,
// before we have a result. An empty record indicates that this decryption is pending.
// The new record will have a nil Value and the current timestamp.
// This method is thread-safe.
func (dr *DecryptionResults) CreateEmptyRecord(key PendingDecryption) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	if _, exists := dr.data[key]; !exists {
		dr.data[key] = DecryptionRecord{Value: nil, Timestamp: time.Now()}
	}
}

func (dr *DecryptionResults) SetValue(key PendingDecryption, value any) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	switch key.Type {
	case SealOutput:
		if _, ok := value.([]byte); !ok {
			return fmt.Errorf("value for SealOutput must be []byte")
		}
	case Require:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("value for Require must be bool")
		}
	case Decrypt:
		if _, ok := value.(*big.Int); !ok {
			return fmt.Errorf("value for Decrypt must be *big.Int")
		}
	default:
		return fmt.Errorf("unknown PrecompileName")
	}

	dr.data[key] = DecryptionRecord{Value: value, Timestamp: time.Now()}
	return nil
}

// SetRecord is just like setValue but sets the complete record, including timestamp
// This way timestamps could be synchronized between different nodes
func (dr *DecryptionResults) SetRecord(key PendingDecryption, record DecryptionRecord) error {
	switch key.Type {
	case SealOutput:
		if _, ok := record.Value.([]byte); !ok {
			return fmt.Errorf("value for SealOutput must be []byte")
		}
	case Require:
		if _, ok := record.Value.(bool); !ok {
			return fmt.Errorf("value for Require must be bool")
		}
	case Decrypt:
		if _, ok := record.Value.(*big.Int); !ok {
			return fmt.Errorf("value for Decrypt must be *big.Int")
		}
	default:
		return fmt.Errorf("unknown PrecompileName")
	}

	dr.mu.Lock()
	defer dr.mu.Unlock()

	dr.data[key] = record
	return nil
}

func (dr *DecryptionResults) Get(key PendingDecryption) (DecryptionRecord, bool) { //(any, bool, time.Time, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	record, exists := dr.data[key]
	if !exists {
		return DecryptionRecord{}, false
	}

	return record, true
}

func (dr *DecryptionResults) Remove(key PendingDecryption) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	delete(dr.data, key)
}
