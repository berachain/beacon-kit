if [ -z "$CHAINID" ]; then
    CHAINID="brickchain-666"
fi
if [ -z "$KEYRING" ]; then
    KEYRING="test"
fi
if [ -z "$KEYALGO" ]; then
    KEYALGO="secp256k1"
fi
if [ -z "$LOGLEVEL" ]; then
    LOGLEVEL="info"
fi
if [ -z "$HOMEDIR" ]; then
    HOMEDIR="/.polard"
fi

KEY1="seed-0"
MONIKER1="seed-0"
TRACE=""
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json


polard init $MONIKER1 -o --chain-id $CHAINID --home "$HOMEDIR"

jq '.app_state["staking"]["params"]["bond_denom"]="abera"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.app_state["crisis"]["constant_fee"]["denom"]="abera"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="abera"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.app_state["gov"]["params"]["min_deposit"][0]["denom"]="abera"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.app_state["gov"]["params"]["min_deposit"][0]["amount"]="10000000000000000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.app_state["gov"]["params"]["expedited_min_deposit"][0]["denom"]="abera"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.app_state["gov"]["params"]["expedited_min_deposit"][0]["amount"]="20000000000000000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.app_state["gov"]["params"]["max_deposit_period"]="300s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.app_state["gov"]["params"]["voting_period"]="300s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.app_state["gov"]["params"]["expedited_voting_period"]="240s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.app_state["gov"]["constitution"]="Honey is money."' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.consensus["params"]["block"]["max_gas"]="30000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.app_state["mint"]["params"]["mint_denom"]="abera"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";
jq '.consensus["params"]["abci"]["vote_extensions_enable_height"] = "2"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS";

polard config set client chain-id $CHAINID --home "$HOMEDIR"
polard config set client keyring-backend $KEYRING --home "$HOMEDIR"

polard keys add $KEY1 --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"

polard genesis add-genesis-account $KEY1 100000000000000000000000000abera,100000000000000000000000000stgusdc --keyring-backend $KEYRING --home "$HOMEDIR"

polard genesis gentx $KEY1 1000000000000000000000abera --keyring-backend $KEYRING --chain-id $CHAINID --home "$HOMEDIR" \
    --moniker="seed-0" \
    --identity="identity of seed-0" \
    --details="This is seed-0" \
    --security-contact="brick@berachain.com" \
    --website="https://quantumwn.org/"
