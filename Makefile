-include docker/env
export $(shell sed 's/=.*//' docker/env)

COMMIT_HASH=`git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE=`date -u +%FT%T%z`
LDFLAGS=-ldflags "-X github.com/nawa/cryptoexchange-dashboard/cmd.CommitHash=${COMMIT_HASH} -X github.com/nawa/cryptoexchange-dashboard/cmd.BuildDate=${BUILD_DATE}"
MY_UID = $(shell id -u)
WORKDIR := $(PWD)
COVERAGE_DIR=$(CURDIR)/coverage

build:
	@ echo "-> Building binary ..."
	@ go build ${LDFLAGS} -o bin/cryptoexchange-dashboard main.go
.PHONY: build

linter:
	@ echo "-> Running linters ..."
	@ gometalinter --vendor --config=.gometalinter.json --enable=goimports ./...
.PHONY: linter

mockgen:
	@ echo "-> Generate mocks for tests ..."
	# @ mockgen -source storage/balance.go -package mock_data -destination storage/mock_data/balance_mock.go
	# @ mockgen -source storage/exchange.go -package mock_data -destination storage/mock_data/exchange_mock.go
	# mockgen -source usecase/balance.go -package mock_data -destination usecase/mock_data/balance_mock.go
.PHONY: mockgen

test:
	@ echo "-> Run tests ..."

	go test ./...
.PHONY: test

test-coverage:
	@ echo "-> Running tests with coverage ..."
	@rm -rf $(COVERAGE_DIR)
	@mkdir -p $(COVERAGE_DIR)

	@go list ./... | grep -v "/testdata" | grep -v "/mocks" | xargs -I {} mkdir -p $(COVERAGE_DIR)/{}
	@go list ./... | grep -v "/testdata" | grep -v "/mocks" | xargs -I {} go test -v -coverprofile $(COVERAGE_DIR)/{}/cover.out $(GOTEST_PARAM) {}

	@echo "mode: set" > $(COVERAGE_DIR)/coverage-total.out
	@go list ./... | grep -v "/testdata" | grep -v "/mocks" | xargs -I {} cat $(COVERAGE_DIR)/{}/cover.out {} 2>/dev/null | grep -v "mode: set" >> $(COVERAGE_DIR)/coverage-total.out

	@go tool cover -func=$(COVERAGE_DIR)/coverage-total.out | tail -n 1 | xargs -I {} echo "TOTAL COVERAGE. "{}

.PHONY: coverage-gen

coverage-open:
	@ echo "-> Opening coverage report ..."
	go tool cover -html=$(COVERAGE_DIR)/coverage-total.out -o $(COVERAGE_DIR)/coverage-total.html
	open $(COVERAGE_DIR)/coverage-total.html
.PHONY: coverage-open

docker-image-build-x86:
	@ echo "-> Building Docker image ..."
	docker rmi -f $(CRWI_IMAGENAME_X86):bak || true
	docker tag $(CRWI_IMAGENAME_X86) $(CRWI_IMAGENAME_X86):bak || true
	docker rmi -f $(CRWI_IMAGENAME_X86) || true
	docker run --rm -v "$(WORKDIR)":/go/src/github.com/nawa/cryptoexchange-dashboard -w /go/src/github.com/nawa/cryptoexchange-dashboard $(CRWI_BUILDER_IMAGE) /bin/bash -c "CGO_ENABLED=0 GOOS=linux make build && chown -R $(MY_UID) bin"
	docker build -f $(WORKDIR)/$(CRWI_DOCKERFILE_X86) -t $(CRWI_IMAGENAME_X86) $(WORKDIR)

docker-image-build-armhf:
	@ echo "-> Building Docker image ..."
	docker rmi -f $(CRWI_IMAGENAME_ARMHF):bak || true
	docker tag $(CRWI_IMAGENAME_ARMHF) $(CRWI_IMAGENAME_ARMHF):bak || true
	docker rmi -f $(CRWI_IMAGENAME_ARMHF) || true
	docker run --rm -v "$(WORKDIR)":/go/src/github.com/nawa/cryptoexchange-dashboard -w /go/src/github.com/nawa/cryptoexchange-dashboard $(CRWI_BUILDER_IMAGE) /bin/bash -c "CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 make build && chown -R $(MY_UID) bin"
	docker build -f $(WORKDIR)/$(CRWI_DOCKERFILE_ARMHF) -t $(CRWI_IMAGENAME_ARMHF) $(WORKDIR)

#make docker-compose-x86 DCO_ARGS="up -d"
docker-compose-x86:
	docker-compose --file $(WORKDIR)/$(CRWI_DCO_FILE_X86) $(DCO_ARGS)

docker-compose-armhf:
	docker-compose --file $(WORKDIR)/$(CRWI_DCO_FILE_ARMHF) $(DCO_ARGS)

docker-rmi-x86:
	docker rmi $(CRWI_IMAGENAME_X86)

docker-rmi-armhf:
	docker rmi $(CRWI_IMAGENAME_ARMHF)

docker-run-sync:
	docker rm $(CRWI_SYNC_CONTAINER_NAME) || true
	docker run -d --name $(CRWI_SYNC_CONTAINER_NAME) $(CRWI_SYNC_ENVIRONMENT) $(CRWI_SYNC_RESTART) $(CRWI_IMAGENAME_X86)

docker-stop-sync:
	docker stop $(CRWI_SYNC_CONTAINER_NAME)

docker-start-sync:
	docker start $(CRWI_SYNC_CONTAINER_NAME)

docker-rmf-sync:
	docker rm -f $(CRWI_SYNC_CONTAINER_NAME)
