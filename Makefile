NAME := lambda-go-echo-function

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -extldflags \"-static\""

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	docker-compose run --rm -e CGO_ENABLED=0 go build $(LDFLAGS) -a -tags netgo -installsuffix netgo -o bin/$(NAME) github.com/dtan4/lambda-go-echo-function

.PHONY: deploy
deploy: bin/$(NAME)
ifeq ($(AWS_ACCOUNT_ID),)
	@echo "AWS_ACCOUNT_ID must be set" >&2
	@exit 1
endif
ifeq ($(AWS_S3_BUCKET),)
	@echo "AWS_S3_BUCKET must be set" >&2
	@exit 1
endif
ifeq ($(AWS_CLOUDFORMATION_STACK_NAME),)
	@echo "AWS_CLOUDFORMATION_STACK_NAME must be set" >&2
	@exit 1
endif
	docker-compose run --rm -e AWS_ACCOUNT_ID=$(AWS_ACCOUNT_ID) envsubst < template.yaml.template > template.yaml
	docker-compose run --rm sam package --template-file template.yaml --s3-bucket $(AWS_S3_BUCKET) --output-template-file packaged.yaml
	docker-compose run --rm sam deploy --template-file packaged.yaml --stack-name $(AWS_CLOUDFORMATION_STACK_NAME) --capabilities CAPABILITY_IAM

.PHONY: generate
generate:
	docker-compose run --rm go generate -v ./...

.PHONY: setup
setup: setup-envsubst setup-go setup-sam

.PHONY: setup-envsubst
setup-envsubst:
	docker-compose build envsubst

.PHONY: setup-go
setup-go:
	docker-compose build go

.PHONY: setup-sam
setup-sam:
	docker-compose build sam

.PHONY: test
test:
	docker-compose run --rm go test -coverprofile=coverage.txt -v `docker-compose run -T --rm go list ./... | grep -v mock`
