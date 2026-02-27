include test.env
DATABASE_URL := $(shell echo $(DATABASE_URL) | sed 's/"//g')
export DATABASE_URL

MAGIC_NUMBER := $(shell echo $(MAGIC_NUMBER) | sed 's/"//g')
export MAGIC_NUMBER

docker-images:
	docker compose up -d --build

test: docker-images
	docker compose restart database # Makes sure the db is empty
	go test ./...

format:
	gofmt -l -w .
