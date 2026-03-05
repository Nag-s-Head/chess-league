# chess-league

[Now live!](https://nagsknights.co.uk)

[![Go CI](https://github.com/Nag-s-Head/chess-league/actions/workflows/go.yml/badge.svg)](https://github.com/Nag-s-Head/chess-league/actions/workflows/go.yml)

Want to run a chess league? - use this.

## Environment Variables

| Variable     | Usage                                               | Example                                                                                       |
| ------------ | --------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| DATABASE_URL | The full URL of the Postgres database               | `user=magnus password=bong-cloud dbname=chess-league host=database port=5432 sslmode=disable` |
| MAGIC_NUMBER | The magic number that clients need to submit a game | 3743289472-does-not-actually-need-to-be-a-number                                              |
