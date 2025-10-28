package ante

import (
	"math/big"

	evmante "github.com/zenanetwork/zena/ante/evm"
	"github.com/zenanetwork/zena/testutil"
	testutiltx "github.com/zenanetwork/zena/testutil/tx"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *EvmAnteTestSuite) TestEthSetupContextDecorator() {
	dec := evmante.NewEthSetUpContextDecorator(s.GetNetwork().App.GetEVMKeeper())
	ethContractCreationTxParams := &evmtypes.EvmTxArgs{
		ChainID:  evmtypes.GetEthChainConfig().ChainID,
		Nonce:    1,
		Amount:   big.NewInt(10),
		GasLimit: 1000,
		GasPrice: big.NewInt(1),
	}
	tx := evmtypes.NewTx(ethContractCreationTxParams)

	testCases := []struct {
		name    string
		tx      sdk.Tx
		expPass bool
	}{
		{"invalid transaction type - does not implement GasTx", &testutiltx.InvalidTx{}, false},
		{
			"success - transaction implement GasTx",
			tx,
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, err := dec.AnteHandle(s.GetNetwork().GetContext(), tc.tx, false, testutil.NoOpNextFn)

			if tc.expPass {
				s.Require().NoError(err)
				s.Equal(storetypes.GasConfig{}, ctx.KVGasConfig())
				s.Equal(storetypes.GasConfig{}, ctx.TransientKVGasConfig())
			} else {
				s.Require().Error(err)
			}
		})
	}
}
