module github.com/berachain/beacon-kit/mod/consensus

go 1.22.4

replace (
	// The following are required to build with the latest version of the cosmos-sdk main branch:
	cosmossdk.io/api => cosmossdk.io/api v0.7.3-0.20240623110059-dec2d5583e39
	cosmossdk.io/core => cosmossdk.io/core v0.12.1-0.20240623110059-dec2d5583e39
	cosmossdk.io/core/testing => cosmossdk.io/core/testing v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/server/v2 => cosmossdk.io/server/v2 v2.0.0-20240624195302-5762b0b817f2
	cosmossdk.io/server/v2/appmanager => cosmossdk.io/server/v2/appmanager v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/server/v2/stf => cosmossdk.io/server/v2/stf v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/x/auth => cosmossdk.io/x/auth v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/x/consensus => cosmossdk.io/x/consensus v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/x/staking => cosmossdk.io/x/staking v0.0.0-20240623110059-dec2d5583e39
	github.com/cosmos/cosmos-sdk => github.com/berachain/cosmos-sdk v0.46.0-beta2.0.20240624014538-75ba469b1881
)

require (
	cosmossdk.io/core v0.12.1-0.20240623110059-dec2d5583e39
	github.com/berachain/beacon-kit/mod/async v0.0.0-20240624041151-58d9ba1eb6b4
	github.com/berachain/beacon-kit/mod/errors v0.0.0-20240624041151-58d9ba1eb6b4
	github.com/berachain/beacon-kit/mod/log v0.0.0-20240624041151-58d9ba1eb6b4
	github.com/berachain/beacon-kit/mod/p2p v0.0.0-20240624041151-58d9ba1eb6b4
	github.com/berachain/beacon-kit/mod/primitives v0.0.0-20240624041151-58d9ba1eb6b4
	github.com/berachain/beacon-kit/mod/runtime v0.0.0-20240624041151-58d9ba1eb6b4
	github.com/cometbft/cometbft v1.0.0-alpha.2.0.20240613135100-716d8f8c592d
	github.com/cosmos/gogoproto v1.5.0
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8
	golang.org/x/sync v0.7.0
)

require (
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cometbft/cometbft/api v1.0.0-rc.1 // indirect
	github.com/ethereum/go-ethereum v1.14.5 // indirect
	github.com/ferranbt/fastssz v0.1.4-0.20240422063434-a4db75388da1 // indirect
	github.com/getsentry/sentry-go v0.28.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/holiman/uint256 v1.2.5-0.20240612125212-75a520988c94 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prysmaticlabs/gohashtree v0.0.4-beta // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
