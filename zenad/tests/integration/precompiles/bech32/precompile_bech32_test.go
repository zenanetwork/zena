package bech32

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zenanetwork/zena/zenad/tests/integration"
	"github.com/zenanetwork/zena/tests/integration/precompiles/bech32"
)

func TestBech32PrecompileTestSuite(t *testing.T) {
	s := bech32.NewPrecompileTestSuite(integration.CreateEvmd)
	suite.Run(t, s)
}
