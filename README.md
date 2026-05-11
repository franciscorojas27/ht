# HTTP Trace (HT)

```
 _   _ _____
| | | |_   _|
| |_| | | |
|  _  | | |
|_| |_| |_|
HTTP Trace
```

HT is a command-line HTTP client built for quick testing, debugging, and reproducible HTTP flows. It works with YAML definitions and a fast cURL-like mode.

**Key features**
- Execute multiple requests defined in YAML.
- Variable substitution (`${VAR}`) via `vars` or `env_file`.
- Quick mode with `--url` and cURL-like flags.
- `--dry-run` prints the request without sending it.
- `--verbose` includes detailed timing (DNS, TCP, TLS, first byte).
- Pretty JSON output with colors.
- Global and per-request headers, inline body or `body_file`.

**Installation**
```bash
# from releases (recommended)
# download the binary for your OS and add it to your PATH
```

```bash
# install with Go
go install github.com/franciscorojas27/ht/cmd/ht@latest
```

```bash
# or build locally
go build -o ht ./cmd/ht
```

**Quick mode (cURL-like)**
Quick mode is enabled when you omit `--yml` or when you pass `--url`/`--method`.

```bash
# simple GET
ht --url https://example.com

# POST with body
ht --url https://example.com --data "a=1&b=2"

# multiple headers
ht --url https://example.com --header "Accept: application/json" --header "X-Token: 123"

# HEAD
ht --url https://example.com --head

# basic auth
ht --url https://example.com --user user:pass

# insecure TLS
ht --url https://example.com --insecure
```

**YAML mode**
```bash
# run a request by name
ht --yml http/dev.yml --name get-user

# list requests defined in the YAML
ht --yml http/dev.yml --ls

# print the request without sending it
ht --yml http/dev.yml --name create-post --dry-run

# enable detailed trace
ht --yml http/dev.yml --name get-user --verbose
```

**Main flags**
- `--yml` : Path to the YAML configuration file.
- `--ls` : List requests defined in the YAML file.
- `--name` : Request name to execute.
- `--dry-run` : Print the request without sending it.
- `--verbose` : Enable detailed timing output.
- `--url` : Direct URL for quick mode.
- `--method`, `-X` : HTTP method for quick mode.
- `--header`, `-H` : Extra header (repeatable).
- `--data`, `-d` : Request body for quick mode.
- `--head`, `-I` : Fetch headers only (HEAD).
- `--insecure`, `-k` : Allow insecure TLS connections.
- `--user`, `-u` : Basic auth `user:pass`.

**YAML format**
Basic example (sanitized):

```yaml
config:
  base_url: https://api.example.com
  timeout: 30
  headers:
    Content-Type: application/json; charset=UTF-8

vars:
  token: "MY_TOKEN"

env_file: ./http/.env

requests:
  - name: get-user
    method: GET
    url: /
    headers:
      User-Agent: "MyClient/1.0"

  - name: create-post
    method: POST
    url: /posts
    body:
      title: "hello"
      body: "example"
      userId: 1
```

Main fields:
- `config.base_url` : Optional base URL prepended to each request `url`.
- `config.timeout` : Timeout in seconds (per request).
- `config.headers` : Headers applied to all requests.
- `vars` : Variables that can be referenced as `${var}` in the YAML.
- `env_file` : Path to a `.env` file with `key=value` per line. Values are merged with `vars`.
- `requests[]` : List of requests with `name`, `method`, `url`, `headers`, `body` or `body_file`.

**`env_file` format**
Lines in the form `KEY=value`. Lines starting with `#` and empty lines are ignored.

**License**
MIT
