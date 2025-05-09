package ante_test

import (
	"github.com/zenanetwork/zena/testutil/integration/os/network"
	evmante "github.com/zenanetwork/zena/x/vm/ante"

	storetypes "cosmossdk.io/store/types"
)

func (suite *EvmAnteTestSuite) TestBuildEvmExecutionCtx() {
	network := network.New()

	ctx := evmante.BuildEvmExecutionCtx(network.GetContext())

	suite.Equal(storetypes.GasConfig{}, ctx.KVGasConfig())
	suite.Equal(storetypes.GasConfig{}, ctx.TransientKVGasConfig())
}
