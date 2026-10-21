# Repository guide

- This is a single Go module (`github.com/Nag-s-Head/chess-league`): `main.go` starts the app, `app/` wires the database and HTTP handlers, and handler HTML/assets are embedded into the binary.
- PostgreSQL is required for DB-backed tests and app startup. `make test -j` is the CI test command: it starts the Docker Compose stack and runs `go generate ./...` before `go test ./...`. The Makefile loads `test.env`; for a focused test, start the DB with `docker compose up -d database`, load the env (`set -a; . ./test.env; set +a`), then run `go test ./<package> -run '^TestName$'`.
- `go generate ./...` runs `pnpm i`, copies HTMX scripts, and rebuilds `handlers/assets/tailwind.css`; run it after changing frontend sources before building/testing. `make build` also runs generation.
- Database schema changes belong in `db/migrations/` (migration definitions are in Go); the app applies migrations when connecting. The migration tests also launch a disposable PostgreSQL Docker container on port 5444.
- `make format -j` runs `gofmt` and Prettier across the repo; it installs pnpm dependencies and tolerates Prettier failures on Go templates, so review its diff rather than assuming it formatted everything cleanly.
