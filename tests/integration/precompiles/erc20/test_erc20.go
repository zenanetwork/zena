package erc20

import (
	"math/big"

	"github.com/zenanetwork/zena/precompiles/erc20"
)

func (s *PrecompileTestSuite) TestIsTransaction() {
	s.SetupTest()

	// Queries
	method := s.precompile.Methods[erc20.BalanceOfMethod]
	s.Require().False(s.precompile.IsTransaction(&method))
	method = s.precompile.Methods[erc20.DecimalsMethod]
	s.Require().False(s.precompile.IsTransaction(&method))
	method = s.precompile.Methods[erc20.NameMethod]
	s.Require().False(s.precompile.IsTransaction(&method))
	method = s.precompile.Methods[erc20.SymbolMethod]
	s.Require().False(s.precompile.IsTransaction(&method))
	method = s.precompile.Methods[erc20.TotalSupplyMethod]
	s.Require().False(s.precompile.IsTransaction(&method))

	// Transactions
	method = s.precompile.Methods[erc20.ApproveMethod]
	s.Require().True(s.precompile.IsTransaction(&method))
	method = s.precompile.Methods[erc20.TransferMethod]
	s.Require().True(s.precompile.IsTransaction(&method))
	method = s.precompile.Methods[erc20.TransferFromMethod]
	s.Require().True(s.precompile.IsTransaction(&method))
}

func (s *PrecompileTestSuite) TestRequiredGas() {
	s.SetupTest()

	testcases := []struct {
		name     string
		malleate func() []byte
		expGas   uint64
	}{
		{
			name: erc20.BalanceOfMethod,
			malleate: func() []byte {
				bz, err := s.precompile.Pack(erc20.BalanceOfMethod, s.keyring.GetAddr(0))
				s.Require().NoError(err, "expected no error packing ABI")
				return bz
			},
			expGas: erc20.GasBalanceOf,
		},
		{
			name: erc20.DecimalsMethod,
			malleate: func() []byte {
				bz, err := s.precompile.Pack(erc20.DecimalsMethod)
				s.Require().NoError(err, "expected no error packing ABI")
				return bz
			},
			expGas: erc20.GasDecimals,
		},
		{
			name: erc20.NameMethod,
			malleate: func() []byte {
				bz, err := s.precompile.Pack(erc20.NameMethod)
				s.Require().NoError(err, "expected no error packing ABI")
				return bz
			},
			expGas: erc20.GasName,
		},
		{
			name: erc20.SymbolMethod,
			malleate: func() []byte {
				bz, err := s.precompile.Pack(erc20.SymbolMethod)
				s.Require().NoError(err, "expected no error packing ABI")
				return bz
			},
			expGas: erc20.GasSymbol,
		},
		{
			name: erc20.TotalSupplyMethod,
			malleate: func() []byte {
				bz, err := s.precompile.Pack(erc20.TotalSupplyMethod)
				s.Require().NoError(err, "expected no error packing ABI")
				return bz
			},
			expGas: erc20.GasTotalSupply,
		},
		{
			name: erc20.ApproveMethod,
			malleate: func() []byte {
				bz, err := s.precompile.Pack(erc20.ApproveMethod, s.keyring.GetAddr(0), big.NewInt(1))
				s.Require().NoError(err, "expected no error packing ABI")
				return bz
			},
			expGas: erc20.GasApprove,
		},
		{
			name: erc20.TransferMethod,
			malleate: func() []byte {
				bz, err := s.precompile.Pack(erc20.TransferMethod, s.keyring.GetAddr(0), big.NewInt(1))
				s.Require().NoError(err, "expected no error packing ABI")
				return bz
			},
			expGas: erc20.GasTransfer,
		},
		{
			name: erc20.TransferFromMethod,
			malleate: func() []byte {
				bz, err := s.precompile.Pack(erc20.TransferFromMethod, s.keyring.GetAddr(0), s.keyring.GetAddr(0), big.NewInt(1))
				s.Require().NoError(err, "expected no error packing ABI")
				return bz
			},
			expGas: erc20.GasTransferFrom,
		},
		{
			name: erc20.AllowanceMethod,
			malleate: func() []byte {
				bz, err := s.precompile.Pack(erc20.AllowanceMethod, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NoError(err, "expected no error packing ABI")
				return bz
			},
			expGas: erc20.GasAllowance,
		},
		{
			name: "invalid method",
			malleate: func() []byte {
				return []byte("invalid method")
			},
			expGas: 0,
		},
		{
			name: "input bytes too short",
			malleate: func() []byte {
				return []byte{0x00, 0x00, 0x00}
			},
			expGas: 0,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			input := tc.malleate()

			s.Require().Equal(tc.expGas, s.precompile.RequiredGas(input))
		})
	}
}
