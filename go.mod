module github.com/berachain/beacon-kit

go 1.22.1

replace (
	// The following are required to build with the lastest version of the cosmos-sdk main branch:
	cosmossdk.io/api => github.com/berachain/cosmos-sdk/api v0.4.2-0.20240325184954-fae9fbd80a0a
	cosmossdk.io/client/v2 => cosmossdk.io/client/v2 v2.0.0-20240325170708-c46688729fd9
	cosmossdk.io/core => cosmossdk.io/core v0.12.1-0.20240325170708-c46688729fd9
	cosmossdk.io/x/accounts => cosmossdk.io/x/accounts v0.0.0-20240325170708-c46688729fd9
	cosmossdk.io/x/auth => cosmossdk.io/x/auth v0.0.0-20240325170708-c46688729fd9
	cosmossdk.io/x/bank => cosmossdk.io/x/bank v0.0.0-20240325170708-c46688729fd9
	cosmossdk.io/x/gov => cosmossdk.io/x/gov v0.0.0-20240325170708-c46688729fd9
	cosmossdk.io/x/staking => cosmossdk.io/x/staking v0.0.0-20240325170708-c46688729fd9
	github.com/cometbft/cometbft => github.com/berachain/cometbft v0.0.0-20240312055307-dff5fd68a3b0
	github.com/cosmos/cosmos-sdk => github.com/berachain/cosmos-sdk v0.46.0-beta2.0.20240325184954-fae9fbd80a0a

)

require (
	cosmossdk.io/api v0.7.3
	cosmossdk.io/client/v2 v2.0.0-20240221095859-541df89f2bb4
	cosmossdk.io/collections v0.4.0
	cosmossdk.io/core v0.12.1-0.20231114100755-569e3ff6a0d7
	cosmossdk.io/depinject v1.0.0-alpha.4.0.20240221095859-541df89f2bb4
	cosmossdk.io/log v1.3.1
	cosmossdk.io/store v1.1.0
	cosmossdk.io/tools/confix v0.1.1
	cosmossdk.io/x/auth v0.0.0-00010101000000-000000000000
	github.com/bazelbuild/buildtools v0.0.0-20240207142252-03bf520394af
	github.com/bufbuild/buf v1.30.0
	github.com/cockroachdb/errors v1.11.1
	github.com/cometbft/cometbft v0.38.6
	github.com/cosmos/cosmos-db v1.0.2
	github.com/cosmos/cosmos-proto v1.0.0-beta.4
	github.com/cosmos/cosmos-sdk v0.51.0
	github.com/cosmos/gosec/v2 v2.0.0-20230124142343-bf28a33fadf2
	github.com/ethereum/go-ethereum v1.13.5-0.20240328163540-a3829178af6c
	github.com/ferranbt/fastssz v0.1.4-0.20240325182853-3fad96355b01
	github.com/fjl/gencodec v0.0.0-20230517082657-f9840df7b83e
	github.com/golangci/golangci-lint v1.57.2
	github.com/google/addlicense v1.1.1
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/holiman/uint256 v1.2.4
	github.com/itsdevbear/comet-bls12-381 v0.0.0-20240226135442-10e4707bd0ca
	github.com/kurtosis-tech/kurtosis/api/golang v0.88.11
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/minio/sha256-simd v1.0.1
	github.com/protolambda/ztyp v0.2.2
	github.com/prysmaticlabs/gohashtree v0.0.4-beta
	github.com/segmentio/golines v0.12.2
	github.com/sourcegraph/conc v0.3.0
	github.com/spf13/afero v1.11.0
	github.com/spf13/cast v1.6.0
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
	github.com/stretchr/testify v1.9.0
	github.com/vektra/mockery/v2 v2.42.1
	go.uber.org/automaxprocs v1.5.3
	go.uber.org/nilaway v0.0.0-20240224031343-67945fb5199f
	golang.org/x/sync v0.6.0
	google.golang.org/protobuf v1.33.0
)

require github.com/go-faster/xor v1.0.0 // indirect
