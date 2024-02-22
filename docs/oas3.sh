#!/bin/bash
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


# Function to run oas3-mdx command on a file
run_oas3_mdx_on_file() {
    echo "Running oas3-mdx on $1"
    oas3-mdx --spec "$1" --target docs/pages --templates docs/templates
}

# Function to iterate over files in a directory
process_directory() {
    local directory="$1"
    local files=$(find "$directory" -type f \( -name "*.yaml" -o -name "*.json" \))

    for file in $files; do
        run_oas3_mdx_on_file "$file"
    done
}

# Main script starts here
if [ $# -eq 0 ]; then
    echo "Usage: $0 <directory>"
    exit 1
fi

# Check if oas3-mdx is installed
if ! command -v oas3-mdx &> /dev/null; then
    echo "Error: oas3-mdx command not found. Make sure it's installed and in your PATH."
    exit 1
fi

# Process the specified directory
process_directory "$1"
