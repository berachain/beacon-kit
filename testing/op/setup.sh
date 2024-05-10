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

# Clone Optimism repos
cd ~/
mkdir op-stack-deployment
cd ~/op-stack-deployment
git clone --branch tutorials/chain --single-branch https://github.com/ethereum-optimism/optimism.git
git clone https://github.com/ethereum-optimism/op-geth.git --depth 1
 
# Check if ~/.nvm directory doesn't exist
if [ ! -d "$nvm_dir" ]; then
    mkdir "$nvm_dir"
    echo "Created ~/.nvm directory."
fi
. ~/.nvm/nvm.sh
. ~/.zshrc
. $(brew --prefix nvm)/nvm.sh  # if installed via Brew
nvm install v20.11.0

# Display the versions of the required packages
printf "\nMinimum versions of the required packages:\n"
sh ./packages/contracts-bedrock/scripts/getting-started/versions.sh

# Build op-node op-batcher op-proposer
cd ~/op-stack-deployment/optimism
npm i -g pnpm
pnpm install
make op-node op-batcher op-proposer
pnpm build

# Install and build op-geth
cd ~/op-stack-deployment/op-geth/
make geth

# Install direnv
cd ~/op-stack-deployment
brew install direnv
