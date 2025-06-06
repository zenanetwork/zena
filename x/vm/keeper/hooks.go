package keeper

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/zenanetwork/zena/x/vm/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Event Hooks
// These can be utilized to customize evm transaction processing.

var _ types.EvmHooks = MultiEvmHooks{}

// MultiEvmHooks combine multiple evm hooks, all hook functions are run in array sequence
type MultiEvmHooks []types.EvmHooks

// NewMultiEvmHooks combine multiple evm hooks
func NewMultiEvmHooks(hooks ...types.EvmHooks) MultiEvmHooks {
	return hooks
}

// PostTxProcessing delegate the call to underlying hooks
func (mh MultiEvmHooks) PostTxProcessing(ctx sdk.Context, sender common.Address, msg core.Message, receipt *ethtypes.Receipt) error {
	for i := range mh {
		if err := mh[i].PostTxProcessing(ctx, sender, msg, receipt); err != nil {
			return errorsmod.Wrapf(err, "EVM hook %T failed", mh[i])
		}
	}
	return nil
}
