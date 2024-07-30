// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eip4788

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

// BeaconVerifierMetaData contains all meta data concerning the BeaconVerifier contract.
var BeaconVerifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_zeroValidatorPubkeyGIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_executionNumberGIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BEACON_ROOTS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelOwnershipHandover\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"completeOwnershipHandover\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"executionNumberGIndex\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParentBeaconBlockRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParentBeaconBlockRootAt\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"result\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownershipHandoverExpiresAt\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"requestOwnershipHandover\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setExecutionNumberGIndex\",\"inputs\":[{\"name\":\"_executionNumberGIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setZeroValidatorPubkeyGIndex\",\"inputs\":[{\"name\":\"_zeroValidatorPubkeyGIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"verifyBeaconBlockProposer\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"validatorPubkeyProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"proposerIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyExecutionNumber\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"executionNumberProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"zeroValidatorPubkeyGIndex\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ExecutionNumberGIndexChanged\",\"inputs\":[{\"name\":\"newExecutionNumberGIndex\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipHandoverCanceled\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipHandoverRequested\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"oldOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ZeroValidatorPubkeyGIndexChanged\",\"inputs\":[{\"name\":\"newZeroValidatorPubkeyGIndex\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidValidatorPubkeyLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NewOwnerIsZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoHandoverRequest\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RootNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
	Bin: "0x6080604052348015600e575f80fd5b50604051610bde380380610bde833981016040819052602b916098565b5f8290556001819055603b336041565b505060b9565b638b78c6d819805415605a57630dc149f05f526004601cfd5b6001600160a01b03909116801560ff1b8117909155805f7f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08180a350565b5f806040838503121560a8575f80fd5b505080516020909101519092909150565b610b18806100c65f395ff3fe6080604052600436106100ef575f3560e01c80638170577111610087578063f04e283e11610057578063f04e283e14610257578063f2fde38b1461026a578063f769afd11461027d578063fee81cf41461029c575f80fd5b806381705771146101d25780638da5cb5b146101f1578063ae9af3e314610224578063efcff00a14610243575f80fd5b806354d1f13d116100c257806354d1f13d1461015857806356d7e8fd146101605780635e67d452146101ab578063715018a6146101ca575f80fd5b806325692962146100f3578063305ea416146100fd57806335222ff114610125578063491920ab14610139575b5f80fd5b6100fb6102cd565b005b348015610108575f80fd5b5061011260015481565b6040519081526020015b60405180910390f35b348015610130575f80fd5b506101125f5481565b348015610144575f80fd5b506100fb6101533660046108f0565b61031a565b6100fb610338565b34801561016b575f80fd5b50610186720f3df6d732807ef1319fb7b8bb8522d0beac0281565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161011c565b3480156101b6575f80fd5b506101126101c53660046109ad565b610371565b6100fb610381565b3480156101dd575f80fd5b506100fb6101ec3660046109cd565b610394565b3480156101fc575f80fd5b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffff7487392754610186565b34801561022f575f80fd5b506100fb61023e3660046109e4565b6103d7565b34801561024e575f80fd5b506101126103f1565b6100fb610265366004610a45565b610400565b6100fb610278366004610a45565b61043d565b348015610288575f80fd5b506100fb6102973660046109cd565b610463565b3480156102a7575f80fd5b506101126102b6366004610a45565b63389a75e1600c9081525f91909152602090205490565b5f6202a30067ffffffffffffffff164201905063389a75e1600c52335f52806020600c2055337fdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d5f80a250565b610330610326876104a0565b86868686866104db565b505050505050565b63389a75e1600c52335f525f6020600c2055337ffa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c925f80a2565b5f61037b826104a0565b92915050565b6103896105d6565b6103925f61060b565b565b61039c6105d6565b5f8190556040518181527fe2e34d74290e04aa3646f259594633252b8ecbbc3bb59351491198ab661eab21906020015b60405180910390a150565b6103eb6103e3856104a0565b848484610679565b50505050565b5f6103fb426104a0565b905090565b6104086105d6565b63389a75e1600c52805f526020600c20805442111561042e57636f5e88185f526004601cfd5b5f905561043a8161060b565b50565b6104456105d6565b8060601b61045a57637448fbae5f526004601cfd5b61043a8161060b565b61046b6105d6565b60018190556040518181527febb01e1290204d373c7af00fec97d08cfacb9d72df4842842680e74f9ff76264906020016103cc565b5f815f5260205f60205f720f3df6d732807ef1319fb7b8bb8522d0beac025afa806104d257633033b0ff5f526004601cfd5b50505f51919050565b6501000000000067ffffffffffffffff821610610524576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f61056384848080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525061079f92505050565b90505f610571836008610aa5565b67ffffffffffffffff165f546105879190610acf565b905061059687878a85856107f2565b6105cc576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927543314610392576382b429005f526004601cfd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927805473ffffffffffffffffffffffffffffffffffffffff9092169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3811560ff1b8217905550565b5f7bffffffff000000000000000000000000ffffffff000000000000000066ff000000ff0000600884811c91821667ff000000ff0000009186901b91821617601090811c64ff000000ff90931665ff000000ff0090921691909117901b17602081811c9283167fffffffff000000000000000000000000ffffffff0000000000000000000000009290911b91821617604090811c73ffffffff000000000000000000000000ffffffff90931677ffffffff000000000000000000000000ffffffff0000000090921691909117901b17608081811c91901b179050610762848487846001546107f2565b610798576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050565b80515f906030146107dc576040517f5f4167e900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60205f60406020850160025afa806104d2575f80fd5b5f841561084e578460051b8601865b6001841660051b8460011c94508461082057635849603f5f526004601cfd5b85815281356020918218525f60408160025afa8061083c575f80fd5b505f5194506020018181106108015750505b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82011561088357631b6661c35f526004601cfd5b50501492915050565b803567ffffffffffffffff811681146108a3575f80fd5b919050565b5f8083601f8401126108b8575f80fd5b50813567ffffffffffffffff8111156108cf575f80fd5b6020830191508360208260051b85010111156108e9575f80fd5b9250929050565b5f805f805f8060808789031215610905575f80fd5b61090e8761088c565b9550602087013567ffffffffffffffff811115610929575f80fd5b61093589828a016108a8565b909650945050604087013567ffffffffffffffff811115610954575f80fd5b8701601f81018913610964575f80fd5b803567ffffffffffffffff81111561097a575f80fd5b89602082840101111561098b575f80fd5b602091909101935091506109a16060880161088c565b90509295509295509295565b5f602082840312156109bd575f80fd5b6109c68261088c565b9392505050565b5f602082840312156109dd575f80fd5b5035919050565b5f805f80606085870312156109f7575f80fd5b610a008561088c565b9350602085013567ffffffffffffffff811115610a1b575f80fd5b610a27878288016108a8565b9094509250610a3a90506040860161088c565b905092959194509250565b5f60208284031215610a55575f80fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146109c6575f80fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b67ffffffffffffffff8181168382160290811690818114610ac857610ac8610a78565b5092915050565b8082018082111561037b5761037b610a7856fea264697066735822122039550e460ebcb1acfc71b3cb3b60a1ce8cdc0ed5d2f44202f2b24c8df0b48eca64736f6c634300081a0033",
}

// BeaconVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use BeaconVerifierMetaData.ABI instead.
var BeaconVerifierABI = BeaconVerifierMetaData.ABI

// BeaconVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BeaconVerifierMetaData.Bin instead.
var BeaconVerifierBin = BeaconVerifierMetaData.Bin

// DeployBeaconVerifier deploys a new Ethereum contract, binding an instance of BeaconVerifier to it.
func DeployBeaconVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, _zeroValidatorPubkeyGIndex *big.Int, _executionNumberGIndex *big.Int) (common.Address, *types.Transaction, *BeaconVerifier, error) {
	parsed, err := BeaconVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BeaconVerifierBin), backend, _zeroValidatorPubkeyGIndex, _executionNumberGIndex)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BeaconVerifier{BeaconVerifierCaller: BeaconVerifierCaller{contract: contract}, BeaconVerifierTransactor: BeaconVerifierTransactor{contract: contract}, BeaconVerifierFilterer: BeaconVerifierFilterer{contract: contract}}, nil
}

// BeaconVerifier is an auto generated Go binding around an Ethereum contract.
type BeaconVerifier struct {
	BeaconVerifierCaller     // Read-only binding to the contract
	BeaconVerifierTransactor // Write-only binding to the contract
	BeaconVerifierFilterer   // Log filterer for contract events
}

// BeaconVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type BeaconVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BeaconVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BeaconVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BeaconVerifierSession struct {
	Contract     *BeaconVerifier   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BeaconVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BeaconVerifierCallerSession struct {
	Contract *BeaconVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// BeaconVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BeaconVerifierTransactorSession struct {
	Contract     *BeaconVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// BeaconVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type BeaconVerifierRaw struct {
	Contract *BeaconVerifier // Generic contract binding to access the raw methods on
}

// BeaconVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BeaconVerifierCallerRaw struct {
	Contract *BeaconVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// BeaconVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BeaconVerifierTransactorRaw struct {
	Contract *BeaconVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBeaconVerifier creates a new instance of BeaconVerifier, bound to a specific deployed contract.
func NewBeaconVerifier(address common.Address, backend bind.ContractBackend) (*BeaconVerifier, error) {
	contract, err := bindBeaconVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BeaconVerifier{BeaconVerifierCaller: BeaconVerifierCaller{contract: contract}, BeaconVerifierTransactor: BeaconVerifierTransactor{contract: contract}, BeaconVerifierFilterer: BeaconVerifierFilterer{contract: contract}}, nil
}

// NewBeaconVerifierCaller creates a new read-only instance of BeaconVerifier, bound to a specific deployed contract.
func NewBeaconVerifierCaller(address common.Address, caller bind.ContractCaller) (*BeaconVerifierCaller, error) {
	contract, err := bindBeaconVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconVerifierCaller{contract: contract}, nil
}

// NewBeaconVerifierTransactor creates a new write-only instance of BeaconVerifier, bound to a specific deployed contract.
func NewBeaconVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*BeaconVerifierTransactor, error) {
	contract, err := bindBeaconVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconVerifierTransactor{contract: contract}, nil
}

// NewBeaconVerifierFilterer creates a new log filterer instance of BeaconVerifier, bound to a specific deployed contract.
func NewBeaconVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*BeaconVerifierFilterer, error) {
	contract, err := bindBeaconVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BeaconVerifierFilterer{contract: contract}, nil
}

// bindBeaconVerifier binds a generic wrapper to an already deployed contract.
func bindBeaconVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BeaconVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconVerifier *BeaconVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconVerifier.Contract.BeaconVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconVerifier *BeaconVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.BeaconVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconVerifier *BeaconVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.BeaconVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconVerifier *BeaconVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconVerifier *BeaconVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconVerifier *BeaconVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.contract.Transact(opts, method, params...)
}

// BEACONROOTS is a free data retrieval call binding the contract method 0x56d7e8fd.
//
// Solidity: function BEACON_ROOTS() view returns(address)
func (_BeaconVerifier *BeaconVerifierCaller) BEACONROOTS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "BEACON_ROOTS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BEACONROOTS is a free data retrieval call binding the contract method 0x56d7e8fd.
//
// Solidity: function BEACON_ROOTS() view returns(address)
func (_BeaconVerifier *BeaconVerifierSession) BEACONROOTS() (common.Address, error) {
	return _BeaconVerifier.Contract.BEACONROOTS(&_BeaconVerifier.CallOpts)
}

// BEACONROOTS is a free data retrieval call binding the contract method 0x56d7e8fd.
//
// Solidity: function BEACON_ROOTS() view returns(address)
func (_BeaconVerifier *BeaconVerifierCallerSession) BEACONROOTS() (common.Address, error) {
	return _BeaconVerifier.Contract.BEACONROOTS(&_BeaconVerifier.CallOpts)
}

// ExecutionNumberGIndex is a free data retrieval call binding the contract method 0x305ea416.
//
// Solidity: function executionNumberGIndex() view returns(uint256)
func (_BeaconVerifier *BeaconVerifierCaller) ExecutionNumberGIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "executionNumberGIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExecutionNumberGIndex is a free data retrieval call binding the contract method 0x305ea416.
//
// Solidity: function executionNumberGIndex() view returns(uint256)
func (_BeaconVerifier *BeaconVerifierSession) ExecutionNumberGIndex() (*big.Int, error) {
	return _BeaconVerifier.Contract.ExecutionNumberGIndex(&_BeaconVerifier.CallOpts)
}

// ExecutionNumberGIndex is a free data retrieval call binding the contract method 0x305ea416.
//
// Solidity: function executionNumberGIndex() view returns(uint256)
func (_BeaconVerifier *BeaconVerifierCallerSession) ExecutionNumberGIndex() (*big.Int, error) {
	return _BeaconVerifier.Contract.ExecutionNumberGIndex(&_BeaconVerifier.CallOpts)
}

// GetParentBeaconBlockRoot is a free data retrieval call binding the contract method 0xefcff00a.
//
// Solidity: function getParentBeaconBlockRoot() view returns(bytes32)
func (_BeaconVerifier *BeaconVerifierCaller) GetParentBeaconBlockRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "getParentBeaconBlockRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetParentBeaconBlockRoot is a free data retrieval call binding the contract method 0xefcff00a.
//
// Solidity: function getParentBeaconBlockRoot() view returns(bytes32)
func (_BeaconVerifier *BeaconVerifierSession) GetParentBeaconBlockRoot() ([32]byte, error) {
	return _BeaconVerifier.Contract.GetParentBeaconBlockRoot(&_BeaconVerifier.CallOpts)
}

// GetParentBeaconBlockRoot is a free data retrieval call binding the contract method 0xefcff00a.
//
// Solidity: function getParentBeaconBlockRoot() view returns(bytes32)
func (_BeaconVerifier *BeaconVerifierCallerSession) GetParentBeaconBlockRoot() ([32]byte, error) {
	return _BeaconVerifier.Contract.GetParentBeaconBlockRoot(&_BeaconVerifier.CallOpts)
}

// GetParentBeaconBlockRootAt is a free data retrieval call binding the contract method 0x5e67d452.
//
// Solidity: function getParentBeaconBlockRootAt(uint64 timestamp) view returns(bytes32)
func (_BeaconVerifier *BeaconVerifierCaller) GetParentBeaconBlockRootAt(opts *bind.CallOpts, timestamp uint64) ([32]byte, error) {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "getParentBeaconBlockRootAt", timestamp)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetParentBeaconBlockRootAt is a free data retrieval call binding the contract method 0x5e67d452.
//
// Solidity: function getParentBeaconBlockRootAt(uint64 timestamp) view returns(bytes32)
func (_BeaconVerifier *BeaconVerifierSession) GetParentBeaconBlockRootAt(timestamp uint64) ([32]byte, error) {
	return _BeaconVerifier.Contract.GetParentBeaconBlockRootAt(&_BeaconVerifier.CallOpts, timestamp)
}

// GetParentBeaconBlockRootAt is a free data retrieval call binding the contract method 0x5e67d452.
//
// Solidity: function getParentBeaconBlockRootAt(uint64 timestamp) view returns(bytes32)
func (_BeaconVerifier *BeaconVerifierCallerSession) GetParentBeaconBlockRootAt(timestamp uint64) ([32]byte, error) {
	return _BeaconVerifier.Contract.GetParentBeaconBlockRootAt(&_BeaconVerifier.CallOpts, timestamp)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_BeaconVerifier *BeaconVerifierCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_BeaconVerifier *BeaconVerifierSession) Owner() (common.Address, error) {
	return _BeaconVerifier.Contract.Owner(&_BeaconVerifier.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_BeaconVerifier *BeaconVerifierCallerSession) Owner() (common.Address, error) {
	return _BeaconVerifier.Contract.Owner(&_BeaconVerifier.CallOpts)
}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_BeaconVerifier *BeaconVerifierCaller) OwnershipHandoverExpiresAt(opts *bind.CallOpts, pendingOwner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "ownershipHandoverExpiresAt", pendingOwner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_BeaconVerifier *BeaconVerifierSession) OwnershipHandoverExpiresAt(pendingOwner common.Address) (*big.Int, error) {
	return _BeaconVerifier.Contract.OwnershipHandoverExpiresAt(&_BeaconVerifier.CallOpts, pendingOwner)
}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_BeaconVerifier *BeaconVerifierCallerSession) OwnershipHandoverExpiresAt(pendingOwner common.Address) (*big.Int, error) {
	return _BeaconVerifier.Contract.OwnershipHandoverExpiresAt(&_BeaconVerifier.CallOpts, pendingOwner)
}

// VerifyBeaconBlockProposer is a free data retrieval call binding the contract method 0x491920ab.
//
// Solidity: function verifyBeaconBlockProposer(uint64 timestamp, bytes32[] validatorPubkeyProof, bytes validatorPubkey, uint64 proposerIndex) view returns()
func (_BeaconVerifier *BeaconVerifierCaller) VerifyBeaconBlockProposer(opts *bind.CallOpts, timestamp uint64, validatorPubkeyProof [][32]byte, validatorPubkey []byte, proposerIndex uint64) error {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "verifyBeaconBlockProposer", timestamp, validatorPubkeyProof, validatorPubkey, proposerIndex)

	if err != nil {
		return err
	}

	return err

}

// VerifyBeaconBlockProposer is a free data retrieval call binding the contract method 0x491920ab.
//
// Solidity: function verifyBeaconBlockProposer(uint64 timestamp, bytes32[] validatorPubkeyProof, bytes validatorPubkey, uint64 proposerIndex) view returns()
func (_BeaconVerifier *BeaconVerifierSession) VerifyBeaconBlockProposer(timestamp uint64, validatorPubkeyProof [][32]byte, validatorPubkey []byte, proposerIndex uint64) error {
	return _BeaconVerifier.Contract.VerifyBeaconBlockProposer(&_BeaconVerifier.CallOpts, timestamp, validatorPubkeyProof, validatorPubkey, proposerIndex)
}

// VerifyBeaconBlockProposer is a free data retrieval call binding the contract method 0x491920ab.
//
// Solidity: function verifyBeaconBlockProposer(uint64 timestamp, bytes32[] validatorPubkeyProof, bytes validatorPubkey, uint64 proposerIndex) view returns()
func (_BeaconVerifier *BeaconVerifierCallerSession) VerifyBeaconBlockProposer(timestamp uint64, validatorPubkeyProof [][32]byte, validatorPubkey []byte, proposerIndex uint64) error {
	return _BeaconVerifier.Contract.VerifyBeaconBlockProposer(&_BeaconVerifier.CallOpts, timestamp, validatorPubkeyProof, validatorPubkey, proposerIndex)
}

// VerifyExecutionNumber is a free data retrieval call binding the contract method 0xae9af3e3.
//
// Solidity: function verifyExecutionNumber(uint64 timestamp, bytes32[] executionNumberProof, uint64 blockNumber) view returns()
func (_BeaconVerifier *BeaconVerifierCaller) VerifyExecutionNumber(opts *bind.CallOpts, timestamp uint64, executionNumberProof [][32]byte, blockNumber uint64) error {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "verifyExecutionNumber", timestamp, executionNumberProof, blockNumber)

	if err != nil {
		return err
	}

	return err

}

// VerifyExecutionNumber is a free data retrieval call binding the contract method 0xae9af3e3.
//
// Solidity: function verifyExecutionNumber(uint64 timestamp, bytes32[] executionNumberProof, uint64 blockNumber) view returns()
func (_BeaconVerifier *BeaconVerifierSession) VerifyExecutionNumber(timestamp uint64, executionNumberProof [][32]byte, blockNumber uint64) error {
	return _BeaconVerifier.Contract.VerifyExecutionNumber(&_BeaconVerifier.CallOpts, timestamp, executionNumberProof, blockNumber)
}

// VerifyExecutionNumber is a free data retrieval call binding the contract method 0xae9af3e3.
//
// Solidity: function verifyExecutionNumber(uint64 timestamp, bytes32[] executionNumberProof, uint64 blockNumber) view returns()
func (_BeaconVerifier *BeaconVerifierCallerSession) VerifyExecutionNumber(timestamp uint64, executionNumberProof [][32]byte, blockNumber uint64) error {
	return _BeaconVerifier.Contract.VerifyExecutionNumber(&_BeaconVerifier.CallOpts, timestamp, executionNumberProof, blockNumber)
}

// ZeroValidatorPubkeyGIndex is a free data retrieval call binding the contract method 0x35222ff1.
//
// Solidity: function zeroValidatorPubkeyGIndex() view returns(uint256)
func (_BeaconVerifier *BeaconVerifierCaller) ZeroValidatorPubkeyGIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "zeroValidatorPubkeyGIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ZeroValidatorPubkeyGIndex is a free data retrieval call binding the contract method 0x35222ff1.
//
// Solidity: function zeroValidatorPubkeyGIndex() view returns(uint256)
func (_BeaconVerifier *BeaconVerifierSession) ZeroValidatorPubkeyGIndex() (*big.Int, error) {
	return _BeaconVerifier.Contract.ZeroValidatorPubkeyGIndex(&_BeaconVerifier.CallOpts)
}

// ZeroValidatorPubkeyGIndex is a free data retrieval call binding the contract method 0x35222ff1.
//
// Solidity: function zeroValidatorPubkeyGIndex() view returns(uint256)
func (_BeaconVerifier *BeaconVerifierCallerSession) ZeroValidatorPubkeyGIndex() (*big.Int, error) {
	return _BeaconVerifier.Contract.ZeroValidatorPubkeyGIndex(&_BeaconVerifier.CallOpts)
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_BeaconVerifier *BeaconVerifierTransactor) CancelOwnershipHandover(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconVerifier.contract.Transact(opts, "cancelOwnershipHandover")
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_BeaconVerifier *BeaconVerifierSession) CancelOwnershipHandover() (*types.Transaction, error) {
	return _BeaconVerifier.Contract.CancelOwnershipHandover(&_BeaconVerifier.TransactOpts)
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_BeaconVerifier *BeaconVerifierTransactorSession) CancelOwnershipHandover() (*types.Transaction, error) {
	return _BeaconVerifier.Contract.CancelOwnershipHandover(&_BeaconVerifier.TransactOpts)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_BeaconVerifier *BeaconVerifierTransactor) CompleteOwnershipHandover(opts *bind.TransactOpts, pendingOwner common.Address) (*types.Transaction, error) {
	return _BeaconVerifier.contract.Transact(opts, "completeOwnershipHandover", pendingOwner)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_BeaconVerifier *BeaconVerifierSession) CompleteOwnershipHandover(pendingOwner common.Address) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.CompleteOwnershipHandover(&_BeaconVerifier.TransactOpts, pendingOwner)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_BeaconVerifier *BeaconVerifierTransactorSession) CompleteOwnershipHandover(pendingOwner common.Address) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.CompleteOwnershipHandover(&_BeaconVerifier.TransactOpts, pendingOwner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_BeaconVerifier *BeaconVerifierTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconVerifier.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_BeaconVerifier *BeaconVerifierSession) RenounceOwnership() (*types.Transaction, error) {
	return _BeaconVerifier.Contract.RenounceOwnership(&_BeaconVerifier.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_BeaconVerifier *BeaconVerifierTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BeaconVerifier.Contract.RenounceOwnership(&_BeaconVerifier.TransactOpts)
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_BeaconVerifier *BeaconVerifierTransactor) RequestOwnershipHandover(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconVerifier.contract.Transact(opts, "requestOwnershipHandover")
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_BeaconVerifier *BeaconVerifierSession) RequestOwnershipHandover() (*types.Transaction, error) {
	return _BeaconVerifier.Contract.RequestOwnershipHandover(&_BeaconVerifier.TransactOpts)
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_BeaconVerifier *BeaconVerifierTransactorSession) RequestOwnershipHandover() (*types.Transaction, error) {
	return _BeaconVerifier.Contract.RequestOwnershipHandover(&_BeaconVerifier.TransactOpts)
}

// SetExecutionNumberGIndex is a paid mutator transaction binding the contract method 0xf769afd1.
//
// Solidity: function setExecutionNumberGIndex(uint256 _executionNumberGIndex) returns()
func (_BeaconVerifier *BeaconVerifierTransactor) SetExecutionNumberGIndex(opts *bind.TransactOpts, _executionNumberGIndex *big.Int) (*types.Transaction, error) {
	return _BeaconVerifier.contract.Transact(opts, "setExecutionNumberGIndex", _executionNumberGIndex)
}

// SetExecutionNumberGIndex is a paid mutator transaction binding the contract method 0xf769afd1.
//
// Solidity: function setExecutionNumberGIndex(uint256 _executionNumberGIndex) returns()
func (_BeaconVerifier *BeaconVerifierSession) SetExecutionNumberGIndex(_executionNumberGIndex *big.Int) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.SetExecutionNumberGIndex(&_BeaconVerifier.TransactOpts, _executionNumberGIndex)
}

// SetExecutionNumberGIndex is a paid mutator transaction binding the contract method 0xf769afd1.
//
// Solidity: function setExecutionNumberGIndex(uint256 _executionNumberGIndex) returns()
func (_BeaconVerifier *BeaconVerifierTransactorSession) SetExecutionNumberGIndex(_executionNumberGIndex *big.Int) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.SetExecutionNumberGIndex(&_BeaconVerifier.TransactOpts, _executionNumberGIndex)
}

// SetZeroValidatorPubkeyGIndex is a paid mutator transaction binding the contract method 0x81705771.
//
// Solidity: function setZeroValidatorPubkeyGIndex(uint256 _zeroValidatorPubkeyGIndex) returns()
func (_BeaconVerifier *BeaconVerifierTransactor) SetZeroValidatorPubkeyGIndex(opts *bind.TransactOpts, _zeroValidatorPubkeyGIndex *big.Int) (*types.Transaction, error) {
	return _BeaconVerifier.contract.Transact(opts, "setZeroValidatorPubkeyGIndex", _zeroValidatorPubkeyGIndex)
}

// SetZeroValidatorPubkeyGIndex is a paid mutator transaction binding the contract method 0x81705771.
//
// Solidity: function setZeroValidatorPubkeyGIndex(uint256 _zeroValidatorPubkeyGIndex) returns()
func (_BeaconVerifier *BeaconVerifierSession) SetZeroValidatorPubkeyGIndex(_zeroValidatorPubkeyGIndex *big.Int) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.SetZeroValidatorPubkeyGIndex(&_BeaconVerifier.TransactOpts, _zeroValidatorPubkeyGIndex)
}

// SetZeroValidatorPubkeyGIndex is a paid mutator transaction binding the contract method 0x81705771.
//
// Solidity: function setZeroValidatorPubkeyGIndex(uint256 _zeroValidatorPubkeyGIndex) returns()
func (_BeaconVerifier *BeaconVerifierTransactorSession) SetZeroValidatorPubkeyGIndex(_zeroValidatorPubkeyGIndex *big.Int) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.SetZeroValidatorPubkeyGIndex(&_BeaconVerifier.TransactOpts, _zeroValidatorPubkeyGIndex)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_BeaconVerifier *BeaconVerifierTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BeaconVerifier.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_BeaconVerifier *BeaconVerifierSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.TransferOwnership(&_BeaconVerifier.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_BeaconVerifier *BeaconVerifierTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.TransferOwnership(&_BeaconVerifier.TransactOpts, newOwner)
}

// BeaconVerifierExecutionNumberGIndexChangedIterator is returned from FilterExecutionNumberGIndexChanged and is used to iterate over the raw logs and unpacked data for ExecutionNumberGIndexChanged events raised by the BeaconVerifier contract.
type BeaconVerifierExecutionNumberGIndexChangedIterator struct {
	Event *BeaconVerifierExecutionNumberGIndexChanged // Event containing the contract specifics and raw log

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
func (it *BeaconVerifierExecutionNumberGIndexChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVerifierExecutionNumberGIndexChanged)
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
		it.Event = new(BeaconVerifierExecutionNumberGIndexChanged)
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
func (it *BeaconVerifierExecutionNumberGIndexChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconVerifierExecutionNumberGIndexChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconVerifierExecutionNumberGIndexChanged represents a ExecutionNumberGIndexChanged event raised by the BeaconVerifier contract.
type BeaconVerifierExecutionNumberGIndexChanged struct {
	NewExecutionNumberGIndex *big.Int
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterExecutionNumberGIndexChanged is a free log retrieval operation binding the contract event 0xebb01e1290204d373c7af00fec97d08cfacb9d72df4842842680e74f9ff76264.
//
// Solidity: event ExecutionNumberGIndexChanged(uint256 newExecutionNumberGIndex)
func (_BeaconVerifier *BeaconVerifierFilterer) FilterExecutionNumberGIndexChanged(opts *bind.FilterOpts) (*BeaconVerifierExecutionNumberGIndexChangedIterator, error) {

	logs, sub, err := _BeaconVerifier.contract.FilterLogs(opts, "ExecutionNumberGIndexChanged")
	if err != nil {
		return nil, err
	}
	return &BeaconVerifierExecutionNumberGIndexChangedIterator{contract: _BeaconVerifier.contract, event: "ExecutionNumberGIndexChanged", logs: logs, sub: sub}, nil
}

// WatchExecutionNumberGIndexChanged is a free log subscription operation binding the contract event 0xebb01e1290204d373c7af00fec97d08cfacb9d72df4842842680e74f9ff76264.
//
// Solidity: event ExecutionNumberGIndexChanged(uint256 newExecutionNumberGIndex)
func (_BeaconVerifier *BeaconVerifierFilterer) WatchExecutionNumberGIndexChanged(opts *bind.WatchOpts, sink chan<- *BeaconVerifierExecutionNumberGIndexChanged) (event.Subscription, error) {

	logs, sub, err := _BeaconVerifier.contract.WatchLogs(opts, "ExecutionNumberGIndexChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconVerifierExecutionNumberGIndexChanged)
				if err := _BeaconVerifier.contract.UnpackLog(event, "ExecutionNumberGIndexChanged", log); err != nil {
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

// ParseExecutionNumberGIndexChanged is a log parse operation binding the contract event 0xebb01e1290204d373c7af00fec97d08cfacb9d72df4842842680e74f9ff76264.
//
// Solidity: event ExecutionNumberGIndexChanged(uint256 newExecutionNumberGIndex)
func (_BeaconVerifier *BeaconVerifierFilterer) ParseExecutionNumberGIndexChanged(log types.Log) (*BeaconVerifierExecutionNumberGIndexChanged, error) {
	event := new(BeaconVerifierExecutionNumberGIndexChanged)
	if err := _BeaconVerifier.contract.UnpackLog(event, "ExecutionNumberGIndexChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconVerifierOwnershipHandoverCanceledIterator is returned from FilterOwnershipHandoverCanceled and is used to iterate over the raw logs and unpacked data for OwnershipHandoverCanceled events raised by the BeaconVerifier contract.
type BeaconVerifierOwnershipHandoverCanceledIterator struct {
	Event *BeaconVerifierOwnershipHandoverCanceled // Event containing the contract specifics and raw log

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
func (it *BeaconVerifierOwnershipHandoverCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVerifierOwnershipHandoverCanceled)
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
		it.Event = new(BeaconVerifierOwnershipHandoverCanceled)
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
func (it *BeaconVerifierOwnershipHandoverCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconVerifierOwnershipHandoverCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconVerifierOwnershipHandoverCanceled represents a OwnershipHandoverCanceled event raised by the BeaconVerifier contract.
type BeaconVerifierOwnershipHandoverCanceled struct {
	PendingOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOwnershipHandoverCanceled is a free log retrieval operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_BeaconVerifier *BeaconVerifierFilterer) FilterOwnershipHandoverCanceled(opts *bind.FilterOpts, pendingOwner []common.Address) (*BeaconVerifierOwnershipHandoverCanceledIterator, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _BeaconVerifier.contract.FilterLogs(opts, "OwnershipHandoverCanceled", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BeaconVerifierOwnershipHandoverCanceledIterator{contract: _BeaconVerifier.contract, event: "OwnershipHandoverCanceled", logs: logs, sub: sub}, nil
}

// WatchOwnershipHandoverCanceled is a free log subscription operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_BeaconVerifier *BeaconVerifierFilterer) WatchOwnershipHandoverCanceled(opts *bind.WatchOpts, sink chan<- *BeaconVerifierOwnershipHandoverCanceled, pendingOwner []common.Address) (event.Subscription, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _BeaconVerifier.contract.WatchLogs(opts, "OwnershipHandoverCanceled", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconVerifierOwnershipHandoverCanceled)
				if err := _BeaconVerifier.contract.UnpackLog(event, "OwnershipHandoverCanceled", log); err != nil {
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
func (_BeaconVerifier *BeaconVerifierFilterer) ParseOwnershipHandoverCanceled(log types.Log) (*BeaconVerifierOwnershipHandoverCanceled, error) {
	event := new(BeaconVerifierOwnershipHandoverCanceled)
	if err := _BeaconVerifier.contract.UnpackLog(event, "OwnershipHandoverCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconVerifierOwnershipHandoverRequestedIterator is returned from FilterOwnershipHandoverRequested and is used to iterate over the raw logs and unpacked data for OwnershipHandoverRequested events raised by the BeaconVerifier contract.
type BeaconVerifierOwnershipHandoverRequestedIterator struct {
	Event *BeaconVerifierOwnershipHandoverRequested // Event containing the contract specifics and raw log

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
func (it *BeaconVerifierOwnershipHandoverRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVerifierOwnershipHandoverRequested)
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
		it.Event = new(BeaconVerifierOwnershipHandoverRequested)
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
func (it *BeaconVerifierOwnershipHandoverRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconVerifierOwnershipHandoverRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconVerifierOwnershipHandoverRequested represents a OwnershipHandoverRequested event raised by the BeaconVerifier contract.
type BeaconVerifierOwnershipHandoverRequested struct {
	PendingOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOwnershipHandoverRequested is a free log retrieval operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_BeaconVerifier *BeaconVerifierFilterer) FilterOwnershipHandoverRequested(opts *bind.FilterOpts, pendingOwner []common.Address) (*BeaconVerifierOwnershipHandoverRequestedIterator, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _BeaconVerifier.contract.FilterLogs(opts, "OwnershipHandoverRequested", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BeaconVerifierOwnershipHandoverRequestedIterator{contract: _BeaconVerifier.contract, event: "OwnershipHandoverRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipHandoverRequested is a free log subscription operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_BeaconVerifier *BeaconVerifierFilterer) WatchOwnershipHandoverRequested(opts *bind.WatchOpts, sink chan<- *BeaconVerifierOwnershipHandoverRequested, pendingOwner []common.Address) (event.Subscription, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _BeaconVerifier.contract.WatchLogs(opts, "OwnershipHandoverRequested", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconVerifierOwnershipHandoverRequested)
				if err := _BeaconVerifier.contract.UnpackLog(event, "OwnershipHandoverRequested", log); err != nil {
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
func (_BeaconVerifier *BeaconVerifierFilterer) ParseOwnershipHandoverRequested(log types.Log) (*BeaconVerifierOwnershipHandoverRequested, error) {
	event := new(BeaconVerifierOwnershipHandoverRequested)
	if err := _BeaconVerifier.contract.UnpackLog(event, "OwnershipHandoverRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconVerifierOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BeaconVerifier contract.
type BeaconVerifierOwnershipTransferredIterator struct {
	Event *BeaconVerifierOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BeaconVerifierOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVerifierOwnershipTransferred)
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
		it.Event = new(BeaconVerifierOwnershipTransferred)
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
func (it *BeaconVerifierOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconVerifierOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconVerifierOwnershipTransferred represents a OwnershipTransferred event raised by the BeaconVerifier contract.
type BeaconVerifierOwnershipTransferred struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_BeaconVerifier *BeaconVerifierFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, oldOwner []common.Address, newOwner []common.Address) (*BeaconVerifierOwnershipTransferredIterator, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BeaconVerifier.contract.FilterLogs(opts, "OwnershipTransferred", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BeaconVerifierOwnershipTransferredIterator{contract: _BeaconVerifier.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_BeaconVerifier *BeaconVerifierFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BeaconVerifierOwnershipTransferred, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BeaconVerifier.contract.WatchLogs(opts, "OwnershipTransferred", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconVerifierOwnershipTransferred)
				if err := _BeaconVerifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_BeaconVerifier *BeaconVerifierFilterer) ParseOwnershipTransferred(log types.Log) (*BeaconVerifierOwnershipTransferred, error) {
	event := new(BeaconVerifierOwnershipTransferred)
	if err := _BeaconVerifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconVerifierZeroValidatorPubkeyGIndexChangedIterator is returned from FilterZeroValidatorPubkeyGIndexChanged and is used to iterate over the raw logs and unpacked data for ZeroValidatorPubkeyGIndexChanged events raised by the BeaconVerifier contract.
type BeaconVerifierZeroValidatorPubkeyGIndexChangedIterator struct {
	Event *BeaconVerifierZeroValidatorPubkeyGIndexChanged // Event containing the contract specifics and raw log

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
func (it *BeaconVerifierZeroValidatorPubkeyGIndexChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVerifierZeroValidatorPubkeyGIndexChanged)
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
		it.Event = new(BeaconVerifierZeroValidatorPubkeyGIndexChanged)
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
func (it *BeaconVerifierZeroValidatorPubkeyGIndexChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconVerifierZeroValidatorPubkeyGIndexChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconVerifierZeroValidatorPubkeyGIndexChanged represents a ZeroValidatorPubkeyGIndexChanged event raised by the BeaconVerifier contract.
type BeaconVerifierZeroValidatorPubkeyGIndexChanged struct {
	NewZeroValidatorPubkeyGIndex *big.Int
	Raw                          types.Log // Blockchain specific contextual infos
}

// FilterZeroValidatorPubkeyGIndexChanged is a free log retrieval operation binding the contract event 0xe2e34d74290e04aa3646f259594633252b8ecbbc3bb59351491198ab661eab21.
//
// Solidity: event ZeroValidatorPubkeyGIndexChanged(uint256 newZeroValidatorPubkeyGIndex)
func (_BeaconVerifier *BeaconVerifierFilterer) FilterZeroValidatorPubkeyGIndexChanged(opts *bind.FilterOpts) (*BeaconVerifierZeroValidatorPubkeyGIndexChangedIterator, error) {

	logs, sub, err := _BeaconVerifier.contract.FilterLogs(opts, "ZeroValidatorPubkeyGIndexChanged")
	if err != nil {
		return nil, err
	}
	return &BeaconVerifierZeroValidatorPubkeyGIndexChangedIterator{contract: _BeaconVerifier.contract, event: "ZeroValidatorPubkeyGIndexChanged", logs: logs, sub: sub}, nil
}

// WatchZeroValidatorPubkeyGIndexChanged is a free log subscription operation binding the contract event 0xe2e34d74290e04aa3646f259594633252b8ecbbc3bb59351491198ab661eab21.
//
// Solidity: event ZeroValidatorPubkeyGIndexChanged(uint256 newZeroValidatorPubkeyGIndex)
func (_BeaconVerifier *BeaconVerifierFilterer) WatchZeroValidatorPubkeyGIndexChanged(opts *bind.WatchOpts, sink chan<- *BeaconVerifierZeroValidatorPubkeyGIndexChanged) (event.Subscription, error) {

	logs, sub, err := _BeaconVerifier.contract.WatchLogs(opts, "ZeroValidatorPubkeyGIndexChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconVerifierZeroValidatorPubkeyGIndexChanged)
				if err := _BeaconVerifier.contract.UnpackLog(event, "ZeroValidatorPubkeyGIndexChanged", log); err != nil {
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

// ParseZeroValidatorPubkeyGIndexChanged is a log parse operation binding the contract event 0xe2e34d74290e04aa3646f259594633252b8ecbbc3bb59351491198ab661eab21.
//
// Solidity: event ZeroValidatorPubkeyGIndexChanged(uint256 newZeroValidatorPubkeyGIndex)
func (_BeaconVerifier *BeaconVerifierFilterer) ParseZeroValidatorPubkeyGIndexChanged(log types.Log) (*BeaconVerifierZeroValidatorPubkeyGIndexChanged, error) {
	event := new(BeaconVerifierZeroValidatorPubkeyGIndexChanged)
	if err := _BeaconVerifier.contract.UnpackLog(event, "ZeroValidatorPubkeyGIndexChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
