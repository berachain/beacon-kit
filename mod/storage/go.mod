module github.com/berachain/beacon-kit/mod/storage

go 1.22.5

replace (
	cosmossdk.io/api => cosmossdk.io/api v0.7.3-0.20240731205446-aee9803a0af6
	cosmossdk.io/client/v2 => cosmossdk.io/client/v2 v2.0.0-20240726151427-2e0884564fdb
	cosmossdk.io/collections => github.com/berachain/cosmos-sdk/collections v0.0.0-20240725053043-79fa56d34c79
	cosmossdk.io/core => cosmossdk.io/core v0.12.1-0.20240731205446-aee9803a0af6
	cosmossdk.io/core/testing => cosmossdk.io/core/testing v0.0.0-20240731202123-43dd23137e9d
	cosmossdk.io/log => cosmossdk.io/log v1.3.2-0.20240729192831-b2989459ae91

	cosmossdk.io/store/v2 => cosmossdk.io/store/v2 v2.0.0-20240731205446-aee9803a0af6
	cosmossdk.io/x/auth => cosmossdk.io/x/auth v0.0.0-20240731202123-43dd23137e9d
	cosmossdk.io/x/staking => cosmossdk.io/x/staking v0.0.0-20240731202123-43dd23137e9d

	github.com/cosmos/cosmos-sdk => github.com/berachain/cosmos-sdk v0.46.0-beta2.0.20240729221345-3df18b5c3b34
)

require (
	cosmossdk.io/collections v0.4.0
	cosmossdk.io/core v0.12.1-0.20240623110059-dec2d5583e39
	cosmossdk.io/log v1.3.2-0.20240530141513-465410c75bce
	cosmossdk.io/store/v2 v2.0.0-00010101000000-000000000000
	github.com/berachain/beacon-kit/mod/errors v0.0.0-20240617161612-ab1257fcf5a1
	github.com/berachain/beacon-kit/mod/log v0.0.0-20240610210054-bfdc14c4013c
	github.com/berachain/beacon-kit/mod/primitives v0.0.0-20240726210727-594bfb4e7157
	github.com/cometbft/cometbft v1.0.0-rc1.0.20240729121641-d06d2e8229ee
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/spf13/afero v1.11.0
	github.com/stretchr/testify v1.9.0
	github.com/tidwall/btree v1.7.0
)

require (
	cosmossdk.io/errors/v2 v2.0.0-20240731132947-df72853b3ca5 // indirect
	github.com/berachain/beacon-kit/mod/chain-spec v0.0.0-20240703145037-b5612ab256db // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/prysmaticlabs/gohashtree v0.0.4-beta // indirect
)

require (
	github.com/berachain/beacon-kit/mod/async v0.0.0-20240618214413-d5ec0e66b3dd
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cosmos/gogoproto v1.5.0 // indirect
	github.com/cosmos/ics23/go v0.10.0 // indirect
	github.com/ethereum/go-ethereum v1.14.6 // indirect
	github.com/getsentry/sentry-go v0.28.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-metrics v0.5.3 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/holiman/uint256 v1.3.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240711142825-46eb208f015d // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
