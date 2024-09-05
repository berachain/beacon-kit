// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package deposit

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
	ABI: "[{\"type\":\"function\",\"name\":\"allowDeposit\",\"inputs\":[{\"name\":\"depositor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelOwnershipHandover\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"completeOwnershipHandover\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"withdrawal_credentials\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositAuth\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depositCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"result\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownershipHandoverExpiresAt\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"requestOwnershipHandover\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"credentials\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"index\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipHandoverCanceled\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipHandoverRequested\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"oldOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DepositNotMultipleOfGwei\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DepositValueTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientDeposit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCredentialsLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPubKeyLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignatureLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NewOwnerIsZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoHandoverRequest\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedDeposit\",\"inputs\":[]}]",
	Bin: "0x60806040525f805460ff60401b19169055348015601a575f80fd5b50610cac806100285f395ff3fe6080604052600436106100ce575f3560e01c80635b70fa291161007c5780638da5cb5b116100575780638da5cb5b146101de578063f04e283e14610231578063f2fde38b14610244578063fee81cf414610257575f80fd5b80635b70fa29146101a4578063715018a6146101b75780638c5f36bb146101bf575f80fd5b80633198a6b8116100ac5780633198a6b81461014857806354d1f13d1461017d5780635a7517ad14610185575f80fd5b806301ffc9a7146100d257806325692962146101065780632dfdf0b514610110575b5f80fd5b3480156100dd575f80fd5b506100f16100ec366004610947565b610296565b60405190151581526020015b60405180910390f35b61010e61032e565b005b34801561011b575f80fd5b505f5461012f9067ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100fd565b348015610153575f80fd5b5061012f6101623660046109b5565b60016020525f908152604090205467ffffffffffffffff1681565b61010e61037b565b348015610190575f80fd5b5061010e61019f3660046109e5565b6103b4565b61010e6101b2366004610a5b565b61041a565b61010e6104cf565b3480156101ca575f80fd5b5061010e6101d93660046109b5565b6104e2565b3480156101e9575f80fd5b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffff748739275460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100fd565b61010e61023f3660046109b5565b61059c565b61010e6102523660046109b5565b6105d9565b348015610262575f80fd5b506102886102713660046109b5565b63389a75e1600c9081525f91909152602090205490565b6040519081526020016100fd565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000148061032857507fffffffff0000000000000000000000000000000000000000000000000000000082167f5b70fa2900000000000000000000000000000000000000000000000000000000145b92915050565b5f6202a30067ffffffffffffffff164201905063389a75e1600c52335f52806020600c2055337fdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d5f80a250565b63389a75e1600c52335f525f6020600c2055337ffa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c925f80a2565b6103bc6105ff565b73ffffffffffffffffffffffffffffffffffffffff919091165f90815260016020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff909216919091179055565b335f9081526001602052604081205467ffffffffffffffff16900361046b576040517fce7ccd9600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b335f90815260016020526040812080549091906104919067ffffffffffffffff16610b0b565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506104c687878787878787610634565b50505050505050565b6104d76105ff565b6104e05f6107c4565b565b5f5468010000000000000000900460ff161561055e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f416c726561647920696e697469616c697a656400000000000000000000000000604482015260640160405180910390fd5b61056781610829565b505f80547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff1668010000000000000000179055565b6105a46105ff565b63389a75e1600c52805f526020600c2080544211156105ca57636f5e88185f526004601cfd5b5f90556105d6816107c4565b50565b6105e16105ff565b8060601b6105f657637448fbae5f526004601cfd5b6105d6816107c4565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffff748739275433146104e0576382b429005f526004601cfd5b6030861461066e576040517f9f10647200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602084146106a8576040517fb39bca1600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606081146106e2576040517f4be6321b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6106ec8461088c565b905064077359400067ffffffffffffffff82161015610737576040517f0e1eddda00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f80547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000008116600167ffffffffffffffff928316908101909216179091556040517f68af751683498a9f9be59fe8b0d52a64dd155255d85cdb29fea30b1e3f891d46916107b2918b918b918b918b9188918b918b9190610bb8565b60405180910390a15050505050505050565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927805473ffffffffffffffffffffffffffffffffffffffff9092169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a355565b73ffffffffffffffffffffffffffffffffffffffff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927819055805f7f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08180a350565b5f61089b633b9aca0034610c50565b156108d2576040517f40567b3800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6108e1633b9aca0034610c63565b905067ffffffffffffffff811115610925576040517f2aa6673400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103285f345f385f3884865af16109435763b12d13eb5f526004601cfd5b5050565b5f60208284031215610957575f80fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610986575f80fd5b9392505050565b803573ffffffffffffffffffffffffffffffffffffffff811681146109b0575f80fd5b919050565b5f602082840312156109c5575f80fd5b6109868261098d565b803567ffffffffffffffff811681146109b0575f80fd5b5f80604083850312156109f6575f80fd5b6109ff8361098d565b9150610a0d602084016109ce565b90509250929050565b5f8083601f840112610a26575f80fd5b50813567ffffffffffffffff811115610a3d575f80fd5b602083019150836020828501011115610a54575f80fd5b9250929050565b5f805f805f805f6080888a031215610a71575f80fd5b873567ffffffffffffffff811115610a87575f80fd5b610a938a828b01610a16565b909850965050602088013567ffffffffffffffff811115610ab2575f80fd5b610abe8a828b01610a16565b9096509450610ad19050604089016109ce565b9250606088013567ffffffffffffffff811115610aec575f80fd5b610af88a828b01610a16565b989b979a50959850939692959293505050565b5f67ffffffffffffffff821680610b49577f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0192915050565b81835281816020850137505f602082840101525f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60a081525f610bcb60a083018a8c610b71565b8281036020840152610bde81898b610b71565b905067ffffffffffffffff871660408401528281036060840152610c03818688610b71565b91505067ffffffffffffffff831660808301529998505050505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f82610c5e57610c5e610c23565b500690565b5f82610c7157610c71610c23565b50049056fea26469706673582212205b3330c66c9ba2a1794047e550acba3acfe1954245109488435a8dcab9112dc964736f6c634300081a0033",
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

// DepositAuth is a free data retrieval call binding the contract method 0x3198a6b8.
//
// Solidity: function depositAuth(address ) view returns(uint64)
func (_BeaconDepositContract *BeaconDepositContractCaller) DepositAuth(opts *bind.CallOpts, arg0 common.Address) (uint64, error) {
	var out []interface{}
	err := _BeaconDepositContract.contract.Call(opts, &out, "depositAuth", arg0)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// DepositAuth is a free data retrieval call binding the contract method 0x3198a6b8.
//
// Solidity: function depositAuth(address ) view returns(uint64)
func (_BeaconDepositContract *BeaconDepositContractSession) DepositAuth(arg0 common.Address) (uint64, error) {
	return _BeaconDepositContract.Contract.DepositAuth(&_BeaconDepositContract.CallOpts, arg0)
}

// DepositAuth is a free data retrieval call binding the contract method 0x3198a6b8.
//
// Solidity: function depositAuth(address ) view returns(uint64)
func (_BeaconDepositContract *BeaconDepositContractCallerSession) DepositAuth(arg0 common.Address) (uint64, error) {
	return _BeaconDepositContract.Contract.DepositAuth(&_BeaconDepositContract.CallOpts, arg0)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint64)
func (_BeaconDepositContract *BeaconDepositContractCaller) DepositCount(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _BeaconDepositContract.contract.Call(opts, &out, "depositCount")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint64)
func (_BeaconDepositContract *BeaconDepositContractSession) DepositCount() (uint64, error) {
	return _BeaconDepositContract.Contract.DepositCount(&_BeaconDepositContract.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint64)
func (_BeaconDepositContract *BeaconDepositContractCallerSession) DepositCount() (uint64, error) {
	return _BeaconDepositContract.Contract.DepositCount(&_BeaconDepositContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_BeaconDepositContract *BeaconDepositContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BeaconDepositContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_BeaconDepositContract *BeaconDepositContractSession) Owner() (common.Address, error) {
	return _BeaconDepositContract.Contract.Owner(&_BeaconDepositContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_BeaconDepositContract *BeaconDepositContractCallerSession) Owner() (common.Address, error) {
	return _BeaconDepositContract.Contract.Owner(&_BeaconDepositContract.CallOpts)
}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_BeaconDepositContract *BeaconDepositContractCaller) OwnershipHandoverExpiresAt(opts *bind.CallOpts, pendingOwner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BeaconDepositContract.contract.Call(opts, &out, "ownershipHandoverExpiresAt", pendingOwner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_BeaconDepositContract *BeaconDepositContractSession) OwnershipHandoverExpiresAt(pendingOwner common.Address) (*big.Int, error) {
	return _BeaconDepositContract.Contract.OwnershipHandoverExpiresAt(&_BeaconDepositContract.CallOpts, pendingOwner)
}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_BeaconDepositContract *BeaconDepositContractCallerSession) OwnershipHandoverExpiresAt(pendingOwner common.Address) (*big.Int, error) {
	return _BeaconDepositContract.Contract.OwnershipHandoverExpiresAt(&_BeaconDepositContract.CallOpts, pendingOwner)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_BeaconDepositContract *BeaconDepositContractCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BeaconDepositContract.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_BeaconDepositContract *BeaconDepositContractSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BeaconDepositContract.Contract.SupportsInterface(&_BeaconDepositContract.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_BeaconDepositContract *BeaconDepositContractCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BeaconDepositContract.Contract.SupportsInterface(&_BeaconDepositContract.CallOpts, interfaceId)
}

// AllowDeposit is a paid mutator transaction binding the contract method 0x5a7517ad.
//
// Solidity: function allowDeposit(address depositor, uint64 number) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) AllowDeposit(opts *bind.TransactOpts, depositor common.Address, number uint64) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "allowDeposit", depositor, number)
}

// AllowDeposit is a paid mutator transaction binding the contract method 0x5a7517ad.
//
// Solidity: function allowDeposit(address depositor, uint64 number) returns()
func (_BeaconDepositContract *BeaconDepositContractSession) AllowDeposit(depositor common.Address, number uint64) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.AllowDeposit(&_BeaconDepositContract.TransactOpts, depositor, number)
}

// AllowDeposit is a paid mutator transaction binding the contract method 0x5a7517ad.
//
// Solidity: function allowDeposit(address depositor, uint64 number) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) AllowDeposit(depositor common.Address, number uint64) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.AllowDeposit(&_BeaconDepositContract.TransactOpts, depositor, number)
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) CancelOwnershipHandover(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "cancelOwnershipHandover")
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_BeaconDepositContract *BeaconDepositContractSession) CancelOwnershipHandover() (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.CancelOwnershipHandover(&_BeaconDepositContract.TransactOpts)
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) CancelOwnershipHandover() (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.CancelOwnershipHandover(&_BeaconDepositContract.TransactOpts)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) CompleteOwnershipHandover(opts *bind.TransactOpts, pendingOwner common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "completeOwnershipHandover", pendingOwner)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_BeaconDepositContract *BeaconDepositContractSession) CompleteOwnershipHandover(pendingOwner common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.CompleteOwnershipHandover(&_BeaconDepositContract.TransactOpts, pendingOwner)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) CompleteOwnershipHandover(pendingOwner common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.CompleteOwnershipHandover(&_BeaconDepositContract.TransactOpts, pendingOwner)
}

// Deposit is a paid mutator transaction binding the contract method 0x5b70fa29.
//
// Solidity: function deposit(bytes pubkey, bytes withdrawal_credentials, uint64 amount, bytes signature) payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) Deposit(opts *bind.TransactOpts, pubkey []byte, withdrawal_credentials []byte, amount uint64, signature []byte) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "deposit", pubkey, withdrawal_credentials, amount, signature)
}

// Deposit is a paid mutator transaction binding the contract method 0x5b70fa29.
//
// Solidity: function deposit(bytes pubkey, bytes withdrawal_credentials, uint64 amount, bytes signature) payable returns()
func (_BeaconDepositContract *BeaconDepositContractSession) Deposit(pubkey []byte, withdrawal_credentials []byte, amount uint64, signature []byte) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.Deposit(&_BeaconDepositContract.TransactOpts, pubkey, withdrawal_credentials, amount, signature)
}

// Deposit is a paid mutator transaction binding the contract method 0x5b70fa29.
//
// Solidity: function deposit(bytes pubkey, bytes withdrawal_credentials, uint64 amount, bytes signature) payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) Deposit(pubkey []byte, withdrawal_credentials []byte, amount uint64, signature []byte) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.Deposit(&_BeaconDepositContract.TransactOpts, pubkey, withdrawal_credentials, amount, signature)
}

// InitializeOwner is a paid mutator transaction binding the contract method 0x8c5f36bb.
//
// Solidity: function initializeOwner(address owner) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) InitializeOwner(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "initializeOwner", owner)
}

// InitializeOwner is a paid mutator transaction binding the contract method 0x8c5f36bb.
//
// Solidity: function initializeOwner(address owner) returns()
func (_BeaconDepositContract *BeaconDepositContractSession) InitializeOwner(owner common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.InitializeOwner(&_BeaconDepositContract.TransactOpts, owner)
}

// InitializeOwner is a paid mutator transaction binding the contract method 0x8c5f36bb.
//
// Solidity: function initializeOwner(address owner) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) InitializeOwner(owner common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.InitializeOwner(&_BeaconDepositContract.TransactOpts, owner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_BeaconDepositContract *BeaconDepositContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.RenounceOwnership(&_BeaconDepositContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.RenounceOwnership(&_BeaconDepositContract.TransactOpts)
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) RequestOwnershipHandover(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "requestOwnershipHandover")
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_BeaconDepositContract *BeaconDepositContractSession) RequestOwnershipHandover() (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.RequestOwnershipHandover(&_BeaconDepositContract.TransactOpts)
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) RequestOwnershipHandover() (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.RequestOwnershipHandover(&_BeaconDepositContract.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_BeaconDepositContract *BeaconDepositContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.TransferOwnership(&_BeaconDepositContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.TransferOwnership(&_BeaconDepositContract.TransactOpts, newOwner)
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
	Pubkey      []byte
	Credentials []byte
	Amount      uint64
	Signature   []byte
	Index       uint64
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x68af751683498a9f9be59fe8b0d52a64dd155255d85cdb29fea30b1e3f891d46.
//
// Solidity: event Deposit(bytes pubkey, bytes credentials, uint64 amount, bytes signature, uint64 index)
func (_BeaconDepositContract *BeaconDepositContractFilterer) FilterDeposit(opts *bind.FilterOpts) (*BeaconDepositContractDepositIterator, error) {

	logs, sub, err := _BeaconDepositContract.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractDepositIterator{contract: _BeaconDepositContract.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x68af751683498a9f9be59fe8b0d52a64dd155255d85cdb29fea30b1e3f891d46.
//
// Solidity: event Deposit(bytes pubkey, bytes credentials, uint64 amount, bytes signature, uint64 index)
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

// ParseDeposit is a log parse operation binding the contract event 0x68af751683498a9f9be59fe8b0d52a64dd155255d85cdb29fea30b1e3f891d46.
//
// Solidity: event Deposit(bytes pubkey, bytes credentials, uint64 amount, bytes signature, uint64 index)
func (_BeaconDepositContract *BeaconDepositContractFilterer) ParseDeposit(log types.Log) (*BeaconDepositContractDeposit, error) {
	event := new(BeaconDepositContractDeposit)
	if err := _BeaconDepositContract.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconDepositContractOwnershipHandoverCanceledIterator is returned from FilterOwnershipHandoverCanceled and is used to iterate over the raw logs and unpacked data for OwnershipHandoverCanceled events raised by the BeaconDepositContract contract.
type BeaconDepositContractOwnershipHandoverCanceledIterator struct {
	Event *BeaconDepositContractOwnershipHandoverCanceled // Event containing the contract specifics and raw log

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
func (it *BeaconDepositContractOwnershipHandoverCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconDepositContractOwnershipHandoverCanceled)
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
		it.Event = new(BeaconDepositContractOwnershipHandoverCanceled)
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
func (it *BeaconDepositContractOwnershipHandoverCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconDepositContractOwnershipHandoverCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconDepositContractOwnershipHandoverCanceled represents a OwnershipHandoverCanceled event raised by the BeaconDepositContract contract.
type BeaconDepositContractOwnershipHandoverCanceled struct {
	PendingOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOwnershipHandoverCanceled is a free log retrieval operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_BeaconDepositContract *BeaconDepositContractFilterer) FilterOwnershipHandoverCanceled(opts *bind.FilterOpts, pendingOwner []common.Address) (*BeaconDepositContractOwnershipHandoverCanceledIterator, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _BeaconDepositContract.contract.FilterLogs(opts, "OwnershipHandoverCanceled", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractOwnershipHandoverCanceledIterator{contract: _BeaconDepositContract.contract, event: "OwnershipHandoverCanceled", logs: logs, sub: sub}, nil
}

// WatchOwnershipHandoverCanceled is a free log subscription operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_BeaconDepositContract *BeaconDepositContractFilterer) WatchOwnershipHandoverCanceled(opts *bind.WatchOpts, sink chan<- *BeaconDepositContractOwnershipHandoverCanceled, pendingOwner []common.Address) (event.Subscription, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _BeaconDepositContract.contract.WatchLogs(opts, "OwnershipHandoverCanceled", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconDepositContractOwnershipHandoverCanceled)
				if err := _BeaconDepositContract.contract.UnpackLog(event, "OwnershipHandoverCanceled", log); err != nil {
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

// ParseOwnershipHandoverCanceled is a log parse operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_BeaconDepositContract *BeaconDepositContractFilterer) ParseOwnershipHandoverCanceled(log types.Log) (*BeaconDepositContractOwnershipHandoverCanceled, error) {
	event := new(BeaconDepositContractOwnershipHandoverCanceled)
	if err := _BeaconDepositContract.contract.UnpackLog(event, "OwnershipHandoverCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconDepositContractOwnershipHandoverRequestedIterator is returned from FilterOwnershipHandoverRequested and is used to iterate over the raw logs and unpacked data for OwnershipHandoverRequested events raised by the BeaconDepositContract contract.
type BeaconDepositContractOwnershipHandoverRequestedIterator struct {
	Event *BeaconDepositContractOwnershipHandoverRequested // Event containing the contract specifics and raw log

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
func (it *BeaconDepositContractOwnershipHandoverRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconDepositContractOwnershipHandoverRequested)
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
		it.Event = new(BeaconDepositContractOwnershipHandoverRequested)
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
func (it *BeaconDepositContractOwnershipHandoverRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconDepositContractOwnershipHandoverRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconDepositContractOwnershipHandoverRequested represents a OwnershipHandoverRequested event raised by the BeaconDepositContract contract.
type BeaconDepositContractOwnershipHandoverRequested struct {
	PendingOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOwnershipHandoverRequested is a free log retrieval operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_BeaconDepositContract *BeaconDepositContractFilterer) FilterOwnershipHandoverRequested(opts *bind.FilterOpts, pendingOwner []common.Address) (*BeaconDepositContractOwnershipHandoverRequestedIterator, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _BeaconDepositContract.contract.FilterLogs(opts, "OwnershipHandoverRequested", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractOwnershipHandoverRequestedIterator{contract: _BeaconDepositContract.contract, event: "OwnershipHandoverRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipHandoverRequested is a free log subscription operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_BeaconDepositContract *BeaconDepositContractFilterer) WatchOwnershipHandoverRequested(opts *bind.WatchOpts, sink chan<- *BeaconDepositContractOwnershipHandoverRequested, pendingOwner []common.Address) (event.Subscription, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _BeaconDepositContract.contract.WatchLogs(opts, "OwnershipHandoverRequested", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconDepositContractOwnershipHandoverRequested)
				if err := _BeaconDepositContract.contract.UnpackLog(event, "OwnershipHandoverRequested", log); err != nil {
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

// ParseOwnershipHandoverRequested is a log parse operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_BeaconDepositContract *BeaconDepositContractFilterer) ParseOwnershipHandoverRequested(log types.Log) (*BeaconDepositContractOwnershipHandoverRequested, error) {
	event := new(BeaconDepositContractOwnershipHandoverRequested)
	if err := _BeaconDepositContract.contract.UnpackLog(event, "OwnershipHandoverRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconDepositContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BeaconDepositContract contract.
type BeaconDepositContractOwnershipTransferredIterator struct {
	Event *BeaconDepositContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BeaconDepositContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconDepositContractOwnershipTransferred)
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
		it.Event = new(BeaconDepositContractOwnershipTransferred)
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
func (it *BeaconDepositContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconDepositContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconDepositContractOwnershipTransferred represents a OwnershipTransferred event raised by the BeaconDepositContract contract.
type BeaconDepositContractOwnershipTransferred struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_BeaconDepositContract *BeaconDepositContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, oldOwner []common.Address, newOwner []common.Address) (*BeaconDepositContractOwnershipTransferredIterator, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BeaconDepositContract.contract.FilterLogs(opts, "OwnershipTransferred", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractOwnershipTransferredIterator{contract: _BeaconDepositContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_BeaconDepositContract *BeaconDepositContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BeaconDepositContractOwnershipTransferred, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BeaconDepositContract.contract.WatchLogs(opts, "OwnershipTransferred", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconDepositContractOwnershipTransferred)
				if err := _BeaconDepositContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_BeaconDepositContract *BeaconDepositContractFilterer) ParseOwnershipTransferred(log types.Log) (*BeaconDepositContractOwnershipTransferred, error) {
	event := new(BeaconDepositContractOwnershipTransferred)
	if err := _BeaconDepositContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
