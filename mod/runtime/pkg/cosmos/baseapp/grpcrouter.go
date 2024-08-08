package baseapp

import (
	"context"

	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/runtime/protoiface"

	"github.com/cosmos/cosmos-sdk/client/grpc/reflection"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type QueryRouter interface {
	HybridHandlerByRequestName(name string) []func(ctx context.Context, req, resp protoiface.MessageV1) error
	RegisterService(sd *grpc.ServiceDesc, handler interface{})
	ResponseNameByRequestName(requestName string) string
	Route(path string) GRPCQueryHandler
	SetInterfaceRegistry(interfaceRegistry codectypes.InterfaceRegistry)
}

// GRPCQueryRouter routes ABCI Query requests to GRPC handlers
type GRPCQueryRouter struct {
	// routes maps query handlers used in ABCIQuery.
	routes map[string]GRPCQueryHandler
	// hybridHandlers maps the request name to the handler. It is a hybrid handler which seamlessly
	// handles both gogo and protov2 messages.
	hybridHandlers map[string][]func(ctx context.Context, req, resp protoiface.MessageV1) error
	// responseByRequestName maps the request name to the response name.
	responseByRequestName map[string]string
	// binaryCodec is used to encode/decode binary protobuf messages.
	binaryCodec codec.BinaryCodec
	// cdc is the gRPC codec used by the router to correctly unmarshal messages.
	cdc encoding.Codec
	// serviceData contains the gRPC services and their handlers.
	serviceData []serviceData
}

// serviceData represents a gRPC service, along with its handler.
type serviceData struct {
	serviceDesc *grpc.ServiceDesc
	handler     interface{}
}

var (
	_ gogogrpc.Server = &GRPCQueryRouter{}
	_ QueryRouter     = &GRPCQueryRouter{}
)

// NewGRPCQueryRouter creates a new GRPCQueryRouter
func NewGRPCQueryRouter() *GRPCQueryRouter {
	return &GRPCQueryRouter{
		routes:                map[string]GRPCQueryHandler{},
		hybridHandlers:        map[string][]func(ctx context.Context, req, resp protoiface.MessageV1) error{},
		responseByRequestName: map[string]string{},
	}
}

// GRPCQueryHandler defines a function type which handles ABCI Query requests
// using gRPC
type GRPCQueryHandler = func(ctx sdk.Context, req *abci.QueryRequest) (*abci.QueryResponse, error)

// Route returns the GRPCQueryHandler for a given query route path or nil
// if not found
func (qrt *GRPCQueryRouter) Route(path string) GRPCQueryHandler {
	handler, found := qrt.routes[path]
	if !found {
		return nil
	}
	return handler
}

// RegisterService implements the gRPC Server.RegisterService method. sd is a gRPC
// service description, handler is an object which implements that gRPC service/
//
// This functions PANICS:
// - if a protobuf service is registered twice.
func (qrt *GRPCQueryRouter) RegisterService(sd *grpc.ServiceDesc, handler interface{}) {}

func (qrt *GRPCQueryRouter) HybridHandlerByRequestName(name string) []func(ctx context.Context, req, resp protoiface.MessageV1) error {
	return qrt.hybridHandlers[name]
}

func (qrt *GRPCQueryRouter) ResponseNameByRequestName(requestName string) string {
	return qrt.responseByRequestName[requestName]
}

// SetInterfaceRegistry sets the interface registry for the router. This will
// also register the interface reflection gRPC service.
func (qrt *GRPCQueryRouter) SetInterfaceRegistry(interfaceRegistry codectypes.InterfaceRegistry) {
	qrt.binaryCodec = codec.NewProtoCodec(interfaceRegistry)
	// instantiate the codec
	qrt.cdc = codec.NewProtoCodec(interfaceRegistry).GRPCCodec()
	// Once we have an interface registry, we can register the interface
	// registry reflection gRPC service.
	reflection.RegisterReflectionServiceServer(qrt, reflection.NewReflectionServiceServer(interfaceRegistry))
}
