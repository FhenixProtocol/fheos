package precompiles

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
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

func PreProcessOperation1(functionName types.PrecompileName, utype byte, input []byte, tp *TxParams) (uint64, error) {
	uintType := fhe.EncryptionType(utype)

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		return gas, nil
	}

	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return gas, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	if len(input) != 32 {
		msg := functionName.String() + " ct hash len must be 32 bytes"
		logger.Error(msg, "ctHash", hex.EncodeToString(input), "len", len(input))
		return gas, vm.ErrExecutionReverted
	}

	return gas, nil
}

func PreProcessSealOutput(functionName types.PrecompileName, utype byte, ctHash []byte, pk []byte, onResultCallback *SealOutputCallbackFunc) (uint64, error) {
	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return 0, vm.ErrExecutionReverted
	}

	if onResultCallback == nil {
		msg := functionName.String() + " must set callback"
		logger.Error(msg, "callback was not provided")
		return 0, vm.ErrExecutionReverted
	}

	if len(ctHash) != 32 {
		msg := functionName.String() + " ciphertext's hashes need to be 32 bytes long"
		logger.Error(msg, "ciphertext-hash", hex.EncodeToString(ctHash), "hash-len", len(ctHash))
		return 0, vm.ErrExecutionReverted
	}

	if len(pk) != 32 {
		msg := functionName.String() + " public key need to be 32 bytes long"
		logger.Error(msg, "public-key", hex.EncodeToString(pk), "len", len(pk))
		return 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	return gas, nil
}

func ProcessOperation1(functionName types.PrecompileName, utype byte, input []byte, tp *TxParams) (*fhe.FheEncrypted, uint64, error) {
	gas, err := PreProcessOperation1(functionName, utype, input, tp)
	if err != nil {
		return nil, gas, err
	}

	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)
	ct := awaitCtResult(storage, input, tp)
	if ct == nil {
		msg := functionName.String() + " unverified ciphertext handle"
		logger.Error(msg, " input ", input)
		return nil, gas, vm.ErrExecutionReverted
	}
	return ct, gas, nil
}

func ProcessOperation2(functionName types.PrecompileName, mathOp TwoOperationFunc, utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	placeholderCt := fhe.CreateFheEncryptedWithData(CreatePlaceHolderData(), fhe.EncryptionType(utype), true)
	placeholderKey := fhe.CalcBinaryPlaceholderValueHash(lhsHash, rhsHash, int(functionName))
	placeholderCt.Hash = placeholderKey

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug(functionName.String(), "lhs", hex.EncodeToString(lhsHash), "rhs", hex.EncodeToString(rhsHash), "placeholderKey", hex.EncodeToString(placeholderKey))
	}

	err := storeCipherText(storage, placeholderCt, tp.ContractAddress)
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
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("fn", functionName.String(), "Storing async ciphertext", "placeholderKey", hex.EncodeToString(placeholderKey))
	}
	err = storage.SetAsyncCtStart(types.Hash(placeholderKey))
	if err != nil {
		logger.Error(functionName.String()+" failed to set async value start", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("Starting new precompiled contract function: " + functionName.String())
	}

	lhsHashCopy := CopySlice(lhsHash)
	rhsHashCopy := CopySlice(rhsHash)
	placeholderKeyCopy := CopySlice(placeholderKey)

	go func(lhsHash, rhsHash, resultHash []byte) {
		lhs, rhs := blockUntilBinaryOperandsAvailable(storage, lhsHash, rhsHash, tp)

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

		result.Hash = resultHash

		err2 = storeCipherText(storage, result, tp.ContractAddress)
		if err2 != nil {
			logger.Error(functionName.String()+" failed", "err", err2)
			return
		}

		_ = storage.SetAsyncCtDone(types.Hash(resultHash))
		if callback != nil {
			url := (*callback).CallbackUrl
			(*callback).Callback(url, placeholderKeyCopy, realResultHash)
		}
		logger.Info(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.GetHash().Hex(), "rhs", rhs.GetHash().Hex(), "result", result.GetHash().Hex())
	}(lhsHashCopy, rhsHashCopy, placeholderKeyCopy)
	return placeholderKey[:], gas, err
}
