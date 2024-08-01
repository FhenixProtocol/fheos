package precompiles

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	"math"
	"math/big"
	"os"
	"time"
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
	ErrChannel      chan error
}

type NotPlaceholderKeyError struct{}

func (err *NotPlaceholderKeyError) Error() string {
	return "Placeholder key in incorrect format."
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
	tp.ErrChannel = evm.ErrorChannel
	return tp
}

type Precompile struct {
	Metadata *bind.MetaData
	Address  common.Address
}

func InitLogger() {
	logger = log.Root().New("module", "fheos")
	warpDriveLogger = log.Root().New("module", "warp-drive")
	fhe.SetLogger(warpDriveLogger)
}

func InitFheConfig(fheConfig *fhe.Config) error {
	// Itzik: I'm not sure if this is the right way to initialize the logger
	handler := log.NewTerminalHandlerWithLevel(os.Stderr, log.FromLegacyLevel(fheConfig.LogLevel), true)
	glogger := log.NewGlogHandler(handler)

	logger = log.NewLogger(glogger)
	fhe.SetLogger(log.NewLogger(glogger))

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
	}

	// can speed this up to be concurrent, but for now this is fine I guess?

	lhsValue = awaitCtResult(storage, lhsHash, tp)
	rhsValue = awaitCtResult(storage, rhsHash, tp)
	//rhsValue = getCiphertext(storage, fhe.Hash(rhsHash), tp.ContractAddress)
	//for rhsValue.IsPlaceholderValue() {
	//	rhsValue = getCiphertext(storage, fhe.Hash(rhsHash), tp.ContractAddress)
	//	time.Sleep(1 * time.Millisecond)
	//}

	//var lhsCheck, rhsCheck [32]byte
	////copy(lhsCheck[:], lhsValue.Data)
	////copy(rhsCheck[:], rhsAddress.Data)
	////for now like this but needs a refactor
	//for (lhsValue.IsPlaceholderValue() || (rhsValue.IsPlaceholderValue() && !fhe.IsCtHash(rhsCheck)) {
	//	rhsValue = getCiphertext(storage, fhe.Hash(rhsHash), tp.ContractAddress)
	//	//copy(lhsCheck[:], lhsAddress.Data)
	//	//copy(rhsCheck[:], rhsAddress.Data)
	//}
	//var lhsData []byte = lhsHash
	//var rhsData []byte = rhsHash
	//

	//if fhe.IsPlaceholderValue(lhsHash) {
	//	copy(lhsData, lhsCheck[:])
	//}
	//if fhe.IsPlaceholderValue(rhsHash) {
	//	copy(rhsData, rhsCheck[:])
	//}
	return lhsValue, rhsValue
}

func awaitCtResult(storage *storage.MultiStore, lhsHash []byte, tp *TxParams) *fhe.FheEncrypted {
	lhsValue := getCiphertext(storage, fhe.Hash(lhsHash), tp.ContractAddress)
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

func storePlaceholderValue(storage *storage.MultiStore, key types.Hash, ct *fhe.FheEncrypted, owner common.Address) error {
	var err error
	if fhe.IsPlaceholderValue(key[:]) {
		err = storage.AppendPhV(key, (*types.FheEncrypted)(ct), owner)
	} else {
		err = &NotPlaceholderKeyError{}
	}
	if err != nil {
		logger.Error("failed importing placeholder value to storage: ", err)
		return err
	}
	return nil
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
func CopySlice(original []byte) []byte {
	copied := make([]byte, len(original))
	copy(copied, original)
	return copied
}
