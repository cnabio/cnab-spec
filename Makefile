BASE_DIR         := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VALIDATOR_IMG    := deislabs/cnab-spec.ajv
VALIDATOR_SCRIPT := ./validate.sh

.PHONY: build-validator
build-validator:
	@docker build -f Dockerfile.ajv -t $(VALIDATOR_IMG) .

.PHONY: validate
validate: build-validator
	@docker run --rm \
		-v $(BASE_DIR):/root \
		-w /root \
		$(VALIDATOR_IMG) $(VALIDATOR_SCRIPT)

.PHONY: build-validator-local
build-validator-local:
	@npm install -g ajv-cli

.PHONY: validate-local
validate-local: build-validator-local
	@$(VALIDATOR_SCRIPT)

