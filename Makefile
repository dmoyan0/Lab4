DIRECTOR_DOCKER_IMAGE = director-server
DATANODE1_DOCKER_IMAGE = datanode1-server
DATANODE2_DOCKER_IMAGE = datanode2-server
DATANODE3_DOCKER_IMAGE = datanode3-server
MERCENARY_DOCKER_IMAGE = mercenarios-server
DOSHBANK_DOCKER_IMAGE = doshbank-server
NAMENODE_DOCKER_IMAGE = namenode-server

DIRECTOR_PORT = 50050
DATANODE1_PORT = 50052
DATANODE2_PORT = 50053
DATANODE3_PORT = 50054
MERCENARY_PORT = 50055
DOSHBANK_PORT = 50056
NAMENODE_PORT = 50051

# Define the default target
all: help

# Build the director server Docker image
docker-director:
    docker build -t $(DIRECTOR_DOCKER_IMAGE) --build-arg SERVER_TYPE=Director .
    docker run -d --name $(DIRECTOR_DOCKER_IMAGE) -e SERVER_TYPE=Director -p $(DIRECTOR_PORT):$(DIRECTOR_PORT) $(DIRECTOR_DOCKER_IMAGE)

# Build the datanode server Docker images
docker-datanode:
    @echo "SERVER_TYPE is set to: $(SERVER_TYPE)"
    @if [ "$(SERVER_TYPE)" = "datanode1" ]; then \
        docker build -t $(DATANODE1_DOCKER_IMAGE) --build-arg SERVER_TYPE=datanode1 .; \
        docker run -d --name $(DATANODE1_DOCKER_IMAGE) -e SERVER_TYPE=datanode1 -p $(DATANODE1_PORT):$(DATANODE1_PORT) $(DATANODE1_DOCKER_IMAGE); \
    elif [ "$(SERVER_TYPE)" = "datanode2" ]; then \
        docker build -t $(DATANODE2_DOCKER_IMAGE) --build-arg SERVER_TYPE=datanode2 .; \
        docker run -d --name $(DATANODE2_DOCKER_IMAGE) -e SERVER_TYPE=datanode2 -p $(DATANODE2_PORT):$(DATANODE2_PORT) $(DATANODE2_DOCKER_IMAGE); \
    elif [ "$(SERVER_TYPE)" = "datanode3" ]; then \
        docker build -t $(DATANODE3_DOCKER_IMAGE) --build-arg SERVER_TYPE=datanode3 .; \
        docker run -d --name $(DATANODE3_DOCKER_IMAGE) -e SERVER_TYPE=datanode3 -p $(DATANODE3_PORT):$(DATANODE3_PORT) $(DATANODE3_DOCKER_IMAGE); \
	else \
        echo "Invalid SERVER_TYPE argument. Use 'datanode1' or 'datanode2' or 'datanode3'."; \
        exit 1; \
    fi

# Build the mercenary server Docker image
docker-mercenary:
    docker build -t $(MERCENARY_DOCKER_IMAGE) --build-arg SERVER_TYPE=mercenary .
    docker run -d --name $(MERCENARY_DOCKER_IMAGE) -e SERVER_TYPE=mercenary -p $(MERCENARY_PORT):$(MERCENARY_PORT) $(MERCENARY_DOCKER_IMAGE)

# Build the dosh bank server Docker image
docker-doshbank:
    docker build -t $(DOSHBANK_DOCKER_IMAGE) --build-arg SERVER_TYPE=doshbank .
    docker run -d --name $(DOSHBANK_DOCKER_IMAGE) -e SERVER_TYPE=doshbank -p $(DOSHBANK_PORT):$(DOSHBANK_PORT) $(DOSHBANK_DOCKER_IMAGE)

# Build the namenode server Docker image
docker-namenode:
    docker build -t $(NAMENODE_DOCKER_IMAGE) --build-arg SERVER_TYPE=namenode .
    docker run -d --name $(NAMENODE_DOCKER_IMAGE) -e SERVER_TYPE=namenode -p $(NAMENODE_PORT):$(NAMENODE_PORT) $(NAMENODE_DOCKER_IMAGE)

# Usage: make help
help:
    @echo "Available targets:"
    @echo "  docker-director  - Build and run the director server Docker container"
    @echo "  docker-datanode  SERVER_TYPE={datanode1,datanode2,datanode3}  - Build and run the specified datanode server Docker container"
    @echo "  docker-mercenary  - Build and run the mercenary server Docker container"
    @echo "  docker-doshbank   - Build and run the dosh bank server Docker container"
    @echo "  docker-namenode   - Build and run the namenode server Docker container"
    @echo "  help              - Display this help message"
