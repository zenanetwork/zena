package staking_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zenanetwork/zena/precompiles/staking"
	testconstants "github.com/zenanetwork/zena/testutil/constants"
	"github.com/zenanetwork/zena/testutil/integration/os/factory"
	"github.com/zenanetwork/zena/testutil/integration/os/grpc"
	testkeyring "github.com/zenanetwork/zena/testutil/integration/os/keyring"
	"github.com/zenanetwork/zena/testutil/integration/os/network"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type PrecompileTestSuite struct {
	suite.Suite

	network     *network.UnitTestNetwork
	factory     factory.TxFactory
	grpcHandler grpc.Handler
	keyring     testkeyring.Keyring

	bondDenom     string
	precompile    *staking.Precompile
	customGenesis bool
}

func TestPrecompileUnitTestSuite(t *testing.T) {
	suite.Run(t, new(PrecompileTestSuite))
}

func (s *PrecompileTestSuite) SetupTest() {
	keyring := testkeyring.New(2)
	customGenesis := network.CustomGenesisState{}
	// mint some coin to fee collector
	coins := sdk.NewCoins(sdk.NewCoin(testconstants.ExampleAttoDenom, sdkmath.NewInt(1000000000000000)))
	balances := []banktypes.Balance{
		{
			Address: authtypes.NewModuleAddress(authtypes.FeeCollectorName).String(),
			Coins:   coins,
		},
	}
	bankGenesis := banktypes.DefaultGenesisState()
	bankGenesis.Balances = balances
	customGenesis[banktypes.ModuleName] = bankGenesis
	cfgOpts := []network.ConfigOption{
		network.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
	}
	if s.customGenesis {
		cfgOpts = append(cfgOpts, network.WithCustomGenesis(customGenesis))
	}
	nw := network.NewUnitTestNetwork(
		cfgOpts...,
	)
	grpcHandler := grpc.NewIntegrationHandler(nw)
	txFactory := factory.New(nw, grpcHandler)

	ctx := nw.GetContext()
	sk := nw.App.StakingKeeper
	bondDenom, err := sk.BondDenom(ctx)
	if err != nil {
		panic(err)
	}

	s.bondDenom = bondDenom
	s.factory = txFactory
	s.grpcHandler = grpcHandler
	s.keyring = keyring
	s.network = nw

	if s.precompile, err = staking.NewPrecompile(
		*s.network.App.StakingKeeper,
	); err != nil {
		panic(err)
	}
}
