module github.com/berachain/beacon-kit/mod/consensus

go 1.22.4

replace (
	// The following are required to build with the latest version of the cosmos-sdk main branch:
	cosmossdk.io/api => cosmossdk.io/api v0.7.3-0.20240623110059-dec2d5583e39
	cosmossdk.io/core => github.com/berachain/cosmos-sdk/core v0.7.1-0.20240627214151-8adce88544c9
	cosmossdk.io/core/testing => cosmossdk.io/core/testing v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/server/v2 => github.com/berachain/cosmos-sdk/server/v2 v2.0.0-20240628175732-ab2d37478785
	cosmossdk.io/server/v2/appmanager => github.com/berachain/cosmos-sdk/server/v2/appmanager v0.0.0-20240628175732-ab2d37478785
	cosmossdk.io/server/v2/cometbft => github.com/berachain/cosmos-sdk/server/v2/cometbft v0.0.0-20240628175732-ab2d37478785
	cosmossdk.io/server/v2/stf => github.com/berachain/cosmos-sdk/server/v2/stf v0.0.0-20240628175732-ab2d37478785
	cosmossdk.io/x/auth => cosmossdk.io/x/auth v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/x/consensus => cosmossdk.io/x/consensus v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/x/staking => cosmossdk.io/x/staking v0.0.0-20240623110059-dec2d5583e39
	github.com/cosmos/cosmos-sdk => github.com/berachain/cosmos-sdk v0.46.0-beta2.0.20240628175732-ab2d37478785
	// tac0turtle branch for the rand package import.
	github.com/cosmos/crypto => github.com/cosmos/crypto v0.0.0-20240312084433-de8f9c76030d
)

require (
	cosmossdk.io/core v0.12.1-0.20240623110059-dec2d5583e39
	cosmossdk.io/server/v2/cometbft v0.0.0-00010101000000-000000000000
	github.com/berachain/beacon-kit/mod/primitives v0.0.0-20240627134700-de48919ec4d6
	github.com/cometbft/cometbft v1.0.0-alpha.2.0.20240627055712-4f91afce3247
	github.com/cosmos/gogoproto v1.5.0
	github.com/sourcegraph/conc v0.3.0
)

require (
	buf.build/gen/go/cometbft/cometbft/protocolbuffers/go v1.34.2-20240312114316-c0d3497e35d6.2 // indirect
	buf.build/gen/go/cosmos/gogo-proto/protocolbuffers/go v1.34.2-20240130113600-88ef6483f90f.2 // indirect
	cosmossdk.io/api v0.7.5 // indirect
	github.com/berachain/beacon-kit/mod/errors v0.0.0-20240624210933-3baa52ebb587 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.3 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cometbft/cometbft/api v1.0.0-rc.1 // indirect
	github.com/cosmos/cosmos-proto v1.0.0-beta.5 // indirect
	github.com/cosmos/crypto v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
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
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/supranational/blst v0.3.12 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240624140628-dc46fd24d27d // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
