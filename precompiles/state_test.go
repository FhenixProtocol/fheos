package precompiles

import (
	"bytes"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"

	"github.com/fhenixprotocol/fheos/precompiles/types"
)

func TestDecryptionResultsSerialization(t *testing.T) {
	t.Run("Serialize Deserialize", func(t *testing.T) {
		dr := types.NewDecryptionResultsMap()
		dr2 := types.NewDecryptionResultsMap()

		// Set SealOutput
		key := types.PendingDecryption{Hash: fhe.Hash{1, 2, 3}, Type: types.SealOutput}
		dr.SetValue(key, []byte{4, 5, 6})
		drSerialized, err := dr.GetSerializedDecryptionResult(key)
		assert.NoError(t, err)

		err = dr2.LoadResolvedDecryption(bytes.NewReader(drSerialized))
		assert.NoError(t, err)

		record2, ok := dr2.Get(key)
		assert.True(t, ok)

		expected, ok := dr.Get(key)
		assert.Equal(t, expected.Value, record2.Value)
		assert.True(t, expected.Timestamp.Equal(record2.Timestamp))

		//Set Require
		keyRequire := types.PendingDecryption{Hash: fhe.Hash{4, 5, 6}, Type: types.Require}
		dr.SetValue(keyRequire, true)
		drSerialized, err = dr.GetSerializedDecryptionResult(keyRequire)
		assert.NoError(t, err)

		err = dr2.LoadResolvedDecryption(bytes.NewReader(drSerialized))
		assert.NoError(t, err)

		record2, ok = dr2.Get(keyRequire)
		assert.True(t, ok)

		expected, ok = dr.Get(keyRequire)
		assert.Equal(t, expected.Value, record2.Value)
		assert.True(t, expected.Timestamp.Equal(record2.Timestamp))

		// Set Decrypt
		keyDecrypt := types.PendingDecryption{Hash: fhe.Hash{7, 8, 9}, Type: types.Decrypt}
		dr.SetValue(keyDecrypt, big.NewInt(123))
		drSerialized, err = dr.GetSerializedDecryptionResult(keyDecrypt)
		assert.NoError(t, err)

		err = dr2.LoadResolvedDecryption(bytes.NewReader(drSerialized))
		assert.NoError(t, err)

		record2, ok = dr2.Get(keyDecrypt)
		assert.True(t, ok)

		expected, ok = dr.Get(keyDecrypt)
		assert.Equal(t, expected.Value, record2.Value)
		assert.True(t, expected.Timestamp.Equal(record2.Timestamp))
	})
}
