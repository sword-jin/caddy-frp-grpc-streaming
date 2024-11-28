CADDY_VERSION := 2.8.4
BUF_VERSION = 1.43.0
PROTOC_GEN_CONNECT_GO_VERSION = 1.17.0
PROTOC_GEN_GO_VERSION = 1.35.2

BIN_DIR := $(shell pwd)/bin

CADDY := $(BIN_DIR)/caddy
BUF = $(PWD)/bin/buf
RUN_BUF = PATH=$(PWD)/bin:$$PATH $(BUF)
PROTOC_GEN_GO = $(PWD)/bin/protoc-gen-go
PROTOC_GEN_CONNECT_GO = $(PWD)/bin/protoc-gen-connect-go

protoc:
	@echo "Generating Go files"
	cd proto && protoc --go_out=plugins=grpc:. *.proto

build:
	docker build -t grpc-server:local .

generate: $(BUF)
	$(RUN_BUF) generate

generate-certs:
	openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=example Inc./CN=example.local' -keyout example.local.key -out example.local.crt

run-client:
	go run cmd/client/main.go

$(CADDY):
	mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install github.com/caddyserver/caddy/v2/cmd/caddy@v$(CADDY_VERSION)

$(PROTOC_GEN_GO):
	mkdir -p bin
	GOBIN=$(PWD)/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@v$(PROTOC_GEN_GO_VERSION)

$(PROTOC_GEN_CONNECT_GO):
	mkdir -p bin
	GOBIN=$(PWD)/bin go install connectrpc.com/connect/cmd/protoc-gen-connect-go@v$(PROTOC_GEN_CONNECT_GO_VERSION)

$(BUF): $(PROTOC_GEN_GO) $(PROTOC_GEN_CONNECT_GO)
	mkdir -p bin
	GOBIN=$(PWD)/bin go install github.com/bufbuild/buf/cmd/buf@v$(BUF_VERSION)

clean:
	rm -rf bin
	rm *.crt
	rm *.key
	rm *.csr
