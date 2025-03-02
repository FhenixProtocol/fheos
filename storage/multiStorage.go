package storage

import (
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/fhenixprotocol/fheos/precompiles/types"
)

type IMultiStore interface {
	types.FheCipherTextStorage
}

type MultiStore struct {
	disk *FheosStorage
}

func (ms *MultiStore) PutCt(h types.Hash, cipher *types.FheEncrypted) error {
	err := ms.disk.PutCt(h, cipher)
	return err
}

func (ms *MultiStore) PutCtIfNotExist(h types.Hash, cipher *types.FheEncrypted) error {
	ct, _ := ms.GetCt(h)
	if ct != nil && ct.Placeholder == cipher.Placeholder {
		return nil
	}

	return ms.PutCt(h, cipher)
}

func (ms *MultiStore) GetCt(h types.Hash) (*types.FheEncrypted, error) {
	return ms.disk.GetCt(h)
}
func (ms *MultiStore) Has(h types.Hash) bool {
	return ms.disk.HasCt(h)
}

func (ms *MultiStore) DeleteCt(h types.Hash) error {
	return ms.disk.DeleteCt(h)
}

func NewMultiStore(db *memorydb.Database, disk *FheosStorage) *MultiStore {
	return &MultiStore{
		disk: disk,
	}
}
