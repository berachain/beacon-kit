# PoL deployment on devnet

We will need two repositories:

- BeaconKit repo, `main`
  [branch](https://github.com/berachain/beacon-kit/tree/main).
- Contracts repo, `main`
  [branch](https://github.com/berachain/contracts/tree/main), to deploy PoL
  smart contracts.

## Run BeaconKit

Compile and run BeaconKit from the branch mentioned above, with the usual
`make start` and `make start-geth` or `make start-reth`.

## Deploy PoL contracts

First off make sure you have run `forge i` and `bun i` in the contract repo.

On the contracts repo, run the following commands:

```bash
export FOUNDRY_PROFILE="deploy";
export IS_TESTNET=false;
export USE_SOFTWARE_WALLET=true;
export ETH_FROM="0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4";
export RPC_URL="http://localhost:8545";
export ETH_FROM_PK="0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306";
```

where `ETH_FROM` and `ETH_FROM_PK` are the preloaded EVM keys as described in
`BeaconKit`
[README](https://github.com/berachain/beacon-kit/blob/main/README.md).

First off check address for relevant commands. Run

```bash
forge script script/pol/POLPredictAddresses.s.sol -vv
```

and copy the output in the file `script/pol/POLAddresses.sol` file. Usually
just two contracts addresses, `BERACHEF_ADDRESS` and
`REWARD_VAULT_FACTORY_ADDRESS` should have changed, but there may be more.

**Very important check**: when deploying PoL from scratch make sure that:

- `<BGT_CONTRACT_ADDRESS>` matches with `specData.EVMInflationAddressDeneb1` in `BeaconKit`.
- `<DISTRIBUTOR_ADDRESS>` matches with `polDistributorAddress` in the node `eth-genesis.json`.

Moreover run

```bash
cast balance <BGT_CONTRACT_ADDRESS> --rpc-url $RPC_URL
```

and observe the balance going up every time the validator produces a block.

Once that is done run the following commands:

```bash
forge script script/pol/deployment/2_DeployBGT.s.sol \
  --private-key $ETH_FROM_PK --sender $ETH_FROM \
  --rpc-url $RPC_URL --broadcast -vv;
forge script script/pol/deployment/3_DeployPoL.s.sol \
  --private-key $ETH_FROM_PK --sender $ETH_FROM \
  --rpc-url $RPC_URL --broadcast -vv;
forge script script/pol/actions/ChangePOLParameters.s.sol \
  --private-key $ETH_FROM_PK --sender $ETH_FROM \
  --rpc-url $RPC_URL --broadcast -vv;
```

Now we need to generate 5 tokens to be associated with the 5 reward vaults we
will have in the default reward allocations. Run:

```bash
forge script script/misc/testnet/DeployToken.s.sol \
  --sig "deployBST(uint256)" 1 --sender $ETH_FROM \
  --private-key $ETH_FROM_PK --rpc-url $RPC_URL --broadcast;
forge script script/misc/testnet/DeployToken.s.sol \
  --sig "deployBST(uint256)" 2 --sender $ETH_FROM \
  --private-key $ETH_FROM_PK --rpc-url $RPC_URL --broadcast;
forge script script/misc/testnet/DeployToken.s.sol \
  --sig "deployBST(uint256)" 3 --sender $ETH_FROM \
  --private-key $ETH_FROM_PK --rpc-url $RPC_URL --broadcast;
forge script script/misc/testnet/DeployToken.s.sol \
  --sig "deployBST(uint256)" 4 --sender $ETH_FROM \
  --private-key $ETH_FROM_PK --rpc-url $RPC_URL --broadcast;
forge script script/misc/testnet/DeployToken.s.sol \
  --sig "deployBST(uint256)" 5 --sender $ETH_FROM \
  --private-key $ETH_FROM_PK --rpc-url $RPC_URL --broadcast;
```

and for each command note down the token address, which is indicated in the
logs at the line `BST deployed at: <TOKEN_ADDRESS>`.

These addresses must be written in
`script/pol/actions/DeployRewardVault.s.sol`. You will find in the script
tokens named `LP_BERA_HONEY`, `LP_BERA_ETH`, ... `LP_BEE_HONEY`[^1]. Finally
run:

```bash
forge script script/pol/actions/DeployRewardVault.s.sol \
  --private-key $ETH_FROM_PK --sender $ETH_FROM \
  --rpc-url $RPC_URL --broadcast -vv;
```

Once you run this file note down the Reward Vaults addresses, which should have
appear in the log lines like:
`RewardVault deployed at <REWARD_VAULT_ADDRESS> for staking token <TOKEN_ADDRESS>`.

Now consider `script/pol/actions/WhitelistRewardVault.s.sol`. On this file:

- Drop anything related to `REWARD_VAULT_USDS_HONEY`, since we only need 5 tokens.
- Assign the reward vaults addresses from above to `REWARD_VAULT_BERA_HONEY/ETH/WBTC...`[^2].

Finally run:

```bash
forge script script/pol/actions/WhitelistRewardVault.s.sol \
  --private-key $ETH_FROM_PK --sender $ETH_FROM \
  --rpc-url $RPC_URL --broadcast -vv;
cast send <BERACHEF_ADDRESS> "setMaxWeightPerVault(uint96)" 2000 \
  --private-key $ETH_FROM_PK --rpc-url $RPC_URL -vv;
```

Now consider `script/pol/actions/SetDefaultRewardAllocation`. On this file:

- Set `REWARD_VAULT_BERA_HONEY/ETH/WBTC..._WEIGHT`s to `2000`.
- Assign reward vaults addresses from above in
  `REWARD_VAULT_BERA_HONEY/ETH/WBTC...`[^3].

Finally run:

```bash
forge script \
  script/pol/actions/SetDefaultRewardAllocation.s.sol:\
WhitelistIncentiveTokenScript \
  --private-key $ETH_FROM_PK --sender $ETH_FROM \
  --rpc-url $RPC_URL --broadcast -vv;
```

## Check BGT distribution is carried out

Post Pectra11 fork, BGT distribution is automatically carried out by the
execution layer. A way to check distribution is happening is checking that the
validator producing blocks is receiving its base-rate BGTs via

```bash
cast call <BGT_CONTRACT_ADDRESS> "balanceOf(address)(uint256)" \
  <OPERATOR_ADDRESS> --rpc-url $RPC_URL
```

[^1]: You may consider renaming the tokens with `LP_TOKEN_1`, `LP_TOKEN_2`, ... `LP_TOKEN_5`.
[^2]: Again you may consider renaming the Reward vaults with `REWARD_VAULT_1/2/3/4/5`.
[^3]: Again you may consider renaming `REWARD_VAULT_BERA_HONEY/ETH/WBTC`... with `REWARD_VAULT_1/2/3`... for every variable and variable content in the file.
