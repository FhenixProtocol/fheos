package precompiles

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

func UtypeToString(utype byte) string {
	switch fhe.EncryptionType(utype) {
	case fhe.Uint8:
		return "uint8"
	case fhe.Uint16:
		return "uint16"
	case fhe.Uint32:
		return "uint32"
	case fhe.Uint64:
		return "uint64"
	case fhe.Uint128:
		return "uint128"
	case fhe.Uint256:
		return "uint256"
	case fhe.Address:
		return "address"
	case fhe.Bool:
		return "bool"
	default:
		return "unknown"
	}
}

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

func shouldPrintPrecompileInfo(tp *TxParams) bool {
	return tp.Commit && !tp.GasEstimation
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

func InitLogger() {
	// default log level "debug", we don't have access to args here
	logLevel := 4
	handler := log.NewTerminalHandlerWithLevel(os.Stderr, log.FromLegacyLevel(logLevel), true)
	glogger := log.NewGlogHandler(handler)
	glogger.Verbosity(log.FromLegacyLevel(logLevel))
	logger = log.NewLogger(glogger).New("module", "fheos")
	fhe.SetLogger(log.NewLogger(glogger).New("module", "warp-drive"))
}

func InitFheConfig(fheConfig *fhe.Config) error {
	// Itzik: I'm not sure if this is the right way to initialize the logger
	handler := log.NewTerminalHandlerWithLevel(os.Stderr, log.FromLegacyLevel(fheConfig.LogLevel), true)
	glogger := log.NewGlogHandler(handler)
	glogger.Verbosity(log.FromLegacyLevel(fheConfig.LogLevel))

	logger = log.NewLogger(glogger).New("module", "fheos")
	fhe.SetLogger(log.NewLogger(glogger).New("module", "warp-drive"))

	err := fhe.Init(fheConfig)

	if err != nil {
		logger.Error("Failed to init fhe config with", "error:", err)
		return err
	}

	logger.Info("Successfully initialized fhe config", "config", fheConfig)

	return nil
}

func InitFheos(tfheConfig *fhe.Config) error {
	err := InitFheConfig(tfheConfig)
	if err != nil {
		return err
	}

	err = InitializeFheosState()
	if err != nil {
		return err
	}

	return nil
}

func CreatePlaceHolderData() []byte {
	return make([]byte, 32)[:]
}

func blockUntilBinaryOperandsAvailable(storage *storage.MultiStore, lhsHash, rhsHash []byte, tp *TxParams) (*fhe.FheEncrypted, *fhe.FheEncrypted) {
	var lhsValue *fhe.FheEncrypted
	var rhsValue *fhe.FheEncrypted

	if !fhe.IsCtHash([32]byte(lhsHash)) || !fhe.IsCtHash([32]byte(rhsHash)) {
		// return error
		return nil, nil
	}

	// can speed this up to be concurrent, but for now this is fine I guess?
	lhsValue = awaitCtResult(storage, lhsHash, tp)
	rhsValue = awaitCtResult(storage, rhsHash, tp)

	return lhsValue, rhsValue
}

func awaitCtResult(storage *storage.MultiStore, lhsHash []byte, tp *TxParams) *fhe.FheEncrypted {
	lhsValue := getCiphertext(storage, fhe.Hash(lhsHash), tp.ContractAddress)
	if lhsValue == nil {
		return nil
	}

	for lhsValue.IsPlaceholderValue() {
		lhsValue = getCiphertext(storage, fhe.Hash(lhsHash), tp.ContractAddress)
		time.Sleep(1 * time.Millisecond)
	}
	return lhsValue
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

func get3VerifiedOperands(storage *storage.MultiStore, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte, tp *TxParams) (control *fhe.FheEncrypted, ifTrue *fhe.FheEncrypted, ifFalse *fhe.FheEncrypted, err error) {
	if len(controlHash) != 32 || len(ifTrueHash) != 32 || len(ifFalseHash) != 32 {
		return nil, nil, nil, errors.New("ciphertext's hashes need to be 32 bytes long")
	}

	control = awaitCtResult(storage, controlHash, tp)
	if control == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifTrue = awaitCtResult(storage, ifTrueHash, tp)
	if ifTrue == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifFalse = awaitCtResult(storage, ifFalseHash, tp)
	if ifFalse == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func storeCipherText(storage *storage.MultiStore, ct *fhe.FheEncrypted, owner common.Address) error {
	err := storage.AppendCt(types.Hash(ct.GetHash()), (*types.FheEncrypted)(ct), owner)
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
func CopySlice(original []byte) []byte {
	copied := make([]byte, len(original))
	copy(copied, original)
	return copied
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
