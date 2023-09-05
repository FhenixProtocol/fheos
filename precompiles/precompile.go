// Copyright 2021-2023, Offchain Labs, Inc.
// For license information, see https://github.com/OffchainLabs/nitro/blob/master/LICENSE

package precompiles

import (
	template "github.com/FhenixProtocol/fheos/go/contractsgen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type addr = common.Address
type Precompile struct {
	Metadata *bind.MetaData
	Address  addr
}

func GetPrecompilesList() []Precompile {
	precompiles := []Precompile{}
	append := func(metadata *bind.MetaData, a addr) {
		var precompile Precompile
		precompile.Metadata = metadata
		precompile.Address = a
		precompiles = append(precompiles, precompile)
	}

	hex := func(s string) addr {
		return common.HexToAddress(s)
	}

	append(template.FheOpsMetaData, hex("80"))
	return precompiles

}
