package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

type DataType uint64

type Hash fhe.Hash
type FheEncrypted fhe.FheEncrypted

func IsValidType(t fhe.EncryptionType) bool {
	return t >= fhe.Uint8 && t <= fhe.Bool
}

type CipherTextRepresentation struct {
	Data     *FheEncrypted
	Owners   []common.Address
	RefCount uint64
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
	PutCt(h Hash, cipher *CipherTextRepresentation) error
	GetCt(h Hash) (*CipherTextRepresentation, error)

	HasCt(h Hash) bool

	DeleteCt(h Hash) error
}

type PrecompileName int

const (
	GetNetworkKey PrecompileName = iota
	Verify
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
	Verify:         "verify",
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
	"verify":         Verify,
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
