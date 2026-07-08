# active-cookies

A small CLI that reads a cookie log file and prints the most active cookie(s) for a given day (UTC).

## Requirements

- [Go](https://go.dev/) 1.25+
- [Make](https://www.gnu.org/software/make/)
- [Docker](https://www.docker.com/) (for `make format`, `make check`, `make test`, and Docker image targets)

## Cookie log format

The input file is a CSV with a header row:

```csv
cookie,timestamp
AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-09T10:13:00+00:00
```

## Run without Docker

Build the binary:

```bash
make build
```

Run it with a log file and date:

```bash
./bin/main -f cookie_log.csv -d 2018-12-09
```

Or use the Makefile shortcut:

```bash
make run FILE=cookie_log.csv DATE=2018-12-09
```

Both flags are required:

- `-f` — path to the cookie log file
- `-d` — date to query (`YYYY-MM-DD`, UTC)

## Run with Docker

Build the image (compiles a Linux binary into `bin/main`, then packages it):

```bash
make docker-build
```

Run with a log file from your machine. Mount a host directory so the container can read it:

```bash
docker run --rm -v "$(pwd):/data" activecookie:latest \
  -f /data/cookie_log.csv -d 2018-12-09
```

Or use the Makefile helper:

```bash
make docker-run FILE=cookie_log.csv DATE=2018-12-09
```

Output is printed to stdout. Save it on the host:

```bash
docker run --rm -v "$(pwd):/data" activecookie:latest \
  -f /data/cookie_log.csv -d 2018-12-09 > results.txt
```

## Build and push to Docker Hub

Build the image:

```bash
make docker-build
```

Push to Docker Hub (replace `your-dockerhub-username`):

```bash
make docker-push DOCKER_USER=your-dockerhub-username
```

That tags and pushes `your-dockerhub-username/activecookie:latest`.

Others can pull and run it:

```bash
docker pull your-dockerhub-username/activecookie:latest
docker run --rm -v "$(pwd):/data" your-dockerhub-username/activecookie:latest \
  -f /data/cookie_log.csv -d 2018-12-09
```

## Development commands

These targets use a dev Docker image (`Dockerfile.dev`) so everyone runs the same tooling.

| Command | What it does |
|---------|--------------|
| `make format` | Formats Go source files with `go fmt`. |
| `make check` | Runs static analysis (`go vet` and `golangci-lint`). |
| `make test` | Runs all unit tests. |
| `make cover` | Runs `make check`, then reports test coverage. |

Other useful targets:

```bash
make clean        # remove binaries, tmp files, and local Docker images
```
