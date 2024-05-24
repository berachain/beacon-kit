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

CHAINID="beacond-2061"
MONIKER="localtestnet"
LOGLEVEL="info"
CONSENSUS_KEY_ALGO="bls12_381"
HOMEDIR="./.tmp/beacond"

# Path variables
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json
ETH_GENESIS=./testing/files/eth-genesis.json # TODO: Fix this to not use a relative path or make it configurable

# used to exit on first error (any non-zero exit code)
set -e

# Reinstall daemon
make build

overwrite="N"
if [ -d $HOMEDIR ]; then
	printf "\nAn existing folder at '%s' was found. You can choose to delete this folder and start a new local node with new keys from genesis. When declined, the existing local node is started. \n" $HOMEDIR
	echo "Overwrite the existing configuration and start a new local node? [y/n]"
	read -r overwrite
else	
overwrite="Y"
fi

# Setup local node if overwrite is set to Yes, otherwise skip setup
if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
	rm -rf $HOMEDIR
	./build/bin/beacond init $MONIKER \
		--chain-id $CHAINID \
		--home $HOMEDIR \
		--beacon-kit.accept-tos \
		--consensus-key-algo $CONSENSUS_KEY_ALGO
	./build/bin/beacond genesis add-premined-deposit --home $HOMEDIR
	./build/bin/beacond genesis collect-premined-deposits --home $HOMEDIR 
	./build/bin/beacond genesis execution-payload "$ETH_GENESIS" --home $HOMEDIR
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)m
./build/bin/beacond start --pruning=nothing "$TRACE" \
--log_level $LOGLEVEL --api.enabled-unsafe-cors \
--api.enable --api.swagger --minimum-gas-prices=0.0001abgt \
--home $HOMEDIR --beacon-kit.engine.jwt-secret-path ${JWT_SECRET_PATH}
