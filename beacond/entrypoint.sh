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
# Set dedicated home directory for the ./build/bin/beacond instance
HOMEDIR="./.tmp/beacond"

# Path variables
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json

# used to exit on first error (any non-zero exit code)
set -e

# Reinstall daemon
make build

if !command -v jq &> /dev/null
then
    echo "jq could not be found. Please install jq to continue."
    exit 1
fi


overwrite="N"
if [ -d "$HOMEDIR" ]; then
	printf "\nAn existing folder at '%s' was found. You can choose to delete this folder and start a new local node with new keys from genesis. When declined, the existing local node is started. \n" "$HOMEDIR"
	echo "Overwrite the existing configuration and start a new local node? [y/n]"
	read -r overwrite
else	
overwrite="Y"
fi

# Setup local node if overwrite is set to Yes, otherwise skip setup
if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
	# # Remove the previous folder
	rm -rf "$HOMEDIR"

	# # Set moniker and chain-id (Moniker can be anything, chain-id must be an integer)
	./build/bin/beacond init $MONIKER -o --chain-id $CHAINID --home "$HOMEDIR" --beacon-kit.accept-tos
	
	# Change parameter token denominations to abgt
	jq '.consensus.params.validator.pub_key_types += ["bls12_381"] | .consensus.params.validator.pub_key_types -= ["ed25519"]' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

	./build/bin/beacond genesis add-validator --home "$HOMEDIR"
	./build/bin/beacond genesis collect-validators --home "$HOMEDIR" 
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)m
./build/bin/beacond start --pruning=nothing "$TRACE" \
--log_level $LOGLEVEL --api.enabled-unsafe-cors \
--api.enable --api.swagger --minimum-gas-prices=0.0001abgt \
--home "$HOMEDIR" --beacon-kit.engine.jwt-secret-path ${JWT_SECRET_PATH}
