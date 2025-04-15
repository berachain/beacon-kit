// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ssztest

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

// SSZTestMetaData contains all meta data concerning the SSZTest contract.
var SSZTestMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"mustVerifyProof\",\"inputs\":[{\"name\":\"proof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leaf\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyProof\",\"inputs\":[{\"name\":\"proof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"leaf\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"isValid\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x6080604052348015600e575f80fd5b5061025f8061001c5f395ff3fe608060405234801561000f575f80fd5b5060043610610034575f3560e01c806341b703ff146100385780634fc36be61461004d575b5f80fd5b61004b6100463660046101a4565b610074565b005b61006061005b3660046101a4565b6100f2565b604051901515815260200160405180910390f35b610081858585858561010a565b6100eb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f50726f6f6620697320696e76616c696400000000000000000000000000000000604482015260640160405180910390fd5b5050505050565b5f610100868686868661010a565b9695505050505050565b5f8415610166578460051b8601865b6001841660051b8460011c94508461013857635849603f5f526004601cfd5b85815281356020918218525f60408160025afa80610154575f80fd5b505f5194506020018181106101195750505b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82011561019b57631b6661c35f526004601cfd5b50501492915050565b5f805f805f608086880312156101b8575f80fd5b853567ffffffffffffffff8111156101ce575f80fd5b8601601f810188136101de575f80fd5b803567ffffffffffffffff8111156101f4575f80fd5b8860208260051b8401011115610208575f80fd5b6020918201999098509087013596604081013596506060013594509250505056fea26469706673582212208778da1308154d8fae10a8ae616eb37887fec9be573397aab1d43ed15bb221ab64736f6c634300081a0033",
}

// SSZTestABI is the input ABI used to generate the binding from.
// Deprecated: Use SSZTestMetaData.ABI instead.
var SSZTestABI = SSZTestMetaData.ABI

// SSZTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SSZTestMetaData.Bin instead.
var SSZTestBin = SSZTestMetaData.Bin

// DeploySSZTest deploys a new Ethereum contract, binding an instance of SSZTest to it.
func DeploySSZTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SSZTest, error) {
	parsed, err := SSZTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SSZTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SSZTest{SSZTestCaller: SSZTestCaller{contract: contract}, SSZTestTransactor: SSZTestTransactor{contract: contract}, SSZTestFilterer: SSZTestFilterer{contract: contract}}, nil
}

// SSZTest is an auto generated Go binding around an Ethereum contract.
type SSZTest struct {
	SSZTestCaller     // Read-only binding to the contract
	SSZTestTransactor // Write-only binding to the contract
	SSZTestFilterer   // Log filterer for contract events
}

// SSZTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type SSZTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SSZTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SSZTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SSZTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SSZTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SSZTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SSZTestSession struct {
	Contract     *SSZTest          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SSZTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SSZTestCallerSession struct {
	Contract *SSZTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// SSZTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SSZTestTransactorSession struct {
	Contract     *SSZTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// SSZTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type SSZTestRaw struct {
	Contract *SSZTest // Generic contract binding to access the raw methods on
}

// SSZTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SSZTestCallerRaw struct {
	Contract *SSZTestCaller // Generic read-only contract binding to access the raw methods on
}

// SSZTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SSZTestTransactorRaw struct {
	Contract *SSZTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSSZTest creates a new instance of SSZTest, bound to a specific deployed contract.
func NewSSZTest(address common.Address, backend bind.ContractBackend) (*SSZTest, error) {
	contract, err := bindSSZTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SSZTest{SSZTestCaller: SSZTestCaller{contract: contract}, SSZTestTransactor: SSZTestTransactor{contract: contract}, SSZTestFilterer: SSZTestFilterer{contract: contract}}, nil
}

// NewSSZTestCaller creates a new read-only instance of SSZTest, bound to a specific deployed contract.
func NewSSZTestCaller(address common.Address, caller bind.ContractCaller) (*SSZTestCaller, error) {
	contract, err := bindSSZTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SSZTestCaller{contract: contract}, nil
}

// NewSSZTestTransactor creates a new write-only instance of SSZTest, bound to a specific deployed contract.
func NewSSZTestTransactor(address common.Address, transactor bind.ContractTransactor) (*SSZTestTransactor, error) {
	contract, err := bindSSZTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SSZTestTransactor{contract: contract}, nil
}

// NewSSZTestFilterer creates a new log filterer instance of SSZTest, bound to a specific deployed contract.
func NewSSZTestFilterer(address common.Address, filterer bind.ContractFilterer) (*SSZTestFilterer, error) {
	contract, err := bindSSZTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SSZTestFilterer{contract: contract}, nil
}

// bindSSZTest binds a generic wrapper to an already deployed contract.
func bindSSZTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SSZTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SSZTest *SSZTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SSZTest.Contract.SSZTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SSZTest *SSZTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SSZTest.Contract.SSZTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SSZTest *SSZTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SSZTest.Contract.SSZTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SSZTest *SSZTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SSZTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SSZTest *SSZTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SSZTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SSZTest *SSZTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SSZTest.Contract.contract.Transact(opts, method, params...)
}

// MustVerifyProof is a free data retrieval call binding the contract method 0x41b703ff.
//
// Solidity: function mustVerifyProof(bytes32[] proof, bytes32 root, bytes32 leaf, uint256 index) view returns()
func (_SSZTest *SSZTestCaller) MustVerifyProof(opts *bind.CallOpts, proof [][32]byte, root [32]byte, leaf [32]byte, index *big.Int) error {
	var out []interface{}
	err := _SSZTest.contract.Call(opts, &out, "mustVerifyProof", proof, root, leaf, index)

	if err != nil {
		return err
	}

	return err

}

// MustVerifyProof is a free data retrieval call binding the contract method 0x41b703ff.
//
// Solidity: function mustVerifyProof(bytes32[] proof, bytes32 root, bytes32 leaf, uint256 index) view returns()
func (_SSZTest *SSZTestSession) MustVerifyProof(proof [][32]byte, root [32]byte, leaf [32]byte, index *big.Int) error {
	return _SSZTest.Contract.MustVerifyProof(&_SSZTest.CallOpts, proof, root, leaf, index)
}

// MustVerifyProof is a free data retrieval call binding the contract method 0x41b703ff.
//
// Solidity: function mustVerifyProof(bytes32[] proof, bytes32 root, bytes32 leaf, uint256 index) view returns()
func (_SSZTest *SSZTestCallerSession) MustVerifyProof(proof [][32]byte, root [32]byte, leaf [32]byte, index *big.Int) error {
	return _SSZTest.Contract.MustVerifyProof(&_SSZTest.CallOpts, proof, root, leaf, index)
}

// VerifyProof is a free data retrieval call binding the contract method 0x4fc36be6.
//
// Solidity: function verifyProof(bytes32[] proof, bytes32 root, bytes32 leaf, uint256 index) view returns(bool isValid)
func (_SSZTest *SSZTestCaller) VerifyProof(opts *bind.CallOpts, proof [][32]byte, root [32]byte, leaf [32]byte, index *big.Int) (bool, error) {
	var out []interface{}
	err := _SSZTest.contract.Call(opts, &out, "verifyProof", proof, root, leaf, index)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyProof is a free data retrieval call binding the contract method 0x4fc36be6.
//
// Solidity: function verifyProof(bytes32[] proof, bytes32 root, bytes32 leaf, uint256 index) view returns(bool isValid)
func (_SSZTest *SSZTestSession) VerifyProof(proof [][32]byte, root [32]byte, leaf [32]byte, index *big.Int) (bool, error) {
	return _SSZTest.Contract.VerifyProof(&_SSZTest.CallOpts, proof, root, leaf, index)
}

// VerifyProof is a free data retrieval call binding the contract method 0x4fc36be6.
//
// Solidity: function verifyProof(bytes32[] proof, bytes32 root, bytes32 leaf, uint256 index) view returns(bool isValid)
func (_SSZTest *SSZTestCallerSession) VerifyProof(proof [][32]byte, root [32]byte, leaf [32]byte, index *big.Int) (bool, error) {
	return _SSZTest.Contract.VerifyProof(&_SSZTest.CallOpts, proof, root, leaf, index)
}
