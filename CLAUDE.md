# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Introduction & Overview

BeaconKit is a modular consensus client implementation that uses a modified CometBFT (Tendermint) for consensus instead of the standard Ethereum beacon chain consensus. It implements the Ethereum consensus layer specification while supporting all major Ethereum execution clients through the Engine API.

### Key Differences from Standard Ethereum
- **Consensus**: Uses CometBFT's Tendermint consensus instead of Ethereum's LMD-GHOST + Casper FFG
- **Block Timing**: No fixed 12-second slots; uses timeout-based rounds
- **Validator Set**: Custom validator set cap implementation
- **Block Heights**: Sequential without gaps (CometBFT height = beacon slot)
- **No Missed Slots**: Every height produces a block eventually through round-based consensus

## Project Structure

### Key Directories
- `/beacon/` - Core beacon chain logic (blockchain service, validator service, payload coordination)
- `/chain/` - Chain specifications and helpers (spec interfaces, network configs, genesis generation)
- `/cli/` - CLI commands and configuration (all beacond commands, flag definitions)
- `/config/` - Configuration management (TOML templates, default configs)
- `/consensus-types/` - Consensus layer types (blocks, states, validators, SSZ serialization)
- `/contracts/` - Solidity contracts (deposit contract, staking, Go bindings)
- `/da/` - Data availability layer (blob management, KZG commitments, proofs)
- `/engine-primitives/` - Execution engine types (Engine API types, payloads, withdrawals)
- `/execution/` - Execution client integration (Engine API client, deposit syncing)
- `/node-api/` - REST API implementation (beacon API handlers, server setup)
- `/node-core/` - Core node infrastructure (DI components, service registry, node builder)
- `/primitives/` - Basic types and constants (Slot, Gwei, ValidatorIndex, crypto primitives)
- `/state-transition/` - State transition logic (state machine, fork transitions, operations)
- `/storage/` - Database and storage backends (beacon DB, block store, blob store, pruning)
- `/testing/` - Test utilities and networks (e2e tests, simulations, test fixtures)

### Configuration Files
- Home directory: `~/.beacond` (or `.tmp/beacond` for local testing)
- Config files: `config.toml`, `app.toml`, `client.toml`, `genesis.json`
- Configuration sections in `config.toml`:
  - `beacon-kit.engine` - Execution client settings
  - `beacon-kit.logger` - Logging configuration
  - `beacon-kit.kzg` - KZG trusted setup
  - `beacon-kit.payload-builder` - Block building
  - `beacon-kit.validator` - Validator settings
  - `beacon-kit.block-store-service` - Block storage
  - `beacon-kit.node-api` - API server settings

## Development Commands

### Building
```bash
make build                    # Build beacond binary to build/bin/beacond
make build-docker            # Build Docker image
make install                 # Install beacond to $GOPATH/bin
```

### Running
```bash
make start                   # Start ephemeral devnet node (chain ID: 80087)
make start-custom <spec>     # Start with custom chain spec TOML file
make start-<client>          # Start execution client (reth/geth/nethermind/besu/erigon/ethereumjs)
```
Note: Always start the beacon node before the execution client, as it generates the required genesis configuration.

### Testing
```bash
make test                    # Run all tests (unit + forge)
make test-unit               # Run unit tests with coverage
make test-unit-no-coverage   # Run unit tests without coverage
make test-unit-bench         # Run benchmarks
make test-unit-fuzz          # Run Go fuzz tests
make test-simulated          # Run simulation tests (chaos, forks)
make test-e2e                # Run e2e tests (builds Docker first)
make test-forge-cover        # Run Solidity tests with coverage
```

### Linting & Formatting
```bash
make lint                    # Run all linters
make format                  # Run all formatters
make golangci-fix            # Auto-fix linting issues
make gosec                   # Run security scanner
make nilaway                 # Run nil pointer checker
```

### Code Generation
```bash
make generate                # Run all code generation
make proto                   # Generate protobuf code
make generate-check          # Verify generated code is up to date
```

## Environment Variables

BeaconKit uses Viper for configuration management, which supports environment variables with automatic binding.

### Environment Variable Configuration
- **Prefix**: Configured per application (e.g., `BEACOND_` for the main beacon daemon)
- **Key Mapping**: Configuration keys are transformed for environment variables:
  - Dots (`.`) are replaced with underscores (`_`)
  - Hyphens (`-`) are replaced with underscores (`_`)
  - All keys are uppercased
- **Auto-binding**: Environment variables are automatically bound using `viper.AutomaticEnv()`

### Configuration Precedence
The configuration system follows this precedence order (highest to lowest):
1. **CLI flags** - Command-line arguments override everything
2. **Environment variables** - Override config file values
3. **Config file** - Values from `config.toml`, `app.toml`, etc.
4. **Default values** - Hardcoded defaults in the application

### Implementation Details

The environment variable system is initialized in `cli/commands/server/cmd/execute.go`:

```go
viper.SetEnvPrefix(envPrefix)
viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
viper.AutomaticEnv()
```

This setup ensures that:
- All environment variables with the configured prefix are automatically recognized
- Configuration keys are properly mapped to environment variable names
- Values from environment variables override configuration file values

## CLI Usage

### beacond Commands
```bash
beacond init                                    # Initialize a new node
beacond start                                   # Start the beacon node
beacond rollback                                # Rollback blockchain state
beacond genesis add-premined-deposit            # Add premined deposits to genesis
beacond genesis collect-premined-deposits       # Collect premined deposits
beacond genesis set-deposit-storage             # Set deposit contract storage
beacond genesis execution-payload               # Generate execution payload
beacond deposit create-validator                # Create validator deposit
```

### Key Flags
```bash
--beacon-kit.chain-spec <spec>                  # Chain spec: devnet/testnet/mainnet
--beacon-kit.chain-spec-file <path>             # Custom chain spec TOML file
--beacon-kit.engine.jwt-secret-path <path>      # JWT secret for EL auth
--beacon-kit.engine.rpc-dial-url <url>          # Execution client RPC URL
--beacon-kit.kzg.trusted-setup-path <path>      # KZG trusted setup file
--beacon-kit.node-api.enabled                   # Enable REST API
--home <path>                                   # Node home directory (default: ~/.beacond)
```

## Network Specifications

### Supported Networks
- **Devnet** - Chain ID: 80087 (local development)
- **Testnet/Bepolia** - Chain ID: 80069 (public testnet)
- **Mainnet** - Chain ID: 80094 (production)

Network configurations are in `testing/networks/<chain-id>/`

## Architecture Overview

BeaconKit implements a modular EVM consensus client using a modified CometBFT for consensus and supporting all major Ethereum execution clients through the Engine API.

### System Architecture

```
┌───────────────────────────────────────────────────────────────┐
│                      Service Registry                         │
│              (Lifecycle orchestrator for all services)        │
│  Startup Order:                                               │
│  1. ShutdownService → 2. ValidatorService → 3. NodeAPIServer  │
│  4. ReportingService → 5. TelemetryService → 6. EngineClient  │
│  7. ChainService → 8. CometBFTService                         │
└──────────────────────────────┬────────────────────────────────┘
                               │ manages lifecycle
                               ▼
                ┌─────────────────────────────────────┐
                │         CometBFT Service            │
                │  (Consensus orchestrator via P2P)   │
                │  ┌─────────────────────────────┐    │
                │  │ ABCI++ Interface:           │    │
                │  │ • InitChain                 │    │
                │  │ • PrepareProposal           │    │
                │  │ • ProcessProposal           │    │
                │  │ • FinalizeBlock             │    │
                │  └─────────────────────────────┘    │
                └───────────────┬─────────────────────┘
                                │
            PrepareProposal     │    ProcessProposal/FinalizeBlock
                    ┌───────────┴─────────┐
                    ▼                     ▼
        ┌─────────────────────┐ ┌────────────────────┐
        │  Validator Service  │ │ Blockchain Service │
        │  (Block builder)    │ │ (Block processor)  │
        │ ┌─────────────────┐ │ │ ┌────────────────┐ │
        │ │ • StateProcessor│ │ │ │• StateProcessor│ │
        │ │ • BlobFactory   │ │ │ │• BlobProcessor │ │
        │ │ • PayloadBuilder│ │ │ │• LocalBuilder  │ │
        │ │ • Signer        │ │ │ │• DepositFetcher│ │
        │ └─────────────────┘ │ │ └────────────────┘ │
        └──────────┬──────────┘ └────────┬───────────┘
                   │                     │
                   │ both use            │ both use
                   ▼                     ▼
        ┌────────────────────────────────────────┐
        │         Execution Engine               │
        │   (Engine API client wrapper)          │
        │  • forkchoiceUpdate                    │
        │  • newPayload / getPayload             │
        └──────────────┬─────────────────────────┘
                       │ communicates with
                       ▼
                External EL Client
                (Geth/Reth/etc.)

┌───────────────────────────────────────────────────────────────┐
│                      Storage Backend                          │
│                 (Shared data layer - DI injected)             │
│  ┌────────────┐ ┌───────────┐ ┌─────────────────┐ ┌────────┐  │
│  │ BlockStore │ │ BeaconDB  │ │AvailabilityStore│ │Deposit │  │
│  │            │ │ (StateDB) │ │ (Blob storage)  │ │Store   │  │
│  └────────────┘ └───────────┘ └─────────────────┘ └────────┘  │
└────────────────────────┬──────────────────────────────────────┘
                         │ accessed by
     ┌───────────────────┼───────────────────┐
     │                   │                   │
     ▼                   ▼                   ▼
┌─────────────┐ ┌─────────────────┐ ┌───────────────┐
│ Node API    │ │ Validator       │ │ Blockchain    │
│ Server      │ │ Service         │ │ Service       │
│(REST/HTTP)  │ │                 │ │               │
└─────────────┘ └─────────────────┘ └───────────────┘

Async Background Process:
┌────────────────────────────────────────────────────┐
│            Deposit Monitoring Flow                 │
│  Execution Layer → Deposit Contract → DepositStore │
└────────────────────────────────────────────────────┘
```

**System Flow:**
1. **Service Registry** manages the lifecycle of all services in dependency order
2. **CometBFT Service** drives the entire system through ABCI callbacks
3. **Block Building**: CometBFT → PrepareProposal → Validator Service
4. **Block Processing**: CometBFT → ProcessProposal/FinalizeBlock → Blockchain Service
5. **State Updates**: Blockchain Service → State Processor → Storage Backend
6. **External Queries**: Node API → Storage Backend (independent of consensus flow)

**Key Facts:**
- Service Registry orchestrates startup/shutdown in correct order
- CometBFT is the consensus driver, not a middle layer
- Storage Backend is a passive resource, not an active service
- Node API runs independently, only queries storage
- Execution Engine communicates with external EL clients
- Services must wait for dependencies (e.g., EngineClient blocks until EL is ready)

### Core Components

**beacon/** - Consensus layer implementation
- `blockchain/`: Core service handling block processing, state transitions, fork choice
  - Integrates with execution engine via Engine API
  - Manages block and blob storage
  - Handles deposit processing from execution layer
- `validator/`: Block building and validation
  - Payload building coordination with execution client
  - Blob sidecar creation and bundling
  - Block proposal generation for CometBFT
- `payload-time/`: Timing utilities for payload handling

**state-transition/** - Beacon chain state transitions
- Core state machine implementing Ethereum consensus specs
- Handles validator lifecycle: activation, exits, slashing
- Processes operations: deposits, withdrawals, attestations
- Fork transition logic (Deneb → Electra)
- Custom modifications: validator set cap, churn rules

**node-core/** - Infrastructure and dependency injection
- `components/`: All service providers for DI
  - Each component has a `Provide*` function
  - Dependencies declared via struct tags
- `services/registry/`: Service lifecycle management
  - Start/stop ordering based on dependencies
  - Graceful shutdown handling
- `builder/`: Node assembly and configuration

**execution/** - Execution client integration
- `client/`: Engine API implementation
  - JWT authenticated HTTP client
  - Retry logic with exponential backoff
  - Error classification (fatal vs retryable)
- `deposit/`: Deposit contract monitoring
  - Syncs deposits from execution layer
  - Manages deposit Merkle tree
- Supports all major EL clients via standard Engine API

**storage/** - Multi-layered persistence
- `beacondb/`: State storage with context management
  - Fork-aware state queries
  - Validator and balance lookups
- `blockstore/`: KV store for beacon blocks
  - Indexed by slot number
  - Range query support
- `availabilitystore/`: Blob sidecar storage
  - TTL-based pruning (availability window)
  - Indexed by block root and index
- `depositstore/`: Deposit event tracking
  - Synced from execution layer logs

**consensus-types/** - Core type definitions
- Beacon chain types: blocks, states, validators
- Fork-specific types with version handling
- SSZ serialization for all consensus types
- Generic interfaces for fork compatibility

**da/** - Data availability layer
- `blob/`: Blob and sidecar management
  - KZG commitment verification
  - Blob to sidecar transformation
- `kzg/`: KZG ceremony integration
  - Trusted setup loading
  - Proof generation and verification

### Data Flow Patterns

#### Block Production Flow
1. CometBFT calls `PrepareProposal` when node is proposer
2. Validator service initiates payload building:
   - Sends FCU with payload attributes to execution engine
   - Receives payload ID for tracking
3. Parallel block assembly:
   - Calls `getPayload` to retrieve execution payload
   - Creates blob sidecars from blob transactions
   - Assembles beacon block with payload
4. Returns complete block to CometBFT for proposal

#### Block Processing Flow
1. CometBFT calls `ProcessProposal` with new block
2. Blockchain service validates block structure
3. Blob processor verifies:
   - KZG proofs for all blobs
   - Blob count within limits
4. State processor performs lightweight validation
5. Execution engine validates via `newPayload`
6. Vote to accept/reject returned to CometBFT

#### State Finalization Flow
1. CometBFT calls `FinalizeBlock` after consensus
2. State processor executes full state transition:
   - Process slots up to block slot
   - Apply block operations
   - Update validator balances
3. Storage backends persist:
   - Updated beacon state
   - Block and sidecars
   - State root mappings
4. Post-block FCU updates execution engine head

### Critical Paths

**Block Production (must complete in ~1s):**
1. Initiate payload building immediately
2. Parallel assembly of beacon block components
3. Timeout handling for slow execution clients

**Block Validation (must complete quickly):**
1. Structural validation first (fail fast)
2. Parallel blob verification
3. Execution payload validation last

**State Transitions (deterministic execution):**
1. Slot processing (RANDAO, proposer selection)
2. Block processing (operations in order)
3. Epoch processing (validator updates)

## Key Types and Interfaces

### Core Interfaces
- `ChainSpec` - Chain specification interface
- `BeaconState` - Beacon chain state
- `ExecutionPayload` - Execution layer payload
- `AvailabilityStore` - Blob storage interface
- `DepositStore` - Deposit storage interface
- `BlockStore` - Block storage interface
- `StateProcessor` - State transition processor

### Important Types
- `types.BeaconBlock` - Beacon block structure
- `types.Validator` - Validator information
- `types.Deposit` - Deposit data
- `engine.ExecutionEngine` - Execution engine client
- `payload.PayloadBuilder` - Local payload builder

### Key Abstractions

**StorageBackend Interface:**
```go
type StorageBackend interface {
    AvailabilityStore() AvailabilityStore
    BlockStore() BlockStore
    DepositStore() DepositStore
    StateFromContext(ctx context.Context) BeaconState
}
```

**StateProcessor Interface:**
```go
type StateProcessor interface {
    ProcessSlot(BeaconState) (TransitionResult, error)
    ProcessBlock(BeaconState, BeaconBlock) (BeaconState, error)
    ProcessEpoch(BeaconState) (TransitionResult, error)
}
```

**ExecutionEngine Interface:**
```go
type ExecutionEngine interface {
    NewPayload(ctx, payload, versionedHashes) (PayloadStatus, error)
    ForkchoiceUpdate(ctx, state, attrs) (ForkchoiceResponse, error)
    GetPayload(ctx, payloadID) (ExecutionPayload, error)
}
```

## Integration Patterns

### CometBFT Integration (ABCI++)

BeaconKit integrates with CometBFT through ABCI++ hooks:

```go
// Key integration points in beacon/blockchain/service.go
PrepareProposal(ctx, req) (*ProposalResponse, error)  // Build blocks
ProcessProposal(ctx, req) (*ProcessResponse, error)   // Validate blocks
FinalizeBlock(ctx, req) (*BlockResponse, error)      // Execute state transition
```

**Important:** BeaconKit maps 1 beacon slot = 1 CometBFT height (no missed slots)

### Execution Engine Integration

Communication with execution clients via Engine API:

```go
// Engine API flow for block production
1. forkchoiceUpdate(head, safe, finalized, payloadAttributes) → payloadId
2. getPayload(payloadId) → executionPayload + blobsBundle
3. newPayload(executionPayload) → status
4. forkchoiceUpdate(newHead, safe, finalized) → status
```

**Key Files:**
- `execution/client/client.go` - Main Engine API client
- `execution/client/errors.go` - Error handling and retries
- `execution/pkg/engine/helpers.go` - Request/response helpers

### State Management Integration

State transitions follow a specific pattern:

```go
// State transition pipeline
ctx := state.Context()
state = processSlots(state, targetSlot)
state = processBlock(state, block)
if isEpochEnd(slot) {
    state = processEpoch(state)
}
storage.SetState(ctx, state)
```

**Context Usage:** All state queries use context for fork-awareness

### Storage Layer Integration

Multi-backend storage with clear interfaces:

```go
// Storage access pattern
backend := node.StorageBackend()
state := backend.StateFromContext(ctx)
block := backend.BlockStore().GetBySlot(slot)
blobs := backend.AvailabilityStore().Get(blockRoot)
deposits := backend.DepositStore().GetAll()
```

### Service Lifecycle Integration

Services follow a strict lifecycle pattern:

```go
type Service interface {
    Start(context.Context) error
    Stop() error
    Name() string
}
```

**Startup Order:**
1. Storage backends initialized
2. Core services created via DI
3. Services started in dependency order
4. CometBFT node started last

**Shutdown:** Reverse order with graceful termination

### Configuration Integration

Configuration flows through the system:

```toml
# config.toml structure
[beacon-kit]
  [beacon-kit.engine]
    jwt-secret-path = "path/to/jwt.hex"
    rpc-dial-url = "http://localhost:8551"

  [beacon-kit.payload-builder]
    enabled = true

  [beacon-kit.validator]
    enable-optimistic-payloads = true
```

**Environment Overrides:** CLI flags > env vars > config file > defaults

## Block Height Lifecycle

This section describes how BeaconKit processes blocks using CometBFT's Tendermint consensus algorithm.

### Key Differences from Ethereum
- **Timeout-Based Timing**: No fixed slot duration; block time depends on consensus timeouts and network conditions
- **Sequential Heights**: CometBFT maintains sequential block heights without gaps (unlike Ethereum's slot system)
- **Round-Based Consensus**: Multiple proposers per height if rounds fail
- **Height = Slot**: CometBFT height maps 1:1 to beacon slot number

### Consensus Timing Configuration
BeaconKit enforces minimum timeout values for CometBFT consensus:
- **TimeoutPropose**: 2000ms minimum - Initial proposal timeout
- **TimeoutPrevote**: 2000ms minimum - Prevote collection timeout
- **TimeoutPrecommit**: 2000ms minimum - Precommit collection timeout
- **TimeoutCommit**: 500ms minimum - Post-commit wait time

These are **enforced minimums**, not targets. The cumulative timeouts effectively create a minimum block time, but actual block times depend on:
- Network latency and message propagation
- Number of consensus rounds needed (failed rounds extend time)
- Execution client response times
- Validator participation and vote collection speed

### CometBFT Round-Based Consensus

**Key Concept**: CometBFT uses rounds within each height. If consensus fails, it increments the round (not the height) and selects a new proposer.

```
Height N, Round 0: Proposer A fails/times out
Height N, Round 1: Proposer B tries
Height N, Round 2: Proposer C succeeds → Height N+1
```

### ABCI Method Flow

#### 1. Block Proposal (PrepareProposal)
**When**: Node is selected as proposer for current height/round
**Called on**: Proposer node only
**Can be called**: Multiple times per height (once per round if timeouts occur)

**Steps**:
1. CometBFT calls `PrepareProposal` with height and round info
2. State is reset if this is a subsequent round (handles timeouts)
3. `ValidatorService.BuildBlockAndSidecars()` executes:
   - Sends `forkchoiceUpdate` with attributes to execution client
   - Receives payload ID
   - Calls `getPayload()` to retrieve execution payload
   - `BlobFactory` creates blob sidecars
4. Returns proposed block to CometBFT

#### 2. Block Validation (ProcessProposal)
**When**: Any node (including proposer) receives a proposed block
**Called on**: ALL nodes
**Purpose**: Validate proposal before voting

**Steps**:
1. CometBFT calls `ProcessProposal` with proposed block
2. Resets `finalizeBlockState` in preparation
3. Blockchain Service validates:
   - Block structure and signatures
   - `BlobProcessor` verifies KZG proofs
   - Execution engine validates via `newPayload` (but doesn't commit)
4. Returns ACCEPT or REJECT to CometBFT

#### 3. Consensus Voting Rounds
**After ProcessProposal returns ACCEPT**:

1. **Prevote Phase**:
   - Validators broadcast prevotes for the proposal
   - Wait up to TimeoutPrevote for 2/3+ prevotes

2. **Precommit Phase**:
   - If 2/3+ prevotes received, broadcast precommit
   - Wait up to TimeoutPrecommit for 2/3+ precommits

3. **Commit Decision**:
   - If 2/3+ precommits received, block is decided
   - Wait TimeoutCommit before moving to next height

If any phase fails/times out → new round at same height

#### 4. Block Finalization (FinalizeBlock)
**When**: After consensus achieved (2/3+ precommits)
**Called on**: ALL nodes
**Called**: Exactly once per height (not per round)

**Steps**:
1. CometBFT calls `FinalizeBlock` with decided block
2. State Processor executes state transition:
   - Applies all block operations
   - Updates validator balances and registry
   - Processes epoch transitions if applicable
3. Commits to storage:
   - Beacon state → BeaconDB
   - Block → BlockStore
   - Blob sidecars → AvailabilityStore
4. Sends post-finalization `forkchoiceUpdate` to execution client
5. CometBFT advances to next height

### Important Characteristics

**Height Progression**:
- Heights only increment after successful FinalizeBlock
- Multiple rounds can occur at the same height
- Each round gets a new proposer (deterministic selection)

**State Management**:
- ProcessProposal validates but doesn't commit
- FinalizeBlock performs actual state mutations
- State can be reset between rounds at same height

**Sequential Block Heights**:
- BeaconKit's `ProcessSlots` called for every sequential height
- No gaps in block heights (unlike Ethereum's slot system)
- Failed proposals trigger new rounds at same height, not height skips

### Failure Scenarios

- **Proposer Timeout**: New round with different proposer at same height
- **Invalid Proposal**: Rejected in ProcessProposal, new round begins
- **Insufficient Votes**: Timeout triggers new round
- **Network Partition**: Consensus halts until 2/3+ validators connected
- **Execution Client Issues**: Proposal fails, new round attempted

This architecture ensures continuous block production through CometBFT's robust round-based consensus, maintaining Byzantine fault tolerance with up to 1/3 malicious validators.

### Critical Integration Points

1. **Block Production Timing:**
   - Must complete within slot duration
   - Timeout handling for slow execution clients
   - Parallel preparation of components

2. **State Consistency:**
   - CometBFT app hash = beacon state root
   - Execution payload state root verification
   - Fork choice alignment with execution client

3. **Deposit Bridge:**
   - Monitors execution layer deposit events
   - Maintains deposit Merkle tree
   - Synchronizes with beacon state

4. **Blob Handling:**
   - KZG verification in consensus layer
   - Blob propagation separate from blocks
   - Pruning after availability window

### Component Dependencies

The system uses `cosmossdk.io/depinject` for dependency injection:

```go
// Example component provider
func ProvideBlockchainService(
    in struct {
        depinject.In
        ChainSpec    chain.ChainSpec
        ExecutionEngine ExecutionEngine
        LocalBuilder    LocalBuilder
        StateProcessor  StateProcessor
        StorageBackend  StorageBackend
        TelemetrySink   TelemetrySink
    },
) *Service {
    // Component initialization
}
```

**Dependency Graph:**
- Node → CometBFT Service → Blockchain Service
- Blockchain → State Processor, Execution Engine, Storage
- State Processor → Chain Spec, Execution Engine
- All components → Logger, Metrics, Config

## Storage Architecture

**Multi-Backend Approach:**
1. **CometBFT State**: Application state via IAVL trees
2. **BeaconDB**: Beacon state indexed by StateRoot
3. **BlockStore**: Blocks indexed by slot
4. **AvailabilityStore**: Blobs with pruning
5. **FileDB**: Generic KV store implementation

**State Management:**
- StateDB provides high-level API
- Context-based forking for speculation
- Lazy loading for performance
- Cache layers for hot data

## Error Handling

- Custom error wrapper in `errors/mod.go`
- Fatal vs non-fatal errors with `IsFatal()` checks
- Consistent error wrapping with context
- Detailed error types for different failure modes

## Build Requirements

### Dependencies
- Go 1.23.6+
- Docker (for running EL clients)
- Foundry (for Solidity contracts)
- Make (GNU Make)

### Build Tags
- `bls12381` - BLS cryptography
- `ckzg` - KZG blob commitments
- `test` - Testing utilities
- `e2e` - End-to-end tests
- `simulated` - Simulation tests

### Key Constants
- `RootLength = 32` - Hash tree root length
- Default RPC timeout: 30s
- Default shutdown timeout: 5 minutes
- Block store availability window: configurable via chain spec

## Development Considerations

### Key Design Patterns

1. **Interface Segregation**: Small, focused interfaces for each concern
2. **Dependency Injection**: All wiring via DI container, no global state
3. **Context Propagation**: Request-scoped values and cancellation
4. **Error Wrapping**: Detailed context preservation with error chains
5. **Metrics-First**: Every operation instrumented with Prometheus
6. **Defensive Validation**: Validate early, validate often
7. **Fork Abstraction**: Generic types handle fork differences

### Testing Approach

- Unit tests alongside code files (`*_test.go`)
- Simulated tests for chaos/fork scenarios in `testing/`
- E2E tests with real execution clients in `testing/e2e/`
- Kurtosis for multi-node testing scenarios
