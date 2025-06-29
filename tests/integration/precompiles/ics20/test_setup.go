package ics20

import (
	"testing"

	"github.com/stretchr/testify/suite"

	ibctesting "github.com/cosmos/ibc-go/v10/testing"
	"github.com/zenanetwork/zena/precompiles/ics20"
	evmibctesting "github.com/zenanetwork/zena/testutil/ibc"
	"github.com/zenanetwork/zena/zenad"
)

type PrecompileTestSuite struct {
	suite.Suite
	internalT   *testing.T
	coordinator *evmibctesting.Coordinator

	create           ibctesting.AppCreator
	chainA           *evmibctesting.TestChain
	chainAPrecompile *ics20.Precompile
	chainABondDenom  string
	chainB           *evmibctesting.TestChain
	chainBPrecompile *ics20.Precompile
	chainBBondDenom  string
}

//nolint:thelper // NewPrecompileTestSuite is not a helper function; it's an instantiation function for the test suite.
func NewPrecompileTestSuite(t *testing.T, create ibctesting.AppCreator) *PrecompileTestSuite {
	return &PrecompileTestSuite{
		internalT: t,
		create:    create,
	}
}

func (s *PrecompileTestSuite) SetupTest() {
	// Setup IBC
	if s.internalT == nil {
		s.internalT = s.T()
	}
	s.coordinator = evmibctesting.NewCoordinator(s.internalT, 2, 0, s.create)
	s.chainA = s.coordinator.GetChain(evmibctesting.GetEvmChainID(1))
	s.chainB = s.coordinator.GetChain(evmibctesting.GetEvmChainID(2))

	evmAppA := s.chainA.App.(*zenad.EVMD)
	s.chainAPrecompile, _ = ics20.NewPrecompile(
		evmAppA.BankKeeper,
		*evmAppA.StakingKeeper,
		evmAppA.TransferKeeper,
		evmAppA.IBCKeeper.ChannelKeeper,
		evmAppA.EVMKeeper,
	)
	s.chainABondDenom, _ = evmAppA.StakingKeeper.BondDenom(s.chainA.GetContext())
	evmAppB := s.chainB.App.(*zenad.EVMD)
	s.chainBPrecompile, _ = ics20.NewPrecompile(
		evmAppB.BankKeeper,
		*evmAppB.StakingKeeper,
		evmAppB.TransferKeeper,
		evmAppB.IBCKeeper.ChannelKeeper,
		evmAppB.EVMKeeper,
	)
	s.chainBBondDenom, _ = evmAppB.StakingKeeper.BondDenom(s.chainB.GetContext())
}
