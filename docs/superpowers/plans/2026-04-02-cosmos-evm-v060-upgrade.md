# Cosmos EVM v0.6.0 Upgrade Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Upgrade Zena's Cosmos EVM layer from pre-v0.6.0 to v0.6.0, adding explicit StateDB parameters to all EVM call functions, removing the IBC transfer wrapper, extracting ERC20 conversion logic, and updating all interfaces.

**Architecture:** The core change is making StateDB an explicit parameter to `CallEVM`, `CallEVMWithData`, `ApplyMessage`, and `ApplyMessageWithConfig` instead of creating it internally. This enables precompiles to reuse their existing StateDB and prevents double state creation. The IBC transfer wrapper (`x/ibc/transfer/`) is removed — ERC20 conversions via IBC now go exclusively through the ICS20 precompile. A new public `ConvertERC20IntoCoinsForNativeToken` function is extracted to `convert.go` so both the msg server and precompiles can use it.

**Tech Stack:** Go 1.23.8, Cosmos SDK, go-ethereum, IBC-Go v10

**Working Directory:** `/Users/hwangjeong-yeon/workspace/blockchain/confero-projects/blockchain/zena-security-fix` (branch: `fix/security-audit`)

**Reference:** Upstream v0.6.0 code cloned at `/tmp/cosmos-evm-v060/`

---

## File Map

### Files to Modify (Production Code)

| File | Responsibility | Changes |
|------|---------------|---------|
| `x/vm/types/errors.go` | EVM error definitions | Add `ErrNilStateDB` |
| `x/vm/keeper/call_evm.go` | CallEVM, CallEVMWithData wrappers | Add `stateDB`, `callFromPrecompile` params |
| `x/vm/keeper/state_transition.go` | ApplyMessage, ApplyMessageWithConfig, ApplyTransaction | Add `stateDB`, `callFromPrecompile` params; stateDB created externally |
| `x/vm/keeper/grpc_query.go` | EthCall, EstimateGas, TraceTx | Create stateDB before calling ApplyMessageWithConfig |
| `x/erc20/types/interfaces.go` | EVMKeeper interface | Update function signatures; add new methods |
| `x/erc20/keeper/evm.go` | ERC20 EVM helpers (QueryERC20, BalanceOf) | Create stateDB, pass to CallEVM |
| `x/erc20/keeper/msg_server.go` | ConvertERC20, ConvertCoin handlers | Use new ConvertERC20IntoCoinsForNativeToken; update CallEVM calls |
| `x/ibc/callbacks/types/expected_keepers.go` | IBC callback keeper interfaces | Update EVMKeeper interface signatures |
| `x/ibc/callbacks/keeper/keeper.go` | IBC callback execution | Create stateDB, pass to CallEVM/CallEVMWithData |

### Files to Create

| File | Responsibility |
|------|---------------|
| `x/erc20/keeper/convert.go` | Public `ConvertERC20IntoCoinsForNativeToken` function |

### Files to Delete

| File | Reason |
|------|--------|
| `x/ibc/transfer/keeper/keeper.go` | IBC transfer wrapper removed in v0.6.0 |
| `x/ibc/transfer/keeper/msg_server.go` | IBC transfer wrapper removed in v0.6.0 |
| `x/ibc/transfer/ibc_module.go` | IBC transfer wrapper removed in v0.6.0 |
| `x/ibc/transfer/module.go` | IBC transfer wrapper removed in v0.6.0 |
| `x/ibc/transfer/types/channels.go` | IBC transfer wrapper removed in v0.6.0 |
| `x/ibc/transfer/types/interfaces.go` | IBC transfer wrapper removed in v0.6.0 |
| `x/ibc/transfer/v2/ibc_module.go` | IBC transfer wrapper removed in v0.6.0 |

### Files to Modify (App Wiring)

| File | Changes |
|------|---------|
| `zenad/app.go` | Remove x/ibc/transfer imports, use official IBC transfer keeper |
| `interfaces.go` | Remove TransferKeeper getter if custom |
| `precompiles/ics20/*.go` | Receive ERC20Keeper, update function calls |
| `precompiles/common/interfaces.go` | Update ERC20Keeper interface if needed |

### Files to Modify (Tests)

| File | Changes |
|------|---------|
| `tests/integration/x/vm/test_call_evm.go` | Update CallEVMWithData calls |
| `tests/integration/x/vm/test_state_transition.go` | Update ApplyMessage/ApplyMessageWithConfig calls |
| `tests/integration/x/vm/state_transition_benchmark.go` | Update benchmark calls |
| `tests/integration/x/erc20/test_msg_server.go` | Update for new public conversion function |
| `zenad/tests/testdata/debug/debug.go` | Update CallEVMWithData calls |
| `zenad/tests/integration/balance_handler/helper.go` | Update CallEVMWithData calls |
| `zenad/tests/ibc/helper.go` | Update CallEVMWithData calls |

---

## Task 1: Add ErrNilStateDB Error

**Files:**
- Modify: `x/vm/types/errors.go:15-97`

- [ ] **Step 1: Add error code constant**

In `x/vm/types/errors.go`, add `codeErrNilStateDB` to the const block:

```go
// After line 34 (codeErrInvalidPreinstall):
	codeErrInvalidPreinstall
	codeErrNilStateDB
```

- [ ] **Step 2: Add error variable**

After the `ErrInvalidPreinstall` variable (line 93), before `RevertSelector` (line 96), add:

```go
	// ErrNilStateDB returns an error when a nil stateDB is passed
	ErrNilStateDB = errorsmod.Register(ModuleName, codeErrNilStateDB, "stateDB cannot be nil")
```

- [ ] **Step 3: Verify compilation**

Run: `cd /Users/hwangjeong-yeon/workspace/blockchain/confero-projects/blockchain/zena-security-fix && go build ./x/vm/types/...`
Expected: Success, no errors

- [ ] **Step 4: Commit**

```bash
git add x/vm/types/errors.go
git commit -m "feat(vm): add ErrNilStateDB error type for v0.6.0 upgrade"
```

---

## Task 2: Update Core EVM Function Signatures

**Files:**
- Modify: `x/vm/keeper/call_evm.go`
- Modify: `x/vm/keeper/state_transition.go`

- [ ] **Step 1: Update CallEVM signature**

Replace the entire `CallEVM` function in `x/vm/keeper/call_evm.go`:

```go
// CallEVM performs a smart contract method call using given args.
// Note: if you call this from a precompile context, ensure that
// you use the existing stateDB.
func (k Keeper) CallEVM(ctx sdk.Context, stateDB *statedb.StateDB, abi abi.ABI, from, contract common.Address, commit, callFromPrecompile bool, gasCap *big.Int, method string, args ...interface{}) (*types.MsgEthereumTxResponse, error) {
	data, err := abi.Pack(method, args...)
	if err != nil {
		return nil, errorsmod.Wrap(
			types.ErrABIPack,
			errorsmod.Wrap(err, "failed to create transaction data").Error(),
		)
	}

	resp, err := k.CallEVMWithData(ctx, stateDB, from, &contract, data, commit, callFromPrecompile, gasCap)
	if err != nil {
		return resp, errorsmod.Wrapf(err, "contract call failed: method '%s', contract '%s'", method, contract)
	}
	return resp, nil
}
```

- [ ] **Step 2: Update CallEVMWithData signature**

Replace the entire `CallEVMWithData` function:

```go
// CallEVMWithData performs a smart contract method call using contract data.
// Note: if you call this from a precompile context, ensure that
// you use the existing stateDB.
func (k Keeper) CallEVMWithData(ctx sdk.Context, stateDB *statedb.StateDB, from common.Address, contract *common.Address, data []byte, commit bool, callFromPrecompile bool, gasCap *big.Int) (*types.MsgEthereumTxResponse, error) {
	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return nil, err
	}

	msg := core.Message{
		From:       from,
		To:         contract,
		Nonce:      nonce,
		Value:      big.NewInt(0),
		GasLimit:   config.DefaultGasCap,
		GasPrice:   big.NewInt(0),
		GasTipCap:  big.NewInt(0),
		GasFeeCap:  big.NewInt(0),
		Data:       data,
		AccessList: ethtypes.AccessList{},
	}

	res, err := k.ApplyMessage(ctx, stateDB, msg, nil, commit, callFromPrecompile, true)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		k.ResetGasMeterAndConsumeGas(ctx, ctx.GasMeter().Limit())
		return res, errorsmod.Wrap(types.ErrVMExecution, res.VmError)
	}

	ctx.GasMeter().ConsumeGas(res.GasUsed, "apply evm message")

	return res, nil
}
```

- [ ] **Step 3: Add statedb import to call_evm.go**

Update imports in `call_evm.go` — add:

```go
	"github.com/zenanetwork/zena/x/vm/statedb"
```

- [ ] **Step 4: Update ApplyMessage signature**

In `x/vm/keeper/state_transition.go`, replace the `ApplyMessage` function:

```go
// ApplyMessage calls ApplyMessageWithConfig with an empty TxConfig.
// Note: if you call this from a precompile context, ensure that
// you use the existing stateDB.
func (k *Keeper) ApplyMessage(ctx sdk.Context, stateDB *statedb.StateDB, msg core.Message, tracer *tracing.Hooks, commit, callFromPrecompile, internal bool) (*types.MsgEthereumTxResponse, error) {
	cfg, err := k.EVMConfig(ctx, ctx.BlockHeader().ProposerAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to load evm config")
	}

	txConfig := statedb.NewEmptyTxConfig()
	return k.ApplyMessageWithConfig(ctx, stateDB, msg, tracer, commit, callFromPrecompile, cfg, txConfig, internal, nil)
}
```

- [ ] **Step 5: Update ApplyMessageWithConfig signature**

Replace the `ApplyMessageWithConfig` function signature and add nil check. The key changes are:
1. `stateDB *statedb.StateDB` added as 2nd parameter
2. `callFromPrecompile bool` added after `commit`
3. StateDB is no longer created inside — it's passed in
4. Nil check added at the beginning
5. When `commit=true` and `callFromPrecompile=true`, flush instead of full commit

```go
func (k *Keeper) ApplyMessageWithConfig(ctx sdk.Context, stateDB *statedb.StateDB, msg core.Message, tracer *tracing.Hooks, commit bool, callFromPrecompile bool, cfg *statedb.EVMConfig, txConfig statedb.TxConfig, internal bool, overrides *rpctypes.StateOverride) (*types.MsgEthereumTxResponse, error) {
	var (
		ret   []byte
		vmErr error
	)
	if stateDB == nil {
		return nil, types.ErrNilStateDB
	}
	ethCfg := types.GetEthChainConfig()
	evm := k.NewEVMWithOverridePrecompiles(ctx, msg, cfg, tracer, stateDB, overrides == nil)
	// ... rest of function remains the same except stateDB is no longer created here
```

The internal `stateDB := statedb.New(ctx, k, txConfig)` line must be REMOVED from ApplyMessageWithConfig.

- [ ] **Step 6: Update ApplyTransaction to create stateDB externally**

In `ApplyTransaction`, the stateDB must now be created before calling `ApplyMessageWithConfig`:

Replace lines around 214-218:
```go
	// OLD:
	// res, err := k.ApplyMessageWithConfig(tmpCtx, *msg, nil, true, cfg, txConfig, false, nil)
	
	// NEW:
	stateDB := statedb.New(tmpCtx, k, txConfig)
	res, err := k.ApplyMessageWithConfig(tmpCtx, stateDB, *msg, nil, true, false, cfg, txConfig, false, nil)
```

- [ ] **Step 7: Verify compilation of vm package**

Run: `cd /Users/hwangjeong-yeon/workspace/blockchain/confero-projects/blockchain/zena-security-fix && go build ./x/vm/...`
Expected: Compilation errors in CALLERS (erc20, ibc, grpc_query) — that's expected, we fix those next.

- [ ] **Step 8: Commit**

```bash
git add x/vm/keeper/call_evm.go x/vm/keeper/state_transition.go
git commit -m "feat(vm): add stateDB and callFromPrecompile params to EVM call functions (v0.6.0)"
```

---

## Task 3: Update grpc_query.go Callers

**Files:**
- Modify: `x/vm/keeper/grpc_query.go`

- [ ] **Step 1: Update EthCall**

At line ~266, update the `ApplyMessageWithConfig` call. Create stateDB before the call:

```go
	// Before the call, add:
	stateDB := statedb.New(ctx, k, txConfig)
	// Update the call:
	res, err := k.ApplyMessageWithConfig(ctx, stateDB, *msg, nil, false, false, cfg, txConfig, false, overrides)
```

- [ ] **Step 2: Update EstimateGasInternal**

At line ~406, same pattern:

```go
	stateDB := statedb.New(tmpCtx, k, txConfig)
	rsp, err = k.ApplyMessageWithConfig(tmpCtx, stateDB, *msg, nil, false, false, cfg, txConfig, false, nil)
```

- [ ] **Step 3: Update TraceTx (predecessor tx processing)**

At line ~558:

```go
	stateDB := statedb.New(ctx, k, txConfig)
	rsp, _ := k.ApplyMessageWithConfig(ctx, stateDB, *msg, nil, true, false, cfg, txConfig, false, nil)
```

- [ ] **Step 4: Update traceTx (actual trace)**

At line ~836:

```go
	stateDB := statedb.New(ctx, k, txConfig)
	res, err := k.ApplyMessageWithConfig(ctx, stateDB, *msg, tracer.Hooks, commitMessage, false, cfg, txConfig, false, nil)
```

- [ ] **Step 5: Add statedb import**

Ensure `"github.com/zenanetwork/zena/x/vm/statedb"` is imported in grpc_query.go.

- [ ] **Step 6: Verify compilation**

Run: `go build ./x/vm/...`
Expected: Success for x/vm package

- [ ] **Step 7: Commit**

```bash
git add x/vm/keeper/grpc_query.go
git commit -m "feat(vm): update grpc query handlers with stateDB parameter"
```

---

## Task 4: Update ERC20 Interfaces

**Files:**
- Modify: `x/erc20/types/interfaces.go`

- [ ] **Step 1: Update EVMKeeper interface**

Replace the entire `EVMKeeper` interface with v0.6.0-compatible version:

```go
// EVMKeeper defines the expected EVM keeper interface used on erc20
type EVMKeeper interface {
	GetParams(ctx sdk.Context) evmtypes.Params
	GetAccountWithoutBalance(ctx sdk.Context, addr common.Address) *statedb.Account
	EstimateGasInternal(c context.Context, req *evmtypes.EthCallRequest, fromType evmtypes.CallType) (*evmtypes.EstimateGasResponse, error)
	ApplyMessage(ctx sdk.Context, stateDB *statedb.StateDB, msg core.Message, tracer *tracing.Hooks, commit, callFromPrecompile, internal bool) (*evmtypes.MsgEthereumTxResponse, error)
	DeleteAccount(ctx sdk.Context, addr common.Address) error
	IsAvailableStaticPrecompile(params *evmtypes.Params, address common.Address) bool
	CallEVM(ctx sdk.Context, stateDB *statedb.StateDB, abi abi.ABI, from, contract common.Address, commit, callFromPrecompile bool, gasCap *big.Int, method string, args ...interface{}) (*evmtypes.MsgEthereumTxResponse, error)
	CallEVMWithData(ctx sdk.Context, stateDB *statedb.StateDB, from common.Address, contract *common.Address, data []byte, commit bool, callFromPrecompile bool, gasCap *big.Int) (*evmtypes.MsgEthereumTxResponse, error)
	GetCode(ctx sdk.Context, hash common.Hash) []byte
	SetCode(ctx sdk.Context, hash []byte, bytecode []byte)
	SetAccount(ctx sdk.Context, address common.Address, account statedb.Account) error
	GetAccount(ctx sdk.Context, address common.Address) *statedb.Account
	IsContract(ctx sdk.Context, address common.Address) bool
	GetState(ctx sdk.Context, addr common.Address, key common.Hash) common.Hash
	GetCodeHash(ctx sdk.Context, addr common.Address) common.Hash
	ForEachStorage(ctx sdk.Context, addr common.Address, cb func(key, value common.Hash) bool)
	DeleteState(ctx sdk.Context, addr common.Address, key common.Hash)
	SetState(ctx sdk.Context, addr common.Address, key common.Hash, value []byte)
	DeleteCode(ctx sdk.Context, codeHash []byte)
	KVStoreKeys() map[string]*storetypes.KVStoreKey
}
```

- [ ] **Step 2: Add required imports**

Add to imports:

```go
	storetypes "cosmossdk.io/store/types"
```

- [ ] **Step 3: Verify compilation**

Run: `go build ./x/erc20/types/...`
Expected: Success

- [ ] **Step 4: Commit**

```bash
git add x/erc20/types/interfaces.go
git commit -m "feat(erc20): update EVMKeeper interface for v0.6.0 stateDB params"
```

---

## Task 5: Update IBC Callbacks Interfaces

**Files:**
- Modify: `x/ibc/callbacks/types/expected_keepers.go`

- [ ] **Step 1: Update EVMKeeper interface**

Replace the `EVMKeeper` interface with v0.6.0-compatible version:

```go
// EVMKeeper defines the expected EVM keeper interface used on erc20
type EVMKeeper interface {
	CallEVM(ctx sdk.Context, stateDB *statedb.StateDB, abi abi.ABI, from, contract common.Address, commit bool, callFromPrecompile bool, gasCap *big.Int, method string, args ...interface{}) (*evmtypes.MsgEthereumTxResponse, error)
	CallEVMWithData(ctx sdk.Context, stateDB *statedb.StateDB, from common.Address, contract *common.Address, data []byte, commit bool, callFromPrecompile bool, gasCap *big.Int) (*evmtypes.MsgEthereumTxResponse, error)
	GetAccountOrEmpty(ctx sdk.Context, addr common.Address) statedb.Account
	GetAccount(ctx sdk.Context, addr common.Address) *statedb.Account
	IsContract(ctx sdk.Context, addr common.Address) bool
	GetState(ctx sdk.Context, addr common.Address, key common.Hash) common.Hash
	GetCode(ctx sdk.Context, codeHash common.Hash) []byte
	GetCodeHash(ctx sdk.Context, addr common.Address) common.Hash
	ForEachStorage(ctx sdk.Context, addr common.Address, cb func(key common.Hash, value common.Hash) bool)
	SetAccount(ctx sdk.Context, addr common.Address, account statedb.Account) error
	DeleteState(ctx sdk.Context, addr common.Address, key common.Hash)
	SetState(ctx sdk.Context, addr common.Address, key common.Hash, value []byte)
	DeleteCode(ctx sdk.Context, codeHash []byte)
	SetCode(ctx sdk.Context, codeHash []byte, code []byte)
	DeleteAccount(ctx sdk.Context, addr common.Address) error
	KVStoreKeys() map[string]*storetypes.KVStoreKey
}
```

- [ ] **Step 2: Add required imports**

Add to imports:

```go
	storetypes "cosmossdk.io/store/types"
```

- [ ] **Step 3: Verify compilation**

Run: `go build ./x/ibc/callbacks/types/...`
Expected: Success

- [ ] **Step 4: Commit**

```bash
git add x/ibc/callbacks/types/expected_keepers.go
git commit -m "feat(ibc): update EVMKeeper interface for v0.6.0 stateDB params"
```

---

## Task 6: Update ERC20 Keeper EVM Helpers

**Files:**
- Modify: `x/erc20/keeper/evm.go`

- [ ] **Step 1: Update DeployERC20Contract**

Add stateDB creation and update CallEVMWithData call:

```go
func (k Keeper) DeployERC20Contract(
	ctx sdk.Context,
	coinMetadata banktypes.Metadata,
) (common.Address, error) {
	// ... existing code until data construction ...

	contractAddr := crypto.CreateAddress(types.ModuleAddress, nonce)
	stateDB := statedb.New(ctx, k.evmKeeper, statedb.NewEmptyTxConfig())
	_, err = k.evmKeeper.CallEVMWithData(ctx, stateDB, types.ModuleAddress, nil, data, true, false, nil)
	if err != nil {
		return common.Address{}, errorsmod.Wrapf(err, "failed to deploy contract for %s", coinMetadata.Name)
	}

	return contractAddr, nil
}
```

- [ ] **Step 2: Update QueryERC20**

Add stateDB creation for the decimals call:

```go
	stateDB := statedb.New(ctx, k.evmKeeper, statedb.NewEmptyTxConfig())
	res, err := k.evmKeeper.CallEVM(ctx, stateDB, erc20, types.ModuleAddress, contract, false, false, nil, "decimals")
```

- [ ] **Step 3: Update queryERC20String**

Add stateDB creation:

```go
	stateDB := statedb.New(ctx, k.evmKeeper, statedb.NewEmptyTxConfig())
	res, err := k.evmKeeper.CallEVM(ctx, stateDB, erc20, types.ModuleAddress, contract, false, false, nil, method)
```

- [ ] **Step 4: Update BalanceOf**

Add stateDB creation:

```go
func (k Keeper) BalanceOf(
	ctx sdk.Context,
	abi abi.ABI,
	contract, account common.Address,
) *big.Int {
	stateDB := statedb.New(ctx, k.evmKeeper, statedb.NewEmptyTxConfig())
	res, err := k.evmKeeper.CallEVM(ctx, stateDB, abi, types.ModuleAddress, contract, false, false, nil, "balanceOf", account)
	// ... rest unchanged ...
}
```

- [ ] **Step 5: Add statedb import**

Add to imports:

```go
	"github.com/zenanetwork/zena/x/vm/statedb"
```

Remove unused `banktypes` import if DeployERC20Contract doesn't use it (it does use it for coinMetadata).

- [ ] **Step 6: Verify compilation**

Run: `go build ./x/erc20/keeper/...`
Expected: May still fail due to msg_server.go — that's Task 7

- [ ] **Step 7: Commit**

```bash
git add x/erc20/keeper/evm.go
git commit -m "feat(erc20): update EVM helper calls with stateDB parameter"
```

---

## Task 7: Extract ConvertERC20IntoCoinsForNativeToken to convert.go

**Files:**
- Create: `x/erc20/keeper/convert.go`
- Modify: `x/erc20/keeper/msg_server.go`

- [ ] **Step 1: Create convert.go with public conversion function**

Create `x/erc20/keeper/convert.go` based on v0.6.0 upstream. This extracts the conversion logic from the private `convertERC20IntoCoinsForNativeToken` into a public function that accepts `stateDB` and `callFromPrecompile`:

```go
package keeper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	"github.com/zenanetwork/zena/contracts"
	"github.com/zenanetwork/zena/x/erc20/types"
	"github.com/zenanetwork/zena/x/vm/statedb"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// ConvertERC20IntoCoinsForNativeToken handles the erc20 conversion for a native erc20 token pair.
// This function is used by both the msg server and precompiles (like ICS20).
func (k Keeper) ConvertERC20IntoCoinsForNativeToken(ctx sdk.Context, stateDB *statedb.StateDB, contract common.Address, amount math.Int, receiver sdk.AccAddress, sender common.Address, commit bool, callFromPrecompile bool) (*types.MsgConvertERC20Response, error) {
	// Validate and get token pair
	pair, err := k.MintingEnabled(ctx, receiver, contract.Hex())
	if err != nil {
		return nil, err
	}

	// Check that this is a native ERC20 token
	if !pair.IsNativeERC20() {
		if pair.IsNativeCoin() {
			return nil, types.ErrNativeConversionDisabled
		}
		return nil, types.ErrUndefinedOwner
	}

	// Remove token pair if contract is suicided
	acc := k.evmKeeper.GetAccountWithoutBalance(ctx, pair.GetERC20Contract())
	if acc == nil || !acc.HasCodeHash() {
		k.DeleteTokenPair(ctx, pair)
		k.Logger(ctx).Info(
			"deleting selfdestructed token pair from state",
			"contract", pair.Erc20Address,
		)
		return nil, errors.Wrapf(
			types.ErrContractSelfDestructed,
			"contract %s has been self-destructed; token pair removed from state",
			pair.Erc20Address,
		)
	}

	erc20 := contracts.ERC20MinterBurnerDecimalsContract.ABI
	erc20Contract := pair.GetERC20Contract()
	balanceCoin := k.bankKeeper.GetBalance(ctx, receiver, pair.Denom)
	balanceToken := k.BalanceOf(ctx, erc20, erc20Contract, types.ModuleAddress)
	if balanceToken == nil {
		return nil, errors.Wrap(types.ErrEVMCall, "failed to retrieve balance")
	}

	// Escrow tokens on module account
	transferData, err := erc20.Pack("transfer", types.ModuleAddress, amount.BigInt())
	if err != nil {
		return nil, err
	}

	res, err := k.evmKeeper.CallEVMWithData(ctx, stateDB, sender, &erc20Contract, transferData, commit, callFromPrecompile, nil)
	if err != nil {
		return nil, err
	}

	// Check evm call response
	var unpackedRet types.ERC20BoolResponse
	if len(res.Ret) == 0 {
		if err := validateTransferEventExists(res.Logs, erc20Contract); err != nil {
			return nil, err
		}
	} else {
		if err := erc20.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
			return nil, err
		}
		if !unpackedRet.Value {
			return nil, errors.Wrap(errortypes.ErrLogic, "failed to execute transfer")
		}
	}

	// Check expected escrow balance after transfer execution
	coins := sdk.Coins{sdk.Coin{Denom: pair.Denom, Amount: amount}}
	tokens := coins[0].Amount.BigInt()
	balanceTokenAfter := k.BalanceOf(ctx, erc20, erc20Contract, types.ModuleAddress)
	if balanceTokenAfter == nil {
		return nil, errors.Wrap(types.ErrEVMCall, "failed to retrieve balance")
	}

	expToken := big.NewInt(0).Add(balanceToken, tokens)

	if r := balanceTokenAfter.Cmp(expToken); r != 0 {
		return nil, errors.Wrapf(
			types.ErrBalanceInvariance,
			"invalid token balance - expected: %v, actual: %v",
			expToken, balanceTokenAfter,
		)
	}

	// Mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return nil, err
	}

	// Send minted coins to the receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins); err != nil {
		return nil, err
	}

	// Check expected receiver balance after transfer
	balanceCoinAfter := k.bankKeeper.GetBalance(ctx, receiver, pair.Denom)
	expCoin := balanceCoin.Add(coins[0])

	if ok := balanceCoinAfter.Equal(expCoin); !ok {
		return nil, errors.Wrapf(
			types.ErrBalanceInvariance,
			"invalid coin balance - expected: %v, actual: %v",
			expCoin, balanceCoinAfter,
		)
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tx", "msg", "convert", "erc20", "total"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("coin", pair.Denom),
			},
		)

		if amount.IsInt64() {
			telemetry.IncrCounterWithLabels(
				[]string{"tx", "msg", "convert", "erc20", "amount", "total"},
				float32(amount.Int64()),
				[]metrics.Label{
					telemetry.NewLabel("denom", pair.Denom),
				},
			)
		}
	}()

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeConvertERC20,
				sdk.NewAttribute(sdk.AttributeKeySender, sender.Hex()),
				sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
				sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
				sdk.NewAttribute(types.AttributeKeyERC20Token, contract.Hex()),
			),
		},
	)

	return &types.MsgConvertERC20Response{}, nil
}
```

- [ ] **Step 2: Update ConvertERC20 in msg_server.go**

Replace the `ConvertERC20` function to use the new public function:

```go
func (k Keeper) ConvertERC20(
	goCtx context.Context,
	msg *types.MsgConvertERC20,
) (*types.MsgConvertERC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}
	sender := common.HexToAddress(msg.Sender)
	contract := common.HexToAddress(msg.ContractAddress)

	stateDB := statedb.New(ctx, k.evmKeeper, statedb.NewEmptyTxConfig())

	return k.ConvertERC20IntoCoinsForNativeToken(ctx, stateDB, contract, msg.Amount, receiver, sender, true, false)
}
```

- [ ] **Step 3: Remove old convertERC20IntoCoinsForNativeToken from msg_server.go**

Delete the private `convertERC20IntoCoinsForNativeToken` function (lines 66-199 in current msg_server.go).

- [ ] **Step 4: Update ConvertCoin in msg_server.go**

Update to match v0.6.0 pattern — use `ConvertCoinNativeERC20` with `callFromPrecompile=false`:

```go
func (k Keeper) ConvertCoin(
	goCtx context.Context,
	msg *types.MsgConvertCoin,
) (*types.MsgConvertCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	receiver := common.HexToAddress(msg.Receiver)

	pair, err := k.MintingEnabled(ctx, receiver.Bytes(), msg.Coin.Denom)
	if err != nil {
		return nil, err
	}

	switch {
	case pair.IsNativeERC20():
		acc := k.evmKeeper.GetAccountWithoutBalance(ctx, pair.GetERC20Contract())
		if acc == nil || !acc.HasCodeHash() {
			k.DeleteTokenPair(ctx, pair)
			k.Logger(ctx).Info(
				"deleting selfdestructed token pair from state",
				"contract", pair.Erc20Address,
			)
			return nil, sdkerrors.Wrapf(
				types.ErrContractSelfDestructed,
				"contract %s has been self-destructed; token pair removed from state",
				pair.Erc20Address,
			)
		}

		return nil, k.ConvertCoinNativeERC20(ctx, pair, msg.Coin.Amount, receiver, sender, false)
	case pair.IsNativeCoin():
		return nil, types.ErrNativeConversionDisabled
	}

	return nil, types.ErrUndefinedOwner
}
```

- [ ] **Step 5: Update ConvertCoinNativeERC20 to accept callFromPrecompile**

Update the function signature and add stateDB creation:

```go
func (k Keeper) ConvertCoinNativeERC20(ctx sdk.Context, pair types.TokenPair, amount math.Int, receiver common.Address, sender sdk.AccAddress, callFromPrecompile bool) error {
	if !amount.IsPositive() {
		return sdkerrors.Wrap(types.ErrNegativeToken, "converted coin amount must be positive")
	}

	erc20 := contracts.ERC20MinterBurnerDecimalsContract.ABI
	contract := pair.GetERC20Contract()

	balanceToken := k.BalanceOf(ctx, erc20, contract, receiver)
	if balanceToken == nil {
		return sdkerrors.Wrap(types.ErrEVMCall, "failed to retrieve balance")
	}

	// Escrow Coins on module account
	coins := sdk.Coins{{Denom: pair.Denom, Amount: amount}}
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, coins); err != nil {
		return sdkerrors.Wrap(err, "failed to escrow coins")
	}

	// Unescrow Tokens and send to receiver
	stateDB := statedb.New(ctx, k.evmKeeper, statedb.NewEmptyTxConfig())
	res, err := k.evmKeeper.CallEVM(ctx, stateDB, erc20, types.ModuleAddress, contract, true, callFromPrecompile, nil, "transfer", receiver, amount.BigInt())
	if err != nil {
		return err
	}

	// Check unpackedRet execution
	var unpackedRet types.ERC20BoolResponse
	if len(res.Ret) == 0 {
		if err := validateTransferEventExists(res.Logs, contract); err != nil {
			return err
		}
	} else {
		if err := erc20.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
			return err
		}
		if !unpackedRet.Value {
			return sdkerrors.Wrap(errortypes.ErrLogic, "failed to execute unescrow tokens from user")
		}
	}

	// Check expected Receiver balance after transfer execution
	balanceTokenAfter := k.BalanceOf(ctx, erc20, contract, receiver)
	if balanceTokenAfter == nil {
		return sdkerrors.Wrap(types.ErrEVMCall, "failed to retrieve balance")
	}

	exp := big.NewInt(0).Add(balanceToken, amount.BigInt())

	if r := balanceTokenAfter.Cmp(exp); r != 0 {
		return sdkerrors.Wrapf(
			types.ErrBalanceInvariance,
			"invalid token balance - expected: %v, actual: %v", exp, balanceTokenAfter,
		)
	}

	// Burn escrowed Coins
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to burn coins")
	}

	return nil
}
```

- [ ] **Step 6: Add statedb import to msg_server.go**

Add to imports:

```go
	"github.com/zenanetwork/zena/x/vm/statedb"
```

- [ ] **Step 7: Verify compilation**

Run: `go build ./x/erc20/...`
Expected: Success

- [ ] **Step 8: Commit**

```bash
git add x/erc20/keeper/convert.go x/erc20/keeper/msg_server.go
git commit -m "feat(erc20): extract ConvertERC20IntoCoinsForNativeToken and add stateDB params (v0.6.0)"
```

---

## Task 8: Update IBC Callbacks Keeper

**Files:**
- Modify: `x/ibc/callbacks/keeper/keeper.go`

- [ ] **Step 1: Add statedb import**

Add to imports:

```go
	"github.com/zenanetwork/zena/x/vm/statedb"
```

- [ ] **Step 2: Update IBCReceivePacketCallback**

After the `cachedCtx` setup (around line where `evmante.BuildEvmExecutionCtx` is called), create stateDB:

```go
	cachedCtx, writeFn := ctx.CacheContext()
	cachedCtx = evmante.BuildEvmExecutionCtx(cachedCtx).
		WithGasMeter(evmtypes.NewInfiniteGasMeterWithLimit(cbData.CommitGasLimit))
	stateDB := statedb.New(cachedCtx, k.evmKeeper, statedb.NewEmptyTxConfig())
```

Then update all CallEVM/CallEVMWithData calls in IBCReceivePacketCallback:

For the `approve` call:
```go
	res, err := k.evmKeeper.CallEVM(cachedCtx, stateDB, erc20.ABI, receiverHex, tokenPair.GetERC20Contract(), true, false, remainingGas, "approve", contractAddr, amountInt.BigInt())
```

For the `CallEVMWithData` call:
```go
	res, err = k.evmKeeper.CallEVMWithData(cachedCtx, stateDB, receiverHex, &contractAddr, cbData.Calldata, true, false, remainingGas)
```

- [ ] **Step 3: Update IBCOnAcknowledgementPacketCallback**

Same pattern — create stateDB after `cachedCtx` setup:

```go
	cachedCtx, writeFn := ctx.CacheContext()
	cachedCtx = evmante.BuildEvmExecutionCtx(cachedCtx).
		WithGasMeter(evmtypes.NewInfiniteGasMeterWithLimit(cbData.CommitGasLimit))
	stateDB := statedb.New(cachedCtx, k.evmKeeper, statedb.NewEmptyTxConfig())
```

Update the CallEVM call:
```go
	res, err := k.evmKeeper.CallEVM(cachedCtx, stateDB, *abi, sender, contractAddr, true, false, math.NewIntFromUint64(cachedCtx.GasMeter().GasRemaining()).BigInt(), "onPacketAcknowledgement",
		packet.GetSourceChannel(), packet.GetSourcePort(), packet.GetSequence(), packet.GetData(), acknowledgement)
```

- [ ] **Step 4: Update IBCOnTimeoutPacketCallback**

Same pattern:

```go
	cachedCtx, writeFn := ctx.CacheContext()
	cachedCtx = evmante.BuildEvmExecutionCtx(cachedCtx).
		WithGasMeter(evmtypes.NewInfiniteGasMeterWithLimit(cbData.CommitGasLimit))
	stateDB := statedb.New(cachedCtx, k.evmKeeper, statedb.NewEmptyTxConfig())
```

Update the CallEVM call:
```go
	res, err := k.evmKeeper.CallEVM(ctx, stateDB, *abi, sender, contractAddr, true, false, math.NewIntFromUint64(cachedCtx.GasMeter().GasRemaining()).BigInt(), "onPacketTimeout",
		packet.GetSourceChannel(), packet.GetSourcePort(), packet.GetSequence(), packet.GetData())
```

- [ ] **Step 5: Verify compilation**

Run: `go build ./x/ibc/callbacks/...`
Expected: Success

- [ ] **Step 6: Commit**

```bash
git add x/ibc/callbacks/keeper/keeper.go
git commit -m "feat(ibc): update callback keeper with stateDB parameter (v0.6.0)"
```

---

## Task 9: Remove IBC Transfer Wrapper

**Files:**
- Delete: `x/ibc/transfer/` (entire directory)
- Modify: App wiring files that reference it

**IMPORTANT NOTE:** This task is the most disruptive and may require extensive app.go modifications. The IBC transfer wrapper provides auto-ERC20 conversion in the transfer flow. Removing it means ERC20 transfers via IBC can ONLY go through the ICS20 precompile.

- [ ] **Step 1: Identify all imports of x/ibc/transfer**

Run: `grep -r "x/ibc/transfer" --include="*.go" -l` (excluding test files initially)

Document every file that imports from `x/ibc/transfer`.

- [ ] **Step 2: Delete x/ibc/transfer directory**

```bash
rm -rf x/ibc/transfer/
```

- [ ] **Step 3: Update zenad/app.go**

Replace imports:
```go
// OLD:
// transferkeeper "github.com/zenanetwork/zena/x/ibc/transfer/keeper"
// ibctransfer "github.com/zenanetwork/zena/x/ibc/transfer"

// NEW:
transfer "github.com/cosmos/ibc-go/v10/modules/apps/transfer"
transferkeeper "github.com/cosmos/ibc-go/v10/modules/apps/transfer/keeper"
```

Remove `erc20Keeper` parameter from TransferKeeper initialization.

Update `BasicModuleManager` registration to use `transfer.AppModuleBasic{}`.

- [ ] **Step 4: Update precompiles/ics20 if needed**

If the ICS20 precompile references the custom transfer keeper, update to accept ERC20Keeper for direct conversion.

- [ ] **Step 5: Update interfaces.go**

Remove any TransferKeeper getter that references the custom wrapper.

- [ ] **Step 6: Fix compilation errors iteratively**

Run: `go build ./...`
Fix each compilation error by replacing custom transfer keeper references with official IBC transfer keeper.

- [ ] **Step 7: Verify full build**

Run: `make build`
Expected: Success

- [ ] **Step 8: Commit**

```bash
git add -A
git commit -m "feat(ibc): remove custom IBC transfer wrapper, use official transfer keeper (v0.6.0)"
```

---

## Task 10: Update Test Files

**Files:**
- Modify: `tests/integration/x/vm/test_call_evm.go`
- Modify: `tests/integration/x/vm/test_state_transition.go`
- Modify: `tests/integration/x/vm/state_transition_benchmark.go`
- Modify: `tests/integration/x/erc20/test_msg_server.go`
- Modify: `zenad/tests/testdata/debug/debug.go`
- Modify: `zenad/tests/integration/balance_handler/helper.go`
- Modify: `zenad/tests/ibc/helper.go`

- [ ] **Step 1: Update test_call_evm.go**

For each `CallEVMWithData` call, add stateDB creation:

```go
stateDB := statedb.New(ctx, evmKeeper, statedb.NewEmptyTxConfig())
// Update calls to include stateDB, false (callFromPrecompile)
res, err := evmKeeper.CallEVMWithData(ctx, stateDB, from, nil, data, true, false, nil)
```

- [ ] **Step 2: Update test_state_transition.go**

For `ApplyMessage` calls:
```go
stateDB := statedb.New(ctx, evmKeeper, statedb.NewEmptyTxConfig())
res, err := evmKeeper.ApplyMessage(ctx, stateDB, *coreMsg, tracer, true, false, false)
```

For `ApplyMessageWithConfig` calls:
```go
stateDB := statedb.New(ctx, evmKeeper, txConfig)
res, err := evmKeeper.ApplyMessageWithConfig(ctx, stateDB, msg, nil, true, false, config, txConfig, false, tc.overrides)
```

- [ ] **Step 3: Update benchmarks**

Same pattern for benchmark files.

- [ ] **Step 4: Update zenad test helpers**

For `debug.go`, `helper.go` files in zenad tests:
```go
stateDB := statedb.New(ctx, evmKeeper, statedb.NewEmptyTxConfig())
res, err := evmKeeper.CallEVMWithData(ctx, stateDB, from, nil, data, true, false, nil)
```

- [ ] **Step 5: Update ERC20 integration tests**

Update mock expectations and test calls for the new function signatures.

- [ ] **Step 6: Run unit tests**

Run: `cd /Users/hwangjeong-yeon/workspace/blockchain/confero-projects/blockchain/zena-security-fix && go test -race -tags=test ./x/vm/... ./x/erc20/... ./x/ibc/callbacks/...`
Expected: All pass

- [ ] **Step 7: Run full test suite**

Run: `make test-unit`
Expected: All pass (except known Bech32 prefix issues in integration tests)

- [ ] **Step 8: Commit**

```bash
git add -A
git commit -m "test: update all test files for v0.6.0 stateDB parameter changes"
```

---

## Task 11: Update Mock Files

**Files:**
- Modify: Any mock files that implement EVMKeeper interface

- [ ] **Step 1: Regenerate mocks**

Run: `make mocks`

If `make mocks` is not available or fails, manually update mock files to match the new interface signatures.

- [ ] **Step 2: Verify compilation**

Run: `go build ./...`
Expected: Success

- [ ] **Step 3: Run tests**

Run: `go test -race -tags=test ./...`
Expected: All pass

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "chore: regenerate mocks for v0.6.0 interface changes"
```

---

## Task 12: Final Verification

- [ ] **Step 1: Full build**

Run: `make build`
Expected: Success, `zenad` binary produced

- [ ] **Step 2: Lint check**

Run: `make lint-go`
Expected: No new lint errors

- [ ] **Step 3: Full unit test**

Run: `make test-unit`
Expected: All pass

- [ ] **Step 4: Verify diff is clean**

Run: `git status && git log --oneline -10`
Expected: Clean working tree, all commits present

---

## Important Notes

### validateTransferEventExists Signature Difference

The Zena codebase has an extended `validateTransferEventExists` that takes `(logs, contract, from, to, amount)` while v0.6.0 only takes `(logs, contract)`. In `convert.go`, use the v0.6.0 simplified version. Check if the extended version is still needed elsewhere and keep it if so.

### MintingEnabled Signature Difference

Zena's `MintingEnabled` takes `(ctx, sender, receiver, contractAddr)` while v0.6.0 takes `(ctx, receiver, contractAddr)`. Keep Zena's version in `msg_server.go` but use v0.6.0's version in `convert.go`. Verify which signature the keeper actually implements.

### CacheContext Usage

Zena's current `msg_server.go` uses `CacheContext` for atomicity (added in the security audit fix). The v0.6.0 `convert.go` does NOT use CacheContext — it relies on the stateDB commit/revert mechanism instead. When porting, decide whether to keep the CacheContext safety or follow upstream. Recommendation: follow upstream since stateDB provides atomic state management.

### IBC Transfer Removal Risk

Removing `x/ibc/transfer/` breaks the ability to auto-convert ERC20 tokens during Cosmos IBC transfers. After this change, users must use the ICS20 precompile (EVM) for ERC20 IBC transfers. Ensure the ICS20 precompile is properly updated before removing the wrapper.
