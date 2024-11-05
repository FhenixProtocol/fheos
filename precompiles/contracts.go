package precompiles

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	storage2 "github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	"math/big"
	"strings"

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

func SealOutput(utype byte, ctHash []byte, pk []byte, tp *TxParams, _ *CallbackFunc) (string, uint64, error) {
	//solgen: bool math
	functionName := types.SealOutput

	ct, gas, err := ProcessOperation1(functionName, utype, ctHash, tp, false)
	if err != nil {
		return "", gas, vm.ErrExecutionReverted
	}

	if len(pk) != 32 {
		msg := functionName.String() + " public key need to be 32 bytes long"
		logger.Error(msg, "public-key", hex.EncodeToString(pk), "len", len(pk))
		return "", 0, vm.ErrExecutionReverted
	}

	if tp.GasEstimation {
		return "0x" + strings.Repeat("00", 370), gas, nil
	}

	reencryptedValue, err := fhe.SealOutput(*ct, pk)
	if err != nil {
		logger.Error(functionName.String()+" failed to encrypt to user key", "err", err)
		return "", 0, vm.ErrExecutionReverted
	}

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "ciphertext-hash ", hex.EncodeToString(ctHash), "public-key", hex.EncodeToString(pk))

	return string(reencryptedValue), gas, nil
}

func Decrypt(utype byte, input []byte, defaultValue *big.Int, tp *TxParams, _ *CallbackFunc) (*big.Int, uint64, error) {
	//solgen: output plaintext
	functionName := types.Decrypt

	ct, gas, err := ProcessOperation1(functionName, utype, input, tp, true)
	if err != nil {
		return nil, gas, vm.ErrExecutionReverted
	}

	if tp.GasEstimation {
		return defaultValue, gas, nil
	}

	decryptedValue, err := fhe.Decrypt(*ct)
	if err != nil {
		logger.Error("failed decrypting ciphertext", "error", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "input", hex.EncodeToString(input))
	return decryptedValue, gas, nil
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

func Select(utype byte, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte, tp *TxParams, _ *CallbackFunc) ([]byte, uint64, error) {
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

	control, ifTrue, ifFalse, err := get3VerifiedOperands(storage, controlHash, ifTrueHash, ifFalseHash, tp)
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

	resultHash := result.GetHash()
	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "control", control.GetHash().Hex(), "ifTrue", ifTrue.GetHash().Hex(), "ifFalse", ifTrue.GetHash().Hex(), "result", resultHash.Hex())
	return resultHash[:], gas, nil
}

func Req(utype byte, input []byte, tp *TxParams, _ *CallbackFunc) ([]byte, uint64, error) {
	//solgen: input encrypted
	//solgen: return none
	functionName := types.Require

	ct, gas, err := ProcessOperation1(functionName, utype, input, tp, false)
	if err != nil {
		return nil, gas, vm.ErrExecutionReverted
	}

	if tp.GasEstimation {
		return nil, gas, nil
	}

	ev, err := evaluateRequire(ct)
	if err != nil {
		msg := functionName.String() + " error on evaluation"
		logger.Error(msg, " err ", err)
		return nil, gas, vm.ErrExecutionReverted
	}

	if !ev {
		msg := functionName.String() + " condition not met"
		logger.Error(msg)
		return nil, gas, vm.ErrExecutionReverted
	}

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "input", hex.EncodeToString(input))

	return nil, gas, nil
}

func Cast(utype byte, input []byte, toType byte, tp *TxParams, _ *CallbackFunc) ([]byte, uint64, error) {
	functionName := types.Cast
	storage := storage2.NewMultiStore(tp.CiphertextDb, &State.Storage)

	if !types.IsValidType(fhe.EncryptionType(toType)) {
		logger.Error("invalid type to cast to")
		return nil, 0, vm.ErrExecutionReverted
	}
	castToType := fhe.EncryptionType(toType)

	ct, gas, err := ProcessOperation1(functionName, utype, input, tp, false)
	if err != nil {
		return nil, gas, vm.ErrExecutionReverted
	}

	if tp.GasEstimation {
		randomHash := State.GetRandomForGasEstimation()
		return randomHash[:], gas, nil
	}

	res, err := ct.Cast(castToType)
	if err != nil {
		msg := fmt.Sprintf("failed to cast to type %s", UtypeToString(toType))
		logger.Error(msg, " type ", castToType)
		return nil, 0, vm.ErrExecutionReverted
	}

	resHash := res.GetHash()

	err = storeCipherText(storage, res, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, 0, vm.ErrExecutionReverted
	}

	logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "ctHash", resHash.Hex())

	return resHash[:], gas, nil
}

// TrivialEncrypt takes a plaintext number and encrypts it to a _compact_ ciphertext
// using the server/computation key - obviously this doesn't hide any information as the
// number was known plaintext
func TrivialEncrypt(input []byte, toType byte, securityZone int32, tp *TxParams, _ *CallbackFunc) ([]byte, uint64, error) {
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
		return nil, gas, vm.ErrExecutionReverted
	}

	valueToEncrypt := *new(big.Int).SetBytes(input)

	var ct *fhe.FheEncrypted
	var err error
	// Check if value is not overflowing the type
	maxOfType := fhe.MaxOfType(uintType)
	if maxOfType == nil {
		logger.Error("failed to create trivially encrypted value, type is not supported.")
		return nil, gas, vm.ErrExecutionReverted
	}

	// If value is bigger than the maximal value that is supported by the type
	if valueToEncrypt.Cmp(maxOfType) > 0 {
		logger.Error("failed to create trivially encrypted value, value is too large for type: ", "value", valueToEncrypt, "type", uintType)
		return nil, gas, vm.ErrExecutionReverted
	}

	// we encrypt this using the computation key not the public key. Also, compact to save space in case this gets saved directly
	// to storage
	ct, err = fhe.EncryptPlainText(valueToEncrypt, uintType, securityZone)
	if err != nil {
		logger.Error("failed to create trivial encrypted value")
		return nil, gas, vm.ErrExecutionReverted
	}

	ctHash := ct.GetHash()
	err = storeCipherText(storage, ct, tp.ContractAddress)
	if err != nil {
		logger.Error(functionName.String()+" failed", "err", err)
		return nil, gas, vm.ErrExecutionReverted
	}

	if shouldPrintPrecompileInfo(tp) {
		logger.Debug(functionName.String()+" success", "contractAddress", tp.ContractAddress, "ctHash", ctHash.Hex(), "valueToEncrypt", valueToEncrypt.Uint64(), "securityZone", ct.SecurityZone)
	}
	return ctHash[:], gas, nil
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

	ct, gas, err := ProcessOperation1(functionName, utype, value, tp, false)
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
	//ct := getCiphertext(storage, fhe.BytesToHash(value), tp.ContractAddress)
}
