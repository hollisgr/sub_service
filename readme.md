# Simple CRUDL app

Simple CRUD application built with Golang using Gin framework, PGX library for PostgreSQL connectivity, Swag for generating API documentation, and Logrus for structured logging.

## Overview

This project demonstrates a typical CRUDL implementation with:

- **Gin**: Fast web framework providing routing capabilities.
- **PGX**: High-performance driver for connecting to PostgreSQL.
- **Swag**: Automatic API documentation generator.
- **Logrus**: Structured logger with flexible output formats.
- **Goose**: Database migrations.

## Requirements

- **Go**: Version 1.20+.
- **PostgreSQL**: Version 12+. Set up a database before running the application.
- **Goose**: Version 3.24+.
- **Docker** (optional): For easier setup with Docker images.

## Directory structure

- **cmd**: Contains the main application entry point.
- **build**: Contains the docker images.
- **docs**: Auto-generated Swagger documentation using Swag.
- **internal/config**: Holds configurations including database settings.
- **internal/models**: Defines the entity models.
- **internal/handler**: Handlers for managing API endpoints.
- **internal/subscription**: Service for managing subscription entity.
- **migrations**: This directory stores database migrations.
- **pkg**: Helper utilities like database connections and logging.

## Getting Started

### 1. Cloning the Repository

```
git clone https://github.com/hollisgr/sub_service.git
cd sub_service
```

### 2. Setting Up the Environment

Update the config variables in `.env` file.

Example `.env` content:

```
BIND_IP=127.0.0.1
LISTEN_PORT=8000
PSQL_HOST=127.0.0.1
PSQL_PORT=5432
PSQL_NAME=sub_service
PSQL_USER=postgres
PSQL_PASSWORD=password
```

### 3. Running the Application

database migrations:
```
goose -dir migrations postgres "host=${PSQL_HOST} port=${PSQL_PORT} dbname=${PSQL_NAME} user=${PSQL_USER} password=${PSQL_PASSWORD} sslmode=disable" up
```

run app:
```
go build -o sub_service cmd/app/main.go
./sub_service
```

or run with makefile
```
make build && make run
```

## 4. Running with Docker:

Update the config variables in `.env` file.

Example `.env` content:

```
BIND_IP=127.0.0.1
LISTEN_PORT=8000
PSQL_HOST=psql-db
PSQL_PORT=25432
PSQL_NAME=sub_service
PSQL_USER=postgres
PSQL_PASSWORD=password
```

```
make docker-compose-up
```

Or silent mode:

```
make docker-compose-up-silent
```

## Swagger Docs

Auto-generated API documentation is available at `http://127.0.0.1:8000/swagger/index.html`. These documents are generated dynamically using the Swag package during compilation.
