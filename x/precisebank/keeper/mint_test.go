package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zenanetwork/zena/x/precisebank/types"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func TestMintCoins_PanicValidations(t *testing.T) {
	// panic tests for invalid inputs

	tests := []struct {
		name            string
		recipientModule string
		setupFn         func(td testData)
		mintAmount      sdk.Coins
		wantPanic       string
	}{
		{
			"invalid module",
			"notamodule",
			func(td testData) {
				// Make module not found
				td.ak.EXPECT().
					GetModuleAccount(td.ctx, "notamodule").
					Return(nil).
					Once()
			},
			cs(c(types.IntegerCoinDenom(), 1000)),
			"module account notamodule does not exist: unknown address",
		},
		{
			"no permission",
			minttypes.ModuleName,
			func(td testData) {
				td.ak.EXPECT().
					GetModuleAccount(td.ctx, minttypes.ModuleName).
					Return(authtypes.NewModuleAccount(
						authtypes.NewBaseAccountWithAddress(sdk.AccAddress{1}),
						minttypes.ModuleName,
						// no mint permission
					)).
					Once()
			},
			cs(c(types.IntegerCoinDenom(), 1000)),
			"module account mint does not have permissions to mint tokens: unauthorized",
		},
		{
			"has mint permission",
			minttypes.ModuleName,
			func(td testData) {
				td.ak.EXPECT().
					GetModuleAccount(td.ctx, minttypes.ModuleName).
					Return(authtypes.NewModuleAccount(
						authtypes.NewBaseAccountWithAddress(sdk.AccAddress{1}),
						minttypes.ModuleName,
						// includes minter permission
						authtypes.Minter,
					)).
					Once()

				// Will call x/bank MintCoins coins
				td.bk.EXPECT().
					MintCoins(td.ctx, minttypes.ModuleName, cs(c(types.IntegerCoinDenom(), 1000))).
					Return(nil).
					Once()
			},
			cs(c(types.IntegerCoinDenom(), 1000)),
			"",
		},
		{
			"disallow minting to x/precisebank",
			types.ModuleName,
			func(td testData) {
				// No mock setup needed since this is checked before module
				// account checks
			},
			cs(c(types.IntegerCoinDenom(), 1000)),
			"module account precisebank cannot be minted to: unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := newMockedTestData(t)
			tt.setupFn(td)

			if tt.wantPanic != "" {
				require.PanicsWithError(t, tt.wantPanic, func() {
					_ = td.keeper.MintCoins(td.ctx, tt.recipientModule, tt.mintAmount)
				})
				return
			}

			require.NotPanics(t, func() {
				// Not testing errors, only panics for this test
				_ = td.keeper.MintCoins(td.ctx, tt.recipientModule, tt.mintAmount)
			})
		})
	}
}

func TestMintCoins_Errors(t *testing.T) {
	// returned errors, not panics

	tests := []struct {
		name            string
		recipientModule string
		setupFn         func(td testData)
		mintAmount      sdk.Coins
		wantError       string
	}{
		{
			"invalid coins",
			minttypes.ModuleName,
			func(td testData) {
				// Valid module account minter
				td.ak.EXPECT().
					GetModuleAccount(td.ctx, minttypes.ModuleName).
					Return(authtypes.NewModuleAccount(
						authtypes.NewBaseAccountWithAddress(sdk.AccAddress{1}),
						minttypes.ModuleName,
						// includes minter permission
						authtypes.Minter,
					)).
					Once()
			},
			sdk.Coins{sdk.Coin{
				Denom:  types.IntegerCoinDenom(),
				Amount: sdkmath.NewInt(-1000),
			}},
			fmt.Sprintf("-1000%s: invalid coins", types.IntegerCoinDenom()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := newMockedTestData(t)
			tt.setupFn(td)

			require.NotPanics(t, func() {
				err := td.keeper.MintCoins(td.ctx, tt.recipientModule, tt.mintAmount)

				if tt.wantError != "" {
					require.Error(t, err)
					require.EqualError(t, err, tt.wantError)
					return
				}

				require.NoError(t, err)
			})
		})
	}
}

func TestMintCoins_ExpectedCalls(t *testing.T) {
	// Tests the expected calls to the bank keeper when minting coins

	tests := []struct {
		name string
		// Only care about starting fractional balance.
		// MintCoins() doesn't care about the previous integer balance.
		startFractionalBalance sdkmath.Int
		mintAmount             sdk.Coins
		// account x/precisebank balance (fractional amount)
		wantPreciseBalance sdkmath.Int
	}{
		{
			"passthrough mint - integer denom",
			sdkmath.ZeroInt(),
			cs(c(types.IntegerCoinDenom(), 1000)),
			sdkmath.ZeroInt(),
		},

		{
			"passthrough mint - unrelated denom",
			sdkmath.ZeroInt(),
			cs(c("meow", 1000)),
			sdkmath.ZeroInt(),
		},
		{
			"no carry - 0 starting fractional",
			sdkmath.ZeroInt(),
			cs(c(types.ExtendedCoinDenom(), 1000)),
			sdkmath.NewInt(1000),
		},
		{
			"no carry - non-zero fractional",
			sdkmath.NewInt(1_000_000),
			cs(c(types.ExtendedCoinDenom(), 1000)),
			sdkmath.NewInt(1_001_000),
		},
		{
			"fractional carry",
			// max fractional amount
			types.ConversionFactor().SubRaw(1),
			cs(c(types.ExtendedCoinDenom(), 1)), // +1 to carry
			sdkmath.ZeroInt(),
		},
		{
			"fractional carry max",
			// max fractional amount + max fractional amount
			types.ConversionFactor().SubRaw(1),
			cs(ci(types.ExtendedCoinDenom(), types.ConversionFactor().SubRaw(1))),
			types.ConversionFactor().SubRaw(2),
		},
		{
			"integer with fractional no carry",
			sdkmath.NewInt(1234),
			// mint 100 fractional
			cs(c(types.ExtendedCoinDenom(), 100)),
			sdkmath.NewInt(1234 + 100),
		},
		{
			"integer with fractional carry",
			types.ConversionFactor().SubRaw(100),
			// mint 105 fractional to carry
			cs(c(types.ExtendedCoinDenom(), 105)),
			sdkmath.NewInt(5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := newMockedTestData(t)

			// Set initial fractional balance
			// Initial integer balance doesn't matter for this test
			moduleAddr := sdk.AccAddress{1}
			td.keeper.SetFractionalBalance(
				td.ctx,
				moduleAddr,
				tt.startFractionalBalance,
			)
			fBal := td.keeper.GetFractionalBalance(td.ctx, moduleAddr)
			require.Equal(t, tt.startFractionalBalance, fBal)

			// Always calls GetModuleAccount() to check if module exists &
			// has permission
			td.ak.EXPECT().
				GetModuleAccount(td.ctx, minttypes.ModuleName).
				Return(authtypes.NewModuleAccount(
					authtypes.NewBaseAccountWithAddress(
						moduleAddr,
					),
					minttypes.ModuleName,
					// Include minter permissions - not testing permission in
					// this test
					authtypes.Minter,
				)).
				Once()

			// ----------------------------------------
			// Separate passthrough and extended coins
			// Determine how much is passed through to x/bank
			passthroughCoins := tt.mintAmount

			found, extCoins := tt.mintAmount.Find(types.ExtendedCoinDenom())
			if found {
				// Remove extended coin from passthrough coins
				passthroughCoins = passthroughCoins.Sub(extCoins)
			} else {
				extCoins = sdk.NewCoin(types.ExtendedCoinDenom(), sdkmath.ZeroInt())
			}

			require.Equalf(
				t,
				sdkmath.ZeroInt(),
				passthroughCoins.AmountOf(types.ExtendedCoinDenom()),
				"expected pass through coins should not include %v",
				types.ExtendedCoinDenom(),
			)

			// ----------------------------------------
			// Set expectations for minting passthrough coins
			// Only expect MintCoins to be called with passthrough coins with non-zero amount
			if !passthroughCoins.Empty() {
				t.Logf("Expecting MintCoins(%v)", passthroughCoins)

				td.bk.EXPECT().
					MintCoins(td.ctx, minttypes.ModuleName, passthroughCoins).
					Return(nil).
					Once()
			}

			// ----------------------------------------
			// Set expectations for reserve minting when fractional amounts
			// are minted & remainder is insufficient
			mintFractionalAmount := extCoins.Amount.Mod(types.ConversionFactor())
			currentRemainder := td.keeper.GetRemainderAmount(td.ctx)

			causesIntegerCarry := fBal.Add(mintFractionalAmount).GTE(types.ConversionFactor())
			remainderEnough := currentRemainder.GTE(mintFractionalAmount)

			// Optimization: Carry & insufficient remainder is directly minted
			if causesIntegerCarry && !remainderEnough {
				extCoins = extCoins.AddAmount(types.ConversionFactor())
			}

			// ----------------------------------------
			// Set expectations for minting fractional coins
			if !extCoins.IsNil() && extCoins.IsPositive() {
				td.ak.EXPECT().
					GetModuleAddress(minttypes.ModuleName).
					Return(moduleAddr).
					Once()

				// Initial integer balance is always 0 for this test
				mintIntegerAmount := extCoins.Amount.Quo(types.ConversionFactor())

				// Minted coins does NOT include roll-over, simply excludes
				mintCoins := cs(ci(types.IntegerCoinDenom(), mintIntegerAmount))

				// Only expect MintCoins to be called with mint coins with
				// non-zero amount.
				// Will fail if x/bank MintCoins is called with empty coins
				if !mintCoins.Empty() {
					t.Logf("Expecting MintCoins(%v)", mintCoins)

					td.bk.EXPECT().
						MintCoins(td.ctx, minttypes.ModuleName, mintCoins).
						Return(nil).
						Once()
				}
			}

			if causesIntegerCarry && remainderEnough {
				td.bk.EXPECT().
					SendCoinsFromModuleToModule(
						td.ctx,
						types.ModuleName,
						minttypes.ModuleName,
						cs(c(types.IntegerCoinDenom(), 1)),
					).
					Return(nil).
					Once()
			}

			if !remainderEnough && !causesIntegerCarry {
				reserveMintCoins := cs(c(types.IntegerCoinDenom(), 1))
				td.bk.EXPECT().
					// Mints to x/precisebank
					MintCoins(td.ctx, types.ModuleName, reserveMintCoins).
					Return(nil).
					Once()
			}

			// ----------------------------------------
			// Actual call after all setup and expectations
			require.NotPanics(t, func() {
				err := td.keeper.MintCoins(td.ctx, minttypes.ModuleName, tt.mintAmount)
				require.NoError(t, err)
			})

			// Check final fractional balance
			fBal = td.keeper.GetFractionalBalance(td.ctx, moduleAddr)
			require.Equal(t, tt.wantPreciseBalance, fBal)
		})
	}
}
