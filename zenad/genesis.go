package zenad

import (
	"encoding/json"

	"github.com/zenanetwork/zena/config"
	erc20types "github.com/zenanetwork/zena/x/erc20/types"
	feemarkettypes "github.com/zenanetwork/zena/x/feemarket/types"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// GenesisState of the blockchain is represented here as a map of raw json
// messages key'd by an identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewEVMGenesisState returns the default genesis state for the EVM module.
//
// It uses the Zenanet native denom (aznnt), enables every static precompile
// implemented by this binary and includes the default preinstalls.
func NewEVMGenesisState() *evmtypes.GenesisState {
	evmGenState := evmtypes.DefaultGenesisState()
	evmGenState.Params.EvmDenom = config.ZenaChainDenom
	evmGenState.Params.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: config.ZenaChainDenom}
	evmGenState.Params.ActiveStaticPrecompiles = evmtypes.AvailableStaticPrecompiles
	evmGenState.Preinstalls = evmtypes.DefaultPreinstalls

	return evmGenState
}

// NewErc20GenesisState returns the default genesis state for the ERC20 module.
//
// NOTE: no native token pair is registered by default. The address of the
// native ERC-20 representation (WZENA) is network specific, so it must be set
// explicitly in the network's genesis (see local_node.sh) instead of being
// hard coded here.
func NewErc20GenesisState() *erc20types.GenesisState {
	return erc20types.DefaultGenesisState()
}

// NewMintGenesisState returns the default genesis state for the mint module.
func NewMintGenesisState() *minttypes.GenesisState {
	mintGenState := minttypes.DefaultGenesisState()
	mintGenState.Params.MintDenom = config.ZenaChainDenom

	return mintGenState
}

// NewFeeMarketGenesisState returns the default genesis state for the feemarket module.
//
// NOTE: the base fee is disabled for local test networks.
func NewFeeMarketGenesisState() *feemarkettypes.GenesisState {
	feeMarketGenState := feemarkettypes.DefaultGenesisState()
	feeMarketGenState.Params.NoBaseFee = true

	return feeMarketGenState
}
