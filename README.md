# newbox

Cross-platform machine setup CLI with an interactive TUI. Select profiles and categories of software to install on macOS, Windows, and Linux using native package managers.

## Quick Start

```bash
go run ./cmd/newbox
```

## Development

```bash
make build     # Build binary to bin/newbox
make test      # Run all tests
make test-race # Run tests with race detector
make lint      # Run go vet
make coverage  # Generate coverage report
```

## Architecture

newbox uses **Hexagonal Architecture** (Ports & Adapters):

- `internal/core/domain/` — Pure domain types (Platform, OS, Arch, etc.)
- `internal/core/port/` — Interface definitions (PlatformDetector, etc.)
- `internal/adapter/output/` — Implementations (system detector, package managers, etc.)
- `internal/adapter/input/` — TUI screens (coming soon)

## License

MIT