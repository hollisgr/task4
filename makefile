EXEC := task4
SRC := cmd/app/main.go

CLEANENV := github.com/ilyakaznacheev/cleanenv
GIN := github.com/gin-gonic/gin
PGX := github.com/jackc/pgx github.com/jackc/pgx/v5/pgxpool
CORS := github.com/gin-contrib/cors
SQUIRREL := github.com/Masterminds/squirrel

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
	go get \
		$(GIN) \
		$(CLEANENV) \
		$(PGX) \
		$(CORS) \
		$(SQUIRREL)

docker-compose-up-silent: docker-compose-stop
	sudo docker compose -f docker-compose.yml --env-file=config.env up -d

docker-compose-stop:
	sudo docker compose -f docker-compose.yml --env-file=config.env stop

docker-compose-up: docker-compose-down
	sudo docker compose -f docker-compose.yml --env-file=config.env up

docker-compose-down:
	sudo docker compose -f docker-compose.yml --env-file=config.env down