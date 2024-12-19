package precompiles

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

// OperationFunc is a generic function type that can handle variable number of operands
type OperationFunc interface {
	Execute(inputs []*fhe.FheEncrypted) (*fhe.FheEncrypted, error)
	ValidateTypes(inputs []*fhe.FheEncrypted, utype byte) error
}

// OneOperationFunc wraps a single operand function to match OperationFunc interface
type OneOperationFunc func(input *fhe.FheEncrypted) (*fhe.FheEncrypted, error)

// TwoOperationFunc wraps a two operand function to match OperationFunc interface
type TwoOperationFunc func(first, second *fhe.FheEncrypted) (*fhe.FheEncrypted, error)

// ThreeOperationFunc wraps a three operand function to match OperationFunc interface
type ThreeOperationFunc struct {
	Fn               func(control, ifTrue, ifFalse *fhe.FheEncrypted) (*fhe.FheEncrypted, error)
	CustomValidation func(inputs []*fhe.FheEncrypted, utype byte) error
}

func (f OneOperationFunc) Execute(inputs []*fhe.FheEncrypted) (*fhe.FheEncrypted, error) {
	if len(inputs) != 1 {
		return nil, fmt.Errorf("expected 1 input, got %d", len(inputs))
	}
	return f(inputs[0])
}

func (f OneOperationFunc) ValidateTypes(inputs []*fhe.FheEncrypted, _ byte) error {
	if len(inputs) != 1 {
		return fmt.Errorf("expected 1 input, got %d", len(inputs))
	}

	return nil
}

func (f TwoOperationFunc) Execute(inputs []*fhe.FheEncrypted) (*fhe.FheEncrypted, error) {
	if len(inputs) != 2 {
		return nil, fmt.Errorf("expected 2 inputs, got %d", len(inputs))
	}
	return f(inputs[0], inputs[1])
}

func (f TwoOperationFunc) ValidateTypes(inputs []*fhe.FheEncrypted, _ byte) error {
	if len(inputs) != 2 {
		return fmt.Errorf("expected 2 inputs, got %d", len(inputs))
	}

	if inputs[0].Key.UintType != inputs[1].Key.UintType {
		return fmt.Errorf("inputs type mismatch: expected %v, got %v", inputs[0].Key.UintType.ToString(), inputs[1].Key.UintType.ToString())
	}

	return nil
}

func (f ThreeOperationFunc) Execute(inputs []*fhe.FheEncrypted) (*fhe.FheEncrypted, error) {
	if len(inputs) != 3 {
		return nil, fmt.Errorf("expected 3 inputs, got %d", len(inputs))
	}
	return f.Fn(inputs[0], inputs[1], inputs[2])
}

func (f ThreeOperationFunc) ValidateTypes(inputs []*fhe.FheEncrypted, utype byte) error {
	if f.CustomValidation != nil {
		return f.CustomValidation(inputs, utype)
	}
	return validateAllSameType(inputs, utype)
}

// Helper function for default type validation
func validateAllSameType(inputs []*fhe.FheEncrypted, utype byte) error {
	expectedType := fhe.EncryptionType(utype)
	for i, ct := range inputs {
		if ct.UintType != expectedType {
			return fmt.Errorf("input %d type mismatch: expected %v, got %v", i, expectedType.ToString(), ct.UintType.ToString())
		}
	}
	return nil
}

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

func inputsToString(inputKeys []fhe.CiphertextKey) string {
	var concatenated string
	for i, key := range inputKeys {
		concatenated += fmt.Sprintf("input%d: %s, ", i, hex.EncodeToString(key.Hash[:]))
	}
	return concatenated
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
func createPlaceholder(utype byte, securityZone int32, functionName types.PrecompileName, inputKeys ...[]byte) (*fhe.FheEncrypted, error) {
	placeholderCt := fhe.CreateFheEncryptedWithData(CreatePlaceHolderData(), fhe.EncryptionType(utype), true)

	// Calculate placeholder based on number of inputs
	placeholderKey := fhe.CalcPlaceholderValueHash(int(functionName), fhe.EncryptionType(utype), securityZone, inputKeys...)

	placeholderCt.Key = placeholderKey
	return placeholderCt, nil
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

func SolidityInputsToCiphertextKeys(inputs ...[]byte) ([]fhe.CiphertextKey, error) {
	var inputKeys []fhe.CiphertextKey
	for _, input := range inputs {
		key, err := types.DeserializeCiphertextKey(input)
		if err != nil {
			return nil, err
		}
		inputKeys = append(inputKeys, key)
	}
	return inputKeys, nil
}

func keysToHashes(keys []fhe.CiphertextKey) [][]byte {
	hashes := make([][]byte, len(keys))
	for i, key := range keys {
		hashes[i] = make([]byte, len(key.Hash))
		copy(hashes[i], key.Hash[:])
	}
	return hashes
}

func getUtypeForFunctionName(functionName types.PrecompileName, currentType byte) byte {
	switch functionName {
	case types.Lte, types.Lt, types.Gte, types.Gt, types.Eq, types.Ne:
		return byte(fhe.Bool)
	default:
		return currentType
	}
}

// ProcessOperation handles operations with variable number of inputs
func ProcessOperation(functionName types.PrecompileName, operation OperationFunc, utype byte, securtiyZone int32, inputKeys []fhe.CiphertextKey, tp *TxParams, callback *CallbackFunc) ([]byte, uint64, error) {
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	placeholderCt, err := createPlaceholder(getUtypeForFunctionName(functionName, utype), securtiyZone, functionName, keysToHashes(inputKeys)...)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug(functionName.String(), inputsToString(inputKeys), "placeholderKey", hex.EncodeToString(placeholderCt.Key.Hash[:]))
	}

	uintType := fhe.EncryptionType(utype)
	if !types.IsValidType(uintType) {
		logger.Error("invalid ciphertext", "type", utype)
		return nil, 0, vm.ErrExecutionReverted
	}

	if err := storeCipherText(storage, placeholderCt); err != nil {
		logger.Error(functionName.String()+" failed to store async ciphertext", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	gas := getGasForPrecompile(functionName, uintType)
	if tp.GasEstimation {
		randomHash := State.GetEmptyKeyForGasEstimation()
		return randomHash[:], gas, nil
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("fn", functionName.String(), "Storing async ciphertext", "placeholderKey", hex.EncodeToString(placeholderCt.Key.Hash[:]))
	}

	// Make copies for goroutine
	copiedInputs := make([]fhe.CiphertextKey, len(inputKeys))
	copy(copiedInputs, inputKeys)
	placeholderKeyCopy := placeholderCt.Key

	go func(inputs []fhe.CiphertextKey, resultKey fhe.CiphertextKey) {
		cleanup := func() { deleteCipherText(storage, resultKey.Hash) }
		defer cleanup()
		cts, err := blockUntilInputsAvailable(storage, tp, inputs...)
		if err != nil || len(cts) != len(inputs) {
			logger.Error(functionName.String() + ": inputs not verified")
			return
		}

		for i, ct := range cts {
			if ct == nil {
				logger.Error(functionName.String()+": input not verified", "index", i)
				return
			}
		}

		// Use the operation's custom type validation
		if err := operation.ValidateTypes(cts, utype); err != nil {
			logger.Error(functionName.String()+" type validation failed", "err", err)
			return
		}

		result, err := operation.Execute(cts)
		if err != nil {
			logger.Error(functionName.String()+" failed", "err", err)
			return
		}

		realResultHash, err := hex.DecodeString(result.GetHash().Hex())
		if err != nil {
			logger.Error(functionName.String()+" failed", "err", err)
			return
		}

		result.Key = resultKey

		err = storeCipherText(storage, result)
		if err != nil {
			logger.Error(functionName.String()+" failed", "err", err)
			return
		}
		defer func() { cleanup = func() {} }() // Cancel the cleanup

		if callback != nil {
			url := (*callback).CallbackUrl
			(*callback).Callback(url, placeholderKeyCopy.Hash[:], realResultHash)
		}

		// Log success with all input hashes and result
		logFields := []interface{}{
			"contractAddress", tp.ContractAddress,
			"result", result.GetHash().Hex(),
		}
		for i, ct := range cts {
			logFields = append(logFields, fmt.Sprintf("input%d", i), ct.GetHash().Hex())
		}
		logger.Info("["+functionName.String()+"]: success", logFields...)
	}(copiedInputs, placeholderKeyCopy)

	return types.SerializeCiphertextKey(placeholderCt.Key), gas, nil
}
