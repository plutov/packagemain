// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TwoDABI is the input ABI used to generate the binding from.
const TwoDABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"x\",\"type\":\"uint8\"},{\"name\":\"y\",\"type\":\"uint8\"}],\"name\":\"getValue\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"kill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// TwoDBin is the compiled bytecode used for deploying new contracts.
const TwoDBin = `0x60606040526000805460a060020a60ff021916740a0000000000000000000000000000000000000000179055341561003657600080fd5b60008054600160a060020a03191633600160a060020a0316178155805b60005460ff740100000000000000000000000000000000000000009091048116908216101561010157600091505b60005460ff74010000000000000000000000000000000000000000909104811690831610156100f957818101600160ff8416600a81106100bd57fe5b0160ff8316600a81106100cc57fe5b602091828204019190066101000a81548160ff021916908360ff1602179055508180600101925050610081565b600101610053565b5050610145806101126000396000f30060606040526004361061004b5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630e13b9af811461005057806341c0e1b514610085575b600080fd5b341561005b57600080fd5b61006f60ff6004358116906024351661009a565b60405160ff909116815260200160405180910390f35b341561009057600080fd5b6100986100d8565b005b6000600160ff8416600a81106100ac57fe5b0160ff8316600a81106100bb57fe5b60208082049092015460ff929091066101000a9004169392505050565b6000543373ffffffffffffffffffffffffffffffffffffffff908116911614156101175760005473ffffffffffffffffffffffffffffffffffffffff16ff5b5600a165627a7a723058203cacbb8b93ce91ee0f52d315aa038336420ab25e0befcbb8db69749d165a41ee0029`

// DeployTwoD deploys a new Ethereum contract, binding an instance of TwoD to it.
func DeployTwoD(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TwoD, error) {
	parsed, err := abi.JSON(strings.NewReader(TwoDABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TwoDBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TwoD{TwoDCaller: TwoDCaller{contract: contract}, TwoDTransactor: TwoDTransactor{contract: contract}}, nil
}

// TwoD is an auto generated Go binding around an Ethereum contract.
type TwoD struct {
	TwoDCaller     // Read-only binding to the contract
	TwoDTransactor // Write-only binding to the contract
}

// TwoDCaller is an auto generated read-only Go binding around an Ethereum contract.
type TwoDCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoDTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TwoDTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoDSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TwoDSession struct {
	Contract     *TwoD             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TwoDCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TwoDCallerSession struct {
	Contract *TwoDCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TwoDTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TwoDTransactorSession struct {
	Contract     *TwoDTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TwoDRaw is an auto generated low-level Go binding around an Ethereum contract.
type TwoDRaw struct {
	Contract *TwoD // Generic contract binding to access the raw methods on
}

// TwoDCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TwoDCallerRaw struct {
	Contract *TwoDCaller // Generic read-only contract binding to access the raw methods on
}

// TwoDTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TwoDTransactorRaw struct {
	Contract *TwoDTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTwoD creates a new instance of TwoD, bound to a specific deployed contract.
func NewTwoD(address common.Address, backend bind.ContractBackend) (*TwoD, error) {
	contract, err := bindTwoD(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TwoD{TwoDCaller: TwoDCaller{contract: contract}, TwoDTransactor: TwoDTransactor{contract: contract}}, nil
}

// NewTwoDCaller creates a new read-only instance of TwoD, bound to a specific deployed contract.
func NewTwoDCaller(address common.Address, caller bind.ContractCaller) (*TwoDCaller, error) {
	contract, err := bindTwoD(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &TwoDCaller{contract: contract}, nil
}

// NewTwoDTransactor creates a new write-only instance of TwoD, bound to a specific deployed contract.
func NewTwoDTransactor(address common.Address, transactor bind.ContractTransactor) (*TwoDTransactor, error) {
	contract, err := bindTwoD(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &TwoDTransactor{contract: contract}, nil
}

// bindTwoD binds a generic wrapper to an already deployed contract.
func bindTwoD(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TwoDABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TwoD *TwoDRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TwoD.Contract.TwoDCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TwoD *TwoDRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TwoD.Contract.TwoDTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TwoD *TwoDRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TwoD.Contract.TwoDTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TwoD *TwoDCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TwoD.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TwoD *TwoDTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TwoD.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TwoD *TwoDTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TwoD.Contract.contract.Transact(opts, method, params...)
}

// GetValue is a free data retrieval call binding the contract method 0x0e13b9af.
//
// Solidity: function getValue(x uint8, y uint8) constant returns(uint8)
func (_TwoD *TwoDCaller) GetValue(opts *bind.CallOpts, x uint8, y uint8) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _TwoD.contract.Call(opts, out, "getValue", x, y)
	return *ret0, err
}

// GetValue is a free data retrieval call binding the contract method 0x0e13b9af.
//
// Solidity: function getValue(x uint8, y uint8) constant returns(uint8)
func (_TwoD *TwoDSession) GetValue(x uint8, y uint8) (uint8, error) {
	return _TwoD.Contract.GetValue(&_TwoD.CallOpts, x, y)
}

// GetValue is a free data retrieval call binding the contract method 0x0e13b9af.
//
// Solidity: function getValue(x uint8, y uint8) constant returns(uint8)
func (_TwoD *TwoDCallerSession) GetValue(x uint8, y uint8) (uint8, error) {
	return _TwoD.Contract.GetValue(&_TwoD.CallOpts, x, y)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_TwoD *TwoDTransactor) Kill(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TwoD.contract.Transact(opts, "kill")
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_TwoD *TwoDSession) Kill() (*types.Transaction, error) {
	return _TwoD.Contract.Kill(&_TwoD.TransactOpts)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_TwoD *TwoDTransactorSession) Kill() (*types.Transaction, error) {
	return _TwoD.Contract.Kill(&_TwoD.TransactOpts)
}
