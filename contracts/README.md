# BeaconKit Contracts

## Mock PoL Contracts

Mock Proof-of-Liquidity contracts for testing execution client integration.

### Usage

Generate bytecode for genesis deployment:

```bash
forge script script/GetBytecode.s.sol
```

### Contracts

- `MockPoL.sol` - Basic PoL distributor with multi-contract state changes
- `MockPoLReverting.sol` - PoL distributor that reverts after 10 distributions
- `MockPoLGasEnforcer.sol` - Gas-constrained PoL distributor for gas limit testing
- `MockValidatorRegistry.sol` - Registry contract for testing cross-contract calls