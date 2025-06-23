package evidence

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zenanetwork/zena/tests/integration/precompiles/evidence"
	"github.com/zenanetwork/zena/zenad/tests/integration"
)

func TestEvidencePrecompileTestSuite(t *testing.T) {
	s := evidence.NewPrecompileTestSuite(integration.CreateEvmd)
	suite.Run(t, s)
}
