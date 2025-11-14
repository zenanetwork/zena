# Technology Stack

## Primary Language
- **Go 1.23.8** (main development language)

## Core Dependencies

### Blockchain Framework
- **Cosmos SDK**: v0.53.4 (blockchain application framework)
- **CometBFT**: v0.38.18 (BFT consensus engine, 10k+ TPS)
- **IBC**: v10.3.1 (Inter-Blockchain Communication protocol)

### Ethereum Compatibility
- **go-ethereum**: v1.15.11 (Ethereum implementation in Go)
- **Solidity**: Smart contract development language

### Key Libraries
- **btcsuite/btcd**: Bitcoin libraries for cryptography
- **holiman/uint256**: 256-bit integer operations
- **gorilla/mux**: HTTP routing
- **gorilla/websocket**: WebSocket support
- **spf13/cobra**: CLI framework
- **spf13/viper**: Configuration management

### Testing
- **Ginkgo v2**: BDD testing framework
- **Gomega**: Matcher library for assertions
- **pytest**: Python test framework (for scripts)

### Build Tools
- **Protocol Buffers**: Inter-service communication
- **Docker**: Containerization and build environments
- **Make**: Build automation

## Module Organization
Module path: `github.com/zenanetwork/zena`
