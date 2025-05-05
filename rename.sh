#!/bin/bash
   find . -type f -name "*.go" -exec sed -i '' 's|github.com/cosmos/evm|github.com/zenanetwork/zena|g' {} \;