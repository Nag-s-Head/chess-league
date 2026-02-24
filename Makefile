include test.env

docker-images:
	docker compose up -d --build

test: docker-images
	go test ./...
