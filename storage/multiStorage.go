package storage

import (
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/fhenixprotocol/fheos/precompiles/types"
)

type IMultiStore interface {
	types.FheCipherTextStorage

	//IsEphemeral() bool
	//SetEphemeral()
	//Reset()
}

// MultiStore implements the MultiStore interface and is used by the precompiles to store ciphertexts.
// it is a wrapper around the ephemeral and disk storage and manages the balance between long term storage and ephemeral storage.
type MultiStore struct {
	ephemeral EphemeralStorage
	disk      *FheosStorage
}

// PutCt stores a ciphertext in the ephemeral storage - it does NOT mark it as LTS. The reason is that we want only SSTORE to mark it as LTS, which is
// called not only by the precompiles, only by the EVM hook from the evm interpreter.
func (ms *MultiStore) PutCt(h types.Hash, cipher *types.FheEncrypted) error {
	return ms.ephemeral.PutCt(h, cipher)
}

func (ms *MultiStore) GetCt(h types.Hash) (*types.FheEncrypted, error) {
	ct, err := ms.ephemeral.GetCt(h)
	if err != nil {
		// if we didn't find the ciphertext in the ephemeral storage, we try to get it from the disk storage
		if err.Error() == "not found" {
			return ms.disk.GetCt(h)
		}
	}

	return ct, nil
}

func NewMultiStore(db *memorydb.Database, disk *FheosStorage) *MultiStore {
	return &MultiStore{
		ephemeral: NewEphemeralStorage(db),
		disk:      disk,
	}
}
