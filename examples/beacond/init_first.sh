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

GENESIS=$BEACOND_HOME/config/genesis.json
TMP_GENESIS=$BEACOND_HOME/config/tmp_genesis.json

# Create beacond config directory
if [ ! -d "$BEACOND_HOME/config" ]; then
    # # Init the chain
    /usr/bin/beacond init --chain-id "$BEACOND_CHAIN_ID" "$BEACOND_MONIKER" --home "$BEACOND_HOME" --beacon-kit.accept-tos
	
	# Set client config
	/usr/bin/beacond config set client keyring-backend $BEACOND_KEYRING_BACKEND --home "$BEACOND_HOME"
	/usr/bin/beacond config set client chain-id "$BEACOND_CHAIN_ID" --home "$BEACOND_HOME"

	# If keys exist they should be deleted
	/usr/bin/beacond keys add "$BEACOND_MONIKER" --keyring-backend $BEACOND_KEYRING_BACKEND --home "$BEACOND_HOME" --indiscreet --output json > "$BEACOND_HOME/config/mnemonic.json"


	# Change parameter token denominations to abgt
	jq '.app_state["staking"]["params"]["bond_denom"]="abgt"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="abgt"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.app_state["mint"]["params"]["mint_denom"]="abgt"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.consensus["params"]["block"]["max_gas"]="30000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.consensus.params.validator.pub_key_types += ["bls12_381"] | .consensus.params.validator.pub_key_types -= ["ed25519"]' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	# Allocate genesis accounts (cosmos formatted addresses)
	/usr/bin/beacond genesis add-genesis-account "$BEACOND_MONIKER" 100000000000000000000000000abgt --keyring-backend $BEACOND_KEYRING_BACKEND --home "$BEACOND_HOME"
	
	# Test Account
	# absurd surge gather author blanket acquire proof struggle runway attract cereal quiz tattoo shed almost sudden survey boring film memory picnic favorite verb tank
	# 0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306
	/usr/bin/beacond genesis add-genesis-account cosmos1yrene6g2zwjttemf0c65fscg8w8c55w58yh8rl 69000000000000000000000000abgt --keyring-backend $BEACOND_KEYRING_BACKEND --home "$BEACOND_HOME"

	# Sign genesis transaction
	/usr/bin/beacond genesis gentx "$BEACOND_MONIKER" 1000000000000000000000abgt --keyring-backend $BEACOND_KEYRING_BACKEND --chain-id $BEACOND_CHAIN_ID --home "$BEACOND_HOME" 
	## In case you want to create multiple validators at genesis
	## 1. Back to `/usr/bin/beacond keys add` step, init more keys
	## 2. Back to `/usr/bin/beacond add-genesis-account` step, add balance for those
	## 3. Clone this ~/./usr/bin/beacond home directory into some others, let's say `~/.cloned/usr/bin/beacond`
	## 4. Run `gentx` in each of those folders
	## 5. Copy the `gentx-*` folders under `~/.cloned/usr/bin/beacond/config/gentx/` folders into the original `~/./usr/bin/beacond/config/gentx`
fi
