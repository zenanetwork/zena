package ics20

import (
	"testing"

	"github.com/stretchr/testify/suite"

	ibctesting "github.com/cosmos/ibc-go/v10/testing"
	"github.com/zenanetwork/zena"
	"github.com/zenanetwork/zena/precompiles/ics20"
	evmibctesting "github.com/zenanetwork/zena/testutil/ibc"
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

	evmAppA := s.chainA.App.(zena.EvmApp)
	s.chainAPrecompile, _ = ics20.NewPrecompile(
		evmAppA.GetBankKeeper(),
		*evmAppA.GetStakingKeeper(),
		evmAppA.GetTransferKeeper(),
		evmAppA.GetIBCKeeper().ChannelKeeper,
		evmAppA.GetEVMKeeper(),
	)
	s.chainABondDenom, _ = evmAppA.GetStakingKeeper().BondDenom(s.chainA.GetContext())
	evmAppB := s.chainB.App.(zena.EvmApp)
	s.chainBPrecompile, _ = ics20.NewPrecompile(
		evmAppB.GetBankKeeper(),
		*evmAppB.GetStakingKeeper(),
		evmAppB.GetTransferKeeper(),
		evmAppB.GetIBCKeeper().ChannelKeeper,
		evmAppB.GetEVMKeeper(),
	)
	s.chainBBondDenom, _ = evmAppB.GetStakingKeeper().BondDenom(s.chainB.GetContext())
}
