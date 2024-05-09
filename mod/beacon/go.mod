module github.com/berachain/beacon-kit/mod/beacon

go 1.22.3

replace (
	// The following are required to build with the lastest version of the cosmos-sdk main branch:
	cosmossdk.io/api => cosmossdk.io/api v0.7.3-0.20240421190920-f31a6a3024ba
	cosmossdk.io/x/auth => cosmossdk.io/x/auth v0.0.0-20240506152226-0d32fa94e9eb
	cosmossdk.io/x/consensus => cosmossdk.io/x/consensus v0.0.0-20240506152226-0d32fa94e9eb
	cosmossdk.io/x/staking => cosmossdk.io/x/staking v0.0.0-20240506152226-0d32fa94e9eb

	github.com/berachain/beacon-kit/mod/consensus-types => ../consensus-types
	github.com/berachain/beacon-kit/mod/errors => ../errors
	github.com/berachain/beacon-kit/mod/execution => ../execution
	github.com/berachain/beacon-kit/mod/log => ../log
	github.com/berachain/beacon-kit/mod/primitives => ../primitives
	github.com/berachain/beacon-kit/mod/primitives-engine => ../primitives-engine
	github.com/berachain/beacon-kit/mod/state-transition => ../state-transition
	github.com/berachain/beacon-kit/mod/storage => ../storage
	github.com/cosmos/cosmos-sdk => github.com/berachain/cosmos-sdk v0.46.0-beta2.0.20240506152226-0d32fa94e9eb
)

require (
	github.com/berachain/beacon-kit/mod/consensus-types v0.0.0-20240508035017-2fb637ea5f0a
	github.com/berachain/beacon-kit/mod/errors v0.0.0-20240508035017-2fb637ea5f0a
	github.com/berachain/beacon-kit/mod/execution v0.0.0-00010101000000-000000000000
	github.com/berachain/beacon-kit/mod/log v0.0.0-20240508035017-2fb637ea5f0a
	github.com/berachain/beacon-kit/mod/primitives v0.0.0-20240508035017-2fb637ea5f0a
	github.com/berachain/beacon-kit/mod/primitives-engine v0.0.0-20240508035017-2fb637ea5f0a
	github.com/berachain/beacon-kit/mod/state-transition v0.0.0-00010101000000-000000000000
	github.com/berachain/beacon-kit/mod/storage v0.0.0-00010101000000-000000000000
	github.com/ethereum/go-ethereum v1.14.3
	github.com/ferranbt/fastssz v0.1.4-0.20240422063434-a4db75388da1
	golang.org/x/sync v0.7.0
)

require (
	cosmossdk.io/collections v0.4.0 // indirect
	cosmossdk.io/core v0.12.1-0.20240421190920-f31a6a3024ba // indirect
	cosmossdk.io/log v1.3.1 // indirect
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.13.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cockroachdb/errors v1.11.1 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v1.1.0 // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.12.1 // indirect
	github.com/cosmos/gogoproto v1.4.12 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20231025140028-3c0104f4b233 // indirect
	github.com/crate-crypto/go-kzg-4844 v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/ethereum/c-kzg-4844 v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gballet/go-verkle v0.1.1-0.20231031103413-a67434b50f46 // indirect
	github.com/getsentry/sentry-go v0.27.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/holiman/uint256 v1.2.4 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/linxGnu/grocksdb v1.8.14 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.19.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.53.0 // indirect
	github.com/prometheus/procfs v0.14.0 // indirect
	github.com/prysmaticlabs/gohashtree v0.0.4-beta // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/rs/zerolog v1.32.0 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/supranational/blst v0.3.11 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240415180920-8c6c420018be // indirect
	google.golang.org/grpc v1.63.2 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)
