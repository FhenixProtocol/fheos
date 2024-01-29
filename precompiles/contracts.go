package precompiles

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"

	"github.com/sirupsen/logrus"

	tfhe "github.com/fhenixprotocol/go-tfhe"
)

var logger log.Logger

func InitLogger() {
	logger = log.Root().New("module", "fheos")
	//todo (eshel): remove log
	logger.Info("Initialized logger for fheos")
	tfhe.SetLogger(log.Root().New("module", "go-tfhe"))
}

func initTfheConfig(tfheConfig *tfhe.Config) error {
	err := tfhe.InitTfhe(tfheConfig)
	if err != nil {
		logger.Error("Failed to init tfhe config with", "error:", err)
		return err
	}

	logger.Info("Successfully initialized tfhe config", "config", tfheConfig)

	return nil
}

func InitFheos(tfheConfig *tfhe.Config) error {
	InitLogger()
	err := initTfheConfig(tfheConfig)
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

// ============================
func Add(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "add"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+": Inputs not verified", "err", err)
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

	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()

	logger.Debug(functionName+" success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Verify(utype byte, input []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "verify"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	ct, err := tfhe.NewCipherTextFromBytes(input, uintType, true /* TODO: not sure + shouldn't be hardcoded */)
	if err != nil {
		logger.Error(functionName, " failed to deserialize input ciphertext",
			" err ", err,
			" len ", len(input))
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := ct.Hash()
	err = importCiphertext(state, ct)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	logger.Debug(functionName + " success", "ctHash", ctHash.Hex())
	return ctHash[:], gas, nil
}

func SealOutput(utype byte, ctHash []byte, pk []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := "sealOutput"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	if len(ctHash) != 32 {
		msg := functionName + " ciphertext's hashes need to be 32 bytes long"
		logger.Error(msg, " ciphertxt's hash is: ", hex.EncodeToString(ctHash), " of len ", len(ctHash))
		return nil, 0, vm.ErrExecutionReverted
	}

	if len(pk) != 32 {
		msg := functionName + " public key need to be 32 bytes long"
		logger.Error(msg, " public key is: ", hex.EncodeToString(pk), " of len ", len(pk))
		return nil, 0, vm.ErrExecutionReverted
	}

	ct := getCiphertext(state, tfhe.BytesToHash(ctHash))
	if ct == nil {
		msg := functionName + " unverified ciphertext handle"
		logger.Error(msg, " ciphertext's hash: ", hex.EncodeToString(ctHash))
		return nil, 0, vm.ErrExecutionReverted
	}

	reencryptedValue, err := tfhe.SealOutput(*ct, pk)
	if err != nil {
		logger.Error(functionName, " failed to encrypt to user key", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	logger.Debug(functionName + " success", "ciphertext-hash ", hex.EncodeToString(ctHash), "public-key", hex.EncodeToString(pk))

	return reencryptedValue, gas, nil
}

func Decrypt(utype byte, input []byte, tp *TxParams) (*big.Int, uint64, error) {
	//solgen: output plaintext
	functionName := "decrypt"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.MaxUintValue, gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	if len(input) != 32 {
		msg := functionName + " input len must be 32 bytes"
		logger.Error(msg, " input ", hex.EncodeToString(input), " len ", len(input))
		return nil, 0, vm.ErrExecutionReverted
	}

	ct := getCiphertext(state, tfhe.BytesToHash(input))
	if ct == nil {
		msg := functionName + " unverified ciphertext handle"
		logger.Error(msg, " input ", hex.EncodeToString(input))
		return nil, 0, vm.ErrExecutionReverted
	}

	decryptedValue, err := tfhe.Decrypt(*ct)
	if err != nil {
		logger.Error("failed decrypting ciphertext", " error ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	bgDecrypted := new(big.Int).SetUint64(decryptedValue)

	logger.Debug(functionName + " success", "input", hex.EncodeToString(input))
	return bgDecrypted, gas, nil
}

func Lte(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := "lte"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName, " inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Lte(rhs)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Sub(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "sub"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName, " inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Sub(rhs)
	if err != nil {
		logger.Error(functionName, " failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Mul(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "mul"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName, " inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Mul(rhs)
	if err != nil {
		logger.Error(functionName, " failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	return ctHash[:], gas, nil
}

func Lt(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := "lt"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName, " inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Lt(rhs)
	if err != nil {
		logger.Error(functionName, " failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Select(utype byte, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "select"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	control, ifTrue, ifFalse, err := get3VerifiedOperands(state, controlHash, ifTrueHash, ifFalseHash)
	if err != nil {
		logger.Error(functionName, " inputs not verified control len: ", len(controlHash), " ifTrue len: ", len(ifTrueHash), " ifFalse len: ", len(ifFalseHash), " err: ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if uintType != ifTrue.UintType || ifTrue.UintType != ifFalse.UintType {
		msg := functionName + " operands type mismatch"
		logger.Error(msg, " ifTrue ", ifTrue.UintType, " ifFalse ", ifFalse.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := control.Cmux(ifTrue, ifFalse)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName + " success", "control", control.Hash().Hex(), "ifTrue", ifTrue.Hash().Hex(), "ifFalse", ifTrue.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Req(utype byte, input []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: input encrypted
	//solgen: return none
	functionName := "require"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

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

	ct := getCiphertext(state, tfhe.BytesToHash(input))
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

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}
	if !isValidType(toType) {
		logger.Error("invalid type to cast to")
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[toType], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	ct := getCiphertext(state, tfhe.BytesToHash(input))
	if ct == nil {
		logger.Error(functionName + " input not verified")
		return nil, 0, vm.ErrExecutionReverted
	}

	castToType := tfhe.UintType(toType)

	res, err := ct.Cast(castToType)
	if err != nil {
		msg := functionName + " Run() error casting ciphertext to"
		logger.Error(msg, " type ", castToType)
		return nil, 0, vm.ErrExecutionReverted
	}

	resHash := res.Hash()

	err = importCiphertext(state, res)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug(functionName+" success", "ctHash", resHash.Hex())
	}

	return resHash[:], gas, nil
}

func TrivialEncrypt(input []byte, toType byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "trivialEncrypt"

	if !isValidType(toType) {
		logger.Error("invalid ciphertext", "type", toType)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(toType)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[toType], gas, nil
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
	encryptToType := tfhe.UintType(toType)

	// Optimize trivial encrypts of zero since we already have trivially encrypted zeros
	// Trivial encryption of zero is common because it is done for every uninitialized ciphertext
	if state.EZero != nil && valueToEncrypt.Cmp(big.NewInt(0)) == 0 {
		// return trivial ciphertext
		return state.EZero[toType], gas, nil

	}

	ct, err := tfhe.NewCipherTextTrivial(valueToEncrypt, encryptToType)
	if err != nil {
		logger.Error("failed to create trivial encrypted value")
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := ct.Hash()
	err = importCiphertext(state, ct)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug(functionName + " success", "ctHash", ctHash.Hex(), "valueToEncrypt", valueToEncrypt.Uint64())
	}
	return ctHash[:], gas, nil
}

func Div(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "div"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName, " inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Div(rhs)
	if err != nil {
		logger.Error(functionName, " failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Gt(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := "gt"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName, " inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Gt(rhs)
	if err != nil {
		logger.Error(functionName, " failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Gte(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := "gte"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Gte(rhs)
	if err != nil {
		logger.Error(functionName+" failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Rem(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "rem"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Rem(rhs)
	if err != nil {
		logger.Error(functionName+" failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func And(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := "and"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.And(rhs)
	if err != nil {
		logger.Error(functionName+" failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Or(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := "or"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Or(rhs)
	if err != nil {
		logger.Error(functionName+" failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Xor(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := "xor"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Xor(rhs)
	if err != nil {
		logger.Error(functionName+" failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Eq(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	//solgen: return ebool
	functionName := "eq"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Eq(rhs)
	if err != nil {
		logger.Error(functionName+" failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Ne(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	//solgen: return ebool
	functionName := "ne"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Ne(rhs)
	if err != nil {
		logger.Error(functionName+" failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Min(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "min"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Min(rhs)
	if err != nil {
		logger.Error(functionName+" failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Max(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "max"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Max(rhs)
	if err != nil {
		logger.Error(functionName+" failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Shl(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "shl"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Shl(rhs)
	if err != nil {
		logger.Error(functionName+" failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Shr(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "shr"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	lhs, rhs, err := get2VerifiedOperands(state, lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName+" inputs not verified", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName + " operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Shr(rhs)
	if err != nil {
		logger.Error(functionName, " failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName + " success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Not(utype byte, value []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := "not"

	if !isValidType(utype) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := tfhe.UintType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return state.EZero[utype], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	ct := getCiphertext(state, tfhe.BytesToHash(value))
	if ct == nil {
		msg := "not unverified ciphertext handle"
		logger.Error(msg, "input", hex.EncodeToString(value))
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := ct.Not()
	if err != nil {
		logger.Error("not failed", " err ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = importCiphertext(state, result)
	if err != nil {
		logger.Error(functionName + " failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName + " success", "input", ct.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func GetNetworkPublicKey(tp *TxParams) ([]byte, error) {
	functionName := "getNetworkPublicKey"

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName)
	}

	pk, err := tfhe.PublicKey()
	if err != nil {
		logger.Error("could not get public key", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	return pk, nil
}
