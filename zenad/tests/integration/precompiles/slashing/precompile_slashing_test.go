package slashing

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zenanetwork/zena/tests/integration/precompiles/slashing"
	"github.com/zenanetwork/zena/zenad/tests/integration"
)

func TestSlashingPrecompileTestSuite(t *testing.T) {
	s := slashing.NewPrecompileTestSuite(integration.CreateEvmd)
	suite.Run(t, s)
}
