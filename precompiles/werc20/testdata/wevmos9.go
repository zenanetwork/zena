package testdata

import (
	contractutils "github.com/zenanetwork/zena/contracts/utils"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"
)

// LoadWEVMOS9Contract load the WEVMOS9 contract from the json representation of
// the Solidity contract.
func LoadWEVMOS9Contract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("WEVMOS9.json")
}
