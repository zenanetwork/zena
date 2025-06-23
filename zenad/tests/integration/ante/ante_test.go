package ante

import (
	"testing"

	"github.com/zenanetwork/zena/tests/integration/ante"
	"github.com/zenanetwork/zena/zenad/tests/integration"
)

func TestAnte_Integration(t *testing.T) {
	ante.TestIntegrationAnteHandler(t, integration.CreateEvmd)
}
