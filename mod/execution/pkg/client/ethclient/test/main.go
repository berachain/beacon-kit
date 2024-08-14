// package main

// import (
// 	"context"
// 	"fmt"

// 	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
// 	"github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient"
// 	"github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient/rpc"
// 	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
// 	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
// 	"github.com/davecgh/go-spew/spew"
// )

// func main() {
// 	c := rpc.NewClient("https://restless-thrumming-county.bera-bartio.quiknode.pro/6f5c8dc2120be6048421ac6d84c1f700e5875e50")
// 	cl := ethclient.New[*types.ExecutionPayload](c)

// 	logs, err := cl.GetLogsAtBlockNumber(
// 		context.TODO(),
// 		math.U64(2870518),
// 		common.NewExecutionAddressFromHex("0x4242424242424242424242424242424242424242"),
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	for _, l := range logs {
// 		spew.Dump(l)
// 		deposit := new(types.Deposit)
// 		deposit.UnmarshalLog(l)
// 		spew.Dump(deposit)
// 	}
// }
