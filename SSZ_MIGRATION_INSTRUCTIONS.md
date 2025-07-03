# SSZ Migration Plan

## Overview
Migration from karalabe/ssz to fastssz for BeaconKit consensus types.

### Current Status (As of Phase 9 Completion)
- ✅ **All types that marshal with SSZ now have fastssz methods**
- ✅ **Core types migrated from karalabe/ssz to fastssz**
- ⚠️ **Some packages still use karalabe/ssz (da/types, test files)**
- ⚠️ **Cannot use sszgen until all karalabe/ssz dependencies removed**

## Status Update

### ✅ Completed Phases
- **Phase 1**: Dual interface support added to core types
- **Phase 2**: Generated fastssz code for simple types (Fork, ForkData, AttestationData, etc.)
- **Phase 3**: Added fastssz support to math types (U64, Gwei, etc.)
- **Phase 4**: Added fastssz methods to ExecutionPayload and ExecutionPayloadHeader
- **Phase 5**: Added fastssz support to BeaconState
- **Phase 6**: Added fastssz methods to ExecutionRequests
- **Phase 7**: Added fastssz methods to collection types
- **Phase 8**: Added fastssz support to SignedBeaconBlock
- **Phase 9**: Migrated core types from karalabe/ssz to fastssz (consensus-types, engine-primitives, primitives/common)

### 🚧 Current Phase: Phase 10 - Complete Migration of Remaining Types

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

### Phase 2 Migrations
- [x] SyncAggregate - Manual fastssz implementation with karalabe/ssz compatibility
- [x] WithdrawalCredentials - Inherits fastssz from bytes.B32
- [x] Added fastssz methods to bytes.B32 and bytes.B96 types

## Remaining Work

### Simple Types (Can Migrate Independently)
- [x] WithdrawalCredentials - Type alias to common.Bytes32 (inherits fastssz from B32) ✅

### Types Blocked by UnusedType Dependency
These types are currently type aliases to common.UnusedType, which still uses karalabe/ssz:
- [x] ProposerSlashing - Type alias to UnusedType ✅
- [x] AttesterSlashing - Type alias to UnusedType ✅
- [x] VoluntaryExit - Type alias to UnusedType ✅
- [x] BLSToExecutionChange - Type alias to UnusedType ✅
- [x] Attestation - Type alias to UnusedType (affects Attestations[] in BeaconBlockBody) ✅

### Types Blocked by Other Dependencies
- [x] DepositRequest - Type alias to Deposit ✅

### Types Blocked by Incomplete Primitive Support
These types need additional fastssz support:
- [x] WithdrawalRequest - Has all primitive support ✅
- [x] ConsolidationRequest - Has all primitive support ✅
- [x] math.Gwei (U64) - Now has full fastssz support ✅
- [x] bytes.B48 - Now has full fastssz support ✅
- [x] common.ExecutionAddress - Now has full fastssz support ✅

### Types Still Using karalabe/ssz
- [x] SyncAggregate - Migrated to fastssz ✅
- [x] ExecutionRequests - Migrated to fastssz ✅

### Types with Mixed SSZ Support (Already have fastssz methods)
- [x] Validator - Migrated to fastssz ✅
- [x] Eth1Data - Migrated to fastssz ✅
- [x] Deposit - Migrated to fastssz ✅
- [x] BeaconBlockHeader - Migrated to fastssz ✅
- [x] PendingPartialWithdrawal - Has all fastssz methods needed ✅

### Fork-Specific Types (Need Manual fastssz Implementation)
These types have fork-specific serialization logic that changes based on fork version:

- [x] **ExecutionPayload** - Migrated to fastssz ✅
- [x] **ExecutionPayloadHeader** - Migrated to fastssz ✅
- [x] **BeaconBlockBody** - Migrated to fastssz ✅
- [x] **BeaconState** - Migrated to fastssz ✅

### Additional Complex Types
- [x] BeaconBlock - Migrated to fastssz ✅
- [ ] BlobSidecar - Critical DA type, still uses karalabe/ssz

### Complex Dependency Chains

#### BeaconBlock Chain (Must Migrate Together)
```
BeaconBlock ✅
└── BeaconBlockBody ✅
    ├── ExecutionPayload ✅
    ├── ExecutionRequests ✅
    ├── Eth1Data ✅
    ├── Deposits[] ✅
    ├── ProposerSlashings[] ✅
    ├── AttesterSlashings[] ✅
    ├── Attestations[] ✅
    ├── VoluntaryExits[] ✅
    ├── SyncAggregate ✅
    └── BlobKzgCommitments ✅
```

#### BeaconState Chain
```
BeaconState ✅
├── Fork ✅
├── BeaconBlockHeader ✅
├── Validators[] ✅
├── Eth1Data ✅
├── ExecutionPayloadHeader ✅
└── PendingPartialWithdrawals[] ✅
```

#### Signed Types (Depend on Base Types)
- [x] SignedBeaconBlock ✅
- [x] SignedBeaconBlockHeader ✅

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
- [x] **Phase 1**: Migrate UnusedType and its aliases from karalabe/ssz to fastssz ✅
- [x] **Phase 2**: Migrate truly independent types ✅
- [x] **Phase 3**: Add SSZ methods to math.U64/Gwei to unblock ✅
- [x] **Phase 4**: Complete migration of mixed support types ✅
- [x] **Phase 5**: Implement manual fastssz for fork-specific types ✅
- [x] **Phase 6**: Migrate complex chains ✅
- [x] **Phase 7**: Migrate all collection types and SignedBeaconBlock ✅
- [x] **Phase 8**: Complete fastssz support for remaining types ✅

## Technical Notes
- Types used by karalabe/ssz types must keep karalabe methods
- Fork-specific logic requires manual HashTreeRootWith implementation
- Complex types with many dependencies should be migrated together to avoid build breaks
- **Critical Discovery**: BeaconBlockBody must be migrated to fastssz BEFORE UnusedType can be migrated
- UnusedType is a simple uint8 that enforces zero value - requires manual fastssz implementation
- **Note**: common.Bytes32 (bytes.B32) has custom SSZ methods independent of both libraries
- **Collection Types**: Types like Deposits, Validators, Attestations etc. are slice wrappers
- **Interface Conflicts**: Some types implement both karalabe/ssz and fastssz methods with conflicting signatures

## Current Migration Status
- ✅ **Phase 0-8 Complete**: All core functionality migrated
- ✅ **Phase 9 Complete**: Core types migrated from karalabe/ssz to fastssz
  - Removed all DefineSSZ methods from migrated types
  - Updated SizeSSZ signatures from `(siz *ssz.Sizer, fixed bool) uint32` to `() int`
  - Replaced ssz.HashSequential/HashConcurrent with fastssz hasher pool
  - Migrated BeaconBlockBody with complex fork-specific logic
  - Migrated BeaconState with proper field handling
  - Migrated all collection types (Deposits, Validators, etc.)
  - Migrated engine-primitives (Withdrawals, Transactions)
  - Migrated primitives/common/UnusedType

## Phase 9: Complete Removal of karalabe/ssz Dependency 🚀 ✅ COMPLETED

### Summary of Changes
- Removed all DefineSSZ methods from migrated types
- Updated SizeSSZ signatures from `(siz *ssz.Sizer, fixed bool) uint32` to `() int`
- Replaced ssz.HashSequential/HashConcurrent with fastssz hasher pool
- Migrated BeaconBlockBody with complex fork-specific logic
- Migrated BeaconState with proper field handling
- Migrated all collection types (Deposits, Validators, etc.)
- Migrated engine-primitives (Withdrawals, Transactions)
- Migrated primitives/common/UnusedType

## Phase 10: Complete Migration of Remaining Types 🚧

### Step 1: Migrate Remaining Types
1. **da/types** - Migrate BlobSidecar and Sidecars
2. **Other packages** - Find and migrate any remaining types

### Step 2: Replace All karalabe/ssz Usage with fastssz
1. **Replace ssz.EncodeToBytes calls**
   - Find all `ssz.EncodeToBytes(buf, obj)` calls
   - Replace with `obj.MarshalSSZ()` or `obj.MarshalSSZTo(buf)`

2. **Replace ssz.DecodeFromBytes calls**
   - Find all `ssz.DecodeFromBytes(buf, obj)` calls
   - Replace with `obj.UnmarshalSSZ(buf)`

3. **Replace hash functions**
   - Replace `ssz.HashSequential(obj)` with:
     ```go
     root, _ := obj.HashTreeRoot()
     return common.Root(root)
     ```
   - Replace `ssz.HashConcurrent(obj)` similarly

4. **Update ssz.Size() calls**
   - Replace `ssz.Size(obj)` with `obj.SizeSSZ()`

### Step 3: Remove karalabe/ssz Dependency
1. **Update go.mod:**
   ```bash
   go get -u github.com/karalabe/ssz@none
   ```

2. **Update all imports**

3. **Run tests to ensure everything works**

### Step 4: Run sszgen for All Types
Once karalabe/ssz is removed, we can run sszgen for most types without conflicts.

### Testing Strategy
1. **Create compatibility tests:**
   - Ensure fastssz produces same serialization as karalabe/ssz
   - Test hash tree roots match
   - Test round-trip serialization

2. **Benchmark before and after:**
   - Measure serialization performance
   - Measure deserialization performance
   - Measure hash tree root computation

3. **Integration tests:**
   - Run full consensus tests
   - Test with real chain data
   - Verify no regressions

### Risk Mitigation
1. **Create a migration branch**
2. **Make incremental commits** for easy rollback
3. **Run extensive tests** after each major change
4. **Consider a phased rollout:**
   - Migrate non-critical types first
   - Test in devnet before mainnet
   - Have rollback plan ready