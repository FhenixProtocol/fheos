package precompiles

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

func ProcessOperation2(functionName types.PrecompileName, mathOp func(lhs *fhe.FheEncrypted, rhs *fhe.FheEncrypted) (*fhe.FheEncrypted, error), utype byte, lhsHash []byte, rhsHash []byte, tp *TxParams) ([]byte, uint64, error) {
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	placeholderCt := fhe.CreateFheEncryptedWithData(CreatePlaceHolderData(), fhe.EncryptionType(utype), true)
	placeholderKey := fhe.CalcBinaryPlaceholderValueHash(lhsHash, rhsHash, int(functionName))
	placeholderCt.SetHash(placeholderKey)

	logger.Debug("ProcessOperation: ", "lhs", hex.EncodeToString(lhsHash), "rhs", hex.EncodeToString(rhsHash), "placeholderKey", hex.EncodeToString(placeholderKey))

	err := storeCipherText(storage, placeholderCt, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed to store async ciphertext", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	err = storage.SetAsyncCtStart(types.Hash(placeholderKey))
	if err != nil {
		logger.Error(functionName.String()+" failed to set async value start", "err", err)
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
		logger.Info("Starting new precompiled contract function: " + functionName.String())
	}

	lhsHashCopy := CopySlice(lhsHash)
	rhsHashCopy := CopySlice(rhsHash)
	placeholderKeyCopy := CopySlice(placeholderKey)

	go func(lhsHash, rhsHash, resultHash []byte) {
		println("******** AWAITING LHS & RHS: ")
		lhs, rhs := blockUntilBinaryOperandsAvailable(storage, lhsHash, rhsHash, tp)

		if err != nil {
			logger.Error(functionName.String()+": inputs not verified", "err", err)
			tp.ErrChannel <- vm.ErrExecutionReverted
			return
		}

		if lhs.UintType != rhs.UintType || lhs.UintType != uintType {
			msg := functionName.String() + " operand type mismatch"
			logger.Error(msg, "lhs", lhs.UintType, "rhs", rhs.UintType)
			tp.ErrChannel <- vm.ErrExecutionReverted
			return
		}

		result, err2 := mathOp(lhs, rhs)
		if err2 != nil {
			logger.Error(functionName.String()+" failed", "err", err2)
			tp.ErrChannel <- vm.ErrExecutionReverted
			return
		}

		result.SetHash(resultHash)

		err2 = storeCipherText(storage, result, tp.ContractAddress)
		if err2 != nil {
			logger.Error(functionName.String()+" failed", "err", err2)
			tp.ErrChannel <- vm.ErrExecutionReverted
			return
		}

		_ = storage.SetAsyncCtDone(types.Hash(resultHash))

		logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "lhs", lhs.Hash().Hex(), "rhs", rhs.Hash().Hex(), "result", result.Hash().Hex())
	}(lhsHashCopy, rhsHashCopy, placeholderKeyCopy)

	return placeholderKey[:], gas, err
}
