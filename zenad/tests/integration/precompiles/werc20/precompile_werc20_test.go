package werc20

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zenanetwork/zena/tests/integration/precompiles/werc20"
	"github.com/zenanetwork/zena/zenad/tests/integration"
)

func TestWERC20PrecompileUnitTestSuite(t *testing.T) {
	s := werc20.NewPrecompileUnitTestSuite(integration.CreateEvmd)
	suite.Run(t, s)
}

func TestWERC20PrecompileIntegrationTestSuite(t *testing.T) {
	werc20.TestPrecompileIntegrationTestSuite(t, integration.CreateEvmd)
}
