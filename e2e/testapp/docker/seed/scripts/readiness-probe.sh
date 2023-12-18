#!/bin/bash
# Execute the cURL command and capture the response
response=$(curl -s -X POST -H "Content-Type: application/json" -d '{
    "jsonrpc":"2.0",
    "method": "eth_syncing",
    "params": [],
    "id": 1
}' http://localhost:8545)

# Check if the response contains the "result" field
if echo "$response" | grep -q '"result":.*false'; then
  echo "Syncing is not in progress"
  exit 0 # Exit with success code
else
  echo "Syncing is in progress or port is not up"
  exit 1 # Exit with failed code
fi
