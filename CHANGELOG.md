# CHANGELOG

## UNRELEASED

### NOTICE

- This software is a modified version of [cosmos/evm](https://github.com/cosmos/evm) project
- All modifications are made by [Your Organization Name] for the Zena blockchain
- The original software is licensed under Apache License 2.0
- This modified version maintains the same license terms

### BRANDING CHANGES

- Rebranded entire codebase from cosmos/evm to Zena blockchain
- Changed chain ID conventions to zena_1 (mainnet), zenatest_5/zenatest_11 (testnets)
- Modified token denomination from "atest" to "azena"
- Updated display denomination to "ZENA"
- Renamed binary from "evmd" to "zenad"
- Updated all references to reflect Zena brand identity

### DEPENDENCIES

- [\#31](https://github.com/cosmos/evm/pull/31) Migrated example_chain to evmd
- Migrated evmos/go-ethereum to cosmos/go-ethereum
- Migrated evmos/cosmos-sdk to cosmos/cosmos-sdk
- [\#95](https://github.com/cosmos/evm/pull/95) Bump up ibc-go from v8 to v10

### BUG FIXES

- Fixed example chain's cmd by adding NoOpEVMOptions to tmpApp in root.go
- Added RPC support for `--legacy` transactions (Non EIP-1559)
- Fixed "unknown chain id: zena_1" error by adding CosmosChainID to ChainsCoinInfo map in test environment
- [\#140](https://github.com/cosmos/evm/pull/140) Fixed inconsistent method name in distribution precompile (withdrawDelegatorReward)
- [\#143](https://github.com/cosmos/evm/pull/143) Fixed Equivocation.Time type in evidence precompile (uint64 -> int64) to prevent potential overflow

### IMPROVEMENTS

- [\#138](https://github.com/cosmos/evm/pull/138) Separated Cosmos and EVM Chain IDs for improved compatibility

### FEATURES

- [\#54](https://github.com/cosmos/evm/pull/54) Added EVM post transaction hooks with safety checks for gas usage
- [\#84](https://github.com/cosmos/evm/pull/84) permissionless erc20 registration to cosmos coin conversion
- [\#69](https://github.com/cosmos/evm/pull/69) Add new `x/precisebank` module with bank decimal extension for EVM usage.

### STATE BREAKING

- Refactored evmos/os into cosmos/evm
- Renamed x/evm to x/vm
- Renamed protobuf files from evmos to cosmos org
- [\#83](https://github.com/cosmos/evm/pull/83) Remove base fee v1 from x/feemarket
- [\#93](https://github.com/cosmos/evm/pull/93) Remove legacy subspaces
- [\#95](https://github.com/cosmos/evm/pull/95) Replaced erc20/ with erc20 in native ERC20 denoms prefix for IBC v2
- [\#62](https://github.com/cosmos/evm/pull/62) Remove x/authz dependency from precompiles
- [\#138](https://github.com/cosmos/evm/pull/138) Separated Cosmos and EVM Chain IDs, removed chain config from keeper params

### API-Breaking

- Refactored evmos/os into cosmos/evm
- Renamed x/evm to x/vm
- Renamed protobuf files from evmos to cosmos org
- [\#95](https://github.com/cosmos/evm/pull/95) Updated ics20 precompile to use Denom instead of DenomTrace for IBC v2
- [\#138](https://github.com/cosmos/evm/pull/138) Separated Cosmos and EVM Chain IDs, removed EVM chain config from params
