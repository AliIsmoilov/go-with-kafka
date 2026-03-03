## Overview

The project can be used as a platform for training in performing work tasks. The project has a deliberately poorly designed database subject area. The request path starts with an API written in Golang using SOLID and CLEAN ARCHITECTURE. Upon a successful POST request, the API returns a 202 HTTP status code, and the request goes to a Kafka topic, from where psql-master pulls it and replicates it to psql-replica. The application also uses a Redis cache.

![Go Gopher](https://raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher.png)

---

## What's Inside?

- **`producer/`**: Go application to produce messages to Kafka
- **`consumer/`**: Go application to consume messages from Kafka
- **`psql/`**: Main database
- **`psql replica/`**: Replication of main database
- **`redis/`**: Cache 
- **`promethus/`**: Metrics
- **`grafana/`**: Monitoring
- **`docker-compose.yml`**: Docker setup for local Kafka cluster

---

## Quick Start

### Prerequisites

- **Docker & Docker Compose** - [Get Docker](https://docs.docker.com/get-docker/)
- **Make** - ```sudo apt-get intsall make```
- **golang** - ```sudo apt-get update && sudo apt-get install golang-1.23 && /usr/lib/go-1.23/bin/go version```
  
### Running project with Docker
```bash
git clone https://github.com/FollG/kafka-with-go.git
cd kafka-with-go
go mod init github.com/FollG/kafka-with-go
go mod tidy
make build
make docker-up
```
