module github.com/berachain/beacon-kit/mod/sync

go 1.22.4

replace github.com/berachain/beacon-kit/mod/primitives => ../primitives

require (
	github.com/berachain/beacon-kit/mod/log v0.0.0-20240612175710-7d5f3e4f7041
	github.com/berachain/beacon-kit/mod/primitives v0.0.0-20240612175710-7d5f3e4f7041
	github.com/ethereum/go-ethereum v1.14.5
)
