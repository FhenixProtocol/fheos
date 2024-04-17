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
	GetAllToPersist() []ContractCiphertext
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

func (es *EphemeralStorageImpl) encodePersistingCiphertexts(lts []ContractCiphertext) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(lts)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (es *EphemeralStorageImpl) DecodePersistingCiphertexts(raw []byte) ([]ContractCiphertext, error) {
	var lts []ContractCiphertext
	err := gob.NewDecoder(bytes.NewBuffer(raw)).Decode(&lts)
	if err != nil {
		return nil, err
	}

	return lts, nil
}

func (es *EphemeralStorageImpl) MarkForPersistence(contract common.Address, h types.Hash) error {
	key := es.getPersistKey()

	var parsedLts []ContractCiphertext
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

	parsedLts = append(parsedLts, ContractCiphertext{
		ContractAddress: contract,
		CipherTextHash:  h,
	})

	encodedLts, err := es.encodePersistingCiphertexts(parsedLts)
	if err != nil {
		return err
	}

	return es.memstore.Put(key, encodedLts)
}

func (es *EphemeralStorageImpl) PutCt(h types.Hash, cipher *types.FheEncrypted) error {

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

func (es *EphemeralStorageImpl) GetCt(h types.Hash) (*types.FheEncrypted, error) {

	if es.memstore == nil {
		return nil, errors.New("memstore is nil")
	}

	val, err := es.memstore.Get(h[:])
	if err != nil {
		return nil, err
	}

	var cipher types.FheEncrypted
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
