package types_test

import (
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zenanetwork/zena/config"
	testconstants "github.com/zenanetwork/zena/testutil/constants"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"
)

func TestMain(m *testing.M) {
	// Set bech32 prefixes for address validation
	cfg := sdk.GetConfig()
	config.SetBech32Prefixes(cfg)
	config.SetBip44CoinType(cfg)

	// precisebank uses SixDecimals (ConversionFactor = 10^12)
	coinInfo := testconstants.ExampleChainCoinInfo[testconstants.SixDecimalsChainID]
	if err := evmtypes.NewEVMConfigurator().
		WithEVMCoinInfo(coinInfo).
		Configure(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
