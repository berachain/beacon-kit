KEY1="val0"
KEYRING="test"
HOMEDIR="/.polard"

polard genesis add-genesis-account $KEY1 100000000000000000000000000abera --keyring-backend $KEYRING --home "$HOMEDIR"
