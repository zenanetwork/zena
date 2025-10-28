package contracts

import (
	contractutils "github.com/zenanetwork/zena/contracts/utils"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"
)

func LoadERC20RecursiveReverting() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("solidity/ERC20RecursiveRevertingPrecompileCall.json")
}
