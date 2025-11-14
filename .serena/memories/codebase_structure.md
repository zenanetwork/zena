# Codebase Structure

## Root Directory Layout

```
zena/
├── zenad/              # Main binary application (separate Go module)
├── x/                  # Cosmos SDK modules (custom chain logic)
│   ├── erc20/         # ERC-20 token module
│   ├── precisebank/   # Precise banking operations
│   ├── ibc/           # IBC integration
│   ├── vm/            # EVM virtual machine
│   └── feemarket/     # Fee market mechanism (EIP-1559)
├── contracts/          # Solidity smart contracts
├── proto/              # Protocol buffer definitions
├── tests/              # Integration and E2E tests
├── testutil/           # Testing utilities
├── scripts/            # Build and utility scripts
├── docs/               # Documentation
├── api/                # Generated API code
├── client/             # Client libraries
├── rpc/                # JSON-RPC server implementation
├── server/             # Server components
├── indexer/            # Blockchain indexer
├── ethereum/           # Ethereum compatibility layer
├── encoding/           # Data encoding/decoding
├── ante/               # Ante handlers (transaction preprocessing)
├── crypto/             # Cryptographic utilities
├── utils/              # General utilities
├── wallets/            # Wallet implementations
├── mempool/            # Transaction mempool
├── version/            # Version information
├── metrics/            # Monitoring and metrics
├── config/             # Configuration files
├── precompiles/        # EVM precompiled contracts
├── eips/               # Ethereum Improvement Proposal implementations
├── ibc/                # IBC (Inter-Blockchain Communication) integration
├── contrib/            # Community contributions
└── .github/            # GitHub workflows and configs

## Key Files

### Build & Configuration
- `Makefile`              # Build automation and commands
- `go.mod`, `go.sum`      # Go module dependencies
- `.golangci.yml`         # Go linting configuration
- `buf.gen.proto.yaml`    # Protobuf generation config
- `buf.work.yaml`         # Buf workspace config
- `docker-compose.yml`    # Docker services configuration
- `local_node.sh`         # Local development node script

### Code Quality
- `.pylintrc`             # Python linting config
- `.protolint.yml`        # Protobuf linting config
- `.solhint.json`         # Solidity linting config
- `.markdownlint.yml`     # Markdown linting config
- `.yamllint`             # YAML linting config
- `.gitleaks.toml`        # Secret scanning config
- `codecov.yml`           # Code coverage config

### Documentation
- `README.md`             # Project overview
- `CONTRIBUTING.md`       # Contribution guidelines
- `CHANGELOG.md`          # Version history
- `SECURITY.md`           # Security policy
- `LICENSE`               # Apache 2.0 license

## Module Structure

### Cosmos SDK Modules (`x/`)
Each module follows standard Cosmos SDK structure:
- `keeper/` - State management logic
- `types/` - Type definitions and proto messages
- `module.go` - Module interface implementation
- Tests alongside implementation files

### zenad Application
Separate Go module under `./zenad/` containing:
- `cmd/zenad/` - Main entry point
- Application-specific configurations

## Generated Code
- `.pb.go` - Protobuf-generated Go code
- `.pb.gw.go` - gRPC-gateway generated code
- `.pulsar.go` - Pulsar generated code
- These files are excluded from coverage and manual editing

## Testing Organization
- `tests/` - Integration and E2E tests
- `testutil/` - Shared testing utilities
- Unit tests colocated with source files
- Simulation tests in dedicated directories
