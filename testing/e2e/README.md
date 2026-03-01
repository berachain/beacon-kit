# E2E Tests

End-to-end tests for BeaconKit using [Kurtosis](https://www.kurtosis.com/), a platform for running distributed systems on Docker. Each test suite spins up a full beacon-kit network inside a Kurtosis enclave, runs tests against it, and tears it down automatically.

## Directory Structure

```
testing/e2e/
  config/          # Network configuration structs and defaults
    config.go      #   E2ETestConfig, NetworkConfiguration, NodeSet, etc.
    defaults.go    #   DefaultE2ETestConfig(), PreconfLoadE2ETestConfig()

  suite/           # Shared test framework (Kurtosis orchestration, lifecycle)
    suite.go       #   KurtosisE2ESuite struct and accessors
    setup.go       #   SetupSuite, TearDownSuite, FundAccounts, WaitForFinalizedBlockNumber
    options.go     #   Functional options (WithPreconfLoadConfig)
    logs.go        #   Log fetching, dumping on failure
    constants.go   #   Ether, gas limits, timeouts
    errors.go      #   Shared error variables
    types/         #   Client wrappers (ConsensusClient, RPCClient, EthAccount, etc.)

  standard/        # Standard e2e test suite (package standard_test)
    setup_test.go  #   BeaconKitE2ESuite struct + entry point
    *_test.go      #   One file per feature (blobs, beacon API, staking, ...)

  preconf/         # Preconf e2e test suite (package preconf_test)
    setup_test.go  #   PreconfE2ESuite struct + entry point
    flow_test.go   #   Sequencer serving, validator fetching, fallback
    load_test.go   #   ETH transfers through preconf RPC
```

## How It Works

### Suite Lifecycle

Each test suite (`standard/`, `preconf/`) defines a struct that embeds `suite.KurtosisE2ESuite` and uses [testify suites](https://pkg.go.dev/github.com/stretchr/testify/suite) for lifecycle management.

1. **SetupSuite** (runs once before all tests in a suite):
   - Loads configuration from `config/defaults.go` (or applies custom options)
   - Connects to the local Kurtosis engine
   - Destroys any leftover `e2e-test-enclave` from a previous run
   - Creates a fresh Kurtosis enclave named `e2e-test-enclave`
   - Runs the `kurtosis/` Starlark package to spin up validators, full nodes, seed nodes, and EL clients
   - Sets up consensus clients and a JSON-RPC connection to the execution layer
   - Waits for the network to reach a minimum finalized block number
   - Funds test accounts from the genesis account

2. **Test methods** run against the live network (e.g., `TestBasicStartup`, `Test4844Live`).

3. **TearDownSuite** (runs once after all tests):
   - If any test failed, dumps all service logs to `e2e-logs/<SuiteName>/` (one `.log` file per service)
   - Stops all consensus clients
   - Destroys the Kurtosis enclave

### Enclave = One Suite

Each test suite gets its own Kurtosis enclave. The standard suite uses `DefaultE2ETestConfig()` (5 validators: 3 geth + 2 reth, full nodes, seed nodes). The preconf suite uses `PreconfLoadE2ETestConfig()` which adds a dedicated sequencer node and preconf RPC nodes.

### Build Tags

All e2e test files use `//go:build e2e`. This prevents them from running during `go test ./...`. You must explicitly pass `-tags e2e` (the Makefile targets handle this).

### Log Collection on Failure

When any test in a suite fails, `TearDownSuite` automatically calls `DumpAllServiceLogs()`, which fetches up to 100,000 log lines from every service in the enclave and writes them to:

```
e2e-logs/<TestSuiteName>/<service-name>.log
```

This directory is git-ignored. Logs from previous runs are overwritten.

## Prerequisites

- **Docker** running locally
- **Kurtosis CLI** installed ([instructions](https://docs.kurtosis.com/install)) with the engine started (`kurtosis engine start`)
- **Go 1.25.7+**

## Running Tests

### With Docker Build (builds `beacond:kurtosis-local` image first)

```bash
make test-e2e              # Run ALL e2e tests (standard + preconf, sequentially)
make test-e2e-standard     # Run only standard suite
make test-e2e-preconf      # Run only preconf suite
make test-e2e-4844         # Run only blob tests
make test-e2e-deposits     # Run only deposit tests
```

### Without Docker Build (image must already exist)

```bash
make test-e2e-standard-no-build
make test-e2e-preconf-no-build
```

### Running a Specific Test

```bash
go test -timeout 0 -tags e2e,bls12381,test ./testing/e2e/standard/. -v -testify.m TestBasicStartup
```

## Configuration

Network topology is defined in `config/defaults.go`. Key settings:

- **Validators**: 5 nodes (3 geth + 2 reth) by default
- **Full nodes**: 4 nodes (2 reth + 2 geth)
- **Seed nodes**: 1 geth node
- **Chain ID**: 80087 (devnet). To modify the chain spec, see `config/spec/devnet.go`.
- **EL images**: Configured in `defaultExecutionSettings()` (bera-geth and bera-reth)
- **CL image**: `beacond:kurtosis-local` (built from local source)

To add a new test suite configuration, create a function in `defaults.go` that returns `*E2ETestConfig` and wire it via a functional option in `suite/options.go`.

## Adding New Tests

1. Pick the appropriate suite directory (`standard/` or `preconf/`).
2. Create a new `*_test.go` file with the `//go:build e2e` tag and the matching package name (`standard_test` or `preconf_test`).
3. Add test methods on the suite struct (e.g., `func (s *BeaconKitE2ESuite) TestMyFeature() { ... }`).
4. Use `s.RPCClient()` for EL queries, `s.ConsensusClients()` for CL queries, `s.GenesisAccount()` / `s.TestAccounts()` for funded accounts.

## Debugging

1. Check Kurtosis engine status:
   ```bash
   kurtosis engine status
   kurtosis engine restart
   ```

2. Inspect a running enclave:
   ```bash
   kurtosis enclave inspect e2e-test-enclave
   ```

3. If tests fail, check `e2e-logs/` for service logs.

4. If the enclave wasn't cleaned up (e.g., test was killed), destroy it manually:
   ```bash
   kurtosis enclave rm -f e2e-test-enclave
   ```
