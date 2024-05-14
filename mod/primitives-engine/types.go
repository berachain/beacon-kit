// package engineprimitives

// import (
// 	"math/big"

// 	"github.com/berachain/beacon-kit/mod/primitives/pkg/common/types"
// 	"github.com/ethereum/go-ethereum/common"
// 	ethTypes "github.com/ethereum/go-ethereum/core/types"
// )

// // ExecutableData is the data necessary to execute an EL payload.
// type ExecutableData struct {
// 	ParentHash    types.Hash             `json:"parentHash"    gencodec:"required"`
// 	FeeRecipient  common.Address         `json:"feeRecipient"  gencodec:"required"`
// 	StateRoot     types.Hash             `json:"stateRoot"     gencodec:"required"`
// 	ReceiptsRoot  types.Hash             `json:"receiptsRoot"  gencodec:"required"`
// 	LogsBloom     []byte                 `json:"logsBloom"     gencodec:"required"`
// 	Random        common.Hash            `json:"prevRandao"    gencodec:"required"`
// 	Number        uint64                 `json:"blockNumber"   gencodec:"required"`
// 	GasLimit      uint64                 `json:"gasLimit"      gencodec:"required"`
// 	GasUsed       uint64                 `json:"gasUsed"       gencodec:"required"`
// 	Timestamp     uint64                 `json:"timestamp"     gencodec:"required"`
// 	ExtraData     []byte                 `json:"extraData"     gencodec:"required"`
// 	BaseFeePerGas *big.Int               `json:"baseFeePerGas" gencodec:"required"`
// 	BlockHash     types.Hash             `json:"blockHash"     gencodec:"required"`
// 	Transactions  [][]byte               `json:"transactions"  gencodec:"required"`
// 	Withdrawals   []*ethTypes.Withdrawal `json:"withdrawals"`
// 	BlobGasUsed   *uint64                `json:"blobGasUsed"`
// 	ExcessBlobGas *uint64                `json:"excessBlobGas"`
// }

// type ForkchoiceState struct {
// 	HeadBlockHash      types.Hash `json:"headBlockHash"`
// 	SafeBlockHash      types.Hash `json:"safeBlockHash"`
// 	FinalizedBlockHash types.Hash `json:"finalizedBlockHash"`
// }

// type PayloadStatusV1 struct {
// 	Status          string      `json:"status"`
// 	LatestValidHash *types.Hash `json:"latestValidHash"`
// 	ValidationError *string     `json:"validationError"`
// }
