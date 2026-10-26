package zenad

import (
	"context"
	"fmt"
	"slices"

	"github.com/zenanetwork/zena/config"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	upgradetypes "cosmossdk.io/x/upgrade/types"
)

// UpgradeName defines the on-chain upgrade name for the Zenanet mainnet
// stabilization upgrade.
//
// The upgrade carries the state machine changes accumulated since the
// mainnet genesis binary (Cosmos SDK v0.53.6, Cosmos EVM v0.6.0 StateDB port,
// official ibc-go transfer keeper, security audit fixes) and runs the
// following migrations:
//   - remove the unimplemented static precompiles (0x...0803, 0x...0807)
//     from the EVM active_static_precompiles parameter
//   - make sure the native aznnt bank denom metadata exists
//
// No store keys are added, renamed or deleted, so no store upgrades are needed.
const UpgradeName = "v0.6.0-stabilization"

func (app ZENAD) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)

			if err := app.removeUnimplementedStaticPrecompiles(sdkCtx); err != nil {
				return nil, err
			}

			app.ensureNativeDenomMetadata(sdkCtx)

			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		},
	)
}

// removeUnimplementedStaticPrecompiles drops every address listed in
// evmtypes.UnimplementedStaticPrecompiles from the active static precompiles.
// It is a no-op when none of them are active.
func (app ZENAD) removeUnimplementedStaticPrecompiles(ctx sdk.Context) error {
	params := app.EVMKeeper.GetParams(ctx)

	active := make([]string, 0, len(params.ActiveStaticPrecompiles))
	var removed []string
	for _, addr := range params.ActiveStaticPrecompiles {
		if slices.Contains(evmtypes.UnimplementedStaticPrecompiles, addr) {
			removed = append(removed, addr)
			continue
		}
		active = append(active, addr)
	}

	if len(removed) == 0 {
		ctx.Logger().Info("no unimplemented static precompiles active, skipping removal")
		return nil
	}

	params.ActiveStaticPrecompiles = active
	if err := app.EVMKeeper.SetParams(ctx, params); err != nil {
		return fmt.Errorf("failed to remove unimplemented static precompiles: %w", err)
	}

	ctx.Logger().Info("removed unimplemented static precompiles", "removed", removed, "active", active)
	return nil
}

// ensureNativeDenomMetadata writes the canonical aznnt metadata only when the
// chain has none. Existing metadata is never overwritten, because the EVM coin
// info (decimals, display denom) was derived from it at genesis; a mismatch is
// only logged so operators can review it.
func (app ZENAD) ensureNativeDenomMetadata(ctx sdk.Context) {
	expected := config.NativeDenomMetadata()

	current, found := app.BankKeeper.GetDenomMetaData(ctx, expected.Base)
	if !found {
		app.BankKeeper.SetDenomMetaData(ctx, expected)
		ctx.Logger().Info("set missing native denom metadata", "denom", expected.Base)
		return
	}

	if current.Display != expected.Display || current.Symbol != expected.Symbol {
		ctx.Logger().Error(
			"native denom metadata differs from the canonical values; leaving it unchanged",
			"denom", expected.Base,
			"display", current.Display, "expected_display", expected.Display,
			"symbol", current.Symbol, "expected_symbol", expected.Symbol,
		)
	}
}
