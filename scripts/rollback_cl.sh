#!/bin/bash
# SPDX-License-Identifier: MIT
#
# Copyright (c) 2025 Berachain Foundation
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

# fail immediately on errors including any commands in a pipeline
set -Eeuo pipefail
# include linenumber and command if an error occurs
trap 'echo "Script failed at line $LINENO when running command: \"$BASH_COMMAND\""' ERR

usage() {
    echo "Usage: $(basename "$0")"
    echo ""
    echo "This script compares the Consensus Layer (CL) height with the Execution Layer (EL) height,"
    echo "and performs rollbacks on CL until it is at or below the EL height."
    echo ""
    echo "Environment Variables:"
    echo "  BEACOND_BINARY   (default: 'beacond') - Path to the beacon-kit binary."
    echo "  BEACOND_HOME     (default: '~/.beacond') - Path to the beacon-kit home directory."
    echo "  EL_RPC_URL       (default: '127.0.0.1:8545') - Execution Layer RPC endpoint."
    echo ""
    echo "Example usage:"
    echo "  BEACOND_HOME=/data/.beacond ./$(basename "$0")"
    echo "  BEACOND_BINARY=/custom/path/beacond ./$(basename "$0")"
    exit 0
}

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
    usage
fi

BEACOND_BINARY="${BEACOND_BINARY:-beacond}"
BEACOND_HOME="${BEACOND_HOME:-$HOME/.beacond}"
EL_RPC_URL="${EL_RPC_URL:-127.0.0.1:8545}"

echo "Starting rollback process:"
echo "- BEACOND_BINARY: $BEACOND_BINARY"
echo "- BEACOND_HOME: $BEACOND_HOME"
echo "- EL_RPC_URL: $EL_RPC_URL"

# Validate BEACOND_HOME already exists
if [[ ! -d "$BEACOND_HOME" ]]; then
    echo "Error: BEACOND_HOME is not a valid directory."
    exit 1
fi

# Validate BEACOND_BINARY (must be in PATH or specified)
if ! "$BEACOND_BINARY" version &>/dev/null; then
    echo "Error: BEACOND_BINARY is not a valid executable or not found."
    exit 1
fi

echo "[Fetching EL height...]"
EL_HEX=$(curl -s -X POST --location "$EL_RPC_URL" --header 'Content-Type: application/json' --data '{
    "jsonrpc":"2.0",
    "method":"eth_getBlockByNumber",
    "params":["latest", false],
    "id":1
}' | jq -r .result.number)
[[ -z "$EL_HEX" || "$EL_HEX" == "null" || "$EL_HEX" -le 0 ]] && echo "Error: Invalid Execution Layer height (EL). EL must be greater than zero." && exit 1
EL=$((${EL_HEX}))
echo "-> EL height: $EL ($EL_HEX)"

echo "[Fetching CL height...]"
ROLLBACK_OUTPUT=$("$BEACOND_BINARY" rollback --home="$BEACOND_HOME")
CL=$(echo "$ROLLBACK_OUTPUT" | sed -n 's/.*height=\([0-9]\+\).*/\1/p')
echo "-> CL height: $CL"

# Check if CL is already at or below EL
if (( CL <= EL )); then
    echo "No rollback needed. Consensus Layer height is already at or below Execution Layer height."
    exit 0
fi

# Start the rollback loop from CL down to EL
echo "[Starting rolling back of CL]"
while true; do
    echo "Rolling back CL height $CL..."

    ROLLBACK_OUTPUT=$("$BEACOND_BINARY" rollback --hard --home="$BEACOND_HOME")
    CL=$(echo "$ROLLBACK_OUTPUT" | sed -n 's/.*height=\([0-9]\+\).*/\1/p')
    echo "New CL height after rollback: $CL (required height: $EL)"

    if (( CL <= EL )); then
        echo "Reached target Execution Layer height. Exiting."
        break
    fi
done

echo "Rollback process completed successfully."
