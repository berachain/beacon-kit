package ethclient

import "github.com/ethereum/go-ethereum/ethclient"

type (
	Client = ethclient.Client
)

//nolint:gochecknoglobals // its okay.
var (
	NewClient = ethclient.NewClient
)
