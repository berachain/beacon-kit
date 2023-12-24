# SPDX-License-Identifier: MIT
#
# Copyright (c) 2023 Berachain Foundation
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

if [ -z "$HOMEDIR" ]; then
    HOMEDIR="/.polard"
fi
if [ -z "$KEYRING" ]; then
    KEYRING="test"
fi
if [ -z "$KEYALGO" ]; then
    KEYALGO="secp256k1"
fi

polard genesis collect-gentxs --home "$HOMEDIR"

polard genesis validate-genesis --home "$HOMEDIR"

# # faucet
# polard keys add faucet --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"
# polard genesis add-genesis-account faucet 1000000000000000000000000000abera,1000000000000000000000000000stgusdc --keyring-backend $KEYRING --home "$HOMEDIR"

# # # Test Account
# # absurd surge gather author blanket acquire proof struggle runway attract cereal quiz tattoo shed almost sudden survey boring film memory picnic favorite verb tank
# # 0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306
# polard genesis add-genesis-account cosmos1yrene6g2zwjttemf0c65fscg8w8c55w58yh8rl 1000000000000000000000000000abera,1000000000000000000000000000stgusdc --keyring-backend $KEYRING --home "$HOMEDIR"
