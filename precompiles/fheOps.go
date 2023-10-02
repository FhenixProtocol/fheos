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

func isVerifiedAtCurrentDepth(ct *vm.VerifiedCiphertext) bool {
	return ct.GetVerifiedDepths().Has(interpreter.GetEVM().GetDepth())
}

func getVerifiedCiphertextFromEVM(ciphertextHash tfhe.Hash) *vm.VerifiedCiphertext {
	ct, ok := interpreter.GetVerifiedCiphertexts()[ciphertextHash]
	if ok && isVerifiedAtCurrentDepth(ct) {
		return ct
	}
	return nil
}

func getVerifiedCiphertext(ciphertextHash tfhe.Hash) *vm.VerifiedCiphertext {
	return getVerifiedCiphertextFromEVM(ciphertextHash)
}

func get2VerifiedOperands(input []byte) (lhs *vm.VerifiedCiphertext, rhs *vm.VerifiedCiphertext, err error) {
	if len(input) != 65 {
		return nil, nil, errors.New("input needs to contain two 256-bit sized values and 1 8-bit value")
	}
	lhs = getVerifiedCiphertext(tfhe.BytesToHash(input[0:32]))
	if lhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	rhs = getVerifiedCiphertext(tfhe.BytesToHash(input[32:64]))
	if rhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func importCiphertextToEVMAtDepth(ct *tfhe.Ciphertext, depth int) *vm.VerifiedCiphertext {
	existing, ok := interpreter.GetVerifiedCiphertexts()[ct.Hash()]
	if ok {
		existing.GetVerifiedDepths().Add(depth)
		return existing
	} else {
		verifiedDepths := vm.NewDepthSet()
		verifiedDepths.Add(depth)
		new := vm.NewVerifiedCiphertext(verifiedDepths, ct)
		interpreter.GetVerifiedCiphertexts()[ct.Hash()] = new
		return new
	}
}

func importCiphertextToEVM(ct *tfhe.Ciphertext) *vm.VerifiedCiphertext {
	return importCiphertextToEVMAtDepth(ct, interpreter.GetEVM().GetDepth())
}

func importCiphertext(ct *tfhe.Ciphertext) *vm.VerifiedCiphertext {
	return importCiphertextToEVM(ct)
}

func Lior(a uint32, b uint32) (uint32, error) {
	return a * b, nil
}

func Moshe(input []byte, inputLen uint32) ([1][32]byte, error) {
	// byteSlice := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	// // Convert []byte to [1][32]byte
	// var byteArray [1][32]byte

	// copy(byteArray[0][:], byteSlice)

	// return byteArray, nil

	if interpreter == nil {
		msg := "fheAdd no evm interpreter"
		// logger.Error(msg, "lhs", lhs.ciphertext.UintType, "rhs", rhs.ciphertext.UintType)
		return [1][32]byte{}, errors.New(msg)
	}

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
	lhs, rhs, err := get2VerifiedOperands(input)
	if err != nil {
		logger.Error("fheAdd inputs not verified", "err", err, "input", hex.EncodeToString(input))
		return [1][32]byte{}, err
	}

	if lhs.GetCiphertext().UintType != rhs.GetCiphertext().UintType {
		msg := "fheAdd operand type mismatch"
		logger.Error(msg, "lhs", lhs.GetCiphertext().UintType, "rhs", rhs.GetCiphertext().UintType)
		return [1][32]byte{}, errors.New(msg)
	}

	// If we are doing gas estimation, skip execution and insert a random ciphertext as a result.
	// if !accessibleState.Interpreter().evm.Commit && !accessibleState.Interpreter().evm.EthCall {
	// 	return importRandomCiphertext(accessibleState, lhs.ciphertext.UintType), nil
	// }

	result, err := lhs.GetCiphertext().Add(rhs.GetCiphertext())
	if err != nil {
		logger.Error("fheAdd failed", "err", err)
		return [1][32]byte{}, err
	}
	importCiphertext(result)

	// TODO: for testing
	err = os.WriteFile("/tmp/add_result", result.Serialization, 0644)
	if err != nil {
		logger.Error("fheAdd failed to write /tmp/add_result", "err", err)
		return [1][32]byte{}, err
	}

	resultHash := result.Hash()
	logger.Info("fheAdd success", "lhs", lhs.GetCiphertext().Hash().Hex(), "rhs", rhs.GetCiphertext().Hash().Hex(), "result", resultHash.Hex())
	var byteArray [1][32]byte

	copy(byteArray[0][:], resultHash[:])
	return byteArray, nil
}
