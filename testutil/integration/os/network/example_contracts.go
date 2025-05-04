package network

import (
	testconstants "github.com/cosmos/evm/testutil/constants"
)

// chainsWZENAHex is an utility map used to retrieve the WZENA contract
// address in hex format from the chain ID.
//
// TODO: refactor to define this in the example chain initialization and pass as function argument
var chainsWZENAHex = map[string]string{
	testconstants.ExampleChainID: testconstants.WZENAContractMainnet,
}

// GetWZENAContractHex returns the hex format of address for the WZENA contract
// given the chainID. If the chainID is not found, it defaults to the mainnet
// address.
func GetWZENAContractHex(chainID string) string {
	address, found := chainsWZENAHex[chainID]

	// default to mainnet address
	if !found {
		address = chainsWZENAHex[testconstants.ExampleChainID]
	}

	return address
}
