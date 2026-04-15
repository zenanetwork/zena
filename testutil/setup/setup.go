// Package setup provides shared test initialization for the Zena chain.
// Import this package (blank import) in test files to ensure the SDK global
// config is properly initialized with the correct bech32 prefixes and
// BIP44 coin type.
//
// Usage:
//
//	import _ "github.com/zenanetwork/zena/testutil/setup"
package setup

import (
	"github.com/zenanetwork/zena/config"
	testconstants "github.com/zenanetwork/zena/testutil/constants"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func init() {
	// Set bech32 prefixes for address encoding/decoding
	cfg := sdk.GetConfig()
	config.SetBech32Prefixes(cfg)
	config.SetBip44CoinType(cfg)

	// Set EVM coin info for modules that depend on it (e.g. precisebank)
	coinInfo := testconstants.ExampleChainCoinInfo[testconstants.ExampleChainID]
	_ = evmtypes.NewEVMConfigurator().
		WithEVMCoinInfo(coinInfo).
		Configure()
}
