package constants_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zenanetwork/zena/config"
	"github.com/zenanetwork/zena/testutil/constants"
)

func TestRequireSameTestDenom(t *testing.T) {
	require.Equal(t,
		constants.ExampleAttoDenom,
		config.ExampleChainDenom,
		"test denoms should be the same across the repo",
	)
}

func TestRequireSameTestBech32Prefix(t *testing.T) {
	require.Equal(t,
		constants.ExampleBech32Prefix,
		config.Bech32Prefix,
		"bech32 prefixes should be the same across the repo",
	)
}

func TestRequireSameWEVMOSMainnet(t *testing.T) {
	require.Equal(t,
		constants.WZENAContractMainnet,
		config.WZENAContractMainnet,
		"wevmos contract addresses should be the same across the repo",
	)
}
