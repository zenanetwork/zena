package ante

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zenanetwork/zena/tests/integration/ante"
	"github.com/zenanetwork/zena/zenad/tests/integration"
)

func TestEvmUnitAnteTestSuite(t *testing.T) {
	suite.Run(t, ante.NewEvmUnitAnteTestSuite(integration.CreateEvmd))
}
