# Zena Blockchain - Project Overview

## Purpose
Zena is a **Cosmos EVM implementation** (fork of Cosmos Labs' evmOS project) that adds EVM compatibility and customizability to Cosmos SDK chains. It's a blockchain solution for "Zenanet" that provides:

- **Complete Ethereum Capabilities**: Solidity smart contracts, Ethereum JSON-RPC, EVM wallet/token support
- **Cosmos Integration**: Native IBC support, access to Cosmos SDK modules from EVM
- **Forward Compatibility**: Can run any valid Ethereum smart contract plus new features not yet in standard EVM
- **Customizability**: Permissioned EVM, custom extensions, EIP-1559 fee market, EIP-712 signing

## Key Features
- **Plug-and-play** EVM solution for Cosmos chains
- **IBC Integration**: Use any IBC asset in the EVM
- **ERC-20 Module**: Single token representation across IBC and ERC-20
- **JSON-RPC Server**: Full Ethereum tooling compatibility (MetaMask, Rabby, Blockscout)
- **Custom Opcodes**: Ability to customize EVM opcodes and add new ones
- **On-chain Governance**: All modules controllable by governance

## Current Status
- **Pre-v1.0** (v0.x releases): Code undergoing audits and testing
- Breaking changes may occur before v1.0 stable release
- Based on evmOS open-source codebase (Apache 2.0 license)

## Documentation
- Official Docs: https://evm.cosmos.network/
- GitHub: https://github.com/zenanetwork/zena
