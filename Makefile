GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=btcutil
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
deps:
		$(GOGET) github.com/spf13/viper
		$(GOGET) github.com/urfave/cli
		$(GOGET) github.com/tyler-smith/go-bip39
		$(GOGET) github.com/btcsuite/btcd/chaincfg
		$(GOGET) github.com/btcsuite/btcutil
		$(GOGET) github.com/btcsuite/btclog
		$(GOGET) golang.org/x/crypto/ripemd160

# $(GOGET) github.com/tyler-smith/go-bip32
# $(GOGET) github.com/WeMeetAgain/go-hdwallet

# Cross compilation
# build-linux:
#		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
# docker-build:
#		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
