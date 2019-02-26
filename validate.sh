#!/bin/sh
###
# Use this to validate that the files in examples/* accord with the schema in schema/
#
# To get `ajv`, run `npm install -g ajv-cli`
###

set -eou pipefail

for json in $(ls -1 examples/*); do
  for schema in $(ls -1 schema/*); do
    echo "Testing json '$json' against schema '$schema'"
    ajv test -s $schema -d $json --valid
  done
done
