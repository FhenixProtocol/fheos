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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"add\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"and\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"cast\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"div\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"eq\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"gt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"gte\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"lt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"lte\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"max\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"min\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"mul\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"ne\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"or\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"reencrypt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"rem\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"req\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"shl\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"shr\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"sub\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"trivialEncrypt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"xor\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"inputLen\",\"type\":\"uint32\"}],\"name\":\"cmux\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// And is a free data retrieval call binding the contract method 0xb2a9fcb6.
//
// Solidity: function and(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) And(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "and", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// And is a free data retrieval call binding the contract method 0xb2a9fcb6.
//
// Solidity: function and(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) And(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.And(&_FheOps.CallOpts, input, inputLen)
}

// And is a free data retrieval call binding the contract method 0xb2a9fcb6.
//
// Solidity: function and(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) And(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.And(&_FheOps.CallOpts, input, inputLen)
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

// Div is a free data retrieval call binding the contract method 0xb530d6a2.
//
// Solidity: function div(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Div(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "div", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Div is a free data retrieval call binding the contract method 0xb530d6a2.
//
// Solidity: function div(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Div(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Div(&_FheOps.CallOpts, input, inputLen)
}

// Div is a free data retrieval call binding the contract method 0xb530d6a2.
//
// Solidity: function div(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Div(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Div(&_FheOps.CallOpts, input, inputLen)
}

// Eq is a free data retrieval call binding the contract method 0x8263a274.
//
// Solidity: function eq(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Eq(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "eq", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Eq is a free data retrieval call binding the contract method 0x8263a274.
//
// Solidity: function eq(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Eq(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Eq(&_FheOps.CallOpts, input, inputLen)
}

// Eq is a free data retrieval call binding the contract method 0x8263a274.
//
// Solidity: function eq(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Eq(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Eq(&_FheOps.CallOpts, input, inputLen)
}

// Gt is a free data retrieval call binding the contract method 0xd998a085.
//
// Solidity: function gt(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Gt(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "gt", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Gt is a free data retrieval call binding the contract method 0xd998a085.
//
// Solidity: function gt(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Gt(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Gt(&_FheOps.CallOpts, input, inputLen)
}

// Gt is a free data retrieval call binding the contract method 0xd998a085.
//
// Solidity: function gt(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Gt(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Gt(&_FheOps.CallOpts, input, inputLen)
}

// Gte is a free data retrieval call binding the contract method 0xdb51bd2b.
//
// Solidity: function gte(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Gte(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "gte", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Gte is a free data retrieval call binding the contract method 0xdb51bd2b.
//
// Solidity: function gte(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Gte(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Gte(&_FheOps.CallOpts, input, inputLen)
}

// Gte is a free data retrieval call binding the contract method 0xdb51bd2b.
//
// Solidity: function gte(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Gte(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Gte(&_FheOps.CallOpts, input, inputLen)
}

// Cmux is a free data retrieval call binding the contract method 0x45849cb7.
//
// Solidity: function cmux(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Cmux(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "cmux", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Cmux is a free data retrieval call binding the contract method 0x45849cb7.
//
// Solidity: function cmux(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Cmux(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Cmux(&_FheOps.CallOpts, input, inputLen)
}

// Cmux is a free data retrieval call binding the contract method 0x45849cb7.
//
// Solidity: function cmux(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Cmux(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Cmux(&_FheOps.CallOpts, input, inputLen)
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

// Max is a free data retrieval call binding the contract method 0xe3a3a288.
//
// Solidity: function max(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Max(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "max", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Max is a free data retrieval call binding the contract method 0xe3a3a288.
//
// Solidity: function max(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Max(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Max(&_FheOps.CallOpts, input, inputLen)
}

// Max is a free data retrieval call binding the contract method 0xe3a3a288.
//
// Solidity: function max(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Max(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Max(&_FheOps.CallOpts, input, inputLen)
}

// Min is a free data retrieval call binding the contract method 0x9b4d270b.
//
// Solidity: function min(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Min(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "min", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Min is a free data retrieval call binding the contract method 0x9b4d270b.
//
// Solidity: function min(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Min(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Min(&_FheOps.CallOpts, input, inputLen)
}

// Min is a free data retrieval call binding the contract method 0x9b4d270b.
//
// Solidity: function min(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Min(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Min(&_FheOps.CallOpts, input, inputLen)
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

// Ne is a free data retrieval call binding the contract method 0xc04bbc88.
//
// Solidity: function ne(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Ne(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "ne", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Ne is a free data retrieval call binding the contract method 0xc04bbc88.
//
// Solidity: function ne(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Ne(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Ne(&_FheOps.CallOpts, input, inputLen)
}

// Ne is a free data retrieval call binding the contract method 0xc04bbc88.
//
// Solidity: function ne(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Ne(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Ne(&_FheOps.CallOpts, input, inputLen)
}

// Or is a free data retrieval call binding the contract method 0x47f25d5e.
//
// Solidity: function or(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Or(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "or", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Or is a free data retrieval call binding the contract method 0x47f25d5e.
//
// Solidity: function or(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Or(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Or(&_FheOps.CallOpts, input, inputLen)
}

// Or is a free data retrieval call binding the contract method 0x47f25d5e.
//
// Solidity: function or(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Or(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Or(&_FheOps.CallOpts, input, inputLen)
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

// Rem is a free data retrieval call binding the contract method 0x55bb1845.
//
// Solidity: function rem(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Rem(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "rem", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Rem is a free data retrieval call binding the contract method 0x55bb1845.
//
// Solidity: function rem(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Rem(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Rem(&_FheOps.CallOpts, input, inputLen)
}

// Rem is a free data retrieval call binding the contract method 0x55bb1845.
//
// Solidity: function rem(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Rem(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Rem(&_FheOps.CallOpts, input, inputLen)
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

// Shl is a free data retrieval call binding the contract method 0x0fb30a4a.
//
// Solidity: function shl(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Shl(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "shl", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Shl is a free data retrieval call binding the contract method 0x0fb30a4a.
//
// Solidity: function shl(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Shl(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Shl(&_FheOps.CallOpts, input, inputLen)
}

// Shl is a free data retrieval call binding the contract method 0x0fb30a4a.
//
// Solidity: function shl(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Shl(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Shl(&_FheOps.CallOpts, input, inputLen)
}

// Shr is a free data retrieval call binding the contract method 0x24840a23.
//
// Solidity: function shr(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Shr(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "shr", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Shr is a free data retrieval call binding the contract method 0x24840a23.
//
// Solidity: function shr(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Shr(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Shr(&_FheOps.CallOpts, input, inputLen)
}

// Shr is a free data retrieval call binding the contract method 0x24840a23.
//
// Solidity: function shr(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Shr(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Shr(&_FheOps.CallOpts, input, inputLen)
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

// Xor is a free data retrieval call binding the contract method 0x17a1787d.
//
// Solidity: function xor(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCaller) Xor(opts *bind.CallOpts, input []byte, inputLen uint32) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "xor", input, inputLen)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Xor is a free data retrieval call binding the contract method 0x17a1787d.
//
// Solidity: function xor(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsSession) Xor(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Xor(&_FheOps.CallOpts, input, inputLen)
}

// Xor is a free data retrieval call binding the contract method 0x17a1787d.
//
// Solidity: function xor(bytes input, uint32 inputLen) view returns(bytes)
func (_FheOps *FheOpsCallerSession) Xor(input []byte, inputLen uint32) ([]byte, error) {
	return _FheOps.Contract.Xor(&_FheOps.CallOpts, input, inputLen)
}
