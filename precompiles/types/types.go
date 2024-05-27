package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

type DataType uint64

type Hash fhe.Hash
type FheEncrypted fhe.FheEncrypted

type CipherTextRepresentation struct {
	Data   *FheEncrypted
	Owners []common.Address
}

// SharedCiphertext is a struct that represents a ciphertext with a ref count
type SharedCiphertext struct {
	Ciphertext CipherTextRepresentation
	RefCount   uint64
}

type Storage interface {
	// don't really need these
	// Put(t types.DataType, key []byte, val []byte) error
	// Get(t types.DataType, key []byte) ([]byte, error)
	GetVersion() (uint64, error)
	PutVersion(v uint64) error
	FheInternalCipherTextStorage
}

type FheCipherTextStorage interface {
	PutCt(h Hash, cipher *CipherTextRepresentation) error
	GetCt(h Hash) (*CipherTextRepresentation, error)

	HasCt(h Hash) bool
}

type FheInternalCipherTextStorage interface {
	PutCt(h Hash, cipher *SharedCiphertext) error
	GetCt(h Hash) (*SharedCiphertext, error)

	HasCt(h Hash) bool

	DeleteCt(h Hash) error
}
