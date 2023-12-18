KEY="brick"
CHAINID="berachain-666"
MONIKER="brickchain"
KEYRING="test"
KEYALGO="secp256k1"
LOGLEVEL="info"
HOMEDIR="data/.polard"
TRACE=""
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json

if [ "$(ls -A $HOMEDIR)" ]; then
    echo "$HOMEDIR is not empty"
    polard start --pruning=nothing "$TRACE" --log_level $LOGLEVEL --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "$HOMEDIR"
else
    echo "$HOMEDIR is empty, creating a new network"
    
    polard init $MONIKER -o --chain-id $CHAINID --home "$HOMEDIR"

    jq '.app_state["staking"]["params"]["bond_denom"]="abera"' "$GENESIS" >"$TMP_GENESIS"
    jq '.app_state["crisis"]["constant_fee"]["denom"]="abera"' "$GENESIS" >"$TMP_GENESIS"
    jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="abera"' "$GENESIS" >"$TMP_GENESIS"
    jq '.app_state["evm"]["params"]["evm_denom"]="abera"' "$GENESIS" >"$TMP_GENESIS"
    jq '.app_state["mint"]["params"]["mint_denom"]="abera"' "$GENESIS" >"$TMP_GENESIS"
    jq '.consensus["params"]["block"]["max_gas"]="30000000"' "$GENESIS" >"$TMP_GENESIS"
    mv "$TMP_GENESIS" "$GENESIS"

    polard config set client keyring-backend $KEYRING --home "$HOMEDIR"
    polard config set client chain-id "$CHAINID" --home "$HOMEDIR"

    polard keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"

    polard genesis add-genesis-account $KEY 100000000000000000000000000abera --keyring-backend $KEYRING --home "$HOMEDIR"

    # polard genesis add-genesis-account cosmos1yrene6g2zwjttemf0c65fscg8w8c55w58yh8rl 100000000000000000000000000abera --keyring-backend $KEYRING --home "$HOMEDIR"

    polard genesis gentx $KEY 1000000000000000000000abera --keyring-backend $KEYRING --chain-id $CHAINID --home "$HOMEDIR"

    polard genesis collect-gentxs --home "$HOMEDIR"

    polard genesis validate-genesis --home "$HOMEDIR"

    polard start --pruning=nothing "$TRACE" --log_level $LOGLEVEL --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home "$HOMEDIR"
    polard start --pruning=nothing '' --log_level info --api.enabled-unsafe-cors --api.enable --api.swagger --minimum-gas-prices=0.0001abera --home data/.polard
fi