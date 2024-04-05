package provider

// import (
// 	"strconv"

// 	abci "github.com/cometbft/cometbft/abci/types"

// 	legacyerrors "github.com/cosmos/cosmos-sdk/types/errors"
// 	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"

// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/metadata"
// 	"google.golang.org/grpc/status"
// )

// func sdkErrorToGRPCError(resp abci.ResponseQuery) error {
// 	switch resp.Code {
// 	case legacyerrors.ErrInvalidRequest.ABCICode():
// 		return status.Error(codes.InvalidArgument, resp.Log)
// 	case legacyerrors.ErrUnauthorized.ABCICode():
// 		return status.Error(codes.Unauthenticated, resp.Log)
// 	case legacyerrors.ErrKeyNotFound.ABCICode():
// 		return status.Error(codes.NotFound, resp.Log)
// 	default:
// 		return status.Error(codes.Unknown, resp.Log)
// 	}
// }

// func GetHeightFromMetadata(md metadata.MD) (int64, error) {
// 	height := md.Get(grpctypes.GRPCBlockHeightHeader)
// 	if len(height) == 1 {
// 		return strconv.ParseInt(height[0], 10, 64)
// 	}
// 	return 0, nil
// }
