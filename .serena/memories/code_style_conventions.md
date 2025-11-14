# Code Style & Conventions

## Go Code Style

### Formatting
- **Tool**: `gofumpt` (stricter version of gofmt)
- **Command**: `make format-go`
- Auto-formats all `.go` files (excluding vendor, generated files)

### Linting
- **Tool**: `golangci-lint` v2.2.2
- **Command**: `make lint-go` (check), `make lint-fix` (auto-fix)
- **Timeout**: 15 minutes for linting

### Enabled Linters
- `copyloopvar`: Loop variable copying
- `dogsled`: Blank identifier usage (max 3)
- `errcheck`: Unchecked errors
- `goconst`: Repeated constants
- `gocritic`: Code critique
- `gosec`: Security issues
- `govet`: Go vet analysis
- `ineffassign`: Ineffectual assignments
- `misspell`: Spelling mistakes
- `nakedret`: Naked returns
- `revive`: Fast, configurable, extensible linter
- `staticcheck`: Static analysis
- `thelper`: Test helper detection
- `unconvert`: Unnecessary conversions
- `unparam`: Unused parameters
- `unused`: Unused code

## Python Code Style

### Formatting
- **Tools**: `black` + `isort`
- **Commands**: 
  - `make format-python` (both tools)
  - `make format-black` (code formatting)
  - `make format-isort` (import sorting)

### Linting
- **Tools**: `pylint` + `flake8`
- **Command**: `make lint-python`

## Solidity Code Style

### Linting
- **Tool**: `solhint`
- **Commands**:
  - `make lint-contracts` (check)
  - `make lint-fix-contracts` (auto-fix)

## Shell Script Style

### Formatting
- **Tool**: `shfmt`
- **Command**: `make format-shell`

## Protobuf Style

### Formatting
- **Tool**: `clang-format`
- **Command**: `make proto-format`

### Linting
- **Tools**: `buf` + `protolint`
- **Command**: `make proto-lint`

## General Conventions
- Tests must pass before committing
- All commits must be **GPG-signed**
- Generated files (`.pb.go`, `.pb.gw.go`, `.pulsar.go`) excluded from manual editing
- Conventional commit format for PR titles
