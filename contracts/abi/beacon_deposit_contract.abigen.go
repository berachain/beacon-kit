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

// BeaconDepositContractMetaData contains all meta data concerning the BeaconDepositContract contract.
var BeaconDepositContractMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"NATIVE_ASSET\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"stakingCredentials\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"redirect\",\"inputs\":[{\"name\":\"fromPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"toPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"withdrawalCredentials\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"stakingCredentials\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Redirect\",\"inputs\":[{\"name\":\"fromPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"toPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"stakingCredentials\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"fromPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"withdrawalCredentials\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"DepositNotMultipleOfGwei\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DepositValueTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientDeposit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientRedirectAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientWithdrawAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCredentialsLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPubKeyLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignatureLength\",\"inputs\":[]}]",
	Bin: "0x60806040525f80546001600160a01b03191673eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee179055348015610034575f80fd5b50610a89806100425f395ff3fe60806040526004361061003e575f3560e01c806304a3267f146100425780635b70fa2914610063578063bf53253b14610076578063bf9b6a55146100c6575b5f80fd5b34801561004d575f80fd5b5061006161005c36600461071e565b6100e5565b005b610061610071366004610799565b6101fc565b348015610081575f80fd5b5061009d73eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee81565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b3480156100d1575f80fd5b506100616100e036600461071e565b61034c565b6030841461011f576040517f9f10647200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208214610159576040517fb39bca1600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610169600a64077359400061086c565b67ffffffffffffffff168167ffffffffffffffff1610156101b6576040517febec602100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7fd819a76a9128ab820538179b416ffb491e0fa0b23b2a08b605fba4c2649db9a685858585856040516101ed9594939291906108d9565b60405180910390a15050505050565b60308614610236576040517f9f10647200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208414610270576040517fb39bca1600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606081146102aa576040517f4be6321b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f5473ffffffffffffffffffffffffffffffffffffffff167fffffffffffffffffffffffff1111111111111111111111111111111111111112016102f7576102f06104a1565b9250610300565b61030083610577565b7f1f39b85dd1a529b31e0cd61e5609e1feca0e08e2103fe319fbd3dd5a0c7b68df8787878787878760405161033b979695949392919061091c565b60405180910390a150505050505050565b60308414158061035d575060308214155b15610394576040517f9f10647200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103a4600a64077359400061086c565b67ffffffffffffffff168167ffffffffffffffff1610156103f1576040517f0494a69c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7fe161f5842757f257346b360594d094b7fa530f9404e93a80bf18bd8b14f9258f8585858561048e33604080517f010000000000000000000000000000000000000000000000000000000000000060208201525f6021820152606083811b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016602c83015291016040516020818303038152906040529050919050565b866040516101ed96959493929190610975565b5f67ffffffffffffffff3411156104e4576040517f2aa6673400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b640773594000341015610523576040517f0e1eddda00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610531633b9aca0034610a1a565b15610568576040517f40567b3800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105725f346106a0565b503490565b5f546040517f9dc29fac00000000000000000000000000000000000000000000000000000000815233600482015267ffffffffffffffff8316602482015273ffffffffffffffffffffffffffffffffffffffff90911690639dc29fac906044015f604051808303815f87803b1580156105ee575f80fd5b505af1158015610600573d5f803e3d5ffd5b50505064077359400067ffffffffffffffff83161015905061064e576040517f0e1eddda00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61065c633b9aca0082610a2d565b67ffffffffffffffff161561069d576040517f40567b3800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50565b5f385f3884865af16106b95763b12d13eb5f526004601cfd5b5050565b5f8083601f8401126106cd575f80fd5b50813567ffffffffffffffff8111156106e4575f80fd5b6020830191508360208285010111156106fb575f80fd5b9250929050565b803567ffffffffffffffff81168114610719575f80fd5b919050565b5f805f805f60608688031215610732575f80fd5b853567ffffffffffffffff80821115610749575f80fd5b61075589838a016106bd565b9097509550602088013591508082111561076d575f80fd5b5061077a888289016106bd565b909450925061078d905060408701610702565b90509295509295909350565b5f805f805f805f6080888a0312156107af575f80fd5b873567ffffffffffffffff808211156107c6575f80fd5b6107d28b838c016106bd565b909950975060208a01359150808211156107ea575f80fd5b6107f68b838c016106bd565b909750955085915061080a60408b01610702565b945060608a013591508082111561081f575f80fd5b5061082c8a828b016106bd565b989b979a50959850939692959293505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f67ffffffffffffffff808416806108865761088661083f565b92169190910492915050565b81835281816020850137505f602082840101525f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b606081525f6108ec606083018789610892565b82810360208401526108ff818688610892565b91505067ffffffffffffffff831660408301529695505050505050565b608081525f61092f60808301898b610892565b828103602084015261094281888a610892565b905067ffffffffffffffff861660408401528281036060840152610967818587610892565b9a9950505050505050505050565b608081525f61098860808301888a610892565b60208382038185015261099c82888a610892565b9150838203604085015285518083525f5b818110156109c85787810183015184820184015282016109ad565b505f828285010152817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116840101935050505067ffffffffffffffff83166060830152979650505050505050565b5f82610a2857610a2861083f565b500690565b5f67ffffffffffffffff80841680610a4757610a4761083f565b9216919091069291505056fea264697066735822122037c92cbe359e79eb3b77b23a0f63aafa61719d424026007a4c7591dc77d3d95e64736f6c63430008180033",
}

// BeaconDepositContractABI is the input ABI used to generate the binding from.
// Deprecated: Use BeaconDepositContractMetaData.ABI instead.
var BeaconDepositContractABI = BeaconDepositContractMetaData.ABI

// BeaconDepositContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BeaconDepositContractMetaData.Bin instead.
var BeaconDepositContractBin = BeaconDepositContractMetaData.Bin

// DeployBeaconDepositContract deploys a new Ethereum contract, binding an instance of BeaconDepositContract to it.
func DeployBeaconDepositContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BeaconDepositContract, error) {
	parsed, err := BeaconDepositContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BeaconDepositContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BeaconDepositContract{BeaconDepositContractCaller: BeaconDepositContractCaller{contract: contract}, BeaconDepositContractTransactor: BeaconDepositContractTransactor{contract: contract}, BeaconDepositContractFilterer: BeaconDepositContractFilterer{contract: contract}}, nil
}

// BeaconDepositContract is an auto generated Go binding around an Ethereum contract.
type BeaconDepositContract struct {
	BeaconDepositContractCaller     // Read-only binding to the contract
	BeaconDepositContractTransactor // Write-only binding to the contract
	BeaconDepositContractFilterer   // Log filterer for contract events
}

// BeaconDepositContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type BeaconDepositContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconDepositContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BeaconDepositContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconDepositContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BeaconDepositContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconDepositContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BeaconDepositContractSession struct {
	Contract     *BeaconDepositContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BeaconDepositContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BeaconDepositContractCallerSession struct {
	Contract *BeaconDepositContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// BeaconDepositContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BeaconDepositContractTransactorSession struct {
	Contract     *BeaconDepositContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// BeaconDepositContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type BeaconDepositContractRaw struct {
	Contract *BeaconDepositContract // Generic contract binding to access the raw methods on
}

// BeaconDepositContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BeaconDepositContractCallerRaw struct {
	Contract *BeaconDepositContractCaller // Generic read-only contract binding to access the raw methods on
}

// BeaconDepositContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BeaconDepositContractTransactorRaw struct {
	Contract *BeaconDepositContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBeaconDepositContract creates a new instance of BeaconDepositContract, bound to a specific deployed contract.
func NewBeaconDepositContract(address common.Address, backend bind.ContractBackend) (*BeaconDepositContract, error) {
	contract, err := bindBeaconDepositContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContract{BeaconDepositContractCaller: BeaconDepositContractCaller{contract: contract}, BeaconDepositContractTransactor: BeaconDepositContractTransactor{contract: contract}, BeaconDepositContractFilterer: BeaconDepositContractFilterer{contract: contract}}, nil
}

// NewBeaconDepositContractCaller creates a new read-only instance of BeaconDepositContract, bound to a specific deployed contract.
func NewBeaconDepositContractCaller(address common.Address, caller bind.ContractCaller) (*BeaconDepositContractCaller, error) {
	contract, err := bindBeaconDepositContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractCaller{contract: contract}, nil
}

// NewBeaconDepositContractTransactor creates a new write-only instance of BeaconDepositContract, bound to a specific deployed contract.
func NewBeaconDepositContractTransactor(address common.Address, transactor bind.ContractTransactor) (*BeaconDepositContractTransactor, error) {
	contract, err := bindBeaconDepositContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractTransactor{contract: contract}, nil
}

// NewBeaconDepositContractFilterer creates a new log filterer instance of BeaconDepositContract, bound to a specific deployed contract.
func NewBeaconDepositContractFilterer(address common.Address, filterer bind.ContractFilterer) (*BeaconDepositContractFilterer, error) {
	contract, err := bindBeaconDepositContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractFilterer{contract: contract}, nil
}

// bindBeaconDepositContract binds a generic wrapper to an already deployed contract.
func bindBeaconDepositContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BeaconDepositContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconDepositContract *BeaconDepositContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconDepositContract.Contract.BeaconDepositContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconDepositContract *BeaconDepositContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.BeaconDepositContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconDepositContract *BeaconDepositContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.BeaconDepositContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconDepositContract *BeaconDepositContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconDepositContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconDepositContract *BeaconDepositContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconDepositContract *BeaconDepositContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.contract.Transact(opts, method, params...)
}

// NATIVEASSET is a free data retrieval call binding the contract method 0xbf53253b.
//
// Solidity: function NATIVE_ASSET() view returns(address)
func (_BeaconDepositContract *BeaconDepositContractCaller) NATIVEASSET(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BeaconDepositContract.contract.Call(opts, &out, "NATIVE_ASSET")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NATIVEASSET is a free data retrieval call binding the contract method 0xbf53253b.
//
// Solidity: function NATIVE_ASSET() view returns(address)
func (_BeaconDepositContract *BeaconDepositContractSession) NATIVEASSET() (common.Address, error) {
	return _BeaconDepositContract.Contract.NATIVEASSET(&_BeaconDepositContract.CallOpts)
}

// NATIVEASSET is a free data retrieval call binding the contract method 0xbf53253b.
//
// Solidity: function NATIVE_ASSET() view returns(address)
func (_BeaconDepositContract *BeaconDepositContractCallerSession) NATIVEASSET() (common.Address, error) {
	return _BeaconDepositContract.Contract.NATIVEASSET(&_BeaconDepositContract.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x5b70fa29.
//
// Solidity: function deposit(bytes validatorPubkey, bytes stakingCredentials, uint64 amount, bytes signature) payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) Deposit(opts *bind.TransactOpts, validatorPubkey []byte, stakingCredentials []byte, amount uint64, signature []byte) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "deposit", validatorPubkey, stakingCredentials, amount, signature)
}

// Deposit is a paid mutator transaction binding the contract method 0x5b70fa29.
//
// Solidity: function deposit(bytes validatorPubkey, bytes stakingCredentials, uint64 amount, bytes signature) payable returns()
func (_BeaconDepositContract *BeaconDepositContractSession) Deposit(validatorPubkey []byte, stakingCredentials []byte, amount uint64, signature []byte) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.Deposit(&_BeaconDepositContract.TransactOpts, validatorPubkey, stakingCredentials, amount, signature)
}

// Deposit is a paid mutator transaction binding the contract method 0x5b70fa29.
//
// Solidity: function deposit(bytes validatorPubkey, bytes stakingCredentials, uint64 amount, bytes signature) payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) Deposit(validatorPubkey []byte, stakingCredentials []byte, amount uint64, signature []byte) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.Deposit(&_BeaconDepositContract.TransactOpts, validatorPubkey, stakingCredentials, amount, signature)
}

// Redirect is a paid mutator transaction binding the contract method 0xbf9b6a55.
//
// Solidity: function redirect(bytes fromPubkey, bytes toPubkey, uint64 amount) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) Redirect(opts *bind.TransactOpts, fromPubkey []byte, toPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "redirect", fromPubkey, toPubkey, amount)
}

// Redirect is a paid mutator transaction binding the contract method 0xbf9b6a55.
//
// Solidity: function redirect(bytes fromPubkey, bytes toPubkey, uint64 amount) returns()
func (_BeaconDepositContract *BeaconDepositContractSession) Redirect(fromPubkey []byte, toPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.Redirect(&_BeaconDepositContract.TransactOpts, fromPubkey, toPubkey, amount)
}

// Redirect is a paid mutator transaction binding the contract method 0xbf9b6a55.
//
// Solidity: function redirect(bytes fromPubkey, bytes toPubkey, uint64 amount) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) Redirect(fromPubkey []byte, toPubkey []byte, amount uint64) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.Redirect(&_BeaconDepositContract.TransactOpts, fromPubkey, toPubkey, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x04a3267f.
//
// Solidity: function withdraw(bytes validatorPubkey, bytes withdrawalCredentials, uint64 amount) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) Withdraw(opts *bind.TransactOpts, validatorPubkey []byte, withdrawalCredentials []byte, amount uint64) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "withdraw", validatorPubkey, withdrawalCredentials, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x04a3267f.
//
// Solidity: function withdraw(bytes validatorPubkey, bytes withdrawalCredentials, uint64 amount) returns()
func (_BeaconDepositContract *BeaconDepositContractSession) Withdraw(validatorPubkey []byte, withdrawalCredentials []byte, amount uint64) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.Withdraw(&_BeaconDepositContract.TransactOpts, validatorPubkey, withdrawalCredentials, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x04a3267f.
//
// Solidity: function withdraw(bytes validatorPubkey, bytes withdrawalCredentials, uint64 amount) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) Withdraw(validatorPubkey []byte, withdrawalCredentials []byte, amount uint64) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.Withdraw(&_BeaconDepositContract.TransactOpts, validatorPubkey, withdrawalCredentials, amount)
}

// BeaconDepositContractDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the BeaconDepositContract contract.
type BeaconDepositContractDepositIterator struct {
	Event *BeaconDepositContractDeposit // Event containing the contract specifics and raw log

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
func (it *BeaconDepositContractDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconDepositContractDeposit)
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
		it.Event = new(BeaconDepositContractDeposit)
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
func (it *BeaconDepositContractDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconDepositContractDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconDepositContractDeposit represents a Deposit event raised by the BeaconDepositContract contract.
type BeaconDepositContractDeposit struct {
	ValidatorPubkey    []byte
	StakingCredentials []byte
	Amount             uint64
	Signature          []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x1f39b85dd1a529b31e0cd61e5609e1feca0e08e2103fe319fbd3dd5a0c7b68df.
//
// Solidity: event Deposit(bytes validatorPubkey, bytes stakingCredentials, uint64 amount, bytes signature)
func (_BeaconDepositContract *BeaconDepositContractFilterer) FilterDeposit(opts *bind.FilterOpts) (*BeaconDepositContractDepositIterator, error) {

	logs, sub, err := _BeaconDepositContract.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractDepositIterator{contract: _BeaconDepositContract.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x1f39b85dd1a529b31e0cd61e5609e1feca0e08e2103fe319fbd3dd5a0c7b68df.
//
// Solidity: event Deposit(bytes validatorPubkey, bytes stakingCredentials, uint64 amount, bytes signature)
func (_BeaconDepositContract *BeaconDepositContractFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *BeaconDepositContractDeposit) (event.Subscription, error) {

	logs, sub, err := _BeaconDepositContract.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconDepositContractDeposit)
				if err := _BeaconDepositContract.contract.UnpackLog(event, "Deposit", log); err != nil {
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
// Solidity: event Deposit(bytes validatorPubkey, bytes stakingCredentials, uint64 amount, bytes signature)
func (_BeaconDepositContract *BeaconDepositContractFilterer) ParseDeposit(log types.Log) (*BeaconDepositContractDeposit, error) {
	event := new(BeaconDepositContractDeposit)
	if err := _BeaconDepositContract.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconDepositContractRedirectIterator is returned from FilterRedirect and is used to iterate over the raw logs and unpacked data for Redirect events raised by the BeaconDepositContract contract.
type BeaconDepositContractRedirectIterator struct {
	Event *BeaconDepositContractRedirect // Event containing the contract specifics and raw log

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
func (it *BeaconDepositContractRedirectIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconDepositContractRedirect)
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
		it.Event = new(BeaconDepositContractRedirect)
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
func (it *BeaconDepositContractRedirectIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconDepositContractRedirectIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconDepositContractRedirect represents a Redirect event raised by the BeaconDepositContract contract.
type BeaconDepositContractRedirect struct {
	FromPubkey         []byte
	ToPubkey           []byte
	StakingCredentials []byte
	Amount             uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterRedirect is a free log retrieval operation binding the contract event 0xe161f5842757f257346b360594d094b7fa530f9404e93a80bf18bd8b14f9258f.
//
// Solidity: event Redirect(bytes fromPubkey, bytes toPubkey, bytes stakingCredentials, uint64 amount)
func (_BeaconDepositContract *BeaconDepositContractFilterer) FilterRedirect(opts *bind.FilterOpts) (*BeaconDepositContractRedirectIterator, error) {

	logs, sub, err := _BeaconDepositContract.contract.FilterLogs(opts, "Redirect")
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractRedirectIterator{contract: _BeaconDepositContract.contract, event: "Redirect", logs: logs, sub: sub}, nil
}

// WatchRedirect is a free log subscription operation binding the contract event 0xe161f5842757f257346b360594d094b7fa530f9404e93a80bf18bd8b14f9258f.
//
// Solidity: event Redirect(bytes fromPubkey, bytes toPubkey, bytes stakingCredentials, uint64 amount)
func (_BeaconDepositContract *BeaconDepositContractFilterer) WatchRedirect(opts *bind.WatchOpts, sink chan<- *BeaconDepositContractRedirect) (event.Subscription, error) {

	logs, sub, err := _BeaconDepositContract.contract.WatchLogs(opts, "Redirect")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconDepositContractRedirect)
				if err := _BeaconDepositContract.contract.UnpackLog(event, "Redirect", log); err != nil {
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
// Solidity: event Redirect(bytes fromPubkey, bytes toPubkey, bytes stakingCredentials, uint64 amount)
func (_BeaconDepositContract *BeaconDepositContractFilterer) ParseRedirect(log types.Log) (*BeaconDepositContractRedirect, error) {
	event := new(BeaconDepositContractRedirect)
	if err := _BeaconDepositContract.contract.UnpackLog(event, "Redirect", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconDepositContractWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the BeaconDepositContract contract.
type BeaconDepositContractWithdrawIterator struct {
	Event *BeaconDepositContractWithdraw // Event containing the contract specifics and raw log

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
func (it *BeaconDepositContractWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconDepositContractWithdraw)
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
		it.Event = new(BeaconDepositContractWithdraw)
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
func (it *BeaconDepositContractWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconDepositContractWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconDepositContractWithdraw represents a Withdraw event raised by the BeaconDepositContract contract.
type BeaconDepositContractWithdraw struct {
	FromPubkey            []byte
	WithdrawalCredentials []byte
	Amount                uint64
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xd819a76a9128ab820538179b416ffb491e0fa0b23b2a08b605fba4c2649db9a6.
//
// Solidity: event Withdraw(bytes fromPubkey, bytes withdrawalCredentials, uint64 amount)
func (_BeaconDepositContract *BeaconDepositContractFilterer) FilterWithdraw(opts *bind.FilterOpts) (*BeaconDepositContractWithdrawIterator, error) {

	logs, sub, err := _BeaconDepositContract.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractWithdrawIterator{contract: _BeaconDepositContract.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xd819a76a9128ab820538179b416ffb491e0fa0b23b2a08b605fba4c2649db9a6.
//
// Solidity: event Withdraw(bytes fromPubkey, bytes withdrawalCredentials, uint64 amount)
func (_BeaconDepositContract *BeaconDepositContractFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *BeaconDepositContractWithdraw) (event.Subscription, error) {

	logs, sub, err := _BeaconDepositContract.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconDepositContractWithdraw)
				if err := _BeaconDepositContract.contract.UnpackLog(event, "Withdraw", log); err != nil {
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
// Solidity: event Withdraw(bytes fromPubkey, bytes withdrawalCredentials, uint64 amount)
func (_BeaconDepositContract *BeaconDepositContractFilterer) ParseWithdraw(log types.Log) (*BeaconDepositContractWithdraw, error) {
	event := new(BeaconDepositContractWithdraw)
	if err := _BeaconDepositContract.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
