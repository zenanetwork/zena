#!/usr/bin/env bash
set -euo pipefail

# Usage: ./error_edge_case.sh
# Requires in .env: CUSTOM_RPC, PRIVATE_KEY, ACCOUNT_2 (recipient), CONTRACT
# shellcheck source=../.env
# shellcheck disable=SC1091
source ../.env
export FOUNDRY_DISABLE_NIGHTLY_WARNING=1

RPC_URL=${CUSTOM_RPC:-http://127.0.0.1:8545}
PK=${PRIVATE_KEY:?}
RECIPIENT=${ACCOUNT_2:?}
CHAIN_ID=${CHAIN_ID:-15150}

# Ensure CONTRACT is set
if [ -z "${CONTRACT:-}" ]; then
	echo "Error: CONTRACT environment variable not set."
	exit 1
fi

# Amount exceeding typical balance (2000 tokens)
AMOUNT=2000000000000000000000

echo "🔄 Sending transfer exceeding balance..."

# 1) Send via cast and capture output (suppress exit)
echo "❌ Attempting transfer that should fail:"
echo "$ cast send \"$CONTRACT\" 'transfer(address,uint256)' \"$RECIPIENT\" \"$AMOUNT\" --rpc-url \"$RPC_URL\" --private-key \"[HIDDEN]\" --chain-id \"$CHAIN_ID\" --json"
OUTPUT=$(cast send \
	"$CONTRACT" \
	'transfer(address,uint256)' "$RECIPIENT" "$AMOUNT" \
	--rpc-url "$RPC_URL" \
	--private-key "$PK" \
	--chain-id "$CHAIN_ID" \
	--json 2>&1 || true)

# 2) A revert shows up either in the error text (gas estimation fails) or as a
#    mined receipt with status 0x0. Newer cast versions print errors as JSON
#    ({"level":"error",...,"message":...}), so match on the text first.
if echo "$OUTPUT" | grep -q -e 'execution reverted' -e 'ERC20InsufficientBalance'; then
	echo "✅ Transaction reverted as expected"
	echo
	echo "Revert detail:"
	echo "$OUTPUT"
	exit 0
fi

STATUS=$(echo "$OUTPUT" | jq -r '.status // empty' 2>/dev/null || true)
if [ "$STATUS" = "0x0" ] || [ "$STATUS" = "0" ]; then
	echo "✅ Transaction mined with failed status as expected"
	exit 0
fi

echo "❌ Expected a revert, but got:"
echo
echo "$OUTPUT"
exit 1
