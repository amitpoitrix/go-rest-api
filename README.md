# students-api (go-rest-api)

Quick notes for local development and pushing to GitHub.

## Run locally

The application expects a configuration YAML file. The config loader looks for the `CONFIG_PATH` environment variable or the `-config` flag.

Options:

- Use the `.env` file in the repo (the app loads it automatically). Ensure `CONFIG_PATH` in `.env` points to a valid config file, e.g.:

```
CONFIG_PATH="config/local.yaml"
```

- Or export the env var in your shell before running:

```bash
export CONFIG_PATH="config/local.yaml"
go run ./cmd/students-api
```

- Or pass the config path via the flag:

```bash
go run ./cmd/students-api -config=config/local.yaml
```

Notes:
- The binary uses `github.com/joho/godotenv` to load `.env` (optional). If `.env` isn't present, the app will continue and fall back to environment variables or flags.

## Common git workflow when push is rejected

If `git push` fails with a non-fast-forward error because the remote has new commits, do:

```bash
git fetch origin
git pull --rebase origin main
# resolve any conflicts, then
git push origin main
```

Alternatively you can merge with `git pull origin main` if you prefer merge commits.

## CI

A simple GitHub Actions workflow is included to build the project on `push` and `pull_request`.

---
Small reference commands used during development:

```bash
# build everything
go build ./...

# run the API
go run ./cmd/students-api -config=config/local.yaml
```
