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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_zeroValidatorPubkeyGIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_executionNumberGIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_executionFeeRecipientGIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BEACON_ROOTS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelOwnershipHandover\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"completeOwnershipHandover\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"executionFeeRecipientGIndex\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"executionNumberGIndex\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParentBeaconBlockRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getParentBeaconBlockRootAt\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"result\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownershipHandoverExpiresAt\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"requestOwnershipHandover\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setExecutionFeeRecipientGIndex\",\"inputs\":[{\"name\":\"_executionFeeRecipientGIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setExecutionNumberGIndex\",\"inputs\":[{\"name\":\"_executionNumberGIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setZeroValidatorPubkeyGIndex\",\"inputs\":[{\"name\":\"_zeroValidatorPubkeyGIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"verifyBeaconBlockProposer\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"proposerPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"proposerPubkeyProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyCoinbase\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"coinbase\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"coinbaseProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyExecutionNumber\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"executionNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"executionNumberProof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"zeroValidatorPubkeyGIndex\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ExecutionFeeRecipientGIndexChanged\",\"inputs\":[{\"name\":\"newExecutionFeeRecipientGIndex\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ExecutionNumberGIndexChanged\",\"inputs\":[{\"name\":\"newExecutionNumberGIndex\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipHandoverCanceled\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipHandoverRequested\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"oldOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ZeroValidatorPubkeyGIndexChanged\",\"inputs\":[{\"name\":\"newZeroValidatorPubkeyGIndex\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidValidatorPubkeyLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NewOwnerIsZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoHandoverRequest\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RootNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
	Bin: "0x6080604052348015600e575f80fd5b50604051610ccc380380610ccc833981016040819052602b91609e565b5f839055600182905560028190556040336047565b50505060c8565b638b78c6d819805415606057630dc149f05f526004601cfd5b6001600160a01b03909116801560ff1b8117909155805f7f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08180a350565b5f805f6060848603121560af575f80fd5b5050815160208301516040909301519094929350919050565b610bf7806100d55f395ff3fe60806040526004361061013d575f3560e01c80638da5cb5b116100bb578063efcff00a11610071578063f2fde38b11610057578063f2fde38b1461030b578063f769afd11461031e578063fee81cf41461033d575f80fd5b8063efcff00a146102e4578063f04e283e146102f8575f80fd5b8063c0c27230116100a1578063c0c2723014610291578063cbb60852146102b0578063e701fa76146102cf575f80fd5b80638da5cb5b1461023f578063ab15494a14610272575f80fd5b806356d7e8fd11610110578063715018a6116100f6578063715018a6146101f957806378e97954146102015780638170577114610220575f80fd5b806356d7e8fd1461018f5780635e67d452146101da575f80fd5b80632569296214610141578063305ea4161461014b57806335222ff11461017357806354d1f13d14610187575b5f80fd5b61014961036e565b005b348015610156575f80fd5b5061016060015481565b6040519081526020015b60405180910390f35b34801561017e575f80fd5b506101605f5481565b6101496103bb565b34801561019a575f80fd5b506101b5720f3df6d732807ef1319fb7b8bb8522d0beac0281565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161016a565b3480156101e5575f80fd5b506101606101f4366004610956565b6103f4565b610149610404565b34801561020c575f80fd5b5061014961021b3660046109e1565b610417565b34801561022b575f80fd5b5061014961023a366004610a3e565b610431565b34801561024a575f80fd5b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927546101b5565b34801561027d575f80fd5b5061014961028c366004610a55565b610474565b34801561029c575f80fd5b506101496102ab366004610a3e565b610492565b3480156102bb575f80fd5b506101496102ca366004610b14565b6104cf565b3480156102da575f80fd5b5061016060025481565b3480156102ef575f80fd5b506101606104e3565b610149610306366004610b3e565b6104f2565b610149610319366004610b3e565b61052f565b348015610329575f80fd5b50610149610338366004610a3e565b610555565b348015610348575f80fd5b50610160610357366004610b3e565b63389a75e1600c9081525f91909152602090205490565b5f6202a30067ffffffffffffffff164201905063389a75e1600c52335f52806020600c2055337fdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d5f80a250565b63389a75e1600c52335f525f6020600c2055337ffa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c925f80a2565b5f6103fe82610592565b92915050565b61040c6105cd565b6104155f610602565b565b61042b61042385610592565b838386610670565b50505050565b6104396105cd565b5f8190556040518181527fe2e34d74290e04aa3646f259594633252b8ecbbc3bb59351491198ab661eab21906020015b60405180910390a150565b61048a61048087610592565b838387878a6106e6565b505050505050565b61049a6105cd565b60028190556040518181527f15e42eec45edd1a051ce50f823a5d6482237d402c471639da37b62e85530154890602001610469565b61042b6104db85610592565b8383866107e1565b5f6104ed42610592565b905090565b6104fa6105cd565b63389a75e1600c52805f526020600c20805442111561052057636f5e88185f526004601cfd5b5f905561052c81610602565b50565b6105376105cd565b8060601b61054c57637448fbae5f526004601cfd5b61052c81610602565b61055d6105cd565b60018190556040518181527febb01e1290204d373c7af00fec97d08cfacb9d72df4842842680e74f9ff7626490602001610469565b5f815f5260205f60205f720f3df6d732807ef1319fb7b8bb8522d0beac025afa806105c457633033b0ff5f526004601cfd5b50505f51919050565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927543314610415576382b429005f526004601cfd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927805473ffffffffffffffffffffffffffffffffffffffff9092169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3811560ff1b8217905550565b5f606082901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001690506106a98484878460025461084d565b6106df576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050565b6501000000000067ffffffffffffffff82161061072f576040517f1390f2a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f61076e84848080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152506108e792505050565b90505f61077c836008610b84565b67ffffffffffffffff165f546107929190610bae565b90506107a187878a858561084d565b6107d7576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050565b5f65ff000000ff00600883811b91821664ff000000ff9185901c91821617601090811b67ff000000ff0000009390931666ff000000ff00009290921691909117901c17602081811c63ffffffff1691901b67ffffffff00000000161760c01b90506106a9848487846001545b5f84156108a9578460051b8601865b6001841660051b8460011c94508461087b57635849603f5f526004601cfd5b85815281356020918218525f60408160025afa80610897575f80fd5b505f51945060200181811061085c5750505b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8201156108de57631b6661c35f526004601cfd5b50501492915050565b80515f90603014610924576040517f5f4167e900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60205f60406020850160025afa806105c4575f80fd5b803567ffffffffffffffff81168114610951575f80fd5b919050565b5f60208284031215610966575f80fd5b61096f8261093a565b9392505050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610951575f80fd5b5f8083601f8401126109a9575f80fd5b50813567ffffffffffffffff8111156109c0575f80fd5b6020830191508360208260051b85010111156109da575f80fd5b9250929050565b5f805f80606085870312156109f4575f80fd5b6109fd8561093a565b9350610a0b60208601610976565b9250604085013567ffffffffffffffff811115610a26575f80fd5b610a3287828801610999565b95989497509550505050565b5f60208284031215610a4e575f80fd5b5035919050565b5f805f805f8060808789031215610a6a575f80fd5b610a738761093a565b9550610a816020880161093a565b9450604087013567ffffffffffffffff811115610a9c575f80fd5b8701601f81018913610aac575f80fd5b803567ffffffffffffffff811115610ac2575f80fd5b896020828401011115610ad3575f80fd5b60209190910194509250606087013567ffffffffffffffff811115610af6575f80fd5b610b0289828a01610999565b979a9699509497509295939492505050565b5f805f8060608587031215610b27575f80fd5b610b308561093a565b9350610a0b6020860161093a565b5f60208284031215610b4e575f80fd5b61096f82610976565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b67ffffffffffffffff8181168382160290811690818114610ba757610ba7610b57565b5092915050565b808201808211156103fe576103fe610b5756fea26469706673582212204277a535e42dca01d82f7a275fc22d40d4ab7ae2fdad7309d299900ac418acc664736f6c634300081a0033",
}

// BeaconVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use BeaconVerifierMetaData.ABI instead.
var BeaconVerifierABI = BeaconVerifierMetaData.ABI

// BeaconVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BeaconVerifierMetaData.Bin instead.
var BeaconVerifierBin = BeaconVerifierMetaData.Bin

// DeployBeaconVerifier deploys a new Ethereum contract, binding an instance of BeaconVerifier to it.
func DeployBeaconVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, _zeroValidatorPubkeyGIndex *big.Int, _executionNumberGIndex *big.Int, _executionFeeRecipientGIndex *big.Int) (common.Address, *types.Transaction, *BeaconVerifier, error) {
	parsed, err := BeaconVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BeaconVerifierBin), backend, _zeroValidatorPubkeyGIndex, _executionNumberGIndex, _executionFeeRecipientGIndex)
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

// ExecutionFeeRecipientGIndex is a free data retrieval call binding the contract method 0xe701fa76.
//
// Solidity: function executionFeeRecipientGIndex() view returns(uint256)
func (_BeaconVerifier *BeaconVerifierCaller) ExecutionFeeRecipientGIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "executionFeeRecipientGIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExecutionFeeRecipientGIndex is a free data retrieval call binding the contract method 0xe701fa76.
//
// Solidity: function executionFeeRecipientGIndex() view returns(uint256)
func (_BeaconVerifier *BeaconVerifierSession) ExecutionFeeRecipientGIndex() (*big.Int, error) {
	return _BeaconVerifier.Contract.ExecutionFeeRecipientGIndex(&_BeaconVerifier.CallOpts)
}

// ExecutionFeeRecipientGIndex is a free data retrieval call binding the contract method 0xe701fa76.
//
// Solidity: function executionFeeRecipientGIndex() view returns(uint256)
func (_BeaconVerifier *BeaconVerifierCallerSession) ExecutionFeeRecipientGIndex() (*big.Int, error) {
	return _BeaconVerifier.Contract.ExecutionFeeRecipientGIndex(&_BeaconVerifier.CallOpts)
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

// VerifyBeaconBlockProposer is a free data retrieval call binding the contract method 0xab15494a.
//
// Solidity: function verifyBeaconBlockProposer(uint64 timestamp, uint64 proposerIndex, bytes proposerPubkey, bytes32[] proposerPubkeyProof) view returns()
func (_BeaconVerifier *BeaconVerifierCaller) VerifyBeaconBlockProposer(opts *bind.CallOpts, timestamp uint64, proposerIndex uint64, proposerPubkey []byte, proposerPubkeyProof [][32]byte) error {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "verifyBeaconBlockProposer", timestamp, proposerIndex, proposerPubkey, proposerPubkeyProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyBeaconBlockProposer is a free data retrieval call binding the contract method 0xab15494a.
//
// Solidity: function verifyBeaconBlockProposer(uint64 timestamp, uint64 proposerIndex, bytes proposerPubkey, bytes32[] proposerPubkeyProof) view returns()
func (_BeaconVerifier *BeaconVerifierSession) VerifyBeaconBlockProposer(timestamp uint64, proposerIndex uint64, proposerPubkey []byte, proposerPubkeyProof [][32]byte) error {
	return _BeaconVerifier.Contract.VerifyBeaconBlockProposer(&_BeaconVerifier.CallOpts, timestamp, proposerIndex, proposerPubkey, proposerPubkeyProof)
}

// VerifyBeaconBlockProposer is a free data retrieval call binding the contract method 0xab15494a.
//
// Solidity: function verifyBeaconBlockProposer(uint64 timestamp, uint64 proposerIndex, bytes proposerPubkey, bytes32[] proposerPubkeyProof) view returns()
func (_BeaconVerifier *BeaconVerifierCallerSession) VerifyBeaconBlockProposer(timestamp uint64, proposerIndex uint64, proposerPubkey []byte, proposerPubkeyProof [][32]byte) error {
	return _BeaconVerifier.Contract.VerifyBeaconBlockProposer(&_BeaconVerifier.CallOpts, timestamp, proposerIndex, proposerPubkey, proposerPubkeyProof)
}

// VerifyCoinbase is a free data retrieval call binding the contract method 0x78e97954.
//
// Solidity: function verifyCoinbase(uint64 timestamp, address coinbase, bytes32[] coinbaseProof) view returns()
func (_BeaconVerifier *BeaconVerifierCaller) VerifyCoinbase(opts *bind.CallOpts, timestamp uint64, coinbase common.Address, coinbaseProof [][32]byte) error {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "verifyCoinbase", timestamp, coinbase, coinbaseProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyCoinbase is a free data retrieval call binding the contract method 0x78e97954.
//
// Solidity: function verifyCoinbase(uint64 timestamp, address coinbase, bytes32[] coinbaseProof) view returns()
func (_BeaconVerifier *BeaconVerifierSession) VerifyCoinbase(timestamp uint64, coinbase common.Address, coinbaseProof [][32]byte) error {
	return _BeaconVerifier.Contract.VerifyCoinbase(&_BeaconVerifier.CallOpts, timestamp, coinbase, coinbaseProof)
}

// VerifyCoinbase is a free data retrieval call binding the contract method 0x78e97954.
//
// Solidity: function verifyCoinbase(uint64 timestamp, address coinbase, bytes32[] coinbaseProof) view returns()
func (_BeaconVerifier *BeaconVerifierCallerSession) VerifyCoinbase(timestamp uint64, coinbase common.Address, coinbaseProof [][32]byte) error {
	return _BeaconVerifier.Contract.VerifyCoinbase(&_BeaconVerifier.CallOpts, timestamp, coinbase, coinbaseProof)
}

// VerifyExecutionNumber is a free data retrieval call binding the contract method 0xcbb60852.
//
// Solidity: function verifyExecutionNumber(uint64 timestamp, uint64 executionNumber, bytes32[] executionNumberProof) view returns()
func (_BeaconVerifier *BeaconVerifierCaller) VerifyExecutionNumber(opts *bind.CallOpts, timestamp uint64, executionNumber uint64, executionNumberProof [][32]byte) error {
	var out []interface{}
	err := _BeaconVerifier.contract.Call(opts, &out, "verifyExecutionNumber", timestamp, executionNumber, executionNumberProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyExecutionNumber is a free data retrieval call binding the contract method 0xcbb60852.
//
// Solidity: function verifyExecutionNumber(uint64 timestamp, uint64 executionNumber, bytes32[] executionNumberProof) view returns()
func (_BeaconVerifier *BeaconVerifierSession) VerifyExecutionNumber(timestamp uint64, executionNumber uint64, executionNumberProof [][32]byte) error {
	return _BeaconVerifier.Contract.VerifyExecutionNumber(&_BeaconVerifier.CallOpts, timestamp, executionNumber, executionNumberProof)
}

// VerifyExecutionNumber is a free data retrieval call binding the contract method 0xcbb60852.
//
// Solidity: function verifyExecutionNumber(uint64 timestamp, uint64 executionNumber, bytes32[] executionNumberProof) view returns()
func (_BeaconVerifier *BeaconVerifierCallerSession) VerifyExecutionNumber(timestamp uint64, executionNumber uint64, executionNumberProof [][32]byte) error {
	return _BeaconVerifier.Contract.VerifyExecutionNumber(&_BeaconVerifier.CallOpts, timestamp, executionNumber, executionNumberProof)
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

// SetExecutionFeeRecipientGIndex is a paid mutator transaction binding the contract method 0xc0c27230.
//
// Solidity: function setExecutionFeeRecipientGIndex(uint256 _executionFeeRecipientGIndex) returns()
func (_BeaconVerifier *BeaconVerifierTransactor) SetExecutionFeeRecipientGIndex(opts *bind.TransactOpts, _executionFeeRecipientGIndex *big.Int) (*types.Transaction, error) {
	return _BeaconVerifier.contract.Transact(opts, "setExecutionFeeRecipientGIndex", _executionFeeRecipientGIndex)
}

// SetExecutionFeeRecipientGIndex is a paid mutator transaction binding the contract method 0xc0c27230.
//
// Solidity: function setExecutionFeeRecipientGIndex(uint256 _executionFeeRecipientGIndex) returns()
func (_BeaconVerifier *BeaconVerifierSession) SetExecutionFeeRecipientGIndex(_executionFeeRecipientGIndex *big.Int) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.SetExecutionFeeRecipientGIndex(&_BeaconVerifier.TransactOpts, _executionFeeRecipientGIndex)
}

// SetExecutionFeeRecipientGIndex is a paid mutator transaction binding the contract method 0xc0c27230.
//
// Solidity: function setExecutionFeeRecipientGIndex(uint256 _executionFeeRecipientGIndex) returns()
func (_BeaconVerifier *BeaconVerifierTransactorSession) SetExecutionFeeRecipientGIndex(_executionFeeRecipientGIndex *big.Int) (*types.Transaction, error) {
	return _BeaconVerifier.Contract.SetExecutionFeeRecipientGIndex(&_BeaconVerifier.TransactOpts, _executionFeeRecipientGIndex)
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

// BeaconVerifierExecutionFeeRecipientGIndexChangedIterator is returned from FilterExecutionFeeRecipientGIndexChanged and is used to iterate over the raw logs and unpacked data for ExecutionFeeRecipientGIndexChanged events raised by the BeaconVerifier contract.
type BeaconVerifierExecutionFeeRecipientGIndexChangedIterator struct {
	Event *BeaconVerifierExecutionFeeRecipientGIndexChanged // Event containing the contract specifics and raw log

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
func (it *BeaconVerifierExecutionFeeRecipientGIndexChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconVerifierExecutionFeeRecipientGIndexChanged)
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
		it.Event = new(BeaconVerifierExecutionFeeRecipientGIndexChanged)
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
func (it *BeaconVerifierExecutionFeeRecipientGIndexChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconVerifierExecutionFeeRecipientGIndexChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconVerifierExecutionFeeRecipientGIndexChanged represents a ExecutionFeeRecipientGIndexChanged event raised by the BeaconVerifier contract.
type BeaconVerifierExecutionFeeRecipientGIndexChanged struct {
	NewExecutionFeeRecipientGIndex *big.Int
	Raw                            types.Log // Blockchain specific contextual infos
}

// FilterExecutionFeeRecipientGIndexChanged is a free log retrieval operation binding the contract event 0x15e42eec45edd1a051ce50f823a5d6482237d402c471639da37b62e855301548.
//
// Solidity: event ExecutionFeeRecipientGIndexChanged(uint256 newExecutionFeeRecipientGIndex)
func (_BeaconVerifier *BeaconVerifierFilterer) FilterExecutionFeeRecipientGIndexChanged(opts *bind.FilterOpts) (*BeaconVerifierExecutionFeeRecipientGIndexChangedIterator, error) {

	logs, sub, err := _BeaconVerifier.contract.FilterLogs(opts, "ExecutionFeeRecipientGIndexChanged")
	if err != nil {
		return nil, err
	}
	return &BeaconVerifierExecutionFeeRecipientGIndexChangedIterator{contract: _BeaconVerifier.contract, event: "ExecutionFeeRecipientGIndexChanged", logs: logs, sub: sub}, nil
}

// WatchExecutionFeeRecipientGIndexChanged is a free log subscription operation binding the contract event 0x15e42eec45edd1a051ce50f823a5d6482237d402c471639da37b62e855301548.
//
// Solidity: event ExecutionFeeRecipientGIndexChanged(uint256 newExecutionFeeRecipientGIndex)
func (_BeaconVerifier *BeaconVerifierFilterer) WatchExecutionFeeRecipientGIndexChanged(opts *bind.WatchOpts, sink chan<- *BeaconVerifierExecutionFeeRecipientGIndexChanged) (event.Subscription, error) {

	logs, sub, err := _BeaconVerifier.contract.WatchLogs(opts, "ExecutionFeeRecipientGIndexChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconVerifierExecutionFeeRecipientGIndexChanged)
				if err := _BeaconVerifier.contract.UnpackLog(event, "ExecutionFeeRecipientGIndexChanged", log); err != nil {
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

// ParseExecutionFeeRecipientGIndexChanged is a log parse operation binding the contract event 0x15e42eec45edd1a051ce50f823a5d6482237d402c471639da37b62e855301548.
//
// Solidity: event ExecutionFeeRecipientGIndexChanged(uint256 newExecutionFeeRecipientGIndex)
func (_BeaconVerifier *BeaconVerifierFilterer) ParseExecutionFeeRecipientGIndexChanged(log types.Log) (*BeaconVerifierExecutionFeeRecipientGIndexChanged, error) {
	event := new(BeaconVerifierExecutionFeeRecipientGIndexChanged)
	if err := _BeaconVerifier.contract.UnpackLog(event, "ExecutionFeeRecipientGIndexChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
