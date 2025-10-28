package eip7702

import (
	"testing"

	"github.com/zenanetwork/zena/zenad/tests/integration"
	"github.com/zenanetwork/zena/tests/integration/eip7702"
)

func TestEIP7702IntegrationTestSuite(t *testing.T) {
	eip7702.TestEIP7702IntegrationTestSuite(t, integration.CreateEvmd)
}
