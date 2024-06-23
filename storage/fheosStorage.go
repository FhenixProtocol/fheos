//go:build amd64 || arm64

package storage

import (
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/fheos/storage/pebble"
)

// FheosStorage is a wrapper around the diskStore - it is the main storage interface for the Fheos DB, which stores
// all the ciphertexts that are not ephemeral - i.e. that are stored in the chain state.
type FheosStorage struct {
	diskStore types.Storage
}

func (fs *FheosStorage) DeleteCt(h types.Hash) error {
	return fs.diskStore.DeleteCt(h)
}

func (fs *FheosStorage) PutCt(h types.Hash, cipher *types.CipherTextRepresentation) error {
	return fs.diskStore.PutCt(h, cipher)
}

func (fs *FheosStorage) GetCt(h types.Hash) (*types.CipherTextRepresentation, error) {
	a, e := fs.diskStore.GetCt(h)
	return a, e
}

func (fs *FheosStorage) PutVersion(v uint64) error {
	return fs.diskStore.PutVersion(v)
}

func (fs *FheosStorage) GetVersion() (uint64, error) {
	return fs.diskStore.GetVersion()
}

func newFheosStorage(diskStore types.Storage) *FheosStorage {

	if diskStore == nil {
		panic("failed to initialize FheosStorage: diskStore is nil")
	}

	return &FheosStorage{
		diskStore: diskStore,
		//memStore:    memStore,
		//isEphemeral: false,
	}
}

func InitStorage(path string) (*FheosStorage, error) {
	storage, err := pebble.NewStorage(path)
	if err != nil {
		return nil, err
	}

	fheosStore := newFheosStorage(storage)

	return fheosStore, nil
}
