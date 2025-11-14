# Architecture & Design Patterns

## Project Architecture

### Cosmos EVM Integration Model
Zena follows the **Cosmos EVM plug-and-play architecture**:

```
Cosmos SDK Chain (Base Layer)
    ↓
Cosmos EVM Integration (Compatibility Layer)
    ↓
EVM Runtime + Ethereum JSON-RPC
    ↓
Solidity Smart Contracts + Ethereum Tooling
```

### Key Architectural Decisions

#### 1. Forward-Compatible Design
- **Not "Ethereum Equivalent"**: Doesn't replicate Ethereum transaction execution exactly
- **Ethereum Compatible**: Runs all valid Ethereum transactions
- **Forward-Compatible**: Implements features not yet in standard Ethereum VM
- Can handle both standard Ethereum and divergent transactions

#### 2. Module-Based Architecture
Follows **Cosmos SDK module pattern**:
- Each feature is a separate module in `/x/`
- Modules are composable and independently governable
- Standard keeper pattern for state management
- Clear separation of concerns

#### 3. Dual Module System
- **Root modules**: Core EVM functionality
- **zenad modules**: Application-specific logic
- Separate Go modules for clean separation

## Design Patterns

### Cosmos SDK Patterns

#### Keeper Pattern
```go
// State management through keepers
type Keeper struct {
    storeKey sdk.StoreKey
    cdc      codec.Codec
    // ... other dependencies
}
```
- Centralized state management
- Dependency injection
- Clear access control

#### Module Pattern
Each module implements:
- `module.go` - Module interface
- `keeper/` - State management
- `types/` - Type definitions
- `genesis.go` - Genesis state handling

### EVM Integration Patterns

#### Precompiles for Native Access
- EVM precompiles expose Cosmos SDK functionality to Solidity
- Allows smart contracts to interact with native modules (IBC, governance, etc.)
- Bridge between EVM world and Cosmos world

#### Single Token Representation v2
- Aligns IBC and ERC-20 token representations
- Unified user experience across chains
- Simplifies cross-chain token handling

### Transaction Flow

```
User Transaction
    ↓
JSON-RPC Server (Ethereum compatibility)
    ↓
Ante Handlers (preprocessing, fee validation)
    ↓
EVM Execution (Solidity contracts)
    ↓
State Changes (Cosmos SDK modules)
    ↓
Consensus (CometBFT)
```

## Customization Points

### 1. Permissioned EVM
- Whitelist/blacklist addresses
- Control contract deployment
- Control contract execution

### 2. Custom EVM Extensions
- Add custom business logic
- Extend EVM capabilities
- Use case-specific functionality

### 3. Fee Market Customization
- EIP-1559 base fee mechanism
- Custom surge pricing
- Transaction priority management

### 4. JSON-RPC Configuration
- Namespace exposure control
- Timeout customization
- Connection limits
- Block gas limits

### 5. Custom Opcodes
- Modify existing opcodes
- Add new opcodes
- Use case-specific operations

## Testing Architecture

### Test Layers
1. **Unit Tests**: Individual function/method testing
2. **Integration Tests**: Module interaction testing
3. **E2E Tests**: Full chain testing (`tests/e2e/`)
4. **Solidity Tests**: Smart contract testing
5. **Fuzz Tests**: Random input testing
6. **Simulation**: Stochastic simulation testing

### Coverage Strategy
- Minimum coverage targets (not explicitly stated)
- Excludes generated code, mocks, CLI, test utilities
- Focus on keeper logic and core functionality

## Development Patterns

### Protobuf-First Design
- API definitions in `.proto` files
- Code generation for Go, gRPC, gRPC-gateway
- Cross-language compatibility
- Clear API contracts

### Configuration-Driven
- Extensive use of config files
- Environment-based configuration
- Runtime configurability via governance

### Governance Integration
- All modules controllable by on-chain governance
- Parameter changes via proposals
- Community-driven evolution

## Performance Considerations

### CometBFT Consensus
- 10k+ TPS capability
- Configurable performance parameters
- BFT consensus guarantees

### EVM Optimization
- Custom opcode implementations
- Gas metering optimization
- State access patterns

## Security Patterns

### Defense in Depth
- Multiple validation layers
- Ante handlers for preprocessing
- EVM-level validation
- Cosmos SDK validation

### Access Control
- Module-level permissions
- Governance-controlled parameters
- Optional permissioned EVM

## Interoperability Design

### IBC Integration
- Native IBC support in EVM
- Cross-chain asset transfers
- Inter-chain contract calls

### Ethereum Compatibility
- JSON-RPC endpoints
- EIP-712 signature support
- MetaMask/wallet compatibility
- Block explorer compatibility
