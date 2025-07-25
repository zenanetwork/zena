package bech32

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"

	"github.com/zenanetwork/zena/precompiles/bech32"
	"github.com/zenanetwork/zena/testutil/config"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *PrecompileTestSuite) TestNewPrecompile() {
	testCases := []struct {
		name        string
		baseGas     uint64
		expPass     bool
		errContains string
	}{
		{
			"fail - new precompile with baseGas == 0",
			0,
			false,
			"baseGas cannot be zero",
		},
		{
			"success - new precompile with baseGas > 0",
			10,
			true,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// setup basic test suite
			s.SetupTest()
			p, err := bech32.NewPrecompile(tc.baseGas)
			if tc.expPass {
				s.Require().NoError(err)
				s.Require().NotNil(p)
				s.Require().Equal(tc.baseGas, p.RequiredGas([]byte{}))
			} else {
				s.Require().Error(err)
				s.Require().Nil(p)
				s.Require().Contains(err.Error(), tc.errContains)
			}
		})
	}
}

// TestRun tests the precompile's Run method.
func (s *PrecompileTestSuite) TestRun() {
	contract := vm.NewPrecompile(
		s.keyring.GetAddr(0),
		s.precompile.Address(),
		uint256.NewInt(0),
		uint64(1000000),
	)

	testCases := []struct {
		name        string
		malleate    func() *vm.Contract
		postCheck   func(data []byte)
		expPass     bool
		errContains string
	}{
		{
			"fail - invalid method",
			func() *vm.Contract {
				contract.Input = []byte("invalid")
				return contract
			},
			func([]byte) {},
			false,
			"no method with id",
		},
		{
			"fail - error during unpack",
			func() *vm.Contract {
				// only pass the method ID to the input
				contract.Input = s.precompile.Methods[bech32.HexToBech32Method].ID
				return contract
			},
			func([]byte) {},
			false,
			"abi: attempting to unmarshal an empty string while arguments are expected",
		},
		{
			"fail - HexToBech32 method error",
			func() *vm.Contract {
				input, err := s.precompile.Pack(
					bech32.HexToBech32Method,
					s.keyring.GetAddr(0),
					"",
				)
				s.Require().NoError(err, "failed to pack input")

				// only pass the method ID to the input
				contract.Input = input
				return contract
			},
			func([]byte) {},
			false,
			"invalid bech32 human readable prefix (HRP)",
		},
		{
			"pass - hex to bech32 account (cosmos)",
			func() *vm.Contract {
				input, err := s.precompile.Pack(
					bech32.HexToBech32Method,
					s.keyring.GetAddr(0),
					config.Bech32Prefix,
				)
				s.Require().NoError(err, "failed to pack input")
				contract.Input = input
				return contract
			},
			func(data []byte) {
				args, err := s.precompile.Unpack(bech32.HexToBech32Method, data)
				s.Require().NoError(err, "failed to unpack output")
				s.Require().Len(args, 1)
				addr, ok := args[0].(string)
				s.Require().True(ok)
				s.Require().Equal(s.keyring.GetAccAddr(0).String(), addr)
			},
			true,
			"",
		},
		{
			"pass - hex to bech32 validator operator (cosmosvaloper)",
			func() *vm.Contract {
				valAddrCodec := s.network.App.GetStakingKeeper().ValidatorAddressCodec()
				valAddrBz, err := valAddrCodec.StringToBytes(s.network.GetValidators()[0].GetOperator())
				s.Require().NoError(err, "failed to convert string to bytes")
				input, err := s.precompile.Pack(
					bech32.HexToBech32Method,
					common.BytesToAddress(valAddrBz),
					config.Bech32PrefixValAddr,
				)
				s.Require().NoError(err, "failed to pack input")
				contract.Input = input
				return contract
			},
			func(data []byte) {
				args, err := s.precompile.Unpack(bech32.HexToBech32Method, data)
				s.Require().NoError(err, "failed to unpack output")
				s.Require().Len(args, 1)
				addr, ok := args[0].(string)
				s.Require().True(ok)
				s.Require().Equal(s.network.GetValidators()[0].OperatorAddress, addr)
			},
			true,
			"",
		},
		{
			"pass - hex to bech32 consensus address (cosmosvalcons)",
			func() *vm.Contract {
				input, err := s.precompile.Pack(
					bech32.HexToBech32Method,
					s.keyring.GetAddr(0),
					config.Bech32PrefixConsAddr,
				)
				s.Require().NoError(err, "failed to pack input")
				contract.Input = input
				return contract
			},
			func(data []byte) {
				args, err := s.precompile.Unpack(bech32.HexToBech32Method, data)
				s.Require().NoError(err, "failed to unpack output")
				s.Require().Len(args, 1)
				addr, ok := args[0].(string)
				s.Require().True(ok)
				s.Require().Equal(sdk.ConsAddress(s.keyring.GetAddr(0).Bytes()).String(), addr)
			},
			true,
			"",
		},
		{
			"pass - bech32 to hex account address",
			func() *vm.Contract {
				input, err := s.precompile.Pack(
					bech32.Bech32ToHexMethod,
					s.keyring.GetAccAddr(0).String(),
				)
				s.Require().NoError(err, "failed to pack input")
				contract.Input = input
				return contract
			},
			func(data []byte) {
				args, err := s.precompile.Unpack(bech32.Bech32ToHexMethod, data)
				s.Require().NoError(err, "failed to unpack output")
				s.Require().Len(args, 1)
				addr, ok := args[0].(common.Address)
				s.Require().True(ok)
				s.Require().Equal(s.keyring.GetAddr(0), addr)
			},
			true,
			"",
		},
		{
			"pass - bech32 to hex validator address",
			func() *vm.Contract {
				input, err := s.precompile.Pack(
					bech32.Bech32ToHexMethod,
					s.network.GetValidators()[0].OperatorAddress,
				)
				s.Require().NoError(err, "failed to pack input")
				contract.Input = input
				return contract
			},
			func(data []byte) {
				valAddrCodec := s.network.App.GetStakingKeeper().ValidatorAddressCodec()
				valAddrBz, err := valAddrCodec.StringToBytes(s.network.GetValidators()[0].GetOperator())
				s.Require().NoError(err, "failed to convert string to bytes")

				args, err := s.precompile.Unpack(bech32.Bech32ToHexMethod, data)
				s.Require().NoError(err, "failed to unpack output")
				s.Require().Len(args, 1)
				addr, ok := args[0].(common.Address)
				s.Require().True(ok)
				s.Require().Equal(common.BytesToAddress(valAddrBz), addr)
			},
			true,
			"",
		},
		{
			"pass - bech32 to hex consensus address",
			func() *vm.Contract {
				input, err := s.precompile.Pack(
					bech32.Bech32ToHexMethod,
					sdk.ConsAddress(s.keyring.GetAddr(0).Bytes()).String(),
				)
				s.Require().NoError(err, "failed to pack input")
				contract.Input = input
				return contract
			},
			func(data []byte) {
				args, err := s.precompile.Unpack(bech32.Bech32ToHexMethod, data)
				s.Require().NoError(err, "failed to unpack output")
				s.Require().Len(args, 1)
				addr, ok := args[0].(common.Address)
				s.Require().True(ok)
				s.Require().Equal(s.keyring.GetAddr(0), addr)
			},
			true,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// setup basic test suite
			s.SetupTest()

			// malleate testcase
			contract := tc.malleate()

			// Run precompiled contract

			// NOTE: we can ignore the EVM and readonly args since it's a stateless
			// precompiled contract
			bz, err := s.precompile.Run(nil, contract, true)

			// Check results
			if tc.expPass {
				s.Require().NoError(err, "expected no error when running the precompile")
				s.Require().NotNil(bz, "expected returned bytes not to be nil")
				tc.postCheck(bz)
			} else {
				s.Require().Error(err, "expected error to be returned when running the precompile")
				s.Require().Nil(bz, "expected returned bytes to be nil")
				s.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}
