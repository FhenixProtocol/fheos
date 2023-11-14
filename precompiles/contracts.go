package precompiles

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/vm"
	tfhe "github.com/fhenixprotocol/go-tfhe"
)

var interpreter *vm.EVMInterpreter

// FHENIX: TODO - persist it somehow
var ctHashMap map[tfhe.Hash]*tfhe.Ciphertext

func SetEvmInterpreter(i *vm.EVMInterpreter, tfheConfig *tfhe.Config) error {
	if ctHashMap == nil {
		ctHashMap = make(map[tfhe.Hash]*tfhe.Ciphertext)
	}

	err := tfhe.InitTfhe(tfheConfig)
	if err != nil {
		return err
	}

	interpreter = i
	return nil
}

func validatedInterpreter() error {
	if interpreter == nil {
		msg := "no evm interpreter"
		// logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return errors.New(msg)
	}

	return nil
}

func getLogger() vm.Logger {
	return interpreter.GetEVM().Logger
}

// ============================
func Add(input []byte, inputLen uint32) ([]byte, error) {
	funcName := "Add"
	fmt.Printf("Starting new precompiled contract function %s\n", funcName)
	err := validatedInterpreter()
	if err != nil {
		return nil, err
	}

	logger := getLogger()

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
		return importRandomCiphertext(lhs.UintType), nil
	}

	result, err := lhs.Add(rhs)
	if err != nil {
		logger.Error("fheAdd failed", "err", err)
		return nil, err
	}

	importCiphertext(result)

	// TODO: for testing
	err = os.WriteFile("/tmp/add_result", result.Serialization, 0644)
	if err != nil {
		logger.Error("fheAdd failed to write /tmp/add_result", "err", err)
		return nil, err
	}

	resultHash := result.Hash()
	logger.Info("fheAdd success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], nil
}

func Verify(input []byte, inputLen uint32) ([]byte, error) {
	funcName := "Verify"
	fmt.Printf("Starting new precompiled contract function %s\n", funcName)
	err := validatedInterpreter()
	if err != nil {
		return nil, err
	}

	logger := getLogger()

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
		logger.Info("verifyCiphertext success",
			"ctHash", ctHash.Hex(),
			"ctBytes64", hex.EncodeToString(ctBytes[:minInt(len(ctBytes), 64)]))
	}
	return ctHash[:], nil
}

func Reencrypt(input []byte, inputLen uint32) ([]byte, error) {
	funcName := "Reencrypt"
	fmt.Printf("Starting new precompiled contract function %s\n", funcName)
	err := validatedInterpreter()
	if err != nil {
		return nil, err
	}

	logger := getLogger()

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
			panic(fmt.Sprintf("Failed to decrypt ciphertext: %+v", err))
		}

		pubKey := input[32:64]
		reencryptedValue, err := encryptToUserKey(decryptedValue, pubKey)
		if err != nil {
			logger.Error("reencrypt failed to encrypt to user key", "err", err)
			return nil, err
		}
		logger.Info("reencrypt success", "input", hex.EncodeToString(input))
		// FHENIX: Previously it was "return toEVMBytes(reencryptedValue), nil" but the decrypt function in Fhevm din't support it so we removed the the toEVMBytes
		return reencryptedValue, nil
	}
	msg := "reencrypt unverified ciphertext handle"
	logger.Error(msg, "input", hex.EncodeToString(input))
	return nil, errors.New(msg)
}

func Lte(input []byte, inputLen uint32) ([]byte, error) {
	funcName := "Lte"
	fmt.Printf("Starting new precompiled contract function %s\n", funcName)
	err := validatedInterpreter()
	if err != nil {
		return nil, err
	}

	logger := getLogger()

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
		return importRandomCiphertext(lhs.UintType), nil
	}

	result, err := lhs.Lte(rhs)
	if err != nil {
		logger.Error("fheLte failed", "err", err)
		return nil, err
	}
	importCiphertext(result)

	// TODO: for testing
	err = os.WriteFile("/tmp/lte_result", result.Serialization, 0644)
	if err != nil {
		logger.Error("fheAdd failed to write /tmp/lte_result", "err", err)
		return nil, err
	}

	resultHash := result.Hash()
	logger.Info("fheLte success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], nil
}

func Sub(input []byte, inputLen uint32) ([]byte, error) {
	funcName := "Sub"
	fmt.Printf("Starting new precompiled contract function %s\n", funcName)
	err := validatedInterpreter()
	if err != nil {
		return nil, err
	}

	logger := getLogger()

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
		return importRandomCiphertext(lhs.UintType), nil
	}

	result, err := lhs.Sub(rhs)
	if err != nil {
		logger.Error("fheSub failed", "err", err)
		return nil, err
	}
	importCiphertext(result)

	// TODO: for testing
	err = os.WriteFile("/tmp/sub_result", result.Serialization, 0644)
	if err != nil {
		logger.Error("fheSub failed to write /tmp/sub_result", "err", err)
		return nil, err
	}

	resultHash := result.Hash()
	logger.Info("fheSub success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], nil
}

func Mul(input []byte, inputLen uint32) ([]byte, error) {
	funcName := "Mul"
	fmt.Printf("Starting new precompiled contract function %s\n", funcName)
	err := validatedInterpreter()
	if err != nil {
		return nil, err
	}

	logger := getLogger()

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
		return importRandomCiphertext(lhs.UintType), nil
	}

	result, err := lhs.Mul(rhs)
	if err != nil {
		logger.Error("fheMul failed", "err", err)
		return nil, err
	}
	importCiphertext(result)

	// TODO: for testing
	err = os.WriteFile("/tmp/mul_result", result.Serialization, 0644)
	if err != nil {
		logger.Error("fheMul failed to write /tmp/mul_result", "err", err)
		return nil, err
	}

	ctHash := result.Hash()

	return ctHash[:], nil
}

func Lt(input []byte, inputLen uint32) ([]byte, error) {
	funcName := "Lt"
	fmt.Printf("Starting new precompiled contract function %s\n", funcName)
	err := validatedInterpreter()
	if err != nil {
		return nil, err
	}

	logger := getLogger()

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
		return importRandomCiphertext(lhs.UintType), nil
	}

	result, err := lhs.Lt(rhs)
	if err != nil {
		logger.Error("fheLt failed", "err", err)
		return nil, err
	}
	importCiphertext(result)

	// TODO: for testing
	err = os.WriteFile("/tmp/lt_result", result.Serialization, 0644)
	if err != nil {
		logger.Error("fheLt failed to write /tmp/lt_result", "err", err)
		return nil, err
	}

	resultHash := result.Hash()
	logger.Info("fheLt success", "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], nil
}

func Req(input []byte, inputLen uint32) ([]byte, error) {
	funcName := "Require"
	fmt.Printf("Starting new precompiled contract function %s\n", funcName)
	err := validatedInterpreter()
	if err != nil {
		return nil, err
	}

	logger := getLogger()

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

	if !evaluateRequire(ct, interpreter) {
		logger.Error("require failed to evaluate, reverting")
		return nil, vm.ErrExecutionReverted
	}

	return nil, nil
}

// Implementation of the following is pending and will be completed once TFHE-rs add type casts to their high-level C API.
func Cast(input []byte, inputLen uint32) ([]byte, error) {
	return nil, nil
}

func TrivialEncrypt(input []byte) ([]byte, error) {
	funcName := "TrivialEncrypt"
	fmt.Printf("Starting new precompiled contract function %s\n", funcName)
	err := validatedInterpreter()
	if err != nil {
		return nil, err
	}

	logger := getLogger()

	if len(input) != 33 {
		msg := "trivialEncrypt input len must be 33 bytes"
		logger.Error(msg, "input", hex.EncodeToString(input), "len", len(input))
		return nil, errors.New(msg)
	}

	valueToEncrypt := *new(big.Int).SetBytes(input[0:32])
	encryptToType := tfhe.UintType(input[32])

	ct, err := tfhe.NewCipherTextTrivial(valueToEncrypt, encryptToType)
	if err != nil {
		logger.Error("Failed to create trivial encrypted value")
		return nil, err
	}

	ctHash := ct.Hash()
	importCiphertext(ct)

	if interpreter.GetEVM().Commit {
		logger.Info("trivialEncrypt success",
			"ctHash", ctHash.Hex(),
			"valueToEncrypt", valueToEncrypt.Uint64())
	}
	return ctHash[:], nil
}
