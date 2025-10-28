package testutil

import (
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
)

var (
	UosmoDenom    = transfertypes.NewDenom("uosmo", transfertypes.NewHop(transfertypes.PortID, "channel-0"))
	UosmoIbcDenom = UosmoDenom.IBCDenom()

	UatomDenom    = transfertypes.NewDenom("uznnt", transfertypes.NewHop(transfertypes.PortID, "channel-1"))
	UatomIbcDenom = UatomDenom.IBCDenom()

	UAtomDenom    = transfertypes.NewDenom("aznnt", transfertypes.NewHop(transfertypes.PortID, "channel-0"))
	UAtomIbcDenom = UatomDenom.IBCDenom()

	UatomOsmoDenom = transfertypes.NewDenom(
		"uznnt",
		transfertypes.NewHop(transfertypes.PortID, "channel-0"),
		transfertypes.NewHop(transfertypes.PortID, "channel-1"),
	)
	UatomOsmoIbcDenom = UatomOsmoDenom.IBCDenom()

	AatomDenom    = transfertypes.NewDenom("aznnt", transfertypes.NewHop(transfertypes.PortID, "channel-0"))
	AatomIbcDenom = AatomDenom.IBCDenom()
)
