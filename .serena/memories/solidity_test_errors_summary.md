# Solidity Test Failures Summary

## Test Results
- **Total Tests**: 37
- **Passing**: 10  
- **Failing**: 27
- **Exit Code**: 2 (failure)

## Error Categories

### 1. Insufficient Funds Errors (Most Common - ~15 failures)
**Root Cause**: Test accounts have zero balance when trying to execute transactions.

**Error Pattern**:
```
failed to check sender balance: sender balance < tx cost (0 < [amount]): insufficient funds
```

**Affected Tests**:
- Staking: createValidator, delegate, undelegate, redelegate
- Distribution: withdraw rewards, fund community pool, deposit validator rewards
- ERC20: transferFrom (approve transaction)
- Gov: proposal submission

**Fix Location**: `/Users/hwangjeong-yeon/workspace/blockchain/rebrand-test/zena/tests/solidity/init-node.sh`
- Lines 99-102: Genesis account allocations appear correct but accounts have 0 balance in practice
- Possible issue: accounts not properly funded or denomination mismatch

### 2. Bech32 Address Checksum Errors (~8 failures)
**Root Cause**: Hardcoded validator addresses in tests use incorrect checksums for the `zenanet` prefix.

**Error Pattern**:
```
decoding bech32 failed: invalid checksum (expected 8txm5m got 4xyrql)
decoding bech32 failed: invalid checksum (expected jdxw7v got qqyk2g)
```

**Affected Tests**:
- Staking: delegate, undelegate, redelegate
- Distribution: withdraw rewards, validator queries, deposit rewards

**Current Addresses Used** (INCORRECT):
- `zenanetvaloper10jmp6sgh4cc6zt3e8gw05wavvejgr5pw4xyrql` (checksum: 4xyrql, expected: 8txm5m)
- `zenanetvaloper1cml96vmptgw99syqrrz8az79xer2pcgpqqyk2g` (checksum: qqyk2g, expected: jdxw7v)

**Fix Required**: Recompute correct bech32 addresses with `zenanet` prefix for:
- VAL_KEY (mykey): 0x7cb61d4117ae31a12e393a1cfa3bac666481d02e
- USER1_KEY: 0xc6fe5d33615a1c52c08018c47e8bc53646a0e101

**Files to Fix**:
- `/Users/hwangjeong-yeon/workspace/blockchain/rebrand-test/zena/tests/solidity/suites/precompiles/test/1_staking/0_edge_case_revert.js`
- `/Users/hwangjeong-yeon/workspace/blockchain/rebrand-test/zena/tests/solidity/suites/precompiles/test/1_staking/2_delegate.js`
- `/Users/hwangjeong-yeon/workspace/blockchain/rebrand-test/zena/tests/solidity/suites/precompiles/test/1_staking/3_undelegate_and_cancel.js`
- `/Users/hwangjeong-yeon/workspace/blockchain/rebrand-test/zena/tests/solidity/suites/precompiles/test/1_staking/4_redelegate.js`
- `/Users/hwangjeong-yeon/workspace/blockchain/rebrand-test/zena/tests/solidity/suites/precompiles/test/2_distribution/2_withdraw_delegator_reward.js`
- `/Users/hwangjeong-yeon/workspace/blockchain/rebrand-test/zena/tests/solidity/suites/precompiles/test/2_distribution/4_withdraw_validator_commission.js`
- `/Users/hwangjeong-yeon/workspace/blockchain/rebrand-test/zena/tests/solidity/suites/precompiles/test/2_distribution/6_deposit_validator_rewards_pool.js`
- `/Users/hwangjeong-yeon/workspace/blockchain/rebrand-test/zena/tests/solidity/suites/precompiles/test/2_distribution/7_validator_queries.js`

### 3. Invalid Validator Address (~2 failures)
**Error Pattern**:
```
rpc error: code = InvalidArgument desc = invalid validator address
```

**Affected Tests**:
- Distribution: validatorSlashes query

**Likely Cause**: Related to bech32 checksum issues or validator not existing.

### 4. Validator Not Found (~1 failure)
**Error Pattern**:
```
rpc error: code = NotFound desc = SigningInfo not found for validator zenanetvalcons1qg9q7j9z7n8q7r9xm6lhrkurgaxaw97s4akjaj
```

**Affected Tests**:
- Slashing: getSigningInfo

**Likely Cause**: Consensus address checksum issue or validator not properly initialized.

### 5. ERC20 Balance Issues (~2 failures)
**Error Pattern**:
```
expected 0 to be above 0
ERC20: transfer amount exceeds balance
```

**Affected Tests**:
- ERC20: balance query, transfer

**Likely Cause**: Related to insufficient funds - accounts not funded with ERC20 tokens.

### 6. Empty Validator List (~1 failure)
**Error Pattern**:
```
AssertionError: expected [] to include 'zenanetvaloper1cml96vmptgw99syqrrz8az…'
```

**Affected Tests**:
- Distribution: delegatorValidators query

**Likely Cause**: No delegations due to insufficient funds errors in earlier tests.

### 7. Address Mismatch (~1 failure)
**Error Pattern**:
```
expected '0x2E1cA93285D58942Ba2e22fd893CCaC7c92…' to equal '0x020a0f48a2f4ce0f0cA6debF71DB83474dD…'
```

**Affected Tests**:
- Slashing: getSigningInfos

**Likely Cause**: Validator consensus address mismatch, possibly due to rebrand not updating address derivation.

## Bech32 Configuration
**Current Prefix**: `zenanet` (confirmed in code)
**Validator Prefix**: `zenanetvaloper`
**Consensus Prefix**: `zenanetvalcons`

**Source**: `/Users/hwangjeong-yeon/workspace/blockchain/rebrand-test/zena/zenad/cmd/zenad/config/config.go`

## Next Steps

### Priority 1: Fix Bech32 Checksums
1. Compute correct validator addresses using:
   - Prefix: `zenanetvaloper`
   - From hex addresses in init-node.sh
2. Update all hardcoded addresses in test files

### Priority 2: Debug Insufficient Funds
1. Verify genesis account allocation in init-node.sh works correctly
2. Check if denomination is correct (`atest`)
3. Verify accounts actually receive funds after node initialization
4. May need to add funding step before tests or fix genesis setup

### Priority 3: Verify Validator Initialization
1. Ensure validator is properly created during gentx
2. Verify signing info is created
3. Check consensus address derivation

## Test Files Requiring Updates
All files in `/Users/hwangjeong-yeon/workspace/blockchain/rebrand-test/zena/tests/solidity/suites/precompiles/test/`:
- 1_staking/*.js (4 files)
- 2_distribution/*.js (4 files)
- Potentially others if they reference validator addresses
