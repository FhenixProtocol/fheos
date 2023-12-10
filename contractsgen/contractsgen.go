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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"add\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"and\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"cast\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"cmux\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"div\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"eq\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"gt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"gte\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"lt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"lte\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"max\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"min\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"mul\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"ne\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"or\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"reencrypt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"rem\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"req\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"shl\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"shr\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"sub\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"trivialEncrypt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"xor\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// Add is a free data retrieval call binding the contract method 0xba658111.
//
// Solidity: function add(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Add(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "add", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Add is a free data retrieval call binding the contract method 0xba658111.
//
// Solidity: function add(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Add(input []byte) ([]byte, error) {
	return _FheOps.Contract.Add(&_FheOps.CallOpts, input)
}

// Add is a free data retrieval call binding the contract method 0xba658111.
//
// Solidity: function add(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Add(input []byte) ([]byte, error) {
	return _FheOps.Contract.Add(&_FheOps.CallOpts, input)
}

// And is a free data retrieval call binding the contract method 0x378c56ed.
//
// Solidity: function and(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) And(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "and", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// And is a free data retrieval call binding the contract method 0x378c56ed.
//
// Solidity: function and(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) And(input []byte) ([]byte, error) {
	return _FheOps.Contract.And(&_FheOps.CallOpts, input)
}

// And is a free data retrieval call binding the contract method 0x378c56ed.
//
// Solidity: function and(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) And(input []byte) ([]byte, error) {
	return _FheOps.Contract.And(&_FheOps.CallOpts, input)
}

// Cast is a free data retrieval call binding the contract method 0x756a210d.
//
// Solidity: function cast(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Cast(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "cast", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Cast is a free data retrieval call binding the contract method 0x756a210d.
//
// Solidity: function cast(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Cast(input []byte) ([]byte, error) {
	return _FheOps.Contract.Cast(&_FheOps.CallOpts, input)
}

// Cast is a free data retrieval call binding the contract method 0x756a210d.
//
// Solidity: function cast(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Cast(input []byte) ([]byte, error) {
	return _FheOps.Contract.Cast(&_FheOps.CallOpts, input)
}

// Cmux is a free data retrieval call binding the contract method 0xe70cc6df.
//
// Solidity: function cmux(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Cmux(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "cmux", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Cmux is a free data retrieval call binding the contract method 0xe70cc6df.
//
// Solidity: function cmux(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Cmux(input []byte) ([]byte, error) {
	return _FheOps.Contract.Cmux(&_FheOps.CallOpts, input)
}

// Cmux is a free data retrieval call binding the contract method 0xe70cc6df.
//
// Solidity: function cmux(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Cmux(input []byte) ([]byte, error) {
	return _FheOps.Contract.Cmux(&_FheOps.CallOpts, input)
}

// Div is a free data retrieval call binding the contract method 0xed0dd1f1.
//
// Solidity: function div(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Div(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "div", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Div is a free data retrieval call binding the contract method 0xed0dd1f1.
//
// Solidity: function div(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Div(input []byte) ([]byte, error) {
	return _FheOps.Contract.Div(&_FheOps.CallOpts, input)
}

// Div is a free data retrieval call binding the contract method 0xed0dd1f1.
//
// Solidity: function div(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Div(input []byte) ([]byte, error) {
	return _FheOps.Contract.Div(&_FheOps.CallOpts, input)
}

// Eq is a free data retrieval call binding the contract method 0x1868b889.
//
// Solidity: function eq(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Eq(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "eq", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Eq is a free data retrieval call binding the contract method 0x1868b889.
//
// Solidity: function eq(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Eq(input []byte) ([]byte, error) {
	return _FheOps.Contract.Eq(&_FheOps.CallOpts, input)
}

// Eq is a free data retrieval call binding the contract method 0x1868b889.
//
// Solidity: function eq(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Eq(input []byte) ([]byte, error) {
	return _FheOps.Contract.Eq(&_FheOps.CallOpts, input)
}

// Gt is a free data retrieval call binding the contract method 0x3b902188.
//
// Solidity: function gt(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Gt(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "gt", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Gt is a free data retrieval call binding the contract method 0x3b902188.
//
// Solidity: function gt(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Gt(input []byte) ([]byte, error) {
	return _FheOps.Contract.Gt(&_FheOps.CallOpts, input)
}

// Gt is a free data retrieval call binding the contract method 0x3b902188.
//
// Solidity: function gt(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Gt(input []byte) ([]byte, error) {
	return _FheOps.Contract.Gt(&_FheOps.CallOpts, input)
}

// Gte is a free data retrieval call binding the contract method 0x08d6e8cf.
//
// Solidity: function gte(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Gte(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "gte", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Gte is a free data retrieval call binding the contract method 0x08d6e8cf.
//
// Solidity: function gte(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Gte(input []byte) ([]byte, error) {
	return _FheOps.Contract.Gte(&_FheOps.CallOpts, input)
}

// Gte is a free data retrieval call binding the contract method 0x08d6e8cf.
//
// Solidity: function gte(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Gte(input []byte) ([]byte, error) {
	return _FheOps.Contract.Gte(&_FheOps.CallOpts, input)
}

// Lt is a free data retrieval call binding the contract method 0x9d8a3b5a.
//
// Solidity: function lt(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Lt(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "lt", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Lt is a free data retrieval call binding the contract method 0x9d8a3b5a.
//
// Solidity: function lt(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Lt(input []byte) ([]byte, error) {
	return _FheOps.Contract.Lt(&_FheOps.CallOpts, input)
}

// Lt is a free data retrieval call binding the contract method 0x9d8a3b5a.
//
// Solidity: function lt(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Lt(input []byte) ([]byte, error) {
	return _FheOps.Contract.Lt(&_FheOps.CallOpts, input)
}

// Lte is a free data retrieval call binding the contract method 0xb3dfb138.
//
// Solidity: function lte(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Lte(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "lte", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Lte is a free data retrieval call binding the contract method 0xb3dfb138.
//
// Solidity: function lte(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Lte(input []byte) ([]byte, error) {
	return _FheOps.Contract.Lte(&_FheOps.CallOpts, input)
}

// Lte is a free data retrieval call binding the contract method 0xb3dfb138.
//
// Solidity: function lte(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Lte(input []byte) ([]byte, error) {
	return _FheOps.Contract.Lte(&_FheOps.CallOpts, input)
}

// Max is a free data retrieval call binding the contract method 0xaa11c926.
//
// Solidity: function max(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Max(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "max", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Max is a free data retrieval call binding the contract method 0xaa11c926.
//
// Solidity: function max(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Max(input []byte) ([]byte, error) {
	return _FheOps.Contract.Max(&_FheOps.CallOpts, input)
}

// Max is a free data retrieval call binding the contract method 0xaa11c926.
//
// Solidity: function max(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Max(input []byte) ([]byte, error) {
	return _FheOps.Contract.Max(&_FheOps.CallOpts, input)
}

// Min is a free data retrieval call binding the contract method 0x6583520e.
//
// Solidity: function min(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Min(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "min", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Min is a free data retrieval call binding the contract method 0x6583520e.
//
// Solidity: function min(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Min(input []byte) ([]byte, error) {
	return _FheOps.Contract.Min(&_FheOps.CallOpts, input)
}

// Min is a free data retrieval call binding the contract method 0x6583520e.
//
// Solidity: function min(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Min(input []byte) ([]byte, error) {
	return _FheOps.Contract.Min(&_FheOps.CallOpts, input)
}

// Mul is a free data retrieval call binding the contract method 0x036ad00f.
//
// Solidity: function mul(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Mul(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "mul", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Mul is a free data retrieval call binding the contract method 0x036ad00f.
//
// Solidity: function mul(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Mul(input []byte) ([]byte, error) {
	return _FheOps.Contract.Mul(&_FheOps.CallOpts, input)
}

// Mul is a free data retrieval call binding the contract method 0x036ad00f.
//
// Solidity: function mul(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Mul(input []byte) ([]byte, error) {
	return _FheOps.Contract.Mul(&_FheOps.CallOpts, input)
}

// Ne is a free data retrieval call binding the contract method 0xd903ba51.
//
// Solidity: function ne(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Ne(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "ne", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Ne is a free data retrieval call binding the contract method 0xd903ba51.
//
// Solidity: function ne(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Ne(input []byte) ([]byte, error) {
	return _FheOps.Contract.Ne(&_FheOps.CallOpts, input)
}

// Ne is a free data retrieval call binding the contract method 0xd903ba51.
//
// Solidity: function ne(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Ne(input []byte) ([]byte, error) {
	return _FheOps.Contract.Ne(&_FheOps.CallOpts, input)
}

// Or is a free data retrieval call binding the contract method 0xf081b3dc.
//
// Solidity: function or(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Or(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "or", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Or is a free data retrieval call binding the contract method 0xf081b3dc.
//
// Solidity: function or(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Or(input []byte) ([]byte, error) {
	return _FheOps.Contract.Or(&_FheOps.CallOpts, input)
}

// Or is a free data retrieval call binding the contract method 0xf081b3dc.
//
// Solidity: function or(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Or(input []byte) ([]byte, error) {
	return _FheOps.Contract.Or(&_FheOps.CallOpts, input)
}

// Reencrypt is a free data retrieval call binding the contract method 0xd77357e1.
//
// Solidity: function reencrypt(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Reencrypt(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "reencrypt", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Reencrypt is a free data retrieval call binding the contract method 0xd77357e1.
//
// Solidity: function reencrypt(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Reencrypt(input []byte) ([]byte, error) {
	return _FheOps.Contract.Reencrypt(&_FheOps.CallOpts, input)
}

// Reencrypt is a free data retrieval call binding the contract method 0xd77357e1.
//
// Solidity: function reencrypt(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Reencrypt(input []byte) ([]byte, error) {
	return _FheOps.Contract.Reencrypt(&_FheOps.CallOpts, input)
}

// Rem is a free data retrieval call binding the contract method 0xae07ec6b.
//
// Solidity: function rem(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Rem(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "rem", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Rem is a free data retrieval call binding the contract method 0xae07ec6b.
//
// Solidity: function rem(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Rem(input []byte) ([]byte, error) {
	return _FheOps.Contract.Rem(&_FheOps.CallOpts, input)
}

// Rem is a free data retrieval call binding the contract method 0xae07ec6b.
//
// Solidity: function rem(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Rem(input []byte) ([]byte, error) {
	return _FheOps.Contract.Rem(&_FheOps.CallOpts, input)
}

// Req is a free data retrieval call binding the contract method 0xac6c08dd.
//
// Solidity: function req(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Req(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "req", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Req is a free data retrieval call binding the contract method 0xac6c08dd.
//
// Solidity: function req(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Req(input []byte) ([]byte, error) {
	return _FheOps.Contract.Req(&_FheOps.CallOpts, input)
}

// Req is a free data retrieval call binding the contract method 0xac6c08dd.
//
// Solidity: function req(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Req(input []byte) ([]byte, error) {
	return _FheOps.Contract.Req(&_FheOps.CallOpts, input)
}

// Shl is a free data retrieval call binding the contract method 0xea9cd829.
//
// Solidity: function shl(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Shl(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "shl", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Shl is a free data retrieval call binding the contract method 0xea9cd829.
//
// Solidity: function shl(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Shl(input []byte) ([]byte, error) {
	return _FheOps.Contract.Shl(&_FheOps.CallOpts, input)
}

// Shl is a free data retrieval call binding the contract method 0xea9cd829.
//
// Solidity: function shl(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Shl(input []byte) ([]byte, error) {
	return _FheOps.Contract.Shl(&_FheOps.CallOpts, input)
}

// Shr is a free data retrieval call binding the contract method 0xf8ab927d.
//
// Solidity: function shr(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Shr(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "shr", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Shr is a free data retrieval call binding the contract method 0xf8ab927d.
//
// Solidity: function shr(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Shr(input []byte) ([]byte, error) {
	return _FheOps.Contract.Shr(&_FheOps.CallOpts, input)
}

// Shr is a free data retrieval call binding the contract method 0xf8ab927d.
//
// Solidity: function shr(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Shr(input []byte) ([]byte, error) {
	return _FheOps.Contract.Shr(&_FheOps.CallOpts, input)
}

// Sub is a free data retrieval call binding the contract method 0x67d1438e.
//
// Solidity: function sub(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Sub(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "sub", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Sub is a free data retrieval call binding the contract method 0x67d1438e.
//
// Solidity: function sub(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Sub(input []byte) ([]byte, error) {
	return _FheOps.Contract.Sub(&_FheOps.CallOpts, input)
}

// Sub is a free data retrieval call binding the contract method 0x67d1438e.
//
// Solidity: function sub(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Sub(input []byte) ([]byte, error) {
	return _FheOps.Contract.Sub(&_FheOps.CallOpts, input)
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

// Verify is a free data retrieval call binding the contract method 0x8e760afe.
//
// Solidity: function verify(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Verify(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "verify", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Verify is a free data retrieval call binding the contract method 0x8e760afe.
//
// Solidity: function verify(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Verify(input []byte) ([]byte, error) {
	return _FheOps.Contract.Verify(&_FheOps.CallOpts, input)
}

// Verify is a free data retrieval call binding the contract method 0x8e760afe.
//
// Solidity: function verify(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Verify(input []byte) ([]byte, error) {
	return _FheOps.Contract.Verify(&_FheOps.CallOpts, input)
}

// Xor is a free data retrieval call binding the contract method 0xdded6e15.
//
// Solidity: function xor(bytes input) view returns(bytes)
func (_FheOps *FheOpsCaller) Xor(opts *bind.CallOpts, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "xor", input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Xor is a free data retrieval call binding the contract method 0xdded6e15.
//
// Solidity: function xor(bytes input) view returns(bytes)
func (_FheOps *FheOpsSession) Xor(input []byte) ([]byte, error) {
	return _FheOps.Contract.Xor(&_FheOps.CallOpts, input)
}

// Xor is a free data retrieval call binding the contract method 0xdded6e15.
//
// Solidity: function xor(bytes input) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Xor(input []byte) ([]byte, error) {
	return _FheOps.Contract.Xor(&_FheOps.CallOpts, input)
}
