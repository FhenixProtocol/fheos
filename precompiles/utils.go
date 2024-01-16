package precompiles

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
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
func get2VerifiedOperands(state *FheosState, lhsHash []byte, rhsHash []byte) (lhs *tfhe.Ciphertext, rhs *tfhe.Ciphertext, err error) {
	if len(lhsHash) != 32 || len(rhsHash) != 32 {
		return nil, nil, errors.New("ciphertext's hashes need to be 32 bytes long")
	}

	lhs = getCiphertext(state, tfhe.BytesToHash(lhsHash))
	if lhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	rhs = getCiphertext(state, tfhe.BytesToHash(rhsHash))
	if rhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func get3VerifiedOperands(state *FheosState, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte) (control *tfhe.Ciphertext, ifTrue *tfhe.Ciphertext, ifFalse *tfhe.Ciphertext, err error) {
	if len(controlHash) != 32 || len(ifTrueHash) != 32 || len(ifFalseHash) != 32 {
		return nil, nil, nil, errors.New("ciphertext's hashes need to be 32 bytes long")
	}

	control = getCiphertext(state, tfhe.BytesToHash(controlHash))
	if control == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifTrue = getCiphertext(state, tfhe.BytesToHash(ifTrueHash))
	if ifTrue == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifFalse = getCiphertext(state, tfhe.BytesToHash(ifFalseHash))
	if ifFalse == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func importCiphertext(state *FheosState, ct *tfhe.Ciphertext) error {
	err := state.SetCiphertext(ct)
	if err != nil {
		logger.Error("failed importing ciphertext to state: ", err)
		return err
	}

	return nil
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
	return t >= uint8(FheUint8) && t <= uint8(FheUint32)
}
