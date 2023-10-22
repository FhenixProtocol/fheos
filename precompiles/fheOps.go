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
var verifiedCiphertexts map[tfhe.Hash]*verifiedCiphertext

func SetEvmInterpreter(i *vm.EVMInterpreter) {
	if verifiedCiphertexts == nil {
		verifiedCiphertexts = make(map[tfhe.Hash]*verifiedCiphertext)
	}

	interpreter = i
}

func TrivialEncrypt(input []byte) ([]byte, error) {
	tfheConfig := tfhe.Config{
		IsOracle:             true,
		OracleType:           "local",
		OracleDbPath:         "data/oracle.db",
		OracleAddress:        "http://127.0.0.1:9001",
		ServerKeyPath:        "keys/tfhe/sks",
		ClientKeyPath:        "keys/tfhe/cks",
		PublicKeyPath:        "keys/tfhe/pks",
		OraclePrivateKeyPath: "keys/oracle/private-oracle.key",
		OraclePublicKeyPath:  "keys/oracle/public-oracle.key",
	}

	err := tfhe.InitTfhe(&tfheConfig)
	if err != nil {
		return nil, err
	}

	logger := interpreter.GetEVM().Logger
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

func Add(input []byte, inputLen uint32) ([]byte, error) {
	if interpreter == nil {
		msg := "fheAdd no evm interpreter"
		// logger.Error(msg, "lhs", lhs.ciphertext.UintType, "rhs", rhs.ciphertext.UintType)
		return []byte{}, errors.New(msg)
	}

	tfheConfig := tfhe.Config{
		IsOracle:             true,
		OracleType:           "local",
		OracleDbPath:         "data/oracle.db",
		OracleAddress:        "http://127.0.0.1:9001",
		ServerKeyPath:        "keys/tfhe/sks",
		ClientKeyPath:        "keys/tfhe/cks",
		PublicKeyPath:        "keys/tfhe/pks",
		OraclePrivateKeyPath: "keys/oracle/private-oracle.key",
		OraclePublicKeyPath:  "keys/oracle/public-oracle.key",
	}

	err := tfhe.InitTfhe(&tfheConfig)
	if err != nil {
		return []byte{}, err
	}

	logger := interpreter.GetEVM().Logger
	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheAdd inputs not verified", "err", err, "input", hex.EncodeToString(input))
		return []byte{}, err
	}

	if lhs.ciphertext.UintType != rhs.ciphertext.UintType {
		msg := "fheAdd operand type mismatch"
		logger.Error(msg, "lhs", lhs.ciphertext.UintType, "rhs", rhs.ciphertext.UintType)
		return []byte{}, errors.New(msg)
	}

	result, err := lhs.ciphertext.Add(rhs.ciphertext)
	if err != nil {
		logger.Error("fheAdd failed", "err", err)
		return []byte{}, err
	}

	importCiphertext(result)

	// TODO: for testing
	err = os.WriteFile("/tmp/add_result", result.Serialization, 0644)
	if err != nil {
		logger.Error("fheAdd failed to write /tmp/add_result", "err", err)
		return []byte{}, err
	}

	resultHash := result.Hash()
	logger.Info("fheAdd success", "lhs", lhs.ciphertext.Hash().Hex(), "rhs", rhs.ciphertext.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], nil
}

func Reencrypt(input []byte, inputLen uint32) ([]byte, error) {
	logger := interpreter.GetEVM().Logger
	if len(input) != 64 {
		msg := "reencrypt input len must be 64 bytes"
		logger.Error(msg, "input", hex.EncodeToString(input), "len", len(input))
		return nil, errors.New(msg)
	}
	ct := getVerifiedCiphertext(tfhe.BytesToHash(input[0:32]))
	if ct != nil {
		decryptedValue, err := ct.ciphertext.Decrypt()
		if err != nil {
			panic("Failed to decrypt ciphertext during Run")
		}

		pubKey := input[32:64]
		reencryptedValue, err := encryptToUserKey(decryptedValue, pubKey)
		if err != nil {
			logger.Error("reencrypt failed to encrypt to user key", "err", err)
			return nil, err
		}
		logger.Info("reencrypt success", "input", hex.EncodeToString(input))
		// FHENIX - Previously it was return toEVMBytes(reencryptedValue), nil but the decrypt function in Fhevm din't support it so we removed the the toEVMBytes
		return reencryptedValue, nil
	}
	msg := "reencrypt unverified ciphertext handle"
	logger.Error(msg, "input", hex.EncodeToString(input))
	return nil, errors.New(msg)
}
