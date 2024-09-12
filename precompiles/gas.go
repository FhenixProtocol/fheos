package precompiles

import (
	"github.com/fhenixprotocol/fheos/precompiles/types"
	"github.com/fhenixprotocol/warp-drive/fhe-driver"
)

func getGasForPrecompile(precompileName types.PrecompileName, uintType fhe.EncryptionType) uint64 {
	return getRawPrecompileGas(precompileName, uintType)
}

func getRawPrecompileGas(precompileName types.PrecompileName, uintType fhe.EncryptionType) uint64 {
	switch precompileName {
	case types.Verify:
		switch uintType {
		case fhe.Uint8, fhe.Uint16, fhe.Uint32:
			return 65000
		case fhe.Uint64, fhe.Uint128, fhe.Uint256, fhe.Address:
			return 300000
		}
	case types.Cast:
		switch uintType {
		case fhe.Uint8:
			return 75000
		case fhe.Uint16:
			return 85000
		case fhe.Uint32:
			return 105000
		case fhe.Uint64:
			return 120000
		case fhe.Uint128:
			return 140000
		case fhe.Uint256:
			return 175000
		case fhe.Address:
			return 150000
		}
	case types.SealOutput, types.Require:
		return 150000
	case types.Decrypt:
		switch uintType {
		case fhe.Uint8:
			return 25000
		default:
			return 150000
		}
	case types.Sub, types.Add:
		switch uintType {
		case fhe.Uint8:
			return 50000
		case fhe.Uint16:
			return 65000
		case fhe.Uint32:
			return 120000
		case fhe.Uint64:
			return 175000
		case fhe.Uint128:
			return 290000
		}
	case types.Mul, types.Square:
		switch uintType {
		case fhe.Uint8:
			return 40000
		case fhe.Uint16:
			return 70000
		case fhe.Uint32:
			return 125000
		case fhe.Uint64:
			return 280000
		}
	case types.Select:
		switch uintType {
		case fhe.Uint8, fhe.Uint16:
			return 55000
		case fhe.Uint32:
			return 85000
		case fhe.Uint64:
			return 125000
		case fhe.Uint128:
			return 225000
		case fhe.Bool:
			return 35000
		}
	case types.Div, types.Rem:
		switch uintType {
		case fhe.Uint8:
			return 125000
		case fhe.Uint16:
			return 335000
		case fhe.Uint32:
			return 1003000
		case fhe.Uint64:
			return 2000000
		}
	case types.Gt, types.Lt, types.Gte, types.Lte:
		switch uintType {
		case fhe.Uint8:
			return 40000
		case fhe.Uint16:
			return 50000
		case fhe.Uint32:
			return 75000
		case fhe.Uint64:
			return 125000
		case fhe.Uint128:
			return 190000
		}
	case types.Or, types.Xor, types.And:
		switch uintType {
		case fhe.Uint8:
			return 40000
		case fhe.Uint16:
			return 50000
		case fhe.Uint32:
			return 75000
		case fhe.Uint64:
			return 130000
		case fhe.Uint128:
			return 200000
		case fhe.Bool:
			return 28000
		}
	case types.Eq, types.Ne:
		switch uintType {
		case fhe.Uint8:
			return 40000
		case fhe.Uint16:
			return 50000
		case fhe.Uint32:
			return 65000
		case fhe.Uint64:
			return 120000
		case fhe.Uint128:
			return 180000
		case fhe.Uint256:
			return 260000
		case fhe.Bool:
			return 35000
		case fhe.Address:
			return 210000
		}
	case types.Min, types.Max:
		switch uintType {
		case fhe.Uint8:
			return 45000
		case fhe.Uint16:
			return 55000
		case fhe.Uint32:
			return 100000
		case fhe.Uint64:
			return 145000
		case fhe.Uint128:
			return 250000
		}
	case types.Shl, types.Shr:
		switch uintType {
		case fhe.Uint8:
			return 65000
		case fhe.Uint16:
			return 90000
		case fhe.Uint32:
			return 130000
		case fhe.Uint64:
			return 210000
		case fhe.Uint128:
			return 355000
		}
	case types.Not:
		switch uintType {
		case fhe.Uint8:
			return 42000
		case fhe.Uint16:
			return 35000
		case fhe.Uint32:
			return 49000
		case fhe.Uint64:
			return 85000
		case fhe.Uint128:
			return 120000
		case fhe.Bool:
			return 28000
		}
	case types.GetNetworkKey:
		// this is never meant to be called in the context of a tx, so we give a pretty high gas cost just to avoid DoS
		return 200000
	case types.TrivialEncrypt:
		switch uintType {
		case fhe.Uint8, fhe.Uint16:
			return 20000
		case fhe.Uint32:
			return 30000
		case fhe.Uint64:
			return 35000
		case fhe.Uint128:
			return 65000
		case fhe.Uint256, fhe.Address:
			return 70000
		}
	case types.Random:
		return 120000
	}
	return 0
}
