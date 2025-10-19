//go:build test
// +build test

package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/zenanetwork/zena/testutil/config"
	testconstants "github.com/zenanetwork/zena/testutil/constants"
	"github.com/zenanetwork/zena/x/precisebank/types"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestSumExtendedCoin(t *testing.T) {
	// Use SixDecimalsChainID where IntegerCoinDenom != ExtendedCoinDenom
	coinInfo := testconstants.ExampleChainCoinInfo[testconstants.SixDecimalsChainID]
	configurator := evmtypes.NewEVMConfigurator()
	configurator.ResetTestConfig()
	err := configurator.
		WithEVMCoinInfo(coinInfo).
		Configure()
	require.NoError(t, err)

	// Restore config after test - use new configurator to avoid "sealed" error
	t.Cleanup(func() {
		newConfigurator := evmtypes.NewEVMConfigurator()
		newConfigurator.ResetTestConfig()
		restoreCoinInfo := testconstants.ExampleChainCoinInfo[testconstants.TwelveDecimalsChainID]
		_ = newConfigurator.WithEVMCoinInfo(restoreCoinInfo).Configure()
	})

	tests := []struct {
		name string
		amt  sdk.Coins
		want sdk.Coin
	}{
		{
			"empty",
			sdk.NewCoins(),
			sdk.NewCoin(types.ExtendedCoinDenom(), sdkmath.ZeroInt()),
		},
		{
			"only integer",
			sdk.NewCoins(sdk.NewInt64Coin(types.IntegerCoinDenom(), 100)),
			sdk.NewCoin(types.ExtendedCoinDenom(), types.ConversionFactor().MulRaw(100)),
		},
		{
			"only extended",
			sdk.NewCoins(sdk.NewInt64Coin(types.ExtendedCoinDenom(), 100)),
			sdk.NewCoin(types.ExtendedCoinDenom(), sdkmath.NewInt(100)),
		},
		{
			"integer and extended",
			sdk.NewCoins(
				sdk.NewInt64Coin(types.IntegerCoinDenom(), 100),
				sdk.NewInt64Coin(types.ExtendedCoinDenom(), 100),
			),
			sdk.NewCoin(types.ExtendedCoinDenom(), types.ConversionFactor().MulRaw(100).AddRaw(100)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extVal := types.SumExtendedCoin(tt.amt)
			require.Equal(t, tt.want, extVal)
		})
	}
}
