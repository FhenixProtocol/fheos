package precompiles

import (
	"bytes"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/fheos/storage"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
	"golang.org/x/crypto/sha3"
	"hash"
	"math"
	"math/big"
)

type TxParams struct {
	Commit          bool
	GasEstimation   bool
	EthCall         bool
	CiphertextDb    *memorydb.Database
	ContractAddress common.Address
	TxHash          common.Hash
}

type GasBurner interface {
	Burn(amount uint64) error
	Burned() uint64
}

func TxParamsFromEVM(evm *vm.EVM, callerContract common.Address) TxParams {
	var tp TxParams
	tp.Commit = evm.Commit
	tp.GasEstimation = evm.GasEstimation
	tp.EthCall = evm.EthCall

	tp.CiphertextDb = evm.CiphertextDb
	tp.ContractAddress = callerContract

	tp.TxHash = evm.TransactionHash

	return tp
}

type Precompile struct {
	Metadata *bind.MetaData
	Address  common.Address
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := NewKeccakState()
	for _, datum := range data {
		d.Write(datum)
	}
	return d.Sum(nil)
}

// )(ITZIK)(
var tempValueHashMagicBytes = [...]byte{0x13, 0x37, 0xb0, 0x0b}

func NewKeccakState() hash.Hash {
	return sha3.NewLegacyKeccak256()
}
func CalcTempValueHash(randomBytes []byte) []byte {
	_hash := Keccak256(randomBytes[:])
	copy(_hash[:4], tempValueHashMagicBytes[:])

	return _hash
}
func CheckIfLeetBoob(hashToCheck []byte) bool {
	return bytes.Equal(hashToCheck[:4], tempValueHashMagicBytes[:])
}

func CreateFakeFheEncrypted(bytes []byte) *fhe.FheEncrypted {
	//bytes := [...]byte{0x05}
	var result fhe.FheEncrypted
	if bytes == nil {
		result = fhe.NewFheEncryptedFromBytes(nil, 5, false, false)
		return &result
	}
	result = fhe.NewFheEncryptedFromBytes(bytes[:], 5, false, false)
	return &result
}

func getCiphertext(state *storage.MultiStore, ciphertextHash fhe.Hash, caller common.Address) *fhe.FheEncrypted {
	ct, err := state.GetCt(types.Hash(ciphertextHash), caller)
	if err != nil {
		logger.Error("reading ciphertext from State resulted with error", "error", err.Error())
		return nil
	}

	return (*fhe.FheEncrypted)(ct)
}
func get2VerifiedOperands(storage *storage.MultiStore, lhsHash []byte, rhsHash []byte, caller common.Address) (lhs *fhe.FheEncrypted, rhs *fhe.FheEncrypted, err error) {
	if len(lhsHash) != 32 || len(rhsHash) != 32 {
		return nil, nil, errors.New("ciphertext's hashes need to be 32 bytes long")
	}

	lhs = getCiphertext(storage, fhe.BytesToHash(lhsHash), caller)
	if lhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	rhs = getCiphertext(storage, fhe.BytesToHash(rhsHash), caller)
	if rhs == nil {
		return nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

func get3VerifiedOperands(storage *storage.MultiStore, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte, caller common.Address) (control *fhe.FheEncrypted, ifTrue *fhe.FheEncrypted, ifFalse *fhe.FheEncrypted, err error) {
	if len(controlHash) != 32 || len(ifTrueHash) != 32 || len(ifFalseHash) != 32 {
		return nil, nil, nil, errors.New("ciphertext's hashes need to be 32 bytes long")
	}

	control = getCiphertext(storage, fhe.BytesToHash(controlHash), caller)
	if control == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifTrue = getCiphertext(storage, fhe.BytesToHash(ifTrueHash), caller)
	if ifTrue == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	ifFalse = getCiphertext(storage, fhe.BytesToHash(ifFalseHash), caller)
	if ifFalse == nil {
		return nil, nil, nil, errors.New("unverified ciphertext handle")
	}
	err = nil
	return
}

// )(ITZIK)(
func AllZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

// )(ITZIK)(
func storeTempValue(storage *storage.MultiStore, key types.Hash, ct *fhe.FheEncrypted, owner common.Address) error {
	err := storage.HackAppendCt(key, (*types.FheEncrypted)(ct), owner)
	if err != nil {
		logger.Error("failed importing ciphertext to state: ", err)
		return err
	}

	return nil
}
func storeCipherText(storage *storage.MultiStore, ct *fhe.FheEncrypted, owner common.Address) error {
	err := storage.AppendCt(types.Hash(ct.Hash()), (*types.FheEncrypted)(ct), owner)
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

func evaluateRequire(ct *fhe.FheEncrypted) bool {
	return fhe.Require(ct)
}

func FakeDecryptionResult(encType fhe.EncryptionType) *big.Int {

	switch encType {
	case fhe.Uint8:
		return big.NewInt(math.MaxUint8 / 2)
	case fhe.Uint16:
		return big.NewInt(math.MaxInt16 / 2)
	case fhe.Uint32:
		return big.NewInt(math.MaxUint32 / 2)
	case fhe.Uint64:
		return big.NewInt(math.MaxUint64 / 2)
	case fhe.Uint128:
		value := &big.Int{}
		value.SetString("2222222222222222222222222222222", 16)
		return value
	case fhe.Uint256:
		value := &big.Int{}
		value.SetString("2222222222222222222222222222222222222222222222222", 16)
		return value
	case fhe.Address:
		value := &big.Int{}
		value.SetString("Dd4BEac65bad064932FB21aE7Ba2aa6e8fbc41A4", 16)
		return value
	default:
		return big.NewInt(0)
	}
}
