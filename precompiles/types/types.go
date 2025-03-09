package types

import (
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

type DataType uint64

type Hash fhe.Hash
type FheEncrypted fhe.FheEncrypted

var DeserializeCiphertextKey = fhe.DeserializeCiphertextKey
var GetEmptyCiphertextKey = fhe.GetEmptyCiphertextKey
var SerializeCiphertextKey = fhe.SerializeCiphertextKey

func IsValidType(t fhe.EncryptionType) bool {
	switch t {
	case fhe.Bool,
		fhe.Uint8, fhe.Uint16, fhe.Uint32, fhe.Uint64, fhe.Uint128, fhe.Uint256,
		fhe.Int8, fhe.Int16, fhe.Int32, fhe.Int64, fhe.Int128, fhe.Int256,
		fhe.Address:
		return true
	}
	return false
}

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

	DeleteCt(h Hash) error
}

type PrecompileName int

const (
	GetNetworkKey PrecompileName = iota
	StoreCt
	Cast
	SealOutput
	Select
	Require
	Decrypt
	Sub
	Add
	Xor
	And
	Or
	Not
	Div
	Rem
	Mul
	Shl
	Shr
	Gte
	Lte
	Lt
	Gt
	Min
	Max
	Eq
	Ne
	TrivialEncrypt
	Random
	Rol
	Ror
	Square
)

var precompileNameToString = map[PrecompileName]string{
	GetNetworkKey:  "getNetworkKey",
	StoreCt:        "verify",
	Cast:           "cast",
	SealOutput:     "sealOutput",
	Select:         "select",
	Require:        "require",
	Decrypt:        "decrypt",
	Sub:            "sub",
	Add:            "add",
	Xor:            "xor",
	And:            "and",
	Or:             "or",
	Not:            "not",
	Div:            "div",
	Rem:            "rem",
	Mul:            "mul",
	Shl:            "shl",
	Shr:            "shr",
	Gte:            "gte",
	Lte:            "lte",
	Lt:             "lt",
	Gt:             "gt",
	Min:            "min",
	Max:            "max",
	Eq:             "eq",
	Ne:             "ne",
	Random:         "random",
	TrivialEncrypt: "trivialEncrypt",
	Rol:            "rol",
	Ror:            "ror",
	Square:         "square",
}

var stringToPrecompileName = map[string]PrecompileName{
	"getNetworkKey":  GetNetworkKey,
	"storeCt":        StoreCt,
	"cast":           Cast,
	"sealOutput":     SealOutput,
	"select":         Select,
	"require":        Require,
	"decrypt":        Decrypt,
	"sub":            Sub,
	"add":            Add,
	"xor":            Xor,
	"and":            And,
	"or":             Or,
	"not":            Not,
	"div":            Div,
	"rem":            Rem,
	"mul":            Mul,
	"shl":            Shl,
	"shr":            Shr,
	"gte":            Gte,
	"lte":            Lte,
	"lt":             Lt,
	"gt":             Gt,
	"min":            Min,
	"max":            Max,
	"eq":             Eq,
	"ne":             Ne,
	"random":         Random,
	"trivialEncrypt": TrivialEncrypt,
	"rol":            Rol,
	"ror":            Ror,
	"square":         Square,
}

func (pn PrecompileName) String() string {
	return precompileNameToString[pn]
}

func PrecompileNameFromString(s string) (PrecompileName, bool) {
	pn, ok := stringToPrecompileName[s]
	return pn, ok
}

type ParallelTxProcessingHook interface {
	NotifyCt(*PendingDecryption)
	NotifyDecryptRes(*PendingDecryption) error
	NotifyExistingRes(*PendingDecryption)
}

const (
	TrivialEncryptAndTypeByte = 30
	SecurityZoneByte          = 31
	TypeMask                  = 0x7f
	TrivialEncryptFlag        = 0x80
)
