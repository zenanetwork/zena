package keeper

import (
	"github.com/zenanetwork/zena/x/feemarket/types"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper grants access to the Fee Market module state.
type Keeper struct {
	// Protobuf codec
	cdc codec.BinaryCodec
	// Store key required for the Fee Market Prefix KVStore.
	storeKey     storetypes.StoreKey
	transientKey storetypes.StoreKey
	// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
	authority sdk.AccAddress
}

// NewKeeper generates new fee market module keeper
func NewKeeper(
	cdc codec.BinaryCodec, authority sdk.AccAddress, storeKey, transientKey storetypes.StoreKey,
) Keeper {
	// ensure authority account is correctly formatted
	if err := sdk.VerifyAddressFormat(authority); err != nil {
		panic(err)
	}

	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		authority:    authority,
		transientKey: transientKey,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// ----------------------------------------------------------------------------
// Parent Block Gas Used
// Required by EIP1559 base fee calculation.
// ----------------------------------------------------------------------------

// SetBlockGasWanted sets the block gas wanted to the store.
// CONTRACT: this should be only called during EndBlock.
func (k Keeper) SetBlockGasWanted(ctx sdk.Context, gas uint64) {
	store := ctx.KVStore(k.storeKey)
	gasBz := sdk.Uint64ToBigEndian(gas)
	store.Set(types.KeyPrefixBlockGasWanted, gasBz)
}

// GetBlockGasWanted returns the last block gas wanted value from the store.
func (k Keeper) GetBlockGasWanted(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	return sdk.BigEndianToUint64(store.Get(types.KeyPrefixBlockGasWanted))
}

// GetTransientGasWanted returns the gas wanted in the current block from transient store.
func (k Keeper) GetTransientGasWanted(ctx sdk.Context) uint64 {
	store := ctx.TransientStore(k.transientKey)
	return sdk.BigEndianToUint64(store.Get(types.KeyPrefixTransientBlockGasWanted))
}

// SetTransientBlockGasWanted sets the block gas wanted to the transient store.
func (k Keeper) SetTransientBlockGasWanted(ctx sdk.Context, gasWanted uint64) {
	store := ctx.TransientStore(k.transientKey)
	gasBz := sdk.Uint64ToBigEndian(gasWanted)
	store.Set(types.KeyPrefixTransientBlockGasWanted, gasBz)
}

// AddTransientGasWanted adds the cumulative gas wanted in the transient store
func (k Keeper) AddTransientGasWanted(ctx sdk.Context, gasWanted uint64) (uint64, error) {
	result := k.GetTransientGasWanted(ctx) + gasWanted
	k.SetTransientBlockGasWanted(ctx, result)
	return result, nil
}
