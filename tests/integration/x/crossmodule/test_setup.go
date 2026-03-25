// Package crossmodule provides integration tests that validate state consistency
// across x/erc20, x/precisebank, and x/vm module interactions (H-05).
package crossmodule

import (
	"github.com/stretchr/testify/suite"

	testconstants "github.com/zenanetwork/zena/testutil/constants"
	"github.com/zenanetwork/zena/testutil/integration/evm/factory"
	"github.com/zenanetwork/zena/testutil/integration/evm/grpc"
	"github.com/zenanetwork/zena/testutil/integration/evm/network"
	"github.com/zenanetwork/zena/testutil/keyring"
)

// CrossModuleTestSuite validates cross-module state consistency between
// ERC20, PreciseBank, and VM modules.
type CrossModuleTestSuite struct {
	suite.Suite

	create  network.CreateEvmApp
	options []network.ConfigOption
	network *network.UnitTestNetwork
	handler grpc.Handler
	keyring keyring.Keyring
	factory factory.TxFactory
}

// NewCrossModuleTestSuite creates a new test suite for cross-module integration tests.
func NewCrossModuleTestSuite(create network.CreateEvmApp, options ...network.ConfigOption) *CrossModuleTestSuite {
	return &CrossModuleTestSuite{
		create:  create,
		options: options,
	}
}

func (s *CrossModuleTestSuite) SetupTest() {
	keys := keyring.New(2)

	options := []network.ConfigOption{
		network.WithChainID(testconstants.SixDecimalsChainID),
		network.WithPreFundedAccounts(keys.GetAllAccAddrs()...),
	}
	options = append(options, s.options...)
	nw := network.NewUnitTestNetwork(s.create, options...)
	gh := grpc.NewIntegrationHandler(nw)
	tf := factory.New(nw, gh)

	s.network = nw
	s.factory = tf
	s.handler = gh
	s.keyring = keys
}
