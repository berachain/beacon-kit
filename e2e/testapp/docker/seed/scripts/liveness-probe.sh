#!/bin/bash
# Execute the cURL command and capture the response
response=$(curl -s -X POST -H "Content-Type: application/json" -d '{
    "jsonrpc":"2.0",
    "method": "eth_blockNumber",
    "params": [],
    "id": 1
}' http://localhost:8545)

height=$(echo "$response" | jq -r '.result')

file="last_height.txt"

# Check if the file exists
if [ -e "$file" ]; then
  # Read the contents of the file into the result variable
  last_height=$(cat "$file")
else
  # File does not exist, set result to an empty string
  last_height=""
fi

rm $file
echo "$height" >> $file

# Check if the two input strings are equal
if [ "$height" == "$last_height" ]; then
  # Strings are equal, return 1
  exit 1
else
  # Strings are not equal, return 0
  exit 0
fi
