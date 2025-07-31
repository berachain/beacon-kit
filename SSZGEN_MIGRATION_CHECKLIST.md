# SSZGen Migration Checklist

## Overview
This document tracks the migration from manually written SSZ code to sszgen-generated code for types in consensus-types/types.

## Types That SHOULD Be Migrated to sszgen

### Simple Fixed-Size Types ✅ Ready for sszgen
These have manual implementations but no fork-specific logic:

- [x] **Deposit** (deposit.go) ✅ MIGRATED
  - Fixed size: 192 bytes
  - Simple fields: Pubkey, Credentials, Amount, Signature, Index
  - Migrated to sszgen (with WithdrawalCredentials type alias workaround)
  
- [x] **Eth1Data** (eth1data.go) ✅ MIGRATED
  - Fixed size: 72 bytes
  - Simple fields: DepositRoot, DepositCount, BlockHash
  - Migrated to sszgen

- [x] **Validator** (validator.go) ✅ MIGRATED
  - Fixed size: 121 bytes
  - Simple fields with no conditional logic
  - Migrated to sszgen

- [x] **BeaconBlockHeader** (header.go) ✅ MIGRATED
  - Fixed size structure
  - Simple fields: Slot, ProposerIndex, ParentBlockRoot, StateRoot, BodyRoot
  - Migrated to sszgen

- [x] **PendingPartialWithdrawal** (pending_partial_withdrawal.go) ✅ MIGRATED
  - Fixed size structure
  - Simple fields: ValidatorIndex, Amount, WithdrawableEpoch
  - Migrated to sszgen

- [x] **SyncAggregate** (sync_aggregate.go) ✅ MIGRATED
  - Fixed size: 160 bytes
  - Simple arrays: SyncCommitteeBits, SyncCommitteeSignature
  - Migrated to sszgen

### Variable-Size Types ✅ Ready for sszgen
These have dynamic fields but no fork-specific logic:

- [x] **SignedBeaconBlockHeader** (signed_beacon_block_header.go) ✅ MIGRATED
  - Contains: Header (BeaconBlockHeader) + Signature
  - No fork-specific logic
  - Migrated to sszgen

## Migration Complete ✅

All simple types that can be migrated to sszgen have been successfully migrated. The migration included:
- 7 types migrated from manual SSZ to sszgen
- Added SSZ methods to supporting types (WithdrawalCredentials, KZGCommitment)
- Identified types that cannot be migrated due to architectural constraints

## Types That MUST Keep Manual Implementation

### Fork-Specific Logic Required ❌
These types have conditional logic based on fork version:

- **BeaconBlockBody** (body.go)
  - Conditionally includes `executionRequests` in Electra+
  - Complex offset calculations based on fork
  - **Keep manual implementation**
  - **Reference sszgen**: Create `body_temp_sszgen.go` for verification

- **BeaconState** (state.go)
  - Conditionally includes fields based on fork version
  - Complex state transition logic
  - **Keep manual implementation**
  - **Reference sszgen**: Create `state_temp_sszgen.go` for verification

- **ExecutionPayload** (payload.go)
  - Fork-specific validation (withdrawals in Capella+)
  - Dynamic field handling
  - **Keep manual implementation**
  - **Reference sszgen**: Create `payload_temp_sszgen.go` for verification

- **ExecutionPayloadHeader** (payload_header.go)
  - Similar fork-specific logic as ExecutionPayload
  - **Keep manual implementation**
  - **Reference sszgen**: Create `payload_header_temp_sszgen.go` for verification

### Container Types with References ⚠️
These contain other types that have fork-specific logic:

- **BeaconBlock** (block.go)
  - Contains BeaconBlockBody which has fork-specific logic
  - Might need manual implementation to handle body correctly

- **SignedBeaconBlock** (signed_beacon_block.go)
  - Contains BeaconBlock which contains BeaconBlockBody
  - Might need manual implementation

### Special Collection Types ❌ Cannot Migrate
These are slice type aliases, not structs, so sszgen cannot be used:

- **ExecutionRequests** (execution_requests.go)
  - Custom list encoding for EIP-7685
  - **Keep manual implementation**

- **Deposits** (deposits.go) ❌ CANNOT MIGRATE
  - Type alias for slice: `type Deposits []*Deposit`
  - sszgen only works on structs, not slice types
  - Must keep partial manual implementation

- **Validators** (validators.go) ❌ CANNOT MIGRATE
  - Type alias for slice: `type Validators []*Validator`
  - sszgen only works on structs, not slice types
  - Must keep partial manual implementation

## Types Already Using sszgen ✅
These already have generated code:

- AttestationData (attestation_data_sszgen.go)
- ConsolidationRequest (consolidation_request_sszgen.go)
- DepositMessage (deposit_message_sszgen.go)
- Fork (fork_sszgen.go)
- ForkData (fork_data_sszgen.go)
- SigningData (signing_data_sszgen.go)
- SlashingInfo (slashing_info_sszgen.go)
- WithdrawalRequest (withdrawal_request_sszgen.go)

## Migration Summary

### Successfully Migrated ✅
The following types have been migrated from manual SSZ to sszgen:
1. **Deposit** - Fixed size struct (192 bytes)
2. **Eth1Data** - Fixed size struct (72 bytes)
3. **Validator** - Fixed size struct (121 bytes)
4. **BeaconBlockHeader** - Fixed size struct
5. **PendingPartialWithdrawal** - Fixed size struct
6. **SyncAggregate** - Fixed size struct (160 bytes)
7. **SignedBeaconBlockHeader** - Variable size struct

### Cannot Migrate ❌
The following types cannot be migrated due to architectural constraints:
1. **Deposits** - Slice type alias, not a struct
2. **Validators** - Slice type alias, not a struct
3. **BeaconBlockBody** - Fork-specific conditional logic
4. **BeaconState** - Fork-specific conditional logic
5. **ExecutionPayload** - Fork-specific validation
6. **ExecutionPayloadHeader** - Fork-specific logic
7. **ExecutionRequests** - Custom EIP-7685 encoding
8. **BeaconBlock** - Contains fork-specific BeaconBlockBody
9. **SignedBeaconBlock** - Contains fork-specific BeaconBlock

### Supporting Changes
- **WithdrawalCredentials** - Added SSZ methods to support Deposit/Validator migration
- **KZGCommitment** - Changed from `[48]byte` to `bytes.B48` with proper SSZ methods

## Reference sszgen Files for Manual Implementations

For types that must keep manual implementation due to fork-specific logic, we create temporary reference sszgen files to:
- Verify correctness of manual implementation
- Check field ordering and size calculations
- Identify missing nil checks or error handling
- Serve as documentation for the expected SSZ encoding

### Process for Creating Reference Files:
1. Create a temporary struct (e.g., `BeaconBlockBodyTemp`) with:
   - All fields from the original struct
   - Private fields made public (for sszgen access)
   - Embedded `Versionable` removed (has `ssz:"-"` tag)
   - Same field types and ordering

2. Add go:generate directive:
   ```go
   //go:generate sszgen -path body_temp.go -objs BeaconBlockBodyTemp -output body_temp_sszgen.go -include ...
   ```

3. Run sszgen to generate reference implementation

4. Compare generated code with manual implementation

5. Mark files clearly as "REFERENCE ONLY - NOT FOR PRODUCTION USE"

## Notes

1. Before migrating each type, check for:
   - Any custom validation in ValidateAfterDecodingSSZ
   - Interface requirements (some types implement specific interfaces)
   - Special size constants or calculations
   - Comments about HashTreeRoot compatibility

2. For collection types, verify:
   - Maximum size limits are enforced
   - Proper list encoding is maintained
   - Any special validation requirements

3. Types using UnusedType pattern should NOT be migrated as they're placeholders

4. Always test thoroughly after migration to ensure:
   - Serialization/deserialization works correctly
   - Hash tree root calculations match
   - Size calculations are correct
   - All existing tests pass