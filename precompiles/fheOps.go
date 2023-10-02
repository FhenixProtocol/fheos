package precompiles

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/core/vm"
	tfhe "github.com/fhenixprotocol/go-tfhe"
)

var interpreter *vm.EVMInterpreter

func SetEvmInterpreter(i *vm.EVMInterpreter) {
	interpreter = i
}

func Lior(a uint32, b uint32) (uint32, error) {
	return a * b, nil
}

func Moshe(intput []byte, inputLen uint32) ([1][32]byte, error) {
	// byteSlice := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	// // Convert []byte to [1][32]byte
	// var byteArray [1][32]byte

	// copy(byteArray[0][:], byteSlice)

	// return byteArray, nil
	println("TOMMM load keys in (e *fheAdd) Run()")
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
		fmt.Printf("TOMMM inside error: %s", err)
		return [1][32]byte{}, err
	}

	logger := interpreter.GetEVM().Logger
	lhs, rhs, err := get2VerifiedOperands(accessibleState, input)
	if err != nil {
		logger.Error("fheAdd inputs not verified", "err", err, "input", hex.EncodeToString(input))
		return nil, err
	}

	if lhs.ciphertext.UintType != rhs.ciphertext.UintType {
		msg := "fheAdd operand type mismatch"
		logger.Error(msg, "lhs", lhs.ciphertext.UintType, "rhs", rhs.ciphertext.UintType)
		return nil, errors.New(msg)
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	// if !accessibleState.Interpreter().evm.Commit && !accessibleState.Interpreter().evm.EthCall {
	// 	return importRandomCiphertext(accessibleState, lhs.ciphertext.UintType), nil
	// }

	result, err := lhs.ciphertext.Add(rhs.ciphertext)
	if err != nil {
		logger.Error("fheAdd failed", "err", err)
		return nil, err
	}
	importCiphertext(accessibleState, result)

	// TODO: for testing
	err = os.WriteFile("/tmp/add_result", result.Serialization, 0644)
	if err != nil {
		logger.Error("fheAdd failed to write /tmp/add_result", "err", err)
		return nil, err
	}

	resultHash := result.Hash()
	logger.Info("fheAdd success", "lhs", lhs.ciphertext.Hash().Hex(), "rhs", rhs.ciphertext.Hash().Hex(), "result", resultHash.Hex())
	return resultHash[:], nil
}
