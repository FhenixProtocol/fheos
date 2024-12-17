package storage_test

import (
	"bytes"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"sync"
	"testing"

	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

const storagePath = "/tmp/fheosdb"
const TestFileSize = 1024 * 1024 * 4 // 4MB
func init() {
	err := os.Setenv("FHEOS_DB_PATH", "/tmp/fheosdb")
	if err != nil {
		panic(err)
	}

}

// Helper function to generate a random Ciphertext
func randomCiphertext() *fhe.FheEncrypted {
	// Generate a large serialization
	serialization := make([]byte, TestFileSize)
	_, _ = rand.Read(serialization)

	// Generate a random hash
	hash := make([]byte, 32)
	_, _ = rand.Read(hash)

	ct := &fhe.FheEncrypted{
		Data:       serialization,
		Compact:    true,
		Compressed: true,
		UintType:   2, // Assuming 0 is a valid UintType
	}

	ct.GetHash()

	return ct
}

func TestStorageConcurrency(t *testing.T) {
	storage, err := storage2.InitStorage(storagePath)

	if err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}

	var wg sync.WaitGroup

	// Define the number of concurrent operations
	concurrencyLevel := 10

	wg.Add(concurrencyLevel) // * 2 for both Put and PutCt operations
	ct := randomCiphertext()
	// Test concurrent PutCt operations
	for i := 0; i < concurrencyLevel; i++ {
		go func(i int) {
			defer wg.Done()

			if err := storage.PutCt(types.Hash(fhe.Hash{byte(i)}), (*types.FheEncrypted)(ct)); err != nil {
				t.Errorf("Failed to put ciphertext: %v", err)
			}
		}(i)
	}

	wg.Wait() // Wait for all goroutines to finish
}

func TestStorageEphemeralConcurrency(t *testing.T) {
	storage, err := storage2.InitStorage(storagePath)

	if err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}

	var wg sync.WaitGroup

	// Define the number of concurrent operations
	concurrencyLevel := 10
	ct := randomCiphertext()

	wg.Add(concurrencyLevel)

	//storage.SetEphemeral()

	// Test concurrent PutCt operations
	for i := 0; i < concurrencyLevel; i++ {
		go func(i int) {
			defer wg.Done()

			if err := storage.PutCt(types.Hash(fhe.Hash{byte(i)}), (*types.FheEncrypted)(ct)); err != nil {
				t.Errorf("Failed to put ciphertext: %v", err)
			}
		}(i)
	}

	wg.Wait() // Wait for all goroutines to finish
}

// Add more tests and benchmarks as needed for GetCt, Size, etc.
// Additional tests for the Storage interface

func TestStorageCt(t *testing.T) {
	storage, err := storage2.InitStorage(storagePath)

	if err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}

	ct := randomCiphertext()
	hash := fhe.Hash{0} // Simplified hash for testing

	if err := storage.PutCt(types.Hash(hash), (*types.FheEncrypted)(ct)); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := storage.GetCt(types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	if !bytes.Equal(retrievedCt.Data, ct.Data) || !((*fhe.FheEncrypted)(retrievedCt).GetHash() == ct.GetHash()) {
		t.Errorf("Retrieved ciphertext does not match the original")
	}
	assert.Equal(t, ct.Compact, retrievedCt.Compact)
	assert.Equal(t, ct.Compressed, retrievedCt.Compressed)
}

func TestStorageEphemeralCt(t *testing.T) {
	storage, err := storage2.InitStorage(storagePath)

	if err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}

	ct := randomCiphertext()
	hash := fhe.Hash{0} // Simplified hash for testing

	if err := storage.PutCt(types.Hash(hash), (*types.FheEncrypted)(ct)); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := storage.GetCt(types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	if !bytes.Equal(retrievedCt.Data, ct.Data) || !((*fhe.FheEncrypted)(retrievedCt).GetHash() == ct.GetHash()) {
		t.Errorf("Retrieved ciphertext does not match the original")
	}
}

func TestStorageVersioning(t *testing.T) {
	storage, err := storage2.InitStorage(storagePath)

	if err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}

	initialVersion, err := storage.GetVersion()
	if err != nil {
		t.Fatalf("Failed to get initial version: %v", err)
	}

	newVersion := initialVersion + 1
	if err := storage.PutVersion(newVersion); err != nil {
		t.Fatalf("Failed to update version: %v", err)
	}

	updatedVersion, err := storage.GetVersion()
	if err != nil {
		t.Fatalf("Failed to get updated version: %v", err)
	}

	if updatedVersion != newVersion {
		t.Errorf("Version mismatch: expected %d, got %d", newVersion, updatedVersion)
	}
}

func TestStorageGetSetReset(t *testing.T) {
	storage, err := storage2.InitStorage(storagePath)

	if err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}

	ct := randomCiphertext()
	hash := fhe.Hash{0} // Simplified hash for testing

	if err := storage.PutCt(types.Hash(hash), (*types.FheEncrypted)(ct)); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := storage.GetCt(types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	//storage.SetEphemeral()

	retrievedCt, err = storage.GetCt(types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to get ciphertext from ephemeral (should have fallen back to disk storage): %v", err)
	}

	if !bytes.Equal(retrievedCt.Data, ct.Data) || !((*fhe.FheEncrypted)(retrievedCt).GetHash() == ct.GetHash()) {
		t.Errorf("Retrieved ciphertext does not match the original")
	}

	// todo: when delete is a thing we should test it here and clear the storage
}

func BenchmarkConcurrentPut(b *testing.B) {
	storage, err := storage2.InitStorage(storagePath)

	if err != nil {
		b.Fatalf("Failed to initialize storage: %v", err)
	}

	// Number of goroutines to use for concurrent writes
	concurrencyLevel := 10
	ct := randomCiphertext()

	b.ResetTimer() // Reset the timer to exclude setup time

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(concurrencyLevel)

		for j := 0; j < concurrencyLevel; j++ {
			go func(j int) {
				defer wg.Done()

				if err := storage.PutCt(types.Hash(fhe.Hash{byte(i)}), (*types.FheEncrypted)(ct)); err != nil {
					b.Errorf("Failed to put ciphertext: %v", err)
				}
			}(j)
		}

		wg.Wait() // Wait for all goroutines to finish
	}
}
func TestMultiStore_AppendCt(t *testing.T) {
	diskStorage, err := storage2.InitStorage(storagePath)
	multiStore := storage2.NewMultiStore(nil, diskStorage)

	ct := randomCiphertext()
	hash := fhe.Hash{100} // this key needs to be unique for the test

	if err := multiStore.PutCtIfNotExist(types.Hash(hash), (*types.FheEncrypted)(ct)); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := multiStore.GetCt(types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	assert.Equal(t, (*fhe.FheEncrypted)(retrievedCt).GetHash(), ct.GetHash())

	if !bytes.Equal(retrievedCt.Data, ct.Data) {
		t.Errorf("Retrieved ciphertext does not match the original")
	}
}
func TestMultiStore_AppendCtPlaceholderReplace(t *testing.T) {
	diskStorage, err := storage2.InitStorage(storagePath)
	multiStore := storage2.NewMultiStore(nil, diskStorage)

	ct := randomCiphertext()
	ct.Placeholder = true
	hash := fhe.Hash{104} // this key needs to be unique for the test

	if err := multiStore.PutCtIfNotExist(types.Hash(hash), (*types.FheEncrypted)(ct)); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	ct.Placeholder = false

	if err := multiStore.PutCtIfNotExist(types.Hash(hash), (*types.FheEncrypted)(ct)); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := multiStore.GetCt(types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	assert.Equal(t, (*fhe.FheEncrypted)(retrievedCt).GetHash(), ct.GetHash())
	assert.Equal(t, retrievedCt.Placeholder, false)
}
