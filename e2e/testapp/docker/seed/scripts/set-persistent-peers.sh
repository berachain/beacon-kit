if [ -z "$HOMEDIR" ]; then
    HOMEDIR="/.polard"
fi

CONFIG_TOML=$HOMEDIR/config/config.toml

seed_address=$1
sed -i "s/^persistent_peers = .*/persistent_peers = \"$seed_address\"/" "$CONFIG_TOML"
