package crossmodule

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	pbtypes "github.com/zenanetwork/zena/x/precisebank/types"
)

// TestPreciseBankFractionalBalanceConsistency verifies that the precisebank
// module correctly maintains fractional balances and that the conversion factor
// between 18-decimal EVM representation and 6-decimal Cosmos representation
// is applied consistently.
func (s *CrossModuleTestSuite) TestPreciseBankFractionalBalanceConsistency() {
	ctx := s.network.GetContext()
	pbKeeper := s.network.App.GetPreciseBankKeeper()
	bankKeeper := s.network.App.GetBankKeeper()

	sender := s.keyring.GetAccAddr(0)
	receiver := s.keyring.GetAccAddr(1)

	// Get conversion factor (10^12 for 18-dec to 6-dec)
	convFactor := pbtypes.ConversionFactor()
	extDenom := pbtypes.ExtendedCoinDenom()

	// Record initial balances
	senderBalBefore := bankKeeper.GetBalance(ctx, sender, extDenom)
	receiverBalBefore := bankKeeper.GetBalance(ctx, receiver, extDenom)

	// Test 1: Send an amount that is NOT a clean multiple of the conversion factor.
	// This forces the precisebank to handle fractional balances.
	fractionalAmount := convFactor.QuoRaw(3) // 1/3 of conversion factor
	sendCoin := sdk.NewCoin(extDenom, fractionalAmount)

	err := pbKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(sendCoin))
	s.Require().NoError(err, "precisebank SendCoins should succeed")

	// Verify balances changed by the exact amount
	senderBalAfter := bankKeeper.GetBalance(ctx, sender, extDenom)
	receiverBalAfter := bankKeeper.GetBalance(ctx, receiver, extDenom)

	senderDiff := senderBalBefore.Amount.Sub(senderBalAfter.Amount)
	receiverDiff := receiverBalAfter.Amount.Sub(receiverBalBefore.Amount)

	s.Require().Equal(fractionalAmount.String(), senderDiff.String(),
		"sender balance should decrease by exactly the sent amount")
	s.Require().Equal(fractionalAmount.String(), receiverDiff.String(),
		"receiver balance should increase by exactly the sent amount")

	// Test 2: Send back the same fractional amount — round-trip should result in zero net change
	err = pbKeeper.SendCoins(ctx, receiver, sender, sdk.NewCoins(sendCoin))
	s.Require().NoError(err, "reverse SendCoins should succeed")

	senderBalRoundTrip := bankKeeper.GetBalance(ctx, sender, extDenom)
	receiverBalRoundTrip := bankKeeper.GetBalance(ctx, receiver, extDenom)

	s.Require().Equal(senderBalBefore.Amount.String(), senderBalRoundTrip.Amount.String(),
		"sender balance should return to original after round-trip")
	s.Require().Equal(receiverBalBefore.Amount.String(), receiverBalRoundTrip.Amount.String(),
		"receiver balance should return to original after round-trip")
}

// TestPreciseBankConversionFactorBoundary verifies behavior at the exact
// conversion factor boundary (amounts that are exactly 1 integer unit).
func (s *CrossModuleTestSuite) TestPreciseBankConversionFactorBoundary() {
	ctx := s.network.GetContext()
	pbKeeper := s.network.App.GetPreciseBankKeeper()
	bankKeeper := s.network.App.GetBankKeeper()

	sender := s.keyring.GetAccAddr(0)
	receiver := s.keyring.GetAccAddr(1)

	extDenom := pbtypes.ExtendedCoinDenom()
	convFactor := pbtypes.ConversionFactor()

	receiverBalBefore := bankKeeper.GetBalance(ctx, receiver, extDenom)

	// Send exactly 1 conversion factor (should map to exactly 1 integer coin)
	exactAmount := convFactor
	sendCoin := sdk.NewCoin(extDenom, exactAmount)

	err := pbKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(sendCoin))
	s.Require().NoError(err, "sending exact conversion factor should succeed")

	receiverBalAfter := bankKeeper.GetBalance(ctx, receiver, extDenom)
	diff := receiverBalAfter.Amount.Sub(receiverBalBefore.Amount)
	s.Require().Equal(exactAmount.String(), diff.String(),
		"receiver should gain exactly 1 conversion factor unit")

	// Send exactly (conversion factor - 1) — should be entirely fractional
	almostOne := convFactor.SubRaw(1)
	sendCoin2 := sdk.NewCoin(extDenom, almostOne)

	receiverBal2Before := bankKeeper.GetBalance(ctx, receiver, extDenom)
	err = pbKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(sendCoin2))
	s.Require().NoError(err, "sending (convFactor-1) should succeed")

	receiverBal2After := bankKeeper.GetBalance(ctx, receiver, extDenom)
	diff2 := receiverBal2After.Amount.Sub(receiverBal2Before.Amount)
	s.Require().Equal(almostOne.String(), diff2.String(),
		"receiver should gain exactly (convFactor-1) units")
}

// TestPreciseBankTotalSupplyInvariant verifies that the total supply remains
// consistent across send operations with fractional amounts.
func (s *CrossModuleTestSuite) TestPreciseBankTotalSupplyInvariant() {
	ctx := s.network.GetContext()
	pbKeeper := s.network.App.GetPreciseBankKeeper()
	bankKeeper := s.network.App.GetBankKeeper()

	sender := s.keyring.GetAccAddr(0)
	receiver := s.keyring.GetAccAddr(1)

	extDenom := pbtypes.ExtendedCoinDenom()
	intDenom := pbtypes.IntegerCoinDenom()
	convFactor := pbtypes.ConversionFactor()

	// Record total supply before
	totalSupplyBefore := bankKeeper.GetSupply(ctx, intDenom)

	// Perform multiple fractional transfers that should trigger carry operations
	for i := int64(1); i <= 5; i++ {
		amount := convFactor.QuoRaw(7).MulRaw(i) // Various non-clean fractions
		if amount.IsPositive() {
			sendCoin := sdk.NewCoin(extDenom, amount)
			err := pbKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(sendCoin))
			s.Require().NoError(err, "fractional transfer %d should succeed", i)
		}
	}

	// Total supply of the integer coin should remain unchanged
	// (sends between accounts don't create or destroy coins)
	totalSupplyAfter := bankKeeper.GetSupply(ctx, intDenom)
	s.Require().Equal(totalSupplyBefore.Amount.String(), totalSupplyAfter.Amount.String(),
		"total supply of integer coin must not change during sends")

	// Verify fractional balances are valid (each < convFactor)
	pbKeeper.IterateFractionalBalances(ctx, func(addr sdk.AccAddress, bal sdkmath.Int) bool {
		s.Require().True(bal.LT(convFactor),
			"fractional balance %s for %s must be less than conversion factor %s",
			bal, addr, convFactor)
		s.Require().True(bal.GTE(sdkmath.ZeroInt()),
			"fractional balance %s for %s must be non-negative",
			bal, addr)
		return false
	})
}
