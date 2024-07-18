module github.com/berachain/beacon-kit/mod/runtime

go 1.22.5

replace (
	// The following are required to build with the latest version of the cosmos-sdk main branch:
	cosmossdk.io/api => cosmossdk.io/api v0.7.3-0.20240623110059-dec2d5583e39
	cosmossdk.io/x/staking => cosmossdk.io/x/staking v0.0.0-20240623110059-dec2d5583e39
	github.com/cosmos/cosmos-sdk => github.com/berachain/cosmos-sdk v0.46.0-beta2.0.20240624014538-75ba469b1881
)

require (
	github.com/berachain/beacon-kit/mod/errors v0.0.0-20240617161612-ab1257fcf5a1
	github.com/berachain/beacon-kit/mod/log v0.0.0-20240610210054-bfdc14c4013c
	github.com/berachain/beacon-kit/mod/p2p v0.0.0-20240618214413-d5ec0e66b3dd
	github.com/berachain/beacon-kit/mod/primitives v0.0.0-20240710021622-37e8e3e7e155
	github.com/cometbft/cometbft v1.0.0-rc1.0.20240705183350-8b262bdd1d36
	github.com/cometbft/cometbft/api v1.0.0-rc.1.0.20240702102754-c4fa9a6225b0
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/berachain/beacon-kit/mod/chain-spec v0.0.0-20240703145037-b5612ab256db // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/prysmaticlabs/gohashtree v0.0.4-beta // indirect
)

require (
	github.com/berachain/beacon-kit/mod/async v0.0.0-20240624011057-b0afb8163f14
	github.com/btcsuite/btcd/btcec/v2 v2.3.3 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cosmos/crypto v0.1.1 // indirect
	github.com/cosmos/gogoproto v1.5.0
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/ethereum/go-ethereum v1.14.6 // indirect
	github.com/getsentry/sentry-go v0.28.1 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/holiman/uint256 v1.2.5-0.20240612125212-75a520988c94 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/oasisprotocol/curve25519-voi v0.0.0-20230904125328-1f23a7beb09a // indirect
	github.com/petermattis/goid v0.0.0-20240607163614-bb94eb51e7a7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/supranational/blst v0.3.12 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sync v0.7.0
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
