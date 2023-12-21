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

// RewardVestingV1VestingSchedule is an auto generated low-level Go binding around an user-defined struct.
type RewardVestingV1VestingSchedule struct {
	Beneficiary   common.Address
	Start         *big.Int
	Duration      *big.Int
	DurationUnits uint8
	AmountTotal   *big.Int
	Released      *big.Int
	YieldRate     *big.Int
	Rewarded      *big.Int
}

// StakeMetaData contains all meta data concerning the Stake contract.
var StakeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"days30BaseRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"days90BaseRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"days180BaseRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"days360BaseRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"days720BaseRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"days1080BaseRate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"TokensReleased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumRewardVestingV1.DurationUnits\",\"name\":\"durationUnits\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountTotal\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"yieldRate\",\"type\":\"uint256\"}],\"name\":\"VestingScheduleCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_beneficiary\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_start\",\"type\":\"uint256\"},{\"internalType\":\"enumRewardVestingV1.DurationUnits\",\"name\":\"_durationUnits\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_amountTotal\",\"type\":\"uint256\"}],\"name\":\"createVestingSchedule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"depositPermanently\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumRewardVestingV1.DurationUnits\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"durationUnitRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_beneficiary\",\"type\":\"address\"}],\"name\":\"getAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDurationUnitRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_beneficiary\",\"type\":\"address\"}],\"name\":\"getLockedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_beneficiary\",\"type\":\"address\"}],\"name\":\"getReleasableAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_beneficiary\",\"type\":\"address\"}],\"name\":\"getVestingSchedule\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"enumRewardVestingV1.DurationUnits\",\"name\":\"durationUnits\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"amountTotal\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"released\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"yieldRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rewarded\",\"type\":\"uint256\"}],\"internalType\":\"structRewardVestingV1.VestingSchedule[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"permanentTotal\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_beneficiary\",\"type\":\"address\"}],\"name\":\"release\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"vestingSchedules\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"enumRewardVestingV1.DurationUnits\",\"name\":\"durationUnits\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"amountTotal\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"released\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"yieldRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rewarded\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Stake *StakeCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Stake *StakeSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Stake.Contract.BalanceOf(&_Stake.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Stake *StakeCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Stake.Contract.BalanceOf(&_Stake.CallOpts, account)
}

// DurationUnitRewards is a free data retrieval call binding the contract method 0x55310264.
//
// Solidity: function durationUnitRewards(uint8 ) view returns(uint256)
func (_Stake *StakeCaller) DurationUnitRewards(opts *bind.CallOpts, arg0 uint8) (*big.Int, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "durationUnitRewards", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DurationUnitRewards is a free data retrieval call binding the contract method 0x55310264.
//
// Solidity: function durationUnitRewards(uint8 ) view returns(uint256)
func (_Stake *StakeSession) DurationUnitRewards(arg0 uint8) (*big.Int, error) {
	return _Stake.Contract.DurationUnitRewards(&_Stake.CallOpts, arg0)
}

// DurationUnitRewards is a free data retrieval call binding the contract method 0x55310264.
//
// Solidity: function durationUnitRewards(uint8 ) view returns(uint256)
func (_Stake *StakeCallerSession) DurationUnitRewards(arg0 uint8) (*big.Int, error) {
	return _Stake.Contract.DurationUnitRewards(&_Stake.CallOpts, arg0)
}

// GetAmount is a free data retrieval call binding the contract method 0xf5a79767.
//
// Solidity: function getAmount(address _beneficiary) view returns(uint256, uint256, uint256)
func (_Stake *StakeCaller) GetAmount(opts *bind.CallOpts, _beneficiary common.Address) (*big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "getAmount", _beneficiary)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// GetAmount is a free data retrieval call binding the contract method 0xf5a79767.
//
// Solidity: function getAmount(address _beneficiary) view returns(uint256, uint256, uint256)
func (_Stake *StakeSession) GetAmount(_beneficiary common.Address) (*big.Int, *big.Int, *big.Int, error) {
	return _Stake.Contract.GetAmount(&_Stake.CallOpts, _beneficiary)
}

// GetAmount is a free data retrieval call binding the contract method 0xf5a79767.
//
// Solidity: function getAmount(address _beneficiary) view returns(uint256, uint256, uint256)
func (_Stake *StakeCallerSession) GetAmount(_beneficiary common.Address) (*big.Int, *big.Int, *big.Int, error) {
	return _Stake.Contract.GetAmount(&_Stake.CallOpts, _beneficiary)
}

// GetDurationUnitRewards is a free data retrieval call binding the contract method 0x37dc428c.
//
// Solidity: function getDurationUnitRewards() view returns(uint256, uint256, uint256, uint256, uint256, uint256)
func (_Stake *StakeCaller) GetDurationUnitRewards(opts *bind.CallOpts) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "getDurationUnitRewards")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	out4 := *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	out5 := *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return out0, out1, out2, out3, out4, out5, err

}

// GetDurationUnitRewards is a free data retrieval call binding the contract method 0x37dc428c.
//
// Solidity: function getDurationUnitRewards() view returns(uint256, uint256, uint256, uint256, uint256, uint256)
func (_Stake *StakeSession) GetDurationUnitRewards() (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _Stake.Contract.GetDurationUnitRewards(&_Stake.CallOpts)
}

// GetDurationUnitRewards is a free data retrieval call binding the contract method 0x37dc428c.
//
// Solidity: function getDurationUnitRewards() view returns(uint256, uint256, uint256, uint256, uint256, uint256)
func (_Stake *StakeCallerSession) GetDurationUnitRewards() (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _Stake.Contract.GetDurationUnitRewards(&_Stake.CallOpts)
}

// GetLockedAmount is a free data retrieval call binding the contract method 0x929ec537.
//
// Solidity: function getLockedAmount(address _beneficiary) view returns(uint256)
func (_Stake *StakeCaller) GetLockedAmount(opts *bind.CallOpts, _beneficiary common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "getLockedAmount", _beneficiary)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockedAmount is a free data retrieval call binding the contract method 0x929ec537.
//
// Solidity: function getLockedAmount(address _beneficiary) view returns(uint256)
func (_Stake *StakeSession) GetLockedAmount(_beneficiary common.Address) (*big.Int, error) {
	return _Stake.Contract.GetLockedAmount(&_Stake.CallOpts, _beneficiary)
}

// GetLockedAmount is a free data retrieval call binding the contract method 0x929ec537.
//
// Solidity: function getLockedAmount(address _beneficiary) view returns(uint256)
func (_Stake *StakeCallerSession) GetLockedAmount(_beneficiary common.Address) (*big.Int, error) {
	return _Stake.Contract.GetLockedAmount(&_Stake.CallOpts, _beneficiary)
}

// GetReleasableAmount is a free data retrieval call binding the contract method 0x2afd1a7d.
//
// Solidity: function getReleasableAmount(address _beneficiary) view returns(uint256, uint256)
func (_Stake *StakeCaller) GetReleasableAmount(opts *bind.CallOpts, _beneficiary common.Address) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "getReleasableAmount", _beneficiary)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetReleasableAmount is a free data retrieval call binding the contract method 0x2afd1a7d.
//
// Solidity: function getReleasableAmount(address _beneficiary) view returns(uint256, uint256)
func (_Stake *StakeSession) GetReleasableAmount(_beneficiary common.Address) (*big.Int, *big.Int, error) {
	return _Stake.Contract.GetReleasableAmount(&_Stake.CallOpts, _beneficiary)
}

// GetReleasableAmount is a free data retrieval call binding the contract method 0x2afd1a7d.
//
// Solidity: function getReleasableAmount(address _beneficiary) view returns(uint256, uint256)
func (_Stake *StakeCallerSession) GetReleasableAmount(_beneficiary common.Address) (*big.Int, *big.Int, error) {
	return _Stake.Contract.GetReleasableAmount(&_Stake.CallOpts, _beneficiary)
}

// GetVestingSchedule is a free data retrieval call binding the contract method 0x9f829063.
//
// Solidity: function getVestingSchedule(address _beneficiary) view returns((address,uint256,uint256,uint8,uint256,uint256,uint256,uint256)[])
func (_Stake *StakeCaller) GetVestingSchedule(opts *bind.CallOpts, _beneficiary common.Address) ([]RewardVestingV1VestingSchedule, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "getVestingSchedule", _beneficiary)

	if err != nil {
		return *new([]RewardVestingV1VestingSchedule), err
	}

	out0 := *abi.ConvertType(out[0], new([]RewardVestingV1VestingSchedule)).(*[]RewardVestingV1VestingSchedule)

	return out0, err

}

// GetVestingSchedule is a free data retrieval call binding the contract method 0x9f829063.
//
// Solidity: function getVestingSchedule(address _beneficiary) view returns((address,uint256,uint256,uint8,uint256,uint256,uint256,uint256)[])
func (_Stake *StakeSession) GetVestingSchedule(_beneficiary common.Address) ([]RewardVestingV1VestingSchedule, error) {
	return _Stake.Contract.GetVestingSchedule(&_Stake.CallOpts, _beneficiary)
}

// GetVestingSchedule is a free data retrieval call binding the contract method 0x9f829063.
//
// Solidity: function getVestingSchedule(address _beneficiary) view returns((address,uint256,uint256,uint8,uint256,uint256,uint256,uint256)[])
func (_Stake *StakeCallerSession) GetVestingSchedule(_beneficiary common.Address) ([]RewardVestingV1VestingSchedule, error) {
	return _Stake.Contract.GetVestingSchedule(&_Stake.CallOpts, _beneficiary)
}

// PermanentTotal is a free data retrieval call binding the contract method 0x316d4a17.
//
// Solidity: function permanentTotal() view returns(uint256)
func (_Stake *StakeCaller) PermanentTotal(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "permanentTotal")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PermanentTotal is a free data retrieval call binding the contract method 0x316d4a17.
//
// Solidity: function permanentTotal() view returns(uint256)
func (_Stake *StakeSession) PermanentTotal() (*big.Int, error) {
	return _Stake.Contract.PermanentTotal(&_Stake.CallOpts)
}

// PermanentTotal is a free data retrieval call binding the contract method 0x316d4a17.
//
// Solidity: function permanentTotal() view returns(uint256)
func (_Stake *StakeCallerSession) PermanentTotal() (*big.Int, error) {
	return _Stake.Contract.PermanentTotal(&_Stake.CallOpts)
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

// VestingSchedules is a free data retrieval call binding the contract method 0x45626bd6.
//
// Solidity: function vestingSchedules(address , uint256 ) view returns(address beneficiary, uint256 start, uint256 duration, uint8 durationUnits, uint256 amountTotal, uint256 released, uint256 yieldRate, uint256 rewarded)
func (_Stake *StakeCaller) VestingSchedules(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	Beneficiary   common.Address
	Start         *big.Int
	Duration      *big.Int
	DurationUnits uint8
	AmountTotal   *big.Int
	Released      *big.Int
	YieldRate     *big.Int
	Rewarded      *big.Int
}, error) {
	var out []interface{}
	err := _Stake.contract.Call(opts, &out, "vestingSchedules", arg0, arg1)

	outstruct := new(struct {
		Beneficiary   common.Address
		Start         *big.Int
		Duration      *big.Int
		DurationUnits uint8
		AmountTotal   *big.Int
		Released      *big.Int
		YieldRate     *big.Int
		Rewarded      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Beneficiary = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Start = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Duration = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.DurationUnits = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.AmountTotal = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Released = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.YieldRate = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.Rewarded = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// VestingSchedules is a free data retrieval call binding the contract method 0x45626bd6.
//
// Solidity: function vestingSchedules(address , uint256 ) view returns(address beneficiary, uint256 start, uint256 duration, uint8 durationUnits, uint256 amountTotal, uint256 released, uint256 yieldRate, uint256 rewarded)
func (_Stake *StakeSession) VestingSchedules(arg0 common.Address, arg1 *big.Int) (struct {
	Beneficiary   common.Address
	Start         *big.Int
	Duration      *big.Int
	DurationUnits uint8
	AmountTotal   *big.Int
	Released      *big.Int
	YieldRate     *big.Int
	Rewarded      *big.Int
}, error) {
	return _Stake.Contract.VestingSchedules(&_Stake.CallOpts, arg0, arg1)
}

// VestingSchedules is a free data retrieval call binding the contract method 0x45626bd6.
//
// Solidity: function vestingSchedules(address , uint256 ) view returns(address beneficiary, uint256 start, uint256 duration, uint8 durationUnits, uint256 amountTotal, uint256 released, uint256 yieldRate, uint256 rewarded)
func (_Stake *StakeCallerSession) VestingSchedules(arg0 common.Address, arg1 *big.Int) (struct {
	Beneficiary   common.Address
	Start         *big.Int
	Duration      *big.Int
	DurationUnits uint8
	AmountTotal   *big.Int
	Released      *big.Int
	YieldRate     *big.Int
	Rewarded      *big.Int
}, error) {
	return _Stake.Contract.VestingSchedules(&_Stake.CallOpts, arg0, arg1)
}

// CreateVestingSchedule is a paid mutator transaction binding the contract method 0xb9f8dd0a.
//
// Solidity: function createVestingSchedule(address _beneficiary, uint256 _start, uint8 _durationUnits, uint256 _amountTotal) returns()
func (_Stake *StakeTransactor) CreateVestingSchedule(opts *bind.TransactOpts, _beneficiary common.Address, _start *big.Int, _durationUnits uint8, _amountTotal *big.Int) (*types.Transaction, error) {
	return _Stake.contract.Transact(opts, "createVestingSchedule", _beneficiary, _start, _durationUnits, _amountTotal)
}

// CreateVestingSchedule is a paid mutator transaction binding the contract method 0xb9f8dd0a.
//
// Solidity: function createVestingSchedule(address _beneficiary, uint256 _start, uint8 _durationUnits, uint256 _amountTotal) returns()
func (_Stake *StakeSession) CreateVestingSchedule(_beneficiary common.Address, _start *big.Int, _durationUnits uint8, _amountTotal *big.Int) (*types.Transaction, error) {
	return _Stake.Contract.CreateVestingSchedule(&_Stake.TransactOpts, _beneficiary, _start, _durationUnits, _amountTotal)
}

// CreateVestingSchedule is a paid mutator transaction binding the contract method 0xb9f8dd0a.
//
// Solidity: function createVestingSchedule(address _beneficiary, uint256 _start, uint8 _durationUnits, uint256 _amountTotal) returns()
func (_Stake *StakeTransactorSession) CreateVestingSchedule(_beneficiary common.Address, _start *big.Int, _durationUnits uint8, _amountTotal *big.Int) (*types.Transaction, error) {
	return _Stake.Contract.CreateVestingSchedule(&_Stake.TransactOpts, _beneficiary, _start, _durationUnits, _amountTotal)
}

// DepositPermanently is a paid mutator transaction binding the contract method 0x220ec6be.
//
// Solidity: function depositPermanently(uint256 _amount) returns()
func (_Stake *StakeTransactor) DepositPermanently(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _Stake.contract.Transact(opts, "depositPermanently", _amount)
}

// DepositPermanently is a paid mutator transaction binding the contract method 0x220ec6be.
//
// Solidity: function depositPermanently(uint256 _amount) returns()
func (_Stake *StakeSession) DepositPermanently(_amount *big.Int) (*types.Transaction, error) {
	return _Stake.Contract.DepositPermanently(&_Stake.TransactOpts, _amount)
}

// DepositPermanently is a paid mutator transaction binding the contract method 0x220ec6be.
//
// Solidity: function depositPermanently(uint256 _amount) returns()
func (_Stake *StakeTransactorSession) DepositPermanently(_amount *big.Int) (*types.Transaction, error) {
	return _Stake.Contract.DepositPermanently(&_Stake.TransactOpts, _amount)
}

// Release is a paid mutator transaction binding the contract method 0x19165587.
//
// Solidity: function release(address _beneficiary) returns()
func (_Stake *StakeTransactor) Release(opts *bind.TransactOpts, _beneficiary common.Address) (*types.Transaction, error) {
	return _Stake.contract.Transact(opts, "release", _beneficiary)
}

// Release is a paid mutator transaction binding the contract method 0x19165587.
//
// Solidity: function release(address _beneficiary) returns()
func (_Stake *StakeSession) Release(_beneficiary common.Address) (*types.Transaction, error) {
	return _Stake.Contract.Release(&_Stake.TransactOpts, _beneficiary)
}

// Release is a paid mutator transaction binding the contract method 0x19165587.
//
// Solidity: function release(address _beneficiary) returns()
func (_Stake *StakeTransactorSession) Release(_beneficiary common.Address) (*types.Transaction, error) {
	return _Stake.Contract.Release(&_Stake.TransactOpts, _beneficiary)
}

// StakeTokensReleasedIterator is returned from FilterTokensReleased and is used to iterate over the raw logs and unpacked data for TokensReleased events raised by the Stake contract.
type StakeTokensReleasedIterator struct {
	Event *StakeTokensReleased // Event containing the contract specifics and raw log

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
func (it *StakeTokensReleasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeTokensReleased)
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
		it.Event = new(StakeTokensReleased)
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
func (it *StakeTokensReleasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeTokensReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeTokensReleased represents a TokensReleased event raised by the Stake contract.
type StakeTokensReleased struct {
	Beneficiary common.Address
	Amount      *big.Int
	Reward      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTokensReleased is a free log retrieval operation binding the contract event 0xc5c52c2a9175470464d5ea4429889e7df2ea88630a3d32f4d0d3d2d448656210.
//
// Solidity: event TokensReleased(address indexed beneficiary, uint256 amount, uint256 reward)
func (_Stake *StakeFilterer) FilterTokensReleased(opts *bind.FilterOpts, beneficiary []common.Address) (*StakeTokensReleasedIterator, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}

	logs, sub, err := _Stake.contract.FilterLogs(opts, "TokensReleased", beneficiaryRule)
	if err != nil {
		return nil, err
	}
	return &StakeTokensReleasedIterator{contract: _Stake.contract, event: "TokensReleased", logs: logs, sub: sub}, nil
}

// WatchTokensReleased is a free log subscription operation binding the contract event 0xc5c52c2a9175470464d5ea4429889e7df2ea88630a3d32f4d0d3d2d448656210.
//
// Solidity: event TokensReleased(address indexed beneficiary, uint256 amount, uint256 reward)
func (_Stake *StakeFilterer) WatchTokensReleased(opts *bind.WatchOpts, sink chan<- *StakeTokensReleased, beneficiary []common.Address) (event.Subscription, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}

	logs, sub, err := _Stake.contract.WatchLogs(opts, "TokensReleased", beneficiaryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeTokensReleased)
				if err := _Stake.contract.UnpackLog(event, "TokensReleased", log); err != nil {
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

// ParseTokensReleased is a log parse operation binding the contract event 0xc5c52c2a9175470464d5ea4429889e7df2ea88630a3d32f4d0d3d2d448656210.
//
// Solidity: event TokensReleased(address indexed beneficiary, uint256 amount, uint256 reward)
func (_Stake *StakeFilterer) ParseTokensReleased(log types.Log) (*StakeTokensReleased, error) {
	event := new(StakeTokensReleased)
	if err := _Stake.contract.UnpackLog(event, "TokensReleased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeVestingScheduleCreatedIterator is returned from FilterVestingScheduleCreated and is used to iterate over the raw logs and unpacked data for VestingScheduleCreated events raised by the Stake contract.
type StakeVestingScheduleCreatedIterator struct {
	Event *StakeVestingScheduleCreated // Event containing the contract specifics and raw log

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
func (it *StakeVestingScheduleCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeVestingScheduleCreated)
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
		it.Event = new(StakeVestingScheduleCreated)
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
func (it *StakeVestingScheduleCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeVestingScheduleCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeVestingScheduleCreated represents a VestingScheduleCreated event raised by the Stake contract.
type StakeVestingScheduleCreated struct {
	Beneficiary   common.Address
	Start         *big.Int
	Duration      *big.Int
	DurationUnits uint8
	AmountTotal   *big.Int
	YieldRate     *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterVestingScheduleCreated is a free log retrieval operation binding the contract event 0xce0cdae96f18737abe4b02436163a3d9d15b0cf326715c28b72cd74502c9c424.
//
// Solidity: event VestingScheduleCreated(address indexed beneficiary, uint256 start, uint256 duration, uint8 durationUnits, uint256 amountTotal, uint256 yieldRate)
func (_Stake *StakeFilterer) FilterVestingScheduleCreated(opts *bind.FilterOpts, beneficiary []common.Address) (*StakeVestingScheduleCreatedIterator, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}

	logs, sub, err := _Stake.contract.FilterLogs(opts, "VestingScheduleCreated", beneficiaryRule)
	if err != nil {
		return nil, err
	}
	return &StakeVestingScheduleCreatedIterator{contract: _Stake.contract, event: "VestingScheduleCreated", logs: logs, sub: sub}, nil
}

// WatchVestingScheduleCreated is a free log subscription operation binding the contract event 0xce0cdae96f18737abe4b02436163a3d9d15b0cf326715c28b72cd74502c9c424.
//
// Solidity: event VestingScheduleCreated(address indexed beneficiary, uint256 start, uint256 duration, uint8 durationUnits, uint256 amountTotal, uint256 yieldRate)
func (_Stake *StakeFilterer) WatchVestingScheduleCreated(opts *bind.WatchOpts, sink chan<- *StakeVestingScheduleCreated, beneficiary []common.Address) (event.Subscription, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}

	logs, sub, err := _Stake.contract.WatchLogs(opts, "VestingScheduleCreated", beneficiaryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeVestingScheduleCreated)
				if err := _Stake.contract.UnpackLog(event, "VestingScheduleCreated", log); err != nil {
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

// ParseVestingScheduleCreated is a log parse operation binding the contract event 0xce0cdae96f18737abe4b02436163a3d9d15b0cf326715c28b72cd74502c9c424.
//
// Solidity: event VestingScheduleCreated(address indexed beneficiary, uint256 start, uint256 duration, uint8 durationUnits, uint256 amountTotal, uint256 yieldRate)
func (_Stake *StakeFilterer) ParseVestingScheduleCreated(log types.Log) (*StakeVestingScheduleCreated, error) {
	event := new(StakeVestingScheduleCreated)
	if err := _Stake.contract.UnpackLog(event, "VestingScheduleCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
