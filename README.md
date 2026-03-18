# chess-league

[Now live!](https://nagsknights.co.uk)

[![Go CI](https://github.com/Nag-s-Head/chess-league/actions/workflows/go.yml/badge.svg)](https://github.com/Nag-s-Head/chess-league/actions/workflows/go.yml)

Want to run a chess league? - use this.

## Environment Variables

| Variable            | Usage                                                                         | Example                                                                                       |
| ------------------- | ----------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| DATABASE_URL        | The full URL of the Postgres database                                         | `user=magnus password=bong-cloud dbname=chess-league host=database port=5432 sslmode=disable` |
| MAGIC_NUMBER        | The magic number that clients need to submit a game                           | 3743289472-does-not-actually-need-to-be-a-number                                              |
| GITHUB_ORGANISATION | The full name of the organisation                                             | `Nag-s-Head`, as seen in our repo's URL `https://github.com/Nag-s-Head/chess-league`          |
| GITHUB_API_KEY      | A personal access token to allow private org members                          | `github_pat_...`                                                                              |
| OAUTH_CLIENT_ID     | Used for admin portal authentication, created under Github developer settings | `1234...`                                                                                     |
| OAUTH_CLIENT_SECRET | Used for admin portal authentication, created under Github developer settings | `1234...`                                                                                     |
| APP_BASE_URL        | The base URL of the application, used for OAuth redirects                     | `https://nagsknights.co.uk`                                                                   |
| GITHUB_ORGANISATION | The Github organisation name                                                  | `Nag-s-Head`                                                                                  |
| TEST_MODE           | When enabled this uses a mocked Oauth implementation                          | `true` or omit to disable                                                                     |
