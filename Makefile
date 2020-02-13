BASE_DIR          := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VALIDATOR_IMG     := deislabs/cnab-spec.ajv
VALIDATOR_SCRIPTS := validate.sh validate-url.sh

.PHONY: build-validator
build-validator:
	@docker build -f Dockerfile.ajv -t $(VALIDATOR_IMG) .

.PHONY: validate
validate: build-validator
	@for script in $(VALIDATOR_SCRIPTS); do \
		docker run --rm \
			-v $(BASE_DIR):/root \
			-w /root \
			$(VALIDATOR_IMG) ./scripts/$$script ; \
	done

.PHONY: build-validator-local
build-validator-local:
	@npm install -g ajv-cli

.PHONY: validate-local
validate-local: build-validator-local
	@for script in $(VALIDATOR_SCRIPTS); do \
		./scripts/$$script ; \
	done

