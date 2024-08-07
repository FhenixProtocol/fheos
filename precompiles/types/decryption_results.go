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
	Value     interface{}
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

func (dr *DecryptionResults) CreateEmptyRecord(key PendingDecryption) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	if _, exists := dr.data[key]; !exists {
		dr.data[key] = DecryptionRecord{Value: nil, Timestamp: time.Now()}
	}
}

func (dr *DecryptionResults) SetValue(key PendingDecryption, value interface{}) error {
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

func (dr *DecryptionResults) Get(key PendingDecryption) (interface{}, bool, time.Time, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	record, exists := dr.data[key]
	if !exists {
		return nil, false, time.Time{}, nil
	}

	if record.Value == nil {
		return nil, true, record.Timestamp, nil // Exists but no value
	}

	switch key.Type {
	case SealOutput:
		if bytes, ok := record.Value.([]byte); ok {
			return bytes, true, record.Timestamp, nil
		}
		return nil, true, record.Timestamp, fmt.Errorf("value is not []byte as expected for SealOutput")
	case Require:
		if boolValue, ok := record.Value.(bool); ok {
			return boolValue, true, record.Timestamp, nil
		}
		return nil, true, record.Timestamp, fmt.Errorf("value is not bool as expected for Require")
	case Decrypt:
		if bigInt, ok := record.Value.(*big.Int); ok {
			return bigInt, true, record.Timestamp, nil
		}
		return nil, true, record.Timestamp, fmt.Errorf("value is not *big.Int as expected for Decrypt")
	default:
		return nil, true, record.Timestamp, fmt.Errorf("unknown PrecompileName")
	}
}

func (dr *DecryptionResults) Remove(key PendingDecryption) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	delete(dr.data, key)
}
