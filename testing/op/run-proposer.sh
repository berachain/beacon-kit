#!/bin/bash

cd ~/op-stack-deployment/optimism/op-proposer

# Set based on the wallets from deploy.sh
export PROPOSER_KEY=""
export L1_RPC="http://localhost:64064"
export L2OO_ADDR="" # TODO

./bin/op-proposer \
    --poll-interval 5s \
    --rpc.port 8560 \
    --rollup-rpc http://localhost:8547 \
    --l2oo-address $L2OO_ADDR \
    --private-key $PROPOSER_KEY \
    --l1-eth-rpc $L1_RPC
    