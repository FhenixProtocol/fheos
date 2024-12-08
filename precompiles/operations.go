package precompiles

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	"math/big"
)

type TwoOperationFunc func(lhs *fhe.FheEncrypted, rhs *fhe.FheEncrypted) (*fhe.FheEncrypted, error)
type CallbackFunc struct {
	CallbackUrl string
	Callback    func(url string, ctKey []byte, newCtKey []byte)
}

type DecryptCallbackFunc struct {
	CallbackUrl string
	Callback    func(url string, ctKey []byte, plaintext *big.Int)
}

type SealOutputCallbackFunc struct {
	CallbackUrl string
	Callback    func(url string, ctKey []byte, value string)
}

func DecryptHelper(storage *storage2.MultiStore, ctHash fhe.Hash, tp *TxParams, defaultValue *big.Int) (*big.Int, error) {
	ct := awaitCtResult(storage, ctHash, tp)
	if ct == nil {
		msg := "decrypt unverified ciphertext handle"
		logger.Error(msg, " ctHash ", ctHash)
		return defaultValue, vm.ErrExecutionReverted
	}
	plaintext, err := fhe.Decrypt(*ct)
	if err != nil {
		logger.Error("decrypt failed for ciphertext", "error", err)
		return defaultValue, vm.ErrExecutionReverted
	}

	return plaintext, nil
}

func SealOutputHelper(storage *storage2.MultiStore, ctHash fhe.Hash, pk []byte, tp *TxParams) (string, error) {
	ct := awaitCtResult(storage, ctHash, tp)
	if ct == nil {
		msg := "sealOutput unverified ciphertext handle"
		logger.Error(msg, " ctHash ", ctHash)
		return "", vm.ErrExecutionReverted
	}
	sealed, err := fhe.SealOutput(*ct, pk)
	if err != nil {
		logger.Error("sealOutput failed for ciphertext", "error", err)
		return "", vm.ErrExecutionReverted
	}

	return string(sealed), nil
}

func PreProcessOperation1(functionName types.PrecompileName, utype byte, inputBz []byte, tp *TxParams) (fhe.CiphertextKey, uint64, error) {
	input, err := types.DeserializeCiphertextKey(inputBz)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize input ciphertext key", "err", err)
		return types.GetEmptyCiphertextKey(), 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return types.GetEmptyCiphertextKey(), gas, nil
	}

	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return types.GetEmptyCiphertextKey(), gas, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	return input, gas, nil
}

func ProcessOperation1(functionName types.PrecompileName, utype byte, inputBz []byte, tp *TxParams) (*fhe.FheEncrypted, uint64, error) {
	input, gas, err := PreProcessOperation1(functionName, utype, inputBz, tp)
	if err != nil {
		return nil, gas, err
	}

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	ct := awaitCtResult(storage, input.Hash, tp)
	if ct == nil {
		msg := functionName.String() + " unverified ciphertext handle"
		logger.Error(msg, " input ", input)
		return nil, gas, vm.ErrExecutionReverted
	}
	return ct, gas, nil
}

func ProcessOperation2(functionName types.PrecompileName, mathOp TwoOperationFunc, utype byte, lhsKeyBz, rhsKeyBz []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	lhsKey, err := types.DeserializeCiphertextKey(lhsKeyBz)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize lhs ciphertext key", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}
	rhsKey, err := types.DeserializeCiphertextKey(rhsKeyBz)
	if err != nil {
		logger.Error(functionName.String()+" failed to deserialize rhs ciphertext key", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	placeholderCt := fhe.CreateFheEncryptedWithData(CreatePlaceHolderData(), fhe.EncryptionType(utype), true)
	placeholderKey := fhe.CalcBinaryPlaceholderValueHash(&lhsKey, &rhsKey, int(functionName))
	placeholderCt.Key = placeholderKey

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug(functionName.String(), "lhs", hex.EncodeToString(lhsKey.Hash[:]), "rhs", hex.EncodeToString(rhsKey.Hash[:]), "placeholderKey", hex.EncodeToString(placeholderKey.Hash[:]))
	}

	err = storeCipherText(storage, placeholderCt)
	if err != nil {
		logger.Error(functionName.String()+" failed to store async ciphertext", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetRandomKeyForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("fn", functionName.String(), "Storing async ciphertext", "placeholderKey", hex.EncodeToString(placeholderKey.Hash[:]))
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("Starting new precompiled contract function: " + functionName.String())
	}

	placeholderKeyCopy := placeholderKey
	rhsKeyCopy := rhsKey
	lhsKeyCopy := lhsKey

	go func(lhsKey, rhsKey, resultKey *fhe.CiphertextKey) {
		lhs, rhs := blockUntilBinaryOperandsAvailable(storage, lhsKey, rhsKey, tp)

		if lhs == nil || rhs == nil {
			logger.Error(functionName.String() + ": inputs not verified")
			return
		}

		if lhs.UintType != rhs.UintType || lhs.UintType != uintType {
			msg := functionName.String() + " operand type mismatch"
			logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
			return
		}

		result, err2 := mathOp(lhs, rhs)
		if err2 != nil {
			logger.Error(functionName.String()+" failed", "err", err2)
			return
		}

		realResultHash, err3 := hex.DecodeString(result.GetHash().Hex())
		if err3 != nil {
			logger.Error(functionName.String()+" failed", "err", err3)
			return
		}

		result.Key = *resultKey

		err2 = storeCipherText(storage, result)
		if err2 != nil {
			logger.Error(functionName.String()+" failed", "err", err2)
			return
		}

		if callback != nil {
			url := (*callback).CallbackUrl
			(*callback).Callback(url, placeholderKeyCopy.Hash[:], realResultHash)
		}
		logger.Info(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.GetHash().Hex(), "rhs", rhs.GetHash().Hex(), "result", result.GetHash().Hex())
	}(&lhsKeyCopy, &rhsKeyCopy, &placeholderKeyCopy)
	return types.SerializeCiphertextKey(placeholderKey), gas, err
}
