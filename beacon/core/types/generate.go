package types

//go:generate go run github.com/prysmaticlabs/fastssz/sszgen -path . -objs BeaconBlockDeneb,BeaconBlockBodyDeneb,BlobSidecar,Deposit,BlobSidecars -include ../../../primitives,../../../engine/types,$GOPATH/pkg/mod/github.com/ethereum/go-ethereum@$GETH_GO_GENERATE_VERSION/common -output generated.ssz.go
