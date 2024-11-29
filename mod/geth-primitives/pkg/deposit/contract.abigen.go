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
	ABI: "[{\"type\":\"function\",\"name\":\"acceptOperatorChange\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelOperatorChange\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"credentials\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"genesisDepositsRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperator\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"queuedOperator\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"queuedTimestamp\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"newOperator\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestOperatorChange\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"newOperator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"credentials\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"index\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorChangeCancelled\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorChangeQueued\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"queuedOperator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"currentOperator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"queuedTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorUpdated\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"newOperator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"previousOperator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"DepositNotMultipleOfGwei\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DepositValueTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientDeposit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCredentialsLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPubKeyLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignatureLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotEnoughTime\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotNewOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OperatorAlreadySet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroOperatorOnFirstDeposit\",\"inputs\":[]}]",
	Bin: "0x6080604052348015600e575f80fd5b506110f38061001c5f395ff3fe608060405260043610610093575f3560e01c8063577212fe11610066578063c53925d91161004c578063c53925d914610231578063e12cf4cb14610250578063fea7ab7714610263575f80fd5b8063577212fe146101cc5780639eaffa96146101ed575f80fd5b806301ffc9a7146100975780632dfdf0b5146100cb5780633523f9bd14610103578063560036ec14610126575b5f80fd5b3480156100a2575f80fd5b506100b66100b1366004610c23565b610282565b60405190151581526020015b60405180910390f35b3480156100d6575f80fd5b505f546100ea9067ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100c2565b34801561010e575f80fd5b5061011860015481565b6040519081526020016100c2565b348015610131575f80fd5b50610193610140366004610c96565b80516020818301810180516003825292820191909301209152546bffffffffffffffffffffffff8116906c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1682565b604080516bffffffffffffffffffffffff909316835273ffffffffffffffffffffffffffffffffffffffff9091166020830152016100c2565b3480156101d7575f80fd5b506101eb6101e6366004610dcb565b61031a565b005b3480156101f8575f80fd5b5061020c610207366004610dcb565b6103f0565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100c2565b34801561023c575f80fd5b506101eb61024b366004610dcb565b610431565b6101eb61025e366004610e2d565b610658565b34801561026e575f80fd5b506101eb61027d366004610edc565b6109ac565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000148061031457507fffffffff0000000000000000000000000000000000000000000000000000000082167f136f920d00000000000000000000000000000000000000000000000000000000145b92915050565b6002828260405161032c929190610f2c565b908152604051908190036020019020543373ffffffffffffffffffffffffffffffffffffffff9091161461038c576040517f7c214f0400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6003828260405161039e929190610f2c565b9081526040519081900360200181205f90556103bd9083908390610f2c565b604051908190038120907f1c0a7e1bd09da292425c039309671a03de56b89a0858598aab6df6ce84b006db905f90a25050565b5f60028383604051610403929190610f2c565b9081526040519081900360200190205473ffffffffffffffffffffffffffffffffffffffff16905092915050565b5f60038383604051610444929190610f2c565b908152604051908190036020019020805490915073ffffffffffffffffffffffffffffffffffffffff6c01000000000000000000000000820416906bffffffffffffffffffffffff163382146104c6576040517f819a0d0b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6bffffffffffffffffffffffff42166104e26201518083610f68565b6bffffffffffffffffffffffff161115610528576040517fe8966d7a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6002868660405161053b929190610f2c565b9081526040519081900360200181205473ffffffffffffffffffffffffffffffffffffffff16915083906002906105759089908990610f2c565b908152604051908190036020018120805473ffffffffffffffffffffffffffffffffffffffff939093167fffffffffffffffffffffffff0000000000000000000000000000000000000000909316929092179091556003906105da9088908890610f2c565b9081526040519081900360200181205f90556105f99087908790610f2c565b6040805191829003822073ffffffffffffffffffffffffffffffffffffffff808716845284166020840152917f0adffd98d3072c48341843974dffd7a910bb849ba6ca04163d43bb26feb17403910160405180910390a2505050505050565b60308614610692576040517f9f10647200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602084146106cc576040517fb39bca1600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608214610706576040517f4be6321b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff166002888860405161072f929190610f2c565b9081526040519081900360200190205473ffffffffffffffffffffffffffffffffffffffff16036108765773ffffffffffffffffffffffffffffffffffffffff81166107a7576040517f51969a7a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600288886040516107ba929190610f2c565b908152604051908190036020018120805473ffffffffffffffffffffffffffffffffffffffff939093167fffffffffffffffffffffffff00000000000000000000000000000000000000009093169290921790915561081c9088908890610f2c565b6040805191829003822073ffffffffffffffffffffffffffffffffffffffff841683525f6020840152917f0adffd98d3072c48341843974dffd7a910bb849ba6ca04163d43bb26feb17403910160405180910390a26108c4565b73ffffffffffffffffffffffffffffffffffffffff8116156108c4576040517fc4142b4100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6108cd610b5e565b905064077359400067ffffffffffffffff82161015610918576040517f0e1eddda00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f80547f68af751683498a9f9be59fe8b0d52a64dd155255d85cdb29fea30b1e3f891d46918a918a918a918a9187918b918b9167ffffffffffffffff16908061096083610f8c565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060405161099a989796959493929190610fff565b60405180910390a15050505050505050565b5f600284846040516109bf929190610f2c565b9081526040519081900360200190205473ffffffffffffffffffffffffffffffffffffffff169050338114610a20576040517f7c214f0400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8216610a6d576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60038585604051610a80929190610f2c565b908152604051908190036020018120426bffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff86166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffffff000000000000000000000000161781559150610af89086908690610f2c565b6040805191829003822073ffffffffffffffffffffffffffffffffffffffff8681168452851660208401524283830152905190917f7640ec3c8c4695deadda414dd20400acf275297a7c38715f9237657e97ddba5f919081900360600190a25050505050565b5f610b6d633b9aca0034611097565b15610ba4576040517f40567b3800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610bb3633b9aca00346110aa565b905067ffffffffffffffff811115610bf7576040517f2aa6673400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610c015f34610c06565b919050565b5f385f3884865af1610c1f5763b12d13eb5f526004601cfd5b5050565b5f60208284031215610c33575f80fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610c62575f80fd5b9392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f60208284031215610ca6575f80fd5b813567ffffffffffffffff811115610cbc575f80fd5b8201601f81018413610ccc575f80fd5b803567ffffffffffffffff811115610ce657610ce6610c69565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff82111715610d5257610d52610c69565b604052818152828201602001861015610d69575f80fd5b816020840160208301375f91810160200191909152949350505050565b5f8083601f840112610d96575f80fd5b50813567ffffffffffffffff811115610dad575f80fd5b602083019150836020828501011115610dc4575f80fd5b9250929050565b5f8060208385031215610ddc575f80fd5b823567ffffffffffffffff811115610df2575f80fd5b610dfe85828601610d86565b90969095509350505050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610c01575f80fd5b5f805f805f805f6080888a031215610e43575f80fd5b873567ffffffffffffffff811115610e59575f80fd5b610e658a828b01610d86565b909850965050602088013567ffffffffffffffff811115610e84575f80fd5b610e908a828b01610d86565b909650945050604088013567ffffffffffffffff811115610eaf575f80fd5b610ebb8a828b01610d86565b9094509250610ece905060608901610e0a565b905092959891949750929550565b5f805f60408486031215610eee575f80fd5b833567ffffffffffffffff811115610f04575f80fd5b610f1086828701610d86565b9094509250610f23905060208501610e0a565b90509250925092565b818382375f9101908152919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b6bffffffffffffffffffffffff818116838216019081111561031457610314610f3b565b5f67ffffffffffffffff821667ffffffffffffffff8103610faf57610faf610f3b565b60010192915050565b81835281816020850137505f602082840101525f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60a081525f61101260a083018a8c610fb8565b828103602084015261102581898b610fb8565b905067ffffffffffffffff87166040840152828103606084015261104a818688610fb8565b91505067ffffffffffffffff831660808301529998505050505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f826110a5576110a561106a565b500690565b5f826110b8576110b861106a565b50049056fea26469706673582212207d2196628e1095775bc63abfb75c5f1b26538900aa683597ea8934159f1f0f6d64736f6c634300081a0033",
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

// GenesisDepositsRoot is a free data retrieval call binding the contract method 0x3523f9bd.
//
// Solidity: function genesisDepositsRoot() view returns(bytes32)
func (_BeaconDepositContract *BeaconDepositContractCaller) GenesisDepositsRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BeaconDepositContract.contract.Call(opts, &out, "genesisDepositsRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GenesisDepositsRoot is a free data retrieval call binding the contract method 0x3523f9bd.
//
// Solidity: function genesisDepositsRoot() view returns(bytes32)
func (_BeaconDepositContract *BeaconDepositContractSession) GenesisDepositsRoot() ([32]byte, error) {
	return _BeaconDepositContract.Contract.GenesisDepositsRoot(&_BeaconDepositContract.CallOpts)
}

// GenesisDepositsRoot is a free data retrieval call binding the contract method 0x3523f9bd.
//
// Solidity: function genesisDepositsRoot() view returns(bytes32)
func (_BeaconDepositContract *BeaconDepositContractCallerSession) GenesisDepositsRoot() ([32]byte, error) {
	return _BeaconDepositContract.Contract.GenesisDepositsRoot(&_BeaconDepositContract.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0x9eaffa96.
//
// Solidity: function getOperator(bytes pubkey) view returns(address)
func (_BeaconDepositContract *BeaconDepositContractCaller) GetOperator(opts *bind.CallOpts, pubkey []byte) (common.Address, error) {
	var out []interface{}
	err := _BeaconDepositContract.contract.Call(opts, &out, "getOperator", pubkey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0x9eaffa96.
//
// Solidity: function getOperator(bytes pubkey) view returns(address)
func (_BeaconDepositContract *BeaconDepositContractSession) GetOperator(pubkey []byte) (common.Address, error) {
	return _BeaconDepositContract.Contract.GetOperator(&_BeaconDepositContract.CallOpts, pubkey)
}

// GetOperator is a free data retrieval call binding the contract method 0x9eaffa96.
//
// Solidity: function getOperator(bytes pubkey) view returns(address)
func (_BeaconDepositContract *BeaconDepositContractCallerSession) GetOperator(pubkey []byte) (common.Address, error) {
	return _BeaconDepositContract.Contract.GetOperator(&_BeaconDepositContract.CallOpts, pubkey)
}

// QueuedOperator is a free data retrieval call binding the contract method 0x560036ec.
//
// Solidity: function queuedOperator(bytes ) view returns(uint96 queuedTimestamp, address newOperator)
func (_BeaconDepositContract *BeaconDepositContractCaller) QueuedOperator(opts *bind.CallOpts, arg0 []byte) (struct {
	QueuedTimestamp *big.Int
	NewOperator     common.Address
}, error) {
	var out []interface{}
	err := _BeaconDepositContract.contract.Call(opts, &out, "queuedOperator", arg0)

	outstruct := new(struct {
		QueuedTimestamp *big.Int
		NewOperator     common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.QueuedTimestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.NewOperator = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// QueuedOperator is a free data retrieval call binding the contract method 0x560036ec.
//
// Solidity: function queuedOperator(bytes ) view returns(uint96 queuedTimestamp, address newOperator)
func (_BeaconDepositContract *BeaconDepositContractSession) QueuedOperator(arg0 []byte) (struct {
	QueuedTimestamp *big.Int
	NewOperator     common.Address
}, error) {
	return _BeaconDepositContract.Contract.QueuedOperator(&_BeaconDepositContract.CallOpts, arg0)
}

// QueuedOperator is a free data retrieval call binding the contract method 0x560036ec.
//
// Solidity: function queuedOperator(bytes ) view returns(uint96 queuedTimestamp, address newOperator)
func (_BeaconDepositContract *BeaconDepositContractCallerSession) QueuedOperator(arg0 []byte) (struct {
	QueuedTimestamp *big.Int
	NewOperator     common.Address
}, error) {
	return _BeaconDepositContract.Contract.QueuedOperator(&_BeaconDepositContract.CallOpts, arg0)
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

// AcceptOperatorChange is a paid mutator transaction binding the contract method 0xc53925d9.
//
// Solidity: function acceptOperatorChange(bytes pubkey) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) AcceptOperatorChange(opts *bind.TransactOpts, pubkey []byte) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "acceptOperatorChange", pubkey)
}

// AcceptOperatorChange is a paid mutator transaction binding the contract method 0xc53925d9.
//
// Solidity: function acceptOperatorChange(bytes pubkey) returns()
func (_BeaconDepositContract *BeaconDepositContractSession) AcceptOperatorChange(pubkey []byte) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.AcceptOperatorChange(&_BeaconDepositContract.TransactOpts, pubkey)
}

// AcceptOperatorChange is a paid mutator transaction binding the contract method 0xc53925d9.
//
// Solidity: function acceptOperatorChange(bytes pubkey) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) AcceptOperatorChange(pubkey []byte) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.AcceptOperatorChange(&_BeaconDepositContract.TransactOpts, pubkey)
}

// CancelOperatorChange is a paid mutator transaction binding the contract method 0x577212fe.
//
// Solidity: function cancelOperatorChange(bytes pubkey) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) CancelOperatorChange(opts *bind.TransactOpts, pubkey []byte) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "cancelOperatorChange", pubkey)
}

// CancelOperatorChange is a paid mutator transaction binding the contract method 0x577212fe.
//
// Solidity: function cancelOperatorChange(bytes pubkey) returns()
func (_BeaconDepositContract *BeaconDepositContractSession) CancelOperatorChange(pubkey []byte) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.CancelOperatorChange(&_BeaconDepositContract.TransactOpts, pubkey)
}

// CancelOperatorChange is a paid mutator transaction binding the contract method 0x577212fe.
//
// Solidity: function cancelOperatorChange(bytes pubkey) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) CancelOperatorChange(pubkey []byte) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.CancelOperatorChange(&_BeaconDepositContract.TransactOpts, pubkey)
}

// Deposit is a paid mutator transaction binding the contract method 0xe12cf4cb.
//
// Solidity: function deposit(bytes pubkey, bytes credentials, bytes signature, address operator) payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) Deposit(opts *bind.TransactOpts, pubkey []byte, credentials []byte, signature []byte, operator common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "deposit", pubkey, credentials, signature, operator)
}

// Deposit is a paid mutator transaction binding the contract method 0xe12cf4cb.
//
// Solidity: function deposit(bytes pubkey, bytes credentials, bytes signature, address operator) payable returns()
func (_BeaconDepositContract *BeaconDepositContractSession) Deposit(pubkey []byte, credentials []byte, signature []byte, operator common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.Deposit(&_BeaconDepositContract.TransactOpts, pubkey, credentials, signature, operator)
}

// Deposit is a paid mutator transaction binding the contract method 0xe12cf4cb.
//
// Solidity: function deposit(bytes pubkey, bytes credentials, bytes signature, address operator) payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) Deposit(pubkey []byte, credentials []byte, signature []byte, operator common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.Deposit(&_BeaconDepositContract.TransactOpts, pubkey, credentials, signature, operator)
}

// RequestOperatorChange is a paid mutator transaction binding the contract method 0xfea7ab77.
//
// Solidity: function requestOperatorChange(bytes pubkey, address newOperator) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) RequestOperatorChange(opts *bind.TransactOpts, pubkey []byte, newOperator common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "requestOperatorChange", pubkey, newOperator)
}

// RequestOperatorChange is a paid mutator transaction binding the contract method 0xfea7ab77.
//
// Solidity: function requestOperatorChange(bytes pubkey, address newOperator) returns()
func (_BeaconDepositContract *BeaconDepositContractSession) RequestOperatorChange(pubkey []byte, newOperator common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.RequestOperatorChange(&_BeaconDepositContract.TransactOpts, pubkey, newOperator)
}

// RequestOperatorChange is a paid mutator transaction binding the contract method 0xfea7ab77.
//
// Solidity: function requestOperatorChange(bytes pubkey, address newOperator) returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) RequestOperatorChange(pubkey []byte, newOperator common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.RequestOperatorChange(&_BeaconDepositContract.TransactOpts, pubkey, newOperator)
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

// BeaconDepositContractOperatorChangeCancelledIterator is returned from FilterOperatorChangeCancelled and is used to iterate over the raw logs and unpacked data for OperatorChangeCancelled events raised by the BeaconDepositContract contract.
type BeaconDepositContractOperatorChangeCancelledIterator struct {
	Event *BeaconDepositContractOperatorChangeCancelled // Event containing the contract specifics and raw log

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
func (it *BeaconDepositContractOperatorChangeCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconDepositContractOperatorChangeCancelled)
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
		it.Event = new(BeaconDepositContractOperatorChangeCancelled)
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
func (it *BeaconDepositContractOperatorChangeCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconDepositContractOperatorChangeCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconDepositContractOperatorChangeCancelled represents a OperatorChangeCancelled event raised by the BeaconDepositContract contract.
type BeaconDepositContractOperatorChangeCancelled struct {
	Pubkey common.Hash
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOperatorChangeCancelled is a free log retrieval operation binding the contract event 0x1c0a7e1bd09da292425c039309671a03de56b89a0858598aab6df6ce84b006db.
//
// Solidity: event OperatorChangeCancelled(bytes indexed pubkey)
func (_BeaconDepositContract *BeaconDepositContractFilterer) FilterOperatorChangeCancelled(opts *bind.FilterOpts, pubkey [][]byte) (*BeaconDepositContractOperatorChangeCancelledIterator, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _BeaconDepositContract.contract.FilterLogs(opts, "OperatorChangeCancelled", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractOperatorChangeCancelledIterator{contract: _BeaconDepositContract.contract, event: "OperatorChangeCancelled", logs: logs, sub: sub}, nil
}

// WatchOperatorChangeCancelled is a free log subscription operation binding the contract event 0x1c0a7e1bd09da292425c039309671a03de56b89a0858598aab6df6ce84b006db.
//
// Solidity: event OperatorChangeCancelled(bytes indexed pubkey)
func (_BeaconDepositContract *BeaconDepositContractFilterer) WatchOperatorChangeCancelled(opts *bind.WatchOpts, sink chan<- *BeaconDepositContractOperatorChangeCancelled, pubkey [][]byte) (event.Subscription, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _BeaconDepositContract.contract.WatchLogs(opts, "OperatorChangeCancelled", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconDepositContractOperatorChangeCancelled)
				if err := _BeaconDepositContract.contract.UnpackLog(event, "OperatorChangeCancelled", log); err != nil {
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

// ParseOperatorChangeCancelled is a log parse operation binding the contract event 0x1c0a7e1bd09da292425c039309671a03de56b89a0858598aab6df6ce84b006db.
//
// Solidity: event OperatorChangeCancelled(bytes indexed pubkey)
func (_BeaconDepositContract *BeaconDepositContractFilterer) ParseOperatorChangeCancelled(log types.Log) (*BeaconDepositContractOperatorChangeCancelled, error) {
	event := new(BeaconDepositContractOperatorChangeCancelled)
	if err := _BeaconDepositContract.contract.UnpackLog(event, "OperatorChangeCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconDepositContractOperatorChangeQueuedIterator is returned from FilterOperatorChangeQueued and is used to iterate over the raw logs and unpacked data for OperatorChangeQueued events raised by the BeaconDepositContract contract.
type BeaconDepositContractOperatorChangeQueuedIterator struct {
	Event *BeaconDepositContractOperatorChangeQueued // Event containing the contract specifics and raw log

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
func (it *BeaconDepositContractOperatorChangeQueuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconDepositContractOperatorChangeQueued)
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
		it.Event = new(BeaconDepositContractOperatorChangeQueued)
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
func (it *BeaconDepositContractOperatorChangeQueuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconDepositContractOperatorChangeQueuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconDepositContractOperatorChangeQueued represents a OperatorChangeQueued event raised by the BeaconDepositContract contract.
type BeaconDepositContractOperatorChangeQueued struct {
	Pubkey          common.Hash
	QueuedOperator  common.Address
	CurrentOperator common.Address
	QueuedTimestamp *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterOperatorChangeQueued is a free log retrieval operation binding the contract event 0x7640ec3c8c4695deadda414dd20400acf275297a7c38715f9237657e97ddba5f.
//
// Solidity: event OperatorChangeQueued(bytes indexed pubkey, address queuedOperator, address currentOperator, uint256 queuedTimestamp)
func (_BeaconDepositContract *BeaconDepositContractFilterer) FilterOperatorChangeQueued(opts *bind.FilterOpts, pubkey [][]byte) (*BeaconDepositContractOperatorChangeQueuedIterator, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _BeaconDepositContract.contract.FilterLogs(opts, "OperatorChangeQueued", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractOperatorChangeQueuedIterator{contract: _BeaconDepositContract.contract, event: "OperatorChangeQueued", logs: logs, sub: sub}, nil
}

// WatchOperatorChangeQueued is a free log subscription operation binding the contract event 0x7640ec3c8c4695deadda414dd20400acf275297a7c38715f9237657e97ddba5f.
//
// Solidity: event OperatorChangeQueued(bytes indexed pubkey, address queuedOperator, address currentOperator, uint256 queuedTimestamp)
func (_BeaconDepositContract *BeaconDepositContractFilterer) WatchOperatorChangeQueued(opts *bind.WatchOpts, sink chan<- *BeaconDepositContractOperatorChangeQueued, pubkey [][]byte) (event.Subscription, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _BeaconDepositContract.contract.WatchLogs(opts, "OperatorChangeQueued", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconDepositContractOperatorChangeQueued)
				if err := _BeaconDepositContract.contract.UnpackLog(event, "OperatorChangeQueued", log); err != nil {
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

// ParseOperatorChangeQueued is a log parse operation binding the contract event 0x7640ec3c8c4695deadda414dd20400acf275297a7c38715f9237657e97ddba5f.
//
// Solidity: event OperatorChangeQueued(bytes indexed pubkey, address queuedOperator, address currentOperator, uint256 queuedTimestamp)
func (_BeaconDepositContract *BeaconDepositContractFilterer) ParseOperatorChangeQueued(log types.Log) (*BeaconDepositContractOperatorChangeQueued, error) {
	event := new(BeaconDepositContractOperatorChangeQueued)
	if err := _BeaconDepositContract.contract.UnpackLog(event, "OperatorChangeQueued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BeaconDepositContractOperatorUpdatedIterator is returned from FilterOperatorUpdated and is used to iterate over the raw logs and unpacked data for OperatorUpdated events raised by the BeaconDepositContract contract.
type BeaconDepositContractOperatorUpdatedIterator struct {
	Event *BeaconDepositContractOperatorUpdated // Event containing the contract specifics and raw log

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
func (it *BeaconDepositContractOperatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconDepositContractOperatorUpdated)
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
		it.Event = new(BeaconDepositContractOperatorUpdated)
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
func (it *BeaconDepositContractOperatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconDepositContractOperatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconDepositContractOperatorUpdated represents a OperatorUpdated event raised by the BeaconDepositContract contract.
type BeaconDepositContractOperatorUpdated struct {
	Pubkey           common.Hash
	NewOperator      common.Address
	PreviousOperator common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOperatorUpdated is a free log retrieval operation binding the contract event 0x0adffd98d3072c48341843974dffd7a910bb849ba6ca04163d43bb26feb17403.
//
// Solidity: event OperatorUpdated(bytes indexed pubkey, address newOperator, address previousOperator)
func (_BeaconDepositContract *BeaconDepositContractFilterer) FilterOperatorUpdated(opts *bind.FilterOpts, pubkey [][]byte) (*BeaconDepositContractOperatorUpdatedIterator, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _BeaconDepositContract.contract.FilterLogs(opts, "OperatorUpdated", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return &BeaconDepositContractOperatorUpdatedIterator{contract: _BeaconDepositContract.contract, event: "OperatorUpdated", logs: logs, sub: sub}, nil
}

// WatchOperatorUpdated is a free log subscription operation binding the contract event 0x0adffd98d3072c48341843974dffd7a910bb849ba6ca04163d43bb26feb17403.
//
// Solidity: event OperatorUpdated(bytes indexed pubkey, address newOperator, address previousOperator)
func (_BeaconDepositContract *BeaconDepositContractFilterer) WatchOperatorUpdated(opts *bind.WatchOpts, sink chan<- *BeaconDepositContractOperatorUpdated, pubkey [][]byte) (event.Subscription, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _BeaconDepositContract.contract.WatchLogs(opts, "OperatorUpdated", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconDepositContractOperatorUpdated)
				if err := _BeaconDepositContract.contract.UnpackLog(event, "OperatorUpdated", log); err != nil {
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

// ParseOperatorUpdated is a log parse operation binding the contract event 0x0adffd98d3072c48341843974dffd7a910bb849ba6ca04163d43bb26feb17403.
//
// Solidity: event OperatorUpdated(bytes indexed pubkey, address newOperator, address previousOperator)
func (_BeaconDepositContract *BeaconDepositContractFilterer) ParseOperatorUpdated(log types.Log) (*BeaconDepositContractOperatorUpdated, error) {
	event := new(BeaconDepositContractOperatorUpdated)
	if err := _BeaconDepositContract.contract.UnpackLog(event, "OperatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
