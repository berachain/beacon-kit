# Read the JSON file into a variable
json=$(cat build/scripts/constants.json)

# Use jq to extract all keys and values
keys=$(echo $json | jq -r 'keys[]')
entries=$(echo $json | jq -r 'to_entries[]')
key_to_lookup="block_roots.size"

for file in $(find ./types -name '*.pb.go'); do
    echo $file
    for key in $keys; do
        echo $key
        echo $value
        value=$(echo $json | jq -r ".\"$key\"")
        sed -i '' "s/$key/$value/g" $file
    done
done
