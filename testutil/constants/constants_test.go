package constants_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	config2 "github.com/zenanetwork/zena/cmd/zenad/config"
	"github.com/zenanetwork/zena/testutil/constants"
)

func TestRequireSameTestDenom(t *testing.T) {
	require.Equal(t,
		constants.ExampleAttoDenom,
		config2.ExampleChainDenom,
		"test denoms should be the same across the repo",
	)
}

func TestRequireSameTestBech32Prefix(t *testing.T) {
	require.Equal(t,
		constants.ExampleBech32Prefix,
		config2.Bech32Prefix,
		"bech32 prefixes should be the same across the repo",
	)
}

func TestRequireSameWZENAMainnet(t *testing.T) {
	require.Equal(t,
		constants.WZENAContractMainnet,
		config2.WZENAContractMainnet,
		"wzena contract addresses should be the same across the repo",
	)
}
