package config

const (
	// ZenaChainDenom is the base denomination (18 decimals) of the Zenanet native token.
	ZenaChainDenom = "aznnt"

	// ZenaDisplayDenom is the display denomination of the Zenanet native token.
	ZenaDisplayDenom = "znnt"

	// ZenaTokenSymbol is the ticker symbol of the Zenanet native token.
	ZenaTokenSymbol = "ZNNT"

	// ZenaTokenName is the human readable name of the Zenanet native token.
	ZenaTokenName = "Zena Token"

	// ZenaMainnetEVMChainID is the EIP-155 chain ID of the Zenanet mainnet.
	ZenaMainnetEVMChainID = 15150

	// ExampleChainDenom is kept as an alias of ZenaChainDenom for existing callers.
	ExampleChainDenom = ZenaChainDenom

	// ExampleDisplayDenom is kept as an alias of ZenaDisplayDenom for existing callers.
	ExampleDisplayDenom = ZenaDisplayDenom

	// EighteenDecimalsChainID is the chain ID for the 18 decimals chain (the Zenanet mainnet).
	EighteenDecimalsChainID = ZenaMainnetEVMChainID

	// SixDecimalsChainID is the chain ID for the 6 decimals test chain.
	SixDecimalsChainID = 9002

	// TwelveDecimalsChainID is the chain ID for the 12 decimals test chain.
	TwelveDecimalsChainID = 9003

	// TwoDecimalsChainID is the chain ID for the 2 decimals test chain.
	TwoDecimalsChainID = 9004

	// TestChainID1 is test chain IDs for IBC E2E test
	TestChainID1 = 9005
	// TestChainID2 is test chain IDs for IBC E2E test
	TestChainID2 = 9006
)
