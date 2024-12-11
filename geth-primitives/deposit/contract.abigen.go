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

// DepositContractMetaData contains all meta data concerning the DepositContract contract.
var DepositContractMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"acceptOperatorChange\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelOperatorChange\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"credentials\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"genesisDepositsRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperator\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"queuedOperator\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"queuedTimestamp\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"newOperator\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestOperatorChange\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"newOperator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"credentials\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"index\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorChangeCancelled\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorChangeQueued\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"queuedOperator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"currentOperator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"queuedTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorUpdated\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"newOperator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"previousOperator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"DepositNotMultipleOfGwei\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DepositValueTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientDeposit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCredentialsLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPubKeyLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignatureLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotEnoughTime\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotNewOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OperatorAlreadySet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroOperatorOnFirstDeposit\",\"inputs\":[]}]",
	Bin: "0x6080604052348015600e575f80fd5b506110f28061001c5f395ff3fe608060405260043610610093575f3560e01c8063577212fe11610066578063c53925d91161004c578063c53925d914610231578063e12cf4cb14610250578063fea7ab7714610263575f80fd5b8063577212fe146101cc5780639eaffa96146101ed575f80fd5b806301ffc9a7146100975780632dfdf0b5146100cb5780633523f9bd14610103578063560036ec14610126575b5f80fd5b3480156100a2575f80fd5b506100b66100b1366004610c22565b610282565b60405190151581526020015b60405180910390f35b3480156100d6575f80fd5b505f546100ea9067ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016100c2565b34801561010e575f80fd5b5061011860015481565b6040519081526020016100c2565b348015610131575f80fd5b50610193610140366004610c95565b80516020818301810180516003825292820191909301209152546bffffffffffffffffffffffff8116906c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1682565b604080516bffffffffffffffffffffffff909316835273ffffffffffffffffffffffffffffffffffffffff9091166020830152016100c2565b3480156101d7575f80fd5b506101eb6101e6366004610dca565b61031a565b005b3480156101f8575f80fd5b5061020c610207366004610dca565b6103f0565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100c2565b34801561023c575f80fd5b506101eb61024b366004610dca565b610431565b6101eb61025e366004610e2c565b610658565b34801561026e575f80fd5b506101eb61027d366004610edb565b6109ab565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000148061031457507fffffffff0000000000000000000000000000000000000000000000000000000082167f136f920d00000000000000000000000000000000000000000000000000000000145b92915050565b6002828260405161032c929190610f2b565b908152604051908190036020019020543373ffffffffffffffffffffffffffffffffffffffff9091161461038c576040517f7c214f0400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6003828260405161039e929190610f2b565b9081526040519081900360200181205f90556103bd9083908390610f2b565b604051908190038120907f1c0a7e1bd09da292425c039309671a03de56b89a0858598aab6df6ce84b006db905f90a25050565b5f60028383604051610403929190610f2b565b9081526040519081900360200190205473ffffffffffffffffffffffffffffffffffffffff16905092915050565b5f60038383604051610444929190610f2b565b908152604051908190036020019020805490915073ffffffffffffffffffffffffffffffffffffffff6c01000000000000000000000000820416906bffffffffffffffffffffffff163382146104c6576040517f819a0d0b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6bffffffffffffffffffffffff42166104e26201518083610f67565b6bffffffffffffffffffffffff161115610528576040517fe8966d7a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6002868660405161053b929190610f2b565b9081526040519081900360200181205473ffffffffffffffffffffffffffffffffffffffff16915083906002906105759089908990610f2b565b908152604051908190036020018120805473ffffffffffffffffffffffffffffffffffffffff939093167fffffffffffffffffffffffff0000000000000000000000000000000000000000909316929092179091556003906105da9088908890610f2b565b9081526040519081900360200181205f90556105f99087908790610f2b565b6040805191829003822073ffffffffffffffffffffffffffffffffffffffff808716845284166020840152917f0adffd98d3072c48341843974dffd7a910bb849ba6ca04163d43bb26feb17403910160405180910390a2505050505050565b60308614610692576040517f9f10647200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b602084146106cc576040517fb39bca1600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608214610706576040517f4be6321b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff166002888860405161072f929190610f2b565b9081526040519081900360200190205473ffffffffffffffffffffffffffffffffffffffff16036108765773ffffffffffffffffffffffffffffffffffffffff81166107a7576040517f51969a7a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600288886040516107ba929190610f2b565b908152604051908190036020018120805473ffffffffffffffffffffffffffffffffffffffff939093167fffffffffffffffffffffffff00000000000000000000000000000000000000009093169290921790915561081c9088908890610f2b565b6040805191829003822073ffffffffffffffffffffffffffffffffffffffff841683525f6020840152917f0adffd98d3072c48341843974dffd7a910bb849ba6ca04163d43bb26feb17403910160405180910390a26108c4565b73ffffffffffffffffffffffffffffffffffffffff8116156108c4576040517fc4142b4100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6108cd610b5d565b9050633b9aca0067ffffffffffffffff82161015610917576040517f0e1eddda00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f80547f68af751683498a9f9be59fe8b0d52a64dd155255d85cdb29fea30b1e3f891d46918a918a918a918a9187918b918b9167ffffffffffffffff16908061095f83610f8b565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550604051610999989796959493929190610ffe565b60405180910390a15050505050505050565b5f600284846040516109be929190610f2b565b9081526040519081900360200190205473ffffffffffffffffffffffffffffffffffffffff169050338114610a1f576040517f7c214f0400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8216610a6c576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60038585604051610a7f929190610f2b565b908152604051908190036020018120426bffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff86166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffffff000000000000000000000000161781559150610af79086908690610f2b565b6040805191829003822073ffffffffffffffffffffffffffffffffffffffff8681168452851660208401524283830152905190917f7640ec3c8c4695deadda414dd20400acf275297a7c38715f9237657e97ddba5f919081900360600190a25050505050565b5f610b6c633b9aca0034611096565b15610ba3576040517f40567b3800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610bb2633b9aca00346110a9565b905067ffffffffffffffff811115610bf6576040517f2aa6673400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610c005f34610c05565b919050565b5f385f3884865af1610c1e5763b12d13eb5f526004601cfd5b5050565b5f60208284031215610c32575f80fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610c61575f80fd5b9392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f60208284031215610ca5575f80fd5b813567ffffffffffffffff811115610cbb575f80fd5b8201601f81018413610ccb575f80fd5b803567ffffffffffffffff811115610ce557610ce5610c68565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff82111715610d5157610d51610c68565b604052818152828201602001861015610d68575f80fd5b816020840160208301375f91810160200191909152949350505050565b5f8083601f840112610d95575f80fd5b50813567ffffffffffffffff811115610dac575f80fd5b602083019150836020828501011115610dc3575f80fd5b9250929050565b5f8060208385031215610ddb575f80fd5b823567ffffffffffffffff811115610df1575f80fd5b610dfd85828601610d85565b90969095509350505050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610c00575f80fd5b5f805f805f805f6080888a031215610e42575f80fd5b873567ffffffffffffffff811115610e58575f80fd5b610e648a828b01610d85565b909850965050602088013567ffffffffffffffff811115610e83575f80fd5b610e8f8a828b01610d85565b909650945050604088013567ffffffffffffffff811115610eae575f80fd5b610eba8a828b01610d85565b9094509250610ecd905060608901610e09565b905092959891949750929550565b5f805f60408486031215610eed575f80fd5b833567ffffffffffffffff811115610f03575f80fd5b610f0f86828701610d85565b9094509250610f22905060208501610e09565b90509250925092565b818382375f9101908152919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b6bffffffffffffffffffffffff818116838216019081111561031457610314610f3a565b5f67ffffffffffffffff821667ffffffffffffffff8103610fae57610fae610f3a565b60010192915050565b81835281816020850137505f602082840101525f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60a081525f61101160a083018a8c610fb7565b828103602084015261102481898b610fb7565b905067ffffffffffffffff871660408401528281036060840152611049818688610fb7565b91505067ffffffffffffffff831660808301529998505050505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f826110a4576110a4611069565b500690565b5f826110b7576110b7611069565b50049056fea264697066735822122069227307258cbe8f29985bb4f3e283b1b03d5c0cbab8add81bf3c22e3d13729664736f6c634300081a0033",
}

// DepositContractABI is the input ABI used to generate the binding from.
// Deprecated: Use DepositContractMetaData.ABI instead.
var DepositContractABI = DepositContractMetaData.ABI

// DepositContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DepositContractMetaData.Bin instead.
var DepositContractBin = DepositContractMetaData.Bin

// DeployDepositContract deploys a new Ethereum contract, binding an instance of DepositContract to it.
func DeployDepositContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DepositContract, error) {
	parsed, err := DepositContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DepositContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DepositContract{DepositContractCaller: DepositContractCaller{contract: contract}, DepositContractTransactor: DepositContractTransactor{contract: contract}, DepositContractFilterer: DepositContractFilterer{contract: contract}}, nil
}

// DepositContract is an auto generated Go binding around an Ethereum contract.
type DepositContract struct {
	DepositContractCaller     // Read-only binding to the contract
	DepositContractTransactor // Write-only binding to the contract
	DepositContractFilterer   // Log filterer for contract events
}

// DepositContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositContractSession struct {
	Contract     *DepositContract  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositContractCallerSession struct {
	Contract *DepositContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// DepositContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositContractTransactorSession struct {
	Contract     *DepositContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// DepositContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositContractRaw struct {
	Contract *DepositContract // Generic contract binding to access the raw methods on
}

// DepositContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositContractCallerRaw struct {
	Contract *DepositContractCaller // Generic read-only contract binding to access the raw methods on
}

// DepositContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositContractTransactorRaw struct {
	Contract *DepositContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepositContract creates a new instance of DepositContract, bound to a specific deployed contract.
func NewDepositContract(address common.Address, backend bind.ContractBackend) (*DepositContract, error) {
	contract, err := bindDepositContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DepositContract{DepositContractCaller: DepositContractCaller{contract: contract}, DepositContractTransactor: DepositContractTransactor{contract: contract}, DepositContractFilterer: DepositContractFilterer{contract: contract}}, nil
}

// NewDepositContractCaller creates a new read-only instance of DepositContract, bound to a specific deployed contract.
func NewDepositContractCaller(address common.Address, caller bind.ContractCaller) (*DepositContractCaller, error) {
	contract, err := bindDepositContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositContractCaller{contract: contract}, nil
}

// NewDepositContractTransactor creates a new write-only instance of DepositContract, bound to a specific deployed contract.
func NewDepositContractTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositContractTransactor, error) {
	contract, err := bindDepositContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositContractTransactor{contract: contract}, nil
}

// NewDepositContractFilterer creates a new log filterer instance of DepositContract, bound to a specific deployed contract.
func NewDepositContractFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositContractFilterer, error) {
	contract, err := bindDepositContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositContractFilterer{contract: contract}, nil
}

// bindDepositContract binds a generic wrapper to an already deployed contract.
func bindDepositContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DepositContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositContract *DepositContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DepositContract.Contract.DepositContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositContract *DepositContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositContract.Contract.DepositContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositContract *DepositContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositContract.Contract.DepositContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositContract *DepositContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DepositContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositContract *DepositContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositContract *DepositContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositContract.Contract.contract.Transact(opts, method, params...)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint64)
func (_DepositContract *DepositContractCaller) DepositCount(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _DepositContract.contract.Call(opts, &out, "depositCount")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint64)
func (_DepositContract *DepositContractSession) DepositCount() (uint64, error) {
	return _DepositContract.Contract.DepositCount(&_DepositContract.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint64)
func (_DepositContract *DepositContractCallerSession) DepositCount() (uint64, error) {
	return _DepositContract.Contract.DepositCount(&_DepositContract.CallOpts)
}

// GenesisDepositsRoot is a free data retrieval call binding the contract method 0x3523f9bd.
//
// Solidity: function genesisDepositsRoot() view returns(bytes32)
func (_DepositContract *DepositContractCaller) GenesisDepositsRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DepositContract.contract.Call(opts, &out, "genesisDepositsRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GenesisDepositsRoot is a free data retrieval call binding the contract method 0x3523f9bd.
//
// Solidity: function genesisDepositsRoot() view returns(bytes32)
func (_DepositContract *DepositContractSession) GenesisDepositsRoot() ([32]byte, error) {
	return _DepositContract.Contract.GenesisDepositsRoot(&_DepositContract.CallOpts)
}

// GenesisDepositsRoot is a free data retrieval call binding the contract method 0x3523f9bd.
//
// Solidity: function genesisDepositsRoot() view returns(bytes32)
func (_DepositContract *DepositContractCallerSession) GenesisDepositsRoot() ([32]byte, error) {
	return _DepositContract.Contract.GenesisDepositsRoot(&_DepositContract.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0x9eaffa96.
//
// Solidity: function getOperator(bytes pubkey) view returns(address)
func (_DepositContract *DepositContractCaller) GetOperator(opts *bind.CallOpts, pubkey []byte) (common.Address, error) {
	var out []interface{}
	err := _DepositContract.contract.Call(opts, &out, "getOperator", pubkey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0x9eaffa96.
//
// Solidity: function getOperator(bytes pubkey) view returns(address)
func (_DepositContract *DepositContractSession) GetOperator(pubkey []byte) (common.Address, error) {
	return _DepositContract.Contract.GetOperator(&_DepositContract.CallOpts, pubkey)
}

// GetOperator is a free data retrieval call binding the contract method 0x9eaffa96.
//
// Solidity: function getOperator(bytes pubkey) view returns(address)
func (_DepositContract *DepositContractCallerSession) GetOperator(pubkey []byte) (common.Address, error) {
	return _DepositContract.Contract.GetOperator(&_DepositContract.CallOpts, pubkey)
}

// QueuedOperator is a free data retrieval call binding the contract method 0x560036ec.
//
// Solidity: function queuedOperator(bytes ) view returns(uint96 queuedTimestamp, address newOperator)
func (_DepositContract *DepositContractCaller) QueuedOperator(opts *bind.CallOpts, arg0 []byte) (struct {
	QueuedTimestamp *big.Int
	NewOperator     common.Address
}, error) {
	var out []interface{}
	err := _DepositContract.contract.Call(opts, &out, "queuedOperator", arg0)

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
func (_DepositContract *DepositContractSession) QueuedOperator(arg0 []byte) (struct {
	QueuedTimestamp *big.Int
	NewOperator     common.Address
}, error) {
	return _DepositContract.Contract.QueuedOperator(&_DepositContract.CallOpts, arg0)
}

// QueuedOperator is a free data retrieval call binding the contract method 0x560036ec.
//
// Solidity: function queuedOperator(bytes ) view returns(uint96 queuedTimestamp, address newOperator)
func (_DepositContract *DepositContractCallerSession) QueuedOperator(arg0 []byte) (struct {
	QueuedTimestamp *big.Int
	NewOperator     common.Address
}, error) {
	return _DepositContract.Contract.QueuedOperator(&_DepositContract.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_DepositContract *DepositContractCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _DepositContract.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_DepositContract *DepositContractSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DepositContract.Contract.SupportsInterface(&_DepositContract.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_DepositContract *DepositContractCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DepositContract.Contract.SupportsInterface(&_DepositContract.CallOpts, interfaceId)
}

// AcceptOperatorChange is a paid mutator transaction binding the contract method 0xc53925d9.
//
// Solidity: function acceptOperatorChange(bytes pubkey) returns()
func (_DepositContract *DepositContractTransactor) AcceptOperatorChange(opts *bind.TransactOpts, pubkey []byte) (*types.Transaction, error) {
	return _DepositContract.contract.Transact(opts, "acceptOperatorChange", pubkey)
}

// AcceptOperatorChange is a paid mutator transaction binding the contract method 0xc53925d9.
//
// Solidity: function acceptOperatorChange(bytes pubkey) returns()
func (_DepositContract *DepositContractSession) AcceptOperatorChange(pubkey []byte) (*types.Transaction, error) {
	return _DepositContract.Contract.AcceptOperatorChange(&_DepositContract.TransactOpts, pubkey)
}

// AcceptOperatorChange is a paid mutator transaction binding the contract method 0xc53925d9.
//
// Solidity: function acceptOperatorChange(bytes pubkey) returns()
func (_DepositContract *DepositContractTransactorSession) AcceptOperatorChange(pubkey []byte) (*types.Transaction, error) {
	return _DepositContract.Contract.AcceptOperatorChange(&_DepositContract.TransactOpts, pubkey)
}

// CancelOperatorChange is a paid mutator transaction binding the contract method 0x577212fe.
//
// Solidity: function cancelOperatorChange(bytes pubkey) returns()
func (_DepositContract *DepositContractTransactor) CancelOperatorChange(opts *bind.TransactOpts, pubkey []byte) (*types.Transaction, error) {
	return _DepositContract.contract.Transact(opts, "cancelOperatorChange", pubkey)
}

// CancelOperatorChange is a paid mutator transaction binding the contract method 0x577212fe.
//
// Solidity: function cancelOperatorChange(bytes pubkey) returns()
func (_DepositContract *DepositContractSession) CancelOperatorChange(pubkey []byte) (*types.Transaction, error) {
	return _DepositContract.Contract.CancelOperatorChange(&_DepositContract.TransactOpts, pubkey)
}

// CancelOperatorChange is a paid mutator transaction binding the contract method 0x577212fe.
//
// Solidity: function cancelOperatorChange(bytes pubkey) returns()
func (_DepositContract *DepositContractTransactorSession) CancelOperatorChange(pubkey []byte) (*types.Transaction, error) {
	return _DepositContract.Contract.CancelOperatorChange(&_DepositContract.TransactOpts, pubkey)
}

// Deposit is a paid mutator transaction binding the contract method 0xe12cf4cb.
//
// Solidity: function deposit(bytes pubkey, bytes credentials, bytes signature, address operator) payable returns()
func (_DepositContract *DepositContractTransactor) Deposit(opts *bind.TransactOpts, pubkey []byte, credentials []byte, signature []byte, operator common.Address) (*types.Transaction, error) {
	return _DepositContract.contract.Transact(opts, "deposit", pubkey, credentials, signature, operator)
}

// Deposit is a paid mutator transaction binding the contract method 0xe12cf4cb.
//
// Solidity: function deposit(bytes pubkey, bytes credentials, bytes signature, address operator) payable returns()
func (_DepositContract *DepositContractSession) Deposit(pubkey []byte, credentials []byte, signature []byte, operator common.Address) (*types.Transaction, error) {
	return _DepositContract.Contract.Deposit(&_DepositContract.TransactOpts, pubkey, credentials, signature, operator)
}

// Deposit is a paid mutator transaction binding the contract method 0xe12cf4cb.
//
// Solidity: function deposit(bytes pubkey, bytes credentials, bytes signature, address operator) payable returns()
func (_DepositContract *DepositContractTransactorSession) Deposit(pubkey []byte, credentials []byte, signature []byte, operator common.Address) (*types.Transaction, error) {
	return _DepositContract.Contract.Deposit(&_DepositContract.TransactOpts, pubkey, credentials, signature, operator)
}

// RequestOperatorChange is a paid mutator transaction binding the contract method 0xfea7ab77.
//
// Solidity: function requestOperatorChange(bytes pubkey, address newOperator) returns()
func (_DepositContract *DepositContractTransactor) RequestOperatorChange(opts *bind.TransactOpts, pubkey []byte, newOperator common.Address) (*types.Transaction, error) {
	return _DepositContract.contract.Transact(opts, "requestOperatorChange", pubkey, newOperator)
}

// RequestOperatorChange is a paid mutator transaction binding the contract method 0xfea7ab77.
//
// Solidity: function requestOperatorChange(bytes pubkey, address newOperator) returns()
func (_DepositContract *DepositContractSession) RequestOperatorChange(pubkey []byte, newOperator common.Address) (*types.Transaction, error) {
	return _DepositContract.Contract.RequestOperatorChange(&_DepositContract.TransactOpts, pubkey, newOperator)
}

// RequestOperatorChange is a paid mutator transaction binding the contract method 0xfea7ab77.
//
// Solidity: function requestOperatorChange(bytes pubkey, address newOperator) returns()
func (_DepositContract *DepositContractTransactorSession) RequestOperatorChange(pubkey []byte, newOperator common.Address) (*types.Transaction, error) {
	return _DepositContract.Contract.RequestOperatorChange(&_DepositContract.TransactOpts, pubkey, newOperator)
}

// DepositContractDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the DepositContract contract.
type DepositContractDepositIterator struct {
	Event *DepositContractDeposit // Event containing the contract specifics and raw log

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
func (it *DepositContractDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositContractDeposit)
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
		it.Event = new(DepositContractDeposit)
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
func (it *DepositContractDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositContractDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositContractDeposit represents a Deposit event raised by the DepositContract contract.
type DepositContractDeposit struct {
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
func (_DepositContract *DepositContractFilterer) FilterDeposit(opts *bind.FilterOpts) (*DepositContractDepositIterator, error) {

	logs, sub, err := _DepositContract.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &DepositContractDepositIterator{contract: _DepositContract.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x68af751683498a9f9be59fe8b0d52a64dd155255d85cdb29fea30b1e3f891d46.
//
// Solidity: event Deposit(bytes pubkey, bytes credentials, uint64 amount, bytes signature, uint64 index)
func (_DepositContract *DepositContractFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *DepositContractDeposit) (event.Subscription, error) {

	logs, sub, err := _DepositContract.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositContractDeposit)
				if err := _DepositContract.contract.UnpackLog(event, "Deposit", log); err != nil {
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
func (_DepositContract *DepositContractFilterer) ParseDeposit(log types.Log) (*DepositContractDeposit, error) {
	event := new(DepositContractDeposit)
	if err := _DepositContract.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositContractOperatorChangeCancelledIterator is returned from FilterOperatorChangeCancelled and is used to iterate over the raw logs and unpacked data for OperatorChangeCancelled events raised by the DepositContract contract.
type DepositContractOperatorChangeCancelledIterator struct {
	Event *DepositContractOperatorChangeCancelled // Event containing the contract specifics and raw log

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
func (it *DepositContractOperatorChangeCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositContractOperatorChangeCancelled)
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
		it.Event = new(DepositContractOperatorChangeCancelled)
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
func (it *DepositContractOperatorChangeCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositContractOperatorChangeCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositContractOperatorChangeCancelled represents a OperatorChangeCancelled event raised by the DepositContract contract.
type DepositContractOperatorChangeCancelled struct {
	Pubkey common.Hash
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOperatorChangeCancelled is a free log retrieval operation binding the contract event 0x1c0a7e1bd09da292425c039309671a03de56b89a0858598aab6df6ce84b006db.
//
// Solidity: event OperatorChangeCancelled(bytes indexed pubkey)
func (_DepositContract *DepositContractFilterer) FilterOperatorChangeCancelled(opts *bind.FilterOpts, pubkey [][]byte) (*DepositContractOperatorChangeCancelledIterator, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _DepositContract.contract.FilterLogs(opts, "OperatorChangeCancelled", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return &DepositContractOperatorChangeCancelledIterator{contract: _DepositContract.contract, event: "OperatorChangeCancelled", logs: logs, sub: sub}, nil
}

// WatchOperatorChangeCancelled is a free log subscription operation binding the contract event 0x1c0a7e1bd09da292425c039309671a03de56b89a0858598aab6df6ce84b006db.
//
// Solidity: event OperatorChangeCancelled(bytes indexed pubkey)
func (_DepositContract *DepositContractFilterer) WatchOperatorChangeCancelled(opts *bind.WatchOpts, sink chan<- *DepositContractOperatorChangeCancelled, pubkey [][]byte) (event.Subscription, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _DepositContract.contract.WatchLogs(opts, "OperatorChangeCancelled", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositContractOperatorChangeCancelled)
				if err := _DepositContract.contract.UnpackLog(event, "OperatorChangeCancelled", log); err != nil {
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
func (_DepositContract *DepositContractFilterer) ParseOperatorChangeCancelled(log types.Log) (*DepositContractOperatorChangeCancelled, error) {
	event := new(DepositContractOperatorChangeCancelled)
	if err := _DepositContract.contract.UnpackLog(event, "OperatorChangeCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositContractOperatorChangeQueuedIterator is returned from FilterOperatorChangeQueued and is used to iterate over the raw logs and unpacked data for OperatorChangeQueued events raised by the DepositContract contract.
type DepositContractOperatorChangeQueuedIterator struct {
	Event *DepositContractOperatorChangeQueued // Event containing the contract specifics and raw log

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
func (it *DepositContractOperatorChangeQueuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositContractOperatorChangeQueued)
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
		it.Event = new(DepositContractOperatorChangeQueued)
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
func (it *DepositContractOperatorChangeQueuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositContractOperatorChangeQueuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositContractOperatorChangeQueued represents a OperatorChangeQueued event raised by the DepositContract contract.
type DepositContractOperatorChangeQueued struct {
	Pubkey          common.Hash
	QueuedOperator  common.Address
	CurrentOperator common.Address
	QueuedTimestamp *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterOperatorChangeQueued is a free log retrieval operation binding the contract event 0x7640ec3c8c4695deadda414dd20400acf275297a7c38715f9237657e97ddba5f.
//
// Solidity: event OperatorChangeQueued(bytes indexed pubkey, address queuedOperator, address currentOperator, uint256 queuedTimestamp)
func (_DepositContract *DepositContractFilterer) FilterOperatorChangeQueued(opts *bind.FilterOpts, pubkey [][]byte) (*DepositContractOperatorChangeQueuedIterator, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _DepositContract.contract.FilterLogs(opts, "OperatorChangeQueued", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return &DepositContractOperatorChangeQueuedIterator{contract: _DepositContract.contract, event: "OperatorChangeQueued", logs: logs, sub: sub}, nil
}

// WatchOperatorChangeQueued is a free log subscription operation binding the contract event 0x7640ec3c8c4695deadda414dd20400acf275297a7c38715f9237657e97ddba5f.
//
// Solidity: event OperatorChangeQueued(bytes indexed pubkey, address queuedOperator, address currentOperator, uint256 queuedTimestamp)
func (_DepositContract *DepositContractFilterer) WatchOperatorChangeQueued(opts *bind.WatchOpts, sink chan<- *DepositContractOperatorChangeQueued, pubkey [][]byte) (event.Subscription, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _DepositContract.contract.WatchLogs(opts, "OperatorChangeQueued", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositContractOperatorChangeQueued)
				if err := _DepositContract.contract.UnpackLog(event, "OperatorChangeQueued", log); err != nil {
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
func (_DepositContract *DepositContractFilterer) ParseOperatorChangeQueued(log types.Log) (*DepositContractOperatorChangeQueued, error) {
	event := new(DepositContractOperatorChangeQueued)
	if err := _DepositContract.contract.UnpackLog(event, "OperatorChangeQueued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositContractOperatorUpdatedIterator is returned from FilterOperatorUpdated and is used to iterate over the raw logs and unpacked data for OperatorUpdated events raised by the DepositContract contract.
type DepositContractOperatorUpdatedIterator struct {
	Event *DepositContractOperatorUpdated // Event containing the contract specifics and raw log

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
func (it *DepositContractOperatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositContractOperatorUpdated)
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
		it.Event = new(DepositContractOperatorUpdated)
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
func (it *DepositContractOperatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositContractOperatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositContractOperatorUpdated represents a OperatorUpdated event raised by the DepositContract contract.
type DepositContractOperatorUpdated struct {
	Pubkey           common.Hash
	NewOperator      common.Address
	PreviousOperator common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOperatorUpdated is a free log retrieval operation binding the contract event 0x0adffd98d3072c48341843974dffd7a910bb849ba6ca04163d43bb26feb17403.
//
// Solidity: event OperatorUpdated(bytes indexed pubkey, address newOperator, address previousOperator)
func (_DepositContract *DepositContractFilterer) FilterOperatorUpdated(opts *bind.FilterOpts, pubkey [][]byte) (*DepositContractOperatorUpdatedIterator, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _DepositContract.contract.FilterLogs(opts, "OperatorUpdated", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return &DepositContractOperatorUpdatedIterator{contract: _DepositContract.contract, event: "OperatorUpdated", logs: logs, sub: sub}, nil
}

// WatchOperatorUpdated is a free log subscription operation binding the contract event 0x0adffd98d3072c48341843974dffd7a910bb849ba6ca04163d43bb26feb17403.
//
// Solidity: event OperatorUpdated(bytes indexed pubkey, address newOperator, address previousOperator)
func (_DepositContract *DepositContractFilterer) WatchOperatorUpdated(opts *bind.WatchOpts, sink chan<- *DepositContractOperatorUpdated, pubkey [][]byte) (event.Subscription, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _DepositContract.contract.WatchLogs(opts, "OperatorUpdated", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositContractOperatorUpdated)
				if err := _DepositContract.contract.UnpackLog(event, "OperatorUpdated", log); err != nil {
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
func (_DepositContract *DepositContractFilterer) ParseOperatorUpdated(log types.Log) (*DepositContractOperatorUpdated, error) {
	event := new(DepositContractOperatorUpdated)
	if err := _DepositContract.contract.UnpackLog(event, "OperatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
