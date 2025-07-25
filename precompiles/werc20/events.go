package werc20

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	cmn "github.com/zenanetwork/zena/precompiles/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// EventTypeDeposit is the key of the event type for the Deposit transaction.
	EventTypeDeposit = "Deposit"
	// EventTypeWithdrawal is the key of the event type for the Withdraw transaction.
	EventTypeWithdrawal = "Withdrawal"
)

// EmitDepositEvent creates a new Deposit event emitted after a Deposit transaction.
func (p Precompile) EmitDepositEvent(
	ctx sdk.Context,
	stateDB vm.StateDB,
	caller common.Address,
	amount *big.Int,
) error {
	event := p.Events[EventTypeDeposit]
	return p.createWERC20Event(ctx, stateDB, event, caller, amount)
}

// EmitWithdrawalEvent creates a new Withdrawal event emitted after a Withdraw transaction.
func (p Precompile) EmitWithdrawalEvent(
	ctx sdk.Context,
	stateDB vm.StateDB,
	src common.Address,
	amount *big.Int,
) error {
	event := p.Events[EventTypeWithdrawal]
	return p.createWERC20Event(ctx, stateDB, event, src, amount)
}

// createWERC20Event adds to the StateDB a log representing an event for the
// WERC20 precompile.
func (p Precompile) createWERC20Event(
	ctx sdk.Context,
	stateDB vm.StateDB,
	event abi.Event,
	address common.Address,
	amount *big.Int,
) error {
	// Prepare the event topics
	topics := make([]common.Hash, 2)

	topics[0] = event.ID

	var err error
	topics[1], err = cmn.MakeTopic(address)
	if err != nil {
		return err
	}

	arguments := abi.Arguments{event.Inputs[1]}
	packed, err := arguments.Pack(amount)
	if err != nil {
		return err
	}

	stateDB.AddLog(&ethtypes.Log{
		Address:     p.Address(),
		Topics:      topics,
		Data:        packed,
		BlockNumber: uint64(ctx.BlockHeight()), //nolint:gosec // G115
	})

	return nil
}
