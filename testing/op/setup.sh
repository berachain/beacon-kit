#!/bin/bash

# Stage 1: Clone repos
cd ~/
mkdir op-stack-deployment
cd op-stack-deployment
git clone -b tutorials/chain https://github.com/ethereum-optimism/optimism.git
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

# Stage 3: Install op-node op-batcher op-proposer
cd optimism

# TODO: handle any deps not installed
sh ./packages/contracts-bedrock/scripts/getting-started/versions.sh

npm i -g pnpm
pnpm install
make op-node op-batcher op-proposer
pnpm build

cd packages/contracts-bedrock/
forge install

cd ..

# Stage 4: Install op-geth
cd op-geth
make geth
cd ..

# Stage 5: Install direnv
brew install direnv
direnv_hook='eval "$(direnv hook zsh)"'
zsh_config="$HOME/.zshrc"

# Check if the direnv hook already exists in the file
if ! grep -qF "$direnv_hook" "$zsh_config"; then
    # Append the direnv hook to the file
    echo "$direnv_hook" >> "$zsh_config"
    echo "direnv hook added to $zsh_config"
else
    echo "direnv hook already exists in $zsh_config"
fi
source $zsh_config
