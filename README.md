## Overview

A modern training platform demonstrating event-driven architecture with Golang, built on SOLID and CLEAN ARCHITECTURE principles. The system processes requests asynchronously—successful POST requests return 202 status and are published to Kafka topics, where they are consumed and persisted to PostgreSQL with automatic replication between master and replica instances. The application includes Redis caching for performance optimization and integrated monitoring via Prometheus and Grafana.

<!-- ![Go Gopher](https://raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher.png) -->

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
