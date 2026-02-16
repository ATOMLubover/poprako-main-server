# Prepare development by installing necessary tools.
.PHONY: prepare
prepare:
	@echo Preparing development environment...

	# Check gofumpt installation for formatting.
	gofumpt -h
	
	# Check rust sqlx installation for migration management.
	sqlx --help

	@echo Development environment prepared: gofumpt, sqlx.

# Format the codebase using gofumpt.
.PHONY: fmt
fmt:
	@echo Formatting code...

	gofumpt -w .

	@echo Code formatted.
	
# Create a new up/down migration pair.
.PHONY: migadd
migadd:
	@echo Creating new migration files...

	sqlx migrate add -r $(name)

	@echo Migration files created.

# Apply all up migrations to the database.
.PHONY: migrun
migrun:
	@echo Applying up migrations...

	sqlx migrate run

	@echo Up migrations applied.
	
# Revert ONE last migration from the database.
.PHONY: migrev
migrev:
	@echo Reverting last migration...

	sqlx migrate revert

	@echo Last migration reverted.

# Revert ALL migrations from the database.
.PHONY: migrev-all
migrev-all:
	@echo Reverting all migrations...

	sqlx migrate revert --target-version 0

	@echo All migrations reverted.

.PHONY: vet
vet:
	@echo Vetting code...

	go vet ./...

	@echo Code vetted.

.PHONY: test-all
test-all:
	@echo Running tests...

	go test ./...

	@echo Tests completed.


.PHNOY: run
run:
	go run main.go

.PHONY: build
build:
	$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o poprako-main-server main.go