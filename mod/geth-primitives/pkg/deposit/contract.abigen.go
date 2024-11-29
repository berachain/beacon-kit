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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOperatorChange\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowDeposit\",\"inputs\":[{\"name\":\"depositor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"number\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelOperatorChange\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelOwnershipHandover\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"completeOwnershipHandover\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"withdrawal_credentials\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositAuth\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depositCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"genesisDepositsRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperator\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"result\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownershipHandoverExpiresAt\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"queuedOperator\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"queuedTimestamp\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"newOperator\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"requestOperatorChange\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"newOperator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestOwnershipHandover\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"credentials\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"signature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"index\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorChangeCancelled\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorChangeQueued\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"queuedOperator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"currentOperator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"queuedTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorUpdated\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"newOperator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"previousOperator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipHandoverCanceled\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipHandoverRequested\",\"inputs\":[{\"name\":\"pendingOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"oldOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DepositNotMultipleOfGwei\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DepositValueTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientDeposit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCredentialsLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPubKeyLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignatureLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NewOwnerIsZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoHandoverRequest\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotEnoughTime\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotNewOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OperatorAlreadySet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedDeposit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroOperatorOnFirstDeposit\",\"inputs\":[]}]",
	Bin: "0x6080604052348015600e575f80fd5b50604051611690380380611690833981016040819052602b91608e565b6032816037565b5060b9565b638b78c6d819805415605057630dc149f05f526004601cfd5b6001600160a01b03909116801560ff1b8117909155805f7f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08180a350565b5f60208284031215609d575f80fd5b81516001600160a01b038116811460b2575f80fd5b9392505050565b6115ca806100c65f395ff3fe608060405260043610610123575f3560e01c8063715018a6116100a1578063e12cf4cb11610071578063f2fde38b11610057578063f2fde38b146103a1578063fea7ab77146103b4578063fee81cf4146103d3575f80fd5b8063e12cf4cb1461037b578063f04e283e1461038e575f80fd5b8063715018a6146102e15780638da5cb5b146102e95780639eaffa961461033d578063c53925d91461035c575f80fd5b80633523f9bd116100f6578063560036ec116100dc578063560036ec146101fd578063577212fe146102a35780635a7517ad146102c2575f80fd5b80633523f9bd146101d257806354d1f13d146101f5575f80fd5b806301ffc9a714610127578063256929621461015b5780632dfdf0b5146101655780633198a6b81461019d575b5f80fd5b348015610132575f80fd5b5061014661014136600461105f565b610404565b60405190151581526020015b60405180910390f35b61016361049c565b005b348015610170575f80fd5b505f546101849067ffffffffffffffff1681565b60405167ffffffffffffffff9091168152602001610152565b3480156101a8575f80fd5b506101846101b73660046110c8565b60046020525f908152604090205467ffffffffffffffff1681565b3480156101dd575f80fd5b506101e760015481565b604051908152602001610152565b6101636104e9565b348015610208575f80fd5b5061026a61021736600461110e565b80516020818301810180516003825292820191909301209152546bffffffffffffffffffffffff8116906c01000000000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1682565b604080516bffffffffffffffffffffffff909316835273ffffffffffffffffffffffffffffffffffffffff909116602083015201610152565b3480156102ae575f80fd5b506101636102bd366004611243565b610522565b3480156102cd575f80fd5b506101636102dc366004611282565b6105f8565b61016361065e565b3480156102f4575f80fd5b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927545b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610152565b348015610348575f80fd5b50610318610357366004611243565b610671565b348015610367575f80fd5b50610163610376366004611243565b6106b2565b6101636103893660046112c3565b6108d9565b61016361039c3660046110c8565b61098e565b6101636103af3660046110c8565b6109cb565b3480156103bf575f80fd5b506101636103ce366004611372565b6109f1565b3480156103de575f80fd5b506101e76103ed3660046110c8565b63389a75e1600c9081525f91909152602090205490565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000148061049657507fffffffff0000000000000000000000000000000000000000000000000000000082167f136f920d00000000000000000000000000000000000000000000000000000000145b92915050565b5f6202a30067ffffffffffffffff164201905063389a75e1600c52335f52806020600c2055337fdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d5f80a250565b63389a75e1600c52335f525f6020600c2055337ffa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c925f80a2565b600282826040516105349291906113c2565b908152604051908190036020019020543373ffffffffffffffffffffffffffffffffffffffff90911614610594576040517f7c214f0400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600382826040516105a69291906113c2565b9081526040519081900360200181205f90556105c590839083906113c2565b604051908190038120907f1c0a7e1bd09da292425c039309671a03de56b89a0858598aab6df6ce84b006db905f90a25050565b610600610ba3565b73ffffffffffffffffffffffffffffffffffffffff919091165f90815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff909216919091179055565b610666610ba3565b61066f5f610bd8565b565b5f600283836040516106849291906113c2565b9081526040519081900360200190205473ffffffffffffffffffffffffffffffffffffffff16905092915050565b5f600383836040516106c59291906113c2565b908152604051908190036020019020805490915073ffffffffffffffffffffffffffffffffffffffff6c01000000000000000000000000820416906bffffffffffffffffffffffff16338214610747576040517f819a0d0b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6bffffffffffffffffffffffff421661076362015180836113fe565b6bffffffffffffffffffffffff1611156107a9576040517fe8966d7a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f600286866040516107bc9291906113c2565b9081526040519081900360200181205473ffffffffffffffffffffffffffffffffffffffff16915083906002906107f690899089906113c2565b908152604051908190036020018120805473ffffffffffffffffffffffffffffffffffffffff939093167fffffffffffffffffffffffff00000000000000000000000000000000000000009093169290921790915560039061085b90889088906113c2565b9081526040519081900360200181205f905561087a90879087906113c2565b6040805191829003822073ffffffffffffffffffffffffffffffffffffffff808716845284166020840152917f0adffd98d3072c48341843974dffd7a910bb849ba6ca04163d43bb26feb17403910160405180910390a2505050505050565b335f9081526004602052604081205467ffffffffffffffff16900361092a576040517fce7ccd9600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b335f90815260046020526040812080549091906109509067ffffffffffffffff16611422565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555061098587878787878787610c46565b50505050505050565b610996610ba3565b63389a75e1600c52805f526020600c2080544211156109bc57636f5e88185f526004601cfd5b5f90556109c881610bd8565b50565b6109d3610ba3565b8060601b6109e857637448fbae5f526004601cfd5b6109c881610bd8565b5f60028484604051610a049291906113c2565b9081526040519081900360200190205473ffffffffffffffffffffffffffffffffffffffff169050338114610a65576040517f7c214f0400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8216610ab2576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60038585604051610ac59291906113c2565b908152604051908190036020018120426bffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff86166c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffffff000000000000000000000000161781559150610b3d90869086906113c2565b6040805191829003822073ffffffffffffffffffffffffffffffffffffffff8681168452851660208401524283830152905190917f7640ec3c8c4695deadda414dd20400acf275297a7c38715f9237657e97ddba5f919081900360600190a25050505050565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffff7487392754331461066f576382b429005f526004601cfd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927805473ffffffffffffffffffffffffffffffffffffffff9092169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a3811560ff1b8217905550565b60308614610c80576040517f9f10647200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208414610cba576040517fb39bca1600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60608214610cf4576040517f4be6321b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff1660028888604051610d1d9291906113c2565b9081526040519081900360200190205473ffffffffffffffffffffffffffffffffffffffff1603610e645773ffffffffffffffffffffffffffffffffffffffff8116610d95576040517f51969a7a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060028888604051610da89291906113c2565b908152604051908190036020018120805473ffffffffffffffffffffffffffffffffffffffff939093167fffffffffffffffffffffffff000000000000000000000000000000000000000090931692909217909155610e0a90889088906113c2565b6040805191829003822073ffffffffffffffffffffffffffffffffffffffff841683525f6020840152917f0adffd98d3072c48341843974dffd7a910bb849ba6ca04163d43bb26feb17403910160405180910390a2610eb2565b73ffffffffffffffffffffffffffffffffffffffff811615610eb2576040517fc4142b4100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610ebb610f9a565b905064077359400067ffffffffffffffff82161015610f06576040517f0e1eddda00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f80547f68af751683498a9f9be59fe8b0d52a64dd155255d85cdb29fea30b1e3f891d46918a918a918a918a9187918b918b9167ffffffffffffffff169080610f4e83611463565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550604051610f889897969594939291906114d6565b60405180910390a15050505050505050565b5f610fa9633b9aca003461156e565b15610fe0576040517f40567b3800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610fef633b9aca0034611581565b905067ffffffffffffffff811115611033576040517f2aa6673400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61103d5f34611042565b919050565b5f385f3884865af161105b5763b12d13eb5f526004601cfd5b5050565b5f6020828403121561106f575f80fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461109e575f80fd5b9392505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461103d575f80fd5b5f602082840312156110d8575f80fd5b61109e826110a5565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f6020828403121561111e575f80fd5b813567ffffffffffffffff811115611134575f80fd5b8201601f81018413611144575f80fd5b803567ffffffffffffffff81111561115e5761115e6110e1565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff821117156111ca576111ca6110e1565b6040528181528282016020018610156111e1575f80fd5b816020840160208301375f91810160200191909152949350505050565b5f8083601f84011261120e575f80fd5b50813567ffffffffffffffff811115611225575f80fd5b60208301915083602082850101111561123c575f80fd5b9250929050565b5f8060208385031215611254575f80fd5b823567ffffffffffffffff81111561126a575f80fd5b611276858286016111fe565b90969095509350505050565b5f8060408385031215611293575f80fd5b61129c836110a5565b9150602083013567ffffffffffffffff811681146112b8575f80fd5b809150509250929050565b5f805f805f805f6080888a0312156112d9575f80fd5b873567ffffffffffffffff8111156112ef575f80fd5b6112fb8a828b016111fe565b909850965050602088013567ffffffffffffffff81111561131a575f80fd5b6113268a828b016111fe565b909650945050604088013567ffffffffffffffff811115611345575f80fd5b6113518a828b016111fe565b90945092506113649050606089016110a5565b905092959891949750929550565b5f805f60408486031215611384575f80fd5b833567ffffffffffffffff81111561139a575f80fd5b6113a6868287016111fe565b90945092506113b99050602085016110a5565b90509250925092565b818382375f9101908152919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b6bffffffffffffffffffffffff8181168382160190811115610496576104966113d1565b5f67ffffffffffffffff82168061143b5761143b6113d1565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0192915050565b5f67ffffffffffffffff821667ffffffffffffffff8103611486576114866113d1565b60010192915050565b81835281816020850137505f602082840101525f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60a081525f6114e960a083018a8c61148f565b82810360208401526114fc81898b61148f565b905067ffffffffffffffff87166040840152828103606084015261152181868861148f565b91505067ffffffffffffffff831660808301529998505050505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f8261157c5761157c611541565b500690565b5f8261158f5761158f611541565b50049056fea26469706673582212203d2af01eee1bb2e5809047b6601bca7c0f96684b662b8fc3e450d1487af19c0a64736f6c634300081a0033",
}

// BeaconDepositContractABI is the input ABI used to generate the binding from.
// Deprecated: Use BeaconDepositContractMetaData.ABI instead.
var BeaconDepositContractABI = BeaconDepositContractMetaData.ABI

// BeaconDepositContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BeaconDepositContractMetaData.Bin instead.
var BeaconDepositContractBin = BeaconDepositContractMetaData.Bin

// DeployBeaconDepositContract deploys a new Ethereum contract, binding an instance of BeaconDepositContract to it.
func DeployBeaconDepositContract(auth *bind.TransactOpts, backend bind.ContractBackend, owner common.Address) (common.Address, *types.Transaction, *BeaconDepositContract, error) {
	parsed, err := BeaconDepositContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BeaconDepositContractBin), backend, owner)
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

// Deposit is a paid mutator transaction binding the contract method 0xe12cf4cb.
//
// Solidity: function deposit(bytes pubkey, bytes withdrawal_credentials, bytes signature, address operator) payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactor) Deposit(opts *bind.TransactOpts, pubkey []byte, withdrawal_credentials []byte, signature []byte, operator common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.contract.Transact(opts, "deposit", pubkey, withdrawal_credentials, signature, operator)
}

// Deposit is a paid mutator transaction binding the contract method 0xe12cf4cb.
//
// Solidity: function deposit(bytes pubkey, bytes withdrawal_credentials, bytes signature, address operator) payable returns()
func (_BeaconDepositContract *BeaconDepositContractSession) Deposit(pubkey []byte, withdrawal_credentials []byte, signature []byte, operator common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.Deposit(&_BeaconDepositContract.TransactOpts, pubkey, withdrawal_credentials, signature, operator)
}

// Deposit is a paid mutator transaction binding the contract method 0xe12cf4cb.
//
// Solidity: function deposit(bytes pubkey, bytes withdrawal_credentials, bytes signature, address operator) payable returns()
func (_BeaconDepositContract *BeaconDepositContractTransactorSession) Deposit(pubkey []byte, withdrawal_credentials []byte, signature []byte, operator common.Address) (*types.Transaction, error) {
	return _BeaconDepositContract.Contract.Deposit(&_BeaconDepositContract.TransactOpts, pubkey, withdrawal_credentials, signature, operator)
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
