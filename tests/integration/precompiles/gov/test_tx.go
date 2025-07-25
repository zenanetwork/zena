package gov

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	cmn "github.com/zenanetwork/zena/precompiles/common"
	"github.com/zenanetwork/zena/precompiles/gov"
	"github.com/zenanetwork/zena/precompiles/testutil"
	utiltx "github.com/zenanetwork/zena/testutil/tx"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *PrecompileTestSuite) TestVote() {
	var ctx sdk.Context
	method := s.precompile.Methods[gov.VoteMethod]
	newVoterAddr := utiltx.GenerateAddress()
	const proposalID uint64 = 1
	const option uint8 = 1
	const metadata = "metadata"

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
			fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 4, 0),
		},
		{
			"fail - invalid voter address",
			func() []interface{} {
				return []interface{}{
					"",
					proposalID,
					option,
					metadata,
				}
			},
			func() {},
			200000,
			true,
			"invalid voter address",
		},
		{
			"fail - invalid voter address",
			func() []interface{} {
				return []interface{}{
					common.Address{},
					proposalID,
					option,
					metadata,
				}
			},
			func() {},
			200000,
			true,
			"invalid voter address",
		},
		{
			"fail - using a different voter address",
			func() []interface{} {
				return []interface{}{
					newVoterAddr,
					proposalID,
					option,
					metadata,
				}
			},
			func() {},
			200000,
			true,
			"does not match the requester address",
		},
		{
			"fail - invalid vote option",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					proposalID,
					option + 10,
					metadata,
				}
			},
			func() {},
			200000,
			true,
			"invalid vote option",
		},
		{
			"success - vote proposal success",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					proposalID,
					option,
					metadata,
				}
			},
			func() {
				proposal, _ := s.network.App.GetGovKeeper().Proposals.Get(ctx, proposalID)
				_, _, tallyResult, err := s.network.App.GetGovKeeper().Tally(ctx, proposal)
				s.Require().NoError(err)
				s.Require().Equal(math.NewInt(3e18).String(), tallyResult.YesCount)
			},
			200000,
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

			_, err := s.precompile.Vote(ctx, contract, s.network.GetStateDB(), &method, tc.malleate())

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
			} else {
				s.Require().NoError(err)
				tc.postCheck()
			}
		})
	}
}

func (s *PrecompileTestSuite) TestVoteWeighted() {
	var ctx sdk.Context
	method := s.precompile.Methods[gov.VoteWeightedMethod]
	newVoterAddr := utiltx.GenerateAddress()
	const proposalID uint64 = 1
	const metadata = "metadata"

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
			fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 4, 0),
		},
		{
			"fail - invalid voter address",
			func() []interface{} {
				return []interface{}{
					"",
					proposalID,
					[]gov.WeightedVoteOption{},
					metadata,
				}
			},
			func() {},
			200000,
			true,
			"invalid voter address",
		},
		{
			"fail - using a different voter address",
			func() []interface{} {
				return []interface{}{
					newVoterAddr,
					proposalID,
					[]gov.WeightedVoteOption{},
					metadata,
				}
			},
			func() {},
			200000,
			true,
			"does not match the requester address",
		},
		{
			"fail - invalid vote option",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					proposalID,
					[]gov.WeightedVoteOption{{Option: 10, Weight: "1.0"}},
					metadata,
				}
			},
			func() {},
			200000,
			true,
			"invalid vote option",
		},
		{
			"fail - invalid weight sum",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					proposalID,
					[]gov.WeightedVoteOption{
						{Option: 1, Weight: "0.5"},
						{Option: 2, Weight: "0.6"},
					},
					metadata,
				}
			},
			func() {},
			200000,
			true,
			"total weight overflow 1.00",
		},
		{
			"success - vote weighted proposal",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					proposalID,
					[]gov.WeightedVoteOption{
						{Option: 1, Weight: "0.7"},
						{Option: 2, Weight: "0.3"},
					},
					metadata,
				}
			},
			func() {
				proposal, _ := s.network.App.GetGovKeeper().Proposals.Get(ctx, proposalID)
				_, _, tallyResult, err := s.network.App.GetGovKeeper().Tally(ctx, proposal)
				s.Require().NoError(err)
				s.Require().Equal("2100000000000000000", tallyResult.YesCount)
				s.Require().Equal("900000000000000000", tallyResult.AbstainCount)
			},
			200000,
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

			_, err := s.precompile.VoteWeighted(ctx, contract, s.network.GetStateDB(), &method, tc.malleate())

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
			} else {
				s.Require().NoError(err)
				tc.postCheck()
			}
		})
	}
}
