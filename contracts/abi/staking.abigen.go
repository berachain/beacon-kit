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
	ABI: "[{\"type\":\"function\",\"name\":\"NATIVE_ASSET\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"stakingCredentials\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"redirect\",\"inputs\":[{\"name\":\"fromPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"toPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"withdrawalCredentials\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"stakingCredentials\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Redirect\",\"inputs\":[{\"name\":\"fromPubKey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"toPubKey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"stakingCredentials\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"fromPubKey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"withdrawalCredentials\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"DepositNotMultipleOfGwei\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DepositValueTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientDeposit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientRedirectAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientWithdrawAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCredentialsLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPubKeyLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignatureLength\",\"inputs\":[]}]",
	Bin: "0x60806040525f80546001600160a01b03191673eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee179055348015610034575f80fd5b506109f9806100425f395ff3fe60806040526004361061003e575f3560e01c806304a3267f146100425780635b70fa2914610063578063bf53253b14610076578063bf9b6a55146100c6575b5f80fd5b34801561004d575f80fd5b5061006161005c36600461077e565b6100e5565b005b6100616100713660046107f9565b610229565b348015610081575f80fd5b5061009d73eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee81565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b3480156100d1575f80fd5b506100616100e036600461077e565b6103a1565b6030841461011f576040517f9f10647200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208214610159576040517fb39bca1600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610169600a6407735940006108cc565b67ffffffffffffffff168167ffffffffffffffff1610156101b6576040517febec602100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82826040516101c69291906108f2565b604051809103902085856040516101de9291906108f2565b60405190819003812067ffffffffffffffff84168252907fd819a76a9128ab820538179b416ffb491e0fa0b23b2a08b605fba4c2649db9a69060200160405180910390a35050505050565b60308614610263576040517f9f10647200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6020841461029d576040517fb39bca1600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606081146102d7576040517f4be6321b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f5473ffffffffffffffffffffffffffffffffffffffff167fffffffffffffffffffffffff1111111111111111111111111111111111111112016103245761031d610501565b925061032d565b61032d836105d7565b848460405161033d9291906108f2565b604051809103902087876040516103559291906108f2565b60405180910390207f1f39b85dd1a529b31e0cd61e5609e1feca0e08e2103fe319fbd3dd5a0c7b68df85858560405161039093929190610901565b60405180910390a350505050505050565b6030841415806103b2575060308214155b156103e9576040517f9f10647200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103f9600a6407735940006108cc565b67ffffffffffffffff168167ffffffffffffffff161015610446576040517f0494a69c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516020808252337f01000000000000000000000000000000000000000000000000000000000000001790820152808201918290526104869161095e565b6040518091039020838360405161049e9291906108f2565b604051809103902086866040516104b69291906108f2565b60405190819003812067ffffffffffffffff85168252907fe161f5842757f257346b360594d094b7fa530f9404e93a80bf18bd8b14f9258f9060200160405180910390a45050505050565b5f67ffffffffffffffff341115610544576040517f2aa6673400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b640773594000341015610583576040517f0e1eddda00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610591633b9aca003461098a565b156105c8576040517f40567b3800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105d25f34610700565b503490565b5f546040517f9dc29fac00000000000000000000000000000000000000000000000000000000815233600482015267ffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff90911690639dc29fac906044015f604051808303815f87803b15801561064e575f80fd5b505af1158015610660573d5f803e3d5ffd5b50505064077359400067ffffffffffffffff8316101590506106ae576040517f0e1eddda00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106bc633b9aca008261099d565b67ffffffffffffffff16156106fd576040517f40567b3800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50565b5f385f3884865af16107195763b12d13eb5f526004601cfd5b5050565b5f8083601f84011261072d575f80fd5b50813567ffffffffffffffff811115610744575f80fd5b60208301915083602082850101111561075b575f80fd5b9250929050565b803567ffffffffffffffff81168114610779575f80fd5b919050565b5f805f805f60608688031215610792575f80fd5b853567ffffffffffffffff808211156107a9575f80fd5b6107b589838a0161071d565b909750955060208801359150808211156107cd575f80fd5b506107da8882890161071d565b90945092506107ed905060408701610762565b90509295509295909350565b5f805f805f805f6080888a03121561080f575f80fd5b873567ffffffffffffffff80821115610826575f80fd5b6108328b838c0161071d565b909950975060208a013591508082111561084a575f80fd5b6108568b838c0161071d565b909750955085915061086a60408b01610762565b945060608a013591508082111561087f575f80fd5b5061088c8a828b0161071d565b989b979a50959850939692959293505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f67ffffffffffffffff808416806108e6576108e661089f565b92169190910492915050565b818382375f9101908152919050565b67ffffffffffffffff8416815260406020820152816040820152818360608301375f818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b5f82515f5b8181101561097d5760208186018101518583015201610963565b505f920191825250919050565b5f826109985761099861089f565b500690565b5f67ffffffffffffffff808416806109b7576109b761089f565b9216919091069291505056fea26469706673582212204fdd8f7199f85c1743e6d116f8056d134bb5c31360147c960586016acb6abf9964736f6c63430008180033",
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

// NATIVEASSET is a free data retrieval call binding the contract method 0xbf53253b.
//
// Solidity: function NATIVE_ASSET() view returns(address)
func (_Staking *StakingCaller) NATIVEASSET(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "NATIVE_ASSET")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NATIVEASSET is a free data retrieval call binding the contract method 0xbf53253b.
//
// Solidity: function NATIVE_ASSET() view returns(address)
func (_Staking *StakingSession) NATIVEASSET() (common.Address, error) {
	return _Staking.Contract.NATIVEASSET(&_Staking.CallOpts)
}

// NATIVEASSET is a free data retrieval call binding the contract method 0xbf53253b.
//
// Solidity: function NATIVE_ASSET() view returns(address)
func (_Staking *StakingCallerSession) NATIVEASSET() (common.Address, error) {
	return _Staking.Contract.NATIVEASSET(&_Staking.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x5b70fa29.
//
// Solidity: function deposit(bytes validatorPubKey, bytes stakingCredentials, uint64 amount, bytes signature) payable returns()
func (_Staking *StakingTransactor) Deposit(opts *bind.TransactOpts, validatorPubKey []byte, stakingCredentials []byte, amount uint64, signature []byte) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "deposit", validatorPubKey, stakingCredentials, amount, signature)
}

// Deposit is a paid mutator transaction binding the contract method 0x5b70fa29.
//
// Solidity: function deposit(bytes validatorPubKey, bytes stakingCredentials, uint64 amount, bytes signature) payable returns()
func (_Staking *StakingSession) Deposit(validatorPubKey []byte, stakingCredentials []byte, amount uint64, signature []byte) (*types.Transaction, error) {
	return _Staking.Contract.Deposit(&_Staking.TransactOpts, validatorPubKey, stakingCredentials, amount, signature)
}

// Deposit is a paid mutator transaction binding the contract method 0x5b70fa29.
//
// Solidity: function deposit(bytes validatorPubKey, bytes stakingCredentials, uint64 amount, bytes signature) payable returns()
func (_Staking *StakingTransactorSession) Deposit(validatorPubKey []byte, stakingCredentials []byte, amount uint64, signature []byte) (*types.Transaction, error) {
	return _Staking.Contract.Deposit(&_Staking.TransactOpts, validatorPubKey, stakingCredentials, amount, signature)
}

// Redirect is a paid mutator transaction binding the contract method 0xbf9b6a55.
//
// Solidity: function redirect(bytes fromPubKey, bytes toPubKey, uint64 amount) returns()
func (_Staking *StakingTransactor) Redirect(opts *bind.TransactOpts, fromPubKey []byte, toPubKey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "redirect", fromPubKey, toPubKey, amount)
}

// Redirect is a paid mutator transaction binding the contract method 0xbf9b6a55.
//
// Solidity: function redirect(bytes fromPubKey, bytes toPubKey, uint64 amount) returns()
func (_Staking *StakingSession) Redirect(fromPubKey []byte, toPubKey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.Redirect(&_Staking.TransactOpts, fromPubKey, toPubKey, amount)
}

// Redirect is a paid mutator transaction binding the contract method 0xbf9b6a55.
//
// Solidity: function redirect(bytes fromPubKey, bytes toPubKey, uint64 amount) returns()
func (_Staking *StakingTransactorSession) Redirect(fromPubKey []byte, toPubKey []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.Redirect(&_Staking.TransactOpts, fromPubKey, toPubKey, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x04a3267f.
//
// Solidity: function withdraw(bytes validatorPubKey, bytes withdrawalCredentials, uint64 amount) returns()
func (_Staking *StakingTransactor) Withdraw(opts *bind.TransactOpts, validatorPubKey []byte, withdrawalCredentials []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdraw", validatorPubKey, withdrawalCredentials, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x04a3267f.
//
// Solidity: function withdraw(bytes validatorPubKey, bytes withdrawalCredentials, uint64 amount) returns()
func (_Staking *StakingSession) Withdraw(validatorPubKey []byte, withdrawalCredentials []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.Withdraw(&_Staking.TransactOpts, validatorPubKey, withdrawalCredentials, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x04a3267f.
//
// Solidity: function withdraw(bytes validatorPubKey, bytes withdrawalCredentials, uint64 amount) returns()
func (_Staking *StakingTransactorSession) Withdraw(validatorPubKey []byte, withdrawalCredentials []byte, amount uint64) (*types.Transaction, error) {
	return _Staking.Contract.Withdraw(&_Staking.TransactOpts, validatorPubKey, withdrawalCredentials, amount)
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
	ValidatorPubKey    common.Hash
	StakingCredentials common.Hash
	Amount             uint64
	Signature          []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x1f39b85dd1a529b31e0cd61e5609e1feca0e08e2103fe319fbd3dd5a0c7b68df.
//
// Solidity: event Deposit(bytes indexed validatorPubKey, bytes indexed stakingCredentials, uint64 amount, bytes signature)
func (_Staking *StakingFilterer) FilterDeposit(opts *bind.FilterOpts, validatorPubKey [][]byte, stakingCredentials [][]byte) (*StakingDepositIterator, error) {

	var validatorPubKeyRule []interface{}
	for _, validatorPubKeyItem := range validatorPubKey {
		validatorPubKeyRule = append(validatorPubKeyRule, validatorPubKeyItem)
	}
	var stakingCredentialsRule []interface{}
	for _, stakingCredentialsItem := range stakingCredentials {
		stakingCredentialsRule = append(stakingCredentialsRule, stakingCredentialsItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Deposit", validatorPubKeyRule, stakingCredentialsRule)
	if err != nil {
		return nil, err
	}
	return &StakingDepositIterator{contract: _Staking.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x1f39b85dd1a529b31e0cd61e5609e1feca0e08e2103fe319fbd3dd5a0c7b68df.
//
// Solidity: event Deposit(bytes indexed validatorPubKey, bytes indexed stakingCredentials, uint64 amount, bytes signature)
func (_Staking *StakingFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *StakingDeposit, validatorPubKey [][]byte, stakingCredentials [][]byte) (event.Subscription, error) {

	var validatorPubKeyRule []interface{}
	for _, validatorPubKeyItem := range validatorPubKey {
		validatorPubKeyRule = append(validatorPubKeyRule, validatorPubKeyItem)
	}
	var stakingCredentialsRule []interface{}
	for _, stakingCredentialsItem := range stakingCredentials {
		stakingCredentialsRule = append(stakingCredentialsRule, stakingCredentialsItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Deposit", validatorPubKeyRule, stakingCredentialsRule)
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

// ParseDeposit is a log parse operation binding the contract event 0x1f39b85dd1a529b31e0cd61e5609e1feca0e08e2103fe319fbd3dd5a0c7b68df.
//
// Solidity: event Deposit(bytes indexed validatorPubKey, bytes indexed stakingCredentials, uint64 amount, bytes signature)
func (_Staking *StakingFilterer) ParseDeposit(log types.Log) (*StakingDeposit, error) {
	event := new(StakingDeposit)
	if err := _Staking.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingRedirectIterator is returned from FilterRedirect and is used to iterate over the raw logs and unpacked data for Redirect events raised by the Staking contract.
type StakingRedirectIterator struct {
	Event *StakingRedirect // Event containing the contract specifics and raw log

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
func (it *StakingRedirectIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingRedirect)
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
		it.Event = new(StakingRedirect)
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
func (it *StakingRedirectIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingRedirectIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingRedirect represents a Redirect event raised by the Staking contract.
type StakingRedirect struct {
	FromPubKey         common.Hash
	ToPubKey           common.Hash
	StakingCredentials common.Hash
	Amount             uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterRedirect is a free log retrieval operation binding the contract event 0xe161f5842757f257346b360594d094b7fa530f9404e93a80bf18bd8b14f9258f.
//
// Solidity: event Redirect(bytes indexed fromPubKey, bytes indexed toPubKey, bytes indexed stakingCredentials, uint64 amount)
func (_Staking *StakingFilterer) FilterRedirect(opts *bind.FilterOpts, fromPubKey [][]byte, toPubKey [][]byte, stakingCredentials [][]byte) (*StakingRedirectIterator, error) {

	var fromPubKeyRule []interface{}
	for _, fromPubKeyItem := range fromPubKey {
		fromPubKeyRule = append(fromPubKeyRule, fromPubKeyItem)
	}
	var toPubKeyRule []interface{}
	for _, toPubKeyItem := range toPubKey {
		toPubKeyRule = append(toPubKeyRule, toPubKeyItem)
	}
	var stakingCredentialsRule []interface{}
	for _, stakingCredentialsItem := range stakingCredentials {
		stakingCredentialsRule = append(stakingCredentialsRule, stakingCredentialsItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Redirect", fromPubKeyRule, toPubKeyRule, stakingCredentialsRule)
	if err != nil {
		return nil, err
	}
	return &StakingRedirectIterator{contract: _Staking.contract, event: "Redirect", logs: logs, sub: sub}, nil
}

// WatchRedirect is a free log subscription operation binding the contract event 0xe161f5842757f257346b360594d094b7fa530f9404e93a80bf18bd8b14f9258f.
//
// Solidity: event Redirect(bytes indexed fromPubKey, bytes indexed toPubKey, bytes indexed stakingCredentials, uint64 amount)
func (_Staking *StakingFilterer) WatchRedirect(opts *bind.WatchOpts, sink chan<- *StakingRedirect, fromPubKey [][]byte, toPubKey [][]byte, stakingCredentials [][]byte) (event.Subscription, error) {

	var fromPubKeyRule []interface{}
	for _, fromPubKeyItem := range fromPubKey {
		fromPubKeyRule = append(fromPubKeyRule, fromPubKeyItem)
	}
	var toPubKeyRule []interface{}
	for _, toPubKeyItem := range toPubKey {
		toPubKeyRule = append(toPubKeyRule, toPubKeyItem)
	}
	var stakingCredentialsRule []interface{}
	for _, stakingCredentialsItem := range stakingCredentials {
		stakingCredentialsRule = append(stakingCredentialsRule, stakingCredentialsItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Redirect", fromPubKeyRule, toPubKeyRule, stakingCredentialsRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingRedirect)
				if err := _Staking.contract.UnpackLog(event, "Redirect", log); err != nil {
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

// ParseRedirect is a log parse operation binding the contract event 0xe161f5842757f257346b360594d094b7fa530f9404e93a80bf18bd8b14f9258f.
//
// Solidity: event Redirect(bytes indexed fromPubKey, bytes indexed toPubKey, bytes indexed stakingCredentials, uint64 amount)
func (_Staking *StakingFilterer) ParseRedirect(log types.Log) (*StakingRedirect, error) {
	event := new(StakingRedirect)
	if err := _Staking.contract.UnpackLog(event, "Redirect", log); err != nil {
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
	FromPubKey            common.Hash
	WithdrawalCredentials common.Hash
	Amount                uint64
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xd819a76a9128ab820538179b416ffb491e0fa0b23b2a08b605fba4c2649db9a6.
//
// Solidity: event Withdraw(bytes indexed fromPubKey, bytes indexed withdrawalCredentials, uint64 amount)
func (_Staking *StakingFilterer) FilterWithdraw(opts *bind.FilterOpts, fromPubKey [][]byte, withdrawalCredentials [][]byte) (*StakingWithdrawIterator, error) {

	var fromPubKeyRule []interface{}
	for _, fromPubKeyItem := range fromPubKey {
		fromPubKeyRule = append(fromPubKeyRule, fromPubKeyItem)
	}
	var withdrawalCredentialsRule []interface{}
	for _, withdrawalCredentialsItem := range withdrawalCredentials {
		withdrawalCredentialsRule = append(withdrawalCredentialsRule, withdrawalCredentialsItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Withdraw", fromPubKeyRule, withdrawalCredentialsRule)
	if err != nil {
		return nil, err
	}
	return &StakingWithdrawIterator{contract: _Staking.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xd819a76a9128ab820538179b416ffb491e0fa0b23b2a08b605fba4c2649db9a6.
//
// Solidity: event Withdraw(bytes indexed fromPubKey, bytes indexed withdrawalCredentials, uint64 amount)
func (_Staking *StakingFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *StakingWithdraw, fromPubKey [][]byte, withdrawalCredentials [][]byte) (event.Subscription, error) {

	var fromPubKeyRule []interface{}
	for _, fromPubKeyItem := range fromPubKey {
		fromPubKeyRule = append(fromPubKeyRule, fromPubKeyItem)
	}
	var withdrawalCredentialsRule []interface{}
	for _, withdrawalCredentialsItem := range withdrawalCredentials {
		withdrawalCredentialsRule = append(withdrawalCredentialsRule, withdrawalCredentialsItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Withdraw", fromPubKeyRule, withdrawalCredentialsRule)
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

// ParseWithdraw is a log parse operation binding the contract event 0xd819a76a9128ab820538179b416ffb491e0fa0b23b2a08b605fba4c2649db9a6.
//
// Solidity: event Withdraw(bytes indexed fromPubKey, bytes indexed withdrawalCredentials, uint64 amount)
func (_Staking *StakingFilterer) ParseWithdraw(log types.Log) (*StakingWithdraw, error) {
	event := new(StakingWithdraw)
	if err := _Staking.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
