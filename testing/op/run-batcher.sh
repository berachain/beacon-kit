#!/bin/bash

cd ~/op-stack-deployment/optimism/op-batcher

# Set based on the wallets from deploy.sh
export BATCHER_KEY=""
export L1_RPC="http://localhost:64064"

./bin/op-batcher \
    --l2-eth-rpc=http://localhost:7545 \
    --rollup-rpc=http://localhost:8547 \
    --poll-interval=1s \
    --sub-safety-margin=6 \
    --num-confirmations=1 \
    --safe-abort-nonce-too-low-count=3 \
    --resubmission-timeout=30s \
    --rpc.addr=0.0.0.0 \
    --rpc.port=8548 \
    --rpc.enable-admin \
    --max-channel-duration=1 \
    --l1-eth-rpc=$L1_RPC \
    --private-key=$BATCHER_KEY
