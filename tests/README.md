# Integration Tests

## Overview

This directory contains app-level integration tests that simulate a complete comic translation workflow, from workset creation through comic publication and cleanup.

## Prerequisites

1. **Build the application**:

   ```bash
   cd ..
   make build
   # or
   go build -o bin/project .
   ```

2. **Configure environment**:

   - Ensure `.env` file exists in the project root
   - Verify `MOCK_AUTH_TOKEN` is set (should already exist)
   - Database should be accessible via `DATABASE_URL`

3. **Database state**:
   - Tests will create and delete test data
   - No manual cleanup needed (tests clean up after themselves)
   - Can be run multiple times

## Running Tests

### Run all integration tests:

```bash
cd tests
go test -v
```

### Run with detailed output:

```bash
go test -v -count=1
```

The `-count=1` flag disables test caching, ensuring fresh execution each time.

## Test Flow

The integration test simulates a complete translation team workflow:

### 1. **Create Workset and Comic** (`testCreateWorksetAndComic`)

- Creates a workset via `POST /api/v1/worksets`
- Verifies workset creation via `GET /api/v1/worksets/{id}`
- Creates a comic in the workset via `POST /api/v1/comics`
- Retrieves comic details via `GET /api/v1/comics/{id}`

### 2. **Create Pages and Units** (`testCreatePagesAndUnits`)

- Batch creates pages via `POST /api/v1/pages`
- Marks pages as uploaded via `PATCH /api/v1/pages/{id}`
- Creates translation units via `POST /api/v1/pages/{page_id}/units`
- Retrieves units to get their IDs

### 3. **Translation Workflow** (`testTranslationWorkflow`)

- Updates units with translated text via `PATCH /api/v1/pages/{page_id}/units`
- Simulates proofreading by updating with proved text
- Verifies units were updated correctly

### 4. **Assignments and Status** (`testAssignmentsAndStatus`)

- Creates assignments via `POST /api/v1/assignments`
- Retrieves assignment by ID via `GET /api/v1/assignments/{id}`
- Lists comic assignments via `GET /api/v1/comics/{comic_id}/assignments`
- Lists user assignments via `GET /api/v1/users/{user_id}/assignments`

### 5. **Retrieval and Filtering** (`testRetrievalAndFiltering`)

- Tests pagination with `GET /api/v1/worksets?limit=10&offset=0`
- Retrieves workset comics via `GET /api/v1/worksets/{id}/comics`
- Filters comics by author using query params
- Lists comic pages via `GET /api/v1/comics/{id}/pages`

### 6. **Cleanup** (`testCleanup`)

- Deletes units via `DELETE /api/v1/pages/{page_id}/units`
- Deletes assignment via `DELETE /api/v1/assignments/{id}`
- Deletes pages via `DELETE /api/v1/pages/{id}`
- Deletes comic via `DELETE /api/v1/comics/{id}`
- Deletes workset via `DELETE /api/v1/worksets/{id}`

## Test Execution Order

Tests run **serially** in a single `TestIntegrationFlow` to maintain strict ordering:

```
TestMain (setup)
  ↓
TestIntegrationFlow
  ├── 01_CreateWorksetAndComic
  ├── 02_CreatePagesAndUnits
  ├── 03_TranslationWorkflow
  ├── 04_AssignmentsAndStatus
  ├── 05_RetrievalAndFiltering
  └── 06_Cleanup
  ↓
TestMain (teardown)
```

Each subtest depends on data created by previous subtests, ensuring a realistic workflow simulation.

## Server Management

- **Automatic startup**: Server starts automatically via `TestMain`
- **Health check**: Tests wait for `/api/v1/check-update` to respond before proceeding
- **Automatic shutdown**: Server stops after all tests complete
- **Port**: Server runs on `localhost:8080` (configured in `app_config.json`)

## Authentication

All requests use the `MOCK_AUTH_TOKEN` from `.env`, which represents a user with:

- All roles (translator, proofreader, typesetter, redrawer, reviewer, uploader)
- Admin privileges
- User ID: `019bbf6f-fa6b-7119-b1cd-c961a808c864`

## Troubleshooting

### Server fails to start

- Check if port 8080 is available
- Verify `bin/project` binary exists
- Ensure database connection is configured

### Auth token errors

- Verify `MOCK_AUTH_TOKEN` exists in `.env`
- Check token expiration date (should be valid until 2026)

### Database errors

- Ensure migrations have been applied
- Check `DATABASE_URL` in `.env`
- Verify database is accessible

### Tests fail partway through

- Check server logs for errors
- Database constraints may have been violated
- Run tests again (cleanup should handle partial state)

## Manual Database Cleanup

If tests fail before cleanup completes, you can manually clean the database:

```sql
-- Find test data by name pattern
SELECT * FROM workset_tbl WHERE name LIKE 'Test Workset%';
SELECT * FROM comic_tbl WHERE author = '测试作者';

-- Delete manually if needed
DELETE FROM workset_tbl WHERE name LIKE 'Test Workset%';
```

Or drop and recreate the database schema using migrations.

## Notes

- Tests modify real database data
- No data isolation between test runs
- Only one test flow can run at a time
- Server output goes to stdout/stderr for debugging
