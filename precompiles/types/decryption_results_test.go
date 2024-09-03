package types

import (
	"math/big"
	"testing"
	"time"

	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	"github.com/stretchr/testify/assert"
)

func TestDecryptionResults(t *testing.T) {
	t.Run("NewDecryptionResultsMap", func(t *testing.T) {
		dr := NewDecryptionResultsMap()
		assert.NotNil(t, dr)
		assert.Empty(t, dr.data)
	})

	t.Run("CreateEmptyRecord", func(t *testing.T) {
		dr := NewDecryptionResultsMap()
		key := PendingDecryption{Hash: fhe.Hash{1, 2, 3}, Type: SealOutput}

		dr.CreateEmptyRecord(key)
		record, exists := dr.data[key]
		assert.True(t, exists)
		assert.Nil(t, record.Value)
		assert.WithinDuration(t, time.Now(), record.Timestamp, time.Second)

		// Creating again should not overwrite
		time.Sleep(time.Millisecond * 10)
		dr.CreateEmptyRecord(key)
		newRecord, _ := dr.data[key]
		assert.Equal(t, record.Timestamp, newRecord.Timestamp)
	})

	t.Run("SetValue", func(t *testing.T) {
		dr := NewDecryptionResultsMap()
		key := PendingDecryption{Hash: fhe.Hash{1, 2, 3}, Type: SealOutput}

		// Set for SealOutput
		err := dr.SetValue(key, []byte{4, 5, 6})
		assert.NoError(t, err)
		record, exists := dr.data[key]
		assert.True(t, exists)
		assert.Equal(t, []byte{4, 5, 6}, record.Value)

		// Set for Require
		keyRequire := PendingDecryption{Hash: fhe.Hash{4, 5, 6}, Type: Require}
		err = dr.SetValue(keyRequire, true)
		assert.NoError(t, err)

		// Set for Decrypt
		keyDecrypt := PendingDecryption{Hash: fhe.Hash{7, 8, 9}, Type: Decrypt}
		err = dr.SetValue(keyDecrypt, big.NewInt(123))
		assert.NoError(t, err)

		// Set with wrong type
		err = dr.SetValue(key, true)
		assert.Error(t, err)
	})

	t.Run("Get", func(t *testing.T) {
		dr := NewDecryptionResultsMap()
		key := PendingDecryption{Hash: fhe.Hash{1, 2, 3}, Type: SealOutput}

		// Get non-existent key
		record, exists := dr.Get(key)
		assert.False(t, exists)
		assert.Nil(t, record.Value)
		assert.True(t, record.Timestamp.IsZero())

		// Get empty record
		dr.CreateEmptyRecord(key)
		record, exists = dr.Get(key)
		assert.True(t, exists)
		assert.Nil(t, record.Value)
		assert.False(t, record.Timestamp.IsZero())

		// Get SealOutput
		dr.SetValue(key, []byte{4, 5, 6})
		record, exists = dr.Get(key)
		assert.True(t, exists)
		assert.Equal(t, []byte{4, 5, 6}, record.Value)
		assert.False(t, record.Timestamp.IsZero())

		// Get Require
		keyRequire := PendingDecryption{Hash: fhe.Hash{4, 5, 6}, Type: Require}
		dr.SetValue(keyRequire, true)
		record, exists = dr.Get(keyRequire)
		assert.True(t, exists)
		assert.Equal(t, true, record.Value)
		assert.False(t, record.Timestamp.IsZero())

		// Get Decrypt
		keyDecrypt := PendingDecryption{Hash: fhe.Hash{7, 8, 9}, Type: Decrypt}
		dr.SetValue(keyDecrypt, big.NewInt(123))
		record, exists = dr.Get(keyDecrypt)
		assert.True(t, exists)
		assert.Equal(t, big.NewInt(123), record.Value)
		assert.False(t, record.Timestamp.IsZero())

		// Get with wrong type
		keyWrong := PendingDecryption{Hash: fhe.Hash{10, 11, 12}, Type: PrecompileName(99)}
		dr.data[keyWrong] = DecryptionRecord{Value: "wrong", Timestamp: time.Now()}
		record, exists = dr.Get(keyWrong)
		assert.True(t, exists)
		assert.Equal(t, "wrong", record.Value)
		assert.False(t, record.Timestamp.IsZero())
	})

	t.Run("Remove", func(t *testing.T) {
		dr := NewDecryptionResultsMap()
		key := PendingDecryption{Hash: fhe.Hash{1, 2, 3}, Type: SealOutput}

		dr.SetValue(key, []byte{4, 5, 6})
		assert.Len(t, dr.data, 1)

		dr.Remove(key)
		assert.Len(t, dr.data, 0)

		// Removing non-existent key should not panic
		dr.Remove(key)
	})

	t.Run("Concurrency", func(t *testing.T) {
		dr := NewDecryptionResultsMap()
		key := PendingDecryption{Hash: fhe.Hash{1, 2, 3}, Type: SealOutput}

		done := make(chan bool)
		go func() {
			for i := 0; i < 1000; i++ {
				dr.CreateEmptyRecord(key)
				dr.SetValue(key, []byte{byte(i)})
				dr.Get(key)
			}
			done <- true
		}()

		go func() {
			for i := 0; i < 1000; i++ {
				dr.CreateEmptyRecord(key)
				dr.SetValue(key, []byte{byte(i)})
				dr.Get(key)
			}
			done <- true
		}()

		<-done
		<-done

		// No race condition should occur
	})
}
