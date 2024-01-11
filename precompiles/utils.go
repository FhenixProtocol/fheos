package precompiles

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	tfhe "github.com/fhenixprotocol/go-tfhe"
	"golang.org/x/crypto/nacl/box"
	"math/big"
)

type TxParams struct {
	Commit        bool
	GasEstimation bool
	EthCall       bool
}

type GasBurner interface {
	Burn(amount uint64) error
	Burned() uint64
}

type DataType uint64

const (
	version DataType = iota
	ct
)

type Storage interface {
	Put(t DataType, key []byte, val []byte) error
	Get(t DataType, key []byte) ([]byte, error)
	GetVersion() (uint64, error)
	PutVersion(v uint64) error
	PutCt(h tfhe.Hash, cipher *tfhe.Ciphertext) error
	GetCt(h tfhe.Hash) (*tfhe.Ciphertext, error)
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
	EphemPublicKey string `json:"ephemPublicKey"`
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

	return ct, nil
}

func getCiphertext(state *FheosState, ciphertextHash tfhe.Hash) *tfhe.Ciphertext {
	ct, err := state.GetCiphertext(ciphertextHash)
	if err != nil {
		logger.Error("reading ciphertext from state resulted with error: ", err)
		return nil
	}

	return ct
}

func get2VerifiedOperands(state *FheosState, input []byte) (lhs *tfhe.Ciphertext, rhs *tfhe.Ciphertext, err error) {
	if len(input) != 64 {
		return nil, nil, errors.New("input needs to contain two 256-bit sized values")
	}
	lhs = getCiphertext(state, tfhe.BytesToHash(input[0:32]))
	if lhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	rhs = getCiphertext(state, tfhe.BytesToHash(input[32:64]))
	if rhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func get3VerifiedOperands(state *FheosState, input []byte) (control *tfhe.Ciphertext, ifTrue *tfhe.Ciphertext, ifFalse *tfhe.Ciphertext, err error) {
	if len(input) != 96 {
		return nil, nil, nil, errors.New("input needs to contain three 256-bit sized values and 1 8-bit value")
	}
	control = getCiphertext(state, tfhe.BytesToHash(input[0:32]))
	if control == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifTrue = getCiphertext(state, tfhe.BytesToHash(input[32:64]))
	if ifTrue == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifFalse = getCiphertext(state, tfhe.BytesToHash(input[64:96]))
	if ifFalse == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func importCiphertext(state *FheosState, ct *tfhe.Ciphertext) *tfhe.Ciphertext {
	err := state.SetCiphertext(ct)
	if err != nil {
		logger.Error("failed importing ciphertext to state: ", err)
		return err
	}

	return nil
}

func importRandomCiphertext(state *FheosState, t tfhe.UintType) ([]byte, error) {
	ct, err := tfhe.NewRandomCipherText(t)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed creating random ciphertext of size: %d", t))
	}

	err = importCiphertext(state, ct)
	if err != nil {
		return nil, err
	}

	ctHash := ct.Hash()
	return ctHash[:], nil
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

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
