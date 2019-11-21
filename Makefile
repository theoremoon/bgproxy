.PHONY: help
help:
	@cat Makefile | grep -E "^[A-Za-z0-9-]+:"

.PHONY: build
build: generate
	go build -o . ./cmd/...

.PHONY: server
server: cmd/pbproxy/main.go
	go build -trimpath ./cmd/pbproxy

.PHONY: client
client: cmd/pbproxyctl/main.go
	go build -trimpath ./cmd/pbproxyctl

.PHONY: generate
generate:
	go generate ./...

.PHONY: install
install:
	go install ./cmd/...

.PHONY: clean
clean:
	rm ./bgproxy ./bgproxyctl
