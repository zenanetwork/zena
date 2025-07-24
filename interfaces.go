package zena

import (
	"encoding/json"

	ibctesting "github.com/cosmos/ibc-go/v10/testing"
	erc20keeper "github.com/zenanetwork/zena/x/erc20/keeper"
	feemarketkeeper "github.com/zenanetwork/zena/x/feemarket/keeper"
	"github.com/zenanetwork/zena/x/ibc/callbacks/keeper"
	transferkeeper "github.com/zenanetwork/zena/x/ibc/transfer/keeper"
	precisebankkeeper "github.com/zenanetwork/zena/x/precisebank/keeper"
	evmkeeper "github.com/zenanetwork/zena/x/vm/keeper"

	storetypes "cosmossdk.io/store/types"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

// EvmApp defines the interface for an EVM application.
type EvmApp interface { //nolint:revive
	ibctesting.TestingApp
	runtime.AppI
	InterfaceRegistry() types.InterfaceRegistry
	ChainID() string
	GetEVMKeeper() *evmkeeper.Keeper
	GetErc20Keeper() *erc20keeper.Keeper
	SetErc20Keeper(erc20keeper.Keeper)
	GetGovKeeper() govkeeper.Keeper
	GetSlashingKeeper() slashingkeeper.Keeper
	GetEvidenceKeeper() *evidencekeeper.Keeper
	GetBankKeeper() bankkeeper.Keeper
	GetFeeMarketKeeper() *feemarketkeeper.Keeper
	GetAccountKeeper() authkeeper.AccountKeeper
	GetAuthzKeeper() authzkeeper.Keeper
	GetDistrKeeper() distrkeeper.Keeper
	GetStakingKeeper() *stakingkeeper.Keeper
	GetMintKeeper() mintkeeper.Keeper
	GetPreciseBankKeeper() *precisebankkeeper.Keeper
	GetFeeGrantKeeper() feegrantkeeper.Keeper
	GetCallbackKeeper() keeper.ContractKeeper
	GetTransferKeeper() transferkeeper.Keeper
	SetTransferKeeper(transferKeeper transferkeeper.Keeper)
	DefaultGenesis() map[string]json.RawMessage
	GetKey(storeKey string) *storetypes.KVStoreKey
	GetAnteHandler() sdk.AnteHandler
	GetSubspace(moduleName string) paramstypes.Subspace
	MsgServiceRouter() *baseapp.MsgServiceRouter
}
