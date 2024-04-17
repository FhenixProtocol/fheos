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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"add\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"and\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"toType\",\"type\":\"uint8\"}],\"name\":\"cast\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"decrypt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"div\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"eq\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNetworkPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"gt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"gte\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"s\",\"type\":\"string\"}],\"name\":\"log\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"lt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"lte\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"max\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"min\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"mul\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"ne\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"not\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"or\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"rem\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"req\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"ctHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pk\",\"type\":\"bytes\"}],\"name\":\"sealOutput\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"controlHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"ifTrueHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"ifFalseHash\",\"type\":\"bytes\"}],\"name\":\"select\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"shl\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"shr\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"sub\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"toType\",\"type\":\"uint8\"}],\"name\":\"trivialEncrypt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"utype\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"lhsHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rhsHash\",\"type\":\"bytes\"}],\"name\":\"xor\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
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

// Add is a free data retrieval call binding the contract method 0x002df619.
//
// Solidity: function add(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Add(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "add", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Add is a free data retrieval call binding the contract method 0x002df619.
//
// Solidity: function add(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Add(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Add(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Add is a free data retrieval call binding the contract method 0x002df619.
//
// Solidity: function add(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Add(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Add(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// And is a free data retrieval call binding the contract method 0xae104cfd.
//
// Solidity: function and(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) And(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "and", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// And is a free data retrieval call binding the contract method 0xae104cfd.
//
// Solidity: function and(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) And(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.And(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// And is a free data retrieval call binding the contract method 0xae104cfd.
//
// Solidity: function and(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) And(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.And(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Cast is a free data retrieval call binding the contract method 0x4a5a1117.
//
// Solidity: function cast(uint8 utype, bytes input, uint8 toType) pure returns(bytes)
func (_FheOps *FheOpsCaller) Cast(opts *bind.CallOpts, utype uint8, input []byte, toType uint8) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "cast", utype, input, toType)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Cast is a free data retrieval call binding the contract method 0x4a5a1117.
//
// Solidity: function cast(uint8 utype, bytes input, uint8 toType) pure returns(bytes)
func (_FheOps *FheOpsSession) Cast(utype uint8, input []byte, toType uint8) ([]byte, error) {
	return _FheOps.Contract.Cast(&_FheOps.CallOpts, utype, input, toType)
}

// Cast is a free data retrieval call binding the contract method 0x4a5a1117.
//
// Solidity: function cast(uint8 utype, bytes input, uint8 toType) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Cast(utype uint8, input []byte, toType uint8) ([]byte, error) {
	return _FheOps.Contract.Cast(&_FheOps.CallOpts, utype, input, toType)
}

// Decrypt is a free data retrieval call binding the contract method 0x73cc0154.
//
// Solidity: function decrypt(uint8 utype, bytes input) pure returns(uint256)
func (_FheOps *FheOpsCaller) Decrypt(opts *bind.CallOpts, utype uint8, input []byte) (*big.Int, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "decrypt", utype, input)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Decrypt is a free data retrieval call binding the contract method 0x73cc0154.
//
// Solidity: function decrypt(uint8 utype, bytes input) pure returns(uint256)
func (_FheOps *FheOpsSession) Decrypt(utype uint8, input []byte) (*big.Int, error) {
	return _FheOps.Contract.Decrypt(&_FheOps.CallOpts, utype, input)
}

// Decrypt is a free data retrieval call binding the contract method 0x73cc0154.
//
// Solidity: function decrypt(uint8 utype, bytes input) pure returns(uint256)
func (_FheOps *FheOpsCallerSession) Decrypt(utype uint8, input []byte) (*big.Int, error) {
	return _FheOps.Contract.Decrypt(&_FheOps.CallOpts, utype, input)
}

// Div is a free data retrieval call binding the contract method 0x1f4cda2f.
//
// Solidity: function div(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Div(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "div", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Div is a free data retrieval call binding the contract method 0x1f4cda2f.
//
// Solidity: function div(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Div(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Div(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Div is a free data retrieval call binding the contract method 0x1f4cda2f.
//
// Solidity: function div(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Div(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Div(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Eq is a free data retrieval call binding the contract method 0x92348b34.
//
// Solidity: function eq(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Eq(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "eq", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Eq is a free data retrieval call binding the contract method 0x92348b34.
//
// Solidity: function eq(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Eq(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Eq(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Eq is a free data retrieval call binding the contract method 0x92348b34.
//
// Solidity: function eq(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Eq(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Eq(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// GetNetworkPublicKey is a free data retrieval call binding the contract method 0x44e21dd2.
//
// Solidity: function getNetworkPublicKey() pure returns(bytes)
func (_FheOps *FheOpsCaller) GetNetworkPublicKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "getNetworkPublicKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetNetworkPublicKey is a free data retrieval call binding the contract method 0x44e21dd2.
//
// Solidity: function getNetworkPublicKey() pure returns(bytes)
func (_FheOps *FheOpsSession) GetNetworkPublicKey() ([]byte, error) {
	return _FheOps.Contract.GetNetworkPublicKey(&_FheOps.CallOpts)
}

// GetNetworkPublicKey is a free data retrieval call binding the contract method 0x44e21dd2.
//
// Solidity: function getNetworkPublicKey() pure returns(bytes)
func (_FheOps *FheOpsCallerSession) GetNetworkPublicKey() ([]byte, error) {
	return _FheOps.Contract.GetNetworkPublicKey(&_FheOps.CallOpts)
}

// Gt is a free data retrieval call binding the contract method 0x874b1c10.
//
// Solidity: function gt(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Gt(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "gt", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Gt is a free data retrieval call binding the contract method 0x874b1c10.
//
// Solidity: function gt(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Gt(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Gt(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Gt is a free data retrieval call binding the contract method 0x874b1c10.
//
// Solidity: function gt(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Gt(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Gt(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Gte is a free data retrieval call binding the contract method 0x650de1cf.
//
// Solidity: function gte(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Gte(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "gte", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Gte is a free data retrieval call binding the contract method 0x650de1cf.
//
// Solidity: function gte(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Gte(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Gte(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Gte is a free data retrieval call binding the contract method 0x650de1cf.
//
// Solidity: function gte(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Gte(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Gte(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Log is a free data retrieval call binding the contract method 0x41304fac.
//
// Solidity: function log(string s) pure returns()
func (_FheOps *FheOpsCaller) Log(opts *bind.CallOpts, s string) error {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "log", s)

	if err != nil {
		return err
	}

	return err

}

// Log is a free data retrieval call binding the contract method 0x41304fac.
//
// Solidity: function log(string s) pure returns()
func (_FheOps *FheOpsSession) Log(s string) error {
	return _FheOps.Contract.Log(&_FheOps.CallOpts, s)
}

// Log is a free data retrieval call binding the contract method 0x41304fac.
//
// Solidity: function log(string s) pure returns()
func (_FheOps *FheOpsCallerSession) Log(s string) error {
	return _FheOps.Contract.Log(&_FheOps.CallOpts, s)
}

// Lt is a free data retrieval call binding the contract method 0xb9c7a54b.
//
// Solidity: function lt(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Lt(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "lt", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Lt is a free data retrieval call binding the contract method 0xb9c7a54b.
//
// Solidity: function lt(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Lt(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Lt(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Lt is a free data retrieval call binding the contract method 0xb9c7a54b.
//
// Solidity: function lt(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Lt(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Lt(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Lte is a free data retrieval call binding the contract method 0xeb274b77.
//
// Solidity: function lte(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Lte(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "lte", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Lte is a free data retrieval call binding the contract method 0xeb274b77.
//
// Solidity: function lte(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Lte(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Lte(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Lte is a free data retrieval call binding the contract method 0xeb274b77.
//
// Solidity: function lte(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Lte(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Lte(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Max is a free data retrieval call binding the contract method 0x0b80518e.
//
// Solidity: function max(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Max(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "max", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Max is a free data retrieval call binding the contract method 0x0b80518e.
//
// Solidity: function max(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Max(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Max(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Max is a free data retrieval call binding the contract method 0x0b80518e.
//
// Solidity: function max(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Max(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Max(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Min is a free data retrieval call binding the contract method 0x5211c679.
//
// Solidity: function min(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Min(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "min", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Min is a free data retrieval call binding the contract method 0x5211c679.
//
// Solidity: function min(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Min(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Min(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Min is a free data retrieval call binding the contract method 0x5211c679.
//
// Solidity: function min(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Min(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Min(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Mul is a free data retrieval call binding the contract method 0x4284f765.
//
// Solidity: function mul(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Mul(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "mul", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Mul is a free data retrieval call binding the contract method 0x4284f765.
//
// Solidity: function mul(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Mul(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Mul(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Mul is a free data retrieval call binding the contract method 0x4284f765.
//
// Solidity: function mul(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Mul(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Mul(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Ne is a free data retrieval call binding the contract method 0x13c0c9ae.
//
// Solidity: function ne(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Ne(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "ne", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Ne is a free data retrieval call binding the contract method 0x13c0c9ae.
//
// Solidity: function ne(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Ne(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Ne(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Ne is a free data retrieval call binding the contract method 0x13c0c9ae.
//
// Solidity: function ne(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Ne(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Ne(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Not is a free data retrieval call binding the contract method 0xd260d9ab.
//
// Solidity: function not(uint8 utype, bytes value) pure returns(bytes)
func (_FheOps *FheOpsCaller) Not(opts *bind.CallOpts, utype uint8, value []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "not", utype, value)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Not is a free data retrieval call binding the contract method 0xd260d9ab.
//
// Solidity: function not(uint8 utype, bytes value) pure returns(bytes)
func (_FheOps *FheOpsSession) Not(utype uint8, value []byte) ([]byte, error) {
	return _FheOps.Contract.Not(&_FheOps.CallOpts, utype, value)
}

// Not is a free data retrieval call binding the contract method 0xd260d9ab.
//
// Solidity: function not(uint8 utype, bytes value) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Not(utype uint8, value []byte) ([]byte, error) {
	return _FheOps.Contract.Not(&_FheOps.CallOpts, utype, value)
}

// Or is a free data retrieval call binding the contract method 0x72d456f5.
//
// Solidity: function or(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Or(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "or", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Or is a free data retrieval call binding the contract method 0x72d456f5.
//
// Solidity: function or(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Or(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Or(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Or is a free data retrieval call binding the contract method 0x72d456f5.
//
// Solidity: function or(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Or(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Or(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Rem is a free data retrieval call binding the contract method 0xeb376804.
//
// Solidity: function rem(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Rem(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "rem", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Rem is a free data retrieval call binding the contract method 0xeb376804.
//
// Solidity: function rem(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Rem(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Rem(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Rem is a free data retrieval call binding the contract method 0xeb376804.
//
// Solidity: function rem(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Rem(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Rem(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Req is a free data retrieval call binding the contract method 0x7d23f1db.
//
// Solidity: function req(uint8 utype, bytes input) pure returns(bytes)
func (_FheOps *FheOpsCaller) Req(opts *bind.CallOpts, utype uint8, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "req", utype, input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Req is a free data retrieval call binding the contract method 0x7d23f1db.
//
// Solidity: function req(uint8 utype, bytes input) pure returns(bytes)
func (_FheOps *FheOpsSession) Req(utype uint8, input []byte) ([]byte, error) {
	return _FheOps.Contract.Req(&_FheOps.CallOpts, utype, input)
}

// Req is a free data retrieval call binding the contract method 0x7d23f1db.
//
// Solidity: function req(uint8 utype, bytes input) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Req(utype uint8, input []byte) ([]byte, error) {
	return _FheOps.Contract.Req(&_FheOps.CallOpts, utype, input)
}

// SealOutput is a free data retrieval call binding the contract method 0xa1848ff3.
//
// Solidity: function sealOutput(uint8 utype, bytes ctHash, bytes pk) pure returns(string)
func (_FheOps *FheOpsCaller) SealOutput(opts *bind.CallOpts, utype uint8, ctHash []byte, pk []byte) (string, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "sealOutput", utype, ctHash, pk)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// SealOutput is a free data retrieval call binding the contract method 0xa1848ff3.
//
// Solidity: function sealOutput(uint8 utype, bytes ctHash, bytes pk) pure returns(string)
func (_FheOps *FheOpsSession) SealOutput(utype uint8, ctHash []byte, pk []byte) (string, error) {
	return _FheOps.Contract.SealOutput(&_FheOps.CallOpts, utype, ctHash, pk)
}

// SealOutput is a free data retrieval call binding the contract method 0xa1848ff3.
//
// Solidity: function sealOutput(uint8 utype, bytes ctHash, bytes pk) pure returns(string)
func (_FheOps *FheOpsCallerSession) SealOutput(utype uint8, ctHash []byte, pk []byte) (string, error) {
	return _FheOps.Contract.SealOutput(&_FheOps.CallOpts, utype, ctHash, pk)
}

// Select is a free data retrieval call binding the contract method 0xc2d96952.
//
// Solidity: function select(uint8 utype, bytes controlHash, bytes ifTrueHash, bytes ifFalseHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Select(opts *bind.CallOpts, utype uint8, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "select", utype, controlHash, ifTrueHash, ifFalseHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Select is a free data retrieval call binding the contract method 0xc2d96952.
//
// Solidity: function select(uint8 utype, bytes controlHash, bytes ifTrueHash, bytes ifFalseHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Select(utype uint8, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte) ([]byte, error) {
	return _FheOps.Contract.Select(&_FheOps.CallOpts, utype, controlHash, ifTrueHash, ifFalseHash)
}

// Select is a free data retrieval call binding the contract method 0xc2d96952.
//
// Solidity: function select(uint8 utype, bytes controlHash, bytes ifTrueHash, bytes ifFalseHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Select(utype uint8, controlHash []byte, ifTrueHash []byte, ifFalseHash []byte) ([]byte, error) {
	return _FheOps.Contract.Select(&_FheOps.CallOpts, utype, controlHash, ifTrueHash, ifFalseHash)
}

// Shl is a free data retrieval call binding the contract method 0xae42450a.
//
// Solidity: function shl(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Shl(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "shl", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Shl is a free data retrieval call binding the contract method 0xae42450a.
//
// Solidity: function shl(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Shl(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Shl(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Shl is a free data retrieval call binding the contract method 0xae42450a.
//
// Solidity: function shl(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Shl(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Shl(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Shr is a free data retrieval call binding the contract method 0x9944d12d.
//
// Solidity: function shr(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Shr(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "shr", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Shr is a free data retrieval call binding the contract method 0x9944d12d.
//
// Solidity: function shr(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Shr(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Shr(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Shr is a free data retrieval call binding the contract method 0x9944d12d.
//
// Solidity: function shr(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Shr(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Shr(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Sub is a free data retrieval call binding the contract method 0xcc2cbeff.
//
// Solidity: function sub(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Sub(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "sub", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Sub is a free data retrieval call binding the contract method 0xcc2cbeff.
//
// Solidity: function sub(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Sub(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Sub(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Sub is a free data retrieval call binding the contract method 0xcc2cbeff.
//
// Solidity: function sub(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Sub(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Sub(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// TrivialEncrypt is a free data retrieval call binding the contract method 0x19e1c5c4.
//
// Solidity: function trivialEncrypt(bytes input, uint8 toType) pure returns(bytes)
func (_FheOps *FheOpsCaller) TrivialEncrypt(opts *bind.CallOpts, input []byte, toType uint8) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "trivialEncrypt", input, toType)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// TrivialEncrypt is a free data retrieval call binding the contract method 0x19e1c5c4.
//
// Solidity: function trivialEncrypt(bytes input, uint8 toType) pure returns(bytes)
func (_FheOps *FheOpsSession) TrivialEncrypt(input []byte, toType uint8) ([]byte, error) {
	return _FheOps.Contract.TrivialEncrypt(&_FheOps.CallOpts, input, toType)
}

// TrivialEncrypt is a free data retrieval call binding the contract method 0x19e1c5c4.
//
// Solidity: function trivialEncrypt(bytes input, uint8 toType) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) TrivialEncrypt(input []byte, toType uint8) ([]byte, error) {
	return _FheOps.Contract.TrivialEncrypt(&_FheOps.CallOpts, input, toType)
}

// Verify is a free data retrieval call binding the contract method 0x5fa55ca7.
//
// Solidity: function verify(uint8 utype, bytes input) pure returns(bytes)
func (_FheOps *FheOpsCaller) Verify(opts *bind.CallOpts, utype uint8, input []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "verify", utype, input)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Verify is a free data retrieval call binding the contract method 0x5fa55ca7.
//
// Solidity: function verify(uint8 utype, bytes input) pure returns(bytes)
func (_FheOps *FheOpsSession) Verify(utype uint8, input []byte) ([]byte, error) {
	return _FheOps.Contract.Verify(&_FheOps.CallOpts, utype, input)
}

// Verify is a free data retrieval call binding the contract method 0x5fa55ca7.
//
// Solidity: function verify(uint8 utype, bytes input) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Verify(utype uint8, input []byte) ([]byte, error) {
	return _FheOps.Contract.Verify(&_FheOps.CallOpts, utype, input)
}

// Xor is a free data retrieval call binding the contract method 0x5e639f19.
//
// Solidity: function xor(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCaller) Xor(opts *bind.CallOpts, utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	var out []interface{}
	err := _FheOps.contract.Call(opts, &out, "xor", utype, lhsHash, rhsHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Xor is a free data retrieval call binding the contract method 0x5e639f19.
//
// Solidity: function xor(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsSession) Xor(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Xor(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// Xor is a free data retrieval call binding the contract method 0x5e639f19.
//
// Solidity: function xor(uint8 utype, bytes lhsHash, bytes rhsHash) pure returns(bytes)
func (_FheOps *FheOpsCallerSession) Xor(utype uint8, lhsHash []byte, rhsHash []byte) ([]byte, error) {
	return _FheOps.Contract.Xor(&_FheOps.CallOpts, utype, lhsHash, rhsHash)
}

// PrecompilesMetaData contains all meta data concerning the Precompiles contract.
var PrecompilesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"Fheos\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e3610052600b82828239805160001a607314610045577f4e487b7100000000000000000000000000000000000000000000000000000000600052600060045260246000fd5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361060335760003560e01c806313d8d464146038575b600080fd5b603e6052565b604051604991906094565b60405180910390f35b608081565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006080826057565b9050919050565b608e816077565b82525050565b600060208201905060a760008301846087565b9291505056fea2646970667358221220bd6912a57e1c0144f6338aab2df9f00ce1990b41aa80387ed44cee41f192de9564736f6c63430008130033",
}

// PrecompilesABI is the input ABI used to generate the binding from.
// Deprecated: Use PrecompilesMetaData.ABI instead.
var PrecompilesABI = PrecompilesMetaData.ABI

// PrecompilesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PrecompilesMetaData.Bin instead.
var PrecompilesBin = PrecompilesMetaData.Bin

// DeployPrecompiles deploys a new Ethereum contract, binding an instance of Precompiles to it.
func DeployPrecompiles(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Precompiles, error) {
	parsed, err := PrecompilesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PrecompilesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Precompiles{PrecompilesCaller: PrecompilesCaller{contract: contract}, PrecompilesTransactor: PrecompilesTransactor{contract: contract}, PrecompilesFilterer: PrecompilesFilterer{contract: contract}}, nil
}

// Precompiles is an auto generated Go binding around an Ethereum contract.
type Precompiles struct {
	PrecompilesCaller     // Read-only binding to the contract
	PrecompilesTransactor // Write-only binding to the contract
	PrecompilesFilterer   // Log filterer for contract events
}

// PrecompilesCaller is an auto generated read-only Go binding around an Ethereum contract.
type PrecompilesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrecompilesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PrecompilesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrecompilesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PrecompilesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrecompilesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PrecompilesSession struct {
	Contract     *Precompiles      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PrecompilesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PrecompilesCallerSession struct {
	Contract *PrecompilesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// PrecompilesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PrecompilesTransactorSession struct {
	Contract     *PrecompilesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// PrecompilesRaw is an auto generated low-level Go binding around an Ethereum contract.
type PrecompilesRaw struct {
	Contract *Precompiles // Generic contract binding to access the raw methods on
}

// PrecompilesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PrecompilesCallerRaw struct {
	Contract *PrecompilesCaller // Generic read-only contract binding to access the raw methods on
}

// PrecompilesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PrecompilesTransactorRaw struct {
	Contract *PrecompilesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPrecompiles creates a new instance of Precompiles, bound to a specific deployed contract.
func NewPrecompiles(address common.Address, backend bind.ContractBackend) (*Precompiles, error) {
	contract, err := bindPrecompiles(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Precompiles{PrecompilesCaller: PrecompilesCaller{contract: contract}, PrecompilesTransactor: PrecompilesTransactor{contract: contract}, PrecompilesFilterer: PrecompilesFilterer{contract: contract}}, nil
}

// NewPrecompilesCaller creates a new read-only instance of Precompiles, bound to a specific deployed contract.
func NewPrecompilesCaller(address common.Address, caller bind.ContractCaller) (*PrecompilesCaller, error) {
	contract, err := bindPrecompiles(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PrecompilesCaller{contract: contract}, nil
}

// NewPrecompilesTransactor creates a new write-only instance of Precompiles, bound to a specific deployed contract.
func NewPrecompilesTransactor(address common.Address, transactor bind.ContractTransactor) (*PrecompilesTransactor, error) {
	contract, err := bindPrecompiles(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PrecompilesTransactor{contract: contract}, nil
}

// NewPrecompilesFilterer creates a new log filterer instance of Precompiles, bound to a specific deployed contract.
func NewPrecompilesFilterer(address common.Address, filterer bind.ContractFilterer) (*PrecompilesFilterer, error) {
	contract, err := bindPrecompiles(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PrecompilesFilterer{contract: contract}, nil
}

// bindPrecompiles binds a generic wrapper to an already deployed contract.
func bindPrecompiles(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PrecompilesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Precompiles *PrecompilesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Precompiles.Contract.PrecompilesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Precompiles *PrecompilesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Precompiles.Contract.PrecompilesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Precompiles *PrecompilesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Precompiles.Contract.PrecompilesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Precompiles *PrecompilesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Precompiles.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Precompiles *PrecompilesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Precompiles.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Precompiles *PrecompilesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Precompiles.Contract.contract.Transact(opts, method, params...)
}

// Fheos is a free data retrieval call binding the contract method 0x13d8d464.
//
// Solidity: function Fheos() view returns(address)
func (_Precompiles *PrecompilesCaller) Fheos(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Precompiles.contract.Call(opts, &out, "Fheos")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Fheos is a free data retrieval call binding the contract method 0x13d8d464.
//
// Solidity: function Fheos() view returns(address)
func (_Precompiles *PrecompilesSession) Fheos() (common.Address, error) {
	return _Precompiles.Contract.Fheos(&_Precompiles.CallOpts)
}

// Fheos is a free data retrieval call binding the contract method 0x13d8d464.
//
// Solidity: function Fheos() view returns(address)
func (_Precompiles *PrecompilesCallerSession) Fheos() (common.Address, error) {
	return _Precompiles.Contract.Fheos(&_Precompiles.CallOpts)
}
