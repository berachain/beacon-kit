if [ -z "$CHAINID" ]; then
    CHAINID="brickchain-666"
fi
if [ -z "$KEYRING" ]; then
    KEYRING="test"
fi
if [ -z "$HOMEDIR" ]; then
    HOMEDIR="/.polard"
fi

KEY="$1"

polard genesis add-genesis-account $KEY 100000000000000000000000000abera,100000000000000000000000000stgusdc --keyring-backend $KEYRING --home "$HOMEDIR"

polard genesis gentx $KEY 1000000000000000000000abera --keyring-backend $KEYRING --chain-id $CHAINID --home "$HOMEDIR" \
    --moniker="$KEY" \
    --identity="identity of $KEY" \
    --details="This is $KEY" \
    --security-contact="brick@berachain.com" \
    --website="https://quantumwn.org/"
