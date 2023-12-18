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

KEY="$1"
MONIKER="$1"
TRACE=""
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json


polard init $MONIKER -o --chain-id $CHAINID --home "$HOMEDIR"

polard config set client keyring-backend $KEYRING --home "$HOMEDIR"

polard keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"
