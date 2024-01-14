package precompiles

import "github.com/fhenixprotocol/go-tfhe"

// we calculated gas using the results of the following test results: https://github.com/FhenixProtocol/fheos/actions/runs/7475858194/job/20344926903
// we considered the time that it took to run division of two euint32s as a pivot that will take 1M gas and calculated how much gas is if for 1ms of runtime
// 1M / 4800ms = 209 gas per 1ms

func getGasForPrecompile(precompileName string, uintType tfhe.UintType) uint64 {
	switch precompileName {
	case "add":
		switch uintType {
		case tfhe.Uint8:
			return 7315 // 35 * 209
		case tfhe.Uint16:
			return 17257 // 73 * 209
		case tfhe.Uint32:
			return 34067 // 163 * 209
		}
	case "verify":
		switch uintType {
		case tfhe.Uint8:
			return 9196 // 44 * 209
		case tfhe.Uint16:
			return 9196 // 44 * 209
		case tfhe.Uint32:
			return 9196 // 44 * 209
		}
	case "sealOutput":
		switch uintType {
		case tfhe.Uint8:
			return 4807 // 23 * 209
		case tfhe.Uint16:
			return 4807 // 23 * 209
		case tfhe.Uint32:
			return 4807 // 23 * 209
		}
	case "decrypt": // No test, roughly estimated as sealOutput - trivialEncrypt
		switch uintType {
		case tfhe.Uint8:
			return 3762 // 18 * 209
		case tfhe.Uint16:
			return 3762 // 18 * 209
		case tfhe.Uint32:
			return 3762 // 18 * 209
		}
	case "lte":
		switch uintType {
		case tfhe.Uint8:
			return 3762 // 18 * 209
		case tfhe.Uint16:
			return 6061 // 29 * 209
		case tfhe.Uint32:
			return 8987 // 43 * 209
		}
	case "sub":
		switch uintType {
		case tfhe.Uint8:
			return 7315 // 35 * 209
		case tfhe.Uint16:
			return 17257 // 73 * 209
		case tfhe.Uint32:
			return 34067 // 163 * 209
		}
	case "mul":
		switch uintType {
		case tfhe.Uint8:
			return 19646 // 94 * 209
		case tfhe.Uint16:
			return 64999 // 311 * 209
		case tfhe.Uint32:
			return 235543 // 1127 * 209
		}
	case "lt":
		switch uintType {
		case tfhe.Uint8:
			return 4598 // 22 * 209
		case tfhe.Uint16:
			return 7524 // 36 * 209
		case tfhe.Uint32:
			return 12540 // 60 * 209
		}
	case "select":
		switch uintType {
		case tfhe.Uint8:
			return 44726 // 214 * 209
		case tfhe.Uint16:
			return 66044 // 316 * 209
		case tfhe.Uint32:
			return 114741 // 549 * 209
		}
	case "require": // Took the values when there was no crash as for crash gas is irrelevant as it will be reverted
		switch uintType {
		case tfhe.Uint8:
			return 13585 // 65 * 209
		case tfhe.Uint16:
			return 13585 // 65 * 209
		case tfhe.Uint32:
			return 13585 // 65 * 209
		}
	case "trivialEncrypt":
		switch uintType {
		case tfhe.Uint8:
			return 1045 // 5 * 209
		case tfhe.Uint16:
			return 1045 // 5 * 209
		case tfhe.Uint32:
			return 1045 // 5 * 209
		}
	case "div":
		switch uintType {
		case tfhe.Uint8:
			return 93423 // 447 * 209
		case tfhe.Uint16:
			return 273790 // 1310 * 209
		case tfhe.Uint32:
			return 1003200 // 4800 * 209
		}
	case "gt":
		switch uintType {
		case tfhe.Uint8:
			return 4389 // 21 * 209
		case tfhe.Uint16:
			return 6061 // 29 * 209
		case tfhe.Uint32:
			return 9405 // 45 * 209
		}
	case "gte":
		switch uintType {
		case tfhe.Uint8:
			return 4389 // 21 * 209
		case tfhe.Uint16:
			return 6061 // 29 * 209
		case tfhe.Uint32:
			return 9405 // 45 * 209
		}
	case "rem":
		switch uintType {
		case tfhe.Uint8:
			return 93423 // 447 * 209
		case tfhe.Uint16:
			return 273790 // 1310 * 209
		case tfhe.Uint32:
			return 1003200 // 4800 * 209
		}
	case "and":
		switch uintType {
		case tfhe.Uint8:
			return 2717 // 13 * 209
		case tfhe.Uint16:
			return 4389 // 21 * 209
		case tfhe.Uint32:
			return 7942 // 38 * 209
		}
	case "or":
		switch uintType {
		case tfhe.Uint8:
			return 2717 // 13 * 209
		case tfhe.Uint16:
			return 4389 // 21 * 209
		case tfhe.Uint32:
			return 7942 // 38 * 209
		}
	case "xor":
		switch uintType {
		case tfhe.Uint8:
			return 2717 // 13 * 209
		case tfhe.Uint16:
			return 4389 // 21 * 209
		case tfhe.Uint32:
			return 7942 // 38 * 209
		}
	case "eq":
		switch uintType {
		case tfhe.Uint8:
			return 3762 // 18 * 209
		case tfhe.Uint16:
			return 5225 // 25 * 209
		case tfhe.Uint32:
			return 10450 // 50 * 209
		}
	case "ne":
		switch uintType {
		case tfhe.Uint8:
			return 3762 // 18 * 209
		case tfhe.Uint16:
			return 5225 // 25 * 209
		case tfhe.Uint32:
			return 10450 // 50 * 209
		}
	case "min":
		switch uintType {
		case tfhe.Uint8:
			return 8569 // 41 * 209
		case tfhe.Uint16:
			return 15675 // 75 * 209
		case tfhe.Uint32:
			return 28215 // 135 * 209
		}
	case "max":
		switch uintType {
		case tfhe.Uint8:
			return 8569 // 41 * 209
		case tfhe.Uint16:
			return 15675 // 75 * 209
		case tfhe.Uint32:
			return 88198 // 135 * 209
		}
	case "shl":
		switch uintType {
		case tfhe.Uint8:
			return 17138 // 82 * 209
		case tfhe.Uint16:
			return 39710 // 190 * 209
		case tfhe.Uint32:
			return 88198 // 422 * 209
		}
	case "shr":
		switch uintType {
		case tfhe.Uint8:
			return 17138 // 82 * 209
		case tfhe.Uint16:
			return 39710 // 190 * 209
		case tfhe.Uint32:
			return 88198 // 422 * 209
		}
	case "not":
		switch uintType {
		case tfhe.Uint8:
			return 2508 // 12 * 209
		case tfhe.Uint16:
			return 4598 // 22 * 209
		case tfhe.Uint32:
			return 7524 // 36 * 209
		}
	default:
		panic("invalid precompile name")
	}

	return 0
}
