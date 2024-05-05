#!/bin/bash

# Check if the name argument is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <migration_name>"
    exit 1
fi

# Get the current timestamp in the format YYYYMMDDHHMMSS
timestamp=$(date +'%Y%m%d%H%M%S')

# Convert the provided name to snake case
snake_case_name=$(echo "$1" | tr '[:upper:]' '[:lower:]' | tr ' ' '_')

# Create the migration file with the correct name and extension
filename="db/migrations/${timestamp}_${snake_case_name}"

# Create an empty migration file with the appropriate name
touch "$filename.up.sql"
touch "$filename.down.sql"
