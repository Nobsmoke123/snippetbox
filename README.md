# Snippetbox

Snippetbox is a server-rendered web application written in Go for creating and managing text snippets.
It includes user authentication, CSRF protection, secure session handling, and a PostgreSQL-backed persistence layer.

## Features

- Create, view, and list text snippets.
- User signup, login, logout, and account pages.
- Protected routes for authenticated-only actions.
- Password update workflow from account settings.
- Secure session storage in PostgreSQL via `scs` + `pgxstore`.
- CSRF protection using `nosurf`.
- Embedded HTML templates and static assets using Go `embed`.
- Structured request/error logging with `slog`.
- Optional debug mode for detailed 500 responses.

## Tech Stack

- **Language:** Go 1.26+
- **Router/HTTP:** Standard library `net/http` (`http.ServeMux` patterns)
- **Database:** PostgreSQL (`pgx/v5` connection pool)
- **Sessions:** `github.com/alexedwards/scs/v2` + PostgreSQL store
- **Forms:** `github.com/go-playground/form/v4`
- **Security:** `nosurf` (CSRF), secure session cookies, security headers, TLS
- **Frontend:** Server-rendered templates + static CSS/JS

## Project Structure

```text
.
├── cmd/web/                  # Application entrypoint, handlers, middleware, routing
├── internal/models/          # Data access layer for snippets/users
│   └── testdata/             # SQL setup/teardown fixtures
├── internal/validator/       # Form validation helpers
├── ui/
│   ├── html/                 # Base, partials, and page templates
│   ├── static/               # CSS/JS assets
│   └── efs.go                # Embedded file system declaration
├── .air.toml                 # Optional live-reload config
└── go.mod                    # Go module dependencies
```

## Prerequisites

Before running locally, make sure you have:

- Go `1.26` or newer
- PostgreSQL (local or remote)
- OpenSSL (or another way to generate local TLS certs)

## Environment Variables

The app expects a `.env` file at the project root and reads:

- `DATABASE_URL` (required): PostgreSQL DSN used by `pgxpool`

Example:

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/snippetbox?sslmode=disable
```

> Note: The application exits at startup if `.env` is missing.

## Database Setup

Create a PostgreSQL database, then execute the schema from `internal/models/testdata/setup.sql`.

Example using `psql`:

```bash
createdb snippetbox
psql "postgres://postgres:postgres@localhost:5432/snippetbox?sslmode=disable" -f internal/models/testdata/setup.sql
```

Tables created:

- `snippets`
- `users` (with unique email constraint)

A seed user is also inserted by the setup script:

- **Email:** `alice@example.com`
- **Password:** `pa55word`

## TLS Certificates (Required)

The server is configured to run HTTPS and expects:

- `./tls/cert.pem`
- `./tls/key.pem`

Generate self-signed certs for local development:

```bash
mkdir -p tls
openssl req -x509 -newkey rsa:4096 -sha256 -days 365 -nodes \
  -keyout tls/key.pem -out tls/cert.pem \
  -subj "/CN=localhost"
```

## Running the App

From the project root:

```bash
go run ./cmd/web
```

Optional flags:

- `-addr` (default `:4000`) - HTTP(S) listen address
- `-debug` (default `false`) - include stack traces in 500 responses

Example:

```bash
go run ./cmd/web -addr :4000 -debug
```

Then open:

- <https://localhost:4000>

You may need to accept your local self-signed certificate warning in the browser.

## Development with Air (Optional)

This repo includes `.air.toml` for live reload.

If Air is installed:

```bash
air
```

## Routes

### Public

- `GET /` - home page (latest snippets)
- `GET /about` - about page
- `GET /snippet/view/{id}` - view a snippet
- `GET /user/signup` - signup form
- `POST /user/signup` - create account
- `GET /user/login` - login form
- `POST /user/login` - authenticate user
- `GET /ping` - health check (`OK`)
- `GET /static/*` - static assets (embedded FS)

### Authenticated

- `GET /snippet/create` - create snippet form
- `POST /snippet/create` - create snippet
- `GET /account/view` - account details
- `GET /account/settings` - account settings form
- `POST /account/settings` - update password
- `POST /user/logout` - logout

## Security Notes

- Session cookie is configured with:
  - `SameSite=Strict`
  - `Secure=true` (HTTPS only)
- CSRF protection is enabled for dynamic routes.
- Common security headers are set (CSP, frame/options, no-sniff, etc.).
- Auth-required pages send `Cache-Control: no-store`.

## Testing

Run all tests:

```bash
go test ./...
```

Some tests may require PostgreSQL availability depending on package and setup.

## Troubleshooting

- **`Error loading .env file`**: create `.env` in project root.
- **TLS cert errors on startup**: ensure `tls/cert.pem` and `tls/key.pem` exist.
- **DB connection failures**: verify `DATABASE_URL`, DB is running, and schema is loaded.
- **Session/auth issues**: confirm the app is served over HTTPS, as secure cookies are enforced.

## License

No license file is currently included in this repository.
