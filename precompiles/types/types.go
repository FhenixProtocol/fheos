package types

import "github.com/fhenixprotocol/warp-drive/fhe-driver"

type DataType uint64

type Hash fhe.Hash
type FheEncrypted fhe.FheEncrypted

type Storage interface {
	// don't really need these
	// Put(t types.DataType, key []byte, val []byte) error
	// Get(t types.DataType, key []byte) ([]byte, error)
	GetVersion() (uint64, error)
	PutVersion(v uint64) error
	FheCipherTextStorage
}

type FheCipherTextStorage interface {
	PutCt(h Hash, cipher *FheEncrypted) error
	GetCt(h Hash) (*FheEncrypted, error)

	HasCt(h Hash) bool
}
