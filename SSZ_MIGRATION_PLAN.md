# SSZ Migration Plan: karalabe/ssz → fastssz

## Overview
The migration involves completing the transition from karalabe/ssz to fastssz implementations across the codebase. Most files already use a hybrid approach with karalabe/ssz for marshaling/unmarshaling and fastssz for hash tree root operations. We'll complete this migration by replacing the remaining karalabe/ssz methods with fastssz equivalents in-place.

## Current State
- **Hybrid Implementation**: Most consensus types already import both libraries
- **karalabe/ssz**: Used for `MarshalSSZ()`, `UnmarshalSSZ()`, `DefineSSZ()`, and size calculations
- **fastssz**: Already used for `HashTreeRootWith()` and `GetTree()` methods
- **Migration Focus**: Replace karalabe/ssz marshal/unmarshal with fastssz equivalents

## Migration Approach: In-Place Migration

### Implementation Pattern
Since files already contain both karalabe/ssz and fastssz imports with fastssz methods implemented, we'll complete the migration by:
1. **Replace karalabe/ssz methods** with fastssz implementations:
   - `MarshalSSZ()` - implement using fastssz
   - `UnmarshalSSZ()` - implement using fastssz
   - `HashTreeRoot()` - update to use fastssz hasher
   - Remove `DefineSSZ()` method (not needed for fastssz)
2. **Remove karalabe/ssz import** and related code
3. **Keep existing fastssz methods** that are already implemented

### Testing Strategy
```bash
# Create compatibility tests BEFORE making changes
# to ensure identical behavior

# Test current implementation
go test -tags=test ./consensus-types/types/...

# After migration, run same tests
go test -tags=test ./consensus-types/types/...

# Run compatibility tests if created
go test -tags="test,ssz_compat" ./consensus-types/types/...
```

## Migration Status

## Phase 1: Simple Static Types (Low Risk)

### 1. Fork Type ✅ COMPLETED
**File**: `/consensus-types/types/fork.go`
- **Fields**: PreviousVersion (4 bytes), CurrentVersion (4 bytes), Epoch (8 bytes) = 16 bytes total
- **Migration Status**: ✅ Completed - Now uses fastssz for marshaling/unmarshaling
- **Implementation Notes**:
  - `MarshalSSZ()` now uses fastssz via `MarshalSSZTo()`
  - `UnmarshalSSZ()` implemented using fastssz
  - `HashTreeRoot()` now uses fastssz hasher
  - `DefineSSZ()` kept for backward compatibility with generic unmarshal
  - `SizeSSZ()` kept as required by SSZMarshallable interface
  - All existing tests pass with identical serialization output
- [x] Create compatibility test BEFORE changes
- [x] Replace `MarshalSSZ()` to use fastssz (currently uses karalabe)
- [x] Add `UnmarshalSSZ()` method using fastssz
- [x] Update `HashTreeRoot()` to use fastssz (currently uses karalabe)
- [x] Remove `DefineSSZ()` and `SizeSSZ()` methods (kept for compatibility)
- [x] Remove karalabe/ssz import (kept for interface compatibility)
- [x] Keep existing fastssz methods (`MarshalSSZTo`, `HashTreeRootWith`, `GetTree`)
- [x] Run all tests to ensure nothing breaks
- [x] ✅ All tests pass before proceeding

### 2. ForkData Type
**File**: `/consensus-types/types/fork_data.go`
- **Fields**: CurrentVersion (4 bytes), GenesisValidatorsRoot (32 bytes) = 36 bytes total
- **Current State**: Only uses karalabe/ssz, NO fastssz methods yet
- **karalabe/ssz Methods**: `SizeSSZ()`, `DefineSSZ()`, `MarshalSSZ()`, `MarshalSSZTo()`, `HashTreeRoot()`
- **Missing fastssz Methods**: Need to add all fastssz methods
- [ ] Create compatibility test BEFORE changes
- [ ] Replace `MarshalSSZ()` to use fastssz
- [ ] Update `MarshalSSZTo()` to use fastssz (currently calls MarshalSSZ)
- [ ] Add `UnmarshalSSZ()` method using fastssz
- [ ] Update `HashTreeRoot()` to use fastssz
- [ ] Add `HashTreeRootWith()` fastssz method
- [ ] Add `GetTree()` fastssz method
- [ ] Remove `DefineSSZ()` and `SizeSSZ()` methods
- [ ] Remove karalabe/ssz import and add fastssz import
- [ ] Test `ComputeDomain()` and `ComputeRandaoSigningRoot()` still work correctly
- [ ] Run all tests to ensure nothing breaks
- [ ] ✅ All tests pass before proceeding

### 3. Eth1Data Type
**File**: `/consensus-types/types/eth1data.go`
- **Fields**: DepositRoot (32 bytes), DepositCount (8 bytes), BlockHash (32 bytes) = 72 bytes total
- **Current State**: Hybrid implementation with both karalabe/ssz and fastssz
- **karalabe/ssz Methods**: `SizeSSZ()`, `DefineSSZ()`, `MarshalSSZ()`, `HashTreeRoot()`
- **fastssz Methods**: `MarshalSSZTo()`, `HashTreeRootWith()`, `GetTree()`
- [ ] Create compatibility test BEFORE changes
- [ ] Replace `MarshalSSZ()` to use fastssz (currently uses karalabe)
- [ ] Update `MarshalSSZTo()` to use fastssz directly (currently calls MarshalSSZ)
- [ ] Add `UnmarshalSSZ()` method using fastssz
- [ ] Update `HashTreeRoot()` to use fastssz (currently uses karalabe)
- [ ] Remove `DefineSSZ()` and `SizeSSZ()` methods
- [ ] Remove karalabe/ssz import
- [ ] Keep existing fastssz methods (`HashTreeRootWith`, `GetTree`)
- [ ] Test `GetDepositCount()` method still works
- [ ] Run all tests to ensure nothing breaks
- [ ] ✅ All tests pass before proceeding

### 4. AttestationData Type
**File**: `/consensus-types/types/attestation_data.go`
- [ ] Create compatibility test BEFORE changes
- [ ] Replace `MarshalSSZ()` to use fastssz
- [ ] Replace `UnmarshalSSZ()` to use fastssz
- [ ] Update `HashTreeRoot()` to use fastssz
- [ ] Remove `DefineSSZ()` method
- [ ] Remove karalabe/ssz import
- [ ] Note: Already has `HashTreeRootWith()` fastssz method
- [ ] Run all tests to ensure nothing breaks
- [ ] ✅ All tests pass before proceeding

### 5. SlashingInfo Type
**File**: `/consensus-types/types/slashing_info.go`
- [ ] Create compatibility test BEFORE changes
- [ ] Replace `MarshalSSZ()` - handle RoundSlashed, SlashingThreshold
- [ ] Replace `UnmarshalSSZ()` for two uint64 fields
- [ ] Update `HashTreeRoot()` to use fastssz
- [ ] Remove `DefineSSZ()` method
- [ ] Remove karalabe/ssz import
- [ ] Note: Already has `HashTreeRootWith()` fastssz method
- [ ] Test boundary values for uint64 fields
- [ ] Run all tests to ensure nothing breaks
- [ ] ✅ All tests pass before proceeding

## Phase 2: Dynamic Simple Types (Medium Risk)

### 6. Deposit Type
**File**: `/consensus-types/types/deposit.go`
- [ ] Create compatibility test BEFORE changes
- [ ] Replace `MarshalSSZ()` - handle Proof, Data
- [ ] Replace `UnmarshalSSZ()` with dynamic data handling
- [ ] Update `HashTreeRoot()` to use fastssz
- [ ] Remove `DefineSSZ()` method
- [ ] Remove karalabe/ssz import
- [ ] Note: Already has `HashTreeRootWith()` fastssz method
- [ ] Handle `DepositData` embedded type
- [ ] Test proof validation
- [ ] Run all tests to ensure nothing breaks
- [ ] ✅ All tests pass before proceeding

### 7. Validator Type
**File**: `/consensus-types/types/validator.go`
- [ ] Create compatibility test BEFORE changes
- [ ] Replace `MarshalSSZ()` for all 8 fields
- [ ] Replace `UnmarshalSSZ()` with proper parsing
- [ ] Update `HashTreeRoot()` to use fastssz
- [ ] Remove `DefineSSZ()` method
- [ ] Remove karalabe/ssz import
- [ ] Note: Already has `HashTreeRootWith()` fastssz method
- [ ] Test effective balance calculations
- [ ] Run all tests to ensure nothing breaks
- [ ] ✅ All tests pass before proceeding

### 8. BeaconBlockHeader Type
**File**: `/consensus-types/types/header.go`
- [ ] Create compatibility test BEFORE changes
- [ ] Replace `MarshalSSZ()` for all header fields
- [ ] Replace `UnmarshalSSZ()` with proper parsing
- [ ] Update `HashTreeRoot()` to use fastssz
- [ ] Remove `DefineSSZ()` method
- [ ] Remove karalabe/ssz import
- [ ] Note: Already has `HashTreeRootWith()` fastssz method
- [ ] Test header validation methods
- [ ] Run all tests to ensure nothing breaks
- [ ] ✅ All tests pass before proceeding

## Phase 3: Fork-Aware Types (High Risk)

### 9. BeaconBlock Type
**File**: `/consensus-types/types/block.go`
- [ ] Create `block_fastssz.go`
- [ ] Implement fork version detection
- [ ] Handle different block structures per fork
- [ ] Implement marshaling with fork logic
- [ ] Note: Does not currently have fastssz methods
- [ ] Preserve `SigningMessage` behavior
- [ ] Add build tags to original file
- [ ] Create comprehensive fork-specific tests
- [ ] Test state root validation
- [ ] Run all tests with both implementations
- [ ] ✅ All tests pass before proceeding

### 10. BeaconBlockBody Type
**File**: `/consensus-types/types/body.go`
- [ ] Create `body_fastssz.go`
- [ ] Handle 12 fields (Deneb) vs 13 fields (Electra)
- [ ] Implement fork-conditional field encoding
- [ ] Handle execution requests (Electra only)
- [ ] Note: Does not currently have fastssz methods
- [ ] Implement all embedded types
- [ ] Add build tags to original file
- [ ] Create fork-specific compatibility tests
- [ ] Test deposit/withdrawal limits
- [ ] Run all tests with both implementations
- [ ] ✅ All tests pass before proceeding

### 11. ExecutionPayload Type
**File**: `/consensus-types/types/payload.go`
- [ ] Create `payload_fastssz.go`
- [ ] Handle Capella+ withdrawals requirement
- [ ] Implement fork-aware marshaling
- [ ] Handle dynamic transactions field
- [ ] Note: Already has `HashTreeRootWith()` fastssz method
- [ ] Implement proper validation
- [ ] Add build tags to original file
- [ ] Create compatibility tests for each fork
- [ ] Test transaction size limits
- [ ] Run all tests with both implementations
- [ ] ✅ All tests pass before proceeding

### 12. ExecutionPayloadHeader Type
**File**: `/consensus-types/types/payload_header.go`
- [ ] Create `payload_header_fastssz.go`
- [ ] Mirror ExecutionPayload structure
- [ ] Handle fork-specific fields
- [ ] Implement all hash fields
- [ ] Note: Already has `HashTreeRootWith()` fastssz method
- [ ] Add build tags to original file
- [ ] Create compatibility tests
- [ ] Test against known headers
- [ ] Run all tests with both implementations
- [ ] ✅ All tests pass before proceeding

### 13. BeaconState Type (Most Complex)
**File**: `/consensus-types/types/state.go`
- [ ] Create `state_fastssz.go`
- [ ] Handle all 27+ fields
- [ ] Implement Electra's PendingPartialWithdrawals
- [ ] Handle all dynamic slices with limits
- [ ] Note: Already has `HashTreeRootWith()` fastssz method
- [ ] Implement fork version logic
- [ ] Add build tags to original file
- [ ] Create extensive compatibility tests
- [ ] Test state transitions
- [ ] Test validator registry operations
- [ ] Run all tests with both implementations
- [ ] ✅ All tests pass before proceeding

## Phase 4: Engine Primitives

### 14. Engine Withdrawal Type
**File**: `/engine-primitives/engine-primitives/withdrawal.go`
- [ ] Create `withdrawal_fastssz.go`
- [ ] Implement engine-specific withdrawal format
- [ ] Note: Already imports both karalabe/ssz and fastssz
- [ ] Add build tags
- [ ] Create compatibility tests
- [ ] ✅ All tests pass before proceeding

### 15. Additional Types to Migrate (Found during analysis)
**Files**: Types that use karalabe/ssz but weren't in original plan
- [ ] `/consensus-types/types/pending_partial_withdrawal.go` - Already has fastssz
- [ ] `/consensus-types/types/deposit.go` - Already has fastssz methods
- [ ] Various request types (DepositRequest, WithdrawalRequest, etc.)
- [ ] Add build tags and compatibility tests for each
- [ ] ✅ All tests pass before proceeding

## Phase 5: Infrastructure Components

### 16. SSZ Utils
**File**: `/primitives/encoding/ssz/utils.go`
- [ ] Create `utils_fastssz.go`
- [ ] Port generic unmarshal function
- [ ] Port EIP-7685 methods
- [ ] Add build tags
- [ ] Update all callers
- [ ] ✅ All tests pass before proceeding

### 17. Storage SSZ Codec
**File**: `/storage/encoding/ssz.go`
- [ ] Create `ssz_fastssz.go`
- [ ] Update encoder/decoder to use fastssz directly
- [ ] Maintain storage compatibility
- [ ] Add build tags
- [ ] Test storage round-trips
- [ ] ✅ All tests pass before proceeding

## Phase 6: Final Migration

### Enable FastSSZ by Default
- [ ] Change default build tags to use fastssz
- [ ] Run full test suite
- [ ] Run benchmarks
- [ ] Deploy to testnet
- [ ] Monitor for issues
- [ ] Remove old implementation files
- [ ] Remove karalabe/ssz dependency

## Benefits of In-Place Migration

1. **No Code Duplication**: Avoids duplicating struct definitions, constructors, and getter methods
2. **Aligns with Existing Pattern**: Files already use both libraries, we're just completing the transition
3. **Simpler Migration Path**: Just replace method implementations
4. **Cleaner End State**: No temporary files to clean up later
5. **Easier to Review**: Changes are in-place, making diffs clearer

## Compatibility Test Template

For in-place migration, create a compatibility test BEFORE making changes:

```go
// Save current behavior before migration
package types_test

import (
    "testing"
    "github.com/berachain/beacon-kit/consensus-types/types"
    "github.com/karalabe/ssz"
    "github.com/stretchr/testify/require"
)

func TestTypeCompatibility_BeforeMigration(t *testing.T) {
    // Create test objects
    testCases := []struct{
        name string
        obj  *types.YourType
    }{
        // Add various test cases
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Save current marshaling behavior
            currentBytes, err := tc.obj.MarshalSSZ()
            require.NoError(t, err)

            // Save current hash tree root
            currentRoot := tc.obj.HashTreeRoot()

            // After migration, these same tests should pass
            // demonstrating identical behavior
        })
    }
}
```

## Success Metrics

For each struct migration:
- [ ] All existing tests pass with both implementations
- [ ] Compatibility tests show identical behavior
- [ ] No performance regression
- [ ] No memory leaks
- [ ] Clean code coverage

## Rollback Plan

If issues arise with the new in-place migration approach:
1. Use git to revert the changes: `git revert <commit>`
2. Review what went wrong and fix the implementation
3. Re-run all tests before attempting migration again
4. Consider creating a feature branch for larger migrations

For the already-migrated types (Fork, ForkData, Eth1Data) using build tags:
1. Remove `fastssz` build tag from deployments
2. Original implementation automatically takes over
3. Fix issues in `_fastssz.go` files
4. Re-test before re-enabling

## Important Considerations

### Current State Analysis
1. **Hybrid Implementation Already Exists**: Most types already use both libraries:
   - `karalabe/ssz` for `MarshalSSZ()`, `UnmarshalSSZ()`, `DefineSSZ()`, `SizeSSZ()`
   - `fastssz` for `HashTreeRootWith()`, `GetTree()`, `MarshalSSZTo()`
   - This makes in-place migration the natural next step

2. **Missing Types**: Several files mentioned don't exist or have different names:
   - `withdrawal.go` (consensus-types) doesn't exist
   - `execution_payload.go` → actually `payload.go`
   - `execution_payload_header.go` → actually `payload_header.go`
   - `eth1.go` → actually `eth1data.go`
   - `withdrawal_credentials.go` exists but is a type alias, not a struct

3. **Migration Approach Evolution**:
   - Started with build tag approach for Fork, ForkData, Eth1Data
   - Switching to in-place migration for remaining types
   - May consolidate the already-migrated files later to remove build tags

## Notes
- Always run full test suite after each struct migration
- Keep original and new implementations in sync during transition
- Document any behavioral differences discovered
- Consider running both implementations in parallel in staging for validation
- Check if fastssz supports all features used by karalabe/ssz (e.g., `DefineSSZ` pattern)
