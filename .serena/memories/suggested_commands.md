# Suggested Commands for Zena Development

## Build Commands

### Build Binary
```bash
make build           # Build zenad to ./build/zenad
make build-linux     # Cross-compile for Linux AMD64
make install         # Install zenad to $GOPATH/bin
```

### Dependencies
```bash
make go.sum         # Verify and tidy dependencies
make vulncheck      # Check for security vulnerabilities
```

## Testing Commands

### Unit Tests
```bash
make test-unit              # Run all unit tests (timeout: 15m)
make test-zenad             # Run zenad module tests
make test-unit-cover        # Run tests with coverage report
make test-race              # Run tests with race detector
make test-all               # Run all tests (evm + zenad modules)
```

### Specialized Tests
```bash
make test-fuzz              # Run fuzz tests
make test-solidity          # Run Solidity contract tests
make test-scripts           # Run Python script tests (pytest)
make benchmark              # Run benchmark tests
```

### Coverage
```bash
make test-unit-cover        # Generates coverage.txt with filtered results
# Coverage excludes: /cmd/, /client/, /proto/, /testutil/, /mocks/, test files, generated files
```

## Linting Commands

### Check Code Quality
```bash
make lint                   # Lint all (Go, Python, Solidity)
make lint-go               # Lint Go code only
make lint-python           # Lint Python scripts (pylint + flake8)
make lint-contracts        # Lint Solidity contracts (solhint)
```

### Auto-fix Issues
```bash
make lint-fix              # Auto-fix Go linting issues
make lint-fix-contracts    # Auto-fix Solidity linting issues
```

## Formatting Commands

### Format Code
```bash
make format                # Format all code (Go, Python, Shell)
make format-go            # Format Go code (gofumpt)
make format-python        # Format Python (black + isort)
make format-black         # Format Python with black only
make format-isort         # Sort Python imports with isort
make format-shell         # Format shell scripts (shfmt)
```

## Protobuf Commands

### Generate & Update Protobuf
```bash
make proto-all            # Format, lint, and generate proto files
make proto-gen            # Generate Go code from proto files
make proto-format         # Format proto files (clang-format)
make proto-lint           # Lint proto files (buf + protolint)
make proto-check-breaking # Check for breaking changes against main
```

## Running the Node

### Local Development Node
```bash
./local_node.sh           # Run local development node (zenad)
```

## Workflow Best Practices

### Before Committing
```bash
make format               # Format all code
make lint                 # Check code quality
make test-unit            # Run unit tests
# Ensure all commands pass before committing
```

### Full Quality Check
```bash
make format && make lint && make test-unit-cover
# Run formatting, linting, and coverage tests
```

## macOS-Specific Notes
- System: Darwin (macOS)
- Uses old Apple linker for fuzz tests (workaround for xcode issue)
- CGO_ENABLED=1 required for builds
