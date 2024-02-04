// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// StakingMetaData contains all meta data concerning the Staking contract.
var StakingMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"delegateFn\",\"inputs\":[{\"name\":\"operatorAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"undelegateFn\",\"inputs\":[{\"name\":\"operatorAddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Delegate\",\"inputs\":[{\"name\":\"operatorAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Undelegate\",\"inputs\":[{\"name\":\"operatorAddress\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561000f575f80fd5b506102b88061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610034575f3560e01c80630feae41b14610038578063b9da3dc214610054575b5f80fd5b610052600480360381019061004d919061018c565b610070565b005b61006e6004803603810190610069919061018c565b6100b0565b005b7f5e3b2166f1015d5880455f71c4745d2b76ed6139db23aef0d5eb6ae60592a3098383836040516100a393929190610252565b60405180910390a1505050565b7fe7e9c5a6880a7c4d7787acf85fb43d54e090f90ba1abec26d8d5d9cb9d1a09048383836040516100e393929190610252565b60405180910390a1505050565b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f8083601f840112610119576101186100f8565b5b8235905067ffffffffffffffff811115610136576101356100fc565b5b60208301915083600182028301111561015257610151610100565b5b9250929050565b5f819050919050565b61016b81610159565b8114610175575f80fd5b50565b5f8135905061018681610162565b92915050565b5f805f604084860312156101a3576101a26100f0565b5b5f84013567ffffffffffffffff8111156101c0576101bf6100f4565b5b6101cc86828701610104565b935093505060206101df86828701610178565b9150509250925092565b5f82825260208201905092915050565b828183375f83830152505050565b5f601f19601f8301169050919050565b5f61022283856101e9565b935061022f8385846101f9565b61023883610207565b840190509392505050565b61024c81610159565b82525050565b5f6040820190508181035f83015261026b818587610217565b905061027a6020830184610243565b94935050505056fea26469706673582212205792c69d80ff3a3588c287f008dbb2ad5d81b9006f9d916f3591fe798899f98864736f6c63430008180033",
}

// StakingABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingMetaData.ABI instead.
var StakingABI = StakingMetaData.ABI

// StakingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakingMetaData.Bin instead.
var StakingBin = StakingMetaData.Bin

// DeployStaking deploys a new Ethereum contract, binding an instance of Staking to it.
func DeployStaking(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Staking, error) {
	parsed, err := StakingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// Staking is an auto generated Go binding around an Ethereum contract.
type Staking struct {
	StakingCaller     // Read-only binding to the contract
	StakingTransactor // Write-only binding to the contract
	StakingFilterer   // Log filterer for contract events
}

// StakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingSession struct {
	Contract     *Staking          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingCallerSession struct {
	Contract *StakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingTransactorSession struct {
	Contract     *StakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingRaw struct {
	Contract *Staking // Generic contract binding to access the raw methods on
}

// StakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingCallerRaw struct {
	Contract *StakingCaller // Generic read-only contract binding to access the raw methods on
}

// StakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingTransactorRaw struct {
	Contract *StakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStaking creates a new instance of Staking, bound to a specific deployed contract.
func NewStaking(address common.Address, backend bind.ContractBackend) (*Staking, error) {
	contract, err := bindStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// NewStakingCaller creates a new read-only instance of Staking, bound to a specific deployed contract.
func NewStakingCaller(address common.Address, caller bind.ContractCaller) (*StakingCaller, error) {
	contract, err := bindStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingCaller{contract: contract}, nil
}

// NewStakingTransactor creates a new write-only instance of Staking, bound to a specific deployed contract.
func NewStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingTransactor, error) {
	contract, err := bindStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTransactor{contract: contract}, nil
}

// NewStakingFilterer creates a new log filterer instance of Staking, bound to a specific deployed contract.
func NewStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingFilterer, error) {
	contract, err := bindStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingFilterer{contract: contract}, nil
}

// bindStaking binds a generic wrapper to an already deployed contract.
func bindStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.StakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transact(opts, method, params...)
}

// DelegateFn is a paid mutator transaction binding the contract method 0xb9da3dc2.
//
// Solidity: function delegateFn(string operatorAddress, uint256 amount) returns()
func (_Staking *StakingTransactor) DelegateFn(opts *bind.TransactOpts, operatorAddress string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "delegateFn", operatorAddress, amount)
}

// DelegateFn is a paid mutator transaction binding the contract method 0xb9da3dc2.
//
// Solidity: function delegateFn(string operatorAddress, uint256 amount) returns()
func (_Staking *StakingSession) DelegateFn(operatorAddress string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.DelegateFn(&_Staking.TransactOpts, operatorAddress, amount)
}

// DelegateFn is a paid mutator transaction binding the contract method 0xb9da3dc2.
//
// Solidity: function delegateFn(string operatorAddress, uint256 amount) returns()
func (_Staking *StakingTransactorSession) DelegateFn(operatorAddress string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.DelegateFn(&_Staking.TransactOpts, operatorAddress, amount)
}

// UndelegateFn is a paid mutator transaction binding the contract method 0x0feae41b.
//
// Solidity: function undelegateFn(string operatorAddress, uint256 amount) returns()
func (_Staking *StakingTransactor) UndelegateFn(opts *bind.TransactOpts, operatorAddress string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "undelegateFn", operatorAddress, amount)
}

// UndelegateFn is a paid mutator transaction binding the contract method 0x0feae41b.
//
// Solidity: function undelegateFn(string operatorAddress, uint256 amount) returns()
func (_Staking *StakingSession) UndelegateFn(operatorAddress string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.UndelegateFn(&_Staking.TransactOpts, operatorAddress, amount)
}

// UndelegateFn is a paid mutator transaction binding the contract method 0x0feae41b.
//
// Solidity: function undelegateFn(string operatorAddress, uint256 amount) returns()
func (_Staking *StakingTransactorSession) UndelegateFn(operatorAddress string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.UndelegateFn(&_Staking.TransactOpts, operatorAddress, amount)
}

// StakingDelegateIterator is returned from FilterDelegate and is used to iterate over the raw logs and unpacked data for Delegate events raised by the Staking contract.
type StakingDelegateIterator struct {
	Event *StakingDelegate // Event containing the contract specifics and raw log

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
func (it *StakingDelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingDelegate)
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
		it.Event = new(StakingDelegate)
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
func (it *StakingDelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingDelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingDelegate represents a Delegate event raised by the Staking contract.
type StakingDelegate struct {
	OperatorAddress string
	Amount          *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterDelegate is a free log retrieval operation binding the contract event 0xe7e9c5a6880a7c4d7787acf85fb43d54e090f90ba1abec26d8d5d9cb9d1a0904.
//
// Solidity: event Delegate(string operatorAddress, uint256 amount)
func (_Staking *StakingFilterer) FilterDelegate(opts *bind.FilterOpts) (*StakingDelegateIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Delegate")
	if err != nil {
		return nil, err
	}
	return &StakingDelegateIterator{contract: _Staking.contract, event: "Delegate", logs: logs, sub: sub}, nil
}

// WatchDelegate is a free log subscription operation binding the contract event 0xe7e9c5a6880a7c4d7787acf85fb43d54e090f90ba1abec26d8d5d9cb9d1a0904.
//
// Solidity: event Delegate(string operatorAddress, uint256 amount)
func (_Staking *StakingFilterer) WatchDelegate(opts *bind.WatchOpts, sink chan<- *StakingDelegate) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Delegate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingDelegate)
				if err := _Staking.contract.UnpackLog(event, "Delegate", log); err != nil {
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

// ParseDelegate is a log parse operation binding the contract event 0xe7e9c5a6880a7c4d7787acf85fb43d54e090f90ba1abec26d8d5d9cb9d1a0904.
//
// Solidity: event Delegate(string operatorAddress, uint256 amount)
func (_Staking *StakingFilterer) ParseDelegate(log types.Log) (*StakingDelegate, error) {
	event := new(StakingDelegate)
	if err := _Staking.contract.UnpackLog(event, "Delegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingUndelegateIterator is returned from FilterUndelegate and is used to iterate over the raw logs and unpacked data for Undelegate events raised by the Staking contract.
type StakingUndelegateIterator struct {
	Event *StakingUndelegate // Event containing the contract specifics and raw log

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
func (it *StakingUndelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUndelegate)
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
		it.Event = new(StakingUndelegate)
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
func (it *StakingUndelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUndelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUndelegate represents a Undelegate event raised by the Staking contract.
type StakingUndelegate struct {
	OperatorAddress string
	Amount          *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUndelegate is a free log retrieval operation binding the contract event 0x5e3b2166f1015d5880455f71c4745d2b76ed6139db23aef0d5eb6ae60592a309.
//
// Solidity: event Undelegate(string operatorAddress, uint256 amount)
func (_Staking *StakingFilterer) FilterUndelegate(opts *bind.FilterOpts) (*StakingUndelegateIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Undelegate")
	if err != nil {
		return nil, err
	}
	return &StakingUndelegateIterator{contract: _Staking.contract, event: "Undelegate", logs: logs, sub: sub}, nil
}

// WatchUndelegate is a free log subscription operation binding the contract event 0x5e3b2166f1015d5880455f71c4745d2b76ed6139db23aef0d5eb6ae60592a309.
//
// Solidity: event Undelegate(string operatorAddress, uint256 amount)
func (_Staking *StakingFilterer) WatchUndelegate(opts *bind.WatchOpts, sink chan<- *StakingUndelegate) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Undelegate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUndelegate)
				if err := _Staking.contract.UnpackLog(event, "Undelegate", log); err != nil {
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

// ParseUndelegate is a log parse operation binding the contract event 0x5e3b2166f1015d5880455f71c4745d2b76ed6139db23aef0d5eb6ae60592a309.
//
// Solidity: event Undelegate(string operatorAddress, uint256 amount)
func (_Staking *StakingFilterer) ParseUndelegate(log types.Log) (*StakingUndelegate, error) {
	event := new(StakingUndelegate)
	if err := _Staking.contract.UnpackLog(event, "Undelegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
