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

func blockUntilInputsAvailable(storage *storage.MultiStore, tp *TxParams, inputKeys ...fhe.CiphertextKey) ([]*fhe.FheEncrypted, error) {
	cts := make([]*fhe.FheEncrypted, len(inputKeys))
	results := make(chan struct {
		index int
		ct    *fhe.FheEncrypted
	}, len(inputKeys))

	// Launch goroutines for each hash
	for i, key := range inputKeys {
		go func(index int, key fhe.CiphertextKey) {
			ct := awaitCtResult(storage, key.Hash, tp)
			results <- struct {
				index int
				ct    *fhe.FheEncrypted
			}{index, ct}
		}(i, key)
	}

	// Collect results
	for i := 0; i < len(inputKeys); i++ {
		result := <-results
		cts[result.index] = result.ct
		if result.ct == nil {
			return nil, errors.New("unverified ciphertext handle, hash: " + fhe.Hash(inputKeys[result.index].Hash).Hex())
		}
	}

	return cts, nil
}

func awaitCtResult(storage *storage.MultiStore, lhsHash fhe.Hash, tp *TxParams) *fhe.FheEncrypted {
	lhsValue := getCiphertext(storage, lhsHash)
	if lhsValue == nil {
		return nil
	}

	for lhsValue.IsPlaceholderValue() {
		lhsValue = getCiphertext(storage, lhsHash)
		time.Sleep(1 * time.Millisecond)
	}
	return lhsValue
}

func getCiphertext(state *storage.MultiStore, ciphertextHash fhe.Hash) *fhe.FheEncrypted {
	logger.Error("LIORRRRRR getting ciphertext", "hash", ciphertextHash.Hex())
	ct, err := state.GetCt(types.Hash(ciphertextHash))
	if err != nil {
		logger.Error("reading ciphertext from State resulted with error", "hash", ciphertextHash.Hex(), "error", err.Error())
		return nil
	}

	return (*fhe.FheEncrypted)(ct)
}

func storeCipherText(storage *storage.MultiStore, ct *fhe.FheEncrypted) error {
	logger.Error("LIORRRRRR Putting ciphertext", "hash", ct.GetHash().Hex())
	err := storage.PutCtIfNotExist(types.Hash(ct.GetHash()), (*types.FheEncrypted)(ct))
	if err != nil {
		logger.Error("failed importing ciphertext to state: ", err)
		return err
	}

	return nil
}

func ByteToUint256(b byte) []byte {
	var uint256 [32]byte
	uint256[31] = b // Place the byte in the least significant position
	return uint256[:]
}

func Int32ToUint256(i int32) []byte {
	var uint256 [32]byte

	// Convert the int32 value to an unsigned 4-byte slice
	var bytes [4]byte
	binary.BigEndian.PutUint32(bytes[:], uint32(i))

	// Copy the 4-byte slice into the least significant part of the 32-byte array
	copy(uint256[28:], bytes[:]) // Place at the end (big-endian)

	return uint256[:]
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
