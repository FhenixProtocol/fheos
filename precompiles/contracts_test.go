package precompiles

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	fhedriver "github.com/fhenixprotocol/warp-drive/fhe-driver"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"strconv"
	"testing"
)

var tp TxParams

func getenvInt(key string, defaultValue int) (int, error) {
	s := os.Getenv(key)
	if s == "" {
		return defaultValue, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func init() {
	if os.Getenv("FHEOS_DB_PATH") == "" {
		if err := os.Setenv("FHEOS_DB_PATH", "./fheosdb"); err != nil {
			panic(err)
		}
	}

	if err := InitFheConfig(&fhedriver.ConfigDefault); err != nil {
		panic(err)
	}

	if err := InitializeFheosState(); err != nil {
		panic(err)
	}

	if err := os.Setenv("FHEOS_DB_PATH", ""); err != nil {
		panic(err)
	}

	tp = TxParams{
		Commit:          true,
		GasEstimation:   false,
		EthCall:         false,
		CiphertextDb:    memorydb.New(),
		ContractAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		GetBlockHash:    vm.GetHashFunc(nil),
		BlockNumber:     nil,
		ParallelTxHooks: nil,
	}
}

func trivialEncrypt(t *testing.T, number *big.Int, uintType uint8, securityZone int32) []byte {
	ct, _, err := TrivialEncrypt(number.Bytes(), uintType, securityZone, &tp, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, ct, nil)

	ctKey, err := fhedriver.DeserializeCiphertextKey(ct)
	assert.NoError(t, err)
	assert.True(t, ctKey.IsTriviallyEncrypted)
	assert.Equal(t, ctKey.UintType, fhedriver.EncryptionType(uintType))
	assert.Equal(t, ctKey.SecurityZone, securityZone)

	return ct
}

func forEveryUintType(t *testing.T, testName string, f func(t *testing.T, uintType uint8)) {
	for _, uintType := range []fhedriver.EncryptionType{fhedriver.Uint8, fhedriver.Uint16, fhedriver.Uint32, fhedriver.Uint64, fhedriver.Uint128, fhedriver.Uint256} {
		t.Run(fmt.Sprintf("Running %s test with %s", testName, uintType.ToString()), func(t *testing.T) {
			f(t, uint8(uintType))
		})
	}
}

func forEveryUintTypeAndBool(t *testing.T, testName string, f func(t *testing.T, uintType uint8)) {
	forEveryUintType(t, testName, f)
	t.Run(fmt.Sprintf("Running %s test with %s", testName, fhedriver.Bool.ToString()), func(t *testing.T) {
		f(t, uint8(fhedriver.Bool))
	})
}

func forEveryEncryptedType(t *testing.T, testName string, f func(t *testing.T, uintType uint8)) {
	forEveryUintType(t, testName, f)
	for _, ty := range []fhedriver.EncryptionType{fhedriver.Address, fhedriver.Bool} {
		t.Run(fmt.Sprintf("Running %s test with %s", testName, ty.ToString()), func(t *testing.T) {
			f(t, uint8(ty))
		})
	}
}

func expectPlaintext(t *testing.T, ct []byte, uintType uint8, expected *big.Int) {
	plaintext, _, err := Decrypt(uintType, ct, nil, &tp, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, plaintext, nil)
	assert.Equal(t, plaintext, expected)
}

func generalTwoOpTest(t *testing.T, lhs, rhs *big.Int, uintType uint8, plaintextFunc func(*big.Int, *big.Int) *big.Int, encryptedFunc func(byte, []byte, []byte, *TxParams, *CallbackFunc) ([]byte, uint64, error)) {
	ctLhs := trivialEncrypt(t, lhs, uintType, 0)
	ctRhs := trivialEncrypt(t, rhs, uintType, 0)

	plaintextResult := plaintextFunc(lhs, rhs)
	ctResult, _, err := encryptedFunc(uintType, ctLhs, ctRhs, &tp, nil)
	assert.NoError(t, err)
	expectPlaintext(t, ctResult, uintType, plaintextResult)
}

func generalOneOpTest(t *testing.T, val *big.Int, uintType uint8, plaintextFunc func(*big.Int) *big.Int, encryptedFunc func(byte, []byte, *TxParams, *CallbackFunc) ([]byte, uint64, error)) {
	ctVal := trivialEncrypt(t, val, uintType, 0)

	plaintextResult := plaintextFunc(val)
	ctResult, _, err := encryptedFunc(uintType, ctVal, &tp, nil)
	assert.NoError(t, err)
	expectPlaintext(t, ctResult, uintType, plaintextResult)
}

func TestTrivialEncrypt(t *testing.T) {
	forEveryEncryptedType(t, "TrivialEncrypt", func(t *testing.T, uintType uint8) {
		ct := trivialEncrypt(t, big.NewInt(1), uintType, 0)
		expectPlaintext(t, ct, uintType, big.NewInt(1))
	})
}
func TestAdd(t *testing.T) {
	lhsVal := big.NewInt(120)
	rhsVal := big.NewInt(2)
	forEveryUintType(t, "Add", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return new(big.Int).Add(lhs, rhs)
		}, Add)
	})
}

func TestLte(t *testing.T) {
	lhsVal := big.NewInt(120)
	rhsVal := big.NewInt(2)

	trueVal := big.NewInt(1)
	falseVal := big.NewInt(0)
	forEveryUintType(t, "Lte 1", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if lhs.Cmp(rhs) <= 0 {
				return trueVal
			}

			return falseVal
		}, Lte)
	})
	forEveryUintType(t, "Lte 2", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, rhsVal, lhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if lhs.Cmp(rhs) <= 0 {
				return trueVal
			}
			return falseVal
		}, Lte)
	})
	forEveryUintType(t, "Lte 3", func(t *testing.T, uintType uint8) {

		generalTwoOpTest(t, lhsVal, lhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return trueVal
		}, Lte)
	})
}

func TestSub(t *testing.T) {
	lhsVal := big.NewInt(120)
	rhsVal := big.NewInt(2)
	forEveryUintType(t, "Sub", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return new(big.Int).Sub(lhs, rhs)
		}, Sub)
	})
}

func TestMul(t *testing.T) {
	lhsVal := big.NewInt(120)
	rhsVal := big.NewInt(2)
	forEveryUintType(t, "Mul", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return new(big.Int).Mul(lhs, rhs)
		}, Mul)
	})
}

func TestLt(t *testing.T) {
	lhsVal := big.NewInt(120)
	rhsVal := big.NewInt(2)

	trueVal := big.NewInt(1)
	falseVal := big.NewInt(0)
	forEveryUintType(t, "Lt 1", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if lhs.Cmp(rhs) < 0 {
				return trueVal
			}

			return falseVal
		}, Lt)
	})
	forEveryUintType(t, "Lt 2", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, rhsVal, lhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if lhs.Cmp(rhs) < 0 {
				return trueVal
			}
			return falseVal
		}, Lt)
	})
	forEveryUintType(t, "Lt 3", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, lhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return falseVal
		}, Lt)
	})
}

func TestDiv(t *testing.T) {
	lhsVal := big.NewInt(120)
	rhsVal := big.NewInt(2)
	forEveryUintType(t, "Div", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return new(big.Int).Div(lhs, rhs)
		}, Div)
	})
}

func TestGt(t *testing.T) {
	lhsVal := big.NewInt(120)
	rhsVal := big.NewInt(2)

	trueVal := big.NewInt(1)
	falseVal := big.NewInt(0)
	forEveryUintType(t, "Gt 1", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if lhs.Cmp(rhs) > 0 {
				return trueVal
			}

			return falseVal
		}, Gt)
	})
	forEveryUintType(t, "Gt 2", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, rhsVal, lhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if lhs.Cmp(rhs) > 0 {
				return trueVal
			}
			return falseVal
		}, Gt)
	})
	forEveryUintType(t, "Gt 3", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, lhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return falseVal
		}, Gt)
	})
}

func TestGte(t *testing.T) {
	lhsVal := big.NewInt(120)
	rhsVal := big.NewInt(2)

	trueVal := big.NewInt(1)
	falseVal := big.NewInt(0)
	forEveryUintType(t, "Gte 1", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if lhs.Cmp(rhs) >= 0 {
				return trueVal
			}

			return falseVal
		}, Gte)
	})
	forEveryUintType(t, "Gte 2", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, rhsVal, lhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if lhs.Cmp(rhs) >= 0 {
				return trueVal
			}
			return falseVal
		}, Gte)
	})
	forEveryUintType(t, "Gte 3", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, lhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return trueVal
		}, Gte)
	})
}

func TestRem(t *testing.T) {
	lhsVal := big.NewInt(120)
	rhsVal := big.NewInt(9)
	forEveryUintType(t, "Rem", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return new(big.Int).Rem(lhs, rhs)
		}, Rem)
	})
}

func TestAnd(t *testing.T) {
	lhsVal := big.NewInt(66)
	rhsVal := big.NewInt(2)

	trueVal := big.NewInt(1)
	forEveryUintTypeAndBool(t, "And", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if uintType == uint8(fhedriver.Bool) {
				return trueVal
			}

			return new(big.Int).And(lhs, rhs)
		}, And)
	})
}

func TestOr(t *testing.T) {
	lhsVal := big.NewInt(32)
	rhsVal := big.NewInt(2)

	forEveryUintTypeAndBool(t, "Or", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if uintType == uint8(fhedriver.Bool) {
				return big.NewInt(1)
			}

			return new(big.Int).Or(lhs, rhs)
		}, Or)
	})
}

func TestXor(t *testing.T) {
	lhsVal := big.NewInt(32)
	rhsVal := big.NewInt(2)

	forEveryUintTypeAndBool(t, "Xor", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if uintType == uint8(fhedriver.Bool) {
				return big.NewInt(0)
			}

			return new(big.Int).Xor(lhs, rhs)
		}, Xor)
	})
}

func TestEq(t *testing.T) {
	lhsVal := big.NewInt(32)
	rhsVal := big.NewInt(0)

	trueVal := big.NewInt(1)
	falseVal := big.NewInt(0)
	forEveryEncryptedType(t, "Eq 1", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return falseVal
		}, Eq)
	})
	forEveryEncryptedType(t, "Eq 2", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, rhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return trueVal
		}, Eq)
	})
}

func TestNe(t *testing.T) {
	lhsVal := big.NewInt(32)
	rhsVal := big.NewInt(0)

	trueVal := big.NewInt(1)
	falseVal := big.NewInt(0)
	forEveryEncryptedType(t, "Ne", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return trueVal
		}, Ne)
	})
	forEveryEncryptedType(t, "Ne", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, rhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return falseVal
		}, Ne)
	})
}

func TestMin(t *testing.T) {
	lhsVal := big.NewInt(32)
	rhsVal := big.NewInt(2)

	forEveryUintType(t, "Min", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if lhs.Cmp(rhs) < 0 {
				return lhsVal
			}

			return rhsVal
		}, Min)
	})
}

func TestMax(t *testing.T) {
	lhsVal := big.NewInt(32)
	rhsVal := big.NewInt(2)

	forEveryUintType(t, "Max", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			if lhs.Cmp(rhs) > 0 {
				return lhsVal
			}

			return rhsVal
		}, Max)
	})
}

func TestShl(t *testing.T) {
	lhsVal := big.NewInt(2)
	rhsVal := big.NewInt(2)
	forEveryUintType(t, "Shl", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return new(big.Int).Lsh(lhs, uint(rhs.Uint64()))
		}, Shl)
	})
}

func TestShr(t *testing.T) {
	lhsVal := big.NewInt(16)
	rhsVal := big.NewInt(2)
	forEveryUintType(t, "Shr", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			return new(big.Int).Rsh(lhs, uint(rhs.Uint64()))
		}, Shr)
	})
}

func TestRol(t *testing.T) {
	lhsVal := big.NewInt(2)
	rhsVal := big.NewInt(2)
	forEveryUintType(t, "Rol", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			// with the number we are working with this is the same as Rol
			return new(big.Int).Lsh(lhs, uint(rhs.Uint64()))
		}, Rol)
	})
}
func TestRor(t *testing.T) {
	lhsVal := big.NewInt(16)
	rhsVal := big.NewInt(2)
	forEveryUintType(t, "Ror", func(t *testing.T, uintType uint8) {
		generalTwoOpTest(t, lhsVal, rhsVal, uintType, func(lhs, rhs *big.Int) *big.Int {
			// with the number we are working with this is the same as Rol
			return new(big.Int).Rsh(lhs, uint(rhs.Uint64()))
		}, Ror)
	})
}

func maxBigInt(bits int) *big.Int {
	two := big.NewInt(2)
	exp := new(big.Int).Exp(two, big.NewInt(int64(bits)), nil)

	maxValue := new(big.Int).Sub(exp, big.NewInt(1))

	return maxValue
}

func TestNot(t *testing.T) {
	val := big.NewInt(16)
	forEveryUintTypeAndBool(t, "Not", func(t *testing.T, uintType uint8) {
		generalOneOpTest(t, val, uintType, func(val *big.Int) *big.Int {
			if uintType == uint8(fhedriver.Bool) {
				return big.NewInt(0)
			}
			return new(big.Int).Sub(maxBigInt(1<<(3+uintType)), val)
		}, Not)
	})
}

func TestSelect(t *testing.T) {
	trueCt := trivialEncrypt(t, big.NewInt(1), uint8(fhedriver.Bool), 0)
	falseCt := trivialEncrypt(t, big.NewInt(0), uint8(fhedriver.Bool), 0)

	ifTrueVal := big.NewInt(100)
	ifFalseVal := big.NewInt(0)

	forEveryUintTypeAndBool(t, "Select", func(t *testing.T, uintType uint8) {
		ifTrue := trivialEncrypt(t, ifTrueVal, uintType, 0)
		ifFalse := trivialEncrypt(t, ifFalseVal, uintType, 0)

		plaintextResult := ifTrueVal
		if uintType == uint8(fhedriver.Bool) {
			plaintextResult = big.NewInt(1)
		}
		ctResult, _, err := Select(uintType, trueCt, ifTrue, ifFalse, &tp, nil)
		assert.NoError(t, err)
		expectPlaintext(t, ctResult, uintType, plaintextResult)

		plaintextResult = ifFalseVal
		if uintType == uint8(fhedriver.Bool) {
			plaintextResult = big.NewInt(0)
		}
		ctResult, _, err = Select(uintType, falseCt, ifTrue, ifFalse, &tp, nil)
		assert.NoError(t, err)
		expectPlaintext(t, ctResult, uintType, plaintextResult)
	})
}
