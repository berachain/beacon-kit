#!/usr/bin/env bash

export GETH=/Users/rezbera/Code/ethereum-pos-testnet/dependencies/go-ethereum/build/bin/geth

$GETH init --datadir .tmp/gethdata .tmp/beacond/eth-genesis.json

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