package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zenanetwork/zena/tests/integration/x/crossmodule"
)

func TestCrossModuleIntegrationSuite(t *testing.T) {
	suite.Run(t, crossmodule.NewCrossModuleTestSuite(CreateEvmd))
}
