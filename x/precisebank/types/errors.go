package types

import errorsmod "cosmossdk.io/errors"

var (
	// ErrInsufficientReserve is returned when the reserve module account
	// has insufficient balance for a carry operation.
	ErrInsufficientReserve = errorsmod.Register(ModuleName, 2, "insufficient reserve balance")
)
