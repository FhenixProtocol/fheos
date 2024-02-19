package storage_test

import (
	"bytes"
	"fmt"
	storage2 "github.com/fhenixprotocol/fheos/precompiles/storage"
	"math/rand"
	"os"
	"sync"
	"testing"

	"github.com/fhenixprotocol/fheos/precompiles"
	"github.com/fhenixprotocol/go-tfhe" // Update with the actual package path
)

const storagePath = "/tmp/fheosdb"
const TestFileSize = 1024 * 1024 * 4 // 4MB

func generateKeys() error {
	if _, err := os.Stat("./keys/"); os.IsNotExist(err) {
		err := os.Mkdir("./keys/", 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat("./keys/tfhe/"); os.IsNotExist(err) {
		err := os.Mkdir("./keys/tfhe/", 0755)
		if err != nil {
			return err
		}
	}

	err :=
		tfhe.GenerateFheKeys("./keys/tfhe/", "./sks", "./cks", "./pks")
	if err != nil {
		return fmt.Errorf("error from tfhe GenerateFheKeys: %s", err)
	}
	return nil
}

func init() {
	err := generateKeys()
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	err = os.Setenv("FHEOS_DB_PATH", "/tmp/fheosdb")
	if err != nil {
		panic(err)
	}

	err = precompiles.InitFheos(&tfhe.ConfigDefault)
	if err != nil {
		panic(err)
	}

}

// Helper function to generate a random Ciphertext
func randomCiphertext() *tfhe.Ciphertext {
	// Generate a large serialization
	serialization := make([]byte, TestFileSize)
	_, _ = rand.Read(serialization)

	// Generate a random hash
	hash := make([]byte, 32)
	_, _ = rand.Read(hash)

	ct := &tfhe.Ciphertext{
		Serialization: serialization,
		UintType:      2, // Assuming 0 is a valid UintType
	}

	ct.Hash()

	return ct
}

func TestStorageConcurrency(t *testing.T) {
	storage := storage2.InitStorage(storagePath)
	var wg sync.WaitGroup

	// Define the number of concurrent operations
	concurrencyLevel := 10

	wg.Add(concurrencyLevel * 2) // * 2 for both Put and PutCt operations

	// Test concurrent Put operations
	for i := 0; i < concurrencyLevel; i++ {
		go func(i int) {
			defer wg.Done()
			key := []byte{byte(i)}
			val := []byte("test value")
			if err := storage.Put(0, key, val); err != nil {
				t.Errorf("Failed to put data: %v", err)
			}
		}(i)
	}

	// Test concurrent PutCt operations
	for i := 0; i < concurrencyLevel; i++ {
		go func(i int) {
			defer wg.Done()
			ct := randomCiphertext()
			if err := storage.PutCt(tfhe.Hash{byte(i)}, ct); err != nil {
				t.Errorf("Failed to put ciphertext: %v", err)
			}
		}(i)
	}

	wg.Wait() // Wait for all goroutines to finish
}

func BenchmarkStoragePut(b *testing.B) {
	storage := storage2.InitStorage(storagePath)
	key := []byte("benchmarkKey")
	val := make([]byte, TestFileSize) // 64KB value to simulate large data
	rand.Read(val)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := storage.Put(0, key, val); err != nil {
			b.Fatalf("Failed to put data: %v", err)
		}
	}
}

// Add more tests and benchmarks as needed for Get, GetCt, Size, etc.
// Additional tests for the Storage interface

func TestStorageGet(t *testing.T) {
	storage := storage2.InitStorage(storagePath)
	key := []byte("testKey")
	expectedVal := []byte("testValue")

	// Pre-populate the storage with a known key-value pair
	if err := storage.Put(0, key, expectedVal); err != nil {
		t.Fatalf("Failed to put data: %v", err)
	}

	val, err := storage.Get(0, key)
	if err != nil {
		t.Fatalf("Failed to get data: %v", err)
	}

	if !bytes.Equal(val, expectedVal) {
		t.Errorf("Got unexpected value: got %v, want %v", val, expectedVal)
	}
}

func TestStorageGetCt(t *testing.T) {
	storage := storage2.InitStorage(storagePath)
	ct := randomCiphertext()
	hash := tfhe.Hash{0} // Simplified hash for testing

	if err := storage.PutCt(hash, ct); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := storage.GetCt(hash)
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	if !bytes.Equal(retrievedCt.Serialization, ct.Serialization) || !(retrievedCt.Hash() == ct.Hash()) {
		t.Errorf("Retrieved ciphertext does not match the original")
	}
}

func TestStorageVersioning(t *testing.T) {
	storage := storage2.InitStorage(storagePath)
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

//func TestStorageSize(t *testing.T) {
//	storage := InitStorage()
//	initialSize := storage.Size()
//
//	key := []byte("sizeTestKey")
//	val := []byte("sizeTestValue")
//
//	if err := storage.Put(0, key, val); err != nil {
//		t.Fatalf("Failed to put data: %v", err)
//	}
//
//	if newSize := storage.Size(); newSize <= initialSize {
//		t.Errorf("Storage size did not increase after put operation")
//	}
//}

// Benchmark for Get operation - focuses on reading performance
func BenchmarkStorageGet(b *testing.B) {
	storage := storage2.InitStorage(storagePath)
	key := []byte("benchmarkKey")
	val := make([]byte, TestFileSize) // 64KB value to simulate large data
	rand.Read(val)

	// Ensure the key exists for the benchmark
	if err := storage.Put(0, key, val); err != nil {
		b.Fatalf("Failed to put data for benchmark setup: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := storage.Get(0, key)
		if err != nil {
			b.Fatalf("Failed to get data: %v", err)
		}
	}
}

func BenchmarkConcurrentPut(b *testing.B) {
	storage := storage2.InitStorage(storagePath)

	key := []byte("benchmarkKey")
	val := make([]byte, TestFileSize) // 64KB value to simulate large data
	rand.Read(val)
	// Number of goroutines to use for concurrent writes
	concurrencyLevel := 10

	b.ResetTimer() // Reset the timer to exclude setup time

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(concurrencyLevel)

		for j := 0; j < concurrencyLevel; j++ {
			go func(j int) {
				defer wg.Done()
				// Generate unique key and value for each goroutine
				err := storage.Put(0, key, val)
				if err != nil {
					b.Errorf("Failed to put data: %v", err)
				}
			}(j)
		}

		wg.Wait() // Wait for all goroutines to finish
	}
}
