package mempool

import (
	"testing"

	"github.com/zenanetwork/zena/zenad/tests/integration"

	"github.com/stretchr/testify/suite"

	"github.com/zenanetwork/zena/tests/integration/mempool"
)

func TestMempoolIntegrationTestSuite(t *testing.T) {
	suite.Run(t, mempool.NewMempoolIntegrationTestSuite(integration.CreateEvmd))
}
