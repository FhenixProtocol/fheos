package storage

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"golang.org/x/exp/slices"
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
func (ms *MultiStore) PutCt(h types.Hash, cipher *types.FheEncrypted, owner common.Address) error {
	owners := []common.Address{owner}
	ciphertext := &types.CipherTextRepresentation{
		Data:   cipher,
		Owners: owners,
	}
	return ms.ephemeral.PutCt(h, ciphertext)
}

// AppendCt stores a ciphertext in the ephemeral storage - it does NOT mark it as LTS. The reason is that we want only SSTORE to mark it as LTS, which is
// called not only by the precompiles, only by the EVM hook from the evm interpreter.
// When a CT with the same hash exists in the disk storage, an owner is added to the list of owners.
func (ms *MultiStore) AppendCt(h types.Hash, cipher *types.FheEncrypted, owner common.Address) error {
	ct, _ := ms.getCtHelper(h)
	if ct != nil {
		ct.Owners = append(ct.Owners, owner)
		return ms.ephemeral.PutCt(h, ct)
	}

	return ms.PutCt(h, cipher, owner)
}

// )(ITZIK)(
func (ms *MultiStore) HackAppendCt(h types.Hash, cipher *types.FheEncrypted, owner common.Address) error {
	return ms.PutCt(h, cipher, owner)
}

func (ms *MultiStore) getCtHelper(h types.Hash) (*types.CipherTextRepresentation, error) {
	ct, err := ms.ephemeral.GetCt(h)

	if err != nil && err.Error() == "not found" {
		ct, err = ms.disk.GetCt(h)
	}

	if err != nil {
		return nil, err
	}

	return ct, err
}

func (ms *MultiStore) GetCtRepresentation(h types.Hash, caller common.Address) (*types.CipherTextRepresentation, error) {
	ct, err := ms.getCtHelper(h)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(ct.Owners, caller) {
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
func (ms *MultiStore) isOwner(ct *types.CipherTextRepresentation, owner common.Address) (bool, error) {
	if ct == nil {
		return false, fmt.Errorf("ciphertext not found")
	}

	return slices.Contains(ct.Owners, owner), nil
}

func (ms *MultiStore) AddOwner(h types.Hash, ct *types.CipherTextRepresentation, owner common.Address) error {
	if ct == nil {
		return fmt.Errorf("ciphertext not found")
	}

	isOwner, err := ms.isOwner(ct, owner)
	if err != nil {
		return err
	}

	if isOwner {
		return nil
	}

	ct.Owners = append(ct.Owners, owner)
	return ms.ephemeral.PutCt(h, ct)
}

func NewMultiStore(db *memorydb.Database, disk *FheosStorage) *MultiStore {
	return &MultiStore{
		ephemeral: NewEphemeralStorage(db),
		disk:      disk,
	}
}
