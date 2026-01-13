#!/usr/bin/env pwsh
<#
Rename migration files in the migrations/ folder to follow the
TIMESTAMP_tbl-<table-name>.(up|down).sql pattern.

Behavior:
- For each .sql file named like <timestamp>_*.up.sql or .down.sql
  the script will read the file and try to extract the table name
  from a CREATE TABLE or DROP TABLE IF EXISTS statement.
- If a table name is found, the file will be renamed to
  <timestamp>_tbl-<table-name-with-hyphens>.up|down.sql
- If the target name already exists, the file is skipped.

Run: pwsh -NoProfile -ExecutionPolicy Bypass -File .\scripts\rename_migrations.ps1
#>

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$MigrationsDir = Resolve-Path -LiteralPath (Join-Path $ScriptDir '..\migrations')

Write-Output "Migrations dir: $MigrationsDir"

Get-ChildItem -Path $MigrationsDir -Filter '*.sql' | ForEach-Object {
    $file = $_
    $name = $file.Name

    if ($name -notmatch '^(\d{14})_(.+)\.(up|down)\.sql$') {
        Write-Output "[SKIP] filename doesn't match timestamp pattern: $name"
        return
    }

    $timestamp = $matches[1]
    $direction = $matches[3]

    $content = Get-Content -Raw -LiteralPath $file.FullName

    $table = $null
    if ($content -match '(?i)CREATE\s+TABLE\s+"?([A-Za-z0-9_]+)"?') { $table = $matches[1] }
    elseif ($content -match '(?i)DROP\s+TABLE\s+IF\s+EXISTS\s+"?([A-Za-z0-9_]+)"?') { $table = $matches[1] }

    if (-not $table) {
        Write-Output "[SKIP] no table name found in $name"
        return
    }

    $kebab = ($table -replace '_','-').ToLower()
    $newName = "${timestamp}_tbl-${kebab}.${direction}.sql"
    $newPath = Join-Path $file.DirectoryName $newName

    if ($name -eq $newName) {
        Write-Output "[OK] already named: $name"
        return
    }

    if (Test-Path -LiteralPath $newPath) {
        Write-Output "[CONFLICT] target exists, skipping: $newName"
        return
    }

    Write-Output "[RENAME] '$name' -> '$newName'"
    Rename-Item -LiteralPath $file.FullName -NewName $newName
}

Write-Output "Done."
