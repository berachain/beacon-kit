# YOUR_ETH_WALLET_PRIVATE_KEY=""
YOUR_BEACON_HOME_DIR=".tmp/beacond" # set to your beacond home folder (like $HOME/.beacond)

GENESIS_VALIDATORS_ROOT="0x9147586693b6e8faa837715c0f3071c2000045b54233901c2e7871b15872bc43" # do not change
GENESIS_FORK_VERSION="0x04000000" # do not change
VAL_DEPOSIT_GWEI_AMOUNT=32000000000 # do not change
VAL_WITHDRAW_CREDENTIAL="0x0100000000000000000000000000000000000000000000000000000000000000" # do not change

output=$(./build/bin/beacond deposit create-validator $VAL_WITHDRAW_CREDENTIAL $VAL_DEPOSIT_GWEI_AMOUNT $GENESIS_FORK_VERSION $GENESIS_VALIDATORS_ROOT --home $YOUR_BEACON_HOME_DIR)
echo "output: $output"
VAL_PUB_KEY=$(echo "$output" | awk -F'pubkey=' '{print $2}' | awk '{print $1}' | sed -r 's/\x1B\[[0-9;]*[mK]//g')
SEND_DEPOSIT_SIGNATURE=$(echo "$output" | awk -F'signature=' '{print $2}' | awk '{print $1}' | sed -r 's/\x1B\[[0-9;]*[mK]//g')
echo "SEND_DEPOSIT_SIGNATURE: $SEND_DEPOSIT_SIGNATURE"
# DEPOSIT_CONTRACT_ADDRESS="0x4242424242424242424242424242424242424242"
# ETH_RPC_URL="http://localhost:8545"

# cast send "$DEPOSIT_CONTRACT_ADDRESS" 'deposit(bytes,bytes,uint64,bytes)' \
#     "$VAL_PUB_KEY" "$VAL_WITHDRAW_CREDENTIAL" 32 "$SEND_DEPOSIT_SIGNATURE" \
#     --private-key "$YOUR_ETH_WALLET_PRIVATE_KEY" --value 32ether  --chain-id 80084 -r $ETH_RPC_URL --legacy
# # Command line we were using before:
# # beacond deposit create-validator $VAL_WITHDRAW_CREDENTIAL $VAL_DEPOSIT_GWEI_AMOUNT $GENESIS_FORK_VERSION $GENESIS_VALIDATORS_ROOT --private-key $YOUR_ETH_WALLET_PRIVATE_KEY  --home $YOUR_BEACON_HOME_DIR