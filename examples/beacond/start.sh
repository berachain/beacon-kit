#!/bin/bash
set -x

/usr/bin/beacond start --beacon-kit.engine.jwt-secret-path=/root/app/jwtsecret \
	--beacon-kit.accept-tos --beacon-kit.engine.rpc-dial-url $BEACOND_ENGINE_DIAL_URL \
	--beacon-kit.engine.required-chain-id $BEACOND_ETH_CHAIN_ID \
    --p2p.persistent_peers "$BEACOND_PERSISTENT_PEERS" \