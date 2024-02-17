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
	ABI: "[{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"withdrawalCredentials\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"withdrawalCredentials\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"nonce\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"nonce\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561000f575f80fd5b506105468061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610034575f3560e01c80636993e9a2146100385780639813954e1461004d575b5f80fd5b61004b6100463660046102a2565b610060565b005b61004b61005b36600461031d565b610157565b7f9d0a206c338cfcaf9198c04fe61b39a988e26b623ef97cb2f72bafcfbf8bb93e85858585856040516020016100c1919060c09190911b7fffffffffffffffff00000000000000000000000000000000000000000000000016815260080190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282525f5460208401529101604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052610135969594939291610415565b60405180910390a15f8054908061014b83610470565b91905055505050505050565b6040517fffffffffffffffff00000000000000000000000000000000000000000000000060c083901b1660208201527f6b62663c2f1027e5a27acccfc1125fd051e7cd14f2273bad45c79fc78535201a9084908490602801604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181528282525f5460208401529101604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052610221949392916104cc565b60405180910390a15f8054908061023783610470565b9190505550505050565b5f8083601f840112610251575f80fd5b50813567ffffffffffffffff811115610268575f80fd5b60208301915083602082850101111561027f575f80fd5b9250929050565b803567ffffffffffffffff8116811461029d575f80fd5b919050565b5f805f805f606086880312156102b6575f80fd5b853567ffffffffffffffff808211156102cd575f80fd5b6102d989838a01610241565b909750955060208801359150808211156102f1575f80fd5b506102fe88828901610241565b9094509250610311905060408701610286565b90509295509295909350565b5f805f6040848603121561032f575f80fd5b833567ffffffffffffffff811115610345575f80fd5b61035186828701610241565b9094509250610364905060208501610286565b90509250925092565b81835281816020850137505f602082840101525f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b5f81518084525f5b818110156103d8576020818501810151868301820152016103bc565b505f6020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b608081525f61042860808301888a61036d565b828103602084015261043b81878961036d565b9050828103604084015261044f81866103b4565b9050828103606084015261046381856103b4565b9998505050505050505050565b5f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036104c5577f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5060010190565b606081525f6104df60608301868861036d565b82810360208401526104f181866103b4565b9050828103604084015261050581856103b4565b97965050505050505056fea26469706673582212206d12c9477aa932b19819d13857428e51179bb382b62a20ac5f07143e1c31ecbb64736f6c63430008180033",
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

// Deposit is a paid mutator transaction binding the contract method 0x6993e9a2.
//
// Solidity: function deposit(bytes validatorPubkey, bytes withdrawalCredentials, uint64 amount) returns()
func (_Staking *StakingTransactor) Deposit(opts *bind.TransactOpts, validatorPubkey []byte, withdrawalCredentials []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "deposit", validatorPubkey, withdrawalCredentials, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x6993e9a2.
//
// Solidity: function deposit(bytes validatorPubkey, bytes withdrawalCredentials, uint64 amount) returns()
func (_Staking *StakingSession) Deposit(validatorPubkey []byte, withdrawalCredentials []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.Deposit(&_Staking.TransactOpts, validatorPubkey, withdrawalCredentials, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x6993e9a2.
//
// Solidity: function deposit(bytes validatorPubkey, bytes withdrawalCredentials, uint64 amount) returns()
func (_Staking *StakingTransactorSession) Deposit(validatorPubkey []byte, withdrawalCredentials []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.Deposit(&_Staking.TransactOpts, validatorPubkey, withdrawalCredentials, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x9813954e.
//
// Solidity: function withdraw(bytes validatorPubkey, uint64 amount) returns()
func (_Staking *StakingTransactor) Withdraw(opts *bind.TransactOpts, validatorPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdraw", validatorPubkey, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x9813954e.
//
// Solidity: function withdraw(bytes validatorPubkey, uint64 amount) returns()
func (_Staking *StakingSession) Withdraw(validatorPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.Withdraw(&_Staking.TransactOpts, validatorPubkey, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x9813954e.
//
// Solidity: function withdraw(bytes validatorPubkey, uint64 amount) returns()
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
	WithdrawalCredentials []byte
	Amount                []byte
	Nonce                 []byte
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x9d0a206c338cfcaf9198c04fe61b39a988e26b623ef97cb2f72bafcfbf8bb93e.
//
// Solidity: event Deposit(bytes validatorPubkey, bytes withdrawalCredentials, bytes amount, bytes nonce)
func (_Staking *StakingFilterer) FilterDeposit(opts *bind.FilterOpts) (*StakingDepositIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &StakingDepositIterator{contract: _Staking.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x9d0a206c338cfcaf9198c04fe61b39a988e26b623ef97cb2f72bafcfbf8bb93e.
//
// Solidity: event Deposit(bytes validatorPubkey, bytes withdrawalCredentials, bytes amount, bytes nonce)
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

// ParseDeposit is a log parse operation binding the contract event 0x9d0a206c338cfcaf9198c04fe61b39a988e26b623ef97cb2f72bafcfbf8bb93e.
//
// Solidity: event Deposit(bytes validatorPubkey, bytes withdrawalCredentials, bytes amount, bytes nonce)
func (_Staking *StakingFilterer) ParseDeposit(log types.Log) (*StakingDeposit, error) {
	event := new(StakingDeposit)
	if err := _Staking.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Staking contract.
type StakingWithdrawIterator struct {
	Event *StakingWithdraw // Event containing the contract specifics and raw log

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
func (it *StakingWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingWithdraw)
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
		it.Event = new(StakingWithdraw)
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
func (it *StakingWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingWithdraw represents a Withdraw event raised by the Staking contract.
type StakingWithdraw struct {
	ValidatorPubkey []byte
	Amount          []byte
	Nonce           []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x6b62663c2f1027e5a27acccfc1125fd051e7cd14f2273bad45c79fc78535201a.
//
// Solidity: event Withdraw(bytes validatorPubkey, bytes amount, bytes nonce)
func (_Staking *StakingFilterer) FilterWithdraw(opts *bind.FilterOpts) (*StakingWithdrawIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &StakingWithdrawIterator{contract: _Staking.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x6b62663c2f1027e5a27acccfc1125fd051e7cd14f2273bad45c79fc78535201a.
//
// Solidity: event Withdraw(bytes validatorPubkey, bytes amount, bytes nonce)
func (_Staking *StakingFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *StakingWithdraw) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingWithdraw)
				if err := _Staking.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x6b62663c2f1027e5a27acccfc1125fd051e7cd14f2273bad45c79fc78535201a.
//
// Solidity: event Withdraw(bytes validatorPubkey, bytes amount, bytes nonce)
func (_Staking *StakingFilterer) ParseWithdraw(log types.Log) (*StakingWithdraw, error) {
	event := new(StakingWithdraw)
	if err := _Staking.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
