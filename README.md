# chess-league

[Now live!](https://nagsknights.co.uk)

[![Go CI](https://github.com/Nag-s-Head/chess-league/actions/workflows/go.yml/badge.svg)](https://github.com/Nag-s-Head/chess-league/actions/workflows/go.yml)

Want to run a chess league? - use this.

## Developer Guide

This is a Go application that hosts a web server that serves dynamically generated HTML that Go template
renders, HTMX is used extensively to keep forms nice and easy to manage, and PostgresQL is used as a
database, all truth is stored in the database.

### Tooling

You need the following tools instealled to use this project:

- pnpm (for JS parts)
- go (>= 1.26)
- GNU make
- docker

This project uses Docker to create, and tear down production like environments. You can use a non-Docker setup,
but that is not the intended dev environment.

### Scripts

- To start a local environment that updates when changes are made use the following command:
  `docker compose up --watch`

- To build an executable use the following command:
  `make build`

- To execute all of the tests use the following command:
  `make test -j`

- To reset the db use the following command:
  `make nuke-db`

- To format the codebase use the following command:
  `make format -j`

## Environment Variables

| Variable            | Usage                                                                         | Example                                                                                       |
| ------------------- | ----------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| DATABASE_URL        | The full URL of the Postgres database                                         | `user=magnus password=bong-cloud dbname=chess-league host=database port=5432 sslmode=disable` |
| MAGIC_NUMBER        | The magic number that clients need to submit a game                           | 3743289472-does-not-actually-need-to-be-a-number                                              |
| GITHUB_ORGANISATION | The full name of the organisation                                             | `Nag-s-Head`, as seen in our repo's URL `https://github.com/Nag-s-Head/chess-league`          |
| GITHUB_API_KEY      | A personal access token to allow read of the private org members              | `github_pat_...`                                                                              |
| OAUTH_CLIENT_ID     | Used for admin portal authentication, created under Github developer settings | `1234...`                                                                                     |
| OAUTH_CLIENT_SECRET | Used for admin portal authentication, created under Github developer settings | `1234...`                                                                                     |
| APP_BASE_URL        | The base URL of the application, used for OAuth redirects                     | `https://nagsknights.co.uk`                                                                   |
| TEST_MODE           | When enabled this uses a mocked Oauth implementation                          | `true` or omit to disable                                                                     |

## AI Policy

Whilst AI coding agents have helped to increase the accessibility of development work, it remains very important
that any code written is of high quality, well tested, and deoes not violate other licences. I think an AI ban is
a negative and reactive policy, and that proper code review, and disciplined use of AI agents is a good way to
ensure it is used responibly. So by all means do use any AI use see fit, but be aware that is might cause issues,
and will work best if you can guide the agent to writing quality code with a good architecture.
