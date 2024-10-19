package precompiles

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

type TxParams struct {
	Commit          bool
	GasEstimation   bool
	EthCall         bool
	CiphertextDb    *memorydb.Database
	ContractAddress common.Address
	GetBlockHash    vm.GetHashFunc
	BlockNumber     *big.Int
	ParallelTxHooks types.ParallelTxProcessingHook
	vm.TxContext
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
	tp.BlockNumber = evm.Context.BlockNumber
	tp.GetBlockHash = evm.Context.GetHash

	// If this is running in a sequencer, this should not be nil
	if parallelHook, ok := evm.ProcessingHook.(types.ParallelTxProcessingHook); ok {
		tp.ParallelTxHooks = parallelHook
	} else {
		tp.ParallelTxHooks = nil
	}

	tp.TxContext = evm.TxContext

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

func evaluateRequire(ct *fhe.FheEncrypted) (bool, error) {
	return fhe.Require(ct)
}

func GenerateSeedFromEntropy(contractAddress common.Address, hash common.Hash, randomCounter uint64) uint64 {
	data := make([]byte, 0, len(contractAddress)+len(hash)+8) // 8 bytes for uint64

	data = append(data, contractAddress[:]...)
	data = append(data, hash[:]...)

	uint64Bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(uint64Bytes, randomCounter)
	data = append(data, uint64Bytes...)

	hashResult := crypto.Keccak256Hash(data)

	result := binary.LittleEndian.Uint64(hashResult[:])
	logger.Debug(fmt.Sprintf("generated seed for random: %d", result))
	return result
}

// SAFETY NOTE: this function assumes input length validity (i.e. that ctHash and pk are 32 bytes long)
// since the SealOutput precompile is doing these checks before calling this function. Be extra careful
// when using this function in other places.
func genSealedKey(ctHash, pk []byte, functionName types.PrecompileName) types.PendingDecryption {
	var hash [32]byte
	for i := 0; i < 32; i++ {
		// Assumes input length validity
		hash[i] = ctHash[i] ^ pk[i]
	}

	return types.PendingDecryption{
		Hash: hash,
		Type: functionName,
	}
}
