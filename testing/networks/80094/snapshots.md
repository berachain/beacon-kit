# Berachain Mainnet Snapshots

The following are snapshots provided by the community.

## beacond snapshots

| Provider | URL | Database |
| -------- | --- | -------- |
| GhostGraph - TryGhosst.XYZ | https://public-snapshots.ghostgraph.xyz/bera/snapshot-beacond-mainnet.tgz | pebbledb  |

### Installation advice

1. Stop `beacond`
2. Move your existing `data` directory out of the way.
3. Uncompress the snapshot in the directory that held `data`
4. Copy your `priv_validator_state.json` file into the new `data` directory
5. Restart `beacond`
6. Only delete the old `data` directory.

Typical: 
```
sudo systemctl stop beacond
mv data data-old
tar -xzf snapshot-beacond-mainnet.tgz
cp data-old/priv_validator_state.json data/priv_validator_state.json
sudo systemctl start beacond
```
