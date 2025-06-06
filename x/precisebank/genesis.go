package precisebank

import (
	"fmt"

	"github.com/zenanetwork/zena/x/precisebank/keeper"
	"github.com/zenanetwork/zena/x/precisebank/types"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the store state from a genesis state.
func InitGenesis(
	ctx sdk.Context,
	keeper keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	gs *types.GenesisState,
) {
	// Ensure the genesis state is valid
	if err := gs.Validate(); err != nil {
		panic(fmt.Sprintf("failed to validate %s genesis state: %s", types.ModuleName, err))
	}

	// Initialize module account
	if moduleAcc := ak.GetModuleAccount(ctx, types.ModuleName); moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// Check module balance matches sum of fractional balances + remainder
	// This is always a whole integer amount, as previously verified in
	// GenesisState.Validate()
	totalAmt := gs.TotalAmountWithRemainder()

	moduleAddr := ak.GetModuleAddress(types.ModuleName)
	moduleBal := bk.GetBalance(ctx, moduleAddr, types.IntegerCoinDenom())
	moduleBalExtended := moduleBal.Amount.Mul(types.ConversionFactor())

	// Compare balances in full precise extended amounts
	if !totalAmt.Equal(moduleBalExtended) {
		panic(fmt.Sprintf(
			"module account balance does not match sum of fractional balances and remainder, balance is %s but expected %v%s (%v%s)",
			moduleBal,
			totalAmt, types.ExtendedCoinDenom(),
			totalAmt.Quo(types.ConversionFactor()), types.IntegerCoinDenom(),
		))
	}

	// Set FractionalBalances in state
	for _, bal := range gs.Balances {
		addr := sdk.MustAccAddressFromBech32(bal.Address)

		keeper.SetFractionalBalance(ctx, addr, bal.Amount)
	}

	// Set remainder amount in state
	keeper.SetRemainderAmount(ctx, gs.Remainder)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	balances := types.FractionalBalances{}
	keeper.IterateFractionalBalances(ctx, func(addr sdk.AccAddress, amount sdkmath.Int) bool {
		balances = append(balances, types.NewFractionalBalance(addr.String(), amount))

		return false
	})

	remainder := keeper.GetRemainderAmount(ctx)

	return types.NewGenesisState(balances, remainder)
}
