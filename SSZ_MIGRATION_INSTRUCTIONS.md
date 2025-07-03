# SSZ Migration Plan

## Overview
Migration from karalabe/ssz to fastssz for BeaconKit consensus types.

### Current Status (As of Phase 8 Completion)
- ✅ **All types that marshal with SSZ now have fastssz methods**
- ✅ **Dual interface compatibility maintained** - both karalabe/ssz and fastssz work
- ⚠️ **Still using karalabe/ssz for actual encoding/decoding in many places**
- ⚠️ **51 files still import karalabe/ssz**
- ⚠️ **Cannot use sszgen effectively due to method signature conflicts**

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
- [ ] ProposerSlashing - Type alias to UnusedType
- [ ] AttesterSlashing - Type alias to UnusedType
- [ ] VoluntaryExit - Type alias to UnusedType
- [ ] BLSToExecutionChange - Type alias to UnusedType
- [ ] Attestation - Type alias to UnusedType (affects Attestations[] in BeaconBlockBody)

### Types Blocked by Other Dependencies
- [ ] DepositRequest - Type alias to Deposit (must migrate Deposit first)

### Types Blocked by Incomplete Primitive Support
These types need additional fastssz support:
- [x] WithdrawalRequest - Has all primitive support but blocked by ExecutionRequests using karalabe/ssz ⚠️
- [x] ConsolidationRequest - Has all primitive support but blocked by ExecutionRequests using karalabe/ssz ⚠️
- [x] math.Gwei (U64) - Now has full fastssz support ✅
- [x] bytes.B48 - Now has full fastssz support ✅
- [x] common.ExecutionAddress - Now has full fastssz support ✅

### Types Still Using karalabe/ssz
- [x] SyncAggregate - Migrated to fastssz with manual implementation ✅
- [ ] ExecutionRequests - Needs migration to fastssz (currently blocks WithdrawalRequest/ConsolidationRequest from full migration)

### Types with Mixed SSZ Support (Already have fastssz methods)
- [x] Validator - Has HashTreeRootWith, ready for fastssz but needs manual migration ⚠️
- [x] Eth1Data - Has HashTreeRootWith, ready for fastssz but needs manual migration ⚠️
- [x] Deposit - Has HashTreeRootWith, ready for fastssz but needs manual migration ⚠️
- [x] BeaconBlockHeader - Has HashTreeRootWith, ready for fastssz but needs manual migration ⚠️
- [x] PendingPartialWithdrawal - Has all fastssz methods needed (HashTreeRootWith, GetTree, MarshalSSZTo) ✅

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
- [x] **Phase 2**: Migrate truly independent types: ✅
  - WithdrawalCredentials (inherits fastssz from B32)
  - SyncAggregate (manual fastssz implementation)
- [x] **Phase 3**: Add SSZ methods to math.U64/Gwei to unblock: ✅
  - WithdrawalRequest, ConsolidationRequest
  - ExecutionRequests (depends on the above)
  - **Status**: U64/Gwei now has full fastssz support (HashTreeRootWith, GetTree)
  - **Additional Requirements Found**: 
    - bytes.B48 needs fastssz support for BLSPubkey fields
    - common.ExecutionAddress needs fastssz support
- [x] **Phase 4**: Complete migration of mixed support types (including Deposit to unblock DepositRequest) ✅
- [x] **Phase 5**: Implement manual fastssz for fork-specific types (ExecutionPayload, ExecutionPayloadHeader, BeaconState) ✅
- [x] **Phase 6**: Migrate complex chains (BeaconBlock, BeaconState) and BlobSidecar once all dependencies ready ✅
- [x] **Phase 7**: Migrate all collection types and SignedBeaconBlock ✅
- [x] **Phase 8**: Complete fastssz support for remaining types (Deposit, Eth1Data, BeaconBlockHeader, Validator) ✅

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
- **Interface Conflicts**: Some types (Deposit, Eth1Data) implement both karalabe/ssz and fastssz methods with conflicting signatures:
  - `SizeSSZ()` returns `int` in fastssz but `uint32` in karalabe/ssz
  - `HashTreeRoot()` returns `([32]byte, error)` in fastssz but `common.Root` in karalabe/ssz
  - These types need manual migration similar to BeaconBlockBody to maintain dual compatibility

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
- ✅ **Phase 2 Complete**: WithdrawalCredentials and SyncAggregate migrated to fastssz
  - WithdrawalCredentials automatically works with fastssz through B32 type
  - SyncAggregate has manual fastssz implementation maintaining karalabe/ssz compatibility
  - Added fastssz methods to bytes.B32 and bytes.B96 to support type aliases
- ✅ **Phase 3 Complete**: Added fastssz methods to math.U64/Gwei
  - U64 now has HashTreeRootWith and GetTree methods
  - All type aliases (Gwei, Slot, ValidatorIndex, etc.) inherit fastssz support
  - Tests confirm compatibility with existing HashTreeRoot implementation
  - Discovered additional dependencies: B48 and ExecutionAddress need fastssz support
- ✅ **Phase 3.5 Complete**: Added fastssz support to remaining primitives
  - bytes.B48 now has full fastssz support (MarshalSSZTo, UnmarshalSSZ, HashTreeRootWith, GetTree)
  - common.ExecutionAddress now has full fastssz support
  - WithdrawalRequest and ConsolidationRequest are ready for fastssz but blocked by ExecutionRequests
- ✅ **Phase 4 Complete**: Completed work on types with mixed SSZ support
  - All types (Validator, Eth1Data, Deposit, BeaconBlockHeader) have the minimal fastssz methods needed
  - PendingPartialWithdrawal already has complete fastssz support (HashTreeRootWith, GetTree, MarshalSSZTo)
  - DepositRequest inherits from Deposit and is ready
  - Most types need manual migration to handle dual interface compatibility
- ✅ **Phase 5 Complete**: Added fastssz methods to fork-specific types
  - ✅ ExecutionPayload: Added fastssz methods (UnmarshalSSZ, SizeSSZFastSSZ, MarshalSSZTo)
  - ✅ ExecutionPayloadHeader: Added fastssz methods (UnmarshalSSZ, SizeSSZFastSSZ, MarshalSSZTo)
  - ✅ BeaconState: Added fastssz methods; already had fork-specific logic for PendingPartialWithdrawals (Electra+ only)
  - All three types now have complete fastssz support while maintaining backward compatibility
- ✅ **Phase 6 Complete**: Complex dependency chains and critical types
  - ✅ ExecutionRequests: Added full fastssz implementation with proper offset-based encoding
  - ✅ WithdrawalRequest: Generated fastssz code using sszgen
  - ✅ ConsolidationRequest: Generated fastssz code using sszgen
  - ✅ BeaconBlock: Added fastssz methods (MarshalSSZTo, UnmarshalSSZ, SizeSSZFastSSZ, HashTreeRootWith)
  - ✅ SignedBeaconBlockHeader: Added fastssz support
  - ✅ BlobSidecar: Added fastssz support
- ✅ **Phase 7 Complete**: Migrated all slice/collection types and SignedBeaconBlock
  - ✅ Collection types: Deposits, Validators, Attestations, ProposerSlashings, AttesterSlashings, VoluntaryExits, BLSToExecutionChanges
  - ✅ SignedBeaconBlock: Added full fastssz support with dynamic object handling
  - All collection types now have HashTreeRootWith and GetTree methods
- ✅ **Phase 8 Complete**: Added UnmarshalSSZ methods to complete fastssz support
  - ✅ Deposit: Added UnmarshalSSZ and SizeSSZFastSSZ methods
  - ✅ Eth1Data: Added UnmarshalSSZ and SizeSSZFastSSZ methods
  - ✅ BeaconBlockHeader: Added UnmarshalSSZ and SizeSSZFastSSZ methods
  - ✅ Validator: Added UnmarshalSSZ and SizeSSZFastSSZ methods
  - All types now have complete fastssz support while maintaining karalabe/ssz compatibility
- **Note**: All types maintain temporary karalabe/ssz compatibility stubs until full migration is complete

## Phase 9: Complete Removal of karalabe/ssz Dependency 🚀

### Overview
This is the final major phase to completely remove karalabe/ssz and migrate fully to fastssz. Once complete, we can use sszgen for most types without conflicts.

### Step 1: Replace All karalabe/ssz Usage with fastssz
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

### Step 2: Remove karalabe/ssz Methods from Types
1. **Remove from each type:**
   - `DefineSSZ(codec *ssz.Codec)` method
   - `SizeSSZ(siz *ssz.Sizer) uint32` method (keep the fastssz version)
   - Any karalabe/ssz specific validation

2. **Rename temporary methods:**
   - `SizeSSZFastSSZ() int` → `SizeSSZ() int`
   - `HashTreeRootCommon() common.Root` → Keep as wrapper if needed

### Step 3: Update Interfaces and Constraints
1. **Update constraints/ssz.go:**
   - Remove dependency on `ssz.Object`
   - Update `SSZUnmarshaler` interface to not embed `ssz.Object`
   - Create new interface that matches fastssz requirements

2. **Update compile-time assertions:**
   - Remove `_ ssz.StaticObject = (*Type)(nil)`
   - Remove `_ ssz.DynamicObject = (*Type)(nil)`
   - Add fastssz interface assertions if needed

### Step 4: Run sszgen for All Types
Once karalabe/ssz is removed, we can run sszgen for most types:

1. **Types ready for sszgen** (currently have manual implementations):
   - Validator
   - BeaconBlockHeader
   - Eth1Data
   - Deposit
   - SyncAggregate
   - PendingPartialWithdrawal
   - ExecutionPayload (if no fork logic needed)
   - ExecutionPayloadHeader (if no fork logic needed)
   - BeaconBlock
   - BeaconBlockBody (if fork logic can be handled)
   - BeaconState (if fork logic can be handled)

2. **Types that may need manual implementation:**
   - Types with complex fork-specific logic
   - Types with special validation requirements
   - Collection types (may keep manual implementations)

### Step 5: Cleanup and Optimization
1. **Remove karalabe/ssz dependency:**
   ```bash
   go get -u github.com/berachain/karalabe-ssz@none
   ```

2. **Update go.mod and go.sum**

3. **Run tests to ensure everything works**

4. **Performance optimization:**
   - Profile SSZ operations
   - Optimize hot paths
   - Consider using object pools for frequently marshaled types

### Migration Order (Critical Path)
To avoid breaking the build, migrate in this order:

1. **Update all direct usage** (Step 1)
   - Start with leaf types that don't have dependencies
   - Work up to complex types

2. **Update interfaces** (Step 3)
   - This may temporarily break compilation
   - Fix all compile errors before proceeding

3. **Remove karalabe methods** (Step 2)
   - Do this after interfaces are updated
   - Types can be done incrementally

4. **Run sszgen** (Step 4)
   - Start with simple types
   - Test each type after generation

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
