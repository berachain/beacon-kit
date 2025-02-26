#!/usr/bin/env bash

export GETH=/Users/rezbera/Code/ethereum-pos-testnet/dependencies/go-ethereum/build/bin/geth

rm -r .tmp

$GETH init --datadir .tmp/gethdata /var/folders/z8/dnjjtbt10z50dtk_5sfckr7r0000gn/T/TestSimulatedCometComponentTestProcessProposal_BadBlock_IsRejected691375259/001/eth-genesis.json

$GETH --http --http.addr 0.0.0.0 --http.api eth,net,web3,debug \
                                 --authrpc.addr 0.0.0.0 \
                                 --authrpc.jwtsecret ./testing/files/jwt.hex \
                                 --authrpc.vhosts '*' \
                                 --datadir .tmp/gethdata \
                                 --ipcpath .tmp/gethdata/geth.ipc \
                                 --syncmode full \
                                 --verbosity 4 \
                                 --nodiscover

sleep 100000