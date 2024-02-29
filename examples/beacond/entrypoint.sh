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
# Remember to change to other types of keyring like 'file' in-case exposing to outside world,
# otherwise your balance will be wiped quickly
# The keyring test does not require private key to steal tokens from you
KEYRING="test"
KEYALGO="secp256k1"
LOGLEVEL="info"
# Set dedicated home directory for the ./build/bin/beacond instance
HOMEDIR="./.tmp/beacond"
# to trace evm
#TRACE="--trace"
TRACE=""

# Path variables
CONFIG_TOML=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
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
	
	# Set client config
	./build/bin/beacond config set client keyring-backend $KEYRING --home "$HOMEDIR"
	./build/bin/beacond config set client chain-id "$CHAINID" --home "$HOMEDIR"

	# If keys exist they should be deleted
	./build/bin/beacond keys add dev --no-backup --keyring-backend $KEYRING --home "$HOMEDIR" 


	# Change parameter token denominations to abgt
	jq '.app_state["staking"]["params"]["bond_denom"]="abgt"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="abgt"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.app_state["mint"]["params"]["mint_denom"]="abgt"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.consensus["params"]["block"]["max_gas"]="30000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.consensus.params.validator.pub_key_types += ["bls12_381"] | .consensus.params.validator.pub_key_types -= ["ed25519"]' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	# Allocate genesis accounts (cosmos formatted addresses)
	./build/bin/beacond genesis add-genesis-account dev 100000000000000000000000000abgt --keyring-backend $KEYRING --home "$HOMEDIR"
	
	# Test Account
	# absurd surge gather author blanket acquire proof struggle runway attract cereal quiz tattoo shed almost sudden survey boring film memory picnic favorite verb tank
	# 0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306
	./build/bin/beacond genesis add-genesis-account cosmos1yrene6g2zwjttemf0c65fscg8w8c55w58yh8rl 69000000000000000000000000abgt --keyring-backend $KEYRING --home "$HOMEDIR"

	# Sign genesis transaction
	./build/bin/beacond genesis gentx dev 1000000000000000000000abgt --keyring-backend $KEYRING --chain-id $CHAINID --home "$HOMEDIR" 
	## In case you want to create multiple validators at genesis
	## 1. Back to `./build/bin/beacond keys add` step, init more keys
	## 2. Back to `./build/bin/beacond add-genesis-account` step, add balance for those
	## 3. Clone this ~/../build/bin/beacond home directory into some others, let's say `~/.cloned./build/bin/beacond`
	## 4. Run `gentx` in each of those folders
	## 5. Copy the `gentx-*` folders under `~/.cloned./build/bin/beacond/config/gentx/` folders into the original `~/../build/bin/beacond/config/gentx`

	# Collect genesis tx
	./build/bin/beacond genesis collect-gentxs --home "$HOMEDIR" > /dev/null 2>&1

	# Run this to ensure everything worked and that the genesis file is setup correctly
	./build/bin/beacond genesis validate-genesis --home "$HOMEDIR" > /dev/null 2>&1
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)m
./build/bin/beacond start --pruning=nothing "$TRACE" \
--log_level $LOGLEVEL --api.enabled-unsafe-cors \
--api.enable --api.swagger --minimum-gas-prices=0.0001abgt \
--home "$HOMEDIR" --beacon-kit.engine.jwt-secret-path ${JWT_SECRET_PATH}
