module github.com/berachain/beacon-kit

go 1.23.6

replace (
	github.com/cometbft/cometbft => github.com/berachain/cometbft v1.0.1-0.20250416103256-c547460cd751
	github.com/cometbft/cometbft/api => github.com/berachain/cometbft/api v1.0.1-0.20250416103256-c547460cd751
	github.com/cosmos/cosmos-sdk => github.com/cosmos/cosmos-sdk v0.52.0-rc.1
	github.com/karalabe/ssz => github.com/berachain/karalabe-ssz v0.3.0-alpha.0
)

require (
	cosmossdk.io/collections v1.3.0
	cosmossdk.io/core v1.0.0
	cosmossdk.io/depinject v1.2.0
	cosmossdk.io/errors v1.0.2
	cosmossdk.io/log v1.6.0
	cosmossdk.io/math v1.5.3
	cosmossdk.io/store v1.10.0-rc.1.0.20241218084712-ca559989da43
	github.com/cenkalti/backoff/v5 v5.0.2
	github.com/cometbft/cometbft v1.0.1-0.20241220100824-07c737de00ff
	github.com/cometbft/cometbft/api v1.0.1-0.20241220100824-07c737de00ff
	github.com/cosmos/cosmos-db v1.1.2
	github.com/cosmos/cosmos-sdk v0.53.0
	github.com/cosmos/go-bip39 v1.0.0
	github.com/crate-crypto/go-kzg-4844 v1.1.0
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/ethereum/c-kzg-4844 v1.0.3
	github.com/go-faster/xor v1.0.0
	github.com/go-playground/validator/v10 v10.26.0
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/hashicorp/go-metrics v0.5.4
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/holiman/uint256 v1.3.2
	github.com/karalabe/ssz v0.2.1-0.20240724074312-3d1ff7a6f7c4
	github.com/labstack/echo/v4 v4.13.3
	github.com/minio/sha256-simd v1.0.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/phuslu/log v1.0.118
	github.com/pkg/errors v0.9.1
	github.com/prysmaticlabs/gohashtree v0.0.4-beta.0.20240624100937-73632381301b
	github.com/prysmaticlabs/prysm/v5 v5.3.0
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8
	github.com/spf13/afero v1.14.0
	github.com/spf13/cast v1.8.0
	github.com/spf13/cobra v1.9.1
	github.com/spf13/pflag v1.0.6
	github.com/spf13/viper v1.20.1
	github.com/umbracle/fastrlp v0.1.0
	go.uber.org/automaxprocs v1.6.0
	golang.org/x/crypto v0.38.0
	golang.org/x/sync v0.14.0
	sigs.k8s.io/yaml v1.4.0
)

// build/test dependencies
require (
	github.com/attestantio/go-eth2-client v0.25.0
	github.com/ethereum/go-ethereum v1.15.5
	github.com/ferranbt/fastssz v0.1.5-0.20240903094032-455b54c08c81
	github.com/kurtosis-tech/kurtosis/api/golang v1.4.3
	github.com/protolambda/zrnt v0.34.1
	github.com/protolambda/ztyp v0.2.2
	github.com/rs/zerolog v1.34.0
	github.com/stretchr/testify v1.10.0
)

require (
	buf.build/gen/go/cometbft/cometbft/protocolbuffers/go v1.36.0-20241120201313-68e42a58b301.1 // indirect
	buf.build/gen/go/cosmos/gogo-proto/protocolbuffers/go v1.36.0-20240130113600-88ef6483f90f.1 // indirect
	cosmossdk.io/api v0.8.0-rc.3 // indirect
	cosmossdk.io/core/testing v0.0.1 // indirect
	cosmossdk.io/schema v1.1.0 // indirect
	cosmossdk.io/x/bank v0.0.0-20241218110910-47409028a73d // indirect
	cosmossdk.io/x/staking v0.0.0-20241218110910-47409028a73d // indirect
	cosmossdk.io/x/tx v1.0.0-alpha.3 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/DataDog/datadog-go v4.8.3+incompatible // indirect
	github.com/DataDog/zstd v1.5.6 // indirect
	github.com/Masterminds/semver/v3 v3.3.1 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/VictoriaMetrics/fastcache v1.12.2 // indirect
	github.com/adrg/xdg v0.4.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.2.0 // indirect
	github.com/bits-and-blooms/bitset v1.17.0 // indirect
	github.com/bufbuild/protocompile v0.14.1 // indirect
	github.com/bytedance/sonic v1.13.1 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/cockroachdb/errors v1.12.0 // indirect
	github.com/cockroachdb/fifo v0.0.0-20240816210425-c5d0cb0b6fc0 // indirect
	github.com/cockroachdb/logtags v0.0.0-20241215232642-bb51bb14a506 // indirect
	github.com/cockroachdb/pebble v1.1.5 // indirect
	github.com/cockroachdb/redact v1.1.6 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/cometbft/cometbft-db v1.0.4 // indirect
	github.com/consensys/bavard v0.1.22 // indirect
	github.com/consensys/gnark-crypto v0.14.0 // indirect
	github.com/containerd/continuity v0.4.4 // indirect
	github.com/cosmos/btcutil v1.0.5 // indirect
	github.com/cosmos/cosmos-proto v1.0.0-beta.5 // indirect
	github.com/cosmos/gogogateway v1.2.0 // indirect
	github.com/cosmos/gogoproto v1.7.0 // indirect
	github.com/cosmos/iavl v1.3.4 // indirect
	github.com/cosmos/ics23/go v0.11.0 // indirect
	github.com/cosmos/ledger-cosmos-go v0.13.3 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20240724233137-53bbb0ceb27a // indirect
	github.com/danieljoos/wincred v1.2.1 // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/dgraph-io/badger/v4 v4.5.1 // indirect
	github.com/dgraph-io/ristretto/v2 v2.1.0 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.7.0 // indirect
	github.com/emicklei/dot v1.6.4 // indirect
	github.com/ethereum/go-verkle v0.2.2 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/getsentry/sentry-go v0.32.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/go-yaml/yaml v2.1.0+incompatible // indirect
	github.com/goccy/go-yaml v1.9.2 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gofrs/flock v0.12.1 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.5-0.20231225225746-43d5d4cd4e0e // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/flatbuffers v25.1.24+incompatible // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
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
	github.com/hashicorp/go-plugin v1.6.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/hdevalence/ed25519consensus v0.2.0 // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/huandu/go-clone v1.6.0 // indirect
	github.com/huandu/skiplist v1.2.1 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/kilic/bls12-381 v0.1.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/kurtosis-tech/kurtosis-portal/api/golang v0.0.0-20230818182330-1a86869414d2 // indirect
	github.com/kurtosis-tech/kurtosis/contexts-config-store v0.0.0-20230818184218-f4e3e773463b // indirect
	github.com/kurtosis-tech/kurtosis/grpc-file-transfer/golang v0.0.0-20230803130419-099ee7a4e3dc // indirect
	github.com/kurtosis-tech/kurtosis/path-compression v0.0.0-20240307154559-64d2929cd265 // indirect
	github.com/kurtosis-tech/stacktrace v0.0.0-20211028211901-1c67a77b5409 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/linxGnu/grocksdb v1.9.8 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mholt/archiver v3.1.1+incompatible // indirect
	github.com/minio/highwayhash v1.0.3 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/oasisprotocol/curve25519-voi v0.0.0-20230904125328-1f23a7beb09a // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/onsi/gomega v1.34.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/opencontainers/runc v1.1.14 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/petermattis/goid v0.0.0-20240813172612-4fcff4a6cae7 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pk910/dynamic-ssz v0.0.4 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.20.5 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/protolambda/bls12-381-util v0.1.0 // indirect
	github.com/prysmaticlabs/fastssz v0.0.0-20241008181541-518c4ce73516 // indirect
	github.com/prysmaticlabs/go-bitfield v0.0.0-20240618144021-706c95b2dd15 // indirect
	github.com/r3labs/sse/v2 v2.10.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.5 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/supranational/blst v0.3.14 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tendermint/go-amino v0.16.0 // indirect
	github.com/thomaso-mirodin/intmath v0.0.0-20160323211736-5dc6d854e46e // indirect
	github.com/tidwall/btree v1.7.0 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ulikunitz/xz v0.5.11 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zondax/hid v0.9.2 // indirect
	github.com/zondax/ledger-go v0.14.3 // indirect
	gitlab.com/yawning/secp256k1-voi v0.0.0-20230925100816-f2616030848b // indirect
	gitlab.com/yawning/tuplehash v0.0.0-20230713102510-df83abbf9a02 // indirect
	go.etcd.io/bbolt v1.4.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	golang.org/x/arch v0.15.0 // indirect
	golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.32.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/genproto v0.0.0-20240624140628-dc46fd24d27d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250106144421-5f5ef82da422 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/grpc v1.71.1 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.5.2 // indirect
	pgregory.net/rapid v1.2.0 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)
