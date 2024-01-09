package precompiles

import (
	"encoding/hex"
	"math/big"
	"runtime"
	"sync"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/sirupsen/logrus"

	tfhe "github.com/fhenixprotocol/go-tfhe"
)

var logger *logrus.Logger

// FHENIX: TODO - persist it somehow
var ctHashMap map[tfhe.Hash]*tfhe.Ciphertext
var ctHashMapLock *sync.RWMutex

func InitLogger() {
	logger = newLogger()
	tfhe.InitLogger(getDefaultLogLevel())
}

func InitTfheConfig(tfheConfig *tfhe.Config) error {
	err := tfhe.InitTfhe(tfheConfig)
	if err != nil {
		logger.Error("Failed to init tfhe config with error: ", err)
		return err
	}

	if ctHashMap == nil {
		ctHashMap = make(map[tfhe.Hash]*tfhe.Ciphertext)
	}
	if ctHashMapLock == nil {
		ctHashMapLock = &sync.RWMutex{}
	}

	logger.Info("Successfully initialized tfhe config to be: ", tfheConfig)

	return nil
}

func shouldPrintPrecompileInfo(tp *TxParams) bool {
	return tp.Commit && !tp.GasEstimation
}

func getFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	return funcName
}

// ============================
func Add(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheAdd inputs not verified ", " err ", err, " input ", hex.EncodeToString(input))
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheAdd operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Add(rhs)
	if err != nil {
		logger.Error("fheAdd failed ", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	importCiphertext(result)

	resultHash := result.Hash()
	logger.Debug("fheAdd success ", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), " result ", resultHash.Hex())
	return resultHash[:], nil
}

func Verify(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	if len(input) <= 1 {
		msg := "verifyCiphertext RequiredGas() input needs to contain a ciphertext and one byte for its type"
		logger.Error(msg, " len ", len(input))
		return nil, vm.ErrExecutionReverted
	}

	ctBytes := input[:len(input)-1]
	ctType := tfhe.UintType(input[len(input)-1])

	ct, err := tfhe.NewCipherTextFromBytes(ctBytes, ctType, true /* TODO: not sure + shouldn't be hardcoded */)
	if err != nil {
		logger.Error("verifyCiphertext failed to deserialize input ciphertext",
			" err ", err,
			" len ", len(ctBytes),
			" ctBytes64 ", hex.EncodeToString(ctBytes[:minInt(len(ctBytes), 64)]))
		return nil, vm.ErrExecutionReverted
	}
	ctHash := ct.Hash()
	importCiphertext(ct)

	if tp.Commit {
		logger.Debug("verifyCiphertext success ",
			" ctHash ", ctHash.Hex(),
			" ctBytes64 ", hex.EncodeToString(ctBytes[:minInt(len(ctBytes), 64)]))
	}
	return ctHash[:], nil
}

func SealOutput(input []byte, tp *TxParams) ([]byte, error) {
	//solgen: bool math
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	if len(input) != 64 {
		msg := "sealOutput input len must be 64 bytes"
		logger.Error(msg, " input ", hex.EncodeToString(input), " len ", len(input))
		return nil, vm.ErrExecutionReverted
	}

	ct := getCiphertext(tfhe.BytesToHash(input[0:32]))
	if ct == nil {
		msg := "sealOutput unverified ciphertext handle"
		logger.Error(msg, " input ", hex.EncodeToString(input))
		return nil, vm.ErrExecutionReverted
	}

	if tp.GasEstimation {
		return []byte{1}, nil
	}

	pubKey := input[32:64]
	reencryptedValue, err := tfhe.SealOutput(*ct, pubKey)
	if err != nil {
		logger.Error("sealOutput failed to encrypt to user key", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	logger.Debug("sealOutput success", " input ", hex.EncodeToString(input))
	return reencryptedValue, nil
}

func Decrypt(input []byte, tp *TxParams) (*big.Int, error) {
	//solgen: output plaintext
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	if len(input) != 32 {
		msg := "decrypt input len must be 32 bytes"
		logger.Error(msg, " input ", hex.EncodeToString(input), " len ", len(input))
		return nil, vm.ErrExecutionReverted
	}

	ct := getCiphertext(tfhe.BytesToHash(input[0:32]))
	if ct == nil {
		msg := "decrypt unverified ciphertext handle"
		logger.Error(msg, " input ", hex.EncodeToString(input))
		return nil, vm.ErrExecutionReverted
	}

	if tp.GasEstimation {
		return new(big.Int).SetUint64(1), nil
	}

	decryptedValue, err := tfhe.Decrypt(*ct)
	if err != nil {
		logger.Error("failed decrypting ciphertext", " error ", err)
		return nil, vm.ErrExecutionReverted
	}

	bgDecrypted := new(big.Int).SetUint64(decryptedValue)

	logger.Debug("decrypt success", " input ", hex.EncodeToString(input))
	return bgDecrypted, nil

}

func Lte(input []byte, tp *TxParams) ([]byte, error) {
	//solgen: return ebool
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheLte inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheLte operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)

	}

	result, err := lhs.Lte(rhs)
	if err != nil {
		logger.Error("fheLte failed ", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	resultHash := result.Hash()
	logger.Debug("fheLte success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), " result ", resultHash.Hex())
	return resultHash[:], nil
}

func Sub(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheSub inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheSub operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// // If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Sub(rhs)
	if err != nil {
		logger.Error("fheSub failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	resultHash := result.Hash()
	logger.Debug("fheSub success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), " result ", resultHash.Hex())
	return resultHash[:], nil
}

func Mul(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheMul inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheMul operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Mul(rhs)
	if err != nil {
		logger.Error("fheMul failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	return ctHash[:], nil
}

func Lt(input []byte, tp *TxParams) ([]byte, error) {
	//solgen: return ebool
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheLt inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheLt operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Lt(rhs)
	if err != nil {
		logger.Error("fheLt failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	resultHash := result.Hash()
	logger.Debug("fheLt success ", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), " result ", resultHash.Hex())
	return resultHash[:], nil
}

func Select(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	control, ifTrue, ifFalse, err := get3VerifiedOperands(input)
	if err != nil {
		logger.Error("select inputs not verified input len: ", len(input), " err: ", err)
		return nil, vm.ErrExecutionReverted
	}

	if ifTrue.UintType != ifFalse.UintType {
		msg := "select operands type mismatch"
		logger.Error(msg, " ifTrue ", ifTrue.UintType, " ifFalse ", ifFalse.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(ifTrue.UintType)
	}

	result, err := control.Cmux(ifTrue, ifFalse)
	if err != nil {
		logger.Error("select failed ", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	resultHash := result.Hash()
	logger.Debug("select success ", " control ", control.Hash().Hex(), " ifTrue ", ifTrue.Hash().Hex(), " ifFalse ", ifTrue.Hash().Hex(), " result ", resultHash.Hex())
	return resultHash[:], nil
}

func Req(input []byte, tp *TxParams) ([]byte, error) {
	//solgen: input encrypted
	//solgen: return none
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	if len(input) != 32 {
		msg := "require input len must be 32 bytes"
		logger.Error(msg, " input ", hex.EncodeToString(input), " len ", len(input))
		return nil, vm.ErrExecutionReverted
	}

	ct := getCiphertext(tfhe.BytesToHash(input))
	if ct == nil {
		msg := "require unverified handle"
		logger.Error(msg, " input ", hex.EncodeToString(input))
		return nil, vm.ErrExecutionReverted
	}
	// If we are not committing to state, assume the require is true, avoiding any side effects
	// (i.e. mutatiting the oracle DB).
	if tp.GasEstimation {
		return nil, nil
	}

	ev := evaluateRequire(ct)

	if !ev {
		msg := "require condition not met"
		logger.Error(msg)
		return nil, vm.ErrExecutionReverted
	}

	return nil, nil
}

func Cast(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	if !isValidType(input[32]) {
		logger.Error("invalid type to cast to")
		return nil, vm.ErrExecutionReverted
	}
	castToType := tfhe.UintType(input[32])

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(castToType)
	}

	ct := getCiphertext(tfhe.BytesToHash(input[0:32]))
	if ct == nil {
		logger.Error("cast input not verified")
		return nil, vm.ErrExecutionReverted
	}

	res, err := ct.Cast(castToType)
	if err != nil {
		msg := "cast Run() error casting ciphertext to"
		logger.Error(msg, " type ", castToType)
		return nil, vm.ErrExecutionReverted
	}

	resHash := res.Hash()

	importCiphertext(res)
	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("cast success",
			" ctHash ", resHash.Hex(),
		)
	}

	return resHash[:], nil
}

func TrivialEncrypt(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	if len(input) != 33 {
		msg := "trivialEncrypt input len must be 33 bytes"
		logger.Error(msg, " input ", hex.EncodeToString(input), " len ", len(input))
		return nil, vm.ErrExecutionReverted
	}

	valueToEncrypt := *new(big.Int).SetBytes(input[0:32])
	encryptToType := tfhe.UintType(input[32])

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(encryptToType)
	}

	ct, err := tfhe.NewCipherTextTrivial(valueToEncrypt, encryptToType)
	if err != nil {
		logger.Error("failed to create trivial encrypted value")
		return nil, vm.ErrExecutionReverted
	}

	ctHash := ct.Hash()
	importCiphertext(ct)
	if shouldPrintPrecompileInfo(tp) {
		logger.Debug("trivialEncrypt success",
			" ctHash ", ctHash.Hex(),
			" valueToEncrypt ", valueToEncrypt.Uint64())
	}
	return ctHash[:], nil
}

func Div(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheMul inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheMul operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Div(rhs)
	if err != nil {
		logger.Error("fheDiv failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheDiv success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), " result ", ctHash.Hex())
	return ctHash[:], nil
}

func Gt(input []byte, tp *TxParams) ([]byte, error) {
	//solgen: return ebool
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheGt inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheGt operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Gt(rhs)
	if err != nil {
		logger.Error("fheGt failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheGt success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), " result ", ctHash.Hex())
	return ctHash[:], nil
}

func Gte(input []byte, tp *TxParams) ([]byte, error) {
	//solgen: return ebool
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheGte inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheGte operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Gte(rhs)
	if err != nil {
		logger.Error("fheGte failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheGte success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), " result ", ctHash.Hex())
	return ctHash[:], nil
}

func Rem(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheRem inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheRem operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Rem(rhs)
	if err != nil {
		logger.Error("fheRem failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheRem success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), " result ", ctHash.Hex())
	return ctHash[:], nil
}

func And(input []byte, tp *TxParams) ([]byte, error) {
	//solgen: bool math
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheAnd inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheAnd operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.And(rhs)
	if err != nil {
		logger.Error("fheAnd failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheAnd success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), " result ", ctHash.Hex())
	return ctHash[:], nil
}

func Or(input []byte, tp *TxParams) ([]byte, error) {
	//solgen: bool math
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheOr inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheOr operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Or(rhs)
	if err != nil {
		logger.Error("fheOr failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheOr success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), " result ", ctHash.Hex())
	return ctHash[:], nil
}

func Xor(input []byte, tp *TxParams) ([]byte, error) {
	//solgen: bool math
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheXor inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheXor operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Xor(rhs)
	if err != nil {
		logger.Error("fheXor failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheXor success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), " result ", ctHash.Hex())
	return ctHash[:], nil
}

func Eq(input []byte, tp *TxParams) ([]byte, error) {
	//solgen: bool math
	//solgen: return ebool
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheEq inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheEq operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Eq(rhs)
	if err != nil {
		logger.Error("fheEq failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheEq success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), " result ", ctHash.Hex())
	return ctHash[:], nil
}

func Ne(input []byte, tp *TxParams) ([]byte, error) {
	//solgen: bool math
	//solgen: return ebool
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheNe inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheNe operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Ne(rhs)
	if err != nil {
		logger.Error("fheNe failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheNe success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], nil
}

func Min(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheMin inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheMin operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Min(rhs)
	if err != nil {
		logger.Error("fheMin failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheMin success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], nil
}

func Max(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheMax inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheMax operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Max(rhs)
	if err != nil {
		logger.Error("fheMax failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheMax success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], nil
}

func Shl(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheShl inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheShl operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Shl(rhs)
	if err != nil {
		logger.Error("fheShl failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheShl success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], nil
}

func Shr(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("Starting new precompiled contract function ", getFunctionName())
	}

	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheShr inputs not verified", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	if lhs.UintType != rhs.UintType {
		msg := "fheShr operand type mismatch"
		logger.Error(msg, " lhs ", lhs.UintType, " rhs ", rhs.UintType)
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(lhs.UintType)
	}

	result, err := lhs.Shr(rhs)
	if err != nil {
		logger.Error("fheShr failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}
	importCiphertext(result)

	ctHash := result.Hash()

	logger.Debug("fheShr success", " lhs ", lhs.Hash().Hex(), " rhs ", rhs.Hash().Hex(), "result", ctHash.Hex())
	return ctHash[:], nil
}

func Not(input []byte, tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new precompiled contract function ", getFunctionName())
	}

	ct := getCiphertext(tfhe.BytesToHash(input[0:32]))
	if ct == nil {
		msg := "not unverified ciphertext handle"
		logger.Error(msg, "input", hex.EncodeToString(input))
		return nil, vm.ErrExecutionReverted
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	if tp.GasEstimation {
		return importRandomCiphertext(ct.UintType)
	}

	result, err := ct.Not()
	if err != nil {
		logger.Error("not failed", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	importCiphertext(result)

	resultHash := result.Hash()
	logger.Debug("fheNot success", " in ", ct.Hash().Hex(), " result ", resultHash.Hex())
	return resultHash[:], nil
}

func GetNetworkPublicKey(tp *TxParams) ([]byte, error) {
	if shouldPrintPrecompileInfo(tp) {
		logger.Info("starting new function get network public key:", getFunctionName())
	}

	pk, err := tfhe.PublicKey()
	if err != nil {
		logger.Error("could not get public key", " err ", err)
		return nil, vm.ErrExecutionReverted
	}

	return pk, nil
}
