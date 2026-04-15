package zenad

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/module"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
)

// UpgradeName is the on-chain upgrade identifier for migrating Zena mainnet
// from v0.5.0 to v0.6.0.
//
// v0.6.0 ships:
//   - CRITICAL/HIGH security fixes (see claudedocs/security-audit/)
//   - Rebranding of test data (aatom → aznnt, cosmos → zenanet prefixes)
//   - Cosmos SDK v0.53.4 → v0.53.6 + StateDB event architecture port
//   - Test infrastructure improvements (testutil/setup, proto-gen fix)
//
// No on-chain state migration is currently required: the security fixes do
// not change storage layouts, and rebranding affects test-only data. If
// v0.7.0 introduces storage migrations, extend the handler below and add
// entries to StoreUpgrades.
const UpgradeName = "v0.5.0-to-v0.6.0"

func (app ZENAD) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// v0.5.0 → v0.6.0 does not need custom state migrations. Running
			// the module manager's migrations is sufficient to apply any
			// upstream Cosmos SDK module version bumps picked up by v0.53.6.
			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{},
		}
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
