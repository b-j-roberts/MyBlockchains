# Use golang:1.19 as the base image
FROM golang:1.19

# Set the working directory in the container
WORKDIR /app

# Copy the Go code file to the container
COPY Makefile .
COPY src/ src/
COPY go.mod .
COPY go.sum .

# Copy configs
COPY config/ config/
#TODO: Remove this
COPY accounts/ accounts/

# Define a build argument for the node type
ARG NODE_TYPE
ENV NODE_TYPE=$NODE_TYPE

# Compile the Go file based on the NODE_TYPE build argument
RUN make bin/$NODE_TYPE-node

EXPOSE 8985
EXPOSE 8986
EXPOSE 8987

#TODO: Snapshots
#TODO: Use yaml for configuration and just use a command to copy the node to a universal name to prevent this
#CMD ["./$NODE_TYPE-node", "account-id", "12345"]
CMD if [ "$NODE_TYPE" = "transaction" ] ; \
    then echo $PEER_TO_CALL > config/peer-list-mempool.txt && ./bin/transaction-node; \
    elif [ "$NODE_TYPE" = "miner" ] ; \
    then echo $PEER_TO_CALL > config/peer-list-mempool.txt && ./bin/miner-node --account-id 123456789 ; \
    elif [ "$NODE_TYPE" = "rpc" ] ; \
    then echo $PEER_TO_CALL > config/peer-list-mempool.txt && ./bin/rpc-node --account-id 123456789 ; \
    else echo "Invalid node type specified: $NODE_TYPE (transaction,miner,rpc)" ; \
    fi
