<div align="center">
  <a href="https://github.com/berachain/beacon-kit">
    <img alt="beacon-kit-banner" src=".github/assets/banner.png" width="auto" height="auto">
  </a>
</div>
<h2>
  Berachain's Consensus Client
</h2>

<div>

[![CI status](https://github.com/berachain/beacon-kit/workflows/pipeline/badge.svg)](https://github.com/berachain/beacon-kit/actions/workflows/pipeline.yml)
[![CodeCov](https://codecov.io/gh/berachain/beacon-kit/graph/badge.svg?token=0l5iJ3ZbzV)](https://codecov.io/gh/berachain/beacon-kit)
[![Telegram Chat](https://img.shields.io/endpoint?color=neon&logo=telegram&label=chat&url=https%3A%2F%2Ftg.sumanjay.workers.dev%2Fbeacon_kit)](https://t.me/beacon_kit)
[![X Follow](https://img.shields.io/twitter/follow/berachain)](https://x.com/berachain)
[![Discord](https://img.shields.io/discord/924442927399313448?label=discord)](https://discord.gg/berachain)

</div>

## What is BeaconKit?

BeaconKit is Berachain's consensus client. It implements the Ethereum consensus layer specification with Berachain-specific modifications, using a modified [CometBFT](https://github.com/berachain/cometbft) for consensus instead of Ethereum's standard beacon chain consensus.

Berachain is a high-performance L1 blockchain powered by [Proof of Liquidity](https://docs.berachain.com/general/introduction/what-is-proof-of-liquidity) (PoL), an incentive mechanism that aligns validators, protocols, and users through its three-token system (BERA, BGT, HONEY).

BeaconKit communicates with [Bera-Reth](https://github.com/berachain/bera-reth) (the execution client) via the standard [Engine API](https://github.com/ethereum/execution-apis/blob/main/src/engine), forming the two-client architecture used by Berachain nodes.

### Key Differences from Ethereum

- **CometBFT Consensus** — Uses timeout-based rounds instead of fixed 12-second slots. If a proposer fails, a new round starts at the same height with a different proposer.
- **Single-Slot Finality** — Blocks are finalized immediately when 2/3+ validators commit, unlike Ethereum's ~13-minute finalization through Casper FFG checkpoints.
- **No Missed Slots** — Block heights are strictly sequential with no gaps. Every height eventually produces a block through round-based consensus.
- **No Attestations or Committees** — Consensus is driven by CometBFT's validator voting (prevotes/precommits) rather than Ethereum's attestation committees and sync committees.
- **Capped Validator Set** — The active validator set is capped (currently 69 validators) rather than being open-ended.
- **EVM Inflation Withdrawal** — Every block includes a mandatory EVM inflation withdrawal as its first withdrawal, a Berachain-specific economic mechanism.

## Networks

| Network | Chain ID | Description |
|---------|----------|-------------|
| **Mainnet** | 80094 | Production network |
| **Bepolia** | 80069 | Public testnet |
| **Devnet** | 80087 | Local development |

Network configurations are stored in `testing/networks/<chain-id>/`.

## Execution Client

BeaconKit requires [**Bera-Reth**](https://github.com/berachain/bera-reth), a Berachain-specific fork of [Reth](https://github.com/paradigmxyz/reth). Communication happens over the Engine API with JWT authentication.

## Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/engine/install/)
- [Go 1.25.7+](https://go.dev/doc/install)
- [Foundry](https://getfoundry.sh/)

### Running a Local Devnet

Open two terminals side by side.

**Terminal 1** — Start the consensus client:

```bash
make start
```

**Terminal 2** — Start the execution client (after `beacond` is running, since `make start` generates the execution genesis file):

```bash
make start-reth
```

The devnet runs with chain ID **80087**. The following dev account is preloaded with the native token:

```
Address:     0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4
Private Key: 0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306
```

### Multinode Devnet (Kurtosis)

For running a multinode local devnet using [Kurtosis](https://www.kurtosis.com/):

```bash
make start-devnet
```

See the [Kurtosis README](kurtosis/README.md) for full details on configuration and usage.

## Project Structure

```
beacon/              Core beacon chain logic (blockchain service, validator service)
chain/               Chain specifications and network configs
cli/                 CLI commands and configuration
config/              Configuration templates (TOML)
consensus-types/     Consensus layer types (blocks, states, validators, SSZ)
contracts/           Solidity contracts (deposit contract, staking)
da/                  Data availability (blob management, KZG commitments)
engine-primitives/   Execution engine types (Engine API types, payloads)
execution/           Execution client integration (Engine API client, deposit syncing)
kurtosis/            Kurtosis multinode devnet deployment
node-api/            REST API implementation (beacon API handlers)
node-core/           Core infrastructure (dependency injection, service registry)
primitives/          Basic types and constants (Slot, Gwei, ValidatorIndex)
state-transition/    State transition logic (state machine, fork transitions)
storage/             Database backends (block store, blob store, beacon DB)
testing/             Test utilities, network configs, e2e and simulation tests
```

## Development

### Building

```bash
make build                # Build beacond binary to build/bin/beacond
make install              # Install beacond to $GOPATH/bin
make build-docker         # Build Docker image
```

### Testing

```bash
make test                 # Run all tests (unit + forge)
make test-unit            # Run unit tests with coverage
make test-unit-bench      # Run benchmarks
make test-unit-fuzz       # Run Go fuzz tests
make test-simulated       # Run simulation tests
make test-e2e             # Run e2e tests (builds Docker image first)
make test-forge-cover     # Run Solidity tests with coverage
```

### Linting & Formatting

```bash
make lint                 # Run all linters
make format               # Run all formatters
make golangci-fix         # Auto-fix Go linting issues
make gosec                # Run security scanner
make nilaway              # Run nil pointer checker
make vulncheck            # Run govulncheck vulnerability scanner
```

### Code Generation

```bash
make generate             # Run all code generation
make proto                # Generate protobuf code
make generate-check       # Verify generated code is up to date
```

## Documentation

- [Berachain Docs](https://docs.berachain.com/) — Official documentation
- [What is BeaconKit](https://docs.berachain.com/validators/beaconkit/overview) — Conceptual overview
- [Node Quickstart](https://docs.berachain.com/validators/operations/quickstart) — Deploy a mainnet or testnet node
- [Docker Devnet Guide](https://docs.berachain.com/validators/guides/local-devnet-docker) — Run a devnet with Docker
- [Kurtosis Guide](https://docs.berachain.com/validators/guides/local-devnet-kurtosis) — Multinode deployment with Kurtosis
- [Kurtosis README](kurtosis/README.md) — Local Kurtosis configuration and usage

## License

BeaconKit is licensed under [BUSL-1.1](LICENSE).
