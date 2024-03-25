package precompiles

import (
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

// we calculated gas using the results of the following test results: https://github.com/FhenixProtocol/fheos/actions/runs/7475858194/job/20344926903
// we considered the time that it took to run division of two euint32s as a pivot that will take 1M gas and calculated how much gas is if for 1ms of runtime
// 1M / 4800ms = 209 gas per 1ms

func getGasForPrecompile(precompileName string, uintType fhe.EncryptionType) uint64 {
	return getRawPrecompileGas(precompileName, uintType) * 209
}

func getRawPrecompileGas(precompileName string, uintType fhe.EncryptionType) uint64 {
	switch precompileName {
	case "add":
		switch uintType {
		case fhe.Uint8:
			return 35
		case fhe.Uint16:
			return 73
		case fhe.Uint32:
			return 163
		}
	case "verify":
		switch uintType {
		case fhe.Uint8:
			return 44
		case fhe.Uint16:
			return 44
		case fhe.Uint32:
			return 44
		}
	case "sealOutput":
		switch uintType {
		case fhe.Uint8:
			return 23
		case fhe.Uint16:
			return 23
		case fhe.Uint32:
			return 23
		}
	case "decrypt": // No test, roughly estimated as sealOutput - trivialEncrypt
		switch uintType {
		case fhe.Uint8:
			return 18
		case fhe.Uint16:
			return 18
		case fhe.Uint32:
			return 18
		}
	case "lte":
		switch uintType {
		case fhe.Uint8:
			return 18
		case fhe.Uint16:
			return 29
		case fhe.Uint32:
			return 43
		}
	case "sub":
		switch uintType {
		case fhe.Uint8:
			return 35
		case fhe.Uint16:
			return 73
		case fhe.Uint32:
			return 163
		}
	case "mul":
		switch uintType {
		case fhe.Uint8:
			return 94
		case fhe.Uint16:
			return 311
		case fhe.Uint32:
			return 1127
		}
	case "lt":
		switch uintType {
		case fhe.Uint8:
			return 22
		case fhe.Uint16:
			return 36
		case fhe.Uint32:
			return 60
		}
	case "select":
		switch uintType {
		case fhe.Uint8:
			return 214
		case fhe.Uint16:
			return 316
		case fhe.Uint32:
			return 549
		}
	case "require": // Took the values when there was no crash as for crash gas is irrelevant as it will be reverted
		switch uintType {
		case fhe.Uint8:
			return 65
		case fhe.Uint16:
			return 65
		case fhe.Uint32:
			return 65
		}
	case "trivialEncrypt":
		switch uintType {
		case fhe.Uint8:
			return 5
		case fhe.Uint16:
			return 5
		case fhe.Uint32:
			return 5
		}
	case "cast":
		switch uintType {
		case fhe.Uint8:
			return 5
		case fhe.Uint16:
			return 5
		case fhe.Uint32:
			return 5
		}
	case "div":
		switch uintType {
		case fhe.Uint8:
			return 447
		case fhe.Uint16:
			return 1310
		case fhe.Uint32:
			return 4800
		}
	case "gt":
		switch uintType {
		case fhe.Uint8:
			return 21
		case fhe.Uint16:
			return 29
		case fhe.Uint32:
			return 45
		}
	case "gte":
		switch uintType {
		case fhe.Uint8:
			return 21
		case fhe.Uint16:
			return 29
		case fhe.Uint32:
			return 45
		}
	case "rem":
		switch uintType {
		case fhe.Uint8:
			return 447
		case fhe.Uint16:
			return 1310
		case fhe.Uint32:
			return 4800
		}
	case "and":
		switch uintType {
		case fhe.Uint8:
			return 13
		case fhe.Uint16:
			return 21
		case fhe.Uint32:
			return 38
		}
	case "or":
		switch uintType {
		case fhe.Uint8:
			return 13
		case fhe.Uint16:
			return 21
		case fhe.Uint32:
			return 38
		}
	case "xor":
		switch uintType {
		case fhe.Uint8:
			return 13
		case fhe.Uint16:
			return 21
		case fhe.Uint32:
			return 38
		}
	case "eq":
		switch uintType {
		case fhe.Uint8:
			return 18
		case fhe.Uint16:
			return 25
		case fhe.Uint32:
			return 50
		}
	case "ne":
		switch uintType {
		case fhe.Uint8:
			return 18
		case fhe.Uint16:
			return 25
		case fhe.Uint32:
			return 50
		}
	case "min":
		switch uintType {
		case fhe.Uint8:
			return 41
		case fhe.Uint16:
			return 75
		case fhe.Uint32:
			return 135
		}
	case "max":
		switch uintType {
		case fhe.Uint8:
			return 41
		case fhe.Uint16:
			return 75
		case fhe.Uint32:
			return 135
		}
	case "shl":
		switch uintType {
		case fhe.Uint8:
			return 82
		case fhe.Uint16:
			return 190
		case fhe.Uint32:
			return 422
		}
	case "shr":
		switch uintType {
		case fhe.Uint8:
			return 82
		case fhe.Uint16:
			return 190
		case fhe.Uint32:
			return 422
		}
	case "not":
		switch uintType {
		case fhe.Uint8:
			return 12
		case fhe.Uint16:
			return 22
		case fhe.Uint32:
			return 36
		}
	default:
		panic("invalid precompile name")
	}

	return 0
}
