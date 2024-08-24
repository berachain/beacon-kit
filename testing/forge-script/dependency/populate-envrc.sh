#!/bin/bash

# Function to parse YAML file
function parse_yaml() {
   local prefix=$2
   local s='[[:space:]]*'
   local w='[a-zA-Z0-9_]*'
   local fs=$(echo @|tr @ '\034')
   sed -ne "s|^\($s\):|\1|" \
        -e "s|^\($s\)\($w\)$s:$s\"\(.*\)\"$s\$|\1$fs\2$fs\3|p" \
        -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p" $1 |
   awk -F$fs '{
      indent = length($1)/2;
      vname[indent] = $2;
      for (i in vname) {if (i > indent) {delete vname[i]}}
      if (length($3) > 0) {
         vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
         printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
      }
   }'
}

# Path to YAML file
config_file="config-env.yaml"

# Path to .envrc file
envrc_file=".envrc"

# Clear the existing .envrc file or create a new one
> $envrc_file

# Parse the YAML file and append to .envrc
eval $(parse_yaml $config_file)

# Write the variables to the .envrc file
echo "export HONEY=\"$HONEY\"" >> $envrc_file
echo "export PYTH=\"$PYTH\"" >> $envrc_file
echo "export FEE_COLLECTOR=\"$FEE_COLLECTOR\"" >> $envrc_file
echo "export API_KEY_ROUTESCAN=\"$API_KEY_ROUTESCAN\"" >> $envrc_file
echo "export RPC_URL=\"$RPC_URL\"" >> $envrc_file
echo "export GOV=\"$GOV\"" >> $envrc_file
echo "export GOV_PK=\"$GOV_PK\"" >> $envrc_file
echo "export DEPOSITOR=\"$DEPOSITOR\"" >> $envrc_file
echo "export DEP_PK=\"$DEP_PK\"" >> $envrc_file

# Notify user of completion
echo "The .envrc file has been populated with the values from $config_file."
