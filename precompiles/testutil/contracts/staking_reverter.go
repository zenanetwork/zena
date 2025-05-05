package contracts

import (
	contractutils "github.com/zenanetwork/zena/contracts/utils"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"
)

func LoadStakingReverterContract() (evmtypes.CompiledContract, error) {
	return contractutils.LoadContractFromJSONFile("StakingReverter.json")
}
