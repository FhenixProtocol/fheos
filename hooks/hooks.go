package hooks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	fheos "github.com/fhenixprotocol/fheos/precompiles"
	storage2 "github.com/fhenixprotocol/fheos/storage"
)

type FheOSHooks interface {
	StoreCiphertextHook(contract common.Address, loc [32]byte, val [32]byte) error
	StoreGasHook(contract common.Address, loc [32]byte, val [32]byte) (uint64, uint64)
	LoadCiphertextHook() [32]byte
	EvmCallStart()
	EvmCallEnd(evmSuccess bool)
	ContractCall(callType int, caller common.Address, addr common.Address, input []byte)
}

type FheOSHooksImpl struct {
	evm *vm.EVM
}

func (h FheOSHooksImpl) StoreCiphertextHook(contract common.Address, loc [32]byte, val [32]byte) error {
	// marks the ciphertext as lts - should be stored in long term storage when/if the tx is successful
	// option - better to flush all at the end of the tx from memdb, or define a memdb in fheos that is flushed at the end of the tx?
	storage := storage2.NewEphemeralStorage(h.evm.CiphertextDb)

	// if this value isn't in our storage - i.e. isn't a ciphertext - we just noop
	if !storage.HasCt(val) {
		return nil
	}

	err := storage.MarkForPersistence(contract, val)
	if err != nil {
		log.Crit("Error marking ciphertext as LTS", "err", err)
		return err
	}

	return nil
}

const ExtraGasCost = 0

func (h FheOSHooksImpl) StoreGasHook(contract common.Address, loc [32]byte, val [32]byte) (uint64, uint64) {
	return ExtraGasCost, 0
}

func (h FheOSHooksImpl) LoadCiphertextHook() [32]byte {
	//storage := storage2.NewMultiStore(h.evm.CiphertextDb, &fheos.State.Storage)
	// checks if ciphertext hash is already known (should be in either memory storage or long term storage)
	// doesn't do anything right now
	return [32]byte{}
}

func (h FheOSHooksImpl) EvmCallStart() {
	// don't really need this? Or maybe to start a new ephemeral storage?
	// But how do we know how to keep the context thread safe? Ugh, do we need 2 dbs now?
}

func (h FheOSHooksImpl) EvmCallEnd(evmSuccess bool) {
	if evmSuccess && h.evm.Commit {
		storage := storage2.NewEphemeralStorage(h.evm.CiphertextDb)
		toStore := storage.GetAllToPersist()

		for _, contractCiphertext := range toStore {
			cipherText, err := storage.GetCt(contractCiphertext.CipherTextHash)
			if err != nil {
				// this should probably be a part of the "consensus" - i.e. if we fuck up somewhere and somehow you can reach this failure just revert the tx
				// if we actually quit here all nodes will halt and suckiness will ensue. Though the only way you get here is if we fuck up and let you mark a
				// ciphertext as LTS without it being in memory
				// right now the hook gets called after the evm execution, so I'm not sure that reverting is possible - but we can also probably move this to be
				// inside the evm
				log.Crit("Error getting ciphertext from storage when trying to store in lts - state corruption detected", "err", err)
				continue
			}
			err = fheos.State.SetCiphertext(cipherText)
			if err != nil {
				log.Crit("Error storing ciphertext in LTS - state corruption detected", "err", err)
			}
		}
	}
}

func isCiphertextHash(param [32]byte) bool {
	// Currently we have no indication if 32bytes are representing ciphertext hash or not
	// We are filtering out those params who start with at least 6 zeroes as they are PROBABLY!! size indication and not hashes
	// In the future we might add "deadbeaf" to indicate if hash is a ciphertext, this will be changes accordingly
	// FHENIX: If there is a problem with ciphertext ownership, check for collision here

	// check if param starts with at least 6 zeroes
	for i := 0; i < 6; i++ {
		if param[i] != 0 {
			return true
		}
	}

	return false
}

// ContractCall The purpose of this hook is to be able to pass ownership for a ciphertext to the contract that has been called if the caller is an owner
// The function parses the input for ciphertexts and pass ownership for each ciphertext
func (h FheOSHooksImpl) ContractCall(callType int, caller common.Address, addr common.Address, input []byte) {
	// Input is built as the following:
	//  first 4 bytes are indicating what is the function the is being called
	// 	from now on every param is a 32 byte value
	//  if the param is dynamically sized (string, bytes, etc.) the 32 bytes will only indicate the offset of the actual value in "input"
	//  when going to offset you will find 32 bytes indicating the length of the value
	//  and then the value itself, each value in the array is padded to 32 bytes

	// Skip delegate calls??????????
	if callType == vm.CallTypeDelegateCall {
		return
	}

	inputLen := len(input)
	if inputLen <= 4 {
		return
	}

	inputLen -= 4

	if inputLen%32 != 0 {
		log.Warn("Invalid input length for contract call", "input", input)
	}

	paramsCount := inputLen / 32
	var hash [32]byte
	for i := 0; i < paramsCount; i++ {
		offset := 4 + i*32
		copy(hash[:], input[offset:offset+32])

		if !isCiphertextHash(hash) {
			continue
		}

		log.Error("LIORRRRRRRRR found a ciphertext hash", "hash", hash)

		storage := storage2.NewMultiStore(h.evm.CiphertextDb, &fheos.State.Storage)
		// will return ct if hash exists AND if caller is one of the owners
		// otherwise we have nothing to do anymore
		ct, _ := storage.GetCtRepresentation(hash, caller)
		if ct == nil {
			continue
		}

		_ = storage.AddOwner(hash, ct, addr)

		log.Info("Contract has been added as an owner to the ciphertext", "contract", addr, "ciphertext", hash)
	}

}

func NewFheOSHooks(evm *vm.EVM) FheOSHooksImpl {
	return FheOSHooksImpl{
		evm: evm,
	}
}
