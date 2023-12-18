KEY2="val1"
CHAINID="brickchain-666"
MONIKER1="val-1"
KEYRING="test"
KEYALGO="secp256k1"
HOMEDIR="/.polard"

polard init $MONIKER1 -o --chain-id $CHAINID --home "$HOMEDIR"

polard config set client keyring-backend $KEYRING --home "$HOMEDIR"

polard keys add $KEY1 --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"
  