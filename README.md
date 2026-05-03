# secure-notes-api

A production-ready notes app with JWT authentication, built with Go and PostgreSQL.

![CI](https://github.com/hvnrstrd/secure_notes_api/actions/workflows/ci.yml/badge.svg)

## Stack

- **Go 1.25** — API server
- **PostgreSQL 16** — persistent storage
- **Docker** — containerization (multi-stage build, non-root user)
- **GitHub Actions** — CI/CD pipeline
- **Vanilla JS** — frontend (no frameworks)

## Security

Every push is automatically scanned by:

- **gitleaks** — detects secrets and credentials in code
- **gosec** — static analysis for Go security issues
- **trivy** — scans Docker image for CVEs

## API

### Auth

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /auth/register | Create account |
| POST | /auth/login | Get JWT token |

### Notes (requires Bearer token)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /notes | Get all notes |
| POST | /notes | Create a note |
| GET | /notes/{id} | Get note by ID |
| DELETE | /notes/{id} | Delete a note |

## Quick start

```bash
git clone https://github.com/hvnrstrd/secure_notes_api.git
cd secure_notes_api
docker compose up --build
```

API: `http://localhost:8080`

Frontend:

```bash
python3 -m http.server 3000 --directory frontend
```

Open `http://localhost:3000`

## Examples

```bash
# Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'

# Create note (use token from login)
curl -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"title":"hello","body":"world"}'
```

## Architecture decisions

- **Non-root Docker user** — container runs as `appuser`, not root
- **Multi-stage build** — final image contains only the binary (~10MB vs ~300MB)
- **Server timeouts** — prevents slowloris attacks
- **JWT authentication** — stateless, no sessions stored on server
- **bcrypt passwords** — passwords are hashed, never stored in plaintext
- **Storage interface** — handler is decoupled from PostgreSQL, easy to swap
- **Health check** — API waits for PostgreSQL to be healthy before starting
- **CORS middleware** — controlled cross-origin access