#!/bin/bash

# =============================================================================
# Zena EVM Chain Setup Script
# =============================================================================
# This script sets up a complete Zena EVM chain with all features enabled
# Features:
# - Full EVM functionality with all precompiles
# - JSON-RPC and WebSocket support
# - Ethereum address display
# - Background execution
# - Production-ready configuration
# =============================================================================

set -euo pipefail

# =============================================================================
# Configuration Variables
# =============================================================================

# Chain Configuration
CHAINID="${CHAIN_ID:-4221}"
MONIKER="${MONIKER:-zena-evm-node}"
KEYRING="${KEYRING:-test}"  # Use 'test' for development, 'file' for production
KEYALGO="eth_secp256k1"
LOGLEVEL="${LOGLEVEL:-info}"

# Network Configuration
BASEFEE="${BASEFEE:-10000000000000000}"  # 0.01 zena
DENOM="${DENOM:-azena}"
DISPLAY_DENOM="${DISPLAY_DENOM:-zena}"

# Directories
HOMEDIR="${HOMEDIR:-$HOME/.zenad}"
LOGDIR="${LOGDIR:-$HOMEDIR/logs}"

# File Paths
CONFIG="$HOMEDIR/config/config.toml"
APP_TOML="$HOMEDIR/config/app.toml"
GENESIS="$HOMEDIR/config/genesis.json"
TMP_GENESIS="$HOMEDIR/config/tmp_genesis.json"
PIDFILE="$HOMEDIR/zenad.pid"
LOGFILE="$LOGDIR/zenad.log"

# Network Ports
RPC_PORT="${RPC_PORT:-26657}"
API_PORT="${API_PORT:-1317}"
GRPC_PORT="${GRPC_PORT:-9090}"
JSONRPC_PORT="${JSONRPC_PORT:-8545}"
WS_PORT="${WS_PORT:-8546}"
METRICS_PORT="${METRICS_PORT:-26660}"

# =============================================================================
# Colors for output
# =============================================================================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# =============================================================================
# Helper Functions
# =============================================================================

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
    exit 1
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO: $1${NC}"
}

# Validate dependencies
validate_dependencies() {
    local deps=("jq" "curl" "zenad")
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            error "$dep is required but not installed. Please install it first."
        fi
    done
}

# Create necessary directories
create_directories() {
    mkdir -p "$LOGDIR"
    mkdir -p "$HOMEDIR/config"
    mkdir -p "$HOMEDIR/data"
}

# Check if node is running
is_node_running() {
    if [[ -f "$PIDFILE" ]]; then
        local pid=$(cat "$PIDFILE")
        if kill -0 "$pid" 2>/dev/null; then
            return 0
        else
            rm -f "$PIDFILE"
            return 1
        fi
    fi
    return 1
}

# Stop running node
stop_node() {
    if is_node_running; then
        local pid=$(cat "$PIDFILE")
        log "Stopping zenad node (PID: $pid)..."
        kill -TERM "$pid" 2>/dev/null || true
        
        # Wait for graceful shutdown
        local count=0
        while kill -0 "$pid" 2>/dev/null && [[ $count -lt 30 ]]; do
            sleep 1
            ((count++))
        done
        
        # Force kill if still running
        if kill -0 "$pid" 2>/dev/null; then
            warn "Forcefully killing zenad node..."
            kill -KILL "$pid" 2>/dev/null || true
        fi
        
        rm -f "$PIDFILE"
        log "Node stopped successfully"
    else
        log "No running node found"
    fi
}

# Display Ethereum addresses for all keys
show_ethereum_addresses() {
    log "Ethereum addresses for all keys:"
    echo "=================================="
    
    # Get all keys from keyring
    local keys=$(zenad keys list --keyring-backend "$KEYRING" --home "$HOMEDIR" --output json 2>/dev/null)
    
    if [[ -n "$keys" && "$keys" != "[]" ]]; then
        echo "$keys" | jq -r '.[] | .name' | while read -r keyname; do
            if [[ -n "$keyname" ]]; then
                # Get cosmos address
                local cosmos_addr=$(zenad keys show "$keyname" -a --keyring-backend "$KEYRING" --home "$HOMEDIR" 2>/dev/null)
                
                # Get ethereum address using debug command
                local eth_addr=$(zenad debug addr "$cosmos_addr" 2>/dev/null | grep -E "Address \(hex\):" | awk '{print $3}' || echo "N/A")
                
                echo "  $keyname:"
                echo "    Cosmos:   $cosmos_addr"
                echo "    Ethereum: $eth_addr"
                echo ""
            fi
        done
    else
        warn "No keys found in keyring"
    fi
}

# =============================================================================
# Main Functions
# =============================================================================

# Install zenad binary
install_zenad() {
    log "Installing zenad binary..."
    if [[ "${BUILD_FOR_DEBUG:-false}" == "true" ]]; then
        make install COSMOS_BUILD_OPTIONS=nooptimization,nostrip
    else
        make install
    fi
    log "zenad binary installed successfully"
}

# Initialize chain configuration
init_chain() {
    log "Initializing chain configuration..."
    
    # Set client configuration
    zenad config set client chain-id "$CHAINID" --home "$HOMEDIR"
    zenad config set client keyring-backend "$KEYRING" --home "$HOMEDIR"
    zenad config set client node "tcp://localhost:$RPC_PORT" --home "$HOMEDIR"
    
    # Initialize the chain
    zenad init "$MONIKER" --chain-id "$CHAINID" --home "$HOMEDIR" --overwrite
    
    log "Chain initialized successfully"
}

# Setup keys with predefined mnemonics for quick setup
# ⚠️  SECURITY WARNING: For production use, generate new keys with unique mnemonics!
setup_keys() {
    log "Setting up keys..."
    
    echo -e "${RED}=== SECURITY WARNING ===${NC}"
    echo -e "${RED}Using predefined mnemonics for quick setup${NC}"
    echo -e "${RED}🔐 For PRODUCTION use, generate new keys with:${NC}"
    echo -e "${RED}   zenad keys add <keyname> --keyring-backend file${NC}"
    echo -e "${RED}=========================${NC}"
    echo ""
    
    if [[ "$KEYRING" == "file" ]]; then
        echo -e "${YELLOW}Press Enter to continue with predefined keys or Ctrl+C to cancel${NC}"
        read -r
    fi
    
    # Validator key
    local VAL_KEY="validator"
    local VAL_MNEMONIC="gesture inject test cycle original hollow east ridge hen combine junk child bacon zero hope comfort vacuum milk pitch cage oppose unhappy lunar seat"
    
    # User keys for operation
    local USER1_KEY="user1"
    local USER1_MNEMONIC="copper push brief egg scan entry inform record adjust fossil boss egg comic alien upon aspect dry avoid interest fury window hint race symptom"
    
    local USER2_KEY="user2"
    local USER2_MNEMONIC="maximum display century economy unlock van census kite error heart snow filter midnight usage egg venture cash kick motor survey drastic edge muffin visual"
    
    # Import keys - handle different keyring backends
    if [[ "$KEYRING" == "file" ]]; then
        # For file keyring, recommend generating new keys for security
        echo -e "${YELLOW}File keyring detected for production use!${NC}"
        echo ""
        
        # Check if keys already exist
        keys_exist=0
        for key in "$VAL_KEY" "$USER1_KEY" "$USER2_KEY"; do
            if zenad keys show "$key" --keyring-backend "$KEYRING" --home "$HOMEDIR" &>/dev/null; then
                keys_exist=1
                break
            fi
        done
        
        if [[ $keys_exist -eq 0 ]]; then
            echo -e "${RED}⚠️  No keys found in file keyring.${NC}"
            echo -e "${GREEN}For production security, we recommend generating new keys:${NC}"
            echo ""
            echo "  ./setup_zena_evm.sh generate-keys"
            echo ""
            echo -e "${YELLOW}Or continue with predefined keys (NOT recommended for production):${NC}"
            echo -n "Continue with predefined keys? (y/N): "
            read -r response
            
            if [[ "$response" != "y" && "$response" != "Y" ]]; then
                echo ""
                echo -e "${BLUE}To generate new keys, run:${NC}"
                echo "  ./setup_zena_evm.sh generate-keys"
                echo ""
                echo -e "${BLUE}Then initialize again:${NC}"
                echo "  ./setup_zena_evm.sh --keyring file init"
                exit 0
            fi
            
            echo ""
            echo -e "${YELLOW}Creating keys with predefined mnemonics...${NC}"
        fi
        
        # Try to create keys with predefined mnemonics
        for key_name in "$VAL_KEY" "$USER1_KEY" "$USER2_KEY"; do
            if ! zenad keys show "$key_name" --keyring-backend "$KEYRING" --home "$HOMEDIR" &>/dev/null; then
                echo -e "${BLUE}Creating $key_name key...${NC}"
                
                # Select appropriate mnemonic
                case $key_name in
                    "$VAL_KEY")
                        mnemonic="$VAL_MNEMONIC"
                        ;;
                    "$USER1_KEY")
                        mnemonic="$USER1_MNEMONIC"
                        ;;
                    "$USER2_KEY")
                        mnemonic="$USER2_MNEMONIC"
                        ;;
                esac
                
                # Try to create key with timeout
                if timeout 30s bash -c "echo '$mnemonic' | zenad keys add '$key_name' --recover --keyring-backend '$KEYRING' --algo '$KEYALGO' --home '$HOMEDIR'" 2>/dev/null; then
                    echo -e "${GREEN}✅ $key_name created successfully${NC}"
                else
                    warn "Failed to create $key_name automatically. Please create it manually:"
                    echo "  zenad keys add $key_name --keyring-backend $KEYRING --algo $KEYALGO --home $HOMEDIR"
                    echo ""
                fi
            fi
        done
    else
        # For test keyring, proceed normally
        if ! zenad keys show "$VAL_KEY" --keyring-backend "$KEYRING" --home "$HOMEDIR" &>/dev/null; then
            echo "$VAL_MNEMONIC" | zenad keys add "$VAL_KEY" --recover --keyring-backend "$KEYRING" --algo "$KEYALGO" --home "$HOMEDIR"
        fi
        
        if ! zenad keys show "$USER1_KEY" --keyring-backend "$KEYRING" --home "$HOMEDIR" &>/dev/null; then
            echo "$USER1_MNEMONIC" | zenad keys add "$USER1_KEY" --recover --keyring-backend "$KEYRING" --algo "$KEYALGO" --home "$HOMEDIR"
        fi
        
        if ! zenad keys show "$USER2_KEY" --keyring-backend "$KEYRING" --home "$HOMEDIR" &>/dev/null; then
            echo "$USER2_MNEMONIC" | zenad keys add "$USER2_KEY" --recover --keyring-backend "$KEYRING" --algo "$KEYALGO" --home "$HOMEDIR"
        fi
    fi
    
    log "Keys setup completed"
}

# Configure genesis file with full EVM support
configure_genesis() {
    log "Configuring genesis file..."
    
    # Update denomination parameters
    jq --arg denom "$DENOM" '.app_state.staking.params.bond_denom = $denom' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq --arg denom "$DENOM" '.app_state.gov.params.min_deposit[0].denom = $denom' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq --arg denom "$DENOM" '.app_state.gov.params.expedited_min_deposit[0].denom = $denom' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq --arg denom "$DENOM" '.app_state.mint.params.mint_denom = $denom' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Configure EVM parameters
    jq --arg denom "$DENOM" '.app_state.evm.params.evm_denom = $denom' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Enable ALL available precompiles
    jq '.app_state.evm.params.active_static_precompiles = [
        "0x0000000000000000000000000000000000000100",
        "0x0000000000000000000000000000000000000400",
        "0x0000000000000000000000000000000000000800",
        "0x0000000000000000000000000000000000000801",
        "0x0000000000000000000000000000000000000802",
        "0x0000000000000000000000000000000000000803",
        "0x0000000000000000000000000000000000000804",
        "0x0000000000000000000000000000000000000805",
        "0x0000000000000000000000000000000000000806",
        "0x0000000000000000000000000000000000000807"
    ]' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Configure denomination metadata
    jq --arg denom "$DENOM" --arg display_denom "$DISPLAY_DENOM" '.app_state.bank.denom_metadata = [{
        "description": "The native staking token for zena EVM chain",
        "denom_units": [
            {"denom": $denom, "exponent": 0, "aliases": ["atto" + $display_denom]},
            {"denom": $display_denom, "exponent": 18, "aliases": []}
        ],
        "base": $denom,
        "display": $display_denom,
        "name": "Zena Token",
        "symbol": "ZENA",
        "uri": "",
        "uri_hash": ""
    }]' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Enable ERC20 native precompiles
    jq '.app_state.erc20.params.native_precompiles = ["0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"]' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq --arg denom "$DENOM" '.app_state.erc20.token_pairs = [{
        "contract_owner": 1,
        "erc20_address": "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
        "denom": $denom,
        "enabled": true
    }]' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Configure consensus parameters
    jq '.consensus.params.block.max_gas = "50000000"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Configure fee market
    jq --arg basefee "$BASEFEE" '.app_state.feemarket.params.base_fee = $basefee' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    # Faster governance for development
    jq '.app_state.gov.params.max_deposit_period = "300s"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq '.app_state.gov.params.voting_period = "300s"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    jq '.app_state.gov.params.expedited_voting_period = "150s"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
    
    log "Genesis configuration completed"
}

# Configure app.toml for production
configure_app() {
    log "Configuring app.toml..."
    
    # Enable all APIs
    sed -i.bak 's/enable = false/enable = true/g' "$APP_TOML"
    
    # Configure API settings
    sed -i.bak "s/address = \"tcp:\/\/0.0.0.0:1317\"/address = \"tcp:\/\/0.0.0.0:$API_PORT\"/g" "$APP_TOML"
    sed -i.bak 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' "$APP_TOML"
    
    # Configure gRPC
    sed -i.bak "s/address = \"0.0.0.0:9090\"/address = \"0.0.0.0:$GRPC_PORT\"/g" "$APP_TOML"
    
    # Configure JSON-RPC
    sed -i.bak "s/address = \"127.0.0.1:8545\"/address = \"0.0.0.0:$JSONRPC_PORT\"/g" "$APP_TOML"
    sed -i.bak "s/ws-address = \"127.0.0.1:8546\"/ws-address = \"0.0.0.0:$WS_PORT\"/g" "$APP_TOML"
    sed -i.bak 's/api = "eth,net,web3"/api = "eth,txpool,personal,net,debug,web3,miner"/g' "$APP_TOML"
    sed -i.bak 's/enable-indexer = false/enable-indexer = true/g' "$APP_TOML"
    
    # Configure pruning for production
    sed -i.bak 's/pruning = "default"/pruning = "custom"/g' "$APP_TOML"
    sed -i.bak 's/pruning-keep-recent = "0"/pruning-keep-recent = "100"/g' "$APP_TOML"
    sed -i.bak 's/pruning-interval = "0"/pruning-interval = "10"/g' "$APP_TOML"
    
    # Enable metrics
    sed -i.bak 's/prometheus-retention-time = "0"/prometheus-retention-time = "1000000000000"/g' "$APP_TOML"
    
    # Production security settings
    sed -i.bak 's/max-open-connections = 1000/max-open-connections = 2000/g' "$APP_TOML"
    sed -i.bak 's/max-txs-per-conn = 100/max-txs-per-conn = 200/g' "$APP_TOML"
    
    # Enable state sync for faster sync
    sed -i.bak 's/enable = false/enable = true/g' "$APP_TOML"
    
    # Set minimum gas prices for production
    sed -i.bak "s/minimum-gas-prices = \"\"/minimum-gas-prices = \"0.001${DENOM}\"/g" "$APP_TOML"
    
    log "App configuration completed"
}

# Configure config.toml
configure_config() {
    log "Configuring config.toml..."
    
    # Configure RPC
    sed -i.bak "s/laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/0.0.0.0:$RPC_PORT\"/g" "$CONFIG"
    sed -i.bak 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["*"]/g' "$CONFIG"
    
    # Configure P2P
    sed -i.bak "s/laddr = \"tcp:\/\/0.0.0.0:26656\"/laddr = \"tcp:\/\/0.0.0.0:$((RPC_PORT-1))\"/g" "$CONFIG"
    
    # Enable metrics
    sed -i.bak 's/prometheus = false/prometheus = true/g' "$CONFIG"
    sed -i.bak "s/prometheus_listen_addr = \":26660\"/prometheus_listen_addr = \":$METRICS_PORT\"/g" "$CONFIG"
    
    # Configure logging
    sed -i.bak "s/log_level = \"info\"/log_level = \"$LOGLEVEL\"/g" "$CONFIG"
    
    log "Config configuration completed"
}

# Add genesis accounts
add_genesis_accounts() {
    log "Adding genesis accounts..."
    
    local val_key="validator"
    local user1_key="user1"
    local user2_key="user2"
    
    # Add accounts with substantial balances
    zenad genesis add-genesis-account "$val_key" "100000000000000000000000000$DENOM" --keyring-backend "$KEYRING" --home "$HOMEDIR"
    zenad genesis add-genesis-account "$user1_key" "10000000000000000000000$DENOM" --keyring-backend "$KEYRING" --home "$HOMEDIR"
    zenad genesis add-genesis-account "$user2_key" "10000000000000000000000$DENOM" --keyring-backend "$KEYRING" --home "$HOMEDIR"
    
    log "Genesis accounts added"
}

# Create genesis transaction
create_genesis_tx() {
    log "Creating genesis transaction..."
    
    local val_key="validator"
    local stake_amount="1000000000000000000000$DENOM"
    
    zenad genesis gentx "$val_key" "$stake_amount" \
        --gas-prices "${BASEFEE}${DENOM}" \
        --keyring-backend "$KEYRING" \
        --chain-id "$CHAINID" \
        --home "$HOMEDIR"
    
    # Collect genesis transactions
    zenad genesis collect-gentxs --home "$HOMEDIR"
    
    # Validate genesis
    zenad genesis validate-genesis --home "$HOMEDIR"
    
    log "Genesis transaction created and validated"
}

# Start node in background
start_node_background() {
    log "Starting zenad node in background..."
    
    # Create log directory
    mkdir -p "$LOGDIR"
    
    # Start node in background
    nohup zenad start \
        --home "$HOMEDIR" \
        --log_level "$LOGLEVEL" \
        --minimum-gas-prices="0.0001${DENOM}" \
        --json-rpc.api eth,txpool,personal,net,debug,web3,miner \
        --json-rpc.enable \
        --json-rpc.enable-indexer \
        --chain-id "$CHAINID" \
        > "$LOGFILE" 2>&1 &
    
    # Save PID
    echo $! > "$PIDFILE"
    
    log "Node started in background (PID: $(cat "$PIDFILE"))"
    log "Log file: $LOGFILE"
    
    # Wait for node to start
    sleep 3
    
    # Check if node is running
    if is_node_running; then
        log "Node is running successfully"
    else
        error "Failed to start node. Check log file: $LOGFILE"
    fi
}

# Start node in foreground
start_node_foreground() {
    log "Starting zenad node in foreground..."
    
    zenad start \
        --home "$HOMEDIR" \
        --log_level "$LOGLEVEL" \
        --minimum-gas-prices="0.0001${DENOM}" \
        --json-rpc.api eth,txpool,personal,net,debug,web3,miner \
        --json-rpc.enable \
        --json-rpc.enable-indexer \
        --chain-id "$CHAINID"
}

# Show node status
show_status() {
    log "Node Status:"
    echo "=============="
    
    if is_node_running; then
        echo "Status: Running (PID: $(cat "$PIDFILE"))"
        echo "Log file: $LOGFILE"
    else
        echo "Status: Not running"
    fi
    
    echo ""
    echo "Configuration:"
    echo "  Chain ID: $CHAINID"
    echo "  Moniker: $MONIKER"
    echo "  Home: $HOMEDIR"
    echo "  Keyring: $KEYRING"
    echo ""
    echo "Network Endpoints:"
    echo "  RPC: http://localhost:$RPC_PORT"
    echo "  REST API: http://localhost:$API_PORT"
    echo "  gRPC: localhost:$GRPC_PORT"
    echo "  JSON-RPC: http://localhost:$JSONRPC_PORT"
    echo "  WebSocket: ws://localhost:$WS_PORT"
    echo "  Metrics: http://localhost:$METRICS_PORT"
    echo ""
    
    show_ethereum_addresses
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS] COMMAND"
    echo ""
    echo "Commands:"
    echo "  init          Initialize the chain"
    echo "  start         Start the node in foreground"
    echo "  start-bg      Start the node in background"
    echo "  stop          Stop the running node"
    echo "  restart       Restart the node"
    echo "  status        Show node status"
    echo "  addresses     Show Ethereum addresses"
    echo "  logs          Show logs"
    echo "  clean         Clean all data (use with caution)"
    echo "  backup        Create backup of node data"
    echo "  restore       Restore node data from backup"
    echo "  version       Show version information"
    echo "  generate-keys Generate new keys for production use"
    echo ""
    echo "Options:"
    echo "  --chain-id CHAIN_ID      Chain ID (default: $CHAINID)"
    echo "  --moniker MONIKER        Node moniker (default: $MONIKER)"
    echo "  --home HOME_DIR          Home directory (default: $HOMEDIR)"
    echo "  --keyring KEYRING        Keyring backend (default: $KEYRING)"
    echo "  --log-level LEVEL        Log level (default: $LOGLEVEL)"
    echo "  --no-install             Skip binary installation"
    echo "  --help                   Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 init                  Initialize a new chain"
    echo "  $0 start-bg              Start node in background"
    echo "  $0 --chain-id 1234 init  Initialize with custom chain ID"
}

# =============================================================================
# Main Script Logic
# =============================================================================

# Parse command line arguments
install_binary=true
command=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --chain-id)
            CHAINID="$2"
            shift 2
            ;;
        --moniker)
            MONIKER="$2"
            shift 2
            ;;
        --home)
            HOMEDIR="$2"
            shift 2
            ;;
        --keyring)
            KEYRING="$2"
            shift 2
            ;;
        --log-level)
            LOGLEVEL="$2"
            shift 2
            ;;
        --no-install)
            install_binary=false
            shift
            ;;
        --help|-h)
            show_usage
            exit 0
            ;;
        init|start|start-bg|stop|restart|status|addresses|logs|clean|backup|restore|version|generate-keys)
            command="$1"
            shift
            ;;
        *)
            error "Unknown option: $1"
            ;;
    esac
done

# Validate command
if [[ -z "$command" ]]; then
    error "No command specified. Use --help for usage information."
fi

# Update paths after parsing arguments
CONFIG="$HOMEDIR/config/config.toml"
APP_TOML="$HOMEDIR/config/app.toml"
GENESIS="$HOMEDIR/config/genesis.json"
TMP_GENESIS="$HOMEDIR/config/tmp_genesis.json"
PIDFILE="$HOMEDIR/zenad.pid"
LOGFILE="$LOGDIR/zenad.log"

# Execute commands
case $command in
    init)
        validate_dependencies
        create_directories
        
        if [[ $install_binary == true ]]; then
            install_zenad
        fi
        
        init_chain
        setup_keys
        configure_genesis
        configure_app
        configure_config
        add_genesis_accounts
        create_genesis_tx
        
        log "Chain initialization completed!"
        echo ""
        show_status
        ;;
        
    start)
        if [[ ! -f "$GENESIS" ]]; then
            error "Chain not initialized. Run '$0 init' first."
        fi
        start_node_foreground
        ;;
        
    start-bg)
        if [[ ! -f "$GENESIS" ]]; then
            error "Chain not initialized. Run '$0 init' first."
        fi
        
        if is_node_running; then
            warn "Node is already running. Use 'restart' to restart."
            exit 1
        fi
        
        start_node_background
        show_status
        ;;
        
    stop)
        stop_node
        ;;
        
    restart)
        stop_node
        sleep 2
        start_node_background
        show_status
        ;;
        
    status)
        show_status
        ;;
        
    addresses)
        show_ethereum_addresses
        ;;
        
    logs)
        if [[ -f "$LOGFILE" ]]; then
            tail -f "$LOGFILE"
        else
            error "Log file not found: $LOGFILE"
        fi
        ;;
        
            clean)
        echo -e "${RED}WARNING: This will delete all blockchain data!${NC}"
        echo -n "Are you sure? (type 'yes' to confirm): "
        read -r confirmation
        if [[ "$confirmation" == "yes" ]]; then
            stop_node
            rm -rf "$HOMEDIR"
            log "All data cleaned successfully"
        else
            log "Operation cancelled"
        fi
        ;;
        
    backup)
        if [[ ! -d "$HOMEDIR" ]]; then
            error "Node data directory not found: $HOMEDIR"
        fi
        
        backup_dir="$HOMEDIR/../zena-backup-$(date +%Y%m%d_%H%M%S)"
        log "Creating backup at: $backup_dir"
        
        stop_node
        cp -r "$HOMEDIR" "$backup_dir"
        log "Backup created successfully: $backup_dir"
        ;;
        
    restore)
        echo -n "Enter backup directory path: "
        read -r backup_path
        
        if [[ ! -d "$backup_path" ]]; then
            error "Backup directory not found: $backup_path"
        fi
        
        echo -e "${RED}WARNING: This will replace current node data!${NC}"
        echo -n "Are you sure? (type 'yes' to confirm): "
        read -r confirmation
        if [[ "$confirmation" == "yes" ]]; then
            stop_node
            rm -rf "$HOMEDIR"
            cp -r "$backup_path" "$HOMEDIR"
            log "Data restored successfully from: $backup_path"
        else
            log "Operation cancelled"
        fi
        ;;
        
            version)
        echo "=================================="
        echo "Zena EVM Chain Version Information"
        echo "=================================="
        echo "Script Version: v1.0.0"
        echo "Chain ID: $CHAINID"
        echo "Denomination: $DENOM"
        echo "Display Denomination: $DISPLAY_DENOM"
        echo ""
        if command -v zenad &> /dev/null; then
            echo "zenad version:"
            zenad version
        else
            echo "zenad not found. Run 'make install' to build."
        fi
        echo ""
        echo "Build Information:"
        echo "  Go Version: $(go version 2>/dev/null || echo 'Go not found')"
        echo "  Build Date: $(date)"
        echo "  Platform: $(uname -s)/$(uname -m)"
        ;;
        
    generate-keys)
        echo "=================================="
        echo "🔐 Generating New Keys for Production"
        echo "=================================="
        echo ""
        
        # Force file keyring for production
        PROD_KEYRING="file"
        
        echo -e "${YELLOW}This will generate new keys with file keyring for production use.${NC}"
        echo -e "${YELLOW}You will be prompted to enter a secure passphrase for each key.${NC}"
        echo ""
        
        # Key names
        keys=("validator" "user1" "user2")
        
        for key in "${keys[@]}"; do
            echo -e "${BLUE}Generating key: $key${NC}"
            
            # Check if key already exists
            if zenad keys show "$key" --keyring-backend "$PROD_KEYRING" --home "$HOMEDIR" &>/dev/null; then
                echo -e "${YELLOW}Key '$key' already exists. Skipping...${NC}"
                continue
            fi
            
            # Generate new key
            if zenad keys add "$key" --keyring-backend "$PROD_KEYRING" --algo "$KEYALGO" --home "$HOMEDIR"; then
                echo -e "${GREEN}✅ Key '$key' generated successfully${NC}"
            else
                error "Failed to generate key '$key'"
            fi
            echo ""
        done
        
        echo -e "${GREEN}🎉 All keys generated successfully!${NC}"
        echo ""
        echo "To use these keys, run:"
        echo "  ./setup_zena_evm.sh --keyring file init"
        echo "  ./setup_zena_evm.sh --keyring file start-bg"
        echo ""
        echo -e "${RED}⚠️  IMPORTANT: Make sure to backup your keys!${NC}"
        echo "Key files are stored in: $HOMEDIR/keyring-file/"
        ;;
        
    *)
        error "Unknown command: $command"
        ;;
esac 