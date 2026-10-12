# CLAUDE.md

Guidance for Claude Code (claude.ai/code) when working in this repository. Facts here were verified against commit `2059e2abe` (August 2026). When this file and the code disagree, trust the code.

## Overview

BeaconKit is Berachain's consensus layer client. It implements the Ethereum consensus spec but replaces Ethereum's consensus (LMD-GHOST + Casper FFG) with a modified CometBFT, and drives an execution client over the standard Engine API. bera-reth is the supported execution client.

Key differences from a standard Ethereum beacon chain:

- CometBFT height maps 1:1 to beacon slot. Heights are sequential, with no gaps and no missed slots.
- No fixed 12s slots. Block pacing is driven by SBT (stable block time) via ABCI `NextBlockDelay` (`consensus/cometbft/service/delay/`, target 2s). `timeout_commit` is deprecated and forced to 0. Propose/prevote/precommit timeouts have enforced 2000ms minimums (`consensus/cometbft/service/configs.go`).
- If a consensus round fails, CometBFT retries at the same height with a new proposer. Height advances only after FinalizeBlock.
- Custom validator set cap. At epoch end, the lowest-effective-balance validators are ejected until the cap holds (69 on mainnet, `state-transition/core/state_processor_validators.go`). There is no churn cap.

## Directory map

- `beacon/` - blockchain service (process/finalize blocks), validator service (build blocks), deposit fetching
- `chain/` - chain spec accessors, chain IDs, fork-activation helpers (`ActiveForkVersionForTimestamp`)
- `cli/` - beacond commands and flags
- `cmd/beacond/` - main entrypoint, DI wiring in `defaults.go`
- `config/` - config struct, app.toml template (`config/template/`), per-network chain specs (`config/spec/`)
- `consensus/` - CometBFT service, ABCI++ handlers, SBT block-delay logic, halt-height logic
- `consensus-types/` - blocks, states, validators, SSZ serialization
- `contracts/` - Solidity deposit contract and Go bindings
- `da/` - blob verification (`da/blob/`), KZG verify (`da/kzg/`), blob availability store (`da/store/`)
- `engine-primitives/` - Engine API types, payloads, withdrawals
- `errors/` - error wrapper (`mod.go`)
- `execution/` - Engine API client (`client/`), engine wrapper (`engine/`), deposit contract log reader (`deposit/`)
- `gethlib/` - vendored geth types and bindings, EL chain config (Osaka/Osaka1)
- `kurtosis/` - multi-node devnet
- `node-api/` - REST API server and handlers (beacon, proof, debug, events, ...)
- `node-core/` - DI components (`components/`), service registry, node builder
- `observability/` - telemetry
- `payload/` - local payload builder, payload attributes, payload ID cache
- `primitives/` - basic types, crypto, fork versions (`primitives/version/`)
- `scripts/build/` - all Makefile logic (`*.mk`)
- `state-transition/` - state processor (slots, blocks, epochs, forks)
- `storage/` - `beacondb/` (state), `block/` (blocks by slot), `deposit/`, `filedb/`, plus `db/`, `encoding/`, `interfaces/`
- `testing/` - e2e, simulated tests, network configs (`networks/80069`, `networks/80094`), devnet files (`files/`)

## Development commands

### Build

```bash
make build                   # build beacond to build/bin/beacond
make install                 # install beacond to $GOPATH/bin
make build-docker            # build Docker image
```

### Run

```bash
make start                   # ephemeral devnet CL node (home: .tmp/beacond)
make start-reth              # ephemeral bera-reth EL node (docker, nightly image)
make start-custom <spec>     # devnet with a custom chain spec TOML
make start-bepolia           # testnet node (also start-reth-bepolia)
make start-mainnet           # mainnet node (also start-reth-mainnet)
make start-devnet            # kurtosis multi-node devnet (also stop-devnet, rm-devnet)
```

Start the beacon node before the execution client. It generates the required genesis configuration. bera-reth is the only supported EL, the `start-geth`/`start-besu`/etc targets no longer exist.

### Test

```bash
make test                    # test-unit + test-forge-fuzz
make test-unit               # unit tests with coverage
make test-unit-no-coverage
make test-unit-quick         # quick tests only (build tag: quick)
make test-unit-bench
make test-unit-fuzz
make test-simulated          # simulation tests (chaos, forks)
make test-e2e                # e2e with real EL clients (builds Docker first)
make test-forge-unit         # Solidity tests (also test-forge-fuzz, test-forge-cover)
make test-halt-swap-resume   # halt-height upgrade test (also test-halt-swap-resume-time)
```

Go build tags in use: `test`, `e2e`, `simulated`, `quick`. Test invocations also pass `bls12381`, which gates dependency code, not repo sources.

### Lint, format, codegen

```bash
make lint / format / golangci-fix
make gosec / nilaway / vulncheck / slither
make generate / proto / generate-check
```

## Workflow

- Run `make format` and `make lint` before finishing a change. CI enforces lint, license headers, and generate-check.
- For fast iteration, test a single package with `go test -race -tags bls12381,test ./<package>/...`. A bare `go test` misses tag-gated tests, and `make test-unit` runs everything with coverage and is slow.

## beacond CLI

```bash
beacond init <moniker>       # initialize node (writes config.toml, app.toml, genesis)
beacond start                # run the node
beacond rollback             # roll back state by one height
beacond status               # node status via RPC
beacond version
beacond jwt generate|validate
beacond genesis add-premined-deposit | collect-premined-deposits |
        set-deposit-storage | execution-payload | validator-root
beacond deposit create-validator | validate | validator-keys | db-check
beacond comet show-node-id | show-validator | unsafe-reset-all | bootstrap-state | ...
```

Key flags:

```bash
--home <path>                              # default ~/.beacond
--beacon-kit.chain-spec <spec>             # devnet|testnet|mainnet (default mainnet)
--beacon-kit.chain-spec-file <path>        # custom spec TOML
--beacon-kit.engine.jwt-secret-path <path>
--beacon-kit.engine.rpc-dial-url <url>
--beacon-kit.kzg.trusted-setup-path <path>
--beacon-kit.node-api.enabled
--halt-height <h> / --halt-time <unix>     # graceful halt for coordinated upgrades
```

## Configuration

- `app.toml` holds all BeaconKit settings under `[beacon-kit]` sections: `engine`, `logger`, `kzg`, `payload-builder`, `validator`, `block-store-service`, `node-api`. Template in `config/template/template.go`.
- `config.toml` is the CometBFT config (p2p, rpc, consensus).
- Precedence: CLI flags > env vars > config file > defaults.
- Env vars are prefixed with the uppercased binary basename (`BEACOND_` for beacond), with `.` and `-` mapped to `_`. Example: `BEACOND_BEACON_KIT_ENGINE_RPC_DIAL_URL`. The prefixed viper and per-flag `BindEnv` calls live in `cli/config/server.go`. A near-identical but separate setup in `cli/commands/server/cmd/execute.go` configures the global viper with an empty prefix, the two are not duplicates of each other.

## Networks

| Network | EL chain ID | Config |
|---|---|---|
| Devnet | 80087 | `testing/files/` (used by `make start`) |
| Testnet (Bepolia) | 80069 | `testing/networks/80069/` |
| Mainnet | 80094 | `testing/networks/80094/` |

Chain IDs live in `chain/chain_ids.go`. Per-network fork schedules live in `config/spec/{devnet,testnet,mainnet}.go`.

## Architecture

### ABCI++ flow

CometBFT drives everything. The ABCI handlers live in `consensus/cometbft/service/` (`abci.go`, `prepare_proposal.go`, `process_proposal.go`, `finalize_block.go`) and delegate to beacon services:

1. **PrepareProposal** (proposer only, may run once per round) → `BlockBuilder.BuildBlockAndSidecars` (`beacon/validator/block_builder.go`). Sends FCU with payload attributes, gets payload ID, calls getPayload, builds blob sidecars, returns the block.
2. **ProcessProposal** (all nodes) → `Blockchain.ProcessProposal` (`beacon/blockchain/process_proposal.go`). Verifies structure, signatures, KZG proofs, and the payload via newPayload. Validates only, commits nothing.
3. **FinalizeBlock** (all nodes, once per height) → `Blockchain.FinalizeBlock` (`beacon/blockchain/finalize_block.go`). Runs the full state transition, persists state/block/blobs, sends the post-block FCU.

### Services

The service registry (`node-core/components/service_registry.go`) starts services in order: shutdown, validator, node-api, reporting, telemetry, engine-client, blockchain, cometbft. Shutdown is reverse order with a 5-minute default timeout. Services implement `Start(ctx) error`, `Stop() error`, `Name() string`. All wiring uses `cosmossdk.io/depinject`, with `Provide*` functions in `node-core/components/` assembled in `cmd/beacond/defaults.go`.

### State transition

`state-transition/core/state_processor.go`. Main entry points: `Transition`, `ProcessSlots` (advances exactly one slot), `ProcessBlock`, and `ProcessFork` (`state_processor_forks.go`). Epoch processing is internal to `ProcessSlots`. State changes are cached and flushed on commit.

### Execution engine

`execution/engine/engine.go` wraps the Engine API client (`execution/client/`). Methods: `NotifyNewPayload`, `NotifyForkchoiceUpdate`, `GetPayload`. RPC timeout defaults to 2s and is floored at 2s. Errors are classified fatal vs retryable in `execution/client/errors.go`.

### Forks and deposits

Supported CL fork versions: `deneb`, `deneb1`, `electra`, `electra1`, `fulu` (`primitives/version/`). Activation is timestamp-based (`chain/helpers.go`). Osaka/Osaka1 are EL-side forks configured in the EL genesis (`gethlib/types/config.go`).

Deposits are fork-gated (EIP-6110). From Fulu onward, deposits arrive as execution-layer requests in the payload. Before Fulu, they are read from deposit contract logs (`execution/deposit/`). `validateDepositSource` (`state-transition/core/state_processor_staking.go`) rejects the wrong source for the active fork. `beacon/deposits/` handles the one-time catchup at the Fulu boundary.

### Storage

- CometBFT application state is IAVL-backed. App hash = beacon state root.
- `storage/beacondb/` - beacon state (StateDB), fork-aware, context-scoped
- `storage/block/` - blocks indexed by slot
- `da/store/` - blob sidecars, pruned by the availability window (slots, not TTL)
- `storage/deposit/` - deposit store (pruning currently disabled)

### Halt mechanism

`--halt-height`/`--halt-time` stop the node gracefully at a target point for coordinated upgrades (`consensus/cometbft/service/commit.go`). A halted node refuses to restart past its halt point until the flags are cleared or raised.

## Gotchas

- Every `.go` file needs the BUSL-1.1 license header from `LICENSE.header`. `make lint` fails without it. Run `make license-fix` to add headers to new files.
- `make generate-check` excludes `*.ssz.go` from its staleness check (`scripts/build/codegen.mk`), so stale SSZ codegen passes CI and surfaces at runtime as a serialization mismatch. After changing any SSZ-serialized type, run `make generate` and commit the `.ssz.go` diffs.
- ProcessProposal must never mutate committed state. Only FinalizeBlock commits.
- A failed round re-proposes at the same height. State built during a previous round is reset.
- Deposit handling, hysteresis, and inflation parameters are fork-gated. Check the active fork version before assuming behavior.
- One-off mainnet state fixes live in `state-transition/core/state_processor_fixes.go` and are height-gated. Do not generalize from them.
- `errors.IsFatal` (`errors/mod.go`) exists but the live fatal-error path uses `client.IsFatalError` from `execution/client`.
