KEY1="val1"
KEYRING="test"
HOMEDIR="/.polard"

polard genesis add-genesis-account $KEY1 100000000000000000000000000abera --keyring-backend $KEYRING --home "$HOMEDIR"
