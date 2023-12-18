#!/bin/bash
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
