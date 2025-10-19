package network

import (
	testconstants "github.com/zenanetwork/zena/testutil/constants"
)

// chainsWEVMOSHex is an utility map used to retrieve the WZENA contract
// address in hex format from the chain ID.
//
// TODO: refactor to define this in the example chain initialization and pass as function argument
var chainsWEVMOSHex = map[testconstants.ChainID]string{
	testconstants.ExampleChainID: testconstants.WZENAContractMainnet,
}

// GetWZENAContractHex returns the hex format of address for the WZENA contract
// given the chainID. If the chainID is not found, it defaults to the mainnet
// address.
func GetWZENAContractHex(chainID testconstants.ChainID) string {
	address, found := chainsWEVMOSHex[chainID]

	// default to mainnet address
	if !found {
		address = chainsWEVMOSHex[testconstants.ExampleChainID]
	}

	return address
}
