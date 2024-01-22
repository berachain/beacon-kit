# Bootstrap Mode

- We cannot make `engine_` JSON-RPC requests during `FinalizeBlock()` as this will
cause the sync'ing process to go extremely slowly. The only operations that happen during this
period of ABCI should be updating the Beacon chain's view of what the finalized and safe
blocks are of the execution chain.
- However we need to come up with some sort of bootstrapping mode, to handle this. This is because if we sync the entire beacon chain from genesis and then only start syncing the 
execution chain once the beacon chain is fully sync'd we will run into problems.