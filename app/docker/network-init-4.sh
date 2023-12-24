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

CONTAINER0="polard-node0"
CONTAINER1="polard-node1"
CONTAINER2="polard-node2"
CONTAINER3="polard-node3"

HOMEDIR="/.polard"
SCRIPTS="/scripts"

rm -rf ./temp
mkdir ./temp
mkdir ./temp/seed0
mkdir ./temp/seed1
mkdir ./temp/seed2
mkdir ./temp/seed3
touch ./temp/genesis.json

# init step 1 
docker exec $CONTAINER0 bash -c "$SCRIPTS/seed0-init-step1.sh"
docker exec $CONTAINER1 bash -c "$SCRIPTS/seed1-init-step1.sh seed-1"
docker exec $CONTAINER2 bash -c "$SCRIPTS/seed1-init-step1.sh seed-2"
docker exec $CONTAINER3 bash -c "$SCRIPTS/seed1-init-step1.sh seed-3"

# copy genesis.json from seed-0 to seed-1
docker cp $CONTAINER0:$HOMEDIR/config/genesis.json ./temp/genesis.json
docker cp ./temp/genesis.json $CONTAINER1:$HOMEDIR/config/genesis.json

# init step 2
docker exec $CONTAINER1 bash -c "$SCRIPTS/seed1-init-step2.sh seed-1"

# copy genesis.json from seed-1 to seed-2
docker cp $CONTAINER1:$HOMEDIR/config/genesis.json ./temp/genesis.json
docker cp ./temp/genesis.json $CONTAINER2:$HOMEDIR/config/genesis.json

# init step 2
docker exec $CONTAINER2 bash -c "$SCRIPTS/seed2-init-step2.sh seed-2"

# copy genesis.json from seed-2 to seed-3
docker cp $CONTAINER2:$HOMEDIR/config/genesis.json ./temp/genesis.json
docker cp ./temp/genesis.json $CONTAINER3:$HOMEDIR/config/genesis.json

# init step 2
docker exec $CONTAINER3 bash -c "$SCRIPTS/seed1-init-step2.sh seed-3"


# copy genesis.json from seed-3 to seed-0
docker cp $CONTAINER3:$HOMEDIR/config/genesis.json ./temp/genesis.json
docker cp ./temp/genesis.json $CONTAINER0:$HOMEDIR/config/genesis.json

# copy gentx
docker cp $CONTAINER1:$HOMEDIR/config/gentx ./temp
docker cp $CONTAINER2:$HOMEDIR/config/gentx ./temp
docker cp $CONTAINER3:$HOMEDIR/config/gentx ./temp
docker cp ./temp/gentx $CONTAINER0:$HOMEDIR/config

# init step 2
docker exec $CONTAINER0 bash -c "$SCRIPTS/seed0-init-step2.sh"

# copy genesis.json from seed-0 to seed-1,2,3
docker cp $CONTAINER0:$HOMEDIR/config/genesis.json ./temp/genesis.json
docker cp ./temp/genesis.json $CONTAINER1:$HOMEDIR/config/genesis.json
docker cp ./temp/genesis.json $CONTAINER2:$HOMEDIR/config/genesis.json
docker cp ./temp/genesis.json $CONTAINER3:$HOMEDIR/config/genesis.json

# start
# docker exec -it $CONTAINER0 bash -c "$SCRIPTS/seed-start.sh"
# docker exec -it $CONTAINER1 bash -c "$SCRIPTS/seed-start.sh"
# docker exec -it $CONTAINER2 bash -c "$SCRIPTS/seed-start.sh"
# docker exec -it $CONTAINER3 bash -c "$SCRIPTS/seed-start.sh"

# docker exec -it polard-node0 bash -c "/scripts/seed-start.sh"
# docker exec -it polard-node1 bash -c "/scripts/seed-start.sh"
# docker exec -it polard-node2 bash -c "/scripts/seed-start.sh"
# docker exec -it polard-node3 bash -c "/scripts/seed-start.sh"
