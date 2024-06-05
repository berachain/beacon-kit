#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
#
# Copyright (c) 2024 Berachain Foundation
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


# Define the root directory of your Go project
ROOT_DIR="."

if ! command -v golines &> /dev/null; then
    echo "⚠️ 'golines' not found. Installing..."
    go install github.com/segmentio/golines@latest
    if [ $? -ne 0 ]; then
        echo "Failed to install golines. Manually install to run 'make format'."
        exit 1
    else
        echo "✅ 'golines' successfully installed! Running..."
    fi
else
    echo "✅ 'golines' is already installed. Running..."
fi

# Find all .go files in the project directory and its subdirectories, ignoring .pb.go and .pb_encoding.go files
find "${ROOT_DIR}" -type f -name "*.go" ! -name "*.pb.go" ! -name "*.pb_encoding.go" | while read -r file; do
    echo "Processing $file..."
    golines --reformat-tags --shorten-comments --write-output --max-len=80 "$file"
done

echo "✅ All files processed."
