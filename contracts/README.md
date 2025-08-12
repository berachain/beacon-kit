# BeaconKit Contracts

## Mock PoL Contracts

Mock Proof-of-Liquidity contracts for testing execution client integration.

### Usage

Generate bytecode for genesis deployment:

```bash
# Build contracts
forge build

# Extract SimplePoLDistributor bytecode
cat out/MockPoL.sol/SimplePoLDistributor.json | jq -r .deployedBytecode.object

# Extract ValidatorRegistry bytecode  
cat out/MockValidatorRegistry.sol/ValidatorRegistry.json | jq -r .deployedBytecode.object
```

The same pattern can be used to extract bytecode for other contracts in the `brip0004/` directory.

### Contracts

- `MockPoL.sol` - Basic PoL distributor with multi-contract state changes
- `MockPoLReverting.sol` - PoL distributor that reverts after 10 distributions
- `MockPoLGasEnforcer.sol` - Gas-constrained PoL distributor for gas limit testing
- `MockValidatorRegistry.sol` - Registry contract for testing cross-contract calls
