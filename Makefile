SRC := cmd/app/main.go
EXEC := main

LOGRUS := github.com/sirupsen/logrus github.com/sirupsen/logrus@v1.9.3
CLEANENV := github.com/ilyakaznacheev/cleanenv
GIN := github.com/gin-gonic/gin
GOOSE := github.com/pressly/goose/v3/cmd/goose@latest
PGX := github.com/jackc/pgx github.com/jackc/pgx/v5/pgxpool
SWAG := github.com/swaggo/swag/cmd/swag
GIN_SWAG := github.com/swaggo/gin-swagger github.com/swaggo/files

all: build run

build: clean
	go build -o $(EXEC) $(SRC)

run:
	./$(EXEC)

clean:
	rm -f $(EXEC)

mod:
	go mod init $(EXEC)

get:
	go get $(GIN) \
		$(CLEANENV) \
		$(LOGRUS) \
		$(GOOSE) \
		$(PGX) \
		$(SWAG) \
		$(GIN_SWAG) 

docker-compose-up-silent: docker-compose-stop
	sudo docker compose -f docker-compose.yml up -d

docker-compose-stop:
	sudo docker compose -f docker-compose.yml stop

docker-compose-up: docker-compose-down
	sudo docker compose -f docker-compose.yml up

docker-compose-down:
	sudo docker compose -f docker-compose.yml down

swag:
	swag fmt
	swag init -g common.go -d internal/handler,internal/model