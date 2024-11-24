package hooks

import (
	"encoding/hex"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	fheos "github.com/fhenixprotocol/fheos/precompiles"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

var lock sync.RWMutex

type FheOSHooks interface {
	StoreCiphertextHook(contract common.Address, loc [32]byte, original common.Hash, val [32]byte) error
	StoreGasHook(contract common.Address, loc [32]byte, val [32]byte) (uint64, uint64)
	LoadCiphertextHook() [32]byte
	EvmCallStart()
	EvmCallEnd(evmSuccess bool)
	ContractCall(isSimulation bool, callType int, caller common.Address, addr common.Address, input []byte)
	ContractCallReturn(isSimulation bool, callType int, caller common.Address, addr common.Address, output []byte)
}

type FheOSHooksImpl struct {
	evm *vm.EVM
}

func (h *FheOSHooksImpl) updateCiphertextReferences(original common.Hash, newHash types.Hash) (bool, error) {
	if fheos.State == nil {
		return false, nil
	}

	multiStore := storage2.NewMultiStore(h.evm.CiphertextDb, &fheos.State.Storage)

	if fhe.IsCtHash(newHash) && !fhe.IsTriviallyEncryptedCtHash(newHash) {
		err := multiStore.ReferenceCiphertext(newHash)
		// Check whether the value is not a newly created value (i.e. not in memory but is in store)
		if err != nil {
			log.Warn("Failed to reference a ciphertext", "err", err)
		}
	}

	originalHash := types.Hash(original)
	if fhe.IsCtHash(originalHash) && !fhe.IsTriviallyEncryptedCtHash(originalHash) {
		// if the original hash is not empty, we are updating a value
		// we need to dereference the old value from the storage
		wasDeleted, err := multiStore.DereferenceCiphertext(originalHash)
		if err != nil {
			log.Error("Failed to dereference old ciphertext", "err", err)
			return false, err
		}

		log.Info("Dereferenced old ciphertext", "hash", hex.EncodeToString(originalHash[:]))
		// We want to mark the previously commited ct for persistence only if it wasn't deleted (RefCount == 0)
		return !wasDeleted, nil
	}

	return false, nil
}

// StoreCiphertextHook The purpose of this hook is to mark the ciphertext as LTS if the tx is successful and update reference counts
// contract - The address of the contract in which the ciphertext is stored
// loc - the location (starting from 0) in the storage of the contract
// commited - The previous value (ct hash) that was present in the said location
// val - The new value that is being stored
func (h *FheOSHooksImpl) StoreCiphertextHook(contract common.Address, loc [32]byte, commited common.Hash, val [32]byte) error {
	// marks the ciphertext as lts - should be stored in long term storage when/if the tx is successful
	// option - better to flush all at the end of the tx from memdb, or define a memdb in fheos that is flushed at the end of the tx?
	storage := storage2.NewEphemeralStorage(h.evm.CiphertextDb)

	wasDereferenced, err := h.updateCiphertextReferences(commited, val)
	if err != nil {
		return err
	}

	if wasDereferenced {
		err = storage.MarkForPersistence(contract, types.Hash(commited))
		if err != nil {
			log.Error("Error marking dereferenced ciphertext as LTS", "err", err)
			return err
		}
	}

	if !fhe.IsCtHash(val) {
		return nil
	}

	// Skip for values who are not in memory as there is nothing to persist
	if !storage.HasCt(val) {
		return nil
	}

	err = storage.MarkForPersistence(contract, val)
	if err != nil {
		log.Error("Error marking ciphertext as LTS", "err", err)
		return err
	}

	return nil
}

const ExtraGasCost = 0

func (h *FheOSHooksImpl) StoreGasHook(contract common.Address, loc [32]byte, val [32]byte) (uint64, uint64) {
	return ExtraGasCost, 0
}

func (h *FheOSHooksImpl) LoadCiphertextHook() [32]byte {
	//storage := storage2.NewMultiStore(h.evm.CiphertextDb, &fheos.State.Storage)
	// checks if ciphertext hash is already known (should be in either memory storage or long term storage)
	// doesn't do anything right now
	return [32]byte{}
}

func (h *FheOSHooksImpl) EvmCallStart() {
	// don't really need this? Or maybe to start a new ephemeral storage?
	// But how do we know how to keep the context thread safe? Ugh, do we need 2 dbs now?

	if h.evm.Commit {
		lock.Lock()
	} else {
		lock.RLock()
	}
}

func (h *FheOSHooksImpl) EvmCallEnd(evmSuccess bool) {
	if evmSuccess && h.evm.Commit {
		storage := storage2.NewEphemeralStorage(h.evm.CiphertextDb)

		toStore := storage.GetAllToPersist()
		toDelete := storage.GetAllToDelete()

		for _, contractCiphertext := range toStore {
			cipherText, err := storage.GetCt(contractCiphertext.CipherTextHash)
			if err != nil {
				// this should probably be a part of the "consensus" - i.e. if we fuck up somewhere and somehow you can reach this failure just revert the tx
				// if we actually quit here all nodes will halt and suckiness will ensue. Though the only way you get here is if we fuck up and let you mark a
				// ciphertext as LTS without it being in memory
				// right now the hook gets called after the evm execution, so I'm not sure that reverting is possible - but we can also probably move this to be
				// inside the evm
				log.Error("Error getting ciphertext from storage when trying to store in lts - state corruption detected", "hash", hex.EncodeToString(contractCiphertext.CipherTextHash[:]), "err", err)
				continue
			}
			err = fheos.State.SetCiphertext(cipherText)
			if err != nil {
				log.Error("Error storing ciphertext in LTS - state corruption detected", "err", err)
			}

			log.Info("Added ct to store", "hash", (*fhe.FheEncrypted)(cipherText.Data).GetHash().Hex())
		}

		for hash, _ := range toDelete {
			// TODO: temporarily disabling ct deletion
			// err := fheos.State.Storage.DeleteCt(hash)
			// if err != nil {
			// 	// Deletion failure, bummer but nothing to be worried about
			// 	log.Error("Failed to delete ciphertext", "hash", hex.EncodeToString(hash[:]))
			// }

			log.Info("Deleted ciphertext", "hash", hex.EncodeToString(hash[:]))
		}
	}

	if fheos.State != nil {
		fheos.State.RandomCounter = 0
	}

	if h.evm.Commit {
		lock.Unlock()
	} else {
		lock.RUnlock()
	}
}

func shouldIgnoreContract(caller common.Address, addr common.Address) bool {
	// Address of a user and not a contract
	const NilAddress = "0x0000000000000000000000000000000000000000"
	const FheosPrecompilesAddress = "0x0000000000000000000000000000000000000080"
	userAddress := common.HexToAddress(NilAddress)
	// Address of our precompiled contract - just to be sure we don't waste time on it
	precompilesAddress := common.HexToAddress(FheosPrecompilesAddress)

	if caller.Cmp(userAddress) == 0 {
		return true
	}

	if addr.Cmp(precompilesAddress) == 0 {
		return true
	}

	return false
}

func (h *FheOSHooksImpl) iterateHashes(data []byte, dataType string, owner common.Address, newOwner common.Address) {
	// iterate through the data and check if the hash is a ciphertext hash
	// if it is, add the owner to the ciphertext
	// if not, continue

	const EvmVariableLen = 32
	dataLen := len(data)
	if dataLen%32 != 0 {
		log.Warn("Data is not aligned to 32 bytes", "type", dataType, "length", dataLen)
	}

	paramsCount := dataLen / EvmVariableLen
	var hash [EvmVariableLen]byte

	state := fheos.State
	if state == nil {
		log.Warn("Fheos state is not initialized (can be ignored if it is a part of the unittests)")
		return
	}

	storage := storage2.NewMultiStore(h.evm.CiphertextDb, &state.Storage)

	for i := 0; i < paramsCount; i++ {
		offset := i * EvmVariableLen
		copy(hash[:], data[offset:offset+EvmVariableLen])

		if !fhe.IsCtHash(hash) {
			continue
		}

		// will return ct if hash exists AND if caller is one of the owners
		// otherwise we have nothing to do anymore
		ct, _ := storage.GetCtRepresentation(hash, owner)
		if ct == nil {
			continue
		}

		_ = storage.AddOwner(hash, ct, newOwner)

		log.Info("Contract has been added as an owner to the ciphertext", "contract", newOwner, "ciphertext", hex.EncodeToString(hash[:]))
	}
}

// ContractCall The purpose of this hook is to be able to pass ownership for a ciphertext to the contract that has been called if the caller is an owner
// The function parses the input for ciphertexts and pass ownership for each ciphertext
func (h *FheOSHooksImpl) ContractCall(isSimulation bool, callType int, caller common.Address, addr common.Address, input []byte) {
	// Input is built as the following:
	//  first 4 bytes are indicating what is the function that is being called
	// 	from now on every param is a 32 byte value
	//  if the param is dynamically sized (string, bytes, etc.) the 32 bytes will only indicate the offset of the actual value in "input"
	//  when going to offset you will find 32 bytes indicating the length of the value
	//  and then the value itself, each value in the array is padded to 32 bytes

	// Skip delegate calls - The owner remains the same
	if callType == vm.CallTypeDelegateCall {
		return
	}

	// Skip this logic in simulations for now as it won't affect the gas estimation code
	// FHENIX: When implementing sync decryption should remove this
	if isSimulation {
		return
	}

	if shouldIgnoreContract(caller, addr) {
		return
	}

	if len(input) <= 4 {
		return
	}

	h.iterateHashes(input[4:], "input", caller, addr)
}

func (h *FheOSHooksImpl) ContractCallReturn(isSimulation bool, callType int, caller common.Address, addr common.Address, output []byte) {
	// If a contract returns a value, we should check if it contains any ciphertexts
	// If so, we should pass ownership of the ciphertexts to the caller

	// Skip this logic in simulations for now as it won't affect the gas estimation code
	// FHENIX: When implementing sync decryption should remove this
	if isSimulation {
		return
	}

	if shouldIgnoreContract(caller, addr) {
		return
	}

	h.iterateHashes(output, "input", addr, caller)
}

func NewFheOSHooks(evm *vm.EVM) *FheOSHooksImpl {
	return &FheOSHooksImpl{
		evm: evm,
	}
}
