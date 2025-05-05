package ibctesting

import (
	"encoding/json"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/zenanetwork/zena/zenad"
	"github.com/zenanetwork/zena/zenad/testutil"
	feemarkettypes "github.com/zenanetwork/zena/x/feemarket/types"
	ibctesting "github.com/cosmos/ibc-go/v10/testing"

	"cosmossdk.io/log"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
)

func SetupExampleApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	app := zenad.NewExampleApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		simtestutil.EmptyAppOptions{},
		testutil.NoOpEvmAppOptions,
	)
	// disable base fee for testing
	genesisState := app.DefaultGenesis()
	fmGen := feemarkettypes.DefaultGenesisState()
	fmGen.Params.NoBaseFee = true
	genesisState[feemarkettypes.ModuleName] = app.AppCodec().MustMarshalJSON(fmGen)

	return app, genesisState
}
