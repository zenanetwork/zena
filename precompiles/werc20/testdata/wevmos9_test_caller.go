package testdata

import (
	contractutils "github.com/zenanetwork/zena/contracts/utils"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"
)

func LoadWZENA9TestCaller() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("WZENA9TestCaller.json")
}
