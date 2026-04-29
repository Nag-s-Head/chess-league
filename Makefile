include test.env
DATABASE_URL := $(shell echo $(DATABASE_URL) | sed 's/"//g')
export DATABASE_URL

MAGIC_NUMBER := $(shell echo $(MAGIC_NUMBER) | sed 's/"//g')
export MAGIC_NUMBER

docker-images:
	docker compose up -d --build

nuke-db:
	docker compose down database
	docker container rm nagsknightschessleaguetestserver-database-1 || docker compose up database

test: docker-images
	docker compose restart database
	go test ./... -timeout=60s

format:
	gofmt -l -w .
