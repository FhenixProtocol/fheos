package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

type EphemeralStorage interface {
	types.FheCipherTextStorage

	MarkForPlaceholding(h types.Hash) error
	MarkForPersistence(common.Address, types.Hash) error
	MarkForDeletion(types.Hash) error

	GetAllToPlacehold() []types.Hash
	GetAllToPersist() []ContractCiphertext
	GetAllToDelete() []types.Hash

	BlockUntilPlaceholderValuesRemoved(errorChannel chan error) error
	BlockUntilPlaceHolderValueReplaced(val [32]byte, errorChannel chan error) ([32]byte, error)
}

type ContractCiphertext struct {
	ContractAddress common.Address
	CipherTextHash  types.Hash
}

type AsyncCtList = []types.Hash

const AsyncCtKey = "ASYNCCT"

type EphemeralStorageImpl struct {
	//ltsCache []ContractCiphertext
	memstore *memorydb.Database
	// PlaceholderValues are "leetboobs" i.e. locations that may or may not hold the actual daedbeef addresses
	// They are essentially used as futures for async eval operations
	knownPlaceholderValues [][]byte
}

// ----------------------------------------------------------
// Helper functions for working with the [][32]byte container
// ----------------------------------------------------------
func (es *EphemeralStorageImpl) append(item []byte) {
	es.knownPlaceholderValues = append(es.knownPlaceholderValues, item)
}

func (es *EphemeralStorageImpl) remove(item []byte) {
	for i, v := range es.knownPlaceholderValues {
		if bytes.Equal(v[:], item[:]) {
			es.knownPlaceholderValues = append(es.knownPlaceholderValues[:i], es.knownPlaceholderValues[i+1:]...)
			break
		}
	}
}
func removeFromSliceAt(kpv []types.Hash, index int) []types.Hash {
	return append(kpv[:index], kpv[index+1:]...)
}

func (es *EphemeralStorageImpl) contains(item [32]byte) bool {
	for _, v := range es.knownPlaceholderValues {
		if bytes.Equal(v[:], item[:]) {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------
// ----------------------------------------------------------

//TODO: For now lets stick with it being in FHEencrypted and just getting the fhe module here

//func (es *EphemeralStorageImpl) StorePlaceholderValue(val [32]byte) {
//	log.Info("{{StorePlaceholderValue}}", "VAL:", fhe.Hash(val).Hex())
//	if fhe.IsPlaceholderValue(val[:]) {
//		es.append(val[:])
//	}
//}

// BlockUntilPlaceholderValuesRemoved is a special "blocking" function, and it is not intended to be used
// anywhere where you do not expect this to finish. It requires async operations on the background to finish working.
func (es *EphemeralStorageImpl) BlockUntilPlaceholderValuesRemoved(errorChannel chan error) error {
	knownPlaceholderValues := es.GetAllToPlacehold()

	for len(knownPlaceholderValues) > 0 {
		select {
		case err := <-errorChannel:
			return err
		default:
		}
		for i := 0; i < len(knownPlaceholderValues); {
			hash := knownPlaceholderValues[i]
			addr, _ := es.GetCt(hash)
			var addrToCheck [32]byte
			copy(addrToCheck[:], addr.Data.Data)
			if fhe.IsCtHash(addrToCheck) {
				knownPlaceholderValues = removeFromSliceAt(knownPlaceholderValues, i)
			} else {
				i++
			}
		}
	}
	return nil
}

// BlockUntilPlaceHolderValueReplaced is a special "blocking" function, and it is not intended to be used
// anywhere where you do not expect this to finish. It requires async operations on the background to finish working.
func (es *EphemeralStorageImpl) BlockUntilPlaceHolderValueReplaced(val [32]byte, errorChannel chan error) ([32]byte, error) {
	hash := val

	if fhe.IsPlaceholderValue(hash[:]) {
		addr, _ := es.GetCt(hash)
		var addrToCheck [32]byte
		copy(addrToCheck[:], addr.Data.Data)
		for fhe.IsZero(addrToCheck[:]) {
			select {
			case err := <-errorChannel:
				var nothing [32]byte
				return nothing, err
			default:
			}
			addr, _ = es.GetCt(hash)
			copy(addrToCheck[:], addr.Data.Data)
		}
		copy(hash[:], addrToCheck[:])
	}

	return hash, nil
}

func (es *EphemeralStorageImpl) getPersistKey() []byte {
	return []byte("LTS")
}

func (es *EphemeralStorageImpl) getDeletionKey() []byte {
	return []byte("DEL")
}

func (es *EphemeralStorageImpl) getPlaceholderKey() []byte { return []byte("PHV") }

// I shouldn't repeat myself, but for now lets leave it like this
func encodeKnownPlaceholderValues(phv []types.Hash) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(phv)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (es *EphemeralStorageImpl) encodeKnownPlaceholderValues(phv []types.Hash) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(phv)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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

func (es *EphemeralStorageImpl) DecodeKnownPlaceholderValues(phv []byte) ([]types.Hash, error) {
	var phvHash []types.Hash
	err := gob.NewDecoder(bytes.NewBuffer(phv)).Decode(&phvHash)
	if err != nil {
		return nil, err
	}

	return phvHash, nil
}

func (es *EphemeralStorageImpl) DecodePersistingCiphertexts(raw []byte) ([]ContractCiphertext, error) {
	var lts []ContractCiphertext
	err := gob.NewDecoder(bytes.NewBuffer(raw)).Decode(&lts)
	if err != nil {
		return nil, err
	}

	return lts, nil
}

func (es *EphemeralStorageImpl) DecodeDeletionCiphertexts(raw []byte) ([]types.Hash, error) {
	var deletion []types.Hash
	err := gob.NewDecoder(bytes.NewBuffer(raw)).Decode(&deletion)
	if err != nil {
		return nil, err
	}

	return deletion, nil
}

func (es *EphemeralStorageImpl) GetAsyncList() ([]types.Hash, error) {
	raw, err := es.memstore.Get([]byte(AsyncCtKey))
	if err != nil {
		return nil, err
	}

	var asyncList []types.Hash
	err = gob.NewDecoder(bytes.NewBuffer(raw)).Decode(&asyncList)

	return asyncList, err
}

func (es *EphemeralStorageImpl) SetAsyncList(asyncList []types.Hash) error {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(asyncList)
	if err != nil {
		return err
	}

	return es.memstore.Put([]byte(AsyncCtKey), buf.Bytes())
}

func (es *EphemeralStorageImpl) MarkForPlaceholding(h types.Hash) error {
	key := es.getPlaceholderKey()

	var parsedPlaceholders []types.Hash
	rawPhv, err := es.memstore.Get(key)
	if err != nil {
		if err.Error() != "not found" {
			return err
		}
	} else {
		parsedPlaceholders, err = es.DecodeKnownPlaceholderValues(rawPhv)
		if err != nil {
			return err
		}
	}

	parsedPlaceholders = append(parsedPlaceholders, h)

	encodedLts, err := es.encodeKnownPlaceholderValues(parsedPlaceholders)
	if err != nil {
		return err
	}

	return es.memstore.Put(key, encodedLts)
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

func (es *EphemeralStorageImpl) DeleteCt(h types.Hash) error {
	return es.memstore.Delete(h[:])
}

func (es *EphemeralStorageImpl) MarkForDeletion(h types.Hash) error {
	err := es.DeleteCt(h)
	if err != nil {
		return err
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
func (es *EphemeralStorageImpl) GetAllToPlacehold() []types.Hash {
	key := es.getPlaceholderKey()

	rawPhv, err := es.memstore.Get(key)
	if err != nil {
		if err.Error() == "not found" {
			return []types.Hash{}
		}

		return nil
	}

	parsedPlaceholders, err := es.DecodeDeletionCiphertexts(rawPhv)
	if err != nil {
		return nil
	}

	return parsedPlaceholders
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

func (es *EphemeralStorageImpl) SetAsyncCtStart(h types.Hash) error {
	asyncList, err := es.GetAsyncList()
	if err != nil {
		return err
	}

	asyncList = append(asyncList, h)
	return es.SetAsyncList(asyncList)
}

func (es *EphemeralStorageImpl) SetAsyncCtDone(h types.Hash) error {
	asyncList, err := es.GetAsyncList()
	if err != nil {
		return err
	}

	for i, v := range asyncList {
		if v == h {
			asyncList = append(asyncList[:i], asyncList[i+1:]...)
			break
		}
	}

	return es.SetAsyncList(asyncList)
}

func (es *EphemeralStorageImpl) IsAsyncCtDone() (bool, error) {
	asyncList, err := es.GetAsyncList()
	if err != nil {
		return false, err
	}

	return len(asyncList) == 0, nil
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
