# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

### Building the Binary
```bash
# Build zenad binary to ./build/zenad
make build

# Install zenad to $GOPATH/bin (typically ~/go/bin)
make install

# Build for Linux AMD64
make build-linux

# Build with remote debugging options (disables optimization, preserves debug info)
make install COSMOS_BUILD_OPTIONS=nooptimization,nostrip
```

### Running Local Node
```bash
# Start local node (handles installation and setup)
./local_node.sh

# Available flags:
# -y: Overwrite previous database
# -n: Don't overwrite previous database
# --no-install: Skip binary installation
# --remote-debugging: Build with debugging options
```

The local node script creates a chain with:
- Chain ID: `zena_262144-1`
- Native denom: `atest`
- All precompiles enabled
- EVM permissionless mode
- JSON-RPC at http://localhost:8545

### Testing

```bash
# Run unit tests (excludes e2e and simulation)
make test-unit

# Run unit tests with coverage report
# Generates coverage.txt and prints coverage summary
make test-unit-cover

# Run zenad-specific tests
make test-zenad

# Run fuzz tests (precisebank keeper)
make test-fuzz

# Run Solidity integration tests
# Requires yarn, installs dependencies, runs against local node
make test-solidity

# Run system tests
make test-system

# Run Python script tests
make test-scripts

# Run all tests
make test-all
```

### Code Quality

```bash
# Run all linters (Go, Python, Solidity)
make lint

# Run individual linters
make lint-go        # golangci-lint
make lint-python    # pylint + flake8
make lint-contracts # solhint

# Auto-fix linting issues
make lint-fix
make lint-fix-contracts

# Format code
make format              # All languages
make format-go           # Go files (gofumpt)
make format-python       # Python (black + isort)
make format-shell        # Shell scripts (shfmt)
```

### Protobuf

```bash
# Generate all protobuf files
make proto-all

# Individual operations
make proto-gen     # Generate Go code
make proto-format  # Format .proto files
make proto-lint    # Lint .proto files
make proto-check-breaking  # Check for breaking changes
```

### Solidity Contracts

```bash
# Compile all contracts
make contracts-compile

# Clean compilation artifacts
make contracts-clean

# Add new contract to compilation list
make contracts-add CONTRACT=path/to/contract.sol
```

## Architecture Overview

### Module Structure

The repository uses a **dual-module structure**:
- **Root module** (`github.com/zenanetwork/zena`): Core EVM integration layers (ante, precompiles, x/ modules, rpc, etc.)
- **zenad module** (`./zenad`): Example chain implementation with separate go.mod that imports the root module

**Important**: When working in the zenad directory, remember it's a separate Go module:
```bash
# Working in root module
go test ./x/vm/...

# Working in zenad module (has its own go.mod)
cd zenad && go test ./tests/integration/...
```

To update zenad's dependency on the root module after making changes:
```bash
cd zenad
go mod edit -replace github.com/zenanetwork/zena=../
go mod tidy
```

#### Core Custom Modules (x/)

1. **x/vm**: EVM execution engine
   - Implements go-ethereum StateDB interface via Cosmos SDK KVStore
   - Manages EVM state transitions, gas metering, and transaction execution
   - Handles EVM configuration (chain config, opcodes, permissions)

2. **x/erc20**: Token representation unification
   - Bridges IBC Cosmos coins ↔ ERC20 tokens
   - Eliminates liquidity fragmentation between native and wrapped tokens
   - Dynamic precompile registration for each token pair

3. **x/feemarket**: EIP-1559 fee mechanism
   - Dynamic base fee adjustment based on block gas usage
   - Separate keeper for fee market parameters
   - Integrated with ante handler for fee validation

4. **x/precisebank**: Extended precision banking
   - Fractional coin balances beyond standard Cosmos SDK precision
   - Burn, mint, send operations with extended precision
   - Genesis state management for fractional balances

5. **x/ibc**: Custom IBC transfer integration
   - Overrides standard ICS20 keeper to support ERC20 token transfers
   - Callback system for IBC packet lifecycle events
   - Integrated with x/erc20 for automatic token conversion

### EVM Integration Architecture

#### StateDB Implementation (x/vm/statedb/)
The StateDB bridges go-ethereum's EVM with Cosmos SDK storage:
- **State Objects**: In-memory representation of account state
- **Journal**: Transaction-level state change tracking with revert capability
- **Snapshot System**: Nested snapshots for sub-transaction reverts
- **KVStore Mapping**: Maps Ethereum storage to Cosmos SDK KVStore

Key insight: All EVM state changes are journaled in-memory during transaction execution, then committed atomically to Cosmos SDK store only on successful completion.

#### Ante Handler Chain (ante/)
Dual-path transaction processing:

**Cosmos Transactions** (ante/cosmos_handler.go):
```
SetupContextDecorator
ValidateBasicDecorator
TxTimeoutHeightDecorator
→ EIP712 signature verification
→ Fee validation
→ Gas consumption
→ Sequence increment
```

**EVM Transactions** (ante/evm_handler.go):
```
SetupContextDecorator
MempoolFeeDecorator
→ EVM-specific validations
→ Account verification
→ Can transfer checks
→ Gas consumption (with EVM rules)
→ Nonce increment
→ Event emission
```

Critical: The ante handler selection happens based on transaction type (Cosmos Msg vs EVM MsgEthereumTx).

### Precompiles System

Precompiles expose Cosmos SDK functionality to Solidity contracts at deterministic addresses.

**Architecture** (precompiles/):
- **Stateless Design**: All precompiles are stateless, accessing state via keepers
- **Address Range**: 0x0000000000000000000000000000000000000800+
- **ABI Generation**: Auto-generated from Go interfaces
- **Authorization**: Built-in sender validation and approval mechanisms

**Available Precompiles**:
- bank: Native Cosmos coin operations
- bech32: Address format conversions
- distribution: Staking rewards distribution
- erc20: Token pair operations
- gov: Governance proposals and voting
- ics20: IBC transfers
- p256: NIST P-256 signature verification
- slashing: Validator slashing queries
- staking: Delegation and validator operations
- werc20: Wrapped ERC20 operations

**Registration Flow**:
1. Static precompiles registered at genesis (zenad/precompiles.go)
2. Dynamic ERC20 precompiles registered when token pairs created
3. Precompile availability checked during EVM call execution
4. State changes go through keeper interfaces (never direct storage access)

### JSON-RPC Server (rpc/)

Implements Ethereum JSON-RPC API for EVM compatibility:

**Backend Architecture** (rpc/backend/):
- **QueryClient Pattern**: Uses gRPC clients to query Cosmos SDK modules
- **Block Mapping**: Maps Cosmos block heights ↔ Ethereum block numbers
- **Transaction Indexing**: Custom indexer for EVM transactions
- **Event Filtering**: WebSocket support for eth_subscribe

**Available Namespaces** (rpc/namespaces/ethereum/):
- eth: Standard Ethereum RPC methods
- net: Network information
- web3: Web3 utilities
- debug: Debug and tracing (go-ethereum tracers)
- txpool: Transaction pool inspection
- miner: Mining control (adapted for Cosmos consensus)
- personal: Account management

### Mempool Integration (mempool/)

Custom mempool implementation for EVM transaction ordering:
- **Priority Ordering**: Supports EIP-1559 priority fee ordering
- **Signer Extraction**: Optimized EVM transaction signature recovery
- **Nonce Management**: Per-account nonce tracking
- **Legacy Pool**: Adapted from go-ethereum's legacy transaction pool

### Testing Architecture

**Integration Tests** (tests/integration/):
- Suite-based tests using testutil/integration/evm/network
- Unified test framework with automatic setup/teardown
- Factory pattern for transaction construction
- GRPC client pattern for queries

**Test Categories**:
- ante/: Ante handler decorator tests
- eip712/: EIP-712 signature tests
- mempool/: Mempool ordering tests
- precompiles/: Each precompile has dedicated test suite
- x/: Module-specific integration tests
- rpc/backend/: JSON-RPC backend tests

**Solidity Tests** (tests/solidity/):
- Hardhat-based test suites
- Runs against actual zenad node
- Tests precompile functionality from Solidity perspective

**System Tests** (tests/systemtests/):
- End-to-end chain behavior tests
- Uses compiled zenad binary
- Counter example for basic state transitions

## Key Technical Patterns

### Module Keeper Dependencies
The keeper dependency graph for zenad:
```
EVMKeeper → BankKeeper, StakingKeeper, FeeMarketKeeper, AccountKeeper
ERC20Keeper → EVMKeeper, BankKeeper, AccountKeeper
FeeMarketKeeper → (independent)
PreciseBankKeeper → BankKeeper
IBCTransferKeeper → ERC20Keeper, BankKeeper
```

Critical: Circular dependencies prevented through interface definitions in ante/interfaces/ and x/*/types/interfaces.go

### Transaction Flow

**EVM Transaction Path**:
1. JSON-RPC server receives eth_sendRawTransaction
2. Transaction converted to MsgEthereumTx
3. AnteHandler validates and prepares transaction
4. EVM keeper executes via StateDB
5. State changes committed to Cosmos KVStore
6. Events emitted for indexer
7. Receipt generated and returned

**State Transition Execution** (x/vm/keeper/state_transition.go):
- Constructs EVM with custom configuration
- Applies message (contract call or deployment)
- Refunds unused gas
- Updates block bloom filter
- Returns execution result

### Encoding System (encoding/)

Cosmos EVM uses custom encoding for Ethereum compatibility:
- **Amino**: Legacy Cosmos encoding (deprecated)
- **Protobuf**: Standard Cosmos SDK encoding
- **Ethereum RLP**: For EVM transactions
- **EIP-712**: For typed signature data

Key: MsgEthereumTx can be encoded as both Protobuf (for Cosmos) and RLP (for Ethereum tools).

### IBC Token Transfer Flow

When transferring tokens via IBC:
1. IBC transfer initiated (Cosmos or via precompile)
2. If token has ERC20 pair: burn ERC20 tokens
3. Transfer IBC packet created with native denom
4. On receive: check for ERC20 pair on destination
5. If pair exists: mint ERC20 tokens instead of bank coins
6. All handled transparently via x/erc20 middleware

## Development Notes

### Adding New Precompile
1. Create package in precompiles/
2. Implement Precompile interface (precompiles/common/precompile.go)
3. Define ABI and event structs
4. Register in zenad/precompiles.go
5. Add integration tests in tests/integration/precompiles/
6. Add Solidity tests if exposing to contracts

### Modifying EVM Behavior
- Chain config: zenad/cmd/zenad/config/config.go
- Custom opcodes: zenad/eips/
- EVM params: x/vm/types/params.go
- Gas calculation: ante/evm/fee_checker.go

### Working with StateDB
- State changes during EVM execution are in-memory
- Use Snapshot/RevertToSnapshot for sub-transaction isolation
- Commit() finalizes changes to Cosmos KVStore
- Never bypass StateDB to access storage directly

### Testing Recommendations
- Use testutil/integration for new integration tests
- Follow existing suite patterns in tests/integration/
- Mock keepers defined in x/*/types/mocks/ (generated via mockery)
- Coverage exclusions: cmd/, client/, proto/, testutil/, mocks/

## Common Workflows

### Running Single Integration Test
```bash
# From repository root
go test -v -run TestIntegrationTestSuite/TestSpecificTest ./tests/integration/x/vm/

# From zenad directory
cd zenad && go test -v -run TestSpecificTest ./tests/integration/...
```

### Debugging EVM Execution
1. Enable EVM tracing in config/app.toml: `tracer = "json"`
2. Use debug_traceTransaction RPC method
3. Check logs in ~/.zenad/logs/
4. Inspect StateDB state via printStateDB helper (in test code)

### Regenerating Mocks
```bash
make mocks
# Runs go generate ./... to regenerate mocks from interfaces
```

### Contributing Guidelines

All pull requests must follow these requirements:

1. **Link to GitHub Issue**: PRs without a corresponding issue will not be reviewed
2. **Signed Commits**: All commits must be signed (see [GitHub docs](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits))
3. **Substantial Changes**: Documentation-only PRs must make substantial or impactful changes (minor typo fixes will not be accepted)

When creating issues, include:
- Reproducibility steps (for bugs)
- Context and explanation
- Potential impact (severity, user scope, downstream effects)
