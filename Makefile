VERSION?=$$(cat version.go | grep VERSION | cut -d"=" -f2 | sed 's/"//g' | sed 's/ //g')
GOFMT_FILES?=$$(find . -name '*.go')
PROJECT_BIN?=jailguard
PROJECT_SRC?=github.com/gasiordev/jailguard

default: build

tools:
	GO111MODULE=off go get -u github.com/gasiordev/go-cli

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@ gofmt_files=$$(gofmt -l $(GOFMT_FILES)); \
	if [[ -n $${gofmt_files} ]]; then \
		echo "The following files fail gofmt:"; \
		echo "$${gofmt_files}"; \
		echo "Run \`make fmt\` to fix this."; \
		exit 1; \
	fi

build: guard-GOPATH
	mkdir -p $$GOPATH/bin/freebsd
	GOOS=freebsd GOARCH=amd64 go build -v -o $$GOPATH/bin/freebsd/${PROJECT_BIN} $$GOPATH/src/${PROJECT_SRC}/*.go

release: build
	mkdir -p $$GOPATH/releases
	tar -cvzf $$GOPATH/releases/${PROJECT_BIN}-${VERSION}-freebsd-amd64.tar.gz -C $$GOPATH/bin/freebsd ${PROJECT_BIN}

.NOTPARALLEL:

.PHONY: tools fmt build

