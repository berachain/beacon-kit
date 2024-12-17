# Deposit Merkle Tree

This package implements the [EIP-4881 spec](https://eips.ethereum.org/EIPS/eip-4881) for a minimal sparse Merkle tree.

The format proposed by the EIP allows for the pruning of deposits that are no longer needed to participate fully in consensus.

Thanks to [Prysm](https://github.com/prysmaticlabs/prysm/blob/develop/beacon-chain/cache/depositsnapshot) and Mark Mackey ([@ethDreamer](https://github.com/ethDreamer)) for the reference implementation and tests.
