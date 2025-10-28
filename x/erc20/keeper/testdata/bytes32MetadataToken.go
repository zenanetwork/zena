package testdata

import (
	contractutils "github.com/zenanetwork/zena/contracts/utils"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"
)

// LoadBytes32MetadataTokenContract loads the Bytes32MetadataToken contract
// from the compiled JSON data.
func LoadBytes32MetadataTokenContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("Bytes32MetadataToken.json")
}
