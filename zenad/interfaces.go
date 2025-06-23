package zenad

import (
	cmn "github.com/zenanetwork/zena/precompiles/common"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"
)

type BankKeeper interface {
	evmtypes.BankKeeper
	cmn.BankKeeper
}
