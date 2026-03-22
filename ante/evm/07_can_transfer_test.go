package evm_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	anteinterfaces "github.com/zenanetwork/zena/ante/interfaces"
	evm "github.com/zenanetwork/zena/ante/evm"
	"github.com/zenanetwork/zena/x/vm/statedb"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
)

// mockEVMKeeper embeds the EVMKeeper interface and overrides GetAccount
// for unit testing CanTransfer. Only GetAccount is called by CanTransfer.
type mockEVMKeeper struct {
	anteinterfaces.EVMKeeper
	account *statedb.Account
}

func (m *mockEVMKeeper) GetAccount(_ sdk.Context, _ common.Address) *statedb.Account {
	return m.account
}

func newTestContext() sdk.Context {
	key := storetypes.NewKVStoreKey("test")
	tkey := storetypes.NewTransientStoreKey("test_transient")
	return testutil.DefaultContext(key, tkey)
}

func TestCanTransfer_NilAccount(t *testing.T) {
	ctx := newTestContext()
	keeper := &mockEVMKeeper{account: nil}

	msg := core.Message{
		From:      common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
		Value:     big.NewInt(100),
		GasFeeCap: big.NewInt(1000),
	}

	err := evm.CanTransfer(ctx, keeper, msg, big.NewInt(0), evmtypes.Params{}, false)
	require.Error(t, err)
	require.ErrorIs(t, err, errortypes.ErrInsufficientFunds)
}

func TestCanTransfer_NilBalance(t *testing.T) {
	ctx := newTestContext()
	keeper := &mockEVMKeeper{
		account: &statedb.Account{
			Nonce:   0,
			Balance: nil, // nil balance (SpendableCoin conversion failure)
		},
	}

	msg := core.Message{
		From:      common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
		Value:     big.NewInt(100),
		GasFeeCap: big.NewInt(1000),
	}

	err := evm.CanTransfer(ctx, keeper, msg, big.NewInt(0), evmtypes.Params{}, false)
	require.Error(t, err)
	require.ErrorIs(t, err, errortypes.ErrInsufficientFunds)
}

func TestCanTransfer_InsufficientBalance(t *testing.T) {
	ctx := newTestContext()
	keeper := &mockEVMKeeper{
		account: &statedb.Account{
			Nonce:   0,
			Balance: uint256.NewInt(50), // less than transfer value
		},
	}

	msg := core.Message{
		From:      common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
		Value:     big.NewInt(100),
		GasFeeCap: big.NewInt(1000),
	}

	err := evm.CanTransfer(ctx, keeper, msg, big.NewInt(0), evmtypes.Params{}, false)
	require.Error(t, err)
	require.ErrorIs(t, err, errortypes.ErrInsufficientFunds)
}

func TestCanTransfer_SufficientBalance(t *testing.T) {
	ctx := newTestContext()
	keeper := &mockEVMKeeper{
		account: &statedb.Account{
			Nonce:   0,
			Balance: uint256.NewInt(200),
		},
	}

	msg := core.Message{
		From:      common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
		Value:     big.NewInt(100),
		GasFeeCap: big.NewInt(1000),
	}

	err := evm.CanTransfer(ctx, keeper, msg, big.NewInt(0), evmtypes.Params{}, false)
	require.NoError(t, err)
}

func TestCanTransfer_ZeroValue(t *testing.T) {
	ctx := newTestContext()
	// Even nil account should pass when value is 0
	keeper := &mockEVMKeeper{account: nil}

	msg := core.Message{
		From:      common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
		Value:     big.NewInt(0),
		GasFeeCap: big.NewInt(1000),
	}

	err := evm.CanTransfer(ctx, keeper, msg, big.NewInt(0), evmtypes.Params{}, false)
	require.NoError(t, err)
}
