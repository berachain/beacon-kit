if [ -z "$HOMEDIR" ]; then
    HOMEDIR="/.polard"
fi
if [ -z "$KEYRING" ]; then
    KEYRING="test"
fi
if [ -z "$KEYALGO" ]; then
    KEYALGO="secp256k1"
fi

polard genesis collect-gentxs --home "$HOMEDIR"

polard genesis validate-genesis --home "$HOMEDIR"

# # faucet
# polard keys add faucet --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"
# polard genesis add-genesis-account faucet 1000000000000000000000000000abera,1000000000000000000000000000stgusdc --keyring-backend $KEYRING --home "$HOMEDIR"

# # # Test Account
# # absurd surge gather author blanket acquire proof struggle runway attract cereal quiz tattoo shed almost sudden survey boring film memory picnic favorite verb tank
# # 0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306
# polard genesis add-genesis-account cosmos1yrene6g2zwjttemf0c65fscg8w8c55w58yh8rl 1000000000000000000000000000abera,1000000000000000000000000000stgusdc --keyring-backend $KEYRING --home "$HOMEDIR"
