package precompiles

import (
	"encoding/hex"
	"errors"
	"github.com/sirupsen/logrus"
	"math/big"
	"runtime"

	"github.com/ethereum/go-ethereum/core/vm"
	tfhe "github.com/fhenixprotocol/go-tfhe"
)

var interpreter *vm.EVMInterpreter
var logger *logrus.Logger

// FHENIX: TODO - persist it somehow
var ctHashMap map[tfhe.Hash]*tfhe.Ciphertext

func InitLogger() {
	logger = newLogger()
	tfhe.InitLogger(getDefaultLogLevel())
}

func SetEvmInterpreter(i *vm.EVMInterpreter, tfheConfig *tfhe.Config) error {
	if ctHashMap == nil {
		ctHashMap = make(map[tfhe.Hash]*tfhe.Ciphertext)
	}

	err := tfhe.InitTfhe(tfheConfig)
	if err != nil {
		logger.Error("Failed to init tfhe config with error: ", err)
		return err
	}

	interpreter = i
	return nil
}

func shouldPrintPrecompileInfo() bool {
	return interpreter.GetEVM().Commit && !interpreter.GetEVM().GasEstimation
}

func validateInterpreter() error {
	if interpreter == nil {
		msg := "no evm interpreter"
		return errors.New(msg)
	}

	return nil
}

func getFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	return funcName
}

// ============================
func Add(input []byte) ([]byte, error) {
	err := validateInterpreter()
	if err != nil {
		logger.Error("failed validating evm interpreter for function ", getFunctionName())
		return nil, err
	}

	if shouldPrintPrecompileInfo() {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheAdd inputs not verified", "err", err, "input", hex.EncodeToString(input))
		return nil, err
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheAdd operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, errors.New(msg)
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if interpreter.GetEVM().GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Add(rhs)
	if err != nil {
		logger.Error("fheAdd failed", "err", err)
		return nil, err
	}

	importCiphertext(result)

	resultHash := result.Hash()
	logger.Debug("fheAdd success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], nil
}

func Verify(input []byte) ([]byte, error) {
	err := validateInterpreter()
	if err != nil {
		logger.Error("failed validating evm interpreter for function ", getFunctionName())
		return nil, err
	}
	if shouldPrintPrecompileInfo() {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	if len(input) <= 1 {
		msg := "verifyCiphertext RequiredGas() input needs to contain a ciphertext and one byte for its type"
		logger.Error(msg, "len", len(input))
		return nil, errors.New(msg)
	}

	ctBytes := input[:len(input)-1]
	ctType := tfhe.UintType(input[len(input)-1])

	ct, err := tfhe.NewCipherTextFromBytes(ctBytes, ctType, true /* TODO: not sure + shouldn't be hardcoded */)
	if err != nil {
		logger.Error("verifyCiphertext failed to deserialize input ciphertext",
			"err", err,
			"len", len(ctBytes),
			"ctBytes64", hex.EncodeToString(ctBytes[:minInt(len(ctBytes), 64)]))
		return nil, err
	}
	ctHash := ct.Hash()
	importCiphertext(ct)

	if interpreter.GetEVM().Commit {
		logger.Debug("verifyCiphertext success",
			"ctHash", ctHash.Hex(),
			"ctBytes64", hex.EncodeToString(ctBytes[:minInt(len(ctBytes), 64)]))
	}
	return ctHash[:], nil
}

func Reencrypt(input []byte) ([]byte, error) {
	err := validateInterpreter()
	if err != nil {
		logger.Error("failed validating evm interpreter for function ", getFunctionName())
		return nil, err
	}
	if shouldPrintPrecompileInfo() {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	if !interpreter.GetEVM().EthCall {
		msg := "reencrypt only supported on EthCall"
		logger.Error(msg)
		return nil, errors.New(msg)
	}

	if len(input) != 64 {
		msg := "reencrypt input len must be 64 bytes"
		logger.Error(msg, "input", hex.EncodeToString(input), "len", len(input))
		return nil, errors.New(msg)
	}

	ct := getCiphertext(tfhe.BytesToHash(input[0:32]))
	if ct != nil {
		decryptedValue, err := ct.Decrypt()
		if err != nil {
			logger.Error("failed decrypting ciphertext", "error", err)
			return nil, err
		}

		pubKey := input[32:64]
		reencryptedValue, err := encryptToUserKey(decryptedValue, pubKey)
		if err != nil {
			logger.Error("reencrypt failed to encrypt to user key", "err", err)
			return nil, err
		}
		logger.Debug("reencrypt success", "input", hex.EncodeToString(input))
		// FHENIX: Previously it was "return toEVMBytes(reencryptedValue), nil" but the decrypt function in Fhevm din't support it so we removed the the toEVMBytes
		return reencryptedValue, nil
	}
	msg := "reencrypt unverified ciphertext handle"
	logger.Error(msg, "input", hex.EncodeToString(input))
	return nil, errors.New(msg)
}

func Lte(input []byte) ([]byte, error) {
	err := validateInterpreter()
	if err != nil {
		logger.Error("failed validating evm interpreter for function ", getFunctionName())
		return nil, err
	}
	if shouldPrintPrecompileInfo() {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheLte inputs not verified", "err", err)
		return nil, err
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheLte operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, errors.New(msg)
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if interpreter.GetEVM().GasEstimation {
		return importRandomCiphertext(lhs.UintType)

	}

	result, err := lhs.Lte(rhs)
	if err != nil {
		logger.Error("fheLte failed", "err", err)
		return nil, err
	}
	importCiphertext(result)

	resultHash := result.Hash()
	logger.Debug("fheLte success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], nil
}

func Sub(input []byte) ([]byte, error) {
	err := validateInterpreter()
	if err != nil {
		logger.Error("failed validating evm interpreter for function ", getFunctionName())
		return nil, err
	}
	if shouldPrintPrecompileInfo() {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheSub inputs not verified", "err", err)
		return nil, err
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheSub operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, errors.New(msg)
	}

	// // If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if interpreter.GetEVM().GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Sub(rhs)
	if err != nil {
		logger.Error("fheSub failed", "err", err)
		return nil, err
	}
	importCiphertext(result)

	resultHash := result.Hash()
	logger.Debug("fheSub success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], nil
}

func Mul(input []byte) ([]byte, error) {
	err := validateInterpreter()
	if err != nil {
		logger.Error("failed validating evm interpreter for function ", getFunctionName())
		return nil, err
	}
	if shouldPrintPrecompileInfo() {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheMul inputs not verified", "err", err)
		return nil, err
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheMul operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, errors.New(msg)
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if interpreter.GetEVM().GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Mul(rhs)
	if err != nil {
		logger.Error("fheMul failed", "err", err)
		return nil, err
	}
	importCiphertext(result)

	ctHash := result.Hash()

	return ctHash[:], nil
}

func Lt(input []byte) ([]byte, error) {
	err := validateInterpreter()
	if err != nil {
		logger.Error("failed validating evm interpreter for function ", getFunctionName())
		return nil, err
	}
	if shouldPrintPrecompileInfo() {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheLt inputs not verified", "err", err)
		return nil, err
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheLt operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, errors.New(msg)
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if interpreter.GetEVM().GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Lt(rhs)
	if err != nil {
		logger.Error("fheLt failed", "err", err)
		return nil, err
	}
	importCiphertext(result)

	resultHash := result.Hash()
	logger.Debug("fheLt success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], nil
}

func Cmux(input []byte) ([]byte, error) {
	err := validateInterpreter()
	if err != nil {
		logger.Error("failed validating evm interpreter for function ", getFunctionName())
		return nil, err
	}
	if shouldPrintPrecompileInfo() {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	control, ifTrue, ifFalse, err := get3VerifiedOperands(input)
	if err != nil {
		logger.Error("selector inputs not verified input len:", len(input), " err: ", err)
		return nil, err
	}

	if (ifTrue.UintType != ifFalse.UintType) || (control.UintType != ifTrue.UintType) {
		msg := "selector operands type mismatch"
		logger.Error(msg, "control", control.UintType, "ifTrue", ifTrue.UintType, "ifFalse", ifFalse.UintType)
		return nil, errors.New(msg)
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if interpreter.GetEVM().GasEstimation {
		return importRandomCiphertext(ifTrue.UintType)
	}

	result, err := control.Cmux(ifTrue, ifFalse)
	if err != nil {
		logger.Error("selector failed", "err", err)
		return nil, err
	}
	importCiphertext(result)

	resultHash := result.Hash()
	logger.Debug("selector success", "control", control.Hash().Hex(), "ifTrue", ifTrue.Hash().Hex(), "ifFalse", ifTrue.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], nil
}

func Req(input []byte) ([]byte, error) {
	err := validateInterpreter()
	if err != nil {
		logger.Error("failed validating evm interpreter for function ", getFunctionName())
		return nil, err
	}

	if shouldPrintPrecompileInfo() {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	if interpreter.GetEVM().EthCall {
		msg := "require not supported on EthCall"
		logger.Error(msg)
		return nil, errors.New(msg)
	}

	if len(input) != 32 {
		msg := "require input len must be 32 bytes"
		logger.Error(msg, "input", hex.EncodeToString(input), "len", len(input))
		return nil, errors.New(msg)
	}

	ct := getCiphertext(tfhe.BytesToHash(input))
	if ct == nil {
		msg := "optimisticRequire unverified handle"
		logger.Error(msg, "input", hex.EncodeToString(input))
		return nil, errors.New(msg)
	}
	// If we are not committing to state, assume the require is true, avoiding any side effects
	// (i.e. mutatiting the oracle DB).
	if !interpreter.GetEVM().Commit {
		return nil, nil
	}
	if ct.UintType != tfhe.Uint32 {
		msg := "require ciphertext type is not euint32"
		logger.Error(msg, "type", ct.UintType)
		return nil, errors.New(msg)
	}

	ev := evaluateRequire(ct, interpreter)

	if !ev {
		logger.Error("require failed to evaluate, reverting")
		return nil, vm.ErrExecutionReverted
	}

	return nil, nil
}

func Cast(input []byte) ([]byte, error) {
	err := validateInterpreter()
	if err != nil {
		logger.Error("failed validating evm interpreter for function ", getFunctionName())
		return nil, err
	}

	if shouldPrintPrecompileInfo() {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	if !isValidType(input[32]) {
		logger.Error("invalid type to cast to")
		return nil, errors.New("invalid type provided")
	}
	castToType := tfhe.UintType(input[32])

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if interpreter.GetEVM().GasEstimation {
		return importRandomCiphertext(castToType)
	}

	ct := getCiphertext(tfhe.BytesToHash(input[0:32]))
	if ct == nil {
		logger.Error("cast input not verified")
		return nil, errors.New("unverified ciphertext handle")
	}

	res, err := ct.Cast(castToType)
	if err != nil {
		msg := "cast Run() error casting ciphertext to"
		logger.Error(msg, "type", castToType)
		return nil, errors.New(msg)
	}

	resHash := res.Hash()

	importCiphertext(res)
	if shouldPrintPrecompileInfo() {
		logger.Debug("cast success",
			"ctHash", resHash.Hex(),
		)
	}

	return resHash[:], nil
}

func TrivialEncrypt(input []byte) ([]byte, error) {
	err := validateInterpreter()
	if err != nil {
		logger.Error("failed validating evm interpreter for function ", getFunctionName())
		return nil, err
	}

	if shouldPrintPrecompileInfo() {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	if len(input) != 33 {
		msg := "trivialEncrypt input len must be 33 bytes"
		logger.Error(msg, "input", hex.EncodeToString(input), "len", len(input))
		return nil, errors.New(msg)
	}

	valueToEncrypt := *new(big.Int).SetBytes(input[0:32])
	encryptToType := tfhe.UintType(input[32])

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if interpreter.GetEVM().GasEstimation {
		return importRandomCiphertext(encryptToType)
	}

	ct, err := tfhe.NewCipherTextTrivial(valueToEncrypt, encryptToType)
	if err != nil {
		logger.Error("failed to create trivial encrypted value")
		return nil, err
	}

	ctHash := ct.Hash()
	importCiphertext(ct)
	if shouldPrintPrecompileInfo() {
		logger.Debug("trivialEncrypt success",
			"ctHash", ctHash.Hex(),
			"valueToEncrypt", valueToEncrypt.Uint64())
	}
	return ctHash[:], nil
}
