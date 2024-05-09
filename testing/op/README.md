# Directions

0. Start the L1 chain. If using local, run a kurtosis test environment with `make start-devnet` from the root of the beacon-kit repo.
1. If it's your first time, run `setup.sh` to install Optimism dependencies.
2. Set your L1 (local or remote) values in `deploy.sh`.
3. Run `deploy.sh` to intialize wallets and deploy L1 contracts.
4. Run the run- files in this order to start the OP L2: `geth`, `node`, `batcher`, and `proposer`.
