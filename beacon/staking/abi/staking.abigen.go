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
	ABI: "[{\"type\":\"function\",\"name\":\"delegateFn\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"withdrawalCredentials\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"undelegateFn\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Delegate\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"withdrawalCredentials\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"nonce\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Undelegate\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"nonce\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561000f575f80fd5b506106af8061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610034575f3560e01c80631808312f14610038578063b478fa9c1461004d575b5f80fd5b61004b6100463660046103e1565b610060565b005b61004b61005b36600461044f565b6100cc565b7f300164b233189ce316983af16e299425432bf341ca1c41a2f0439e326f8b6b038585858561008e86610132565b6100985f54610132565b6040516100aa96959493929190610551565b60405180910390a15f805490806100c0836105ac565b91905055505050505050565b7f41384f3d85152f1d3e41c1e8c41a78dc5e52dcd8a340bf5af92212cb9598dc7483836100f884610132565b6101025f54610132565b6040516101129493929190610608565b60405180910390a15f80549080610128836105ac565b9190505550505050565b60408051600880825281830190925260609160208201818036833701905050905060c082901b8060071a60f81b825f815181106101715761017161064c565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a9053508060061a60f81b826001815181106101b9576101b961064c565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a9053508060051a60f81b826002815181106102015761020161064c565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a9053508060041a60f81b826003815181106102495761024961064c565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a9053508060031a60f81b826004815181106102915761029161064c565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a9053508060021a60f81b826005815181106102d9576102d961064c565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a9053508060011a60f81b826006815181106103215761032161064c565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a905350805f1a60f81b826007815181106103685761036861064c565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191690815f1a90535050919050565b5f8083601f8401126103ac575f80fd5b50813567ffffffffffffffff8111156103c3575f80fd5b6020830191508360208285010111156103da575f80fd5b9250929050565b5f805f805f606086880312156103f5575f80fd5b853567ffffffffffffffff8082111561040c575f80fd5b61041889838a0161039c565b90975095506020880135915080821115610430575f80fd5b5061043d8882890161039c565b96999598509660400135949350505050565b5f805f60408486031215610461575f80fd5b833567ffffffffffffffff80821115610478575f80fd5b6104848783880161039c565b909550935060208601359150808216821461049d575f80fd5b50809150509250925092565b81835281816020850137505f602082840101525f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b5f81518084525f5b81811015610514576020818501810151868301820152016104f8565b505f6020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b608081525f61056460808301888a6104a9565b82810360208401526105778187896104a9565b9050828103604084015261058b81866104f0565b9050828103606084015261059f81856104f0565b9998505050505050505050565b5f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610601577f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5060010190565b606081525f61061b6060830186886104a9565b828103602084015261062d81866104f0565b9050828103604084015261064181856104f0565b979650505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffdfea2646970667358221220b71e6951ac337f9da037d639676c80c053d6a2c35ffb7afef9d368f1a5b2ea9d64736f6c63430008180033",
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

// DelegateFn is a paid mutator transaction binding the contract method 0x1808312f.
//
// Solidity: function delegateFn(bytes validatorPubkey, bytes withdrawalCredentials, uint256 amount) returns()
func (_Staking *StakingTransactor) DelegateFn(opts *bind.TransactOpts, validatorPubkey []byte, withdrawalCredentials []byte, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "delegateFn", validatorPubkey, withdrawalCredentials, amount)
}

// DelegateFn is a paid mutator transaction binding the contract method 0x1808312f.
//
// Solidity: function delegateFn(bytes validatorPubkey, bytes withdrawalCredentials, uint256 amount) returns()
func (_Staking *StakingSession) DelegateFn(validatorPubkey []byte, withdrawalCredentials []byte, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.DelegateFn(&_Staking.TransactOpts, validatorPubkey, withdrawalCredentials, amount)
}

// DelegateFn is a paid mutator transaction binding the contract method 0x1808312f.
//
// Solidity: function delegateFn(bytes validatorPubkey, bytes withdrawalCredentials, uint256 amount) returns()
func (_Staking *StakingTransactorSession) DelegateFn(validatorPubkey []byte, withdrawalCredentials []byte, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.DelegateFn(&_Staking.TransactOpts, validatorPubkey, withdrawalCredentials, amount)
}

// UndelegateFn is a paid mutator transaction binding the contract method 0xb478fa9c.
//
// Solidity: function undelegateFn(bytes validatorPubkey, uint64 amount) returns()
func (_Staking *StakingTransactor) UndelegateFn(opts *bind.TransactOpts, validatorPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "undelegateFn", validatorPubkey, amount)
}

// UndelegateFn is a paid mutator transaction binding the contract method 0xb478fa9c.
//
// Solidity: function undelegateFn(bytes validatorPubkey, uint64 amount) returns()
func (_Staking *StakingSession) UndelegateFn(validatorPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.UndelegateFn(&_Staking.TransactOpts, validatorPubkey, amount)
}

// UndelegateFn is a paid mutator transaction binding the contract method 0xb478fa9c.
//
// Solidity: function undelegateFn(bytes validatorPubkey, uint64 amount) returns()
func (_Staking *StakingTransactorSession) UndelegateFn(validatorPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.UndelegateFn(&_Staking.TransactOpts, validatorPubkey, amount)
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
	ValidatorPubkey       []byte
	WithdrawalCredentials []byte
	Amount                []byte
	Nonce                 []byte
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterDelegate is a free log retrieval operation binding the contract event 0x300164b233189ce316983af16e299425432bf341ca1c41a2f0439e326f8b6b03.
//
// Solidity: event Delegate(bytes validatorPubkey, bytes withdrawalCredentials, bytes amount, bytes nonce)
func (_Staking *StakingFilterer) FilterDelegate(opts *bind.FilterOpts) (*StakingDelegateIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Delegate")
	if err != nil {
		return nil, err
	}
	return &StakingDelegateIterator{contract: _Staking.contract, event: "Delegate", logs: logs, sub: sub}, nil
}

// WatchDelegate is a free log subscription operation binding the contract event 0x300164b233189ce316983af16e299425432bf341ca1c41a2f0439e326f8b6b03.
//
// Solidity: event Delegate(bytes validatorPubkey, bytes withdrawalCredentials, bytes amount, bytes nonce)
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

// ParseDelegate is a log parse operation binding the contract event 0x300164b233189ce316983af16e299425432bf341ca1c41a2f0439e326f8b6b03.
//
// Solidity: event Delegate(bytes validatorPubkey, bytes withdrawalCredentials, bytes amount, bytes nonce)
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
	ValidatorPubkey []byte
	Amount          []byte
	Nonce           []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUndelegate is a free log retrieval operation binding the contract event 0x41384f3d85152f1d3e41c1e8c41a78dc5e52dcd8a340bf5af92212cb9598dc74.
//
// Solidity: event Undelegate(bytes validatorPubkey, bytes amount, bytes nonce)
func (_Staking *StakingFilterer) FilterUndelegate(opts *bind.FilterOpts) (*StakingUndelegateIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Undelegate")
	if err != nil {
		return nil, err
	}
	return &StakingUndelegateIterator{contract: _Staking.contract, event: "Undelegate", logs: logs, sub: sub}, nil
}

// WatchUndelegate is a free log subscription operation binding the contract event 0x41384f3d85152f1d3e41c1e8c41a78dc5e52dcd8a340bf5af92212cb9598dc74.
//
// Solidity: event Undelegate(bytes validatorPubkey, bytes amount, bytes nonce)
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

// ParseUndelegate is a log parse operation binding the contract event 0x41384f3d85152f1d3e41c1e8c41a78dc5e52dcd8a340bf5af92212cb9598dc74.
//
// Solidity: event Undelegate(bytes validatorPubkey, bytes amount, bytes nonce)
func (_Staking *StakingFilterer) ParseUndelegate(log types.Log) (*StakingUndelegate, error) {
	event := new(StakingUndelegate)
	if err := _Staking.contract.UnpackLog(event, "Undelegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
