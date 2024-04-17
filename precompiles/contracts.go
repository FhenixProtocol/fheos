package precompiles

import (
	"encoding/hex"
	"fmt"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	"math/big"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
)

var logger log.Logger

func init() {
	InitLogger()
}

func InitLogger() {
	logger = log.Root().New("module", "fheos")
	fhe.SetLogger(log.Root().New("module", "go-tfhe"))
}

func InitFheConfig(fheConfig *fhe.Config) error {
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

func shouldPrintPrecompileInfo(tp *TxParams) bool {
	return tp.Commit && !tp.GasEstimation
}

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
	case fhe.Bool:
		return "bool"
	default:
		return "unknown"
	}
}

// ============================================================

func Log(s string, tp *TxParams) (uint64, error) {
	if tp.GasEstimation {
		return 1, nil
	}

	logger.Debug(fmt.Sprintf("Contract Log: %s", s))
	return 1, nil
}

func Add(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "add"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType || (lhs.UintType != uintType) {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.

	result, err := lhs.Add(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

// Verify takes inputs from the user and runs them through verification. Note that we will always get ciphertexts that
// are public-key encrypted and compressed. Anything else will fail
func Verify(utype byte, input []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "verify"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	ct := fhe.NewFheEncryptedFromBytes(input, uintType, true, false /* TODO: not sure + shouldn't be hardcoded */)
	err := ct.Verify()
	if err != nil {
		logger.Info(fmt.Sprintf("failed to verify ciphertext %s for type %d - was input corrupted?", ct.Hash().Hex(), uintType))
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, &ct)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	logger.Debug(functionName+" success", "ctHash", ct.Hash().Hex())
	return ct.GetHashBytes(), gas, nil
}

func SealOutput(utype byte, ctHash []byte, pk []byte, tp *TxParams) (string, uint64, error) {
	//solgen: bool math
	functionName := "sealOutput"
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return "", 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
    return "0x" + strings.Repeat("00", 100), gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	if len(ctHash) != 32 {
		msg := functionName + " ciphertext's hashes need to be 32 bytes long"
		logger.Error(msg, "ciphertext-hash", hex.EncodeToString(ctHash), "hash-len", len(ctHash))
		return "", 0, vm.ErrExecutionReverted
	}

	if len(pk) != 32 {
		msg := functionName + " public key need to be 32 bytes long"
		logger.Error(msg, "public-key", hex.EncodeToString(pk), "len", len(pk))
		return "", 0, vm.ErrExecutionReverted
	}

	ct := getCiphertext(storage, fhe.BytesToHash(ctHash))
	if ct == nil {
		msg := functionName + " unverified ciphertext handle"
		logger.Error(msg, "ciphertext-hash", hex.EncodeToString(ctHash))
		return "", 0, vm.ErrExecutionReverted
	}

	reencryptedValue, err := fhe.SealOutput(*ct, pk)
	if err != nil {
		logger.Error(functionName+" failed to encrypt to user key", "err", err)
		return "", 0, vm.ErrExecutionReverted
	}

	logger.Debug(functionName+" success", "ciphertext-hash ", hex.EncodeToString(ctHash), "public-key", hex.EncodeToString(pk))

	return "0x" + hex.EncodeToString(reencryptedValue), gas, nil
}

func Decrypt(utype byte, input []byte, tp *TxParams) (*big.Int, uint64, error) {
	//solgen: output plaintext
	functionName := "decrypt"
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return FakeDecryptionResult(uintType), gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	if len(input) != 32 {
		msg := functionName + " input len must be 32 bytes"
		logger.Error(msg, "input", hex.EncodeToString(input), "len", len(input))
		return nil, 0, vm.ErrExecutionReverted
	}

	ct := getCiphertext(storage, fhe.BytesToHash(input))
	if ct == nil {
		msg := functionName + " unverified ciphertext handle"
		logger.Error(msg, " input ", hex.EncodeToString(input))
		return nil, 0, vm.ErrExecutionReverted
	}

	decryptedValue, err := fhe.Decrypt(*ct)
	if err != nil {
		logger.Error("failed decrypting ciphertext", "error", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	bgDecrypted := new(big.Int).SetUint64(decryptedValue)

	logger.Debug(functionName+" success", "input", hex.EncodeToString(input))
	return bgDecrypted, gas, nil
}

func Lte(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := "lte"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Lte(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Sub(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "sub"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Sub(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Mul(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "mul"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Mul(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	return ctHash[:], gas, nil
}

func Lt(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := "lt"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Lt(rhs)
	if err != nil {
		logger.Error(functionName+"  failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Select(utype byte, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "select"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	control, ifTrue, ifFalse, err := get3VerifiedOperands(storage, controlHash, ifTrueHash, ifFalseHash)
	if err != nil {
		logger.Error(functionName+": inputs not verified control len: ", len(controlHash), " ifTrue len: ", len(ifTrueHash), " ifFalse len: ", len(ifFalseHash), " err: ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if uintType != ifTrue.UintType || ifTrue.UintType != ifFalse.UintType {
		msg := functionName + " operands type mismatch"
		logger.Error(msg, " ifTrue ", ifTrue.UintType, " ifFalse ", ifFalse.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := control.Select(ifTrue, ifFalse)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName+" success", "control", control.Hash().Hex(), "ifTrue", ifTrue.Hash().Hex(), "ifFalse", ifTrue.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Req(utype byte, input []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: input encrypted
	//solgen: return none
	functionName := "require"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return nil, gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	if len(input) != 32 {
		msg := functionName + " input len must be 32 bytes"
		logger.Error(msg, " input ", hex.EncodeToString(input), " len ", len(input))
		return nil, 0, vm.ErrExecutionReverted
	}

	ct := getCiphertext(storage, fhe.BytesToHash(input))
	if ct == nil {
		msg := functionName + " unverified handle"
		logger.Error(msg, " input ", hex.EncodeToString(input))
		return nil, 0, vm.ErrExecutionReverted
	}

	ev := evaluateRequire(ct)

	if !ev {
		msg := functionName + " condition not met"
		logger.Error(msg)
		return nil, 0, vm.ErrExecutionReverted
	}

	return nil, gas, nil
}

func Cast(utype byte, input []byte, toType byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "cast"
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}
	if !isValidType(toType) {
		logger.Error("invalid type to cast to")
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	ct := getCiphertext(storage, fhe.BytesToHash(input))
	if ct == nil {
		logger.Error(functionName + " input not verified")
		return nil, 0, vm.ErrExecutionReverted
	}

	castToType := fhe.EncryptionType(toType)

	res, err := ct.Cast(castToType)
	if err != nil {
		msg := fmt.Sprintf("failed to cast to type %s", UtypeToString(toType))
		logger.Error(msg, " type ", castToType)
		return nil, 0, vm.ErrExecutionReverted
	}

	resHash := res.Hash()

	err = storeCipherText(storage, res)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug(functionName+" success", "ctHash", resHash.Hex())
	}

	return resHash[:], gas, nil
}

// TrivialEncrypt takes a plaintext number and encrypts it to a _compact_ ciphertext
// using the server/computation key - obviously this doesn't hide any information as the
// number was known plaintext
func TrivialEncrypt(input []byte, toType byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "trivialEncrypt"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(toType) {
		logger.Error("invalid ciphertext", "type", toType)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(toType)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	if len(input) != 32 {
		msg := functionName + " input len must be 32 bytes"
		logger.Error(msg, " input ", hex.EncodeToString(input), " len ", len(input))
		return nil, 0, vm.ErrExecutionReverted
	}

	valueToEncrypt := *new(big.Int).SetBytes(input)
	encryptToType := fhe.EncryptionType(toType)

	var ct *fhe.FheEncrypted
	var err error
	// Optimize trivial encrypts of zero since we already have trivially encrypted zeros
	// Trivial encryption of zero is common because it is done for every uninitialized ciphertext
	if State.EZero != nil && valueToEncrypt.Cmp(big.NewInt(0)) == 0 {
		// if EZero isn't initialized just initialize it
		ct = State.GetZero(encryptToType)
		if ct == nil {
			logger.Error("failed to create trivial encrypted value")
			return nil, 0, vm.ErrExecutionReverted
		}

	} else {
		// we encrypt this using the computation key not the public key. Also, compact to save space in case this gets saved directly
		// to storage
		ct, err = fhe.EncryptPlainText(valueToEncrypt, uintType)
		if err != nil {
			logger.Error("failed to create trivial encrypted value")
			return nil, 0, vm.ErrExecutionReverted
		}
	}

	ctHash := ct.Hash()
	err = storeCipherText(storage, ct)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug(functionName+" success", "ctHash", ctHash.Hex(), "valueToEncrypt", valueToEncrypt.Uint64())
	}
	return ctHash[:], gas, nil
}

func Div(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "div"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Div(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Gt(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := "gt"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Gt(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Gte(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := "gte"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Gte(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Rem(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "rem"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Rem(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func And(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := "and"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.And(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Or(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := "or"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Or(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Xor(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := "xor"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Xor(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Eq(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	//solgen: return ebool
	functionName := "eq"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Eq(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Ne(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	//solgen: return ebool
	functionName := "ne"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Ne(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Min(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "min"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Min(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Max(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "max"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Max(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Shl(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "shl"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Shl(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Shr(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "shr"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Shr(rhs)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Not(utype byte, value []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "not"

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	ct := getCiphertext(storage, fhe.BytesToHash(value))
	if ct == nil {
		msg := "not unverified ciphertext handle"
		logger.Error(msg, "input", hex.EncodeToString(value))
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := ct.Not()
	if err != nil {
		logger.Error("not failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName+" success", "input", ct.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func GetNetworkPublicKey(tp *TxParams) ([]byte, error) {
	functionName := "getNetworkPublicKey"

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	pk, err := fhe.PublicKey(0)
	if err != nil {
		logger.Error("could not get public key", "err", err)
		return nil, vm.ErrExecutionReverted
	}

	return pk, nil
}
