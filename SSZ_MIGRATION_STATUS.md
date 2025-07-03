# SSZ Migration Status Report

## Overview
The codebase is currently in a mixed state, using both `github.com/karalabe/ssz` and `github.com/prysmaticlabs/fastssz`. There are 51 files still importing karalabe/ssz.

## Files Still Using karalabe/ssz

### Core Types with Both Libraries (Partial Migration)

#### consensus-types/types (23 files):
- `attester_slashings.go` - Uses both, needs DefineSSZ removal
- `attestations.go` - Uses both, needs DefineSSZ removal  
- `block.go` - Uses both, needs DefineSSZ removal
- `bls_to_execution_changes.go` - Uses both, needs DefineSSZ removal
- `body.go` - Uses both, needs DefineSSZ removal
- `consolidation_request.go` - Has sszgen file, needs cleanup
- `deposit.go` - Uses both, needs DefineSSZ removal
- `deposits.go` - Uses both, needs DefineSSZ removal
- `eth1data.go` - Uses both, needs DefineSSZ removal
- `execution_requests.go` - Uses both, needs DefineSSZ removal
- `header.go` - Uses both, needs DefineSSZ removal
- `payload_header.go` - Uses both, needs DefineSSZ removal
- `payload.go` - Uses both, needs DefineSSZ removal
- `pending_partial_withdrawal.go` - Uses both, needs DefineSSZ removal
- `proposer_slashings.go` - Uses both, needs DefineSSZ removal
- `signed_beacon_block_header.go` - Uses both, needs DefineSSZ removal
- `signed_beacon_block.go` - Uses both, needs DefineSSZ removal
- `state.go` - Uses both, needs DefineSSZ removal
- `sync_aggregate.go` - Uses both, needs DefineSSZ removal
- `validator.go` - Uses both, needs DefineSSZ removal
- `validators.go` - Uses both, needs DefineSSZ removal
- `voluntary_exits.go` - Uses both, needs DefineSSZ removal
- `withdrawal_request.go` - Has sszgen file, needs cleanup

#### engine-primitives (3 files):
- `withdrawal.go` - Has fastssz but still uses karalabe methods
- `withdrawals.go` - Needs full migration
- `transactions.go` - Needs full migration

#### da/types (2 files):
- `sidecar.go` - Has fastssz but still uses karalabe methods
- `sidecars.go` - Needs full migration

### Infrastructure Dependencies

#### primitives (3 files):
- `constraints/ssz.go` - Defines interfaces using ssz.Object
- `common/interfaces.go` - Defines interfaces using ssz.Object  
- `common/unused_type.go` - Still uses DefineSSZ

#### utilities:
- `primitives/encoding/sszutil/utils.go` - Uses ssz.DecodeFromBytes

## Types Already Migrated to fastssz

### Fully Generated Types (have _sszgen.go files):
- `AttestationData`
- `ConsolidationRequest`
- `DepositMessage`
- `ForkData`
- `Fork`
- `SigningData`
- `SlashingInfo`
- `WithdrawalRequest`

## Migration Requirements

### 1. Remove karalabe/ssz Methods
All types need to remove:
- `DefineSSZ(codec *ssz.Codec)`
- `SizeSSZ(siz *ssz.Sizer) uint32` (karalabe version)
- Direct calls to `ssz.EncodeToBytes`, `ssz.DecodeFromBytes`, `ssz.HashSequential`, `ssz.HashConcurrent`

### 2. Generate fastssz Methods
Types needing sszgen generation:
- All types in consensus-types/types without _sszgen.go files
- engine-primitives: Withdrawal, Withdrawals, Transactions
- da/types: BlobSidecar, BlobSidecars

### 3. Update Interfaces
The constraint interfaces need updating:
- Remove dependency on `ssz.Object` 
- Update `SSZUnmarshaler` to use fastssz patterns
- Consider if we need these interfaces at all with fastssz

### 4. Update Utility Functions
- `sszutil.DecodeFromBytes` needs to be updated to use fastssz
- Any other utilities using karalabe/ssz directly

## Next Steps

1. **Generate fastssz for remaining types** - Run sszgen on all types that don't have it
2. **Remove dual implementations** - Delete all karalabe/ssz methods from types that have fastssz
3. **Update interfaces** - Refactor constraint interfaces to remove ssz.Object dependency
4. **Update utilities** - Migrate sszutil and other helpers to fastssz
5. **Remove imports** - Clean up all karalabe/ssz imports
6. **Test compatibility** - Ensure all SSZ encoding/decoding produces identical results