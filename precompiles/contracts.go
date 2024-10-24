package precompiles

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
)

var logger log.Logger

func init() {
	InitLogger()
}

// InitLogger initializes a logger which will only be used just before fheos is initialized.
// currently there are no such logs.
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
	case fhe.Address:
		return "address"
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
	functionName := types.Add
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType || (lhs.UintType != uintType) {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.

	result, err := lhs.Add(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

// Verify takes inputs from the user and runs them through verification. Note that we will always get ciphertexts that
// are public-key encrypted and compressed. Anything else will fail
func Verify(utype byte, input []byte, securityZone int32, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Verify

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	ct := fhe.NewFheEncryptedFromBytes(
		input,
		uintType,
		true,
		false, // TODO: not sure + shouldn't be hardcoded
		securityZone,
	)

	err := ct.Verify()
	if err != nil {
		logger.Info(fmt.Sprintf("failed to verify ciphertext %s for type %d - was input corrupted?", ct.Hash().Hex(), uintType))
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, &ct, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "ctHash", ct.Hash().Hex())
	return ct.GetHashBytes(), gas, nil
}

func SealOutput(utype byte, ctHash []byte, pk []byte, tp *TxParams) (string, uint64, error) {
	//solgen: bool math
	functionName := types.SealOutput
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return "", 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)

	if len(ctHash) != 32 {
		msg := functionName.String() + " ciphertext's hashes need to be 32 bytes long"
		logger.Error(msg, "ciphertext-hash", hex.EncodeToString(ctHash), "hash-len", len(ctHash))
		return "", 0, vm.ErrExecutionReverted
	}

	if len(pk) != 32 {
		msg := functionName.String() + " public key need to be 32 bytes long"
		logger.Error(msg, "public-key", hex.EncodeToString(pk), "len", len(pk))
		return "", 0, vm.ErrExecutionReverted
	}

	// The rest of the function performs the following steps:
	// 1. If a result exists, we don't care about mode of execution, just return it.
	// 2. If gas estimation, return default value.
	// 3. If query, try to do sync.
	// 4. Otherwise, we're in a tx, so we are:
	//    a. Trying to asynchronously evaluate the ct.
	//    b. Required to have ParallelTxHooks.
	//    c. Return default value while async evaluation is in progress.
	hash := fhe.BytesToHash(ctHash)
	key := genSealedKey(ctHash, pk, functionName)
	record, exists := State.DecryptResults.Get(key)
	if value, ok := record.Value.(string); exists && ok {
		logger.Debug("found existing sealoutput result, returning..", "value", value)
		return value, gas, nil
	} else if tp.GasEstimation {
		return "0x" + strings.Repeat("00", 370), gas, nil
	} else {
		ct := getCiphertext(storage, hash, tp.ContractAddress)
		if ct == nil {
			msg := functionName.String() + " unverified ciphertext handle"
			logger.Error(msg, "ciphertext-hash", hex.EncodeToString(ctHash))
			return "", 0, vm.ErrExecutionReverted
		}

		if tp.EthCall {
			reencryptedValue, err := fhe.SealOutput(*ct, pk)
			if err != nil {
				logger.Error("failed to seal output", "err", err)
				return "", 0, vm.ErrExecutionReverted
			}
			return string(reencryptedValue), gas, nil
		} else if tp.ParallelTxHooks == nil {
			logger.Error("no decryption result found and no parallel tx hooks were set")
			return "", 0, vm.ErrExecutionReverted
		}

		tp.ParallelTxHooks.NotifyCt(&key)

		if !exists {
			State.DecryptResults.CreateEmptyRecord(key)
		}

		userPk := make([]byte, len(pk))
		copy(userPk, pk)
		go func() {
			logger.Debug("sealing output", "hash", hash)
			reencryptedValue, err := fhe.SealOutput(*ct, userPk)
			if err != nil {
				logger.Error("failed to seal output", "err", err)
				return
			}

			logger.Debug("sealed output", "hash", hash, "value", reencryptedValue)
			err = State.DecryptResults.SetValue(key, string(reencryptedValue))
			if err != nil {
				logger.Error("failed setting sealoutput result", "error", err)
				return
			}
			tp.ParallelTxHooks.NotifyDecryptRes(&key)
		}()

		logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "ciphertext-hash ", hex.EncodeToString(ctHash), "public-key", hex.EncodeToString(pk))
		return "0x" + strings.Repeat("00", 370), gas, nil
	}
}

func Decrypt(utype byte, input []byte, defaultValue *big.Int, tp *TxParams) (*big.Int, uint64, error) {
	//solgen: output plaintext
	functionName := types.Decrypt
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)

	if len(input) != 32 {
		msg := functionName.String() + " input len must be 32 bytes"
		logger.Error(msg, "input", hex.EncodeToString(input), "len", len(input))
		return nil, 0, vm.ErrExecutionReverted
	}

	// The rest of the function performs the following steps:
	// 1. If a result exists, we don't care about mode of execution, just return it.
	// 2. If gas estimation, return default value.
	// 3. If query, try to do sync.
	// 4. Otherwise, we're in a tx, so we are:
	//    a. Trying to asynchronously evaluate the ct.
	//    b. Required to have ParallelTxHooks.
	//    c. Return default value while async evaluation is in progress.
	ctHash := fhe.BytesToHash(input)
	key := types.PendingDecryption{
		Hash: ctHash,
		Type: functionName,
	}
	record, exists := State.DecryptResults.Get(key)
	if value, ok := record.Value.(*big.Int); exists && ok {
		logger.Debug("found existing decryption result, returning..", "value", value)
		return value, gas, nil
	} else if tp.GasEstimation {
		return defaultValue, gas, nil
	} else {
		ct := getCiphertext(storage, ctHash, tp.ContractAddress)
		if ct == nil {
			msg := functionName.String() + " unverified ciphertext handle"
			logger.Error(msg, " input ", hex.EncodeToString(input))
			return nil, 0, vm.ErrExecutionReverted
		}

		if tp.EthCall {
			decryptedValue, err := fhe.Decrypt(*ct)
			if err != nil {
				logger.Error("failed decrypting ciphertext", "error", err)
				return nil, 0, vm.ErrExecutionReverted
			}
			return decryptedValue, gas, nil
		} else if tp.ParallelTxHooks == nil {
			logger.Error("no decryption result found and no parallel tx hooks were set")
			return nil, 0, vm.ErrExecutionReverted
		}

		tp.ParallelTxHooks.NotifyCt(&key)

		if !exists {
			State.DecryptResults.CreateEmptyRecord(key)
		}

		go func() {
			logger.Debug("decrypting ciphertext", "hash", ctHash)
			decryptedValue, err := fhe.Decrypt(*ct)
			if err != nil {
				logger.Error("failed decrypting ciphertext", "error", err)
				return
			}

			logger.Debug("decrypted value", "hash", ctHash, "value", decryptedValue)
			err = State.DecryptResults.SetValue(key, decryptedValue)
			if err != nil {
				logger.Error("failed setting decryption result", "error", err)
				return
			}
			tp.ParallelTxHooks.NotifyDecryptRes(&key)
		}()

		logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "input", hex.EncodeToString(input))
		return defaultValue, gas, nil
	}
}

func Lte(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := types.Lte

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Lte(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Sub(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Sub

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Sub(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Mul(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Mul

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Mul(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	return ctHash[:], gas, nil
}

func Lt(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := types.Lt

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Lt(rhs)
	if err != nil {
		logger.Error(functionName.String()+"  failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Select(utype byte, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Select

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	control, ifTrue, ifFalse, err := get3VerifiedOperands(storage, controlHash, ifTrueHash, ifFalseHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+": inputs not verified control len: ", len(controlHash), " ifTrue len: ", len(ifTrueHash), " ifFalse len: ", len(ifFalseHash), " err: ", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if uintType != ifTrue.UintType || ifTrue.UintType != ifFalse.UintType {
		msg := functionName.String() + " operands type mismatch"
		logger.Error(msg, " ifTrue ", ifTrue.UintType, " ifFalse ", ifFalse.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := control.Select(ifTrue, ifFalse)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "control", control.Hash().Hex(), "ifTrue", ifTrue.Hash().Hex(), "ifFalse", ifTrue.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Req(utype byte, input []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: input encrypted
	//solgen: return none
	functionName := types.Require
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)

	if len(input) != 32 {
		msg := functionName.String() + " input len must be 32 bytes"
		logger.Error(msg, "input", hex.EncodeToString(input), "len", len(input))
		return nil, 0, vm.ErrExecutionReverted
	}

	// The rest of the function performs the following steps:
	// 1. If a result exists, we don't care about mode of execution, just return it.
	// 2. If gas estimation, return default value.
	// 3. If query, try to do sync.
	// 4. Otherwise, we're in a tx, so we are:
	//    a. Trying to asynchronously evaluate the ct.
	//    b. Required to have ParallelTxHooks.
	//    c. Return default value while async evaluation is in progress.
	ctHash := fhe.BytesToHash(input)
	key := types.PendingDecryption{
		Hash: ctHash,
		Type: functionName,
	}
	record, exists := State.DecryptResults.Get(key)
	if value, ok := record.Value.(bool); exists && ok {
		logger.Debug("found existing decryption result, returning..", "value", value)
		if !value {
			return nil, gas, vm.ErrExecutionReverted
		}
		return nil, gas, nil
	} else if tp.GasEstimation {
		return nil, gas, nil
	} else {
		ct := getCiphertext(storage, ctHash, tp.ContractAddress)
		if ct == nil {
			msg := functionName.String() + " unverified ciphertext handle"
			logger.Error(msg, "input", hex.EncodeToString(input))
			return nil, 0, vm.ErrExecutionReverted
		}

		if tp.EthCall {
			result, err := evaluateRequire(ct)
			if err != nil {
				msg := functionName.String() + " error on evaluation"
				logger.Error(msg, " err ", err)
				return nil, gas, vm.ErrExecutionReverted
			}

			if !result {
				return nil, gas, vm.ErrExecutionReverted
			}
			return nil, gas, nil
		} else if tp.ParallelTxHooks == nil {
			logger.Error("no decryption result found and no parallel tx hooks were set")
			return nil, 0, vm.ErrExecutionReverted
		}

		tp.ParallelTxHooks.NotifyCt(&key)

		if !exists {
			State.DecryptResults.CreateEmptyRecord(key)
		}

		go func() {
			logger.Debug("evaluating require condition", "hash", ctHash)
			result, err := evaluateRequire(ct)
			if err != nil {
				msg := functionName.String() + " error on evaluation"
				logger.Error(msg, " err ", err)
				return
			}

			logger.Debug("require condition result", "hash", ctHash, "value", result)
			err = State.DecryptResults.SetValue(key, result)
			if err != nil {
				logger.Error("failed setting require result", "error", err)
				return
			}
			tp.ParallelTxHooks.NotifyDecryptRes(&key)
		}()

		logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "input", hex.EncodeToString(input))

		return nil, gas, nil
	}
}

func Cast(utype byte, input []byte, toType byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Cast
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}
	if !types.IsValidType(fhe.EncryptionType(toType)) {
		logger.Error("invalid type to cast to")
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	ct := getCiphertext(storage, fhe.BytesToHash(input), tp.ContractAddress)
	if ct == nil {
		logger.Error(functionName.String() + " input not verified")
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

	err = storeCipherText(storage, res, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "ctHash", resHash.Hex())
	}

	return resHash[:], gas, nil
}

// TrivialEncrypt takes a plaintext number and encrypts it to a _compact_ ciphertext
// using the server/computation key - obviously this doesn't hide any information as the
// number was known plaintext
func TrivialEncrypt(input []byte, toType byte, securityZone int32, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.TrivialEncrypt

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	uintType := fhe.EncryptionType(toType)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", toType)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	if len(input) != 32 {
		msg := functionName.String() + " input len must be 32 bytes"
		logger.Error(msg, " input ", hex.EncodeToString(input), " len ", len(input))
		return nil, 0, vm.ErrExecutionReverted
	}

	valueToEncrypt := *new(big.Int).SetBytes(input)

	var ct *fhe.FheEncrypted
	var err error
	// Check if value is not overflowing the type
	maxOfType := fhe.MaxOfType(uintType)
	if maxOfType == nil {
		logger.Error("failed to create trivially encrypted value, type is not supported.")
		return nil, 0, vm.ErrExecutionReverted
	}

	// If value is bigger than the maximal value that is supported by the type
	if valueToEncrypt.Cmp(maxOfType) > 0 {
		logger.Error("failed to create trivially encrypted value, value is too large for type.")
		return nil, 0, vm.ErrExecutionReverted
	}

	// we encrypt this using the computation key not the public key. Also, compact to save space in case this gets saved directly
	// to storage
	ct, err = fhe.EncryptPlainText(valueToEncrypt, uintType, securityZone)
	if err != nil {
		logger.Error("failed to create trivial encrypted value")
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := ct.Hash()
	err = storeCipherText(storage, ct, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "ctHash", ctHash.Hex(), "valueToEncrypt", valueToEncrypt.Uint64(), "securityZone", ct.SecurityZone)
	}
	return ctHash[:], gas, nil
}

func Div(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Div

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Div(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Gt(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := types.Gt

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+": inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Gt(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Gte(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := types.Gte

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Gte(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Rem(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Rem

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Rem(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func And(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := types.And

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.And(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Or(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := types.Or

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Or(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Xor(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := types.Xor

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Xor(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Eq(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	//solgen: return ebool
	functionName := types.Eq

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Eq(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Ne(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	//solgen: bool math
	//solgen: return ebool
	functionName := types.Ne

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Ne(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Min(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Min

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Min(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Max(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Max

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Max(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Shl(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Shl

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Shl(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Shr(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Shr

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Shr(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Rol(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Rol

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Rol(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Ror(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Ror

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhs, rhs, err := get2VerifiedOperands(storage, lhsHash, rhsHash, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" inputs not verified", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := functionName.String() + " operand type mismatch"
		logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := lhs.Ror(rhs)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	ctHash := result.Hash()

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], gas, nil
}

func Not(utype byte, value []byte, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Not

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	ct := getCiphertext(storage, fhe.BytesToHash(value), tp.ContractAddress)
	if ct == nil {
		msg := "unverified ciphertext handle"
		logger.Error(msg, "input", hex.EncodeToString(value))
		return nil, 0, vm.ErrExecutionReverted
	}

	result, err := ct.Not()
	if err != nil {
		logger.Error("not failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "input", ct.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Random(utype byte, seed uint64, securityZone int32, tp *TxParams) ([]byte, uint64, error) {
	functionName := types.Random

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid random output type", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	var finalSeed uint64
	if seed != 0 {
		finalSeed = seed
	} else {
		var randomCounter uint64
		var hash common.Hash
		if tp.Commit {
			// We're incrementing before the request for the random number, so that queries
			// that came before this Tx would have received a different seed.
			randomCounter = State.IncRandomCounter()
			hash = tp.TxContext.Hash
		} else {
			randomCounter = State.GetRandomCounter()
			hash = tp.GetBlockHash(tp.BlockNumber.Uint64() - 1) // If no tx hash - use block hash
		}

		finalSeed = GenerateSeedFromEntropy(tp.ContractAddress, hash, randomCounter)
	}

	result, err := fhe.FheRandom(securityZone, uintType, finalSeed)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, result, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	resultHash := result.Hash()
	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func GetNetworkPublicKey(securityZone int32, tp *TxParams) ([]byte, error) {
	functionName := types.GetNetworkKey

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	pk, err := fhe.PublicKey(securityZone)
	if err != nil {
		logger.Error("could not get public key", "err", err, "securityZone", securityZone)
		return nil, vm.ErrExecutionReverted
	}

	return pk, nil
}

func Square(utype byte, value []byte, tp *TxParams) ([]byte, uint64, error) {
	return Mul(utype, value, value, tp)
	// Please don't delete the below comment, this is intentionally left here for code generation.
	//ct := getCiphertext(storage, fhe.BytesToHash(value), tp.ContractAddress)
}
