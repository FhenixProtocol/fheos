package precompiles

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	bridge_types "github.com/fhenixprotocol/warp-drive/fhe-bridge/go/bridgetypes"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
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
	addOp := TwoOperationFunc((*fhe.FheEncrypted).Add)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, addOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

// Verify takes inputs from the user and runs them through verification. Note that we will always get ciphertexts that
// are public-key encrypted and compressed. Anything else will fail
func Verify(utype byte, input []byte, securityZone int32, tp *TxParams, _ *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Verify

	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetEmptyKeyForGasEstimation()
		return randomHash[:], gas, nil
	}

	ct := fhe.NewFheEncryptedFromBytes(
		input,
		uintType,
		true,
		false, // TODO: not sure + shouldn't be hardcoded
		securityZone,
		false,
	)
	hash := adjustHashForMetadata(ct.Key.Hash[:], utype, securityZone, false)
	if hash == nil {
		return nil, 0, vm.ErrExecutionReverted
	}
	copy(ct.Key.Hash[:], hash)

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	err := ct.Verify()
	if err != nil {
		logger.Info(fmt.Sprintf("failed to verify ciphertext %s for type %d - was input corrupted?", ct.GetHash().Hex(), uintType))
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storeCiphertext(storage, &ct)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "ctHash", ct.GetHash().Hex())

	retValue := ct.GetKey().Hash
	return retValue[:], gas, nil
}

func SealOutput(utype byte, inputBz []byte, pk []byte, tp *TxParams, onResultCallback *SealOutputCallbackFunc) (string, uint64, error) {
	//solgen: bool math
	functionName := types.SealOutput
	input, gas, err := PreProcessOperation1(functionName, utype, inputBz, tp)
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
			sealed, err := SealOutputHelper(storage, input.Hash, pk, tp, 0, "")
			return sealed, gas, err
		}

		go func(ctHash [common.HashLength]byte) {
			url := (*onResultCallback).CallbackUrl
			transactionHash := (*onResultCallback).TransactionHash
			chainId := (*onResultCallback).ChainId
			sealed, err := SealOutputHelper(storage, ctHash, pk, tp, chainId, transactionHash)
			if err != nil {
				logger.Error("failed sealing output", "error", err)
				return
			}

			logger.Info("SealOutput callback", "url", url, "ctHash", hex.EncodeToString(ctHash[:]), "pk", hex.EncodeToString(pk), "value", string(sealed), "transactionHash", transactionHash, "chainId", chainId)
			(*onResultCallback).Callback(url, ctHash[:], pk, string(sealed), transactionHash, chainId)
		}(input.Hash)
		logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "ctHash", hex.EncodeToString(input.Hash[:]))
	}

	return "0x" + strings.Repeat("00", 370), gas, nil
}

func Decrypt(utype byte, inputBz []byte, defaultValue *big.Int, tp *TxParams, onResultCallback *DecryptCallbackFunc) (*big.Int, uint64, error) {
	//solgen: output plaintext
	functionName := types.Decrypt

	input, gas, err := PreProcessOperation1(functionName, utype, inputBz, tp)
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
			plaintext, err := DecryptHelper(storage, input.Hash, tp, defaultValue, 0, "")
			return plaintext, gas, err
		}
		go func(ctHash [common.HashLength]byte) {
			url := (*onResultCallback).CallbackUrl
			transactionHash := (*onResultCallback).TransactionHash
			chainId := (*onResultCallback).ChainId

			plaintext, err := DecryptHelper(storage, ctHash, tp, defaultValue, chainId, transactionHash)
			if err != nil {
				logger.Error("failed decrypting ciphertext", "error", err)
				return
			}

			(*onResultCallback).Callback(url, ctHash[:], plaintext, transactionHash, chainId)
		}(input.Hash)
		logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "input", hex.EncodeToString(input.Hash[:]))
	}

	return defaultValue, gas, nil
}

func Lte(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := types.Lte
	lteOp := TwoOperationFunc((*fhe.FheEncrypted).Lte)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, lteOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Sub(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Sub
	subOp := TwoOperationFunc((*fhe.FheEncrypted).Sub)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, subOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Mul(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Mul
	mulOp := TwoOperationFunc((*fhe.FheEncrypted).Mul)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, mulOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Lt(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := types.Lt
	ltOp := TwoOperationFunc((*fhe.FheEncrypted).Lt)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, ltOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Select(utype byte, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Select

	selectOp := ThreeOperationFunc{
		Fn: (*fhe.FheEncrypted).Select,
		CustomValidation: func(inputs []*fhe.FheEncrypted, utype byte) error {
			// Only validate that ifTrue and ifFalse have matching types
			if inputs[1].UintType != inputs[2].UintType {
				return fmt.Errorf("operands type mismatch: ifTrue=%v, ifFalse=%v",
					inputs[1].UintType.ToString(), inputs[2].UintType.ToString())
			}
			return nil
		},
	}

	keys, err := SolidityInputsToCiphertextKeys(controlHash, ifTrueHash, ifFalseHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, selectOp, utype, keys[0].SecurityZone, keys, tp, callback)
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
	var inputSer [common.HashLength]byte
	copy(inputSer[:], input)
	ctHash := fhe.BytesToHash(inputSer)
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
		ct := awaitCtResult(storage, inputSer, tp)
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
		randomHash := State.GetEmptyKeyForGasEstimation()
		return randomHash[:], gas, nil
	}

	keys, err := SolidityInputsToCiphertextKeys(input)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	placeholderCt, err := createPlaceholder(toType, keys[0].SecurityZone, functionName, keys[0].Hash[:], ByteToUint256(toType))
	if err != nil {
		logger.Error(functionName.String()+" failed to create placeholder", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("fn", functionName.String(), "Storing async ciphertext", "placeholderKey", hex.EncodeToString(placeholderCt.Key.Hash[:]))
	}

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	err = storeCiphertext(storage, placeholderCt)
	if err != nil {
		logger.Error(functionName.String()+" failed to store async ciphertext", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	logger.Info(functionName.String(), "stored async ciphertext", "placeholderKey", hex.EncodeToString(placeholderCt.Key.Hash[:]))

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("Starting new async precompiled contract function: " + functionName.String())
	}

	keyCopy := keys[0]
	placeholderKeyCopy := placeholderCt.Key

	go func(inputKey, placeholderKey fhe.CiphertextKey, toType byte) {
		ctReady := false
		defer func() {
			if !ctReady {
				logger.Error(functionName.String() + ": failed, deleting placeholder ciphertext " + hex.EncodeToString(placeholderKey.Hash[:]))
				deleteCiphertext(storage, placeholderKey.Hash)
			}
		}()
		ct, err := blockUntilInputsAvailable(storage, tp, inputKey)
		input := ct[0]
		if err != nil {
			logger.Error(functionName.String()+": input not verified, len: ", len(inputKey.Hash), " err: ", err)
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
		result.Key = placeholderKey
		err = storeCiphertext(storage, result)
		if err != nil {
			logger.Error(functionName.String()+" failed to store result", "err", err)
			return
		}
		ctReady = true // Mark as ready

		if callback != nil {
			url := (*callback).CallbackUrl
			(*callback).Callback(url, placeholderKeyCopy.Hash[:], realResultHash)
		}
		logger.Info(functionName.String()+" success", "contractAddress", tp.ContractAddress, "input", hex.EncodeToString(inputKey.Hash[:]), "result", hex.EncodeToString(realResultHash))
	}(keyCopy, placeholderKeyCopy, toType)

	return fhe.SerializeCiphertextKey(placeholderCt.Key), gas, nil
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
		randomHash := State.GetEmptyKeyForGasEstimation()
		return randomHash[:], gas, nil
	}

	placeholderCt, err := createPlaceholder(toType, securityZone, functionName, input, ByteToUint256(toType), Int32ToUint256(securityZone))
	placeholderCt.Key.IsTriviallyEncrypted = true

	if err != nil {
		logger.Error(functionName.String()+" failed to create placeholder", "err", err)
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

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	err = storeCiphertext(storage, placeholderCt)
	if err != nil {
		logger.Error(functionName.String()+" failed to store async ciphertext", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	logger.Info(functionName.String()+" stored async ciphertext", "placeholderKey", hex.EncodeToString(placeholderCt.Key.Hash[:]))

	placeholderKeyCopy := placeholderCt.Key

	go func(resultKey fhe.CiphertextKey, toType byte) {
		ctReady := false
		defer func() {
			if !ctReady {
				logger.Error(functionName.String() + ": failed, deleting placeholder ciphertext " + hex.EncodeToString(resultKey.Hash[:]))
				deleteCiphertext(storage, resultKey.Hash)
			}
		}()
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
		result.Key = resultKey

		err = storeCiphertext(storage, result)
		if err != nil {
			logger.Error(functionName.String()+" failed to store result", "err", err)
			return
		}

		ctReady = true // Mark as ready

		if callback != nil {
			url := (*callback).CallbackUrl
			(*callback).Callback(url, resultKey.Hash[:], realResultHash)
		}
		logger.Info(functionName.String()+" success", "contractAddress", tp.ContractAddress, "input", hex.EncodeToString(input), "result", hex.EncodeToString(realResultHash))
	}(placeholderKeyCopy, toType)

	return fhe.SerializeCiphertextKey(placeholderCt.Key), gas, nil
}

func Div(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Div
	divOp := TwoOperationFunc((*fhe.FheEncrypted).Div)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, divOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Gt(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := types.Gt
	gtOp := TwoOperationFunc((*fhe.FheEncrypted).Gt)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, gtOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Gte(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: return ebool
	functionName := types.Gte
	gteOp := TwoOperationFunc((*fhe.FheEncrypted).Gte)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, gteOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Rem(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Rem
	remOp := TwoOperationFunc((*fhe.FheEncrypted).Rem)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, remOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func And(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := types.And
	andOp := TwoOperationFunc((*fhe.FheEncrypted).And)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, andOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Or(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := types.Or
	orOp := TwoOperationFunc((*fhe.FheEncrypted).Or)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, orOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Xor(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: bool math
	functionName := types.Xor
	xorOp := TwoOperationFunc((*fhe.FheEncrypted).Xor)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, xorOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Eq(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: bool math
	//solgen: return ebool
	functionName := types.Eq
	eqOp := TwoOperationFunc((*fhe.FheEncrypted).Eq)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, eqOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Ne(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	//solgen: bool math
	//solgen: return ebool
	functionName := types.Ne
	neOp := TwoOperationFunc((*fhe.FheEncrypted).Ne)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, neOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Min(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Min
	minOp := TwoOperationFunc((*fhe.FheEncrypted).Min)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, minOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Max(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Max
	maxOp := TwoOperationFunc((*fhe.FheEncrypted).Max)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, maxOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Shl(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Shl
	shlOp := TwoOperationFunc((*fhe.FheEncrypted).Shl)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, shlOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Shr(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Shr
	shrOp := TwoOperationFunc((*fhe.FheEncrypted).Shr)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, shrOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Rol(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Rol
	rolOp := TwoOperationFunc((*fhe.FheEncrypted).Rol)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, rolOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Ror(utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Ror
	rorOp := TwoOperationFunc((*fhe.FheEncrypted).Ror)

	keys, err := SolidityInputsToCiphertextKeys(lhsHash, rhsHash)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, rorOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Not(utype byte, value []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Not
	notOp := OneOperationFunc((*fhe.FheEncrypted).Not)

	keys, err := SolidityInputsToCiphertextKeys(value)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize inputs", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	return ProcessOperation(functionName, notOp, utype, keys[0].SecurityZone, keys, tp, callback)
}

func Random(utype byte, seed uint64, securityZone int32, tp *TxParams, callback *CallbackFunc, fullSeed *big.Int) ([]byte, uint64, error) {
	functionName := types.Random

	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid random output type", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetEmptyKeyForGasEstimation()
		return randomHash[:], gas, nil
	}

	// todo (eshel) verify that the task manager creates the same placeholder
	placeholderCt, err := createPlaceholder(getUtypeForFunctionName(functionName, utype), securityZone, functionName, fullSeed.Bytes())
	if err != nil {
		logger.Error(functionName.String()+" failed to create placeholder", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("fn", functionName.String(), "Storing async ciphertext", "placeholderKey", hex.EncodeToString(placeholderCt.Key.Hash[:]))
	}

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	if err = storeCiphertext(storage, placeholderCt); err != nil {
		logger.Error(functionName.String()+" failed to store async ciphertext", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	// ====== eshel: verified up to here =======

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new async precompiled contract function: " + functionName.String())
	}

	placeholderKeyCopy := placeholderCt.Key

	go func(securityZone int32, uintType fhe.EncryptionType, seed uint64, resultKey fhe.CiphertextKey) {
		ctReady := false
		defer func() {
			if !ctReady {
				logger.Error(functionName.String() + ": failed, deleting placeholder ciphertext " + hex.EncodeToString(resultKey.Hash[:]))
				deleteCiphertext(storage, resultKey.Hash)
			}
		}()

		var finalSeed uint64
		if seed != 0 {
			finalSeed = seed
		} else {
			// Deprecated: Leaving this here for backwards compatibility,
			// in practice the seed should always be provided

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
			logger.Error(functionName.String()+" failed to generate random value", "err", err)
			return
		}

		realResultHash, err := hex.DecodeString(result.GetHash().Hex())
		if err != nil {
			logger.Error(functionName.String()+" failed to decode result hash", "err", err)
			return
		}

		result.Key = resultKey

		err = storeCiphertext(storage, result)
		if err != nil {
			logger.Error(functionName.String()+" failed to store result", "err", err)
			return
		}

		ctReady = true // Mark as ready

		if callback != nil {
			url := (*callback).CallbackUrl
			(*callback).Callback(url, placeholderKeyCopy.Hash[:], realResultHash)
		}

		logger.Info(functionName.String()+" success", "seed", finalSeed, "securityZone", securityZone, "result", hex.EncodeToString(realResultHash))
	}(securityZone, uintType, seed, placeholderKeyCopy)

	return types.SerializeCiphertextKey(placeholderCt.Key), gas, nil
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

func GetCT(hash []byte, tp *TxParams) (*bridge_types.FheEncrypted, error) {
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	ctHash := fhe.Hash(hash)
	ct, err := getCiphertext(storage, ctHash, true)
	if err != nil {
		return nil, err
	}

	if ct.IsPlaceholderValue() {
		return nil, errors.New("ciphertext is a placeholder value")
	}

	bct, err := fhe.ExpandCompressedValue(ct)
	if err != nil {
		return nil, err
	}

	return bct, nil
}
