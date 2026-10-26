package config

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// NativeDenomMetadata returns the canonical bank metadata of the Zenanet native
// token: base aznnt (exponent 0), display znnt (exponent 18), symbol ZNNT.
//
// It matches the metadata written into the mainnet genesis by local_node.sh.
func NativeDenomMetadata() banktypes.Metadata {
	return banktypes.Metadata{
		Description: "The native staking token for zenad.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    ZenaChainDenom,
				Exponent: 0,
				Aliases:  []string{"attoznnt"},
			},
			{
				Denom:    ZenaDisplayDenom,
				Exponent: 18,
			},
		},
		Base:    ZenaChainDenom,
		Display: ZenaDisplayDenom,
		Name:    ZenaTokenName,
		Symbol:  ZenaTokenSymbol,
	}
}
