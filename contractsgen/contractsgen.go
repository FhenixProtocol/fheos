// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractsgen

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// FheOpsMetaData contains all meta data concerning the FheOps contract.
var FheOpsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"add\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"cast\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"lt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"lte\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"mul\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"reencrypt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"req\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"sub\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"trivialEncrypt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"yONATHAN\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// FheOpsABI is the input ABI used to generate the binding from.
// Deprecated: Use FheOpsMetaData.ABI instead.
var FheOpsABI = FheOpsMetaData.ABI

// FheOps is an auto generated Go binding around an Ethereum contract.
type FheOps struct {
	FheOpsCaller     // Read-only binding to the contract
	FheOpsTransactor // Write-only binding to the contract
	FheOpsFilterer   // Log filterer for contract events
}

// FheOpsCaller is an auto generated read-only Go binding around an Ethereum contract.
type FheOpsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FheOpsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FheOpsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FheOpsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FheOpsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FheOpsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FheOpsSession struct {
	Contract     *FheOps           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FheOpsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FheOpsCallerSession struct {
	Contract *FheOpsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FheOpsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FheOpsTransactorSession struct {
	Contract     *FheOpsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FheOpsRaw is an auto generated low-level Go binding around an Ethereum contract.
type FheOpsRaw struct {
	Contract *FheOps // Generic contract binding to access the raw methods on
}

// FheOpsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FheOpsCallerRaw struct {
	Contract *FheOpsCaller // Generic read-only contract binding to access the raw methods on
}

// FheOpsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FheOpsTransactorRaw struct {
	Contract *FheOpsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFheOps creates a new instance of FheOps, bound to a specific deployed contract.
func NewFheOps(address common.Address, backend bind.ContractBackend) (*FheOps, error) {
	contract, err := bindFheOps(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FheOps{FheOpsCaller: FheOpsCaller{contract: contract}, FheOpsTransactor: FheOpsTransactor{contract: contract}, FheOpsFilterer: FheOpsFilterer{contract: contract}}, nil
}

// NewFheOpsCaller creates a new read-only instance of FheOps, bound to a specific deployed contract.
func NewFheOpsCaller(address common.Address, caller bind.ContractCaller) (*FheOpsCaller, error) {
	contract, err := bindFheOps(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FheOpsCaller{contract: contract}, nil
}

// NewFheOpsTransactor creates a new write-only instance of FheOps, bound to a specific deployed contract.
func NewFheOpsTransactor(address common.Address, transactor bind.ContractTransactor) (*FheOpsTransactor, error) {
	contract, err := bindFheOps(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FheOpsTransactor{contract: contract}, nil
}

// NewFheOpsFilterer creates a new log filterer instance of FheOps, bound to a specific deployed contract.
func NewFheOpsFilterer(address common.Address, filterer bind.ContractFilterer) (*FheOpsFilterer, error) {
	contract, err := bindFheOps(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FheOpsFilterer{contract: contract}, nil
}

// bindFheOps binds a generic wrapper to an already deployed contract.
func bindFheOps(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FheOpsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FheOps *FheOpsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FheOps.Contract.FheOpsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FheOps *FheOpsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FheOps.Contract.FheOpsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FheOps *FheOpsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FheOps.Contract.FheOpsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FheOps *FheOpsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FheOps.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FheOps *FheOpsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FheOps.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FheOps *FheOpsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FheOps.Contract.contract.Transact(opts, method, params...)
}

// Add is a free data retrieval call binding the contract method 0x0512ae91.
//
// Solidity: function add(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Add(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "add", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Add is a free data retrieval call binding the contract method 0x0512ae91.
//
// Solidity: function add(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Add(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Add(&_FheOps.CallOpts, input, inputLen)
}

// Add is a free data retrieval call binding the contract method 0x0512ae91.
//
// Solidity: function add(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Add(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Add(&_FheOps.CallOpts, input, inputLen)
}

// Cast is a free data retrieval call binding the contract method 0xc3bedc57.
//
// Solidity: function cast(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Cast(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "cast", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Cast is a free data retrieval call binding the contract method 0xc3bedc57.
//
// Solidity: function cast(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Cast(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Cast(&_FheOps.CallOpts, input, inputLen)
}

// Cast is a free data retrieval call binding the contract method 0xc3bedc57.
//
// Solidity: function cast(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Cast(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Cast(&_FheOps.CallOpts, input, inputLen)
}

// Lt is a free data retrieval call binding the contract method 0x9fcad060.
//
// Solidity: function lt(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Lt(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "lt", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Lt is a free data retrieval call binding the contract method 0x9fcad060.
//
// Solidity: function lt(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Lt(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Lt(&_FheOps.CallOpts, input, inputLen)
}

// Lt is a free data retrieval call binding the contract method 0x9fcad060.
//
// Solidity: function lt(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Lt(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Lt(&_FheOps.CallOpts, input, inputLen)
}

// Lte is a free data retrieval call binding the contract method 0x406a119f.
//
// Solidity: function lte(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Lte(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "lte", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Lte is a free data retrieval call binding the contract method 0x406a119f.
//
// Solidity: function lte(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Lte(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Lte(&_FheOps.CallOpts, input, inputLen)
}

// Lte is a free data retrieval call binding the contract method 0x406a119f.
//
// Solidity: function lte(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Lte(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Lte(&_FheOps.CallOpts, input, inputLen)
}

// Mul is a free data retrieval call binding the contract method 0x4aea45a8.
//
// Solidity: function mul(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Mul(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "mul", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Mul is a free data retrieval call binding the contract method 0x4aea45a8.
//
// Solidity: function mul(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Mul(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Mul(&_FheOps.CallOpts, input, inputLen)
}

// Mul is a free data retrieval call binding the contract method 0x4aea45a8.
//
// Solidity: function mul(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Mul(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Mul(&_FheOps.CallOpts, input, inputLen)
}

// Reencrypt is a free data retrieval call binding the contract method 0x441a9c62.
//
// Solidity: function reencrypt(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Reencrypt(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "reencrypt", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Reencrypt is a free data retrieval call binding the contract method 0x441a9c62.
//
// Solidity: function reencrypt(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Reencrypt(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Reencrypt(&_FheOps.CallOpts, input, inputLen)
}

// Reencrypt is a free data retrieval call binding the contract method 0x441a9c62.
//
// Solidity: function reencrypt(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Reencrypt(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Reencrypt(&_FheOps.CallOpts, input, inputLen)
}

// Req is a free data retrieval call binding the contract method 0x1bed84bf.
//
// Solidity: function req(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Req(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "req", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Req is a free data retrieval call binding the contract method 0x1bed84bf.
//
// Solidity: function req(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Req(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Req(&_FheOps.CallOpts, input, inputLen)
}

// Req is a free data retrieval call binding the contract method 0x1bed84bf.
//
// Solidity: function req(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Req(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Req(&_FheOps.CallOpts, input, inputLen)
}

// Sub is a free data retrieval call binding the contract method 0x255eeeaf.
//
// Solidity: function sub(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Sub(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "sub", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Sub is a free data retrieval call binding the contract method 0x255eeeaf.
//
// Solidity: function sub(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Sub(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Sub(&_FheOps.CallOpts, input, inputLen)
}

// Sub is a free data retrieval call binding the contract method 0x255eeeaf.
//
// Solidity: function sub(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Sub(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Sub(&_FheOps.CallOpts, input, inputLen)
}

// TrivialEncrypt is a free data retrieval call binding the contract method 0x8a52c8c7.
//
// Solidity: function trivialEncrypt(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) TrivialEncrypt(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "trivialEncrypt", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// TrivialEncrypt is a free data retrieval call binding the contract method 0x8a52c8c7.
//
// Solidity: function trivialEncrypt(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) TrivialEncrypt(input []byte) ([]byte, error) {
	return _FheOps.Contract.TrivialEncrypt(&_FheOps.CallOpts, input)
}

// TrivialEncrypt is a free data retrieval call binding the contract method 0x8a52c8c7.
//
// Solidity: function trivialEncrypt(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) TrivialEncrypt(input []byte) ([]byte, error) {
	return _FheOps.Contract.TrivialEncrypt(&_FheOps.CallOpts, input)
}

// Verify is a free data retrieval call binding the contract method 0xf9229595.
//
// Solidity: function verify(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Verify(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "verify", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Verify is a free data retrieval call binding the contract method 0xf9229595.
//
// Solidity: function verify(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Verify(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Verify(&_FheOps.CallOpts, input, inputLen)
}

// Verify is a free data retrieval call binding the contract method 0xf9229595.
//
// Solidity: function verify(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Verify(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Verify(&_FheOps.CallOpts, input, inputLen)
}

// YONATHAN is a free data retrieval call binding the contract method 0xec4e90dc.
//
// Solidity: function yONATHAN(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) YONATHAN(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "yONATHAN", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// YONATHAN is a free data retrieval call binding the contract method 0xec4e90dc.
//
// Solidity: function yONATHAN(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) YONATHAN(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.YONATHAN(&_FheOps.CallOpts, input, inputLen)
}

// YONATHAN is a free data retrieval call binding the contract method 0xec4e90dc.
//
// Solidity: function yONATHAN(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) YONATHAN(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.YONATHAN(&_FheOps.CallOpts, input, inputLen)
}
