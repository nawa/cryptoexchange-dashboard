COMMIT_HASH=`git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE=`date -u +%FT%T%z`
LDFLAGS=-ldflags "-X github.com/nawa/cryptoexchange-wallet-info/cmd.CommitHash=${COMMIT_HASH} -X github.com/nawa/cryptoexchange-wallet-info/cmd.BuildDate=${BUILD_DATE}"
LINTER_EXCLUDES=--exclude="comment\s+or\s+be\s+unexported"

build:
	@ echo "-> Building binary ..."
	go build ${LDFLAGS} -o cryptoexchange-wallet-info main.go
.PHONY: build

linter:
	@ echo "-> Running linters ..."
	@ gometalinter --vendor ${LINTER_EXCLUDES} ./...
.PHONY: linter