package storage

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	"golang.org/x/exp/slices"
)

type IMultiStore interface {
	types.FheCipherTextStorage
	GetEphemeral() *EphemeralStorage

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

func (ms *MultiStore) GetEphemeral() EphemeralStorage {
	return ms.ephemeral
}

// PutCt stores a ciphertext in the ephemeral storage - it does NOT mark it as LTS. The reason is that we want only SSTORE to mark it as LTS, which is
// called not only by the precompiles, but also by the EVM hook from the evm interpreter.
func (ms *MultiStore) PutCt(h types.Hash, cipher *types.CipherTextRepresentation) error {
	err := ms.ephemeral.PutCt(h, cipher)
	return err
}

func (ms *MultiStore) SetAsyncCtStart(h types.Hash) error {
	return ms.ephemeral.SetAsyncCtStart(h)
}

func (ms *MultiStore) SetAsyncCtDone(h types.Hash) error {
	return ms.ephemeral.SetAsyncCtDone(h)
}

func (ms *MultiStore) IsAsyncCtDone() (bool, error) {
	return ms.ephemeral.IsAsyncCtDone()
}

// AppendPhV stores a PlaceholderValue in ephemeral storage - it does NOT mark it as LTS. The reason is that we want only SSTORE to mark it as LTS, which is
// TODO: Verify that we don't need to worry about RefCount
func (ms *MultiStore) AppendPhV(h types.Hash, cipher *types.FheEncrypted, owner common.Address) error {
	return ms.PutCt(h, &types.CipherTextRepresentation{Data: cipher, Owners: []common.Address{owner}, RefCount: 1})
}

// AppendCt stores a ciphertext in the ephemeral storage - it does NOT mark it as LTS. The reason is that we want only SSTORE to mark it as LTS, which is
// called not only by the precompiles, only by the EVM hook from the evm interpreter.
// When a CT with the same hash exists in the disk storage, an owner is added to the list of owners.
func (ms *MultiStore) AppendCt(h types.Hash, cipher *types.FheEncrypted, owner common.Address) error {
	ct, _ := ms.getCtHelper(h)
	if cipher.IsTriviallyEncrypted {
		// Already exists
		if ct != nil {
			return nil
		}

		return ms.PutCt(h, &types.CipherTextRepresentation{Data: cipher, Owners: []common.Address{}, RefCount: 0})
	}

	// If we're replacing a placeholder value just hard replace the value
	if ct != nil && ct.Data.Placeholder && cipher.Placeholder == false {
		return ms.PutCt(h, &types.CipherTextRepresentation{Data: cipher, Owners: ct.Owners, RefCount: ct.RefCount})
	}

	// Exists but not trivially encrypted
	if ct != nil {
		ct.Owners = append(ct.Owners, owner)
		ct.RefCount++
		return ms.ephemeral.PutCt(h, ct)
	}

	// RefCount is starting at 0, if the value will be stored in the persistent storage it will set to 1
	return ms.PutCt(h, &types.CipherTextRepresentation{Data: cipher, Owners: []common.Address{owner}, RefCount: 0})
}

func (ms *MultiStore) getCtHelper(h types.Hash) (*types.CipherTextRepresentation, error) {
	ct, err := ms.ephemeral.GetCt(h)

	if err != nil && err.Error() == "not found" {
		return ms.disk.GetCt(h)
	}

	return ct, err
}

func (ms *MultiStore) GetCtRepresentation(h types.Hash, caller common.Address) (*types.CipherTextRepresentation, error) {
	ct, err := ms.getCtHelper(h)
	if err != nil {
		return nil, err
	}

	//This is new
	owner, err := ms.isOwner(h, ct, caller)
	if err != nil {
		return nil, err
	}

	if !owner {
		return nil, fmt.Errorf("contract is not allowed to access the ciphertext (ct: %s, contract: %s)", hex.EncodeToString(h[:]), caller.String())
	}

	return ct, nil
}
func (ms *MultiStore) GetCt(h types.Hash, caller common.Address) (*types.FheEncrypted, error) {
	ct, err := ms.GetCtRepresentation(h, caller)
	if (err != nil) || (ct == nil) {
		return nil, err
	}

	return ct.Data, nil
}
func (ms *MultiStore) isOwner(h types.Hash, ct *types.CipherTextRepresentation, owner common.Address) (bool, error) {
	if ct == nil {
		return false, fmt.Errorf("ciphertext not found")
	}
	// No ownership for trivially encrypted values
	if fhe.IsTriviallyEncryptedCtHash(h) {
		return true, nil
	}

	return slices.Contains(ct.Owners, owner), nil
}
func (ms *MultiStore) DereferenceCiphertext(hash types.Hash) (bool, error) {
	ct, _ := ms.getCtHelper(hash)
	if ct == nil {
		return false, nil
	}

	if ct.RefCount <= 1 {
		return true, ms.ephemeral.MarkForDeletion(hash)
	}

	ct.RefCount--

	return false, ms.PutCt(hash, ct)
}
func (ms *MultiStore) ReferenceCiphertext(hash types.Hash) error {
	ct, _ := ms.getCtHelper(hash)
	if ct == nil {
		return fmt.Errorf("can not reference a ciphertext with hash: %s", hex.EncodeToString(hash[:]))
	}

	ct.RefCount++
	return ms.PutCt(hash, ct)
}

func (ms *MultiStore) AddOwner(h types.Hash, ct *types.CipherTextRepresentation, owner common.Address) error {
	if ct == nil {
		return fmt.Errorf("ciphertext not found")
	}

	// isOwner will return true for trivially encrypted values for every contract
	// The meaning is that we won't add any owner for trivially encrypted value
	isOwner, err := ms.isOwner(h, ct, owner)
	if err != nil {
		return err
	}

	if isOwner {
		return nil
	}

	ct.Owners = append(ct.Owners, owner)
	return ms.ephemeral.PutCt(h, ct)
}

func (ms *MultiStore) Has(h types.Hash) bool {
	return ms.ephemeral.HasCt(h)
}

func NewMultiStore(db *memorydb.Database, disk *FheosStorage) *MultiStore {
	return &MultiStore{
		ephemeral: NewEphemeralStorage(db),
		disk:      disk,
	}
}
