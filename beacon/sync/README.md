# sync

## Types of Syncing

### Optimistic

- We forkchoice update as we replay blocks from the head.

### Checkpoint

- Specify a Checkpoint block to sync to.
- Forkchoice update to that checkpoint, still relies on peering at ETH level.

### Regular

- Call newPayload every bloc
