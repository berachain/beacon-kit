if [ -z "$HOMEDIR" ]; then
    HOMEDIR="/.polard"
fi
CONFIG_TOML=$HOMEDIR/config/config.toml

MONIKER=$1
sed -i "s/^moniker = .*/moniker = \"$MONIKER\"/" "$CONFIG_TOML"
