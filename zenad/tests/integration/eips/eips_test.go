package eips

import (
	"testing"

	"github.com/zenanetwork/zena/tests/integration/eips"
	"github.com/zenanetwork/zena/zenad/tests/integration"
)

func Test_EIPs(t *testing.T) {
	eips.TestEIPs(t, integration.CreateEvmd)
}
