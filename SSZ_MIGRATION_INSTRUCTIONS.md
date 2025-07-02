# SSZ Migration Instructions: karalabe/ssz → fastssz

This guide documents the process of migrating types from `karalabe/ssz` to `fastssz` with auto-generated code using `sszgen`. The migration maintains full compatibility while improving performance and reducing manual code maintenance.

## Overview

The migration involves:
1. Removing manual SSZ implementations using karalabe/ssz
2. Adding go:generate directives for sszgen
3. Auto-generating SSZ methods using fastssz
4. Creating compatibility tests to ensure backward compatibility

## Step-by-Step Migration Process

### Step 1: Analyze the Existing Type

First, identify the type to migrate and understand its structure. Look for:
- Field types and their SSZ encodings
- Any custom SSZ tags (though fastssz uses different tags)
- The computed SSZ size of the type

Example from Fork type:
```go
type Fork struct {
    PreviousVersion common.Version `json:"previous_version"`  // 4 bytes
    CurrentVersion  common.Version `json:"current_version"`   // 4 bytes  
    Epoch          math.Epoch     `json:"epoch"`             // 8 bytes
}
// Total size: 16 bytes
```

### Step 2: Remove karalabe/ssz Implementation

Remove all karalabe/ssz-specific code from your type file:

1. **Remove interface assertions:**
   ```go
   // Remove these:
   var (
       _ ssz.StaticObject                    = (*Fork)(nil)
       _ constraints.SSZMarshallableRootable = (*Fork)(nil)
   )
   ```

2. **Remove SSZ method implementations:**
   - `SizeSSZ(*ssz.Sizer) uint32`
   - `DefineSSZ(*ssz.Codec)`
   - `MarshalSSZ() ([]byte, error)`
   - `HashTreeRoot() common.Root`
   - Any other karalabe/ssz specific methods

3. **Remove imports:**
   ```go
   // Remove:
   import "github.com/karalabe/ssz"
   ```

4. **Keep important bits:**
   - Keep the type definition unchanged
   - Keep any constructors (e.g., `NewFork()`)
   - Keep any validation methods (e.g., `ValidateAfterDecodingSSZ()`)
   - Keep size constants if they're used elsewhere

### Step 3: Add go:generate Directive

Add a go:generate directive before your type definition:

```go
//go:generate sszgen --path . --include ../../primitives/common,../../primitives/bytes,../../primitives/math --objs Fork --output fork_sszgen.go
```

Breaking down the directive:
- `--path .`: Generate in the current directory
- `--include`: Add import paths for custom types used in your struct
- `--objs`: The type name(s) to generate SSZ for (comma-separated for multiple)
- `--output`: The output filename (convention: `<type>_sszgen.go`)

### Step 4: Add SSZ Size Tags

fastssz uses struct tags to specify sizes for certain types:

```go
type Fork struct {
    PreviousVersion common.Version `json:"previous_version" ssz-size:"4"`
    CurrentVersion  common.Version `json:"current_version" ssz-size:"4"`
    Epoch          math.Epoch     `json:"epoch"`
}
```

Common tags:
- `ssz-size:"N"`: For fixed-size byte arrays
- `ssz-max:"N"`: For variable-length slices/arrays

### Step 5: Generate SSZ Code

Run the code generation:

```bash
go generate ./...
# or specifically:
go generate ./consensus-types/types
```

This creates a `*_sszgen.go` file with:
- `MarshalSSZ() ([]byte, error)`
- `MarshalSSZTo(buf []byte) ([]byte, error)`
- `UnmarshalSSZ(buf []byte) error`
- `SizeSSZ() int`
- `HashTreeRoot() ([32]byte, error)`
- `HashTreeRootWith(hh ssz.HashWalker) error`
- `GetTree() (*ssz.Node, error)`

### Step 6: Create Compatibility Test

Create a test file `<type>_ssz_compatibility_test.go` to verify the migration maintains compatibility:

```go
package types

import (
    "testing"
    
    "github.com/berachain/beacon-kit/primitives/common"
    "github.com/berachain/beacon-kit/primitives/math"
    "github.com/karalabe/ssz"
    "github.com/stretchr/testify/require"
)

// Create a wrapper type with karalabe/ssz methods
type KaralabeFork struct {
    Fork  // Embed to avoid duplicating fields
}

// Copy the original karalabe/ssz methods
func (f *KaralabeFork) SizeSSZ(*ssz.Sizer) uint32 {
    return 16  // Or use a constant
}

func (f *KaralabeFork) DefineSSZ(codec *ssz.Codec) {
    ssz.DefineStaticBytes(codec, &f.PreviousVersion)
    ssz.DefineStaticBytes(codec, &f.CurrentVersion)
    ssz.DefineUint64(codec, &f.Epoch)
}

func (f *KaralabeFork) MarshalSSZ() ([]byte, error) {
    buf := make([]byte, ssz.Size(f))
    return buf, ssz.EncodeToBytes(buf, f)
}

func (f *KaralabeFork) HashTreeRoot() common.Root {
    return ssz.HashSequential(f)
}

// Test all operations produce identical results
func TestForkSSZCompatibility(t *testing.T) {
    // Test various scenarios...
}
```

Key test cases to include:
1. **Marshaling compatibility** - Both produce identical bytes
2. **Unmarshaling compatibility** - Can unmarshal bytes from either implementation
3. **Size calculation** - Both report the same size
4. **Hash tree root** - Both produce the same root (may need type conversion)
5. **Error cases** - Both handle errors similarly

### Step 7: Handle Type Differences

Common compatibility issues and solutions:

1. **Hash Tree Root return types:**
   - karalabe: Often returns custom types like `common.Root`
   - fastssz: Returns `[32]byte`
   - Solution: Use type conversion in tests

2. **Method signatures:**
   - karalabe: `SizeSSZ(*ssz.Sizer) uint32`
   - fastssz: `SizeSSZ() int`
   - Solution: Adapt in test wrapper

3. **Error handling:**
   - Verify both handle malformed input similarly

### Step 8: Update Build and CI

1. **Add sszgen to build tools:**
   ```bash
   go install github.com/ferranbt/fastssz/sszgen@latest
   ```

2. **Update Makefile/build scripts:**
   ```makefile
   generate:
       go generate ./...
   ```

3. **Add generated files to git:**
   - Commit the `*_sszgen.go` files
   - They should be regenerated in CI to verify consistency

4. **Update .gitignore if needed:**
   - Don't ignore `*_sszgen.go` files

## Migration Checklist

### General Steps for Each Type
- [ ] Identify all karalabe/ssz imports and methods
- [ ] Remove karalabe/ssz implementation
- [ ] Add go:generate directive with correct paths
- [ ] Add ssz struct tags where needed
- [ ] Run go generate
- [ ] Create compatibility test
- [ ] Verify all tests pass
- [ ] Update any code that depends on specific return types
- [ ] Commit both the modified source and generated files
- [ ] Update documentation if the type is part of public API

### Types to Migrate

#### consensus-types/types (29 types)
- [x] Fork (completed)
- [ ] AttestationData
- [ ] AttesterSlashings (slice type)
- [ ] Attestations (slice type)
- [ ] BeaconBlock (complex, has DefineSSZ)
- [ ] BeaconBlockBody (complex, has DefineSSZ)
- [ ] BeaconBlockHeader
- [ ] BlsToExecutionChanges (slice type)
- [ ] ConsolidationRequest
- [ ] Deposit
- [ ] DepositMessage
- [ ] Deposits (slice type)
- [ ] Eth1Data
- [ ] ExecutionPayloadHeader (complex, has DefineSSZ)
- [ ] ExecutionPayload (complex, has DefineSSZ)
- [ ] ExecutionRequests (complex, has DefineSSZ)
- [ ] ForkData
- [ ] PendingPartialWithdrawal
- [ ] ProposerSlashings (slice type)
- [ ] SignedBeaconBlock (complex, has DefineSSZ)
- [ ] SignedBeaconBlockHeader
- [ ] SigningData
- [ ] SlashingInfo
- [ ] BeaconState (very complex, has DefineSSZ)
- [ ] SyncAggregate
- [ ] Validator
- [ ] Validators (slice type)
- [ ] VoluntaryExits (slice type)
- [ ] WithdrawalRequest

#### da/types (2 types)
- [ ] BlobSidecar
- [ ] BlobSidecars (complex, has DefineSSZ)

#### engine-primitives/engine-primitives (3 types)
- [ ] Withdrawal
- [ ] Withdrawals (slice type)
- [ ] Transactions (slice type, has DefineSSZ)

#### primitives/common (1 type)
- [ ] UnusedType (test/utility type)

### Migration Priority

Consider migrating in this order:

1. **Simple types first**: Types that only implement StaticObject without custom DefineSSZ
   - AttestationData, BeaconBlockHeader, ConsolidationRequest, Deposit, DepositMessage, 
   - Eth1Data, ForkData, PendingPartialWithdrawal, SignedBeaconBlockHeader, SigningData,
   - SlashingInfo, SyncAggregate, Validator, WithdrawalRequest, Withdrawal

2. **Slice types**: Dynamic arrays that need ssz-max tags
   - AttesterSlashings, Attestations, BlsToExecutionChanges, Deposits, ProposerSlashings,
   - Validators, VoluntaryExits, Withdrawals, Transactions

3. **Complex types**: Types with custom DefineSSZ methods
   - BeaconBlock, BeaconBlockBody, ExecutionPayloadHeader, ExecutionPayload,
   - ExecutionRequests, SignedBeaconBlock, BlobSidecars

4. **State types last**: Large complex structures
   - BeaconState (most complex, do last)

### Infrastructure Updates
- [x] Create storage/encoding/fastssz.go (in progress)
- [x] Create primitives/constraints/fastssz.go (in progress)
- [ ] Update primitives/encoding/sszutil/utils.go for fastssz
- [ ] Update build tools to include sszgen
- [ ] Update CI/CD pipeline for code generation

## Common Patterns for Complex Types

### Types with Slices
```go
//go:generate sszgen --path . --objs MyType --output mytype_sszgen.go

type MyType struct {
    FixedArray [32]byte     `ssz-size:"32"`
    VarSlice   []byte       `ssz-max:"1024"`  
    Validators []Validator  `ssz-max:"1000000"`
}
```

### Nested Types
Ensure all nested types also implement SSZ. You may need to:
1. Generate SSZ for nested types first
2. Include their paths in the --include flag
3. Order go:generate directives correctly

### Union Types
For types that can be one of several options, you'll need custom handling as sszgen doesn't directly support unions.

## Troubleshooting

### "Type not found" errors
- Ensure --include paths are correct
- Check that imported types have SSZ methods

### Size mismatches
- Verify ssz-size tags match actual type sizes
- Check that all fields are accounted for

### Hash tree root differences
- Ensure field ordering is identical
- Verify no fields were missed
- Check for default values affecting hashing

## Benefits of Migration

1. **Performance**: fastssz is optimized for speed
2. **Maintenance**: Auto-generated code reduces errors
3. **Consistency**: Generated code follows same patterns
4. **Type Safety**: Compile-time verification of sizes
5. **Testing**: Generated code includes comprehensive methods

## Next Steps

After migrating a type:
1. Look for other types that depend on it
2. Consider migrating related types together
3. Update any benchmarks to compare performance
4. Document any breaking changes for API users

Remember: The goal is maintaining 100% compatibility while gaining the benefits of code generation and improved performance.