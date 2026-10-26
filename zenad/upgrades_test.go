package zenad

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zenanetwork/zena/config"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func setupUpgradeTestApp(t *testing.T) *ZENAD {
	t.Helper()
	evmtypes.NewEVMConfigurator().ResetTestConfig()
	return Setup(t, "zena_15150-1", config.ZenaMainnetEVMChainID)
}

func TestRemoveUnimplementedStaticPrecompiles(t *testing.T) {
	app := setupUpgradeTestApp(t)
	ctx := app.BaseApp.NewContext(false)

	// mimic the mainnet genesis created by the previous local_node.sh
	params := app.EVMKeeper.GetParams(ctx)
	params.ActiveStaticPrecompiles = append(
		append([]string{}, evmtypes.AvailableStaticPrecompiles...),
		evmtypes.UnimplementedStaticPrecompiles...,
	)
	require.NoError(t, app.EVMKeeper.SetParams(ctx, params))

	require.NoError(t, app.removeUnimplementedStaticPrecompiles(ctx))

	active := app.EVMKeeper.GetParams(ctx).ActiveStaticPrecompiles
	require.ElementsMatch(t, evmtypes.AvailableStaticPrecompiles, active)
	for _, addr := range evmtypes.UnimplementedStaticPrecompiles {
		require.NotContains(t, active, addr)
	}

	// running it again must be a no-op
	require.NoError(t, app.removeUnimplementedStaticPrecompiles(ctx))
	require.Equal(t, active, app.EVMKeeper.GetParams(ctx).ActiveStaticPrecompiles)
}

func TestEnsureNativeDenomMetadata(t *testing.T) {
	app := setupUpgradeTestApp(t)
	ctx := app.BaseApp.NewContext(false)

	// existing metadata is never overwritten
	custom := config.NativeDenomMetadata()
	custom.Description = "existing on-chain metadata"
	app.BankKeeper.SetDenomMetaData(ctx, custom)

	app.ensureNativeDenomMetadata(ctx)

	got, found := app.BankKeeper.GetDenomMetaData(ctx, config.ZenaChainDenom)
	require.True(t, found)
	require.Equal(t, custom, got)

	// the handler must not write any test denom metadata
	_, found = app.BankKeeper.GetDenomMetaData(ctx, "atest")
	require.False(t, found)
}

func TestEnsureNativeDenomMetadataWhenMissing(t *testing.T) {
	app := setupUpgradeTestApp(t)
	ctx := app.BaseApp.NewContext(false)

	// genesis of the test app stores aznnt metadata; drop it to simulate a missing entry
	baseKeeper, ok := app.BankKeeper.(bankkeeper.BaseKeeper)
	require.True(t, ok)
	require.NoError(t, baseKeeper.BaseSendKeeper.BaseViewKeeper.DenomMetadata.Remove(ctx, config.ZenaChainDenom))
	_, found := app.BankKeeper.GetDenomMetaData(ctx, config.ZenaChainDenom)
	require.False(t, found)

	app.ensureNativeDenomMetadata(ctx)

	got, found := app.BankKeeper.GetDenomMetaData(ctx, config.ZenaChainDenom)
	require.True(t, found)
	require.Equal(t, config.NativeDenomMetadata(), got)
}
