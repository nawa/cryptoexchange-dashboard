-include docker/env
export $(shell sed 's/=.*//' docker/env)

COMMIT_HASH=`git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE=`date -u +%FT%T%z`
LDFLAGS=-ldflags "-X github.com/nawa/cryptoexchange-wallet-info/cmd.CommitHash=${COMMIT_HASH} -X github.com/nawa/cryptoexchange-wallet-info/cmd.BuildDate=${BUILD_DATE}"
MY_UID = $(shell id -u)
WORKDIR := $(PWD)

build:
	@ echo "-> Building binary ..."
	@ go build ${LDFLAGS} -o bin/cryptoexchange-wallet-info main.go
.PHONY: build

linter:
	@ echo "-> Running linters ..."
	@ gometalinter --vendor --config=.gometalinter.json --enable=goimports ./...
.PHONY: linter

docker-image-build-x86:
	@ echo "-> Building Docker image ..."
	docker rmi -f $(CRWI_IMAGENAME_X86):bak || true
	docker tag $(CRWI_IMAGENAME_X86) $(CRWI_IMAGENAME_X86):bak || true
	docker rmi -f $(CRWI_IMAGENAME_X86) || true
	docker run --rm -v "$(WORKDIR)":/go/src/github.com/nawa/cryptoexchange-wallet-info -w /go/src/github.com/nawa/cryptoexchange-wallet-info $(CRWI_BUILDER_IMAGE) /bin/bash -c "CGO_ENABLED=0 GOOS=linux make build && chown -R $(MY_UID) bin"
	docker build -f $(WORKDIR)/$(CRWI_DOCKERFILE_X86) -t $(CRWI_IMAGENAME_X86) $(WORKDIR)

docker-image-build-armhf:
	@ echo "-> Building Docker image ..."
	docker rmi -f $(CRWI_IMAGENAME_ARMHF):bak || true
	docker tag $(CRWI_IMAGENAME_ARMHF) $(CRWI_IMAGENAME_ARMHF):bak || true
	docker rmi -f $(CRWI_IMAGENAME_ARMHF) || true
	docker run --rm -v "$(WORKDIR)":/go/src/github.com/nawa/cryptoexchange-wallet-info -w /go/src/github.com/nawa/cryptoexchange-wallet-info $(CRWI_BUILDER_IMAGE) /bin/bash -c "CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 make build && chown -R $(MY_UID) bin"
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