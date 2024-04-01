#!/bin/bash
# SPDX-License-Identifier: MIT
#
# Copyright (c) 2024 Berachain Foundation
#
# Permission is hereby granted, free of charge, to any person
# obtaining a copy of this software and associated documentation
# files (the "Software"), to deal in the Software without
# restriction, including without limitation the rights to use,
# copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following
# conditions:
#
# The above copyright notice and this permission notice shall be
# included in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
# EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
# OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
# NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
# HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
# WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
# FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

set -x

mv /root/.tmp_genesis/genesis.json /root/.beacond/config/genesis.json

sed -i "s/^prometheus = false$/prometheus = $BEACOND_ENABLE_PROMETHEUS/" $BEACOND_HOME/config/config.toml
sed -i "s/^prometheus_listen_addr = ":26660"$/prometheus_listen_addr = "0.0.0.0:26660"/" $BEACOND_HOME/config/config.toml
sed -i "s/^addr_book_strict = .*/addr_book_strict = false/" "$BEACOND_HOME/config/config.toml"



/usr/bin/beacond start --beacon-kit.engine.jwt-secret-path=/root/app/jwtsecret \
	--beacon-kit.accept-tos --beacon-kit.engine.rpc-dial-url $BEACOND_ENGINE_DIAL_URL \
	--beacon-kit.engine.required-chain-id $BEACOND_ETH_CHAIN_ID \
    --p2p.persistent_peers "$BEACOND_PERSISTENT_PEERS" \
    --rpc.laddr tcp://0.0.0.0:26657 \
    --grpc.address 0.0.0.0:9090 \
    --api.address tcp://0.0.0.0:1317 --api.enable \
    