package integration

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/flags"

	dbm "github.com/cosmos/cosmos-db"
	ibctesting "github.com/cosmos/ibc-go/v10/testing"

	"github.com/zenanetwork/zena"
	"github.com/zenanetwork/zena/config"
	"github.com/zenanetwork/zena/zenad"
	srvflags "github.com/zenanetwork/zena/server/flags"
	"github.com/zenanetwork/zena/testutil/constants"
	feemarkettypes "github.com/zenanetwork/zena/x/feemarket/types"

	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simutils "github.com/cosmos/cosmos-sdk/testutil/sims"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// CreateEvmd creates an evm app for regular integration tests (non-mempool)
// This version uses a noop mempool to avoid state issues during transaction processing
func CreateEvmd(chainID string, evmChainID uint64, customBaseAppOptions ...func(*baseapp.BaseApp)) evm.EvmApp {
	defaultNodeHome, err := clienthelpers.GetNodeHomeDirectory(".zenad")
	if err != nil {
		panic(err)
	}

	db := dbm.NewMemDB()
	logger := log.NewNopLogger()
	loadLatest := true
	appOptions := NewAppOptionsWithFlagHomeAndChainID(defaultNodeHome, evmChainID)

	baseAppOptions := append(customBaseAppOptions, baseapp.SetChainID(chainID))

	return zenad.NewExampleApp(
		logger,
		db,
		nil,
		loadLatest,
		appOptions,
		baseAppOptions...,
	)
}

// SetupEvmd initializes a new zenad app with default genesis state.
// It is used in IBC integration tests to create a new zenad app instance.
func SetupEvmd() (ibctesting.TestingApp, map[string]json.RawMessage) {
	app := zenad.NewExampleApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		NewAppOptionsWithFlagHomeAndChainID("", constants.ExampleEIP155ChainID),
	)
	// disable base fee for testing
	genesisState := app.DefaultGenesis()
	fmGen := feemarkettypes.DefaultGenesisState()
	fmGen.Params.NoBaseFee = true
	genesisState[feemarkettypes.ModuleName] = app.AppCodec().MustMarshalJSON(fmGen)
	stakingGen := stakingtypes.DefaultGenesisState()
	stakingGen.Params.BondDenom = config.ExampleChainDenom
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGen)
	mintGen := minttypes.DefaultGenesisState()
	mintGen.Params.MintDenom = config.ExampleChainDenom
	genesisState[minttypes.ModuleName] = app.AppCodec().MustMarshalJSON(mintGen)

	return app, genesisState
}

func NewAppOptionsWithFlagHomeAndChainID(home string, evmChainID uint64) simutils.AppOptionsMap {
	return simutils.AppOptionsMap{
		flags.FlagHome:      home,
		srvflags.EVMChainID: evmChainID,
	}
}
