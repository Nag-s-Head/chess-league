include test.env

docker-images:
	docker compose up -d --build

test: docker-images
	docker compose restart database # Makes sure the db is empty
	go test ./...
