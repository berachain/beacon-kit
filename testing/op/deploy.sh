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


# Set your L1 values here.
PRIV_KEY="fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306"  # Replace with a funded wallet key on the L1.
RPC_URL="127.0.0.1:54787"  # Replace with your L1 node URL.
RPC_KIND="any" 
CHAINID=80088 # Replace with your L1 chain id.
BLOCK_TIME=6 # Replace with your L1 block time.
FINALIZATION_PERIOD=6 # Replace with your L1 block time.

# Fill out environment variables in .env file
cd ~/op-stack-deployment/optimism
direnv allow
cp .envrc.example .envrc

# Update the PRIVATE_KEY, L1_RPC_URL, and L1_RPC_KIND in the .envrc
if sed --version 2>&1 | grep -q GNU; then
  sed -i 's/^export PRIVATE_KEY=.*/export PRIVATE_KEY='"$PRIV_KEY"'/' .envrc
  sed -i 's/^export L1_RPC_URL=.*/export L1_RPC_URL='"$RPC_URL"'/' .envrc
  sed -i 's/^export L1_RPC_KIND=.*/export L1_RPC_KIND='"$RPC_KIND"'/' .envrc
else 
  sed -i '' 's/^export PRIVATE_KEY=.*/export PRIVATE_KEY='"$PRIV_KEY"'/' .envrc
  sed -i '' 's/^export L1_RPC_URL=.*/export L1_RPC_URL='"$RPC_URL"'/' .envrc
  sed -i '' 's/^export L1_RPC_KIND=.*/export L1_RPC_KIND='"$RPC_KIND"'/' .envrc
fi
direnv allow

# Update the Getting Started wallets in the .envrc file.
wallets=$(sh ./packages/contracts-bedrock/scripts/getting-started/wallets.sh)

update_envrc() {
  local key="$1"
  local value="$2"
  local escaped_value=$(echo "$value" | sed 's_/_\\/_g')
  if sed --version 2>&1 | grep -q GNU; then
    sed -i '' "s/^export $key=.*/export $key=$escaped_value/" .envrc
  else 
    sed -i '' "s/^export $key=.*/export $key=$escaped_value/" .envrc
  fi
}

echo "$wallets" | while IFS= read -r line; do
  if [[ "$line" =~ ^export\ (.*)=(.*)$ ]]; then
    key="${BASH_REMATCH[1]}"
    value="${BASH_REMATCH[2]}"
    update_envrc "$key" "$value"
  fi
done
direnv allow 

echo "Sending 10 ether to admin, proposer, batcher addresses..."
cast send --private-key $PRIVATE_KEY $GS_ADMIN_ADDRESS --value 10ether --rpc-url $L1_RPC_URL --legacy
cast send --private-key $PRIVATE_KEY $GS_BATCHER_ADDRESS --value 10ether --rpc-url $L1_RPC_URL --legacy
cast send --private-key $PRIVATE_KEY $GS_PROPOSER_ADDRESS --value 10ether --rpc-url $L1_RPC_URL --legacy

# # Update deploy-config/getting-started.json with new addresses
# cd packages/contracts-bedrock
# sh ./scripts/getting-started/config.sh

# # Get L1 Info
# output=$(cast block finalized | grep -E "(timestamp|hash|number)")
# # Parse the output using awk and store the values in variables
# hash=$(echo "$output" | awk '/hash/ { print $2 }')
# number=$(echo "$output" | awk '/number/ { print $2 }')
# timestamp=$(echo "$output" | awk '/timestamp/ { print $2 }')

# # Print the variables
# echo "Hash: $hash"
# echo "Number: $number"
# echo "Timestamp: $timestamp"

# # Update deploy-config/getting-started.json file with the values
# L2_CHAINID=42069

# # TODO: make sure this works.
# awk -v hash="$hash" '/"l1StrtingBlockTag": "BLOCKHASH"/{$0="    \"l1StartingBlockTag\": \"" hash "\", "}1' deploy-config/getting-started.json > temp && mv temp deploy-config/getting-started.json
# awk -v timestamp="$timestamp" '/"l2OutputOracleStartingTimestamp": TIMESTAMP/{$0="    \"l2OutputOracleStartingTimestamp\": " timestamp ", "}1' deploy-config/getting-started.json > temp && mv temp deploy-config/getting-started.json
# # TODO: set $number to "l2OutputOracleStartingBlockNumber"
# awk -v L1_CHAINID="$L1_CHAINID" '/"l1ChainID": L1_CHAINID/{$0="    \"l1ChainID\": " L1_CHAINID ", "}1' deploy-config/getting-started.json > temp && mv temp deploy-config/getting-started.json
# awk -v L2_CHAINID="$L2_CHAINID" '/"l2ChainID": L2_CHAINID/{$0="    \"l2ChainID\": " L2_CHAINID ", "}1' deploy-config/getting-started.json > temp && mv temp deploy-config/getting-started.json
# # TODO: set L1_BLOCK_TIME to "l1BlockTime" 
# # TODO: set L1_FINALIZATION_PERIOD to "finalizationPeriodSeconds"

# # Print the updated JSON file
# echo "deploy-config/getting-started.json"
# cat deploy-config/getting-started.json

# # Step 4: Deploy L1 smart contracts
# forge script -vvv scripts/Deploy.s.sol:Deploy --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --legacy --rpc-url $L1_RPC_URL
# forge script -vvv scripts/Deploy.s.sol:Deploy --sig 'sync()' --private-key $admin_private_key --broadcast --legacy --rpc-url $L1_RPC_URL

# # TODO: 
# # - Update the L1 contract addresses in the deployments/getting-started/l1.json
# # - OR figure out why the forge script didn't automatically save to a json

# # Step 5: Run the OP node genesis
# cd ~/op-stack-deployment/optimism/op-node

# go run cmd/main.go genesis l2 \
#   --deploy-config ../packages/contracts-bedrock/deploy-config/getting-started.json \
#   --l1-deployments ../packages/contracts-bedrock/deployments/getting-started/l1.json \
#   --outfile.l2 genesis.json \
#   --outfile.rollup rollup.json \
#   --l1-rpc $L1_RPC_URL

# openssl rand -hex 32 > jwt.txt

# cp genesis.json ~/op-stack-deployment/op-geth
# cp jwt.txt ~/op-stack-deployment/op-geth

# # Step 6: Build OP Geth
# cd ~/op-stack-deployment/op-geth
# mkdir datadir
# build/bin/geth init --datadir=datadir genesis.json
