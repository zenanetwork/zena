package distribution

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	cmn "github.com/zenanetwork/zena/precompiles/common"
	"github.com/zenanetwork/zena/precompiles/distribution"
	"github.com/zenanetwork/zena/precompiles/testutil"
	testconstants "github.com/zenanetwork/zena/testutil/constants"
	"github.com/zenanetwork/zena/testutil/integration/evm/network"
	utiltx "github.com/zenanetwork/zena/testutil/tx"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (s *PrecompileTestSuite) TestSetWithdrawAddress() {
	var ctx sdk.Context
	method := s.precompile.Methods[distribution.SetWithdrawAddressMethod]
	newWithdrawerAddr := utiltx.GenerateAddress()

	testCases := []struct {
		name        string
		malleate    func() []interface{}
		postCheck   func()
		gas         uint64
		expError    bool
		errContains string
	}{
		{
			"fail - empty input args",
			func() []interface{} {
				return []interface{}{}
			},
			func() {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 2, 0),
		},
		{
			"fail - invalid delegator address",
			func() []interface{} {
				return []interface{}{
					"",
					s.keyring.GetAddr(0).String(),
				}
			},
			func() {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidDelegator, ""),
		},
		{
			"fail - invalid withdrawer address",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					nil,
				}
			},
			func() {},
			200000,
			true,
			"invalid withdraw address: empty address string is not allowed: invalid address",
		},
		{
			"success - using the same address withdrawer address",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					s.keyring.GetAddr(0).String(),
				}
			},
			func() {
				withdrawerAddr, err := s.network.App.GetDistrKeeper().GetDelegatorWithdrawAddr(ctx, s.keyring.GetAccAddr(0))
				s.Require().NoError(err)
				s.Require().Equal(withdrawerAddr.String(), s.keyring.GetAccAddr(0).String())
			},
			20000,
			false,
			"",
		},
		{
			"success - using a different withdrawer address",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					newWithdrawerAddr.String(),
				}
			},
			func() {
				withdrawerAddr, err := s.network.App.GetDistrKeeper().GetDelegatorWithdrawAddr(ctx, s.keyring.GetAddr(0).Bytes())
				s.Require().NoError(err)
				s.Require().Equal(withdrawerAddr.Bytes(), newWithdrawerAddr.Bytes())
			},
			20000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			ctx = s.network.GetContext()

			var contract *vm.Contract
			contract, ctx = testutil.NewPrecompileContract(s.T(), ctx, s.keyring.GetAddr(0), s.precompile.Address(), tc.gas)

			_, err := s.precompile.SetWithdrawAddress(ctx, contract, s.network.GetStateDB(), &method, tc.malleate())

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
			} else {
				s.Require().NoError(err)
				tc.postCheck()
			}
		})
	}
}

func (s *PrecompileTestSuite) TestWithdrawDelegatorReward() {
	var (
		ctx sdk.Context
		err error
	)
	method := s.precompile.Methods[distribution.WithdrawDelegatorRewardMethod]

	testCases := []struct {
		name        string
		malleate    func(val stakingtypes.Validator) []interface{}
		postCheck   func(data []byte)
		gas         uint64
		expError    bool
		errContains string
	}{
		{
			"fail - empty input args",
			func(stakingtypes.Validator) []interface{} {
				return []interface{}{}
			},
			func([]byte) {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 2, 0),
		},
		{
			"fail - invalid delegator address",
			func(val stakingtypes.Validator) []interface{} {
				return []interface{}{
					"",
					val.OperatorAddress,
				}
			},
			func([]byte) {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidDelegator, ""),
		},
		{
			"fail - invalid validator address",
			func(stakingtypes.Validator) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					nil,
				}
			},
			func([]byte) {},
			200000,
			true,
			"invalid validator address",
		},
		{
			"success - withdraw rewards from a single validator without commission",
			func(val stakingtypes.Validator) []interface{} {
				ctx, err = s.prepareStakingRewards(
					ctx,
					stakingRewards{
						Validator: val,
						Delegator: s.keyring.GetAccAddr(0),
						RewardAmt: testRewardsAmt,
					},
				)
				s.Require().NoError(err, "failed to unpack output")
				return []interface{}{
					s.keyring.GetAddr(0),
					val.OperatorAddress,
				}
			},
			func(data []byte) {
				var coins []cmn.Coin
				err := s.precompile.UnpackIntoInterface(&coins, distribution.WithdrawDelegatorRewardMethod, data)
				s.Require().NoError(err, "failed to unpack output")
				s.Require().Equal(coins[0].Denom, testconstants.ExampleAttoDenom)
				s.Require().Equal(coins[0].Amount.Int64(), expRewardsAmt.Int64())
				// Check bank balance after the withdrawal of rewards
				balance := s.network.App.GetBankKeeper().GetBalance(ctx, s.keyring.GetAddr(0).Bytes(), testconstants.ExampleAttoDenom)
				s.Require().True(balance.Amount.GT(network.PrefundedAccountInitialBalance))
			},
			20000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			ctx = s.network.GetContext()

			var contract *vm.Contract
			contract, ctx = testutil.NewPrecompileContract(s.T(), ctx, s.keyring.GetAddr(0), s.precompile.Address(), tc.gas)

			args := tc.malleate(s.network.GetValidators()[0])
			bz, err := s.precompile.WithdrawDelegatorReward(ctx, contract, s.network.GetStateDB(), &method, args)

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
			} else {
				s.Require().NoError(err)
				tc.postCheck(bz)
			}
		})
	}
}

func (s *PrecompileTestSuite) TestWithdrawValidatorCommission() {
	var (
		ctx         sdk.Context
		prevBalance sdk.Coin
	)
	method := s.precompile.Methods[distribution.WithdrawDelegatorRewardMethod]

	testCases := []struct {
		name        string
		malleate    func(operatorAddress string) []interface{}
		postCheck   func(data []byte)
		gas         uint64
		expError    bool
		errContains string
	}{
		{
			"fail - empty input args",
			func(string) []interface{} {
				return []interface{}{}
			},
			func([]byte) {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 1, 0),
		},
		{
			"fail - invalid validator address",
			func(string) []interface{} {
				return []interface{}{
					nil,
				}
			},
			func([]byte) {},
			200000,
			true,
			"empty address string is not allowed",
		},
		{
			"success - withdraw all commission from a single validator",
			func(operatorAddress string) []interface{} {
				valAddr, err := sdk.ValAddressFromBech32(operatorAddress)
				s.Require().NoError(err)
				amt := math.LegacyNewDecWithPrec(1000000000000000000, 1)
				valCommission := sdk.DecCoins{sdk.NewDecCoinFromDec(testconstants.ExampleAttoDenom, amt)}
				// set outstanding rewards
				s.Require().NoError(s.network.App.GetDistrKeeper().SetValidatorOutstandingRewards(ctx, valAddr, types.ValidatorOutstandingRewards{Rewards: valCommission}))
				// set commission
				s.Require().NoError(s.network.App.GetDistrKeeper().SetValidatorAccumulatedCommission(ctx, valAddr, types.ValidatorAccumulatedCommission{Commission: valCommission}))

				// fund distr mod to pay for rewards + commission
				coins := sdk.NewCoins(sdk.NewCoin(testconstants.ExampleAttoDenom, amt.Mul(math.LegacyNewDec(2)).RoundInt()))
				err = s.mintCoinsForDistrMod(ctx, coins)
				s.Require().NoError(err)
				return []interface{}{
					operatorAddress,
				}
			},
			func(data []byte) {
				var coins []cmn.Coin
				amt := math.NewInt(100000000000000000)
				err := s.precompile.UnpackIntoInterface(&coins, distribution.WithdrawValidatorCommissionMethod, data)
				s.Require().NoError(err, "failed to unpack output")
				s.Require().Equal(coins[0].Denom, testconstants.ExampleAttoDenom)
				s.Require().Equal(coins[0].Amount, amt.BigInt())

				// Check bank balance after the withdrawal of commission
				valAddr, err := sdk.ValAddressFromBech32(s.network.GetValidators()[0].GetOperator())
				s.Require().NoError(err)
				balance := s.network.App.GetBankKeeper().GetBalance(ctx, valAddr.Bytes(), testconstants.ExampleAttoDenom)
				s.Require().Equal(balance.Amount, prevBalance.Amount.Add(amt))
				s.Require().Equal(balance.Denom, testconstants.ExampleAttoDenom)
			},
			20000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			ctx = s.network.GetContext()

			valAddr, err := sdk.ValAddressFromBech32(s.network.GetValidators()[0].GetOperator())
			s.Require().NoError(err)

			prevBalance = s.network.App.GetBankKeeper().GetBalance(ctx, valAddr.Bytes(), testconstants.ExampleAttoDenom)

			validatorAddress := common.BytesToAddress(valAddr.Bytes())
			var contract *vm.Contract
			contract, ctx = testutil.NewPrecompileContract(s.T(), ctx, validatorAddress, s.precompile.Address(), tc.gas)

			bz, err := s.precompile.WithdrawValidatorCommission(ctx, contract, s.network.GetStateDB(), &method, tc.malleate(s.network.GetValidators()[0].OperatorAddress))

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
			} else {
				s.Require().NoError(err)
				tc.postCheck(bz)
			}
		})
	}
}

func (s *PrecompileTestSuite) TestClaimRewards() {
	var (
		ctx         sdk.Context
		prevBalance sdk.Coin
	)
	method := s.precompile.Methods[distribution.ClaimRewardsMethod]

	testCases := []struct {
		name        string
		malleate    func() []interface{}
		postCheck   func(data []byte)
		gas         uint64
		expError    bool
		errContains string
	}{
		{
			"fail - empty input args",
			func() []interface{} {
				return []interface{}{}
			},
			func([]byte) {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 2, 0),
		},
		{
			"fail - invalid delegator address",
			func() []interface{} {
				return []interface{}{
					nil,
					10,
				}
			},
			func([]byte) {},
			200000,
			true,
			"invalid delegator address",
		},
		{
			"fail - invalid type for maxRetrieve: expected uint32",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(100000000000000000),
				}
			},
			func([]byte) {},
			200000,
			true,
			"invalid type for maxRetrieve: expected uint32",
		},
		{
			"fail - too many retrieved results",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					uint32(32_000_000),
				}
			},
			func([]byte) {},
			200000,
			true,
			"maxRetrieve (32000000) parameter exceeds the maximum number of validators (100)",
		},
		{
			"success - withdraw from all validators - 3",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					uint32(3),
				}
			},
			func(_ []byte) {
				balance := s.network.App.GetBankKeeper().GetBalance(ctx, s.keyring.GetAccAddr(0), testconstants.ExampleAttoDenom)
				// rewards from 3 validators - 5% commission
				expRewards := expRewardsAmt.Mul(math.NewInt(3))
				s.Require().Equal(balance.Amount, prevBalance.Amount.Add(expRewards))
			},
			20000,
			false,
			"",
		},
		{
			"pass - withdraw from validators with maxRetrieve higher than number of validators",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					uint32(10),
				}
			},
			func([]byte) {
				balance := s.network.App.GetBankKeeper().GetBalance(ctx, s.keyring.GetAccAddr(0), testconstants.ExampleAttoDenom)
				// rewards from 3 validators - 5% commission
				expRewards := expRewardsAmt.Mul(math.NewInt(3))
				s.Require().Equal(balance.Amount, prevBalance.Amount.Add(expRewards))
			},
			20000,
			false,
			"",
		},
		{
			"success - withdraw from only 1 validator",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					uint32(1),
				}
			},
			func([]byte) {
				balance := s.network.App.GetBankKeeper().GetBalance(ctx, s.keyring.GetAccAddr(0), testconstants.ExampleAttoDenom)
				s.Require().Equal(balance.Amount, prevBalance.Amount.Add(expRewardsAmt))
			},
			20000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			ctx = s.network.GetContext()

			var (
				contract *vm.Contract
				err      error
			)
			addr := s.keyring.GetAddr(0)
			contract, ctx = testutil.NewPrecompileContract(s.T(), ctx, addr, s.precompile.Address(), tc.gas)

			validators := s.network.GetValidators()
			srs := make([]stakingRewards, len(validators))
			for i, val := range validators {
				srs[i] = stakingRewards{
					Delegator: addr.Bytes(),
					Validator: val,
					RewardAmt: testRewardsAmt,
				}
			}

			ctx, err = s.prepareStakingRewards(ctx, srs...)
			s.Require().NoError(err)

			// get previous balance to compare final balance in the postCheck func
			prevBalance = s.network.App.GetBankKeeper().GetBalance(ctx, addr.Bytes(), testconstants.ExampleAttoDenom)

			bz, err := s.precompile.ClaimRewards(ctx, contract, s.network.GetStateDB(), &method, tc.malleate())

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
			} else {
				s.Require().NoError(err)
				tc.postCheck(bz)
			}
		})
	}
}

func (s *PrecompileTestSuite) TestFundCommunityPool() {
	var ctx sdk.Context
	method := s.precompile.Methods[distribution.FundCommunityPoolMethod]

	testCases := []struct {
		name        string
		malleate    func() []interface{}
		postCheck   func(data []byte)
		gas         uint64
		expError    bool
		errContains string
	}{
		{
			"fail - empty input args",
			func() []interface{} {
				return []interface{}{}
			},
			func([]byte) {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 2, 0),
		},
		{
			"fail - invalid depositor address",
			func() []interface{} {
				return []interface{}{
					nil,
					big.NewInt(1e18),
				}
			},
			func([]byte) {},
			200000,
			true,
			"invalid hex address address",
		},
		{
			"success - fund the community pool 1 ATOM",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					[]cmn.Coin{
						{
							Denom:  testconstants.ExampleAttoDenom,
							Amount: big.NewInt(1e18),
						},
					},
				}
			},
			func([]byte) {
				pool, err := s.network.App.GetDistrKeeper().FeePool.Get(ctx)
				s.Require().NoError(err)
				coins := pool.CommunityPool
				expectedAmount := new(big.Int).Mul(big.NewInt(1e18), new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(math.LegacyPrecision)), nil))
				s.Require().Equal(expectedAmount, coins.AmountOf(testconstants.ExampleAttoDenom).BigInt())
				userBalance := s.network.App.GetBankKeeper().GetBalance(ctx, s.keyring.GetAddr(0).Bytes(), testconstants.ExampleAttoDenom)
				s.Require().Equal(network.PrefundedAccountInitialBalance.Sub(math.NewInt(1e18)), userBalance.Amount)
			},
			20000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			ctx = s.network.GetContext()

			var contract *vm.Contract
			contract, ctx = testutil.NewPrecompileContract(s.T(), ctx, s.keyring.GetAddr(0), s.precompile.Address(), tc.gas)

			// Sanity check to make sure the starting balance is always 100k ATOM
			balance := s.network.App.GetBankKeeper().GetBalance(ctx, s.keyring.GetAddr(0).Bytes(), testconstants.ExampleAttoDenom)
			s.Require().Equal(balance.Amount, network.PrefundedAccountInitialBalance)

			bz, err := s.precompile.FundCommunityPool(ctx, contract, s.network.GetStateDB(), &method, tc.malleate())

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
			} else {
				s.Require().NoError(err)
				tc.postCheck(bz)
			}
		})
	}
}

func (s *PrecompileTestSuite) TestDepositValidatorRewardsPoolMethod() {
	var ctx sdk.Context
	method := s.precompile.Methods[distribution.DepositValidatorRewardsPoolMethod]

	testCases := []struct {
		name        string
		malleate    func(val stakingtypes.Validator) []interface{}
		postCheck   func(data []byte)
		gas         uint64
		expError    bool
		errContains string
	}{
		{
			"fail - empty input args",
			func(_ stakingtypes.Validator) []interface{} {
				return []interface{}{}
			},
			func([]byte) {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 3, 0),
		},
		{
			"fail - invalid depositor address",
			func(val stakingtypes.Validator) []interface{} {
				return []interface{}{
					"invalidAddress",
					val.OperatorAddress,
					[]cmn.Coin{
						{
							Denom:  testconstants.ExampleAttoDenom,
							Amount: big.NewInt(1e18),
						},
					},
				}
			},
			func([]byte) {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidHexAddress, "invalidAddress"),
		},
		{
			"fail - empty validator address",
			func(val stakingtypes.Validator) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					"",
					[]cmn.Coin{
						{
							Denom:  testconstants.ExampleAttoDenom,
							Amount: big.NewInt(1e18),
						},
					},
				}
			},
			func([]byte) {},
			200000,
			true,
			"empty address string is not allowed",
		},
		{
			"fail - invalid amount",
			func(val stakingtypes.Validator) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					val.OperatorAddress,
					"invalidAmount",
				}
			},
			func([]byte) {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidAmount, "invalidAmount"),
		},
		{
			"success - deposit rewards to the validator pool",
			func(val stakingtypes.Validator) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					val.OperatorAddress,
					[]cmn.Coin{
						{
							Denom:  testconstants.ExampleAttoDenom,
							Amount: big.NewInt(1e18),
						},
					},
				}
			},
			func(data []byte) {
				// check data is true
				var success bool
				err := s.precompile.UnpackIntoInterface(&success, distribution.DepositValidatorRewardsPoolMethod, data)
				s.Require().NoError(err, "failed to unpack output")
				s.Require().True(success, "expected true, got false")

				val := s.network.GetValidators()[0]
				valCodec := address.NewBech32Codec("cosmosvaloper")
				valBz, err := valCodec.StringToBytes(val.GetOperator())
				s.Require().NoError(err)

				depositCoins := sdk.DecCoins{
					{Denom: testconstants.ExampleAttoDenom, Amount: math.LegacyNewDecFromBigInt(big.NewInt(1e18))},
				}
				expectedValCommission := depositCoins.MulDec(val.GetCommission())
				expectedCurrentRewards := depositCoins.Sub(expectedValCommission)
				expectedOutstandingRewards := depositCoins

				// check validation commission
				valCommission, err := s.network.App.GetDistrKeeper().GetValidatorAccumulatedCommission(ctx, valBz)
				s.Require().NoError(err)
				s.Require().Equal(expectedValCommission, valCommission.Commission)

				// check current rewards
				currentRewards, err := s.network.App.GetDistrKeeper().GetValidatorCurrentRewards(ctx, valBz)
				s.Require().NoError(err)
				s.Require().Equal(expectedCurrentRewards, currentRewards.Rewards)

				// check outstanding rewards
				outstandingRewards, err := s.network.App.GetDistrKeeper().GetValidatorOutstandingRewards(ctx, valBz)
				s.Require().NoError(err)
				s.Require().Equal(expectedOutstandingRewards, outstandingRewards.Rewards)

				// check bank balance after the deposit
				balance := s.network.App.GetBankKeeper().GetBalance(ctx, s.keyring.GetAddr(0).Bytes(), testconstants.ExampleAttoDenom)
				s.Require().Equal(balance.Amount, network.PrefundedAccountInitialBalance.Sub(math.NewInt(1e18)))
				s.Require().Equal(balance.Denom, testconstants.ExampleAttoDenom)
			},
			20000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			ctx = s.network.GetContext()

			var contract *vm.Contract
			contract, ctx = testutil.NewPrecompileContract(s.T(), ctx, s.keyring.GetAddr(0), s.precompile.Address(), tc.gas)

			args := tc.malleate(s.network.GetValidators()[0])
			bz, err := s.precompile.DepositValidatorRewardsPool(ctx, contract, s.network.GetStateDB(), &method, args)

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
			} else {
				s.Require().NoError(err)
				tc.postCheck(bz)
			}
		})
	}
}
