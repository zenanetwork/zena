package distribution

import (
	"github.com/stretchr/testify/suite"

	"github.com/zenanetwork/zena/precompiles/distribution"
	testconstants "github.com/zenanetwork/zena/testutil/constants"
	"github.com/zenanetwork/zena/testutil/integration/evm/factory"
	"github.com/zenanetwork/zena/testutil/integration/evm/grpc"
	"github.com/zenanetwork/zena/testutil/integration/evm/network"
	testkeyring "github.com/zenanetwork/zena/testutil/keyring"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type PrecompileTestSuite struct {
	suite.Suite

	create      network.CreateEvmApp
	options     []network.ConfigOption
	network     *network.UnitTestNetwork
	factory     factory.TxFactory
	grpcHandler grpc.Handler
	keyring     testkeyring.Keyring

	precompile           *distribution.Precompile
	bondDenom            string
	baseDenom            string
	otherDenoms          []string
	validatorsKeys       []testkeyring.Key
	withValidatorSlashes bool
}

func NewPrecompileTestSuite(
	create network.CreateEvmApp,
	options ...network.ConfigOption,
) *PrecompileTestSuite {
	return &PrecompileTestSuite{
		create:  create,
		options: options,
	}
}

func (s *PrecompileTestSuite) SetupTest() {
	keyring := testkeyring.New(2)
	s.validatorsKeys = generateKeys(3)
	customGen := network.CustomGenesisState{}

	// set some slashing events for integration test
	distrGen := distrtypes.DefaultGenesisState()
	if s.withValidatorSlashes {
		distrGen.ValidatorSlashEvents = []distrtypes.ValidatorSlashEventRecord{
			{
				ValidatorAddress:    sdk.ValAddress(s.validatorsKeys[0].Addr.Bytes()).String(),
				Height:              0,
				Period:              1,
				ValidatorSlashEvent: distrtypes.NewValidatorSlashEvent(1, math.LegacyNewDecWithPrec(5, 2)),
			},
			{
				ValidatorAddress:    sdk.ValAddress(s.validatorsKeys[0].Addr.Bytes()).String(),
				Height:              1,
				Period:              1,
				ValidatorSlashEvent: distrtypes.NewValidatorSlashEvent(1, math.LegacyNewDecWithPrec(5, 2)),
			},
		}
	}
	customGen[distrtypes.ModuleName] = distrGen

	// set non-zero inflation for rewards to accrue (use defaults from SDK for values)
	mintGen := minttypes.DefaultGenesisState()
	mintGen.Params.MintDenom = testconstants.ExampleAttoDenom
	customGen[minttypes.ModuleName] = mintGen

	operatorsAddr := make([]sdk.AccAddress, 3)
	for i, k := range s.validatorsKeys {
		operatorsAddr[i] = k.AccAddr
	}

	s.otherDenoms = []string{
		testconstants.OtherCoinDenoms[0],
		testconstants.OtherCoinDenoms[1],
	}

	options := []network.ConfigOption{
		network.WithOtherDenoms(testconstants.OtherCoinDenoms),
		network.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
		network.WithOtherDenoms(s.otherDenoms),
		network.WithCustomGenesis(customGen),
		network.WithValidatorOperators(operatorsAddr),
	}
	options = append(options, s.options...)
	nw := network.NewUnitTestNetwork(s.create, options...)
	grpcHandler := grpc.NewIntegrationHandler(nw)
	txFactory := factory.New(nw, grpcHandler)

	ctx := nw.GetContext()
	sk := nw.App.GetStakingKeeper()
	bondDenom, err := sk.BondDenom(ctx)
	if err != nil {
		panic(err)
	}

	s.bondDenom = bondDenom
	// TODO: check if this is correct?
	s.baseDenom = evmtypes.GetEVMCoinDenom()

	s.factory = txFactory
	s.grpcHandler = grpcHandler
	s.keyring = keyring
	s.network = nw
	s.precompile, err = distribution.NewPrecompile(
		s.network.App.GetDistrKeeper(),
		*s.network.App.GetStakingKeeper(),
		s.network.App.GetEVMKeeper(),
	)
	if err != nil {
		panic(err)
	}
}
