package constants_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	chainconfig "github.com/zenanetwork/zena/cmd/zenad/config"
	"github.com/zenanetwork/zena/testutil/constants"
	"github.com/zenanetwork/zena/zenad"
)

func TestRequireSameTestDenom(t *testing.T) {
	require.Equal(t,
		constants.ExampleAttoDenom,
		zenad.ExampleChainDenom,
		"test denoms should be the same across the repo",
	)
}

func TestRequireSameTestBech32Prefix(t *testing.T) {
	require.Equal(t,
		constants.ExampleBech32Prefix,
		chainconfig.Bech32Prefix,
		"bech32 prefixes should be the same across the repo",
	)
}

func TestRequireSameWZENAMainnet(t *testing.T) {
	require.Equal(t,
		constants.WZENAContractMainnet,
		zenad.WZENAContractMainnet,
		"wzena contract addresses should be the same across the repo",
	)
}
