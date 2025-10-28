#!/bin/bash

# Stop both zenad and geth nodes

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Stopping both zenad and geth...${NC}"

# Stop zenad
echo -e "${YELLOW}Stopping zenad...${NC}"
"$SCRIPT_DIR/zenad/stop-zenad.sh"

echo
echo -e "${YELLOW}Stopping geth...${NC}"
"$SCRIPT_DIR/geth/stop-geth.sh"

echo
echo -e "${GREEN}Both nodes stopped successfully${NC}"