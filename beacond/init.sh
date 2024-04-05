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

# Initialize first validator flag
flag_first_validator=false

# Process flags
while getopts 'f' flag; do
  case "${flag}" in
    f) flag_first_validator=true ;;
    *) error "Unexpected option ${flag}" ;;
  esac
done

# Init the chain
/usr/bin/beacond init --chain-id "$BEACOND_CHAIN_ID" "$BEACOND_MONIKER" --home "$BEACOND_HOME" --beacon-kit.accept-tos


# Create beacond config directory
if [ "$flag_first_validator" = true ]; then
    GENESIS=$BEACOND_HOME/config/genesis.json
    TMP_GENESIS=$BEACOND_HOME/config/tmp_genesis.json

	jq '.consensus.params.validator.pub_key_types += ["bls12_381"] | .consensus.params.validator.pub_key_types -= ["ed25519"]' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
fi


/usr/bin/beacond genesis add-validator --home "$BEACOND_HOME"
/usr/bin/beacond genesis collect-validators --home "$BEACOND_HOME"



