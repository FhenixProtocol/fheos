package arbitrum

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	template "github.com/fhenixprotocol/fheos/chains/arbitrum/contractsgen"
	types "github.com/fhenixprotocol/fheos/precompiles"
)

func GetPrecompilesList() []types.Precompile {
	precompiles := []types.Precompile{}
	append := func(metadata *bind.MetaData, a common.Address) {
		var precompile types.Precompile
		precompile.Metadata = metadata
		precompile.Address = a
		precompiles = append(precompiles, precompile)
	}

	hex := func(s string) common.Address {
		return common.HexToAddress(s)
	}

	append(template.FheOpsMetaData, hex("80"))
	return precompiles
}
