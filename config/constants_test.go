package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zenanetwork/zena/config"
	serverconfig "github.com/zenanetwork/zena/server/config"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"
)

func TestMainnetDefaultsInSync(t *testing.T) {
	require.Equal(t, uint64(config.ZenaMainnetEVMChainID), evmtypes.DefaultEVMChainID)
	require.Equal(t, uint64(config.ZenaMainnetEVMChainID), uint64(serverconfig.DefaultEVMChainID))
	require.Equal(t, uint64(config.ZenaMainnetEVMChainID), uint64(config.EVMChainID))

	require.Equal(t, config.ZenaChainDenom, evmtypes.DefaultEVMDenom)
	require.Equal(t, config.ZenaChainDenom, evmtypes.DefaultEVMExtendedDenom)
	require.Equal(t, config.ZenaDisplayDenom, evmtypes.DefaultEVMDisplayDenom)
}

func TestNativeDenomMetadata(t *testing.T) {
	md := config.NativeDenomMetadata()
	require.NoError(t, md.Validate())
	require.Equal(t, "aznnt", md.Base)
	require.Equal(t, "znnt", md.Display)
	require.Equal(t, "ZNNT", md.Symbol)
	require.Len(t, md.DenomUnits, 2)
	require.Equal(t, uint32(18), md.DenomUnits[1].Exponent)
}

func TestUnimplementedPrecompilesNotAvailable(t *testing.T) {
	for _, addr := range evmtypes.UnimplementedStaticPrecompiles {
		require.NotContains(t, evmtypes.AvailableStaticPrecompiles, addr)
	}
}
