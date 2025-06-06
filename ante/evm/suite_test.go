package evm_test

import (
	"testing"

	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"

	testconstants "github.com/zenanetwork/zena/testutil/constants"
)

// EvmAnteTestSuite aims to test all EVM ante handler unit functions.
// NOTE: the suite only holds properties related to global execution parameters
// (what type of tx to run the tests with) not independent tests values.
type EvmAnteTestSuite struct {
	suite.Suite

	ethTxType  int
	chainID    string
	evmChainID uint64
}

func TestEvmAnteTestSuite(t *testing.T) {
	txTypes := []int{gethtypes.DynamicFeeTxType, gethtypes.LegacyTxType, gethtypes.AccessListTxType}
	chainIDs := []testconstants.ChainID{testconstants.ExampleChainID, testconstants.SixDecimalsChainID}
	for _, txType := range txTypes {
		for _, chainID := range chainIDs {
			suite.Run(t, &EvmAnteTestSuite{
				ethTxType:  txType,
				chainID:    chainID.ChainID,
				evmChainID: chainID.EVMChainID,
			})
		}
	}
}
