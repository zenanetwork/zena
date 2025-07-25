package ante

import (
	"fmt"
	"math/big"

	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/zenanetwork/zena/ante/evm"
	testconstants "github.com/zenanetwork/zena/testutil/constants"
	"github.com/zenanetwork/zena/testutil/integration/evm/factory"
	"github.com/zenanetwork/zena/testutil/integration/evm/grpc"
	"github.com/zenanetwork/zena/testutil/integration/evm/network"
	testkeyring "github.com/zenanetwork/zena/testutil/keyring"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"

	"cosmossdk.io/math"

	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

func (s *EvmUnitAnteTestSuite) TestCanTransfer() {
	keyring := testkeyring.New(1)
	unitNetwork := network.NewUnitTestNetwork(
		s.create,
		network.WithChainID(testconstants.ChainID{
			ChainID:    s.ChainID,
			EVMChainID: s.EvmChainID,
		}),
		network.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
	)
	grpcHandler := grpc.NewIntegrationHandler(unitNetwork)
	txFactory := factory.New(unitNetwork, grpcHandler)
	senderKey := keyring.GetKey(0)

	testCases := []struct {
		name          string
		expectedError error
		isLondon      bool
		malleate      func(txArgs *evmtypes.EvmTxArgs)
	}{
		{
			name:          "fail: isLondon and insufficient fee",
			expectedError: errortypes.ErrInsufficientFee,
			isLondon:      true,
			malleate: func(txArgs *evmtypes.EvmTxArgs) {
				txArgs.GasFeeCap = big.NewInt(0)
			},
		},
		{
			name:          "fail: invalid tx with insufficient balance",
			expectedError: errortypes.ErrInsufficientFunds,
			isLondon:      true,
			malleate: func(txArgs *evmtypes.EvmTxArgs) {
				balanceResp, err := grpcHandler.GetBalanceFromEVM(senderKey.AccAddr)
				s.Require().NoError(err)

				balance, ok := math.NewIntFromString(balanceResp.Balance)
				s.Require().True(ok)
				invalidAmount := balance.Add(math.NewInt(1)).BigInt()
				txArgs.Amount = invalidAmount
			},
		},
		{
			name:          "success: valid tx and sufficient balance",
			expectedError: nil,
			isLondon:      true,
			malleate: func(*evmtypes.EvmTxArgs) {
			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("%v_%v_%v", evmtypes.GetTxTypeName(s.EthTxType), s.ChainID, tc.name), func() {
			baseFeeResp, err := grpcHandler.GetEvmBaseFee()
			s.Require().NoError(err)
			ethCfg := unitNetwork.GetEVMChainConfig()
			evmParams, err := grpcHandler.GetEvmParams()
			s.Require().NoError(err)
			ctx := unitNetwork.GetContext()
			signer := gethtypes.MakeSigner(ethCfg, big.NewInt(ctx.BlockHeight()), uint64(ctx.BlockTime().Unix())) //#nosec G115 -- int overflow is not a concern here
			txArgs, err := txFactory.GenerateDefaultTxTypeArgs(senderKey.Addr, s.EthTxType)
			s.Require().NoError(err)
			txArgs.Amount = big.NewInt(100)

			tc.malleate(&txArgs)

			msg := evmtypes.NewTx(&txArgs)
			msg.From = senderKey.Addr.String()
			signMsg, err := txFactory.SignMsgEthereumTx(senderKey.Priv, *msg)
			s.Require().NoError(err)
			coreMsg, err := signMsg.AsMessage(signer, baseFeeResp.BaseFee.BigInt())
			s.Require().NoError(err)

			// Function under test
			err = evm.CanTransfer(
				unitNetwork.GetContext(),
				unitNetwork.App.GetEVMKeeper(),
				*coreMsg,
				baseFeeResp.BaseFee.BigInt(),
				evmParams.Params,
				tc.isLondon,
			)

			if tc.expectedError != nil {
				s.Require().Error(err)
				s.Contains(err.Error(), tc.expectedError.Error())

			} else {
				s.Require().NoError(err)
			}
		})
	}
}
