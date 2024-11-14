module github.com/berachain/beacon-kit/mod/node-core

go 1.23.0

replace (
	// The following are required to build with the latest version of the cosmos-sdk main branch:
	cosmossdk.io/api => cosmossdk.io/api v0.7.3-0.20240806152830-8fb47b368cd4
	cosmossdk.io/core => cosmossdk.io/core v0.0.0-20240806152830-8fb47b368cd4
	cosmossdk.io/core/testing => cosmossdk.io/core/testing v0.0.0-20240806152830-8fb47b368cd4
	github.com/berachain/beacon-kit/mod/cli => ../cli
	github.com/berachain/beacon-kit/mod/config => ../config
	github.com/berachain/beacon-kit/mod/consensus => ../consensus
	github.com/berachain/beacon-kit/mod/node-api => ../node-api
	github.com/berachain/beacon-kit/mod/observability => ../observability
	github.com/berachain/beacon-kit/mod/storage => ../storage
	github.com/cosmos/cosmos-sdk => github.com/berachain/cosmos-sdk v0.46.0-beta2.0.20240808182639-7bdbf06a94f2
	ithub.com/berachain/beacon-kit/mod/cli => ../cli
)

require (
	cosmossdk.io/core v1.0.0
	cosmossdk.io/depinject v1.0.0
	cosmossdk.io/store/v2 v2.0.0-20240821144902-e88c138760a3
	github.com/berachain/beacon-kit/mod/beacon v0.0.0-20240821052951-c15422305b4e
	github.com/berachain/beacon-kit/mod/cli v0.0.0-20241107170417-7905e3d59a1d
	github.com/berachain/beacon-kit/mod/config v0.0.0-20241113214258-240f617103ad
	github.com/berachain/beacon-kit/mod/consensus v0.0.0-20241107170417-7905e3d59a1d
	github.com/berachain/beacon-kit/mod/consensus-types v0.0.0-20240904192942-99aeabe6bb1f
	github.com/berachain/beacon-kit/mod/da v0.0.0-20240820191615-398849c34954
	github.com/berachain/beacon-kit/mod/engine-primitives v0.0.0-20240809202957-3e3f169ad720
	github.com/berachain/beacon-kit/mod/errors v0.0.0-20240806211103-d1105603bfc0
	github.com/berachain/beacon-kit/mod/execution v0.0.0-20240820191615-398849c34954
	github.com/berachain/beacon-kit/mod/log v0.0.0-20240821000339-4d4242ba4a50
	github.com/berachain/beacon-kit/mod/node-api v0.0.0-20241107170417-7905e3d59a1d
	github.com/berachain/beacon-kit/mod/node-api/engines v0.0.0-20240806160829-cde2d1347e7e
	github.com/berachain/beacon-kit/mod/observability v0.0.0-unpublished
	github.com/berachain/beacon-kit/mod/payload v0.0.0-20240705193247-d464364483df
	github.com/berachain/beacon-kit/mod/primitives v0.0.0-20240911165923-82f71ec86570
	github.com/berachain/beacon-kit/mod/state-transition v0.0.0-20240717225334-64ec6650da31
	github.com/berachain/beacon-kit/mod/storage v0.0.0-20240822205119-6d7f90fac7d7
	github.com/cometbft/cometbft v1.0.0-rc1.0.20240806094948-2c4293ef36c4
	github.com/cometbft/cometbft/api v1.0.0-rc.1.0.20240806094948-2c4293ef36c4
	github.com/cosmos/cosmos-db v1.0.2
	github.com/cosmos/cosmos-sdk v0.53.0
	github.com/crate-crypto/go-kzg-4844 v1.1.0
	github.com/hashicorp/go-metrics v0.5.3
	github.com/spf13/afero v1.11.0
	github.com/spf13/cast v1.7.0
)

require (
	cosmossdk.io/api v0.7.5 // indirect
	cosmossdk.io/errors/v2 v2.0.0-20240731132947-df72853b3ca5 // indirect
	cosmossdk.io/log v1.4.1 // indirect
	cosmossdk.io/x/tx v0.13.4-0.20240623110059-dec2d5583e39 // indirect
	github.com/VictoriaMetrics/fastcache v1.12.2 // indirect
	github.com/berachain/beacon-kit/mod/chain-spec v0.0.0-20240705193247-d464364483df // indirect
	github.com/berachain/beacon-kit/mod/geth-primitives v0.0.0-20240806160829-cde2d1347e7e // indirect
	github.com/cockroachdb/fifo v0.0.0-20240616162244-4768e80dfb9a // indirect
	github.com/cosmos/cosmos-proto v1.0.0-beta.5 // indirect
	github.com/gabriel-vasile/mimetype v1.4.6 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.22.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/karalabe/ssz v0.2.1-0.20240724074312-3d1ff7a6f7c4 // indirect
	github.com/labstack/echo/v4 v4.12.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/phuslu/log v1.0.110 // indirect
	github.com/prysmaticlabs/go-bitfield v0.0.0-20240618144021-706c95b2dd15 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/umbracle/fastrlp v0.1.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)

require (
	buf.build/gen/go/cometbft/cometbft/protocolbuffers/go v1.34.2-20240701160653-fedbb9acfd2f.2 // indirect
	buf.build/gen/go/cosmos/gogo-proto/protocolbuffers/go v1.34.2-20240130113600-88ef6483f90f.2 // indirect
	cosmossdk.io/collections v0.4.0 // indirect
	cosmossdk.io/errors v1.0.1 // indirect
	cosmossdk.io/math v1.3.0 // indirect
	cosmossdk.io/store v1.1.1-0.20240418092142-896cdf1971bc
	cosmossdk.io/x/auth v0.0.0-20240806152830-8fb47b368cd4 // indirect
	cosmossdk.io/x/bank v0.0.0-20240806152830-8fb47b368cd4 // indirect
	cosmossdk.io/x/consensus v0.0.0-20240806152830-8fb47b368cd4 // indirect
	cosmossdk.io/x/staking v0.0.0-20240806152830-8fb47b368cd4 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/DataDog/datadog-go v4.8.3+incompatible // indirect
	github.com/DataDog/zstd v1.5.6 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/berachain/beacon-kit/mod/async v0.0.0-20240821213929-f32b8e2dc5c8
	github.com/bgentry/speakeasy v0.2.0 // indirect
	github.com/bits-and-blooms/bitset v1.13.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v1.1.1 // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/cometbft/cometbft-db v0.13.0 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.13.0 // indirect
	github.com/cosmos/btcutil v1.0.5 // indirect
	github.com/cosmos/crypto v0.1.2 // indirect
	github.com/cosmos/go-bip39 v1.0.0 // indirect
	github.com/cosmos/gogoproto v1.7.0 // indirect
	github.com/cosmos/iavl v1.2.1-0.20240731145221-594b181f427e // indirect
	github.com/cosmos/ics23/go v0.10.0 // indirect
	github.com/cosmos/ledger-cosmos-go v0.13.3 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20240724233137-53bbb0ceb27a // indirect
	github.com/danieljoos/wincred v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/dgraph-io/badger/v4 v4.2.0 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.7.0 // indirect
	github.com/emicklei/dot v1.6.2 // indirect
	github.com/ethereum/c-kzg-4844 v1.0.3 // indirect
	github.com/ethereum/go-ethereum v1.14.7 // indirect
	github.com/ethereum/go-verkle v0.1.1-0.20240306133620-7d920df305f0 // indirect
	github.com/ferranbt/fastssz v0.1.4-0.20240629094022-eac385e6ee79
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
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hdevalence/ed25519consensus v0.2.0 // indirect
	github.com/holiman/uint256 v1.3.1 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/linxGnu/grocksdb v1.9.2 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/minio/highwayhash v1.0.3 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/oasisprotocol/curve25519-voi v0.0.0-20230904125328-1f23a7beb09a // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/petermattis/goid v0.0.0-20240813172612-4fcff4a6cae7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.20.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/prysmaticlabs/gohashtree v0.0.4-beta.0.20240624100937-73632381301b // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/rs/cors v1.11.0 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.5 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/cobra v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.9.0
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/supranational/blst v0.3.13 // indirect
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
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sync v0.8.0
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/term v0.25.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto v0.0.0-20240624140628-dc46fd24d27d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240701130421-f6361c86f094 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240711142825-46eb208f015d // indirect
	google.golang.org/grpc v1.65.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.5.1 // indirect
	pgregory.net/rapid v1.1.0 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
