module github.com/berachain/beacon-kit/mod/node

go 1.22.5

replace (
	// The following are required to build with the latest version of the cosmos-sdk main branch:
	cosmossdk.io/api => cosmossdk.io/api v0.7.3-0.20240726151427-2e0884564fdb
	cosmossdk.io/core/testing => cosmossdk.io/core/testing v0.0.0-20240726151427-2e0884564fdb
	github.com/berachain/beacon-kit/mod/consensus => ../consensus
	github.com/cosmos/cosmos-sdk => github.com/berachain/cosmos-sdk v0.46.0-beta2.0.20240725053043-79fa56d34c79
)

require (
	github.com/berachain/beacon-kit/mod/log v0.0.0-20240726221339-a8bfeebf8ecf
	github.com/berachain/beacon-kit/mod/runtime v0.0.0-20240723155519-565f208d5482
)

require (
	github.com/berachain/beacon-kit/mod/errors v0.0.0-20240705193247-d464364483df // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/getsentry/sentry-go v0.28.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)
