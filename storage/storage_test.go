package storage_test

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
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

//func generateKeys() error {
//	if _, err := os.Stat("./keys/"); os.IsNotExist(err) {
//		err := os.Mkdir("./keys/", 0755)
//		if err != nil {
//			return err
//		}
//	}
//
//	if _, err := os.Stat("./keys/tfhe/"); os.IsNotExist(err) {
//		err := os.Mkdir("./keys/tfhe/", 0755)
//		if err != nil {
//			return err
//		}
//	}
//
//	err :=
//		tfhe.GenerateFheKeys("./keys/tfhe/", "./sks", "./cks", "./pks")
//	if err != nil {
//		return fmt.Errorf("error from tfhe GenerateFheKeys: %s", err)
//	}
//	return nil
//}

func init() {
	//err := generateKeys()
	//if err != nil {
	//	panic(err)
	//}
	//
	//if err != nil {
	//	panic(err)
	//}

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

	ct.Hash()

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

			ctRep := &types.CipherTextRepresentation{
				Data:     (*types.FheEncrypted)(ct),
				Owners:   nil,
				RefCount: 0,
			}

			if err := storage.PutCt(types.Hash(fhe.Hash{byte(i)}), ctRep); err != nil {
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

			ctRep := &types.CipherTextRepresentation{
				Data:     (*types.FheEncrypted)(ct),
				Owners:   nil,
				RefCount: 0,
			}

			if err := storage.PutCt(types.Hash(fhe.Hash{byte(i)}), ctRep); err != nil {
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

	ctRep := &types.CipherTextRepresentation{
		Data:     (*types.FheEncrypted)(ct),
		Owners:   nil,
		RefCount: 0,
	}

	if err := storage.PutCt(types.Hash(hash), ctRep); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := storage.GetCt(types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	if !bytes.Equal(retrievedCt.Data.Data, ct.Data) || !((*fhe.FheEncrypted)(retrievedCt.Data).Hash() == ct.Hash()) {
		t.Errorf("Retrieved ciphertext does not match the original")
	}
	assert.Equal(t, ct.Compact, retrievedCt.Data.Compact)
	assert.Equal(t, ct.Compressed, retrievedCt.Data.Compressed)
}

func TestStorageEphemeralCt(t *testing.T) {
	storage, err := storage2.InitStorage(storagePath)

	if err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}

	ct := randomCiphertext()
	hash := fhe.Hash{0} // Simplified hash for testing

	//storage.SetEphemeral()

	ctRep := &types.CipherTextRepresentation{
		Data:     (*types.FheEncrypted)(ct),
		Owners:   nil,
		RefCount: 0,
	}

	if err := storage.PutCt(types.Hash(hash), ctRep); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := storage.GetCt(types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	if !bytes.Equal(retrievedCt.Data.Data, ct.Data) || !((*fhe.FheEncrypted)(retrievedCt.Data).Hash() == ct.Hash()) {
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

	ctRep := &types.CipherTextRepresentation{
		Data:     (*types.FheEncrypted)(ct),
		Owners:   nil,
		RefCount: 0,
	}

	if err := storage.PutCt(types.Hash(hash), ctRep); err != nil {
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

	if !bytes.Equal(retrievedCt.Data.Data, ct.Data) || !((*fhe.FheEncrypted)(retrievedCt.Data).Hash() == ct.Hash()) {
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

				ctRep := &types.CipherTextRepresentation{
					Data:     (*types.FheEncrypted)(ct),
					Owners:   nil,
					RefCount: 0,
				}

				if err := storage.PutCt(types.Hash(fhe.Hash{byte(i)}), ctRep); err != nil {
					b.Errorf("Failed to put ciphertext: %v", err)
				}
			}(j)
		}

		wg.Wait() // Wait for all goroutines to finish
	}
}

func TestEphemeralStorageImpl_Basic(t *testing.T) {
	ephemeralStorage := storage2.NewEphemeralStorage(nil)

	ct := randomCiphertext()
	hash := fhe.Hash{0} // Simplified hash for testing

	ctRep := &types.CipherTextRepresentation{
		Data:     (*types.FheEncrypted)(ct),
		Owners:   nil,
		RefCount: 0,
	}

	if err := ephemeralStorage.PutCt(types.Hash(hash), ctRep); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := ephemeralStorage.GetCt(types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	if !bytes.Equal(retrievedCt.Data.Data, ct.Data) || !((*fhe.FheEncrypted)(retrievedCt.Data).Hash() == ct.Hash()) {
		t.Errorf("Retrieved ciphertext does not match the original")
	}
}

func TestEphemeralStorageImpl_Lts(t *testing.T) {
	ephemeralStorage := storage2.NewEphemeralStorage(nil)

	ownerAddress := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	ct := randomCiphertext()
	hash := fhe.Hash{0} // Simplified hash for testing

	ctRep := &types.CipherTextRepresentation{
		Data:     (*types.FheEncrypted)(ct),
		Owners:   nil,
		RefCount: 0,
	}

	if err := ephemeralStorage.PutCt(types.Hash(hash), ctRep); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	err := ephemeralStorage.MarkForPersistence(ownerAddress, types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to mark ciphertext as LTS: %v", err)
	}

	// Test getting all LTS ciphertexts
	ltsCts := ephemeralStorage.GetAllToPersist()
	if len(ltsCts) != 1 {
		t.Fatalf("Expected 1 LTS ciphertext, got %d", len(ltsCts))
	}

	if !bytes.Equal(ltsCts[0].CipherTextHash[:], hash[:]) {
		t.Errorf("LTS ciphertext hash mismatch")
	}

	retrievedCt, err := ephemeralStorage.GetCt(ltsCts[0].CipherTextHash)
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	if !bytes.Equal(retrievedCt.Data.Data, ct.Data) || !((*fhe.FheEncrypted)(retrievedCt.Data).Hash() == ct.Hash()) {
		t.Errorf("Retrieved ciphertext does not match the original")
	}
}

func TestEphemeralStorageImpl_HasCt(t *testing.T) {
	ephemeralStorage := storage2.NewEphemeralStorage(nil)

	ct := randomCiphertext()
	hash := fhe.Hash{0} // Simplified hash for testing

	ctRep := &types.CipherTextRepresentation{
		Data:     (*types.FheEncrypted)(ct),
		Owners:   nil,
		RefCount: 0,
	}

	if err := ephemeralStorage.PutCt(types.Hash(hash), ctRep); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := ephemeralStorage.GetCt(types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	if !bytes.Equal(retrievedCt.Data.Data, ct.Data) || !((*fhe.FheEncrypted)(retrievedCt.Data).Hash() == ct.Hash()) {
		t.Errorf("Retrieved ciphertext does not match the original")
	}

	if !ephemeralStorage.HasCt(types.Hash(hash)) {
		t.Errorf("Expected to have ciphertext")
	}

	hash2 := fhe.Hash{1}
	if ephemeralStorage.HasCt(types.Hash(hash2)) {
		t.Errorf("Expected to not have ciphertext")
	}
}

func TestEphemeralStorageImpl_OverrideValue(t *testing.T) {
	ephemeralStorage := storage2.NewEphemeralStorage(nil)

	ct := randomCiphertext()
	hash := fhe.Hash{0} // Simplified hash for testing

	ct.Placeholder = true

	ctRep := &types.CipherTextRepresentation{
		Data:     (*types.FheEncrypted)(ct),
		Owners:   nil,
		RefCount: 0,
	}

	if err := ephemeralStorage.PutCt(types.Hash(hash), ctRep); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := ephemeralStorage.GetCt(types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	if retrievedCt.Data.Placeholder != true {
		t.Errorf("Retrieved ciphertext does not match the original")
	}

	ctRep.Data.Placeholder = false

	if err := ephemeralStorage.PutCt(types.Hash(hash), ctRep); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err = ephemeralStorage.GetCt(types.Hash(hash))
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	if retrievedCt.Data.Placeholder != false {
		t.Errorf("Retrieved ciphertext does not match the original")
	}
}

func TestMultiStore_AppendCt(t *testing.T) {
	diskStorage, err := storage2.InitStorage(storagePath)
	multiStore := storage2.NewMultiStore(nil, diskStorage)

	ct := randomCiphertext()
	hash := fhe.Hash{100} // this key needs to be unique for the test

	owner := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

	if err := multiStore.AppendCt(types.Hash(hash), (*types.FheEncrypted)(ct), owner); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := multiStore.GetCt(types.Hash(hash), owner)
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	assert.Equal(t, (*fhe.FheEncrypted)(retrievedCt).Hash(), ct.Hash())

	if !bytes.Equal(retrievedCt.Data, ct.Data) {
		t.Errorf("Retrieved ciphertext does not match the original")
	}
}

func TestMultiStore_AppendCtMultipleOwner(t *testing.T) {
	diskStorage, err := storage2.InitStorage(storagePath)
	multiStore := storage2.NewMultiStore(nil, diskStorage)

	ct := randomCiphertext()
	hash := fhe.Hash{101} // this key needs to be unique for the test

	owner := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

	if err := multiStore.AppendCt(types.Hash(hash), (*types.FheEncrypted)(ct), owner); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	owner2 := common.HexToAddress("0x1234567890abcdef1234567890abcdef123456ff")

	if err := multiStore.AppendCt(types.Hash(hash), (*types.FheEncrypted)(ct), owner2); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := multiStore.GetCt(types.Hash(hash), owner)
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	assert.Equal(t, (*fhe.FheEncrypted)(retrievedCt).Hash(), ct.Hash())

	retrievedCt, err = multiStore.GetCt(types.Hash(hash), owner2)
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	assert.Equal(t, (*fhe.FheEncrypted)(retrievedCt).Hash(), ct.Hash())
}

func TestMultiStore_AppendCt_Fail_NotOwner(t *testing.T) {
	diskStorage, err := storage2.InitStorage(storagePath)
	multiStore := storage2.NewMultiStore(nil, diskStorage)

	ct := randomCiphertext()
	hash := fhe.Hash{102} // this key needs to be unique for the test

	owner := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

	if err := multiStore.AppendCt(types.Hash(hash), (*types.FheEncrypted)(ct), owner); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	loser := common.HexToAddress("0x1234567890abcdef1234567890abcdefffffffff")

	_, err = multiStore.GetCt(types.Hash(hash), loser)
	// expect an error, fail the test if not
	if err == nil {
		t.Fatalf("Expected to fail to get ciphertext")
	}
}

func TestMultiStore_AppendCtPlaceholderReplace(t *testing.T) {
	diskStorage, err := storage2.InitStorage(storagePath)
	multiStore := storage2.NewMultiStore(nil, diskStorage)

	ct := randomCiphertext()
	ct.Placeholder = true
	hash := fhe.Hash{104} // this key needs to be unique for the test

	owner := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

	if err := multiStore.AppendCt(types.Hash(hash), (*types.FheEncrypted)(ct), owner); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	owner2 := common.HexToAddress("0x1234567890abcdef1234567890abcdef123456ff")

	ct.Placeholder = false

	if err := multiStore.AppendCt(types.Hash(hash), (*types.FheEncrypted)(ct), owner2); err != nil {
		t.Fatalf("Failed to put ciphertext: %v", err)
	}

	retrievedCt, err := multiStore.GetCt(types.Hash(hash), owner)
	if err != nil {
		t.Fatalf("Failed to get ciphertext: %v", err)
	}

	assert.Equal(t, (*fhe.FheEncrypted)(retrievedCt).Hash(), ct.Hash())
	assert.Equal(t, retrievedCt.Placeholder, false)
}
