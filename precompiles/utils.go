package precompiles

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	tfhe "github.com/fhenixprotocol/go-tfhe"
	"golang.org/x/crypto/nacl/box"
)

type TxParams struct {
	Commit        bool
	GasEstimation bool
	EthCall       bool
}

func TxParamsFromEVM(evm *vm.EVM) TxParams {
	var tp TxParams
	tp.Commit = evm.Commit
	tp.GasEstimation = evm.GasEstimation
	tp.EthCall = evm.EthCall

	return tp
}

type Precompile struct {
	Metadata *bind.MetaData
	Address  common.Address
}

type EthEncryptedReturn struct {
	Version        string `json:"version"`
	Nonce          string `json:"nonce"`
	EphemPublicKey string `json:"ephem_public_key"`
	Ciphertext     string `json:"ciphertext"`
}

func encryptForUser(value *big.Int, userPublicKey []byte) ([]byte, error) {

	ephemeralPub, ephemeralPriv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 24)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	encrypted := box.Seal(nil, value.Bytes(), (*[24]byte)(nonce), (*[32]byte)(userPublicKey), ephemeralPriv)
	////encrypted, err := box.SealAnonymous(nil, value.Bytes(), (*[32]byte)(userPublicKey), rand.Reader)
	//if err != nil {
	//	return nil, err
	//}

	encryptedReturnValue := EthEncryptedReturn{
		Version:        "x25519-xsalsa20-poly1305",
		Nonce:          base64.StdEncoding.EncodeToString(nonce),
		EphemPublicKey: base64.StdEncoding.EncodeToString(ephemeralPub[:]),
		Ciphertext:     base64.StdEncoding.EncodeToString(encrypted),
	}

	return json.Marshal(&encryptedReturnValue)
}

func encryptToUserKey(value *big.Int, pubKey []byte) ([]byte, error) {
	ct, err := encryptForUser(value, pubKey)
	if err != nil {
		return nil, err
	}

	// TODO: for testing
	err = os.WriteFile("/tmp/public_encrypt_result", ct, 0o644)
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func getCiphertext(ciphertextHash tfhe.Hash) *tfhe.Ciphertext {
	ctHashMapLock.RLock()
	defer ctHashMapLock.RUnlock()

	ct, ok := ctHashMap[ciphertextHash]
	if ok {
		return ct
	}
	return nil
}

func get2VerifiedOperands(input []byte) (lhs *tfhe.Ciphertext, rhs *tfhe.Ciphertext, err error) {
	if len(input) != 64 {
		return nil, nil, errors.New("input needs to contain two 256-bit sized values")
	}
	lhs = getCiphertext(tfhe.BytesToHash(input[0:32]))
	if lhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	rhs = getCiphertext(tfhe.BytesToHash(input[32:64]))
	if rhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func get3VerifiedOperands(input []byte) (control *tfhe.Ciphertext, ifTrue *tfhe.Ciphertext, ifFalse *tfhe.Ciphertext, err error) {
	if len(input) != 96 {
		return nil, nil, nil, errors.New("input needs to contain three 256-bit sized values and 1 8-bit value")
	}
	control = getCiphertext(tfhe.BytesToHash(input[0:32]))
	if control == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifTrue = getCiphertext(tfhe.BytesToHash(input[32:64]))
	if ifTrue == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifFalse = getCiphertext(tfhe.BytesToHash(input[64:96]))
	if ifFalse == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func importCiphertext(ct *tfhe.Ciphertext) *tfhe.Ciphertext {
	ctHashMapLock.Lock()
	defer ctHashMapLock.Unlock()

	existing, ok := ctHashMap[ct.Hash()]
	if ok {
		return existing
	} else {
		ctHashMap[ct.Hash()] = ct
		return ct
	}
}

func importRandomCiphertext(t tfhe.UintType) ([]byte, error) {
	ct, err := tfhe.NewRandomCipherText(t)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed creating random ciphertext of size: %d", t))
	}

	importCiphertext(ct)
	ctHash := ct.Hash()
	return ctHash[:], nil
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// Puts the given ciphertext as a require to the oracle DB or exits the process on errors.
// Returns the require value.
//func putRequire(ct *tfhe.Ciphertext, interpreter *vm.EVMInterpreter) (bool, error) {
//	plaintext, err := tfhe.Decrypt(*ct)
//	if err != nil {
//		return false, errors.New(fmt.Sprintf("Failed to decrypt value: %s", err))
//	}
//
//	result, err := tfhe.StoreRequire(ct, plaintext)
//	if err != nil {
//		return false, errors.New("Failed to store require in DB")
//	}
//
//	return result, nil
//}
//
//// Gets the given require from the oracle DB and returns its value.
//// Exits the process on errors or signature verification failure.
//func getRequire(ct *tfhe.Ciphertext) (bool, error) {
//	result, err := tfhe.CheckRequire(ct)
//	if err != nil {
//		return false, errors.New(fmt.Sprintf("Error verifying require", err))
//	}
//
//	return result, nil
//}

func evaluateRequire(ct *tfhe.Ciphertext) bool {
	return tfhe.Require(ct)
}

type fheUintType uint8

const (
	FheUint8  fheUintType = 0
	FheUint16 fheUintType = 1
	FheUint32 fheUintType = 2
)

func isValidType(t byte) bool {
	if uint8(t) < uint8(FheUint8) || uint8(t) > uint8(FheUint32) {
		return false
	}
	return true
}
