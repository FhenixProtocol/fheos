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
	MarkForDeletion(types.Hash) error
	GetAllToPersist() map[types.Hash]ContractCiphertext
	GetAllToDelete() map[types.Hash]uint8 // uint8 here is just to have minimum size map
}

type ContractCiphertext struct {
	ContractAddress common.Address
	CipherTextHash  types.Hash
}

type EphemeralStorageImpl struct {
	//ltsCache []ContractCiphertext
	memstore *memorydb.Database
}

func (es *EphemeralStorageImpl) getPersistKey() []byte {
	return []byte("LTS")
}

func (es *EphemeralStorageImpl) getDeletionKey() []byte {
	return []byte("DEL")
}

func (es *EphemeralStorageImpl) encodePersistingCiphertexts(lts map[types.Hash]ContractCiphertext) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(lts)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (es *EphemeralStorageImpl) encodeDeletionCiphertexts(deletion map[types.Hash]uint8) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(deletion)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (es *EphemeralStorageImpl) DecodePersistingCiphertexts(raw []byte) (map[types.Hash]ContractCiphertext, error) {
	var lts map[types.Hash]ContractCiphertext
	err := gob.NewDecoder(bytes.NewBuffer(raw)).Decode(&lts)
	if err != nil {
		return nil, err
	}

	return lts, nil
}

func (es *EphemeralStorageImpl) DecodeDeletionCiphertexts(raw []byte) (map[types.Hash]uint8, error) {
	var deletion map[types.Hash]uint8
	err := gob.NewDecoder(bytes.NewBuffer(raw)).Decode(&deletion)
	if err != nil {
		return nil, err
	}

	return deletion, nil
}

func (es *EphemeralStorageImpl) MarkForPersistence(contract common.Address, h types.Hash) error {
	// We might have marked this ciphertext for deletion before
	es.unmarkForDeletion(h)

	key := es.getPersistKey()

	var parsedLts = map[types.Hash]ContractCiphertext{}
	rawLts, err := es.memstore.Get(key)
	if err != nil {
		if err.Error() != "not found" {
			return err
		}
	} else {
		parsedLts, err = es.DecodePersistingCiphertexts(rawLts)
		if err != nil {
			return err
		}
	}

	if _, ok := parsedLts[h]; !ok {
		parsedLts[h] = ContractCiphertext{
			ContractAddress: contract,
			CipherTextHash:  h,
		}

		encodedLts, err := es.encodePersistingCiphertexts(parsedLts)
		if err != nil {
			return err
		}

		return es.memstore.Put(key, encodedLts)
	}

	return nil
}

func (es *EphemeralStorageImpl) DeleteCt(h types.Hash) error {
	return es.memstore.Delete(h[:])
}

func (es *EphemeralStorageImpl) MarkForDeletion(h types.Hash) error {
	// We might have marked this ciphertext for persistence before
	es.unmarkForPersistence(h)

	key := es.getDeletionKey()
	var parsedDeletion = map[types.Hash]uint8{}
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

	if _, ok := parsedDeletion[h]; !ok {
		parsedDeletion[h] = 1 // Arbitrary value

		encodedDeletion, err := es.encodeDeletionCiphertexts(parsedDeletion)
		if err != nil {
			return err
		}

		return es.memstore.Put(key, encodedDeletion)
	}

	return nil
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

func (es *EphemeralStorageImpl) GetAllToPersist() map[types.Hash]ContractCiphertext {
	key := es.getPersistKey()

	rawLts, err := es.memstore.Get(key)
	if err != nil {
		if err.Error() == "not found" {
			return map[types.Hash]ContractCiphertext{}
		}

		return nil
	}

	parsedLts, err := es.DecodePersistingCiphertexts(rawLts)
	if err != nil {
		return nil
	}

	return parsedLts
}

func (es *EphemeralStorageImpl) GetAllToDelete() map[types.Hash]uint8 {
	key := es.getDeletionKey()

	rawDeletion, err := es.memstore.Get(key)
	if err != nil {
		if err.Error() == "not found" {
			return map[types.Hash]uint8{}
		}

		return nil
	}

	parsedDeletions, err := es.DecodeDeletionCiphertexts(rawDeletion)
	if err != nil {
		return nil
	}

	return parsedDeletions
}

func (es *EphemeralStorageImpl) unmarkForDeletion(h types.Hash) error {
	key := es.getDeletionKey()

	rawDeletion, err := es.memstore.Get(key)
	if err != nil {
		if err.Error() == "not found" {
			return nil
		}
		return err
	}

	deletions, err := es.DecodeDeletionCiphertexts(rawDeletion)
	if err != nil {
		return err
	}

	delete(deletions, h)

	encoded, err := es.encodeDeletionCiphertexts(deletions)
	if err != nil {
		return err
	}

	return es.memstore.Put(key, encoded)
}

func (es *EphemeralStorageImpl) unmarkForPersistence(h types.Hash) error {
	key := es.getPersistKey()

	rawLts, err := es.memstore.Get(key)
	if err != nil {
		if err.Error() == "not found" {
			return nil
		}
		return err
	}

	lts, err := es.DecodePersistingCiphertexts(rawLts)
	if err != nil {
		return err
	}

	delete(lts, h)

	encoded, err := es.encodePersistingCiphertexts(lts)
	if err != nil {
		return err
	}

	return es.memstore.Put(key, encoded)
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
