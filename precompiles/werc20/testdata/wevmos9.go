package testdata

import (
	contractutils "github.com/cosmos/evm/contracts/utils"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// LoadWZENA9Contract load the WZENA9 contract from the json representation of
// the Solidity contract.
func LoadWZENA9Contract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("WZENA9.json")
}
