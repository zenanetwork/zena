package eips_test

import (
	"github.com/zenanetwork/zena/zenad/tests/integration"
	"github.com/zenanetwork/zena/tests/integration/eips"
	"testing"
	//nolint:revive // dot imports are fine for Ginkgo
	//nolint:revive // dot imports are fine for Ginkgo
)

func TestEIPs(t *testing.T) {
	eips.RunTests(t, integration.CreateEvmd)
}
