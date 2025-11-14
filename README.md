<img
src="rogo-header.png"
alt="Zenanet Blockchain images"
/>

**Please note**: This repository is undergoing changes while the code is being audited and tested. For the time being we will
be making v0.x releases. Some breaking changes might occur. Zenanet will be marked as stable with a v1
release after the audit, key stability features and benchmarking are completed.

**Visit the official Zenanet documentation**: [docs.zenanet.io](https://docs.zenanet.io/) (or [evm.cosmos.network](https://evm.cosmos.network/) for Cosmos EVM documentation)

## What is Zenanet?

Zenanet is a high-performance blockchain network built on Cosmos EVM technology, providing complete Ethereum compatibility within the Cosmos ecosystem. Zenanet combines the best of both worlds: Solidity smart contracts, Ethereum JSON-RPC, native EVM wallet/token/user experience, and seamless access to the Cosmos SDK modules through IBC (Inter-Blockchain Communication).

Powered by Cosmos EVM, Zenanet offers enterprise-grade customization for your business use case, chain architecture, and performance requirements.

## Integration & Architecture

Zenanet is built on Cosmos EVM technology, which can be integrated into existing chains
or added during new chain development by importing as a Go module library.

### Robust Defaults

Zenanet's modules come out of the box with defaults that enable rapid VM deployment. The integrated modules provide:

- Exposed JSON-RPC endpoints for connectivity with EVM tooling like wallets such as [MetaMask](https://metamask.io/) and [Rabby](https://rabby.io/), and block explorers like [Blockscout](https://docs.blockscout.com/).
- EVM extensions that allow functionality that is native to Cosmos SDK modules to be accessible from Solidity smart contracts [Solidity](https://docs.soliditylang.org/en/v0.8.26/) smart contracts.
- Use of any IBC asset in the EVM.

All modules can be controlled by on-chain governance.

### Extensive Customizability

Zenanet provides extensive customization capabilities built on Cosmos EVM's flexible architecture:

- **Permissioned EVM** - Implement customized access controls to either blacklist or whitelist individual addresses for calling and/or creating smart contracts on the network.
- **EVM Extensions** - Use custom EVM extensions to write custom business logic for your specific use case.
- **Single Token Representation v2 & ERC-20 Module** - The Single Token Representation v2 and our `x/erc20` module aligns IBC and ERC-20 token representation to simplify and improve user experience.
- **EIP-1559 Fee Market Mechanism** - Customize fee structures and transaction surge management with the self-regulating fee market mechanism based on [EIP-1559 fee market](https://eips.ethereum.org/EIPS/eip-1559).
- **JSON-RPC Server** - Full control over the exposed namespaces and [JSON-RPC server](https://cosmos-docs.mintlify.app/docs/api-reference/ethereum-json-rpc). Configurable parameters include custom timeouts for EVM calls or HTTP requests, maximum block gas, open connections, and more.
- **EIP-712 Signing** - Integrated [EIP-712 signature](https://eips.ethereum.org/EIPS/eip-712) implementation allows Cosmos SDK messages to be signed with EVM wallets like MetaMask. This supports structured data signing for arbitrary messages.
- **Custom Improvement Proposals (Opcodes)** - Zenanet provides the opportunity to customize EVM opcodes and add new ones. Read more on [custom operations here](https://cosmos-docs.mintlify.app/docs/documentation/smart-contracts/custom-improvement-proposals#custom-improvement-proposals).

## Compatibility with Ethereum

Is Zenanet "Ethereum equivalent"? Ethereum-equivalence describes any EVM solution that is identical in transaction execution to the Ethereum client. On the other hand, Ethereum-compatible means that the EVM implementation can run every transaction that is valid on Ethereum, while also handling divergent transactions that are not valid on Ethereum.

We describe Zenanet as **forward-compatible** with Ethereum. It can run any valid smart contract from Ethereum and also implement new features that are not yet available in the standard Ethereum VM, thus moving the standard forward.

## Getting Started

To run the Zenanet node (`zenad`), execute the following script from the root folder of the repository:

```bash
./local_node.sh
```

### Migrations

We provide upgrade guides [here](./docs/migrations) for upgrading from various versions of Zenanet and Cosmos EVM.

### Testing

All test scripts are found in `Makefile` in the root of the repository.
Listed below are the commands for various tests:

#### Unit Testing

```bash
make test-unit
```

#### Coverage Test

This generates a code coverage file `filtered_coverage.txt` and prints out the
covered code percentage for the working files.

```bash
make test-unit-cover
```

#### Fuzz Testing

```bash
make test-fuzz
```

#### Solidity Tests

```bash
make test-solidity
```

#### Benchmark Tests

```bash
make benchmark
```

## Open-source License & Credits

Zenanet is fully open-source under the Apache 2.0 license. It is built on [Cosmos EVM](https://evm.cosmos.network/), which is a fork of [evmOS](https://github.com/evmos/OS). The Interchain Foundation funded [evmOS developers](https://github.com/evmos/OS) Tharsis to open-source the original evmOS codebase. We acknowledge Tharsis and evmOS for performing the foundational work for EVM compatibility and interoperability in the Cosmos ecosystem.

## Developer Community and Support

The issue list of this repository is exclusively for bug reports and feature requests. For questions and discussions, please join our community channels.

**| Need Help? | Zenanet Community: [Discord](https://discord.gg/zenanet) - [Telegram](https://t.me/zenanet) - [Documentation](https://docs.zenanet.io) |**

For Cosmos ecosystem support: [Discord](https://discord.com/invite/interchain) - [Telegram](https://t.me/CosmosOG) - [#Cosmos-tech Slack](https://forms.gle/A8jawLgB8zuL1FN36)

## Maintainers

Zenanet is maintained by the Zenanet core development team. The project is built on top of the Cosmos Stack, which is maintained by [Cosmos Labs](https://cosmoslabs.io/) and includes core components such as Cosmos SDK, CometBFT, IBC, and Cosmos EVM.

The Cosmos Stack is maintained by Cosmos Labs, a wholly-owned subsidiary of the [Interchain Foundation](https://interchain.io/), and is supported by a robust community of open-source contributors.

## Contributing to Zenanet

We welcome open source contributions and discussions! For more on contributing, read the [guide](./CONTRIBUTING.md).

### Acknowledgments

We would like to thank:

- **Cosmos Labs** for maintaining the Cosmos EVM and Cosmos Stack
- **B-Harvest** and **Mantra** for their key contributions to Cosmos EVM development
- The **Interchain Foundation** for funding the original evmOS development by Tharsis
- The entire **Cosmos community** for their continued support and contributions

## Documentation and Resources

### Zenanet Documentation

- **Official Zenanet Documentation**: [docs.zenanet.io](https://docs.zenanet.io/) (Coming Soon)
- **Cosmos EVM Documentation**: [evm.cosmos.network](https://evm.cosmos.network/)

### Technology Stack

Zenanet is built on the following core technologies:

- **[Cosmos SDK](http://github.com/cosmos/cosmos-sdk)** - A framework for building blockchain applications in Golang
- **[IBC (Inter-Blockchain Communication Protocol)](https://github.com/cosmos/ibc-go/)** - A blockchain interoperability protocol that allows blockchains to transfer any type of data encoded in bytes
- **[CometBFT](https://github.com/cometbft/cometbft)** - High-performance, 10k+ TPS configurable BFT consensus engine
- **[Cosmos EVM](https://evm.cosmos.network/)** - EVM compatibility layer for Cosmos chains
