#!/bin/sh
# A simple script to generate static HTML documentation given an authenticated GraphQL endpoint.

if ! [ -x "$(command -v graphdoc)" ]; then
    echo "Please install graphdoc as per the instructions at https://github.com/2fd/graphdoc."
    exit 1
fi

if [ $# -ne 3 ]; then
    echo "Usage: $0 <graphql_endpoint> <bearer_token> <output_dir>"
    exit 1
fi

graphdoc -e "${1}" -x "Authorization: Bearer ${2}" -o "${3}"
