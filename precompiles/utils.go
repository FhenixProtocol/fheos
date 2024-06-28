package precompiles

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	"math"
	"math/big"
)

type TxParams struct {
	Commit          bool
	GasEstimation   bool
	EthCall         bool
	CiphertextDb    *memorydb.Database
	ContractAddress common.Address
}

type GasBurner interface {
	Burn(amount uint64) error
	Burned() uint64
}

func TxParamsFromEVM(evm *vm.EVM, callerContract common.Address) TxParams {
	var tp TxParams
	tp.Commit = evm.Commit
	tp.GasEstimation = evm.GasEstimation
	tp.EthCall = evm.EthCall

	tp.CiphertextDb = evm.CiphertextDb
	tp.ContractAddress = callerContract

	return tp
}

type Precompile struct {
	Metadata *bind.MetaData
	Address  common.Address
}

func getCiphertext(state *storage.MultiStore, ciphertextHash fhe.Hash, caller common.Address) *fhe.FheEncrypted {
	ct, err := state.GetCt(types.Hash(ciphertextHash), caller)
	if err != nil {

		logger.Error("reading ciphertext from State resulted with error", "hash", ciphertextHash.Hex(), "error", err.Error())
		return nil
	}

	return (*fhe.FheEncrypted)(ct)
}
func get2VerifiedOperands(storage *storage.MultiStore, lhsHash []byte, rhsHash []byte, caller common.Address) (lhs *fhe.FheEncrypted, rhs *fhe.FheEncrypted, err error) {
	if len(lhsHash) != 32 || len(rhsHash) != 32 {
		return nil, nil, errors.New("ciphertext's hashes need to be 32 bytes long")
	}

	lhs = getCiphertext(storage, fhe.BytesToHash(lhsHash), caller)
	if lhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	rhs = getCiphertext(storage, fhe.BytesToHash(rhsHash), caller)
	if rhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func get3VerifiedOperands(storage *storage.MultiStore, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte, caller common.Address) (control *fhe.FheEncrypted, ifTrue *fhe.FheEncrypted, ifFalse *fhe.FheEncrypted, err error) {
	if len(controlHash) != 32 || len(ifTrueHash) != 32 || len(ifFalseHash) != 32 {
		return nil, nil, nil, errors.New("ciphertext's hashes need to be 32 bytes long")
	}

	control = getCiphertext(storage, fhe.BytesToHash(controlHash), caller)
	if control == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifTrue = getCiphertext(storage, fhe.BytesToHash(ifTrueHash), caller)
	if ifTrue == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifFalse = getCiphertext(storage, fhe.BytesToHash(ifFalseHash), caller)
	if ifFalse == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func storeCipherText(storage *storage.MultiStore, ct *fhe.FheEncrypted, owner common.Address) error {
	err := storage.AppendCt(types.Hash(ct.Hash()), (*types.FheEncrypted)(ct), owner)
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

func FakeDecryptionResult(encType fhe.EncryptionType) *big.Int {

	switch encType {
	case fhe.Uint8:
		return big.NewInt(math.MaxUint8 / 2)
	case fhe.Uint16:
		return big.NewInt(math.MaxInt16 / 2)
	case fhe.Uint32:
		return big.NewInt(math.MaxUint32 / 2)
	case fhe.Uint64:
		return big.NewInt(math.MaxUint64 / 2)
	case fhe.Uint128:
		value := &big.Int{}
		value.SetString("2222222222222222222222222222222", 16)
		return value
	case fhe.Uint256:
		value := &big.Int{}
		value.SetString("2222222222222222222222222222222222222222222222222", 16)
		return value
	case fhe.Address:
		value := &big.Int{}
		value.SetString("Dd4BEac65bad064932FB21aE7Ba2aa6e8fbc41A4", 16)
		return value
	default:
		return big.NewInt(0)
	}
}
