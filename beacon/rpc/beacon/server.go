package beacon

import (
	"github.com/berachain/beacon-kit/beacon/rpc"
	"net/http"
	"strconv"

	"github.com/berachain/beacon-kit/runtime/service"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gorilla/mux"
)

// Server defines a server implementation of the gRPC Beacon Chain service,
// providing RPC endpoints to access data relevant to the Ethereum Beacon Chain.
type Server struct {
	ContextGetter func(height int64, prove bool) (sdk.Context, error)
	Service       service.BeaconStorageBackend
}

func (s *Server) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/eth/v1/beacon/states/{state_id}/randao", s.GetRandao).Methods(http.MethodGet)
}

// GetRandao fetches the RANDAO mix for the requested epoch from the state identified by state_id.
// If an epoch is not specified then the RANDAO mix for the state's current epoch will be returned.
// By adjusting the state_id parameter you can query for any historic value of the RANDAO mix.
// Ordinarily states from the same epoch will mutate the RANDAO mix for that epoch as blocks are applied.
func (s *Server) GetRandao(w http.ResponseWriter, r *http.Request) {
	stateId := mux.Vars(r)["state_id"]
	if stateId == "" {
		rpc.HandleError(w, "state_id is required in URL params", http.StatusBadRequest)
		return
	}

	stateIdAsInt, err := strconv.ParseUint(stateId, 10, 64)

	ctx, err := s.ContextGetter(int64(stateIdAsInt), false)
	if err != nil {
		rpc.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	randao, err := s.Service.BeaconState(ctx).RandaoMixAtIndex(stateIdAsInt)
	if err != nil {
		rpc.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &GetRandaoResponse{
		Data:                &Randao{Randao: hexutil.Encode(randao[:])},
		ExecutionOptimistic: true,
		Finalized:           true,
	}

	rpc.WriteJson(w, resp)
}
