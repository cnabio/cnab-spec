BASE_DIR          := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VALIDATOR_IMG     := deislabs/cnab-spec.ajv

.PHONY: build-validator
build-validator:
	@docker build -f Dockerfile.ajv -t $(VALIDATOR_IMG) .

.PHONY: validate
validate: build-validator
	@docker run --rm \
		-v $(BASE_DIR):/root \
		-w /root \
		$(VALIDATOR_IMG) ./scripts/validate.sh

.PHONY: validate-url
validate-url: build-validator
	@docker run --rm \
		-v $(BASE_DIR):/root \
		-w /root \
		$(VALIDATOR_IMG) ./scripts/validate-url.sh

.PHONY: build-validator-local
build-validator-local:
	@npm install -g ajv-cli

.PHONY: validate-local
validate-local: build-validator-local
	./scripts/validate.sh

.PHONY: validate-url-local
validate-url-local: build-validator-local
	./scripts/validate-url.sh

