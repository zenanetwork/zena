package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zenanetwork/zena/tests/integration/testutil"
)

func TestTestUtilTestSuite(t *testing.T) {
	s := testutil.NewTestSuite(CreateEvmd)
	suite.Run(t, s)
}
