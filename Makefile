include test.env
DATABASE_URL := $(shell echo $(DATABASE_URL) | sed 's/"//g')
export DATABASE_URL

MAGIC_NUMBER := $(shell echo $(MAGIC_NUMBER) | sed 's/"//g')
export MAGIC_NUMBER

docker-images:
	docker compose up -d --build

nuke-db:
	docker compose down database
	docker container rm nagsknightschessleaguetestserver-database-1 || docker compose up database -d

test: docker-images
	go test ./... -timeout=60s

gofmt:
	gofmt -l -w .

format: gofmt
	pnpm i
	pnpm format || true # prettier and go templates do not play that well together. but at least it does some formatting

psql:
	docker compose exec -it database psql -U magnus -d chess-league
