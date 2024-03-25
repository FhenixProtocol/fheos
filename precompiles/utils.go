package precompiles

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

type TxParams struct {
	Commit        bool
	GasEstimation bool
	EthCall       bool
}

type GasBurner interface {
	Burn(amount uint64) error
	Burned() uint64
}

func TxParamsFromEVM(evm *vm.EVM) TxParams {
	var tp TxParams
	tp.Commit = evm.Commit
	tp.GasEstimation = evm.GasEstimation
	tp.EthCall = evm.EthCall

	return tp
}

type Precompile struct {
	Metadata *bind.MetaData
	Address  common.Address
}

func getCiphertext(state *FheosState, ciphertextHash fhe.Hash) *fhe.FheEncrypted {
	ct, err := state.GetCiphertext(types.Hash(ciphertextHash))
	if err != nil {
		logger.Error("reading ciphertext from state resulted with error: ", err)
		return nil
	}

	return (*fhe.FheEncrypted)(ct)
}
func get2VerifiedOperands(state *FheosState, lhsHash []byte, rhsHash []byte) (lhs *fhe.FheEncrypted, rhs *fhe.FheEncrypted, err error) {
	if len(lhsHash) != 32 || len(rhsHash) != 32 {
		return nil, nil, errors.New("ciphertext's hashes need to be 32 bytes long")
	}

	lhs = getCiphertext(state, fhe.BytesToHash(lhsHash))
	if lhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	rhs = getCiphertext(state, fhe.BytesToHash(rhsHash))
	if rhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func get3VerifiedOperands(state *FheosState, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte) (control *fhe.FheEncrypted, ifTrue *fhe.FheEncrypted, ifFalse *fhe.FheEncrypted, err error) {
	if len(controlHash) != 32 || len(ifTrueHash) != 32 || len(ifFalseHash) != 32 {
		return nil, nil, nil, errors.New("ciphertext's hashes need to be 32 bytes long")
	}

	control = getCiphertext(state, fhe.BytesToHash(controlHash))
	if control == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifTrue = getCiphertext(state, fhe.BytesToHash(ifTrueHash))
	if ifTrue == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifFalse = getCiphertext(state, fhe.BytesToHash(ifFalseHash))
	if ifFalse == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func importCiphertext(state *FheosState, ct *fhe.FheEncrypted) error {
	err := state.SetCiphertext(ct)
	if err != nil {
		logger.Error("failed importing ciphertext to state: ", err)
		return err
	}

	return nil
}
func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func evaluateRequire(ct *fhe.FheEncrypted) bool {
	return fhe.Require(ct)
}

type fheUintType uint8

const (
	FheUint8  fheUintType = 0
	FheUint16 fheUintType = 1
	FheUint32 fheUintType = 2

	FheBool fheUintType = 13
)

func isValidType(t byte) bool {
	return t >= uint8(FheUint8) && t <= uint8(FheBool)
}
