// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

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
	ABI: "[{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"withdrawalCredentials\",\"type\":\"bytes20\",\"indexed\":false,\"internalType\":\"bytes20\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawal\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"withdrawalCredentials\",\"type\":\"bytes20\",\"indexed\":false,\"internalType\":\"bytes20\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InsufficientDeposit\",\"inputs\":[]}]",
	Bin: "0x608060405234801561000f575f80fd5b506102598061001d5f395ff3fe608060405260043610610028575f3560e01c80632414da0c1461002c5780639813954e14610041575b5f80fd5b61003f61003a366004610119565b610054565b005b61003f61004f366004610119565b6100e1565b633b9aca008167ffffffffffffffff16101561009c576040517f0e1eddda00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7fb0a17f48b507cc393651a59f96277d9c6446062fe85d2754d762b6a2c32d93e583833360601b846040516100d4949392919061019b565b60405180910390a1505050565b7f2f41752ce57ce57aa4e6e12c3f5938cf141797e0d741f4a0abc3d328156652d783833360601b846040516100d4949392919061019b565b5f805f6040848603121561012b575f80fd5b833567ffffffffffffffff80821115610142575f80fd5b818601915086601f830112610155575f80fd5b813581811115610163575f80fd5b876020828501011115610174575f80fd5b6020928301955093509085013590808216821461018f575f80fd5b50809150509250925092565b60608152836060820152838560808301375f608085830101525f60807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f87011683010190507fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008416602083015267ffffffffffffffff831660408301529594505050505056fea26469706673582212207281222e182e7f2e3fc067c4553fd196bb600304c79e7f1be9328a1ecfe7965064736f6c63430008180033",
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

// Deposit is a paid mutator transaction binding the contract method 0x2414da0c.
//
// Solidity: function deposit(bytes validatorPubkey, uint64 amount) payable returns()
func (_Staking *StakingTransactor) Deposit(opts *bind.TransactOpts, validatorPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "deposit", validatorPubkey, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x2414da0c.
//
// Solidity: function deposit(bytes validatorPubkey, uint64 amount) payable returns()
func (_Staking *StakingSession) Deposit(validatorPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.Deposit(&_Staking.TransactOpts, validatorPubkey, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x2414da0c.
//
// Solidity: function deposit(bytes validatorPubkey, uint64 amount) payable returns()
func (_Staking *StakingTransactorSession) Deposit(validatorPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.Deposit(&_Staking.TransactOpts, validatorPubkey, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x9813954e.
//
// Solidity: function withdraw(bytes validatorPubkey, uint64 amount) payable returns()
func (_Staking *StakingTransactor) Withdraw(opts *bind.TransactOpts, validatorPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdraw", validatorPubkey, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x9813954e.
//
// Solidity: function withdraw(bytes validatorPubkey, uint64 amount) payable returns()
func (_Staking *StakingSession) Withdraw(validatorPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.Withdraw(&_Staking.TransactOpts, validatorPubkey, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x9813954e.
//
// Solidity: function withdraw(bytes validatorPubkey, uint64 amount) payable returns()
func (_Staking *StakingTransactorSession) Withdraw(validatorPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.Withdraw(&_Staking.TransactOpts, validatorPubkey, amount)
}

// StakingDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Staking contract.
type StakingDepositIterator struct {
	Event *StakingDeposit // Event containing the contract specifics and raw log

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
func (it *StakingDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingDeposit)
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
		it.Event = new(StakingDeposit)
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
func (it *StakingDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingDeposit represents a Deposit event raised by the Staking contract.
type StakingDeposit struct {
	ValidatorPubkey       []byte
	WithdrawalCredentials [20]byte
	Amount                uint64
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xb0a17f48b507cc393651a59f96277d9c6446062fe85d2754d762b6a2c32d93e5.
//
// Solidity: event Deposit(bytes validatorPubkey, bytes20 withdrawalCredentials, uint64 amount)
func (_Staking *StakingFilterer) FilterDeposit(opts *bind.FilterOpts) (*StakingDepositIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &StakingDepositIterator{contract: _Staking.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xb0a17f48b507cc393651a59f96277d9c6446062fe85d2754d762b6a2c32d93e5.
//
// Solidity: event Deposit(bytes validatorPubkey, bytes20 withdrawalCredentials, uint64 amount)
func (_Staking *StakingFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *StakingDeposit) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingDeposit)
				if err := _Staking.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xb0a17f48b507cc393651a59f96277d9c6446062fe85d2754d762b6a2c32d93e5.
//
// Solidity: event Deposit(bytes validatorPubkey, bytes20 withdrawalCredentials, uint64 amount)
func (_Staking *StakingFilterer) ParseDeposit(log types.Log) (*StakingDeposit, error) {
	event := new(StakingDeposit)
	if err := _Staking.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the Staking contract.
type StakingWithdrawalIterator struct {
	Event *StakingWithdrawal // Event containing the contract specifics and raw log

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
func (it *StakingWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingWithdrawal)
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
		it.Event = new(StakingWithdrawal)
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
func (it *StakingWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingWithdrawal represents a Withdrawal event raised by the Staking contract.
type StakingWithdrawal struct {
	ValidatorPubkey       []byte
	WithdrawalCredentials [20]byte
	Amount                uint64
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x2f41752ce57ce57aa4e6e12c3f5938cf141797e0d741f4a0abc3d328156652d7.
//
// Solidity: event Withdrawal(bytes validatorPubkey, bytes20 withdrawalCredentials, uint64 amount)
func (_Staking *StakingFilterer) FilterWithdrawal(opts *bind.FilterOpts) (*StakingWithdrawalIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return &StakingWithdrawalIterator{contract: _Staking.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x2f41752ce57ce57aa4e6e12c3f5938cf141797e0d741f4a0abc3d328156652d7.
//
// Solidity: event Withdrawal(bytes validatorPubkey, bytes20 withdrawalCredentials, uint64 amount)
func (_Staking *StakingFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *StakingWithdrawal) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingWithdrawal)
				if err := _Staking.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0x2f41752ce57ce57aa4e6e12c3f5938cf141797e0d741f4a0abc3d328156652d7.
//
// Solidity: event Withdrawal(bytes validatorPubkey, bytes20 withdrawalCredentials, uint64 amount)
func (_Staking *StakingFilterer) ParseWithdrawal(log types.Log) (*StakingWithdrawal, error) {
	event := new(StakingWithdrawal)
	if err := _Staking.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
