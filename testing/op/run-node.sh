#!/bin/bash

cd ~/op-stack-deployment/optimism/op-node

# Set based on the wallets from deploy.sh
export SEQ_KEY=""
export RPC_KIND="basic"
export L1_RPC="http://localhost:64064"

./bin/op-node \
	--l2=http://localhost:8551 \
	--l2.jwt-secret=./jwt.txt \
	--sequencer.enabled \
	--sequencer.l1-confs=1 \
	--verifier.l1-confs=1 \
	--rollup.config=./rollup.json \
	--rpc.addr=0.0.0.0 \
	--rpc.port=8547 \
	--p2p.disable \
	--rpc.enable-admin \
	--p2p.sequencer.key=$SEQ_KEY \
	--l1=$L1_RPC \
	--l1.rpckind=$RPC_KIND \
	--l1.trustrpc true
