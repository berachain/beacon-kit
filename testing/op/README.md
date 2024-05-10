# Running an OP Stack L2 on top of a beacon-kit L1

This guide will walk you through setting up a OP stack rollup using a beacon-kit L1 chain. For more information, please visit the [Optimism guide](https://docs.optimism.io/builders/chain-operators/tutorials/create-l2-rollup).

The following directions provide the instructions to run the scripts in this directory.

## Directions

Note: this guide assumes you have `nvm`, `forge`, and `brew` already installed on your machine.

0. Start the L1 chain. If using local, run a kurtosis test environment with `make start-devnet` from the root of the beacon-kit repo.
1. If it's your first time, run `setup.sh` to install Optimism dependencies.
2. Set your L1 (local or remote) values in `deploy.sh`, specifically an `RPC_URL` for the L1.
3. Run `deploy.sh` to intialize wallets and deploy the L1 contracts.
4. Run the run- files in separate processes to start the OP L2: `geth`, `node`, `batcher`, and `proposer` (in this order).

### The L2 is up and running

The L2 has an RPC exposed at `http://localhost:8545`. You can bridge over to the L2 by running `bridge.sh`. From there you can begin interacting with the L2!
