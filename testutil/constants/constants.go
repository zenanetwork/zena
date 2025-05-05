package constants

import (
	"fmt"

	evmtypes "github.com/zenanetwork/zena/x/vm/types"
)

const (
	// DefaultGasPrice is used in testing as the default to use for transactions
	DefaultGasPrice = 20

	// ExampleAttoDenom provides an example denom for use in tests
	ExampleAttoDenom = "azena"

	// ExampleMicroDenom provides an example denom for use in tests
	ExampleMicroDenom = "uzena"

	// ExampleDisplayDenom provides an example display denom for use in tests
	ExampleDisplayDenom = "zena"

	// ExampleBech32Prefix provides an example Bech32 prefix for use in tests
	ExampleBech32Prefix = "zena"

	// ExampleEIP155ChainID provides an example EIP-155 chain ID for use in tests
	ExampleEIP155ChainID = 1

	// WZENAContractMainnet is the WZENA contract address for mainnet
	WZENAContractMainnet = "0xD4949664cD82660AaE99bEdc034a0deA8A0bd517"
	// WZENAContractTestnet is the WZENA contract address for testnet
	WZENAContractTestnet = "0xcc491f589b45d4a3c679016195b3fb87d7848210"
)

var (
	// ExampleChainIDPrefix provides a chain ID prefix for EIP-155 that can be used in tests
	ExampleChainIDPrefix = fmt.Sprintf("zena_%d", ExampleEIP155ChainID)

	// ExampleChainID provides a chain ID that can be used in tests
	ExampleChainID = ExampleChainIDPrefix + "-1"

	// SixDecimalsChainID provides a chain ID which is being set up with 6 decimals
	SixDecimalsChainID = "zenasix_6-1"

	// TwelveDecimalsChainID provides a chain ID which is being set up with 12 decimals
	TwelveDecimalsChainID = "ostwelve_8-2"

	// TwoDecimalsChainID provides a chain ID which is being set up with 2 decimals
	TwoDecimalsChainID = "ostwo_9-3"

	// ExampleChainCoinInfo provides the coin info for the example chain
	//
	// It is a map of the chain id and its corresponding EvmCoinInfo
	// that allows initializing the app with different coin info based on the
	// chain id
	ExampleChainCoinInfo = map[string]evmtypes.EvmCoinInfo{
		ExampleChainID: {
			Denom:         ExampleAttoDenom,
			ExtendedDenom: ExampleAttoDenom,
			DisplayDenom:  ExampleDisplayDenom,
			Decimals:      evmtypes.EighteenDecimals,
		},
		SixDecimalsChainID: {
			Denom:         "utest",
			ExtendedDenom: "atest",
			DisplayDenom:  "test",
			Decimals:      evmtypes.SixDecimals,
		},
		TwelveDecimalsChainID: {
			Denom:         "ptest2",
			ExtendedDenom: "atest2",
			DisplayDenom:  "test2",
			Decimals:      evmtypes.TwelveDecimals,
		},
		TwoDecimalsChainID: {
			Denom:         "ctest3",
			ExtendedDenom: "atest3",
			DisplayDenom:  "test3",
			Decimals:      evmtypes.TwoDecimals,
		},
	}
)
