package hooks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

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

// StoreCiphertextHook The purpose of this hook is to mark the ciphertext as LTS if the tx is successful and update reference counts
// contract - The address of the contract in which the ciphertext is stored
// loc - the location (starting from 0) in the storage of the contract
// committed - The previous value (ct hash) that was present in the said location
// val - The new value that is being stored
func (h *FheOSHooksImpl) StoreCiphertextHook(_ common.Address, _ [32]byte, _ common.Hash, _ [32]byte) error {
	return nil
}

const ExtraGasCost = 0

func (h *FheOSHooksImpl) StoreGasHook(contract common.Address, loc [32]byte, val [32]byte) (uint64, uint64) {
	return ExtraGasCost, 0
}

func (h *FheOSHooksImpl) LoadCiphertextHook() [32]byte {
	//storage := storage2.NewMultiStore(h.evm.CiphertextDb, &fheos.State.Storage)
	// checks if ciphertext hash is already known (should be in either memory storage or long-term storage)
	// doesn't do anything right now
	return [32]byte{}
}

func (h *FheOSHooksImpl) EvmCallStart() {
	// don't really need this? Or maybe to start a new ephemeral storage?
	// But how do we know how to keep the context thread safe? Ugh, do we need 2 dbs now?
}

func (h *FheOSHooksImpl) EvmCallEnd(evmSuccess bool) {
}

// ContractCall The purpose of this hook is to be able to pass ownership for a ciphertext to the contract that has been called if the caller is an owner
// The function parses the input for ciphertexts and passes ownership for each ciphertext
func (h *FheOSHooksImpl) ContractCall(isSimulation bool, callType int, caller common.Address, addr common.Address, input []byte) {
}

func (h *FheOSHooksImpl) ContractCallReturn(isSimulation bool, callType int, caller common.Address, addr common.Address, output []byte) {
}

func NewFheOSHooks(evm *vm.EVM) *FheOSHooksImpl {
	return &FheOSHooksImpl{
		evm: evm,
	}
}
