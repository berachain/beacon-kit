#!/bin/bash

# Set your L1 values here.
L1_RPC_URL="http://localhost:64064"  # Replace with your L1 node URL
L1_RPC_KIND="basic"
L1_CHAINID=80087
L2_CHAINID=42069
# TODO: L1_block_time, l1 finalization period

# Fill out environment variables in .env file
cd ~/op-stack-deployment/optimism
direnv allow
cp .envrc.example .envrc

# TODO: update the values L1_RPC_URL, L1_RPC_KIND in envrc
direnv allow

sh ./packages/contracts-bedrock/scripts/getting-started/wallets.sh
# TODO: update the Getting Started values in envrc
direnv allow

echo "Sending 10 ether to admin, proposer, batcher addresses..."
cast send --private-key fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306 $GS_ADMIN_ADDRESS --value 100ether --rpc-url $L1_RPC_URL
cast send --private-key fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306 $GS_BATCHER_ADDRESS --value 100ether --rpc-url $L1_RPC_URL
cast send --private-key fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306 $GS_PROPOSER_ADDRESS --value 100ether --rpc-url $L1_RPC_URL

# Update deploy-config/getting-started.json with new addresses
cd packages/contracts-bedrock
sh ./scripts/getting-started/config.sh

# Get L1 Info
output=$(cast block finalized | grep -E "(timestamp|hash|number)")
# Parse the output using awk and store the values in variables
hash=$(echo "$output" | awk '/hash/ { print $2 }')
number=$(echo "$output" | awk '/number/ { print $2 }')
timestamp=$(echo "$output" | awk '/timestamp/ { print $2 }')

# Print the variables
echo "Hash: $hash"
echo "Number: $number"
echo "Timestamp: $timestamp"

# Update deploy-config/getting-started.json file with the values
# TODO: make sure this works.
awk -v hash="$hash" '/"l1StrtingBlockTag": "BLOCKHASH"/{$0="    \"l1StartingBlockTag\": \"" hash "\", "}1' deploy-config/getting-started.json > temp && mv temp deploy-config/getting-started.json
awk -v timestamp="$timestamp" '/"l2OutputOracleStartingTimestamp": TIMESTAMP/{$0="    \"l2OutputOracleStartingTimestamp\": " timestamp ", "}1' deploy-config/getting-started.json > temp && mv temp deploy-config/getting-started.json
awk -v L1_CHAINID="$L1_CHAINID" '/"l1ChainID": L1_CHAINID/{$0="    \"l1ChainID\": " L1_CHAINID ", "}1' deploy-config/getting-started.json > temp && mv temp deploy-config/getting-started.json
awk -v L2_CHAINID="$L2_CHAINID" '/"l2ChainID": L2_CHAINID/{$0="    \"l2ChainID\": " L2_CHAINID ", "}1' deploy-config/getting-started.json > temp && mv temp deploy-config/getting-started.json
# TODO: L1_block_time, l1 finalization period

# Print the updated JSON file
echo "deploy-config/getting-started.json"
cat deploy-config/getting-started.json

# Step 4: Deploy L1 smart contracts
forge script -vvv scripts/Deploy.s.sol:Deploy --private-key $GS_ADMIN_PRIVATE_KEY --broadcast --legacy --rpc-url $L1_RPC_URL
forge script -vvv scripts/Deploy.s.sol:Deploy --sig 'sync()' --private-key $admin_private_key --broadcast --legacy --rpc-url $L1_RPC_URL

# TODO: 
# - Update the L1 contract addresses in the deployments/getting-started/l1.json
# - OR figure out why the forge script didn't automatically save to a json

# Step 5: Run the OP node genesis
cd ~/op-stack-deployment/optimism/op-node

go run cmd/main.go genesis l2 \
  --deploy-config ../packages/contracts-bedrock/deploy-config/getting-started.json \
  --l1-deployments ../packages/contracts-bedrock/deployments/getting-started/l1.json \
  --outfile.l2 genesis.json \
  --outfile.rollup rollup.json \
  --l1-rpc $L1_RPC_URL

openssl rand -hex 32 > jwt.txt

cp genesis.json ~/op-stack-deployment/op-geth
cp jwt.txt ~/op-stack-deployment/op-geth

# Step 6: Build OP Geth
cd ~/op-stack-deployment/op-geth
mkdir datadir
build/bin/geth init --datadir=datadir genesis.json
