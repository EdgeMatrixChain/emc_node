// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package helper

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

// StakeMetaData contains all meta data concerning the Stake contract.
var StakeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_releaseContract\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"_canDeposit\",\"type\":\"bool\"}],\"name\":\"CanDepositUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"holder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"nodeId\",\"type\":\"string\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"holder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"nodeId\",\"type\":\"string\"}],\"name\":\"Withdrawed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_beneficiary\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_nodeId\",\"type\":\"string\"}],\"name\":\"balanceOfNode\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"canDeposit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_nodeId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"manager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"nodeInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"accumulated\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"debt\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"releaseContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_canDeposit\",\"type\":\"bool\"}],\"name\":\"setCanDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxLimit\",\"type\":\"uint256\"}],\"name\":\"setLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_manager\",\"type\":\"address\"}],\"name\":\"setManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenInPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_nodeId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_beneficiary\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StakeABI is the input ABI used to generate the binding from.
// Deprecated: Use StakeMetaData.ABI instead.
var StakeABI = StakeMetaData.ABI

// Stake is an auto generated Go binding around an Ethereum contract.
type Stake struct {
	StakeCaller     // Read-only binding to the contract
	StakeTransactor // Write-only binding to the contract
	StakeFilterer   // Log filterer for contract events
}

// StakeCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakeSession struct {
	Contract     *Stake            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakeCallerSession struct {
	Contract *StakeCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// StakeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakeTransactorSession struct {
	Contract     *StakeTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakeRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakeRaw struct {
	Contract *Stake // Generic contract binding to access the raw methods on
}

// StakeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakeCallerRaw struct {
	Contract *StakeCaller // Generic read-only contract binding to access the raw methods on
}

// StakeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakeTransactorRaw struct {
	Contract *StakeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStake creates a new instance of Stake, bound to a specific deployed contract.
func NewStake(address common.Address, backend bind.ContractBackend) (*Stake, error) {
	contract, err := bindStake(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Stake{StakeCaller: StakeCaller{contract: contract}, StakeTransactor: StakeTransactor{contract: contract}, StakeFilterer: StakeFilterer{contract: contract}}, nil
}

// NewStakeCaller creates a new read-only instance of Stake, bound to a specific deployed contract.
func NewStakeCaller(address common.Address, caller bind.ContractCaller) (*StakeCaller, error) {
	contract, err := bindStake(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakeCaller{contract: contract}, nil
}

// NewStakeTransactor creates a new write-only instance of Stake, bound to a specific deployed contract.
func NewStakeTransactor(address common.Address, transactor bind.ContractTransactor) (*StakeTransactor, error) {
	contract, err := bindStake(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakeTransactor{contract: contract}, nil
}

// NewStakeFilterer creates a new log filterer instance of Stake, bound to a specific deployed contract.
func NewStakeFilterer(address common.Address, filterer bind.ContractFilterer) (*StakeFilterer, error) {
	contract, err := bindStake(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakeFilterer{contract: contract}, nil
}

// bindStake binds a generic wrapper to an already deployed contract.
func bindStake(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stake *StakeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Stake.Contract.StakeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stake *StakeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stake.Contract.StakeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stake *StakeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stake.Contract.StakeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stake *StakeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Stake.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stake *StakeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stake.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stake *StakeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stake.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _beneficiary) view returns(uint256)
func (_Stake *StakeCaller) BalanceOf(opts *bind.CallOpts, _beneficiary common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "balanceOf", _beneficiary)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _beneficiary) view returns(uint256)
func (_Stake *StakeSession) BalanceOf(_beneficiary common.Address) (*big.Int, error) {
	return _Stake.Contract.BalanceOf(&_Stake.CallOpts, _beneficiary)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _beneficiary) view returns(uint256)
func (_Stake *StakeCallerSession) BalanceOf(_beneficiary common.Address) (*big.Int, error) {
	return _Stake.Contract.BalanceOf(&_Stake.CallOpts, _beneficiary)
}

// BalanceOfNode is a free data retrieval call binding the contract method 0x4f6d810a.
//
// Solidity: function balanceOfNode(string _nodeId) view returns(uint256)
func (_Stake *StakeCaller) BalanceOfNode(opts *bind.CallOpts, _nodeId string) (*big.Int, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "balanceOfNode", _nodeId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOfNode is a free data retrieval call binding the contract method 0x4f6d810a.
//
// Solidity: function balanceOfNode(string _nodeId) view returns(uint256)
func (_Stake *StakeSession) BalanceOfNode(_nodeId string) (*big.Int, error) {
	return _Stake.Contract.BalanceOfNode(&_Stake.CallOpts, _nodeId)
}

// BalanceOfNode is a free data retrieval call binding the contract method 0x4f6d810a.
//
// Solidity: function balanceOfNode(string _nodeId) view returns(uint256)
func (_Stake *StakeCallerSession) BalanceOfNode(_nodeId string) (*big.Int, error) {
	return _Stake.Contract.BalanceOfNode(&_Stake.CallOpts, _nodeId)
}

// CanDeposit is a free data retrieval call binding the contract method 0xe78a5875.
//
// Solidity: function canDeposit() view returns(bool)
func (_Stake *StakeCaller) CanDeposit(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "canDeposit")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanDeposit is a free data retrieval call binding the contract method 0xe78a5875.
//
// Solidity: function canDeposit() view returns(bool)
func (_Stake *StakeSession) CanDeposit() (bool, error) {
	return _Stake.Contract.CanDeposit(&_Stake.CallOpts)
}

// CanDeposit is a free data retrieval call binding the contract method 0xe78a5875.
//
// Solidity: function canDeposit() view returns(bool)
func (_Stake *StakeCallerSession) CanDeposit() (bool, error) {
	return _Stake.Contract.CanDeposit(&_Stake.CallOpts)
}

// Manager is a free data retrieval call binding the contract method 0x481c6a75.
//
// Solidity: function manager() view returns(address)
func (_Stake *StakeCaller) Manager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "manager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Manager is a free data retrieval call binding the contract method 0x481c6a75.
//
// Solidity: function manager() view returns(address)
func (_Stake *StakeSession) Manager() (common.Address, error) {
	return _Stake.Contract.Manager(&_Stake.CallOpts)
}

// Manager is a free data retrieval call binding the contract method 0x481c6a75.
//
// Solidity: function manager() view returns(address)
func (_Stake *StakeCallerSession) Manager() (common.Address, error) {
	return _Stake.Contract.Manager(&_Stake.CallOpts)
}

// MaxLimit is a free data retrieval call binding the contract method 0x1a861d26.
//
// Solidity: function maxLimit() view returns(uint256)
func (_Stake *StakeCaller) MaxLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "maxLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxLimit is a free data retrieval call binding the contract method 0x1a861d26.
//
// Solidity: function maxLimit() view returns(uint256)
func (_Stake *StakeSession) MaxLimit() (*big.Int, error) {
	return _Stake.Contract.MaxLimit(&_Stake.CallOpts)
}

// MaxLimit is a free data retrieval call binding the contract method 0x1a861d26.
//
// Solidity: function maxLimit() view returns(uint256)
func (_Stake *StakeCallerSession) MaxLimit() (*big.Int, error) {
	return _Stake.Contract.MaxLimit(&_Stake.CallOpts)
}

// MinLimit is a free data retrieval call binding the contract method 0x1fd8088d.
//
// Solidity: function minLimit() view returns(uint256)
func (_Stake *StakeCaller) MinLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "minLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinLimit is a free data retrieval call binding the contract method 0x1fd8088d.
//
// Solidity: function minLimit() view returns(uint256)
func (_Stake *StakeSession) MinLimit() (*big.Int, error) {
	return _Stake.Contract.MinLimit(&_Stake.CallOpts)
}

// MinLimit is a free data retrieval call binding the contract method 0x1fd8088d.
//
// Solidity: function minLimit() view returns(uint256)
func (_Stake *StakeCallerSession) MinLimit() (*big.Int, error) {
	return _Stake.Contract.MinLimit(&_Stake.CallOpts)
}

// NodeInfo is a free data retrieval call binding the contract method 0xe8a0c74e.
//
// Solidity: function nodeInfo(string ) view returns(address beneficiary, uint256 accumulated, uint256 amount, uint256 debt)
func (_Stake *StakeCaller) NodeInfo(opts *bind.CallOpts, arg0 string) (struct {
	Beneficiary common.Address
	Accumulated *big.Int
	Amount      *big.Int
	Debt        *big.Int
}, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "nodeInfo", arg0)

	outstruct := new(struct {
		Beneficiary common.Address
		Accumulated *big.Int
		Amount      *big.Int
		Debt        *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Beneficiary = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Accumulated = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Amount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Debt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// NodeInfo is a free data retrieval call binding the contract method 0xe8a0c74e.
//
// Solidity: function nodeInfo(string ) view returns(address beneficiary, uint256 accumulated, uint256 amount, uint256 debt)
func (_Stake *StakeSession) NodeInfo(arg0 string) (struct {
	Beneficiary common.Address
	Accumulated *big.Int
	Amount      *big.Int
	Debt        *big.Int
}, error) {
	return _Stake.Contract.NodeInfo(&_Stake.CallOpts, arg0)
}

// NodeInfo is a free data retrieval call binding the contract method 0xe8a0c74e.
//
// Solidity: function nodeInfo(string ) view returns(address beneficiary, uint256 accumulated, uint256 amount, uint256 debt)
func (_Stake *StakeCallerSession) NodeInfo(arg0 string) (struct {
	Beneficiary common.Address
	Accumulated *big.Int
	Amount      *big.Int
	Debt        *big.Int
}, error) {
	return _Stake.Contract.NodeInfo(&_Stake.CallOpts, arg0)
}

// ReleaseContract is a free data retrieval call binding the contract method 0xc71ddb9f.
//
// Solidity: function releaseContract() view returns(address)
func (_Stake *StakeCaller) ReleaseContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "releaseContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ReleaseContract is a free data retrieval call binding the contract method 0xc71ddb9f.
//
// Solidity: function releaseContract() view returns(address)
func (_Stake *StakeSession) ReleaseContract() (common.Address, error) {
	return _Stake.Contract.ReleaseContract(&_Stake.CallOpts)
}

// ReleaseContract is a free data retrieval call binding the contract method 0xc71ddb9f.
//
// Solidity: function releaseContract() view returns(address)
func (_Stake *StakeCallerSession) ReleaseContract() (common.Address, error) {
	return _Stake.Contract.ReleaseContract(&_Stake.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Stake *StakeCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Stake *StakeSession) Token() (common.Address, error) {
	return _Stake.Contract.Token(&_Stake.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Stake *StakeCallerSession) Token() (common.Address, error) {
	return _Stake.Contract.Token(&_Stake.CallOpts)
}

// TokenInPool is a free data retrieval call binding the contract method 0xb5e5be5f.
//
// Solidity: function tokenInPool() view returns(uint256)
func (_Stake *StakeCaller) TokenInPool(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "tokenInPool")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenInPool is a free data retrieval call binding the contract method 0xb5e5be5f.
//
// Solidity: function tokenInPool() view returns(uint256)
func (_Stake *StakeSession) TokenInPool() (*big.Int, error) {
	return _Stake.Contract.TokenInPool(&_Stake.CallOpts)
}

// TokenInPool is a free data retrieval call binding the contract method 0xb5e5be5f.
//
// Solidity: function tokenInPool() view returns(uint256)
func (_Stake *StakeCallerSession) TokenInPool() (*big.Int, error) {
	return _Stake.Contract.TokenInPool(&_Stake.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x8e27d719.
//
// Solidity: function deposit(string _nodeId, uint256 _amount) returns()
func (_Stake *StakeTransactor) Deposit(opts *bind.TransactOpts, _nodeId string, _amount *big.Int) (*types.Transaction, error) {
	return _Stake.contract.Transact(opts, "deposit", _nodeId, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x8e27d719.
//
// Solidity: function deposit(string _nodeId, uint256 _amount) returns()
func (_Stake *StakeSession) Deposit(_nodeId string, _amount *big.Int) (*types.Transaction, error) {
	return _Stake.Contract.Deposit(&_Stake.TransactOpts, _nodeId, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x8e27d719.
//
// Solidity: function deposit(string _nodeId, uint256 _amount) returns()
func (_Stake *StakeTransactorSession) Deposit(_nodeId string, _amount *big.Int) (*types.Transaction, error) {
	return _Stake.Contract.Deposit(&_Stake.TransactOpts, _nodeId, _amount)
}

// SetCanDeposit is a paid mutator transaction binding the contract method 0x761b5d91.
//
// Solidity: function setCanDeposit(bool _canDeposit) returns()
func (_Stake *StakeTransactor) SetCanDeposit(opts *bind.TransactOpts, _canDeposit bool) (*types.Transaction, error) {
	return _Stake.contract.Transact(opts, "setCanDeposit", _canDeposit)
}

// SetCanDeposit is a paid mutator transaction binding the contract method 0x761b5d91.
//
// Solidity: function setCanDeposit(bool _canDeposit) returns()
func (_Stake *StakeSession) SetCanDeposit(_canDeposit bool) (*types.Transaction, error) {
	return _Stake.Contract.SetCanDeposit(&_Stake.TransactOpts, _canDeposit)
}

// SetCanDeposit is a paid mutator transaction binding the contract method 0x761b5d91.
//
// Solidity: function setCanDeposit(bool _canDeposit) returns()
func (_Stake *StakeTransactorSession) SetCanDeposit(_canDeposit bool) (*types.Transaction, error) {
	return _Stake.Contract.SetCanDeposit(&_Stake.TransactOpts, _canDeposit)
}

// SetLimit is a paid mutator transaction binding the contract method 0x207add91.
//
// Solidity: function setLimit(uint256 _minLimit, uint256 _maxLimit) returns()
func (_Stake *StakeTransactor) SetLimit(opts *bind.TransactOpts, _minLimit *big.Int, _maxLimit *big.Int) (*types.Transaction, error) {
	return _Stake.contract.Transact(opts, "setLimit", _minLimit, _maxLimit)
}

// SetLimit is a paid mutator transaction binding the contract method 0x207add91.
//
// Solidity: function setLimit(uint256 _minLimit, uint256 _maxLimit) returns()
func (_Stake *StakeSession) SetLimit(_minLimit *big.Int, _maxLimit *big.Int) (*types.Transaction, error) {
	return _Stake.Contract.SetLimit(&_Stake.TransactOpts, _minLimit, _maxLimit)
}

// SetLimit is a paid mutator transaction binding the contract method 0x207add91.
//
// Solidity: function setLimit(uint256 _minLimit, uint256 _maxLimit) returns()
func (_Stake *StakeTransactorSession) SetLimit(_minLimit *big.Int, _maxLimit *big.Int) (*types.Transaction, error) {
	return _Stake.Contract.SetLimit(&_Stake.TransactOpts, _minLimit, _maxLimit)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address _manager) returns()
func (_Stake *StakeTransactor) SetManager(opts *bind.TransactOpts, _manager common.Address) (*types.Transaction, error) {
	return _Stake.contract.Transact(opts, "setManager", _manager)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address _manager) returns()
func (_Stake *StakeSession) SetManager(_manager common.Address) (*types.Transaction, error) {
	return _Stake.Contract.SetManager(&_Stake.TransactOpts, _manager)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address _manager) returns()
func (_Stake *StakeTransactorSession) SetManager(_manager common.Address) (*types.Transaction, error) {
	return _Stake.Contract.SetManager(&_Stake.TransactOpts, _manager)
}

// Withdraw is a paid mutator transaction binding the contract method 0x5a73b0bf.
//
// Solidity: function withdraw(string _nodeId, address _beneficiary, uint256 _amount) returns()
func (_Stake *StakeTransactor) Withdraw(opts *bind.TransactOpts, _nodeId string, _beneficiary common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Stake.contract.Transact(opts, "withdraw", _nodeId, _beneficiary, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x5a73b0bf.
//
// Solidity: function withdraw(string _nodeId, address _beneficiary, uint256 _amount) returns()
func (_Stake *StakeSession) Withdraw(_nodeId string, _beneficiary common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Stake.Contract.Withdraw(&_Stake.TransactOpts, _nodeId, _beneficiary, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x5a73b0bf.
//
// Solidity: function withdraw(string _nodeId, address _beneficiary, uint256 _amount) returns()
func (_Stake *StakeTransactorSession) Withdraw(_nodeId string, _beneficiary common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Stake.Contract.Withdraw(&_Stake.TransactOpts, _nodeId, _beneficiary, _amount)
}

// StakeCanDepositUpdatedIterator is returned from FilterCanDepositUpdated and is used to iterate over the raw logs and unpacked data for CanDepositUpdated events raised by the Stake contract.
type StakeCanDepositUpdatedIterator struct {
	Event *StakeCanDepositUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakeCanDepositUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeCanDepositUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakeCanDepositUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakeCanDepositUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeCanDepositUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeCanDepositUpdated represents a CanDepositUpdated event raised by the Stake contract.
type StakeCanDepositUpdated struct {
	Operator   common.Address
	CanDeposit bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCanDepositUpdated is a free log retrieval operation binding the contract event 0x45abe80e1ccfdfb2be4680ea7a5d3a3a210f24b5e10dc89ea8c4eed1506e6b47.
//
// Solidity: event CanDepositUpdated(address _operator, bool _canDeposit)
func (_Stake *StakeFilterer) FilterCanDepositUpdated(opts *bind.FilterOpts) (*StakeCanDepositUpdatedIterator, error) {

	logs, sub, err := _Stake.contract.FilterLogs(opts, "CanDepositUpdated")
	if err != nil {
		return nil, err
	}
	return &StakeCanDepositUpdatedIterator{contract: _Stake.contract, event: "CanDepositUpdated", logs: logs, sub: sub}, nil
}

// WatchCanDepositUpdated is a free log subscription operation binding the contract event 0x45abe80e1ccfdfb2be4680ea7a5d3a3a210f24b5e10dc89ea8c4eed1506e6b47.
//
// Solidity: event CanDepositUpdated(address _operator, bool _canDeposit)
func (_Stake *StakeFilterer) WatchCanDepositUpdated(opts *bind.WatchOpts, sink chan<- *StakeCanDepositUpdated) (event.Subscription, error) {

	logs, sub, err := _Stake.contract.WatchLogs(opts, "CanDepositUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeCanDepositUpdated)
				if err := _Stake.contract.UnpackLog(event, "CanDepositUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCanDepositUpdated is a log parse operation binding the contract event 0x45abe80e1ccfdfb2be4680ea7a5d3a3a210f24b5e10dc89ea8c4eed1506e6b47.
//
// Solidity: event CanDepositUpdated(address _operator, bool _canDeposit)
func (_Stake *StakeFilterer) ParseCanDepositUpdated(log types.Log) (*StakeCanDepositUpdated, error) {
	event := new(StakeCanDepositUpdated)
	if err := _Stake.contract.UnpackLog(event, "CanDepositUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the Stake contract.
type StakeDepositedIterator struct {
	Event *StakeDeposited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakeDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeDeposited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakeDeposited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakeDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeDeposited represents a Deposited event raised by the Stake contract.
type StakeDeposited struct {
	Holder common.Address
	Amount *big.Int
	NodeId string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x6f85d9948d6ca3dd6ce6ce7d175da22b4e865827ae6fcd530ec7edac1240f928.
//
// Solidity: event Deposited(address holder, uint256 amount, string nodeId)
func (_Stake *StakeFilterer) FilterDeposited(opts *bind.FilterOpts) (*StakeDepositedIterator, error) {

	logs, sub, err := _Stake.contract.FilterLogs(opts, "Deposited")
	if err != nil {
		return nil, err
	}
	return &StakeDepositedIterator{contract: _Stake.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x6f85d9948d6ca3dd6ce6ce7d175da22b4e865827ae6fcd530ec7edac1240f928.
//
// Solidity: event Deposited(address holder, uint256 amount, string nodeId)
func (_Stake *StakeFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *StakeDeposited) (event.Subscription, error) {

	logs, sub, err := _Stake.contract.WatchLogs(opts, "Deposited")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeDeposited)
				if err := _Stake.contract.UnpackLog(event, "Deposited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeposited is a log parse operation binding the contract event 0x6f85d9948d6ca3dd6ce6ce7d175da22b4e865827ae6fcd530ec7edac1240f928.
//
// Solidity: event Deposited(address holder, uint256 amount, string nodeId)
func (_Stake *StakeFilterer) ParseDeposited(log types.Log) (*StakeDeposited, error) {
	event := new(StakeDeposited)
	if err := _Stake.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeWithdrawedIterator is returned from FilterWithdrawed and is used to iterate over the raw logs and unpacked data for Withdrawed events raised by the Stake contract.
type StakeWithdrawedIterator struct {
	Event *StakeWithdrawed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakeWithdrawedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeWithdrawed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakeWithdrawed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakeWithdrawedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeWithdrawedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeWithdrawed represents a Withdrawed event raised by the Stake contract.
type StakeWithdrawed struct {
	Holder common.Address
	Amount *big.Int
	NodeId string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawed is a free log retrieval operation binding the contract event 0x65a13b133ddbf1d2579eb487d69d87473b960f2ae33c9a9cad53b6075f27a1c4.
//
// Solidity: event Withdrawed(address holder, uint256 amount, string nodeId)
func (_Stake *StakeFilterer) FilterWithdrawed(opts *bind.FilterOpts) (*StakeWithdrawedIterator, error) {

	logs, sub, err := _Stake.contract.FilterLogs(opts, "Withdrawed")
	if err != nil {
		return nil, err
	}
	return &StakeWithdrawedIterator{contract: _Stake.contract, event: "Withdrawed", logs: logs, sub: sub}, nil
}

// WatchWithdrawed is a free log subscription operation binding the contract event 0x65a13b133ddbf1d2579eb487d69d87473b960f2ae33c9a9cad53b6075f27a1c4.
//
// Solidity: event Withdrawed(address holder, uint256 amount, string nodeId)
func (_Stake *StakeFilterer) WatchWithdrawed(opts *bind.WatchOpts, sink chan<- *StakeWithdrawed) (event.Subscription, error) {

	logs, sub, err := _Stake.contract.WatchLogs(opts, "Withdrawed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeWithdrawed)
				if err := _Stake.contract.UnpackLog(event, "Withdrawed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawed is a log parse operation binding the contract event 0x65a13b133ddbf1d2579eb487d69d87473b960f2ae33c9a9cad53b6075f27a1c4.
//
// Solidity: event Withdrawed(address holder, uint256 amount, string nodeId)
func (_Stake *StakeFilterer) ParseWithdrawed(log types.Log) (*StakeWithdrawed, error) {
	event := new(StakeWithdrawed)
	if err := _Stake.contract.UnpackLog(event, "Withdrawed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
