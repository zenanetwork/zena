package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zenanetwork/zena/tests/integration/x/ibc/callbacks"
)

func TestIBCCallback(t *testing.T) {
	suite.Run(t, callbacks.NewKeeperTestSuite(CreateEvmd))
}
