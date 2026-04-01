package debug

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zenanetwork/zena/x/vm/statedb"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"
)

type EVMKeeper interface {
	CallEVM(ctx sdk.Context, stateDB *statedb.StateDB, abi abi.ABI, from, contract common.Address, commit, callFromPrecompile bool, gasCap *big.Int, method string, args ...interface{}) (*evmtypes.MsgEthereumTxResponse, error)
	CallEVMWithData(
		ctx sdk.Context,
		stateDB *statedb.StateDB,
		from common.Address,
		contract *common.Address,
		data []byte,
		commit bool,
		callFromPrecompile bool,
		gasCap *big.Int,
	) (*evmtypes.MsgEthereumTxResponse, error)
}
