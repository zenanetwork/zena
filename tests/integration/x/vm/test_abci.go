package vm

import (
	"github.com/zenanetwork/zena/testutil/integration/evm/network"
	testkeyring "github.com/zenanetwork/zena/testutil/keyring"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"
)

func (s *KeeperTestSuite) TestEndBlock() {
	keyring := testkeyring.New(2)
	unitNetwork := network.NewUnitTestNetwork(
		s.Create,
		network.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
	)
	ctx := unitNetwork.GetContext()
	preEventManager := ctx.EventManager()
	s.Require().Equal(0, len(preEventManager.Events()))

	err := unitNetwork.App.GetEVMKeeper().EndBlock(ctx)
	s.Require().NoError(err)

	postEventManager := unitNetwork.GetContext().EventManager()
	// should emit 1 EventTypeBlockBloom event on EndBlock
	s.Require().Equal(1, len(postEventManager.Events()))
	s.Require().Equal(evmtypes.EventTypeBlockBloom, postEventManager.Events()[0].Type)
}
