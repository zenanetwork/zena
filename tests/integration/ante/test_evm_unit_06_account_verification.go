package ante

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/zenanetwork/zena/ante/evm"
	testconstants "github.com/zenanetwork/zena/testutil/constants"
	"github.com/zenanetwork/zena/testutil/integration/evm/factory"
	"github.com/zenanetwork/zena/testutil/integration/evm/grpc"
	"github.com/zenanetwork/zena/testutil/integration/evm/network"
	testkeyring "github.com/zenanetwork/zena/testutil/keyring"
	"github.com/zenanetwork/zena/x/vm/statedb"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"

	"cosmossdk.io/math"

	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

func (s *EvmUnitAnteTestSuite) TestVerifyAccountBalance() {
	// Setup
	keyring := testkeyring.New(2)
	unitNetwork := network.NewUnitTestNetwork(
		s.create,
		network.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
		network.WithChainID(testconstants.ChainID{
			ChainID:    s.ChainID,
			EVMChainID: s.EvmChainID,
		}),
	)
	grpcHandler := grpc.NewIntegrationHandler(unitNetwork)
	txFactory := factory.New(unitNetwork, grpcHandler)
	senderKey := keyring.GetKey(1)

	testCases := []struct {
		name                   string
		expectedError          error
		generateAccountAndArgs func() (*statedb.Account, evmtypes.EvmTxArgs)
	}{
		{
			name:          "fail: sender is not EOA",
			expectedError: errortypes.ErrInvalidType,
			generateAccountAndArgs: func() (*statedb.Account, evmtypes.EvmTxArgs) {
				statedbAccount := getDefaultStateDBAccount(unitNetwork, senderKey.Addr)
				txArgs, err := txFactory.GenerateDefaultTxTypeArgs(senderKey.Addr, s.EthTxType)
				s.Require().NoError(err)

				statedbAccount.CodeHash = []byte("test")
				s.Require().NoError(err)
				return statedbAccount, txArgs
			},
		},
		{
			name:          "fail: sender balance is lower than the transaction cost",
			expectedError: errortypes.ErrInsufficientFunds,
			generateAccountAndArgs: func() (*statedb.Account, evmtypes.EvmTxArgs) {
				statedbAccount := getDefaultStateDBAccount(unitNetwork, senderKey.Addr)
				txArgs, err := txFactory.GenerateDefaultTxTypeArgs(senderKey.Addr, s.EthTxType)
				s.Require().NoError(err)

				// Make tx cost greater than balance
				balanceResp, err := grpcHandler.GetBalanceFromEVM(senderKey.AccAddr)
				s.Require().NoError(err)

				balance, ok := math.NewIntFromString(balanceResp.Balance)
				s.Require().True(ok)
				invalidAmount := balance.Add(math.NewInt(100))
				txArgs.Amount = invalidAmount.BigInt()
				return statedbAccount, txArgs
			},
		},
		{
			name:          "fail: tx cost is negative",
			expectedError: errortypes.ErrInvalidCoins,
			generateAccountAndArgs: func() (*statedb.Account, evmtypes.EvmTxArgs) {
				statedbAccount := getDefaultStateDBAccount(unitNetwork, senderKey.Addr)
				txArgs, err := txFactory.GenerateDefaultTxTypeArgs(senderKey.Addr, s.EthTxType)
				s.Require().NoError(err)

				// Make tx cost negative. This has to be a big value because
				// it has to be bigger than the fee for the full cost to be negative
				invalidAmount := big.NewInt(-1e18)
				txArgs.Amount = invalidAmount
				return statedbAccount, txArgs
			},
		},
		{
			name:          "success: tx is successful and account is created if its nil",
			expectedError: errortypes.ErrInsufficientFunds,
			generateAccountAndArgs: func() (*statedb.Account, evmtypes.EvmTxArgs) {
				txArgs, err := txFactory.GenerateDefaultTxTypeArgs(senderKey.Addr, s.EthTxType)
				s.Require().NoError(err)
				return nil, txArgs
			},
		},
		{
			name:          "success: tx is successful if account is EOA and exists",
			expectedError: nil,
			generateAccountAndArgs: func() (*statedb.Account, evmtypes.EvmTxArgs) {
				statedbAccount := getDefaultStateDBAccount(unitNetwork, senderKey.Addr)
				txArgs, err := txFactory.GenerateDefaultTxTypeArgs(senderKey.Addr, s.EthTxType)
				s.Require().NoError(err)
				return statedbAccount, txArgs
			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("%v_%v_%v", evmtypes.GetTxTypeName(s.EthTxType), s.ChainID, tc.name), func() {
			// Perform test logic
			statedbAccount, txArgs := tc.generateAccountAndArgs()
			txData, err := txArgs.ToTxData()
			s.Require().NoError(err)

			//  Function to be tested
			err = evm.VerifyAccountBalance(
				unitNetwork.GetContext(),
				unitNetwork.App.GetAccountKeeper(),
				statedbAccount,
				senderKey.Addr,
				txData,
			)

			if tc.expectedError != nil {
				s.Require().Error(err)
				s.Contains(err.Error(), tc.expectedError.Error())
			} else {
				s.Require().NoError(err)
			}
			// Make sure the account is created either wa
			acc, err := grpcHandler.GetAccount(senderKey.AccAddr.String())
			s.Require().NoError(err)
			s.Require().NotEmpty(acc)

			// Clean block for next test
			err = unitNetwork.NextBlock()
			s.Require().NoError(err)
		})
	}
}

func getDefaultStateDBAccount(unitNetwork *network.UnitTestNetwork, addr common.Address) *statedb.Account {
	statedb := unitNetwork.GetStateDB()
	return statedb.Keeper().GetAccount(unitNetwork.GetContext(), addr)
}
