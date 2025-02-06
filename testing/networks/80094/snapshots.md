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
5. Restart `beacond`
6. Only after blocks start flowing is it a good idea to delete the old `data` directory.

Typical: 
```
sudo systemctl stop beacond
mv data data-old
tar -xzf snapshot-beacond-mainnet.tgz
sudo systemctl start beacond
```
