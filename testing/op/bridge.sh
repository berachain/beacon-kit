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

# Get the bridge contract from the deployment files
BRIDGE_CONTRACT=$(jq -r '.L1StandardBridgeProxy' ~/op-stack-deployment/optimism/packages/contracts-bedrock/deployments/getting-started/l1.json)
printf "\nBridge contract address: $BRIDGE_CONTRACT\n"

cd ~/op-stack-deployment/optimism
direnv allow 
source .envrc

# Send 10 ETH to the bridge contract to receive on the L2
printf "\nSending 10 ETH to the bridge contract...\n"
cast send $BRIDGE_CONTRACT --value 10ether --private-key $PRIVATE_KEY --legacy --rpc-url $L1_RPC_URL
