# KarisDB

KarisDB is a lightweight, cross-platform local development environment manager. Version 0.1 provides a Go core, a Python CLI frontend, and an alternative Go CLI. Both frontends dispatch to the same Go application router. The core is a one-request process; KarisDB does not run a permanent daemon.

## Requirements

- Go 1.23 or newer
- Python 3.11 or newer for the Python frontend

## Build on Windows PowerShell

```powershell
New-Item -ItemType Directory -Force dist | Out-Null
go build -o dist/karisdb-core.exe ./cmd/karisdb-core
go build -o dist/karisdb-go.exe ./cmd/karisdb-go
```

Run the Go CLI:

```powershell
./dist/karisdb-go.exe help
./dist/karisdb-go.exe version
./dist/karisdb-go.exe status
```

Run the Python CLI from the repository root:

```powershell
$env:PYTHONPATH = "python"
python -m karisdb help
python -m karisdb version
python -m karisdb status
```

The Python frontend locates `dist/karisdb-core.exe`, a sibling executable, or `karisdb-core.exe` on `PATH`. Set `KARISDB_CORE_PATH` to an explicit executable path when needed.

## Tests and checks

```powershell
go test ./...
go vet ./...
$env:PYTHONPATH = "python"
python -m unittest discover -s python/tests -v
```

After building the Go binaries, run the end-to-end CLI checks:

```powershell
$env:PYTHONPATH = "python"
$env:KARISDB_CORE_PATH = (Resolve-Path .\dist\karisdb-core.exe).Path
$env:KARISDB_GO_CLI_PATH = (Resolve-Path .\dist\karisdb-go.exe).Path
python -m unittest discover -s python/tests -v
```

## v0.1 behavior

- `help`, `version`, and `status` are implemented in Go.
- Core requests use one JSON Lines request and one JSON response per process invocation. Protocol input is limited to 1 MiB.
- Configuration is an optional JSON file at `<data_dir>/configs/config.json`. It currently supports only an absolute `data_dir` override.
- Data paths follow `%LOCALAPPDATA%\KarisDB` on Windows, `XDG_DATA_HOME/karisdb` on Linux, and `~/Library/Application Support/KarisDB` on macOS.
- `status` is read-only and does not create data directories or start services.
- Runtime installation, databases, projects, and process management are not part of v0.1.
