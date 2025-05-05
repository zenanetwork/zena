package zenad

import erc20types "github.com/zenanetwork/zena/x/erc20/types"

// WZENAContractMainnet is the WZENA contract address for mainnet
// TODO: Replace with the actual deployed WZENA contract address before mainnet launch
const WZENAContractMainnet = "0x0000000000000000000000000000000000000000"

// ExampleTokenPairs creates a slice of token pairs, that contains a pair for the native denom of the example chain
// implementation.
var ExampleTokenPairs = []erc20types.TokenPair{
	{
		Erc20Address:  WZENAContractMainnet,
		Denom:         ExampleChainDenom,
		Enabled:       true,
		ContractOwner: erc20types.OWNER_MODULE,
	},
}
