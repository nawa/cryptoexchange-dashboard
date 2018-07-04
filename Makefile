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
	@ gometalinter --vendor --skip frontend --config=.gometalinter.json --enable=goimports ./...
.PHONY: linter

mockgen:
	@ echo "-> Generate mocks for tests ..."
	mockgen -source storage/balance.go -package mocks -destination storage/mocks/balance_mock.go
	mockgen -source storage/exchange.go -package mocks -destination storage/mocks/exchange_mock.go
	mockgen -source usecase/balance.go -package mocks -destination usecase/mocks/balance_mock.go
	mockgen -source usecase/order.go -package mocks -destination usecase/mocks/order_mock.go
.PHONY: mockgen

unit-test:
	@ echo "-> Run unit tests ..."

	go test -v ./...
.PHONY: unit-test

unit-test-coverage:
	@ echo "-> Running unit tests with coverage ..."
	@rm -rf $(COVERAGE_DIR)
	@mkdir -p $(COVERAGE_DIR)

	@go list ./... | grep -v "/testdata" | grep -v "/mocks" | xargs -I {} mkdir -p $(COVERAGE_DIR)/{}
	@go list ./... | grep -v "/testdata" | grep -v "/mocks" | xargs -I {} go test -v -coverprofile $(COVERAGE_DIR)/{}/cover.out $(GOTEST_PARAM) {}

	@echo "mode: set" > $(COVERAGE_DIR)/coverage-total.out
	@go list ./... | grep -v "/testdata" | grep -v "/mocks" | xargs -I {} cat $(COVERAGE_DIR)/{}/cover.out {} 2>/dev/null | grep -v "mode: set" >> $(COVERAGE_DIR)/coverage-total.out

	@go tool cover -func=$(COVERAGE_DIR)/coverage-total.out | tail -n 1 | xargs -I {} echo "TOTAL COVERAGE. "{}

.PHONY: unit-test-coverage

test:
	@ echo "-> Run all tests ..."

	docker-compose --file $(WORKDIR)/int-tests/env/docker-compose.yml down; \
    	docker-compose --file $(WORKDIR)/int-tests/env/docker-compose.yml up -d; \
		DB_TEST_URL=localhost:27019/crexd-test go test -v -tags=integration_test ./...; \
		status=$$?; \
		docker-compose --file $(WORKDIR)/int-tests/env/docker-compose.yml down; \
		exit $$status

.PHONY: test

test-coverage:
	@ echo "-> Running all tests with coverage ..."
	@rm -rf $(COVERAGE_DIR)
	@mkdir -p $(COVERAGE_DIR)

	@go list ./... | grep -v "/testdata" | grep -v "/mocks" | xargs -I {} mkdir -p $(COVERAGE_DIR)/{}

	docker-compose --file $(WORKDIR)/int-tests/env/docker-compose.yml down; \
        	docker-compose --file $(WORKDIR)/int-tests/env/docker-compose.yml up -d; \
        	export DB_TEST_URL=localhost:27019/crexd-test; \
    		go list ./... | grep -v "/testdata" | grep -v "/mocks" | xargs -I {} go test -v -tags=integration_test -coverprofile $(COVERAGE_DIR)/{}/cover.out $(GOTEST_PARAM) {}; \
    		status=$$?; \
    		docker-compose --file $(WORKDIR)/int-tests/env/docker-compose.yml down; \
    		echo "mode: set" > $(COVERAGE_DIR)/coverage-total.out; \
    		go list ./... | grep -v "/testdata" | grep -v "/mocks" | xargs -I {} cat $(COVERAGE_DIR)/{}/cover.out {} 2>/dev/null | grep -v "mode: set" >> $(COVERAGE_DIR)/coverage-total.out; \
    		go tool cover -func=$(COVERAGE_DIR)/coverage-total.out | tail -n 1 | xargs -I {} echo "TOTAL COVERAGE. "{}; \
    		exit $$status

.PHONY: test-coverage

coverage-open:
	@ echo "-> Opening coverage report ..."
	go tool cover -html=$(COVERAGE_DIR)/coverage-total.out -o $(COVERAGE_DIR)/coverage-total.html
	open $(COVERAGE_DIR)/coverage-total.html
.PHONY: coverage-open

docker-image-build-x86:
	@ echo "-> Building Docker image ..."
	docker rmi -f $(CREXD_IMAGENAME_X86):bak || true
	docker tag $(CREXD_IMAGENAME_X86) $(CREXD_IMAGENAME_X86):bak || true
	docker rmi -f $(CREXD_IMAGENAME_X86) || true
	docker run --rm -v "$(WORKDIR)":/go/src/github.com/nawa/cryptoexchange-dashboard -w /go/src/github.com/nawa/cryptoexchange-dashboard $(CREXD_BUILDER_IMAGE) /bin/bash -c "CGO_ENABLED=0 GOOS=linux make build && chown -R $(MY_UID) bin"
	docker build -f $(WORKDIR)/$(CREXD_DOCKERFILE_X86) -t $(CREXD_IMAGENAME_X86) $(WORKDIR)

docker-image-build-armhf:
	@ echo "-> Building Docker image ..."
	docker rmi -f $(CREXD_IMAGENAME_ARMHF):bak || true
	docker tag $(CREXD_IMAGENAME_ARMHF) $(CREXD_IMAGENAME_ARMHF):bak || true
	docker rmi -f $(CREXD_IMAGENAME_ARMHF) || true
	docker run --rm -v "$(WORKDIR)":/go/src/github.com/nawa/cryptoexchange-dashboard -w /go/src/github.com/nawa/cryptoexchange-dashboard $(CREXD_BUILDER_IMAGE) /bin/bash -c "CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 make build && chown -R $(MY_UID) bin"
	docker build -f $(WORKDIR)/$(CREXD_DOCKERFILE_ARMHF) -t $(CREXD_IMAGENAME_ARMHF) $(WORKDIR)

#make docker-compose-x86 DCO_ARGS="up -d"
docker-compose-x86:
	docker-compose --file $(WORKDIR)/$(CREXD_DCO_FILE_X86) $(DCO_ARGS)

docker-compose-armhf:
	docker-compose --file $(WORKDIR)/$(CREXD_DCO_FILE_ARMHF) $(DCO_ARGS)

docker-rmi-x86:
	docker rmi $(CREXD_IMAGENAME_X86)

docker-rmi-armhf:
	docker rmi $(CREXD_IMAGENAME_ARMHF)

docker-run-sync:
	docker rm $(CREXD_SYNC_CONTAINER_NAME) || true
	docker run -d --name $(CREXD_SYNC_CONTAINER_NAME) $(CREXD_SYNC_ENVIRONMENT) $(CREXD_SYNC_RESTART) $(CREXD_IMAGENAME_X86)

docker-stop-sync:
	docker stop $(CREXD_SYNC_CONTAINER_NAME)

docker-start-sync:
	docker start $(CREXD_SYNC_CONTAINER_NAME)

docker-rmf-sync:
	docker rm -f $(CREXD_SYNC_CONTAINER_NAME)
