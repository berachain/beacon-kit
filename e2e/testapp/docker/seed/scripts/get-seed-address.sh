if [ -z "$HOMEDIR" ]; then
    HOMEDIR="/.polard"
fi

ip=$1
node_id=$(polard comet show-node-id --home "$HOMEDIR")

echo "$node_id@$ip:26656"