#!/bin/bash

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
