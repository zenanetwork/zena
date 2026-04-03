package keeper

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zenanetwork/zena/x/vm/types"
)

var (
	testTokenAddr = common.HexToAddress("0x1234567890abcdef")
	testFromAddr  = common.HexToAddress("0xaaaaaaaaaaaaaaaa")
	testToAddr    = common.HexToAddress("0xbbbbbbbbbbbbbbbb")
	testAmount    = big.NewInt(1000)
)

// makeTopicHex returns the hex representation of an address zero-padded to 32 bytes,
// matching how ERC-20 Transfer event indexed parameters are encoded.
func makeTopicHex(addr common.Address) string {
	return common.BytesToHash(common.LeftPadBytes(addr.Bytes(), 32)).Hex()
}

// makeAmountData returns the ABI-encoded uint256 representation of amount (32 bytes).
func makeAmountData(amount *big.Int) []byte {
	data := make([]byte, 32)
	b := amount.Bytes()
	copy(data[32-len(b):], b)
	return data
}

func TestValidateApprovalEventDoesNotExist(t *testing.T) {
	tests := []struct {
		name        string
		res         *types.MsgEthereumTxResponse
		expectError bool
	}{
		{
			name: "empty logs",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{},
			},
			expectError: false,
		},
		{
			name: "no approval event",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{
					{
						Topics: []string{"0x1234567890abcdef"},
					},
				},
			},
			expectError: false,
		},
		{
			name: "has approval event",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{
					{
						Topics: []string{logApprovalSigHash.Hex()},
					},
				},
			},
			expectError: true,
		},
		{
			name: "approval event among others",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{
					{
						Topics: []string{"0x1234567890abcdef"},
					},
					{
						Topics: []string{logApprovalSigHash.Hex()},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateApprovalEventDoesNotExist(tt.res.Logs)
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), "unexpected Approval event")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateTransferEventExists(t *testing.T) {
	validTopics := []string{
		logTransferSigHash.Hex(),
		makeTopicHex(testFromAddr),
		makeTopicHex(testToAddr),
	}
	validData := makeAmountData(testAmount)

	tests := []struct {
		name           string
		res            *types.MsgEthereumTxResponse
		tokenAddress   common.Address
		expectedFrom   common.Address
		expectedTo     common.Address
		expectedAmount *big.Int
		expectError    string
	}{
		{
			name: "empty logs",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{},
			},
			expectedFrom:   testFromAddr,
			expectedTo:     testToAddr,
			expectedAmount: testAmount,
			expectError:    "expected Transfer event",
		},
		{
			name: "no transfer event",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{
					{
						Topics: []string{"0x1234567890abcdef"},
					},
				},
			},
			tokenAddress:   testTokenAddr,
			expectedFrom:   testFromAddr,
			expectedTo:     testToAddr,
			expectedAmount: testAmount,
			expectError:    "expected Transfer event",
		},
		{
			name: "transfer event from different contract address",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{
					{
						Address: common.HexToAddress("0x1234567890abcdef").Hex(),
						Topics:  validTopics,
						Data:    validData,
					},
				},
			},
			tokenAddress:   common.HexToAddress("fedcba0987654321"),
			expectedFrom:   testFromAddr,
			expectedTo:     testToAddr,
			expectedAmount: testAmount,
			expectError:    "Transfer event from unexpected address",
		},
		{
			name: "duplicate transfer event",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{
					{
						Address: testTokenAddr.Hex(),
						Topics:  validTopics,
						Data:    validData,
					},
					{
						Address: testTokenAddr.Hex(),
						Topics:  validTopics,
						Data:    validData,
					},
				},
			},
			tokenAddress:   testTokenAddr,
			expectedFrom:   testFromAddr,
			expectedTo:     testToAddr,
			expectedAmount: testAmount,
			expectError:    "duplicate Transfer event",
		},
		{
			name: "topics array too short",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{
					{
						Address: testTokenAddr.Hex(),
						Topics:  []string{logTransferSigHash.Hex()},
						Data:    validData,
					},
				},
			},
			tokenAddress:   testTokenAddr,
			expectedFrom:   testFromAddr,
			expectedTo:     testToAddr,
			expectedAmount: testAmount,
			expectError:    "Transfer event has insufficient topics",
		},
		{
			name: "wrong from address in topics",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{
					{
						Address: testTokenAddr.Hex(),
						Topics: []string{
							logTransferSigHash.Hex(),
							makeTopicHex(common.HexToAddress("0xcccccccccccccccc")),
							makeTopicHex(testToAddr),
						},
						Data: validData,
					},
				},
			},
			tokenAddress:   testTokenAddr,
			expectedFrom:   testFromAddr,
			expectedTo:     testToAddr,
			expectedAmount: testAmount,
			expectError:    "Transfer from mismatch",
		},
		{
			name: "wrong to address in topics",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{
					{
						Address: testTokenAddr.Hex(),
						Topics: []string{
							logTransferSigHash.Hex(),
							makeTopicHex(testFromAddr),
							makeTopicHex(common.HexToAddress("0xdddddddddddddddd")),
						},
						Data: validData,
					},
				},
			},
			tokenAddress:   testTokenAddr,
			expectedFrom:   testFromAddr,
			expectedTo:     testToAddr,
			expectedAmount: testAmount,
			expectError:    "Transfer to mismatch",
		},
		{
			name: "wrong amount in data",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{
					{
						Address: testTokenAddr.Hex(),
						Topics:  validTopics,
						Data:    makeAmountData(big.NewInt(9999)),
					},
				},
			},
			tokenAddress:   testTokenAddr,
			expectedFrom:   testFromAddr,
			expectedTo:     testToAddr,
			expectedAmount: testAmount,
			expectError:    "Transfer value mismatch",
		},
		{
			name: "valid transfer event with all fields matching",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{
					{
						Address: testTokenAddr.Hex(),
						Topics:  validTopics,
						Data:    validData,
					},
				},
			},
			tokenAddress:   testTokenAddr,
			expectedFrom:   testFromAddr,
			expectedTo:     testToAddr,
			expectedAmount: testAmount,
			expectError:    "",
		},
		{
			name: "valid transfer event among other logs",
			res: &types.MsgEthereumTxResponse{
				Logs: []*types.Log{
					{
						Address: testTokenAddr.Hex(),
						Topics:  []string{"0x1234567890abcdef"},
					},
					{
						Address: testTokenAddr.Hex(),
						Topics:  validTopics,
						Data:    validData,
					},
				},
			},
			tokenAddress:   testTokenAddr,
			expectedFrom:   testFromAddr,
			expectedTo:     testToAddr,
			expectedAmount: testAmount,
			expectError:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTransferEventExists(tt.res.Logs, tt.tokenAddress, tt.expectedFrom, tt.expectedTo, tt.expectedAmount)
			if tt.expectError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
