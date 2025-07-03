# SSZ Migration Plan

## Overview
Migration from karalabe/ssz to fastssz for BeaconKit consensus types.

## Completed Migrations ✅

### Fully Migrated to fastssz
- [x] AttestationData
- [x] SlashingInfo
- [x] SigningData
- [x] DepositMessage
- [x] ForkData
- [x] Fork

### Foundation Work
- [x] Created concrete Versionable struct for SSZ compatibility
- [x] All types updated to use concrete Versionable

## Remaining Work

### Simple Types (Can Migrate Independently)
- [ ] WithdrawalCredentials - Type alias to common.Bytes32 (has custom SSZ methods)

### Types Blocked by UnusedType Dependency
These types are currently type aliases to common.UnusedType, which still uses karalabe/ssz:
- [ ] ProposerSlashing - Type alias to UnusedType
- [ ] AttesterSlashing - Type alias to UnusedType
- [ ] VoluntaryExit - Type alias to UnusedType
- [ ] BLSToExecutionChange - Type alias to UnusedType
- [ ] Attestation - Type alias to UnusedType (affects Attestations[] in BeaconBlockBody)

### Types Blocked by Other Dependencies
- [ ] DepositRequest - Type alias to Deposit (must migrate Deposit first)

### Types Blocked by Incomplete Primitive Support
These types use math.Gwei (U64) which lacks MarshalSSZ/UnmarshalSSZ methods:
- [ ] WithdrawalRequest - Uses ExecutionAddress, BLSPubkey, and math.Gwei
- [ ] ConsolidationRequest - Uses ExecutionAddress and BLSPubkey fields

### Types Still Using karalabe/ssz
- [ ] SyncAggregate - No dependencies, used in BeaconBlockBody
- [ ] ExecutionRequests - Blocked by WithdrawalRequest/ConsolidationRequest which need Gwei support

### Types with Mixed SSZ Support (Already have fastssz methods)
- [ ] Validator - Has HashTreeRootWith
- [ ] Eth1Data - Has HashTreeRootWith
- [ ] Deposit - Has HashTreeRootWith
- [ ] BeaconBlockHeader - Has HashTreeRootWith
- [ ] PendingPartialWithdrawal - Has HashTreeRootWith

### Fork-Specific Types (Need Manual fastssz Implementation)
These types have fork-specific serialization logic that changes based on fork version:

- [ ] **ExecutionPayload** - Withdrawals nil check for pre-Capella
- [ ] **ExecutionPayloadHeader** - Same withdrawals handling
- [ ] **BeaconBlockBody** - ExecutionRequests only for Electra+
- [ ] **BeaconState** - PendingPartialWithdrawals only for Electra+

### Additional Complex Types
- [ ] BeaconBlock - Depends on BeaconBlockBody
- [ ] BlobSidecar - Critical DA type, uses karalabe/ssz

### Complex Dependency Chains

#### BeaconBlock Chain (Must Migrate Together)
```
BeaconBlock
└── BeaconBlockBody
    ├── ExecutionPayload (fork-specific)
    ├── ExecutionRequests (Electra+, uses karalabe/ssz)
    ├── Eth1Data (has partial fastssz)
    ├── Deposits[] (Deposit has partial fastssz)
    ├── ProposerSlashings[] (blocked by UnusedType)
    ├── AttesterSlashings[] (blocked by UnusedType)
    ├── Attestations[] (blocked by UnusedType)
    ├── VoluntaryExits[] (blocked by UnusedType)
    ├── SyncAggregate (uses karalabe/ssz)
    └── BlobKzgCommitments ([]eip4844.KZGCommitment - not a custom type)
```

#### BeaconState Chain
```
BeaconState (fork-specific)
├── Fork (✅ already migrated to fastssz)
├── BeaconBlockHeader (has partial fastssz)
├── Validators[] (Validator has partial fastssz)
├── Eth1Data (has partial fastssz)
├── ExecutionPayloadHeader (fork-specific)
└── PendingPartialWithdrawals[] (Electra+, has partial fastssz)
```

#### Signed Types (Depend on Base Types)
- [ ] SignedBeaconBlock → BeaconBlock
- [ ] SignedBeaconBlockHeader → BeaconBlockHeader

## Migration Strategy

### SSZ Code Generation Approach
- **Primary Method**: Use `sszgen` to auto-generate serialization code for all types
- **Manual Implementation**: Only when sszgen cannot handle fork-specific logic:
  1. Run sszgen to generate initial code
  2. Rename generated file (remove `_sszgen` suffix)
  3. Make minimal manual changes for fork-specific logic
  4. Types requiring manual work: ExecutionPayload, ExecutionPayloadHeader, BeaconBlockBody, BeaconState

### Migration Phases
- [x] **Phase 0**: Migrate BeaconBlockBody to fastssz (critical blocker) ✅
  - BeaconBlockBody now supports both karalabe/ssz and fastssz interfaces
  - Manual implementation added to handle fork-specific logic (ExecutionRequests for Electra+)
  - Tests confirm fastssz HashTreeRootWith produces same results as karalabe/ssz
  - BeaconBlockBody must keep karalabe/ssz methods until BeaconBlock is migrated
- [x] **Phase 1**: Migrate UnusedType and its aliases from karalabe/ssz to fastssz ✅
  - UnusedType now has full fastssz support: MarshalSSZTo, UnmarshalSSZ, HashTreeRootWith, etc.
  - All aliases (ProposerSlashing, AttesterSlashing, VoluntaryExit, BLSToExecutionChange, Attestation) automatically inherit fastssz methods
  - Tests confirm all aliases work correctly with fastssz and maintain zero-value enforcement
- [ ] **Phase 2**: Migrate truly independent types:
  - WithdrawalCredentials (simple type alias)
  - SyncAggregate (no type dependencies)
- [ ] **Phase 3**: Add SSZ methods to math.U64/Gwei to unblock:
  - WithdrawalRequest, ConsolidationRequest
  - ExecutionRequests (depends on the above)
- [ ] **Phase 4**: Complete migration of mixed support types (including Deposit to unblock DepositRequest)
- [ ] **Phase 5**: Implement manual fastssz for fork-specific types (ExecutionPayload, ExecutionPayloadHeader, BeaconState)
- [ ] **Phase 6**: Migrate complex chains (BeaconBlock, BeaconState) and BlobSidecar once all dependencies ready

## Technical Notes
- Types used by karalabe/ssz types must keep karalabe methods
- Fork-specific logic requires manual HashTreeRootWith implementation
- Complex types with many dependencies should be migrated together to avoid build breaks
- **Critical Discovery**: BeaconBlockBody must be migrated to fastssz BEFORE UnusedType can be migrated
  - BeaconBlockBody embeds ProposerSlashing, AttesterSlashing, etc. fields
  - These are type aliases to UnusedType
  - BeaconBlockBody's DefineSSZ method expects these types to implement karalabe/ssz.StaticObject
  - If UnusedType is migrated to fastssz first, it breaks BeaconBlockBody compilation
- UnusedType is a simple uint8 that enforces zero value - requires manual fastssz implementation (sszgen doesn't support type aliases)
- **Note**: common.Bytes32 (bytes.B32) has custom SSZ methods independent of both libraries - no migration needed
- **Gwei/U64 Issue**: math.Gwei (alias to math.U64) only has HashTreeRoot() but lacks MarshalSSZ/UnmarshalSSZ methods, blocking migration of WithdrawalRequest, ConsolidationRequest, and ExecutionRequests
- **Collection Types**: Types like Deposits, Validators, Attestations etc. are slice wrappers that will automatically work once their element types are migrated

## Current Migration Status
- ✅ **Phase 0 Complete**: BeaconBlockBody now has fastssz support
  - Added manual fastssz methods: MarshalSSZTo, HashTreeRootWith, GetTree, SizeSSZFastSSZ
  - Maintains backward compatibility with karalabe/ssz for BeaconBlock dependency
  - Fork-specific logic for ExecutionRequests (Electra+) properly handled
  - All tests passing - fastssz produces identical hash roots to karalabe/ssz
- ✅ **Phase 1 Complete**: UnusedType and all aliases migrated to fastssz
  - UnusedType has manual fastssz implementation with zero-value enforcement
  - Type aliases automatically inherit all fastssz methods from UnusedType
  - BeaconBlockBody can now use these types with either karalabe/ssz or fastssz
- **Next**: Phase 2 - Migrate WithdrawalCredentials and SyncAggregate
