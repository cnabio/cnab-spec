#!/bin/sh
###
# Use this to validate that the files in examples/* accord with the schema in schema/
#
# To get `ajv`, run `npm install -g ajv-cli`
###

set -eou pipefail

# Test all of the bundle files against the bundle schema
for json in $(ls -1 examples/*-bundle.json); do
  schema="schema/bundle.schema.json"
  echo "Testing json '$json' against schema '$schema'"
  ajv test -s $schema -d $json --valid
done

# Test all of the claim files against the claim schema.
for json in $(ls -1 examples/*-claim.json); do
  schema="schema/claim.schema.json"
  echo "Testing json '$json' against schema '$schema'"
  ajv test -s $schema -d $json --valid -r schema/bundle.schema.json
done