package types

import (
	"fmt"
	"math/big"

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
	// Rol  // Commented out if not used
	// Ror  // Commented out if not used
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
	TrivialEncrypt: "trivialEncrypt",
	// Rol:          "rol",
	// Ror:          "ror",
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
	"trivialEncrypt": TrivialEncrypt,
	// "rol":          Rol,
	// "ror":          Ror,
}

func (pn PrecompileName) String() string {
	return precompileNameToString[pn]
}

func PrecompileNameFromString(s string) (PrecompileName, bool) {
	pn, ok := stringToPrecompileName[s]
	return pn, ok
}

type PendingDecryption struct {
	Hash fhe.Hash
	Type PrecompileName
}

type DecryptionResults struct {
	data map[PendingDecryption]interface{}
}

func NewDecryptionResultsMap() *DecryptionResults {
	return &DecryptionResults{
		data: make(map[PendingDecryption]interface{}),
	}
}

func (dr *DecryptionResults) CreateEmptyRecord(key PendingDecryption) {
	dr.data[key] = nil
}

func (dr *DecryptionResults) SetValue(key PendingDecryption, value interface{}) error {
	if _, exists := dr.data[key]; !exists {
		return fmt.Errorf("record does not exist, create an empty record first")
	}

	switch key.Type {
	case SealOutput:
		if _, ok := value.([]byte); !ok {
			return fmt.Errorf("value for SealOutput must be []byte")
		}
	case Require:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("value for Require must be bool")
		}
	case Decrypt:
		if _, ok := value.(*big.Int); !ok {
			return fmt.Errorf("value for Decrypt must be *big.Int")
		}
	default:
		return fmt.Errorf("unknown PrecompileName")
	}

	dr.data[key] = value
	return nil
}

func (dr *DecryptionResults) Get(key PendingDecryption) (interface{}, error) {
	value, exists := dr.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found")
	}

	if value == nil {
		return nil, nil // Exists but no value
	}

	switch key.Type {
	case SealOutput:
		if bytes, ok := value.([]byte); ok {
			return bytes, nil
		}
		return nil, fmt.Errorf("value is not []byte as expected for SealOutput")
	case Require:
		if boolValue, ok := value.(bool); ok {
			return boolValue, nil
		}
		return nil, fmt.Errorf("value is not bool as expected for Require")
	case Decrypt:
		if bigInt, ok := value.(*big.Int); ok {
			return bigInt, nil
		}
		return nil, fmt.Errorf("value is not *big.Int as expected for Decrypt")
	default:
		return nil, fmt.Errorf("unknown PrecompileName")
	}
}
