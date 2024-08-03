module github.com/berachain/beacon-kit/mod/cli

go 1.22.5

replace (
	// The following are required to build with the latest version of the cosmos-sdk main branch:
	cosmossdk.io/api => cosmossdk.io/api v0.7.3-0.20240623110059-dec2d5583e39
	cosmossdk.io/core/testing => cosmossdk.io/core/testing v0.0.0-20240623110059-dec2d5583e39
	github.com/cosmos/cosmos-sdk => github.com/berachain/cosmos-sdk v0.46.0-beta2.0.20240624014538-75ba469b1881
)

require (
	cosmossdk.io/client/v2 v2.0.0-20240412212305-037cf98f7eea
	cosmossdk.io/core v0.12.1-0.20240623110059-dec2d5583e39
	cosmossdk.io/depinject v1.0.0-alpha.4.0.20240506202947-fbddf0a55044
	cosmossdk.io/log v1.3.2-0.20240530141513-465410c75bce
	cosmossdk.io/tools/confix v0.1.1
	github.com/berachain/beacon-kit/mod/config v0.0.0-20240705193247-d464364483df
	github.com/berachain/beacon-kit/mod/consensus-types v0.0.0-20240724161918-96b9bf999de2
	github.com/berachain/beacon-kit/mod/engine-primitives v0.0.0-20240710022615-726645827bad
	github.com/berachain/beacon-kit/mod/errors v0.0.0-20240705193247-d464364483df
	github.com/berachain/beacon-kit/mod/geth-primitives v0.0.0-20240722101705-632d23f22727
	github.com/berachain/beacon-kit/mod/log v0.0.0-20240705193247-d464364483df
	github.com/berachain/beacon-kit/mod/node-core v0.0.0-20240710023144-cd73de9f7378
	github.com/berachain/beacon-kit/mod/primitives v0.0.0-20240726210727-594bfb4e7157
	github.com/cometbft/cometbft v1.0.0-rc1.0.20240729121641-d06d2e8229ee
	github.com/cosmos/cosmos-sdk v0.51.0
	github.com/ferranbt/fastssz v0.1.4-0.20240629094022-eac385e6ee79
	github.com/spf13/afero v1.11.0
	github.com/spf13/cobra v1.8.1
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.9.0
	golang.org/x/sync v0.7.0
)

require (
	github.com/berachain/beacon-kit/mod/chain-spec v0.0.0-20240705193247-d464364483df // indirect
	github.com/cockroachdb/fifo v0.0.0-20240616162244-4768e80dfb9a // indirect
	github.com/ethereum/go-ethereum v1.14.6 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/karalabe/ssz v0.2.1-0.20240724074312-3d1ff7a6f7c4 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/phuslu/log v1.0.108-0.20240705160716-a8f8c12ae6c6 // indirect
	github.com/prysmaticlabs/go-bitfield v0.0.0-20240618144021-706c95b2dd15 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
)

require (
	buf.build/gen/go/cometbft/cometbft/protocolbuffers/go v1.34.2-20240312114316-c0d3497e35d6.2 // indirect
	buf.build/gen/go/cosmos/gogo-proto/protocolbuffers/go v1.34.2-20240130113600-88ef6483f90f.2 // indirect
	cosmossdk.io/api v0.7.5 // indirect
	cosmossdk.io/collections v0.4.0 // indirect
	cosmossdk.io/errors v1.0.1 // indirect
	cosmossdk.io/math v1.3.0 // indirect
	cosmossdk.io/store v1.1.1-0.20240418092142-896cdf1971bc // indirect
	cosmossdk.io/store/v2 v2.0.0-20240515130459-16437119e0d8 // indirect
	cosmossdk.io/x/accounts v0.0.0-20240623110059-dec2d5583e39 // indirect
	cosmossdk.io/x/auth v0.0.0-20240623110059-dec2d5583e39 // indirect
	cosmossdk.io/x/bank v0.0.0-20240623110059-dec2d5583e39 // indirect
	cosmossdk.io/x/consensus v0.0.0-20240623110059-dec2d5583e39 // indirect
	cosmossdk.io/x/gov v0.0.0-20240623110059-dec2d5583e39 // indirect
	cosmossdk.io/x/staking v0.0.0-20240623110059-dec2d5583e39 // indirect
	cosmossdk.io/x/tx v0.13.4-0.20240623110059-dec2d5583e39 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/DataDog/datadog-go v4.8.3+incompatible // indirect
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/VictoriaMetrics/fastcache v1.12.2 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/berachain/beacon-kit/mod/async v0.0.0-20240705193247-d464364483df // indirect
	// indirect
	github.com/berachain/beacon-kit/mod/beacon v0.0.0-20240705193247-d464364483df // indirect
	github.com/berachain/beacon-kit/mod/da v0.0.0-20240705193247-d464364483df // indirect
	github.com/berachain/beacon-kit/mod/execution v0.0.0-20240705193247-d464364483df
	github.com/berachain/beacon-kit/mod/p2p v0.0.0-20240618214413-d5ec0e66b3dd // indirect
	github.com/berachain/beacon-kit/mod/payload v0.0.0-20240705193247-d464364483df // indirect
	github.com/berachain/beacon-kit/mod/runtime v0.0.0-20240710022902-57d64ab384c7 // indirect
	github.com/berachain/beacon-kit/mod/state-transition v0.0.0-20240710023256-e9246f6a35e1 // indirect
	github.com/berachain/beacon-kit/mod/storage v0.0.0-20240705193247-d464364483df // indirect
	github.com/bgentry/speakeasy v0.1.1-0.20220910012023-760eaf8b6816 // indirect
	github.com/bits-and-blooms/bitset v1.13.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v1.1.1 // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/cometbft/cometbft-db v0.12.0 // indirect
	github.com/cometbft/cometbft/api v1.0.0-rc.1.0.20240711183925-948692fddcbe // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.12.1 // indirect
	github.com/cosmos/btcutil v1.0.5 // indirect
	github.com/cosmos/cosmos-db v1.0.2 // indirect
	github.com/cosmos/cosmos-proto v1.0.0-beta.5 // indirect
	github.com/cosmos/crypto v0.1.2 // indirect
	github.com/cosmos/go-bip39 v1.0.0 // indirect
	github.com/cosmos/gogogateway v1.2.0 // indirect
	github.com/cosmos/gogoproto v1.5.0 // indirect
	github.com/cosmos/iavl v1.2.0 // indirect
	github.com/cosmos/ics23/go v0.10.0 // indirect
	github.com/cosmos/ledger-cosmos-go v0.13.3 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20240223125850-b1e8a79f509c // indirect
	github.com/crate-crypto/go-kzg-4844 v1.0.0 // indirect
	github.com/creachadair/atomicfile v0.3.1 // indirect
	github.com/creachadair/tomledit v0.0.24 // indirect
	github.com/danieljoos/wincred v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/dgraph-io/badger/v4 v4.2.0 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.7.0 // indirect
	github.com/emicklei/dot v1.6.2 // indirect
	github.com/ethereum/c-kzg-4844 v1.0.2 // indirect
	github.com/ethereum/go-verkle v0.1.1-0.20240306133620-7d920df305f0 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/getsentry/sentry-go v0.28.1 // indirect
	github.com/go-faster/xor v1.0.0 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gofrs/flock v0.12.1 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/glog v1.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/flatbuffers v24.3.25+incompatible // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/orderedcode v0.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/handlers v1.5.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-metrics v0.5.3 // indirect
	github.com/hashicorp/go-plugin v1.6.1 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/hdevalence/ed25519consensus v0.2.0 // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/holiman/uint256 v1.3.0 // indirect
	github.com/huandu/skiplist v1.2.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/itsdevbear/comet-bls12-381 v0.0.0-20240413212931-2ae2f204cde7 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/linxGnu/grocksdb v1.9.2 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/manifoldco/promptui v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mdp/qrterminal/v3 v3.2.0 // indirect
	github.com/minio/highwayhash v1.0.3 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/oasisprotocol/curve25519-voi v0.0.0-20230904125328-1f23a7beb09a // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/petermattis/goid v0.0.0-20240607163614-bb94eb51e7a7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/prysmaticlabs/gohashtree v0.0.4-beta // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/rs/cors v1.11.0 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/supranational/blst v0.3.12 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tendermint/go-amino v0.16.0 // indirect
	github.com/tidwall/btree v1.7.0 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zondax/hid v0.9.2 // indirect
	github.com/zondax/ledger-go v0.14.3 // indirect
	gitlab.com/yawning/secp256k1-voi v0.0.0-20230925100816-f2616030848b // indirect
	gitlab.com/yawning/tuplehash v0.0.0-20230713102510-df83abbf9a02 // indirect
	go.etcd.io/bbolt v1.4.0-alpha.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/term v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240711142825-46eb208f015d // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.5.1 // indirect
	pgregory.net/rapid v1.1.0 // indirect
	rsc.io/qr v0.2.0 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
