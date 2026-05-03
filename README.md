# secure-notes-api

A production-ready REST API for managing notes, built with Go and PostgreSQL.

![CI](https://github.com/hvnrstrd/secure_notes_api/actions/workflows/ci.yml/badge.svg)

## Stack

- **Go 1.25** — API server
- **PostgreSQL 16** — persistent storage
- **Docker** — containerization (multi-stage build, non-root user)
- **GitHub Actions** — CI/CD pipeline

## Security

Every push is automatically scanned by:

- **gitleaks** — detects secrets and credentials in code
- **gosec** — static analysis for Go security issues
- **trivy** — scans Docker image for CVEs

## API

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

API will be available at `http://localhost:8080`.

## Examples

```bash
# Create a note
curl -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"hello","body":"world"}'

# Get all notes
curl http://localhost:8080/notes

# Delete a note
curl -X DELETE http://localhost:8080/notes/{id}
```

## Architecture decisions

- **Non-root Docker user** — container runs as `appuser`, not root
- **Multi-stage build** — final image contains only the binary, not Go toolchain (~10MB vs ~300MB)
- **Server timeouts** — prevents slowloris attacks (`ReadTimeout`, `WriteTimeout`, `IdleTimeout`)
- **Storage interface** — handler doesn't know about PostgreSQL directly, easy to swap implementations
- **Health check** — API waits for PostgreSQL to be healthy before starting