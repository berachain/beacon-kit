# Deploy PoL on devnet

We will need three repositories:

- BeaconKit repo, `pol-related-eth-genesis-data` [branch](https://github.com/berachain/beacon-kit/pull/2842).
- Contracts repo, `abear-PoL-devnet` [branch](https://github.com/berachain/contracts/tree/abear-PoL-devnet), to deploy PoL smart contracts.
- Distributor bot repo, `local-machine-distributeFor` [branch](https://github.com/berachain/berachain-v/tree/local-machine-distributeFor).

## Run BeaconKit

Compile and run BeaconKit from the branch mentioned above, with the usual `make start` and `make start-geth` or `make start-reth`.

Note that for the correct functioning of Distributor bot beaconKit needs to run no stop, without restart (this is because blocks in `BlockStore` are not persisted).

## Deploy PoL contracts

On the contracts repo, run the following commands:

```bash
export FOUNDRY_PROFILE="deploy";
export IS_TESTNET=false;
export USE_SOFTWARE_WALLET=true;
export ETH_FROM="0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4";
export RPC_URL="<http://localhost:8545>";
export ETH_FROM_PK="0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306";
```

where `ETH_FROM` and `ETH_FROM_PK` are the preloaded EVM keys as described in `BeaconKit` [README](https://github.com/berachain/beacon-kit/blob/main/README.md).

First off check address for relevant commands. Run

```bash
forge script script/pol/POLPredictAddresses.s.sol -vv
```

and copy the output in the file `script/pol/POLAddresses.sol` file. Just two contracts addresses, `BERACHEF_ADDRESS` and `REWARD_VAULT_FACTORY_ADDRESS` should have changed.

Once that is done run the following commands:

```bash
forge script script/pol/deployment/2_DeployBGT.s.sol --private-key $ETH_FROM_PK --sender $ETH_FROM --rpc-url $RPC_URL --broadcast -vv;
forge script script/pol/deployment/3_DeployPoL.s.sol --private-key $ETH_FROM_PK --sender $ETH_FROM --rpc-url $RPC_URL --broadcast -vv;
forge script script/pol/actions/ChangePOLParameters.s.sol --private-key $ETH_FROM_PK --sender $ETH_FROM --rpc-url $RPC_URL --broadcast -vv;
```

Now we need to generate 5 tokens to be associated with the 5 reward vaults we will have in the default reward allocations. Run:

```bash
forge script script/misc/testnet/DeployToken.s.sol --sig "deployBST(uint256)" 1 --sender $ETH_FROM --private-key $ETH_FROM_PK --rpc-url $RPC_URL --broadcast;
forge script script/misc/testnet/DeployToken.s.sol --sig "deployBST(uint256)" 2 --sender $ETH_FROM --private-key $ETH_FROM_PK --rpc-url $RPC_URL --broadcast;
forge script script/misc/testnet/DeployToken.s.sol --sig "deployBST(uint256)" 3 --sender $ETH_FROM --private-key $ETH_FROM_PK --rpc-url $RPC_URL --broadcast;
forge script script/misc/testnet/DeployToken.s.sol --sig "deployBST(uint256)" 4 --sender $ETH_FROM --private-key $ETH_FROM_PK --rpc-url $RPC_URL --broadcast;
forge script script/misc/testnet/DeployToken.s.sol --sig "deployBST(uint256)" 5 --sender $ETH_FROM --private-key $ETH_FROM_PK --rpc-url $RPC_URL --broadcast;
```

and for each command note down the token address, which is indicated in the logs at the line `BST deployed at: <TOKEN_ADDRESS>`.

These addresses must be written in `script/pol/actions/DeployRewardVault.s.sol`. You will find in the script tokens named `LP_BERA_HONEY`, `LP_BERA_ETH`, ... `LP_BEE_HONEY`. Rename them to `LP_TOKEN_1`, `LP_TOKEN_2`, ... `LP_TOKEN_5` and assign the address from above to them. Finally run:

```bash
forge script script/pol/actions/DeployRewardVault.s.sol --private-key $ETH_FROM_PK --sender $ETH_FROM --rpc-url $RPC_URL --broadcast -vv;
```

Once you run this file note down the Reward vaults addresses, which should have appear in the log lines like: `RewardVault deployed at <REWARD_VAULT_ADDRESS> for staking token <TOKEN_ADDRESS>`.

Now consider `script/pol/actions/WhitelistRewardVault.s.sol`. On this file:

- Drop anything related to `REWARD_VAULT_USDS_HONEY`, since we only need 5 tokens.
- Replace `REWARD_VAULT_BERA_HONEY/ETH/WBTC`... with `REWARD_VAULT_1/2/3`... for every variable and variable content in the file.
- Assign the reward vaults addresses from above in `REWARD_VAULT_1/2/3/4/5`

Finally run:

```bash
forge script script/pol/actions/WhitelistRewardVault.s.sol --private-key $ETH_FROM_PK --sender $ETH_FROM --rpc-url $RPC_URL --broadcast -vv;
cast send 0x233F5241c0C6d3Ea750d583ab2eD0Ca480446F39 "setMaxWeightPerVault(uint96)" 2000 --private-key $ETH_FROM_PK --rpc-url $RPC_URL -vv;
```

Now consider `script/pol/actions/SetDefaultRewardAllocation`. On this file:

- Replace `REWARD_VAULT_BERA_HONEY/ETH/WBTC`... with `REWARD_VAULT_1/2/3`... for every variable and variable content in the file.
- Set `REWARD_VAULT_1/2/3/4/5_WEIGHT`s to `2000`.
- Assign reward vaults addresses from above in `REWARD_VAULT_1/2/3/4/5`

Finally run:

```bash
forge script script/pol/actions/SetDefaultRewardAllocation.s.sol:WhitelistIncentiveTokenScript --private-key $ETH_FROM_PK --sender $ETH_FROM --rpc-url $RPC_URL --broadcast -vv;
```

Check that `BGT contract` is receiving `BERA`s being minted from `BeaconKit` via the command

```bash
cast balance <BGT_CONTRACT_ADDRESS>
```

## Run the Distributor bot

In the distributor repo, run

```bash
go run ./mod/distributor/cmd/main.go --config-path=./mod/distributor/example.config.toml
```

The log should tell you distribution is properly being executed. Also you can check that the validator producing blocks is receiving its base-rate BGTs via

```bash
cast call <BGT_CONTRACT_ADDRESS> "balanceOf(address)(uint256)" <VALIDATOR_ADDRESS>
```
