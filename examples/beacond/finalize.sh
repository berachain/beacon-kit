#!/bin/bash
set -x

# Collect genesis tx
/usr/bin/beacond genesis collect-gentxs --home "$BEACOND_HOME" > /dev/null 2>&1

# Run this to ensure everything worked and that the genesis file is setup correctly
/usr/bin/beacond genesis validate-genesis --home "$BEACOND_HOME" > /dev/null 2>&1