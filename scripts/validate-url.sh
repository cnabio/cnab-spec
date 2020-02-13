#!/bin/sh
###
# Validates that the schema files are served via the https://cnab.io domain,
# as many of the schema contain $refs to each other via URLs on this domain.
#
# If a new schema is added to this repository, a new redirect rule can be
# added to the cnabio/cnab.io repository.
# See https://github.com/cnabio/cnab.io/blob/master/static/_redirects
#
###

set -eou pipefail

for schema in $(ls -1 schema); do
  url="https://cnab.io/v1/${schema}"
  echo "Validating that ${url} returns well-formed json"
  curl -sfL "${url}" | jq > /dev/null 2>&1 \
    || ( echo "${url} doesn't appear to be available or is invalid" && exit 1 )
done
