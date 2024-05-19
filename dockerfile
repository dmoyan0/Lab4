# Use the specified base image
ARG BASE_IMAGE=golang:1.21.3
FROM ${BASE_IMAGE} AS builder

ARG DIRECTOR_PORT= 50050
ARG DATANODE1_PORT= 50052
ARG DATANODE2_PORT= 50053
ARG DATANODE3_PORT= 50054
ARG MERCENARY_PORT = 50055
ARG DOSHBANK_PORT= 50056
ARG NAMENODE_PORT= 50051

ARG SERVER_TYPE

# Set the working directory inside the container
WORKDIR /app

# Copy the parent directory's go.mod and go.sum files to the container
COPY go.mod .
COPY go.sum .

# Download and install Go dependencies
RUN go mod download

# Copy the rest of your application code to the container
COPY . .

CMD if [ "$SERVER_TYPE" = "Director" ]; then \
        PORT=$DIRECTOR_PORT; \
        cd /app/NameNode; \
        go build -o director-server; \
        ./director-server; \
    elif [ "$SERVER_TYPE" = "Datanode1" ]; then \
        PORT=$DATANODE1_PORT; \
        cd /app/Datanode1; \
        go build -o datanode1-server; \
        ./datanode1-server; \
    elif [ "$SERVER_TYPE" = "Datanode2" ]; then \
        PORT=$DATANODE2_PORT; \
        cd /app/Datanode2; \
        go build -o datanode2-server; \
        ./datanode2-server; \
    elif [ "$SERVER_TYPE" = "Datanode3" ]; then \
        PORT=$DATANODE3_PORT; \
        cd /app/Datanode3; \
        go build -o datanode3-server; \
        ./datanode3-server; \
    elif [ "$SERVER_TYPE" = "NameNode" ]; then \
        PORT=$NAMENODE_PORT; \
        cd /app/NameNode; \
        go build -o namenode-server; \
        ./namenode-server; \
    elif [ "$SERVER_TYPE" = "Dosh_Bank" ]; then \
        PORT=$DOSHBANK_PORT; \
        cd /app/Dosh_Bank; \
        go build -o doshbank-server; \
        ./doshbank-server; \
    elif [ "$SERVER_TYPE" = "mercenarios" ]; then \
        PORT=$MERCENARY_PORT; \
        cd /app/mercenarios; \
        go build -o mercenarios-server; \
        ./mercenarios-server; \
    else \
        echo "Invalid SERVER_TYPE argument."; \
    fi