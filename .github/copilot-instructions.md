# AI Coding Agent Instructions for `saas_template_go`

## Project Overview

This project is a Go-based SaaS template designed with a modular architecture. It includes components for API handling, service logic, data models, and repository patterns. The structure is optimized for scalability and maintainability, with clear separation of concerns.

### Key Components

- **API Layer (`internal/api`)**: Handles HTTP requests and responses. Key files include `http.go`, `result.go`, and `user.go`.
- **Service Layer (`internal/svc`)**: Contains business logic. Example: `user.go` and its corresponding tests in `user_test.go`.
- **Repository Layer (`internal/repo`)**: Manages data persistence and retrieval. Example: `repo.go` and `user.go`.
- **Configuration (`internal/config`)**: Centralized configuration management using `config.go`.
- **Utilities (`internal/util`)**: Shared helper functions.
- **Migrations (`migrations`)**: SQL scripts for database schema management.

### Data Flow

1. **HTTP Request**: Received by the API layer.
2. **Service Logic**: API calls the service layer for business logic.
3. **Data Access**: Service layer interacts with the repository layer for database operations.
4. **Response**: Processed data is returned to the API layer and sent as an HTTP response.

## Developer Workflows

### Building the Project

Use the `Makefile` for common build tasks. Example:

```bash
make build
```

### Running Tests

- To run all tests:

```bash
go test ./...
```

- To run service-specific tests:

```bash
go test ./internal/svc
```

### Adding Dependencies

Use the `go` command to manage dependencies. Example:

```bash
go get github.com/spf13/viper
```

### Database Migrations

- Apply migrations using rust `sqlx`.
- Migration files are located in the `migrations` directory.

## Project-Specific Conventions

### Error Handling

- Use the `result.go` pattern for consistent API responses.
- Centralize error definitions in `internal/repo/error.go`.

### Logging

- Use the `logger` package in `internal/logger` for structured logging.

### Configuration

- Manage all configurations in `internal/config/config.go`.
- Use `viper` for environment-based configuration management.

### Testing

- Follow the structure in `internal/svc/user_test.go` for writing tests.
- Mock dependencies where possible to isolate unit tests.

## Examples

### Adding a New API Endpoint

1. Define the handler in `internal/api`.
2. Add the corresponding service logic in `internal/svc`.
3. Update the repository layer if data access is required.
4. Write tests in the appropriate `*_test.go` file.

### Writing a New Migration

1. Create `.up.sql` and `.down.sql` files in the `migrations` directory.
2. Apply the migration using `sqlx`.

---

This guide is a living document. Update it as the project evolves to ensure AI agents remain productive.
