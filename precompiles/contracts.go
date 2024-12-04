package precompiles

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
)

var logger log.Logger

func init() {
	InitLogger()
}

// ============================================================

func Log(s string, tp *TxParams) (uint64, error) {
	if tp.GasEstimation {
		return 1, nil
	}

	logger.Debug(fmt.Sprintf("Contract Log: %s", s))
	return 1, nil
}

func Add(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Add
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Add, utype, lhsHash, rhsHash, tp, callback)
}

// Verify takes inputs from the user and runs them through verification. Note that we will always get ciphertexts that
// are public-key encrypted and compressed. Anything else will fail
func Verify(utype byte, input []byte, securityZone int32, tp *TxParams, _ *CallbackFunc) ([]byte, uint64, error) {
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
		false,
	)

	err := ct.Verify()
	if err != nil {
		logger.Info(fmt.Sprintf("failed to verify ciphertext %s for type %d - was input corrupted?", ct.GetHash().Hex(), uintType))
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, &ct, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "ctHash", ct.GetHash().Hex())
	return ct.GetHashBytes(), gas, nil
}

func SealOutput(utype byte, ctHash []byte, pk []byte, tp *TxParams, onResultCallback *SealOutputCallbackFunc) (string, uint64, error) {
	//solgen: bool math
	functionName := types.SealOutput
	gas, err := PreProcessOperation1(functionName, utype, ctHash, tp)
	if err != nil {
		return "", gas, vm.ErrExecutionReverted
	}

	if len(pk) != 32 {
		msg := functionName.String() + " public key need to be 32 bytes long"
		logger.Error(msg, "public-key", hex.EncodeToString(pk), "len", len(pk))
		return "", gas, vm.ErrExecutionReverted
	}

	//if onResultCallback == nil {
	//	msg := functionName.String() + " must set callback"
	//	logger.Error(msg, " ctHash ", ctHash)
	//	return "", gas, vm.ErrExecutionReverted
	//}

	if !tp.GasEstimation {
		storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
		if onResultCallback == nil {
			sealed, err := SealOutputHelper(storage, ctHash, pk, tp)
			return sealed, gas, err
		}

		go func(ctHash []byte) {
			sealed, err := SealOutputHelper(storage, ctHash, pk, tp)
			if err != nil {
				return
			}

			url := (*onResultCallback).CallbackUrl
			(*onResultCallback).Callback(url, ctHash, string(sealed))
		}(ctHash)
		logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "ctHash", hex.EncodeToString(ctHash))
	}

	return "0x" + strings.Repeat("00", 370), gas, nil
}

func Decrypt(utype byte, input []byte, defaultValue *big.Int, tp *TxParams, onResultCallback *DecryptCallbackFunc) (*big.Int, uint64, error) {
	//solgen: output plaintext
	functionName := types.Decrypt

	gas, err := PreProcessOperation1(functionName, utype, input, tp)
	if err != nil {
		return nil, gas, vm.ErrExecutionReverted
	}

	//if onResultCallback == nil {
	//	msg := functionName.String() + " must set callback"
	//	logger.Error(msg, " input ", input)
	//	return nil, gas, vm.ErrExecutionReverted
	//}

	if !tp.GasEstimation {
		storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
		if onResultCallback == nil {
			plaintext, err := DecryptHelper(storage, input, tp, defaultValue)
			return plaintext, gas, err
		}
		go func(ctHash []byte) {
			plaintext, err := DecryptHelper(storage, ctHash, tp, defaultValue)
			if err != nil {
				logger.Error("failed decrypting ciphertext", "error", err)
				return
			}

			url := (*onResultCallback).CallbackUrl
			(*onResultCallback).Callback(url, input, plaintext)
		}(input)
		logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "input", hex.EncodeToString(input))
	}

	return defaultValue, gas, nil
}

func Lte(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := types.Lte
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Lte, utype, lhsHash, rhsHash, tp, callback)
}

func Sub(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Sub
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Sub, utype, lhsHash, rhsHash, tp, callback)
}

func Mul(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Mul
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Mul, utype, lhsHash, rhsHash, tp, callback)
}

func Lt(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := types.Lt
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Lt, utype, lhsHash, rhsHash, tp, callback)
}

func Select(utype byte, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Select

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	placeholderCt, err := createPlaceholder(utype, functionName, controlHash, ifTrueHash, ifFalseHash)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

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

	err = storeCipherText(storage, placeholderCt, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("fn", functionName.String(), "Storing async ciphertext", "placeholderKey", hex.EncodeToString(placeholderCt.Hash))
	}

	err = storage.SetAsyncCtStart(types.Hash(placeholderCt.Hash))
	if err != nil {
		logger.Error(functionName.String()+" failed to set async value start", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("Starting new async precompiled contract function: " + functionName.String())
	}

	controlHashCopy := CopySlice(controlHash)
	ifTrueHashCopy := CopySlice(ifTrueHash)
	ifFalseHashCopy := CopySlice(ifFalseHash)
	placeholderKeyCopy := CopySlice(placeholderCt.Hash)

	go func(controlHash, ifTrueHash, ifFalseHash, resultHash []byte) {
		logger.Info("About to block until inputs are available")
		defer logger.Info("Finishing goroutine execution")

		ct, err := blockUntilThreeInputsAvailable(storage, tp, controlHash, ifTrueHash, ifFalseHash)
		if err != nil {
			logger.Error(functionName.String()+": inputs not verified control len: ", len(controlHash), " ifTrue len: ", len(ifTrueHash), " ifFalse len: ", len(ifFalseHash), " err: ", err)
			return
		}

		logger.Info("Inputs are available")

		control, ifTrue, ifFalse := ct[0], ct[1], ct[2]

		if uintType != ifTrue.UintType || ifTrue.UintType != ifFalse.UintType {
			msg := functionName.String() + " operands type mismatch"
			logger.Error(msg, " ifTrue ", ifTrue.UintType, " ifFalse ", ifFalse.UintType)
			return
		}

		logger.Info("About to select")

		result, err := control.Select(ifTrue, ifFalse)
		if err != nil {
			logger.Error(functionName.String()+" failed", "err", err)
			return
		}

		logger.Info("About to decode result hash")

		realResultHash, err := hex.DecodeString(result.GetHash().Hex())
		if err != nil {
			logger.Error(functionName.String()+" failed", "err", err)
			return
		}

		logger.Info("Result hash decoded")

		result.Hash = resultHash

		logger.Info("About to store result")

		err = storeCipherText(storage, result, tp.ContractAddress)
		if err != nil {
			logger.Error(functionName.String()+" failed", "err", err)
			return
		}

		logger.Info("About to set async ct done")

		_ = storage.SetAsyncCtDone(types.Hash(resultHash))
		if callback != nil {
			url := (*callback).CallbackUrl
			(*callback).Callback(url, placeholderKeyCopy, realResultHash)
		}

		logger.Info(functionName.String()+" success", "contractAddress", tp.ContractAddress, "control", control.GetHash().Hex(), "ifTrue", ifTrue.GetHash().Hex(), "ifFalse", ifTrue.GetHash().Hex(), "result", resultHash)
	}(controlHashCopy, ifTrueHashCopy, ifFalseHashCopy, placeholderKeyCopy)
	return placeholderCt.Hash[:], gas, nil
}

func Req(utype byte, input []byte, tp *TxParams, _ *CallbackFunc) ([]byte, uint64, error) {
	//solgen: input encrypted
	//solgen: return none
	// Don't remove the next line
	//gas, err := PreProcessOperation1(functionName, utype, ctHash, tp)
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

		// Only the sequencer need to know about this to include in the L1 message
		if tp.ParallelTxHooks != nil {
			tp.ParallelTxHooks.NotifyExistingRes(&key)
		}

		if !value {
			return nil, gas, vm.ErrExecutionReverted
		}
		return nil, gas, nil
	} else if tp.GasEstimation {
		return nil, gas, nil
	} else {
		ct := awaitCtResult(storage, input, tp)
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

func Cast(utype byte, input []byte, toType byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Cast

	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	if !types.IsValidType(fhe.EncryptionType(toType)) {
		logger.Error("invalid type to cast to")
		return nil, 0, vm.ErrExecutionReverted
	}

	castToType := fhe.EncryptionType(toType)
	gas := getGasForPrecompile(functionName, castToType)
	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	placeholderCt, err := createPlaceholder(utype, functionName, input, []byte{toType})
	if err != nil {
		logger.Error(functionName.String()+" failed to create placeholder", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("fn", functionName.String(), "Storing async ciphertext", "placeholderKey", hex.EncodeToString(placeholderCt.Hash))
	}

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	err = storeCipherText(storage, placeholderCt, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed to store async ciphertext", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storage.SetAsyncCtStart(types.Hash(placeholderCt.Hash))
	if err != nil {
		logger.Error(functionName.String()+" failed to set async value start", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("Starting new async precompiled contract function: " + functionName.String())
	}

	inputCopy := CopySlice(input)
	placeholderKeyCopy := CopySlice(placeholderCt.Hash)

	go func(inputHash, placeholderKey []byte, toType byte) {
		ct, err := blockUntilInputsAvailable(storage, tp, inputHash)
		input := ct[0]
		if err != nil {
			logger.Error(functionName.String()+": input not verified, len: ", len(inputHash), " err: ", err)
			return
		}
		result, err := input.Cast(castToType)
		if err != nil {
			logger.Error("failed to cast to type "+UtypeToString(toType), " err ", err)
			return
		}
		realResultHash, err := hex.DecodeString(result.GetHash().Hex())
		if err != nil {
			logger.Error(functionName.String()+" failed to decode result hash", "err", err)
			return
		}
		result.Hash = realResultHash
		err = storeCipherText(storage, result, tp.ContractAddress)
		if err != nil {
			logger.Error(functionName.String()+" failed to store result", "err", err)
			return
		}
		_ = storage.SetAsyncCtDone(types.Hash(realResultHash))
		if callback != nil {
			url := (*callback).CallbackUrl
			(*callback).Callback(url, placeholderKeyCopy, realResultHash)
		}
		logger.Info(functionName.String()+" success", "contractAddress", tp.ContractAddress, "input", hex.EncodeToString(inputHash), "result", hex.EncodeToString(realResultHash))
	}(inputCopy, placeholderKeyCopy, toType)

	return placeholderCt.Hash[:], gas, nil
}

// TrivialEncrypt takes a plaintext number and encrypts it to a _compact_ ciphertext
// using the server/computation key - obviously this doesn't hide any information as the
// number was known plaintext
func TrivialEncrypt(input []byte, toType byte, securityZone int32, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
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

	if len(input) != 32 {
		logger.Error(functionName.String()+" input len must be 32 bytes", " input ", hex.EncodeToString(input), " len ", len(input))
		return nil, gas, vm.ErrExecutionReverted
	}

	placeholderCt, err := createPlaceholder(toType, functionName, input, []byte{toType}, []byte{byte(securityZone)})
	if err != nil {
		logger.Error(functionName.String()+" failed to create placeholder", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCipherText(storage, placeholderCt, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed to store async ciphertext", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	// Check if value is not overflowing the type
	maxOfType := fhe.MaxOfType(uintType)
	if maxOfType == nil {
		logger.Error("failed to create trivially encrypted value, type is not supported.")
		return nil, gas, vm.ErrExecutionReverted
	}

	valueToEncrypt := *new(big.Int).SetBytes(input)

	// If value is bigger than the maximal value that is supported by the type
	if valueToEncrypt.Cmp(maxOfType) > 0 {
		logger.Error("failed to create trivially encrypted value, value is too large for type: ", "value", valueToEncrypt, "type", uintType)
		return nil, gas, vm.ErrExecutionReverted
	}

	err = storage.SetAsyncCtStart(types.Hash(placeholderCt.Hash))
	if err != nil {
		logger.Error(functionName.String()+" failed to set async value start", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	placeholderKeyCopy := CopySlice(placeholderCt.Hash)

	go func(resultHash []byte, toType byte) {
		// we encrypt this using the computation key not the public key. Also, compact to save space in case this gets saved directly
		// to storage
		result, err := fhe.EncryptPlainText(valueToEncrypt, uintType, securityZone)
		if err != nil {
			logger.Error("failed to create trivial encrypted value")
			return
		}

		realResultHash, err := hex.DecodeString(result.GetHash().Hex())
		if err != nil {
			logger.Error(functionName.String()+" failed to decode result hash", "err", err)
			return
		}
		result.Hash = resultHash

		err = storeCipherText(storage, result, tp.ContractAddress)
		if err != nil {
			logger.Error(functionName.String()+" failed to store result", "err", err)
			return
		}

		_ = storage.SetAsyncCtDone(types.Hash(resultHash))
		if callback != nil {
			url := (*callback).CallbackUrl
			(*callback).Callback(url, resultHash, realResultHash)
		}
		logger.Info(functionName.String()+" success", "contractAddress", tp.ContractAddress, "input", hex.EncodeToString(input), "result", hex.EncodeToString(realResultHash))
	}(placeholderKeyCopy, toType)

	return placeholderCt.Hash[:], gas, nil
}

func Div(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Div
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Div, utype, lhsHash, rhsHash, tp, callback)
}

func Gt(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := types.Gt
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Gt, utype, lhsHash, rhsHash, tp, callback)
}

func Gte(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := types.Gte
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Gte, utype, lhsHash, rhsHash, tp, callback)
}

func Rem(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Rem
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Rem, utype, lhsHash, rhsHash, tp, callback)
}

func And(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := types.And
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).And, utype, lhsHash, rhsHash, tp, callback)
}

func Or(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := types.Or
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Or, utype, lhsHash, rhsHash, tp, callback)
}

func Xor(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := types.Xor
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Xor, utype, lhsHash, rhsHash, tp, callback)
}

func Eq(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: bool math
	//solgen: return ebool
	functionName := types.Eq
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Eq, utype, lhsHash, rhsHash, tp, callback)
}

func Ne(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: bool math
	//solgen: return ebool
	functionName := types.Ne
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Ne, utype, lhsHash, rhsHash, tp, callback)
}

func Min(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Min
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Min, utype, lhsHash, rhsHash, tp, callback)
}

func Max(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Max
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Max, utype, lhsHash, rhsHash, tp, callback)
}

func Shl(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Shl
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Shl, utype, lhsHash, rhsHash, tp, callback)
}

func Shr(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Shr
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Shr, utype, lhsHash, rhsHash, tp, callback)
}

func Rol(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Rol
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Rol, utype, lhsHash, rhsHash, tp, callback)
}

func Ror(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Ror
	return ProcessOperation2(functionName, (*fhe.FheEncrypted).Ror, utype, lhsHash, rhsHash, tp, callback)
}

func Not(utype byte, value []byte, tp *TxParams, _ *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Not
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	ct, gas, err := ProcessOperation1(functionName, utype, value, tp)
	if err != nil {
		return nil, gas, vm.ErrExecutionReverted
	}

	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
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

	resultHash := result.GetHash()
	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "input", ct.GetHash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Random(utype byte, seed uint64, securityZone int32, tp *TxParams, _ *CallbackFunc) ([]byte, uint64, error) {
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

	resultHash := result.GetHash()
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

func Square(utype byte, value []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	return Mul(utype, value, value, tp, callback)
	// Please don't delete the below comment, this is intentionally left here for code generation.
	//ct := ProcessOperation1(storage, fhe.BytesToHash(value), tp.ContractAddress)
}
