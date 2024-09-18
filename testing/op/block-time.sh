#!/bin/bash

# Function to convert RFC3339 timestamp to Unix epoch
rfc3339_to_epoch() {
  gdate -d "$1" +%s
}

# Get the latest block height
latest_block_info=$(curl -s http://localhost:26657/block | jq '.result.block.header.height')
latest_height=$(echo "$latest_block_info" | tr -d '"')

# Initialize variables
total_time_diff=0
prev_timestamp=""

# Loop over the past 100 blocks
for (( i=0; i<30; i++ ))
do
  block_height=$((latest_height - i))
  block_info=$(curl -s http://localhost:26657/block?height=$block_height | jq '.result.block.header.time')
  echo "$block_info"
  block_timestamp=$(echo "$block_info" | tr -d '"')

  if [ -n "$prev_timestamp" ]; then
    # Convert timestamps to Unix epoch
    block_time_epoch=$(rfc3339_to_epoch "$block_timestamp")
    prev_time_epoch=$(rfc3339_to_epoch "$prev_timestamp")

    # Calculate the time difference
    time_diff=$((prev_time_epoch - block_time_epoch))
    total_time_diff=$((total_time_diff + time_diff))
  fi

  prev_timestamp=$block_timestamp
done

# Calculate the average block time
avg_block_time=$(echo "scale=2; $total_time_diff / 29" | bc)

# Output the average block time
echo "Average block time for the past 30 blocks: $avg_block_time seconds"
