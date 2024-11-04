package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/fhenixprotocol/fheos/precompiles/types"
)

type EphemeralStorage interface {
	types.FheCipherTextStorage

	MarkForPersistence(common.Address, types.Hash) error
	MarkForDeletion(common.Address, types.Hash) error
	GetAllToPersist() []ContractCiphertext
	GetAllToDelete() []types.Hash
}

type ContractCiphertext struct {
	ContractAddress common.Address
	CipherTextHash  types.Hash
}

type EphemeralStorageImpl struct {
	//ltsCache []ContractCiphertext
	memstore *memorydb.Database
}

// Function to remove a ContractCiphertext from the slice by matching ContractAddress and CipherTextHash
func removeFromPersistence(slice []ContractCiphertext, address common.Address, hash types.Hash) ([]ContractCiphertext, bool) {
	for i, v := range slice {
		if v.ContractAddress == address && v.CipherTextHash == hash {
			// Replace the element with the last element
			slice[i] = slice[len(slice)-1]
			// Resize the slice to discard the last element
			return slice[:len(slice)-1], true
		}
	}
	// Return the original slice if the element was not found
	return slice, false
}

func removeFromDeletions(slice []types.Hash, hash types.Hash) ([]types.Hash, bool) {
	for i, v := range slice {
		if v == hash {
			// Replace the element with the last element
			slice[i] = slice[len(slice)-1]
			// Resize the slice to discard the last element
			return slice[:len(slice)-1], true
		}
	}
	// Return the original slice if the element was not found
	return slice, false
}

func (es *EphemeralStorageImpl) getPersistKey() []byte {
	return []byte("LTS")
}

func (es *EphemeralStorageImpl) getDeletionKey() []byte {
	return []byte("DEL")
}

func (es *EphemeralStorageImpl) encodePersistingCiphertexts(lts []ContractCiphertext) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(lts)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (es *EphemeralStorageImpl) encodeDeletionCiphertexts(deletion []types.Hash) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(deletion)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (es *EphemeralStorageImpl) DecodePersistingCiphertexts(raw []byte) ([]ContractCiphertext, error) {
	if len(raw) == 0 {
		return []ContractCiphertext{}, nil
	}

	var lts []ContractCiphertext
	err := gob.NewDecoder(bytes.NewBuffer(raw)).Decode(&lts)
	if err != nil {
		return nil, err
	}

	return lts, nil
}

func (es *EphemeralStorageImpl) DecodeDeletionCiphertexts(raw []byte) ([]types.Hash, error) {
	if len(raw) == 0 {
		return []types.Hash{}, nil
	}

	var deletion []types.Hash
	err := gob.NewDecoder(bytes.NewBuffer(raw)).Decode(&deletion)
	if err != nil {
		return nil, err
	}

	return deletion, nil
}

func (es *EphemeralStorageImpl) MarkForPersistence(contract common.Address, h types.Hash) error {
	// check if the contract is already in the deletion list
	deletionKey := es.getDeletionKey()
	rawDeletions, err := es.memstore.Get(deletionKey)
	if err != nil && err.Error() != "not found" {
		return err
	} else if rawDeletions != nil {
		currentDeletions, err := es.DecodeDeletionCiphertexts(rawDeletions)
		if err != nil {
			return err
		}

		newDeletions, isFoundDeletions := removeFromDeletions(currentDeletions, h)
		if isFoundDeletions {
			encodedDeletions, err := es.encodeDeletionCiphertexts(newDeletions)
			if err != nil {
				return err
			}

			// once we remove the deletion, it will stay in the persistent storage, no need to add it again
			return es.memstore.Put(deletionKey, encodedDeletions)
		}
	}

	persistKey := es.getPersistKey()
	rawLts, err := es.memstore.Get(persistKey)
	if err != nil {
		if err.Error() != "not found" {
			return err
		}
	}

	parsedLts, err := es.DecodePersistingCiphertexts(rawLts)
	if err != nil {
		return err
	}

	parsedLts = append(parsedLts, ContractCiphertext{
		ContractAddress: contract,
		CipherTextHash:  h,
	})

	encodedLts, err := es.encodePersistingCiphertexts(parsedLts)
	if err != nil {
		return err
	}

	return es.memstore.Put(persistKey, encodedLts)
}

func (es *EphemeralStorageImpl) DeleteCt(h types.Hash) error {
	return es.memstore.Delete(h[:])
}

func (es *EphemeralStorageImpl) MarkForDeletion(address common.Address, h types.Hash) error {
	err := es.DeleteCt(h)
	if err != nil {
		return err
	}

	// check if the contract is already in the persistence list
	persistKey := es.getPersistKey()
	rawPersists, err := es.memstore.Get(persistKey)
	if err != nil && err.Error() != "not found" {
		return err
	} else if rawPersists != nil {
		currentPersists, err := es.DecodePersistingCiphertexts(rawPersists)
		if err != nil {
			return err
		}

		newPersistence, isFoundPersists := removeFromPersistence(currentPersists, address, h)
		if isFoundPersists {
			encodedLts, err := es.encodePersistingCiphertexts(newPersistence)
			if err != nil {
				return err
			}
			// once we remove the persist task, the ct will never be added to long term storage (no need to delete it)
			return es.memstore.Put(persistKey, encodedLts)
		}
	}

	key := es.getDeletionKey()

	var parsedDeletion []types.Hash
	rawDeletion, err := es.memstore.Get(key)
	if err != nil {
		if err.Error() != "not found" {
			return err
		}
	} else {
		parsedDeletion, err = es.DecodeDeletionCiphertexts(rawDeletion)
		if err != nil {
			return err
		}
	}

	parsedDeletion = append(parsedDeletion, h)

	encodedDeletion, err := es.encodeDeletionCiphertexts(parsedDeletion)
	if err != nil {
		return err
	}

	return es.memstore.Put(key, encodedDeletion)
}

func (es *EphemeralStorageImpl) PutCt(h types.Hash, cipher *types.CipherTextRepresentation) error {
	if es.memstore == nil {
		return errors.New("memstore is nil")
	}

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(cipher)
	if err != nil {
		return err
	}

	// Use hash as key
	return es.memstore.Put(h[:], buf.Bytes())
}

func (es *EphemeralStorageImpl) GetCt(h types.Hash) (*types.CipherTextRepresentation, error) {

	if es.memstore == nil {
		return nil, errors.New("memstore is nil")
	}

	val, err := es.memstore.Get(h[:])
	if err != nil {
		return nil, err
	}

	var cipher types.CipherTextRepresentation
	err = gob.NewDecoder(bytes.NewBuffer(val)).Decode(&cipher)
	if err != nil {
		return nil, err
	}

	return &cipher, nil
}

func (es *EphemeralStorageImpl) HasCt(h types.Hash) bool {
	if es.memstore == nil {
		return false
	}

	ok, err := es.memstore.Has(h[:])
	if err != nil {
		return false
	}

	return ok
}

func (es *EphemeralStorageImpl) GetAllToPersist() []ContractCiphertext {
	key := es.getPersistKey()

	rawLts, err := es.memstore.Get(key)
	if err != nil {
		if err.Error() == "not found" {
			return []ContractCiphertext{}
		}

		return nil
	}

	parsedLts, err := es.DecodePersistingCiphertexts(rawLts)
	if err != nil {
		return nil
	}

	return parsedLts
}

func (es *EphemeralStorageImpl) GetAllToDelete() []types.Hash {
	key := es.getDeletionKey()

	rawDeletion, err := es.memstore.Get(key)
	if err != nil {
		if err.Error() == "not found" {
			return []types.Hash{}
		}

		return nil
	}

	parsedDeletions, err := es.DecodeDeletionCiphertexts(rawDeletion)
	if err != nil {
		return nil
	}

	return parsedDeletions
}

func NewEphemeralStorage(db *memorydb.Database) EphemeralStorage {

	if db == nil {
		return &EphemeralStorageImpl{
			memstore: memorydb.New(),
		}
	}

	return &EphemeralStorageImpl{
		memstore: db,
	}
}
