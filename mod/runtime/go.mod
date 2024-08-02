module github.com/berachain/beacon-kit/mod/runtime

go 1.22.5

replace (
	// The following are required to build with the latest version of the cosmos-sdk main branch:
	cosmossdk.io/api => cosmossdk.io/api v0.7.3-0.20240731205446-aee9803a0af6
	cosmossdk.io/x/staking => cosmossdk.io/x/staking v0.0.0-20240731202123-43dd23137e9d
	github.com/berachain/beacon-kit/mod/consensus => ../consensus
	github.com/cosmos/cosmos-sdk => github.com/berachain/cosmos-sdk v0.46.0-beta2.0.20240729221345-3df18b5c3b34
)

require (
	github.com/berachain/beacon-kit/mod/consensus v0.0.0-00010101000000-000000000000
	github.com/berachain/beacon-kit/mod/errors v0.0.0-20240705193247-d464364483df
	github.com/berachain/beacon-kit/mod/log v0.0.0-20240726221339-a8bfeebf8ecf
	github.com/berachain/beacon-kit/mod/p2p v0.0.0-20240618214413-d5ec0e66b3dd
	github.com/berachain/beacon-kit/mod/primitives v0.0.0-20240726210727-594bfb4e7157
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/berachain/beacon-kit/mod/chain-spec v0.0.0-20240705193247-d464364483df // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/prysmaticlabs/gohashtree v0.0.4-beta // indirect
)

require (
	github.com/berachain/beacon-kit/mod/async v0.0.0-20240624011057-b0afb8163f14
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/ethereum/go-ethereum v1.14.6 // indirect
	github.com/getsentry/sentry-go v0.28.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/holiman/uint256 v1.3.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/sync v0.7.0
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
