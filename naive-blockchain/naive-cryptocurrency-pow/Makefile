# Makefile for building Go client nodes

# Name of the executable files to be built
MINER_NODE = miner-node
TRANSACTION_NODE = transaction-node
RPC_NODE = rpc-node

# Directories
SRC_DIR = src
BUILD_DIR = build
BIN_DIR = bin

# Compiler and linker options
GO = go
GOFLAGS = -ldflags="-s -w" # Reduce binary size
LDFLAGS = -s -w # Strip debug symbols and reduce binary size

# Source files for each executable
MINER_SRC = $(SRC_DIR)/miner-node.go
TRANSACTION_SRC = $(SRC_DIR)/transaction-node.go
RPC_SRC = $(SRC_DIR)/rpc-node.go

# Check if environment variables are set
guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
  fi

# Targets
.PHONY: all clean help debug account

all: $(BIN_DIR)/$(MINER_NODE) $(BIN_DIR)/$(TRANSACTION_NODE) $(BIN_DIR)/$(RPC_NODE)
docker-build: docker-build/${MINER_NODE} docker-build/${TRANSACTION_NODE} docker-build/${RPC_NODE}
docker-push: docker-push/${MINER_NODE} docker-push/${TRANSACTION_NODE} docker-push/${RPC_NODE}
docker-all: docker-build docker-push

# Clean up intermediate build files and executable files
clean:
		rm -rf $(BUILD_DIR) $(BIN_DIR)

# Build rules for each executable
$(BIN_DIR)/$(MINER_NODE): $(MINER_SRC)
		mkdir -p $(BUILD_DIR)
		#GO111MODULE=off $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(MINER_NODE) $<
		$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(MINER_NODE) $<
		mkdir -p $(BIN_DIR)
		cp $(BUILD_DIR)/$(MINER_NODE) $(BIN_DIR)/

$(BIN_DIR)/$(TRANSACTION_NODE): $(TRANSACTION_SRC)
		mkdir -p $(BUILD_DIR)
		$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(TRANSACTION_NODE) $<
		mkdir -p $(BIN_DIR)
		cp $(BUILD_DIR)/$(TRANSACTION_NODE) $(BIN_DIR)/

$(BIN_DIR)/$(RPC_NODE): $(RPC_SRC)
		mkdir -p $(BUILD_DIR)
		$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(RPC_NODE) $<
		mkdir -p $(BIN_DIR)
		cp $(BUILD_DIR)/$(RPC_NODE) $(BIN_DIR)/

docker-build/${MINER_NODE}: $(MINER_SRC)
		docker build --build-arg NODE_TYPE=miner -t brandonjroberts/naive-${MINER_NODE} .

docker-build/${RPC_NODE}: ${RPC_SRC}
		docker build --build-arg NODE_TYPE=rpc -t brandonjroberts/naive-${RPC_NODE} .

docker-build/${TRANSACTION_NODE}: ${TRANSACTION_SRC}
		docker build --build-arg NODE_TYPE=transaction -t brandonjroberts/naive-${TRANSACTION_NODE} .

docker-push/${MINER_NODE}: docker-build/${MINER_NODE}
		docker push brandonjroberts/naive-${MINER_NODE}

docker-push/${RPC_NODE}: docker-build/${RPC_NODE}
		docker push brandonjroberts/naive-${RPC_NODE}

docker-push/${TRANSACTION_NODE}: docker-build/${TRANSACTION_NODE}
		docker push brandonjroberts/naive-${TRANSACTION_NODE}

# Help message for Makefile targets
help:
		@echo "Makefile targets:"
		@echo "  all                    Builds all executables (default target)"
		@echo "  bin/<node-type>-node   Builds executables for <node-type>-node such as miner-node, transaction-node, rpc-node"
		@echo "  clean                  Removes all intermediate build files and executables"
		@echo "  account                Creates an account key pair"
		@echo "  debug                  Builds all executables with debugging information"
		@echo "  help                   Shows this help message"

# Build all executables with debugging information
debug: GOFLAGS = -gcflags="all=-N -l"
debug: LDFLAGS = -X 'main.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)'
debug: clean all

account: guard-ACCOUNT_ID
		@echo "Generating key pair for account $(ACCOUNT_ID)"
		$(GO) run $(SRC_DIR)/scripts/generateAccount.go --accounts accounts/ --account-id $(ACCOUNT_ID)

# Default target is all
.DEFAULT_GOAL := all
