#!/bin/bash

# BitBadges Chain Startup Script
# This script builds, initializes, and starts the BitBadges chain locally
# It handles timeouts and evaluates outputs since chains are daemons

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CHAIN_ID="bitbadges-1"
NODE_HOME="${HOME}/.bitbadgeschain"
BINARY_NAME="bitbadgeschaind"
BINARY_PATH="${SCRIPT_DIR}/${BINARY_NAME}"
BUILD_TIMEOUT=300  # 5 minutes for build
INIT_TIMEOUT=60    # 1 minute for init
START_TIMEOUT=30   # 30 seconds to check if chain started
BLOCK_WAIT_TIMEOUT=60  # 1 minute to wait for first block

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_command() {
    if ! command -v "$1" &> /dev/null; then
        log_error "$1 is not installed or not in PATH"
        return 1
    fi
    return 0
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    check_command "go" || exit 1
    check_command "make" || exit 1
    
    log_info "Prerequisites check passed"
}

# Build the binary
build_binary() {
    # Skip build if SKIP_BUILD is set or binary exists and is executable
    if [ -n "${SKIP_BUILD:-}" ] || ([ -f "${BINARY_PATH}" ] && [ -x "${BINARY_PATH}" ]); then
        if [ -f "${BINARY_PATH}" ]; then
            log_info "Using existing binary at ${BINARY_PATH}"
            return 0
        fi
    fi
    
    log_info "Building ${BINARY_NAME}..."
    
    if [ -f "${BINARY_PATH}" ]; then
        if [ -z "${SKIP_PROMPTS:-}" ]; then
            log_warn "Binary already exists at ${BINARY_PATH}"
            read -p "Rebuild? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                log_info "Skipping build, using existing binary"
                return 0
            fi
        else
            log_info "Binary exists, rebuilding..."
        fi
    fi
    
    cd "${SCRIPT_DIR}"
    
    # Fix GOROOT if incorrectly set
    if [ -n "${GOROOT:-}" ] && [ ! -d "${GOROOT}/src" ]; then
        log_warn "GOROOT is set incorrectly (${GOROOT}), unsetting..."
        unset GOROOT
    fi
    
    log_info "Running 'go build ./cmd/bitbadgeschaind'..."
    if timeout "${BUILD_TIMEOUT}" go build -o "${BINARY_PATH}" ./cmd/bitbadgeschaind; then
        log_info "Build successful: ${BINARY_PATH}"
        chmod +x "${BINARY_PATH}"
    else
        log_error "Build failed or timed out after ${BUILD_TIMEOUT}s"
        log_error "If build issues persist, use SKIP_BUILD=1 to use existing binary"
        exit 1
    fi
}

# Initialize chain if needed
init_chain() {
    log_info "Checking if chain is initialized..."
    
    if [ -d "${NODE_HOME}/config" ] && [ -f "${NODE_HOME}/config/genesis.json" ]; then
        # Check if genesis has validators (gentx files or validators in staking module)
        local has_validators=false
        if [ -d "${NODE_HOME}/config/gentx" ] && [ -n "$(ls -A "${NODE_HOME}/config/gentx" 2>/dev/null)" ]; then
            has_validators=true
        elif grep -q '"validators":\s*\[' "${NODE_HOME}/config/genesis.json" && ! grep -q '"validators":\s*\[\s*\]' "${NODE_HOME}/config/genesis.json"; then
            # Check if validators array is non-empty in staking module
            if python3 -c "import json; d=json.load(open('${NODE_HOME}/config/genesis.json')); vs=d.get('app_state', {}).get('staking', {}).get('validators', []); exit(0 if len(vs) > 0 else 1)" 2>/dev/null; then
                has_validators=true
            fi
        fi
        
        if [ "$has_validators" = false ]; then
            log_warn "Chain initialized but no validators found in genesis"
            if [ -z "${SKIP_PROMPTS:-}" ]; then
                log_warn "Chain will not start without validators. Reinitialization required."
                read -p "Reinitialize? This will delete existing data (y/N): " -n 1 -r
                echo
                if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                    log_error "Cannot proceed without validators. Exiting."
                    exit 1
                fi
            else
                log_info "No validators found, reinitializing..."
            fi
            log_warn "Removing existing chain data..."
            rm -rf "${NODE_HOME}"
        elif [ -z "${SKIP_PROMPTS:-}" ]; then
            log_warn "Chain already initialized at ${NODE_HOME}"
            read -p "Reinitialize? This will delete existing data (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                log_info "Skipping initialization, using existing chain data"
                return 0
            fi
            log_warn "Removing existing chain data..."
            rm -rf "${NODE_HOME}"
        else
            log_info "Chain already initialized, reinitializing..."
            log_warn "Removing existing chain data..."
            rm -rf "${NODE_HOME}"
        fi
    fi
    
    log_info "Initializing chain with chain-id: ${CHAIN_ID}"
    
    # Ensure node home directory exists
    mkdir -p "${NODE_HOME}/config"
    
    # Try to use Ignite to generate everything first (accounts + genesis)
    # This is needed because default genesis doesn't have validators with enough stake
    if [ -f "${SCRIPT_DIR}/config.yml" ] && command -v ignite &> /dev/null; then
        log_info "Found config.yml, using Ignite to generate accounts and genesis..."
        cd "${SCRIPT_DIR}"
        
        # Ignite generates everything in ~/.ignite/chains/${CHAIN_ID}/
        local ignite_chain_dir="${HOME}/.ignite/chains/${CHAIN_ID}"
        local ignite_genesis="${ignite_chain_dir}/config/genesis.json"
        
        # Remove old ignite chain data if it exists
        if [ -d "${ignite_chain_dir}" ]; then
            log_info "Removing old Ignite chain data..."
            rm -rf "${ignite_chain_dir}"
        fi
        
        # Create accounts for Ignite BEFORE running init
        # Try both default location and Ignite's chain directory
        log_info "Creating accounts for Ignite..."
        
        # First, create in Ignite's chain directory (where Ignite will look)
        local ignite_chain_dir="${HOME}/.ignite/chains/${CHAIN_ID}"
        mkdir -p "${ignite_chain_dir}"
        create_accounts_for_ignite_in_dir "${ignite_chain_dir}"
        
        # Also create in default location as fallback
        create_accounts_for_ignite_fixed
        
        # Try Ignite first, but if it fails due to keyring issues, fall back to manual genesis creation
        # NOTE: Ignite has known compatibility issues with EVM keyring (eth_secp256k1 default)
        # Even with accounts created in both default and Ignite chain directories, Ignite may fail
        # The fallback manual genesis creation works perfectly and produces a fully functional chain
        log_info "Attempting to use Ignite for genesis generation..."
        log_info "Note: Ignite may fail due to EVM keyring compatibility - fallback will be used if needed"
        local ignite_output
        if ignite_output=$(timeout "${INIT_TIMEOUT}" ignite chain init --skip-proto 2>&1); then
            # Check if genesis was generated
            if [ -f "${ignite_genesis}" ]; then
                log_info "Ignite generated genesis successfully"
                
                # Copy entire config directory from Ignite
                if [ -d "${ignite_chain_dir}/config" ]; then
                    log_info "Copying Ignite-generated config to node home..."
                    cp -r "${ignite_chain_dir}/config"/* "${NODE_HOME}/config/"
                    log_info "Genesis and config copied successfully"
                    return 0
                else
                    log_warn "Ignite config directory not found"
                fi
            fi
        fi
        
        # Ignite failed - create genesis manually from config.yml
        log_warn "Ignite failed (likely keyring issue), creating genesis manually from config.yml..."
        log_warn "Ignite output: ${ignite_output}"
        create_genesis_from_config
    else
        # Use default initialization
        log_info "No config.yml or Ignite not available, using default initialization"
        log_info "Initializing chain with bitbadgeschaind..."
        # Remove directory if it exists to ensure clean init
        rm -rf "${NODE_HOME}"
        mkdir -p "${NODE_HOME}/config"
        if "${BINARY_PATH}" init "local-node" --chain-id "${CHAIN_ID}" --home "${NODE_HOME}" > "${NODE_HOME}/init_output.json" 2>&1; then
            log_info "Chain initialized successfully"
            log_info "Adding validator to genesis..."
            add_validator_to_genesis
        else
            log_error "Chain initialization failed or timed out"
            if [ -f "${NODE_HOME}/init_output.json" ]; then
                log_error "Init output: $(cat "${NODE_HOME}/init_output.json" | head -10)"
            fi
            exit 1
        fi
    fi
}

# Create genesis manually from config.yml
# This bypasses Ignite's keyring issues by using the binary directly
create_genesis_from_config() {
    log_info "Creating genesis manually from config.yml..."
    
    # Initialize chain first
    log_info "Initializing chain with bitbadgeschaind..."
    rm -rf "${NODE_HOME}"
    mkdir -p "${NODE_HOME}/config"
    if ! "${BINARY_PATH}" init "local-node" --chain-id "${CHAIN_ID}" --home "${NODE_HOME}" > "${NODE_HOME}/init_output.json" 2>&1; then
        log_error "Chain initialization failed"
        if [ -f "${NODE_HOME}/init_output.json" ]; then
            log_error "Init output: $(cat "${NODE_HOME}/init_output.json" | head -10)"
        fi
        exit 1
    fi
    
    # Clean keyring to avoid migration issues
    rm -rf "${NODE_HOME}/keyring-test"
    
    # Get account names from config.yml
    local account_names
    account_names=$(grep -A 20 "^accounts:" "${SCRIPT_DIR}/config.yml" | grep "name:" | awk '{print $3}' | tr -d '"' | grep -v "^$" || echo "")
    
    # Find validator name from config.yml
    local validator_name
    validator_name=$(grep -A 2 "validators:" "${SCRIPT_DIR}/config.yml" | grep "name:" | awk '{print $3}' | tr -d '"' | head -1)
    
    if [ -z "${validator_name}" ]; then
        # Default to alice if no validator found
        validator_name="alice"
    fi
    
    log_info "Using ${validator_name} as validator"
    
    # Create validator account first (this is the most important one)
    # Use standard secp256k1 for validator (not eth_secp256k1) to avoid keyring issues
    log_info "Creating validator account: ${validator_name}"
    local key_output
    key_output=$(echo "test1234" | "${BINARY_PATH}" keys add "${validator_name}" --keyring-backend test --home "${NODE_HOME}" --algo secp256k1 --no-backup 2>&1)
    
    if [ $? -ne 0 ]; then
        log_error "Failed to create validator account ${validator_name}"
        return 1
    fi
    
    # Extract address from key creation output (avoid keyring read issues)
    local validator_addr
    validator_addr=$(echo "${key_output}" | grep -o "address: [a-z0-9]*" | awk '{print $2}' | head -1)
    
    if [ -z "${validator_addr}" ]; then
        log_error "Could not extract validator address from key creation output"
        return 1
    fi
    
    if [ -z "${validator_addr}" ]; then
        log_error "Could not get validator address"
        return 1
    fi
    
    log_info "Validator address: ${validator_addr}"
    
    # Get bonded amount from config.yml
    local bonded_amount
    bonded_amount=$(grep -A 2 "name: ${validator_name}" "${SCRIPT_DIR}/config.yml" | grep "bonded:" | awk '{print $2}' || echo "1000000000000000ustake")
    
    # Get coins for validator from config.yml (simplified extraction)
    local validator_coins="1ubadge,1000000000000000ustake"
    # Try to extract from config, but use defaults if it fails
    local extracted_coins
    extracted_coins=$(awk "/name: ${validator_name}/,/^    - name:/" "${SCRIPT_DIR}/config.yml" | grep "^- " | awk '{print $2}' | tr '\n' ',' | sed 's/,$//' 2>/dev/null || echo "")
    
    if [ -n "${extracted_coins}" ] && [ "${extracted_coins}" != "" ]; then
        validator_coins="${extracted_coins}"
    fi
    
    log_info "Adding validator ${validator_name} to genesis with coins: ${validator_coins}"
    
    # Add validator to genesis
    "${BINARY_PATH}" genesis add-genesis-account "${validator_addr}" "${validator_coins}" --keyring-backend test --home "${NODE_HOME}" > /dev/null 2>&1
    
    # Skip creating other accounts for now due to keyring issues
    # They can be added later if needed
    log_info "Skipping other accounts due to keyring compatibility issues"
    
    # Create validator gentx
    log_info "Creating validator gentx for: ${validator_name}"
    
    if ! (echo "test1234" | "${BINARY_PATH}" genesis gentx "${validator_name}" "${bonded_amount}" --chain-id "${CHAIN_ID}" --keyring-backend test --home "${NODE_HOME}" > /dev/null 2>&1); then
        log_error "Failed to create gentx for ${validator_name}"
        return 1
    fi
    
    # Collect gentxs
    log_info "Collecting gentxs..."
    if ! "${BINARY_PATH}" genesis collect-gentxs --home "${NODE_HOME}" > /dev/null 2>&1; then
        log_error "Failed to collect gentxs"
        return 1
    fi
    
    # Fix denomination mismatch: default genesis uses "stake" but we need "ustake"
    log_info "Fixing denomination in genesis (stake -> ustake)..."
    sed -i 's/"denom": "stake"/"denom": "ustake"/g' "${NODE_HOME}/config/genesis.json"
    sed -i 's/"bond_denom": "stake"/"bond_denom": "ustake"/g' "${NODE_HOME}/config/genesis.json"
    sed -i 's/"mint_denom": "stake"/"mint_denom": "ustake"/g' "${NODE_HOME}/config/genesis.json"

    log_info "Genesis created successfully with validator ${validator_name}"
}

# Create accounts in a specific directory
create_accounts_for_ignite_in_dir() {
    local target_dir="$1"
    log_info "Creating accounts from config.yml in ${target_dir}..."
    
    # Get account names from config.yml
    local account_names
    account_names=$(grep -A 20 "^accounts:" "${SCRIPT_DIR}/config.yml" | grep "name:" | awk '{print $3}' | tr -d '"' | grep -v "^$" || echo "")
    
    if [ -z "${account_names}" ]; then
        log_warn "No accounts found in config.yml"
        return 1
    fi
    
    # Clean keyring in target directory to avoid migration issues
    rm -rf "${target_dir}/keyring-test"
    
    # Create each account in the target directory
    # Use secp256k1 algorithm to avoid keyring migration issues with eth_secp256k1
    for account_name in ${account_names}; do
        log_info "Creating account: ${account_name} in ${target_dir} (using secp256k1)..."
        if ! (echo "test1234" | "${BINARY_PATH}" keys add "${account_name}" --keyring-backend test --home "${target_dir}" --algo secp256k1 --no-backup > /dev/null 2>&1); then
            log_warn "Failed to create account ${account_name} in ${target_dir}, continuing..."
        else
            log_info "Successfully created account ${account_name} in ${target_dir}"
        fi
    done
}

# Create accounts for Ignite in the default keyring location
# Ignite looks for accounts in the default home directory (~/.bitbadgeschain)
# Use secp256k1 (not eth_secp256k1) to avoid keyring migration issues
create_accounts_for_ignite_fixed() {
    log_info "Creating accounts from config.yml for Ignite (default location)..."
    
    # Get account names from config.yml
    local account_names
    account_names=$(grep -A 20 "^accounts:" "${SCRIPT_DIR}/config.yml" | grep "name:" | awk '{print $3}' | tr -d '"' | grep -v "^$" || echo "")
    
    if [ -z "${account_names}" ]; then
        log_warn "No accounts found in config.yml"
        return 1
    fi
    
    # Clean default keyring to avoid migration issues
    rm -rf "${HOME}/.bitbadgeschain/keyring-test"
    
    # Create each account in the DEFAULT keyring location (where Ignite expects them)
    # Use secp256k1 algorithm to avoid keyring migration issues with eth_secp256k1
    for account_name in ${account_names}; do
        log_info "Creating account: ${account_name} (using secp256k1 for Ignite compatibility)..."
        # Create account with secp256k1 (not eth_secp256k1) in default location
        if ! (echo "test1234" | "${BINARY_PATH}" keys add "${account_name}" --keyring-backend test --algo secp256k1 --no-backup > /dev/null 2>&1); then
            log_warn "Failed to create account ${account_name}, continuing..."
        else
            log_info "Successfully created account ${account_name}"
        fi
    done
    
    # Verify accounts were created
    local account_count
    account_count=$("${BINARY_PATH}" keys list --keyring-backend test 2>/dev/null | grep -c "name:" || echo "0")
    log_info "Created ${account_count} account(s) in default keyring"
}

# Create accounts for Ignite in the keyring (deprecated - kept for reference)
# Ignite needs accounts to exist before it can initialize the chain
# Ignite uses the chain's home directory for keyring, which is ~/.ignite/chains/${CHAIN_ID}/
create_accounts_for_ignite() {
    # This function is deprecated - use create_accounts_for_ignite_fixed instead
    create_accounts_for_ignite_fixed
}

# Add a validator to genesis if it doesn't exist
add_validator_to_genesis() {
    local genesis_file="${NODE_HOME}/config/genesis.json"
    
    if [ ! -f "${genesis_file}" ]; then
        log_error "Genesis file not found: ${genesis_file}"
        return 1
    fi
    
    # Check if validators array is non-empty (has at least one validator)
    # Look for validator addresses in the validators array
    if grep -A 20 '"validators":\s*\[' "${genesis_file}" | grep -q '"address"'; then
        log_info "Validators already exist in genesis"
        return 0
    fi
    
    # Create a validator key - use --no-backup to avoid keyring issues
    local validator_key="validator"
    log_info "Creating validator key..."
    
    # Remove old keyring to avoid conflicts
    rm -rf "${NODE_HOME}/keyring-test"
    
    # Create key with password (use yes for confirmation prompt)
    if ! (yes "test1234" | head -2 | "${BINARY_PATH}" keys add "${validator_key}" --keyring-backend test --home "${NODE_HOME}" --no-backup > /dev/null 2>&1); then
        log_error "Failed to create validator key"
        log_warn "You may need to manually add a validator to genesis"
        return 1
    fi
    
    # Get validator address
    local val_addr
    val_addr=$("${BINARY_PATH}" keys show "${validator_key}" -a --keyring-backend test --home "${NODE_HOME}" 2>/dev/null)
    
    if [ -z "${val_addr}" ]; then
        log_error "Could not get validator address"
        return 1
    fi
    
    log_info "Adding validator ${val_addr} to genesis..."
    
    # Add genesis account with coins
    if ! "${BINARY_PATH}" genesis add-genesis-account "${val_addr}" "1000000000000ustake,1000000000000ubadge" --keyring-backend test --home "${NODE_HOME}" > /dev/null 2>&1; then
        log_error "Failed to add genesis account"
        return 1
    fi
    
    # Create validator gentx
    log_info "Creating validator gentx..."
    if ! (yes "test1234" | head -1 | "${BINARY_PATH}" genesis gentx "${validator_key}" "1000000000000ustake" --chain-id "${CHAIN_ID}" --keyring-backend test --home "${NODE_HOME}" > /dev/null 2>&1); then
        log_error "Failed to create gentx"
        return 1
    fi
    
    # Collect gentxs
    log_info "Collecting gentxs..."
    if ! "${BINARY_PATH}" genesis collect-gentxs --home "${NODE_HOME}" > /dev/null 2>&1; then
        log_error "Failed to collect gentxs"
        return 1
    fi
    
    log_info "Validator added to genesis successfully"
}

# Start the chain
start_chain() {
    log_info "Starting chain daemon..."
    
    # Check if chain is already running
    if pgrep -f "${BINARY_NAME} start" > /dev/null; then
        if [ -z "${SKIP_PROMPTS:-}" ]; then
            log_warn "Chain appears to be already running"
            read -p "Kill existing process and restart? (y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                log_info "Stopping existing chain..."
                pkill -f "${BINARY_NAME} start" || true
                sleep 2
            else
                log_info "Exiting, chain is already running"
                return 0
            fi
        else
            log_info "Chain already running, stopping and restarting..."
            pkill -f "${BINARY_NAME} start" || true
            sleep 2
        fi
    fi
    
    # Start the chain in the background
    log_info "Starting ${BINARY_NAME} with chain-id: ${CHAIN_ID}"
    log_info "Node home: ${NODE_HOME}"
    log_info "RPC endpoint will be at: http://localhost:26657"
    log_info "API endpoint will be at: http://localhost:1317"
    
    # Start the daemon with JSON-RPC enabled for EVM compatibility
    "${BINARY_PATH}" start --json-rpc.enable --json-rpc.address 0.0.0.0:8545 --home "${NODE_HOME}" > "${NODE_HOME}/node.log" 2>&1 &
    DAEMON_PID=$!
    
    log_info "Chain daemon started with PID: ${DAEMON_PID}"
    echo "${DAEMON_PID}" > "${NODE_HOME}/daemon.pid"
    
    # Wait a bit for the chain to start
    sleep 5
    
    # Check if process is still running
    if ! kill -0 "${DAEMON_PID}" 2>/dev/null; then
        log_error "Chain daemon process died immediately"
        log_error "Check logs at: ${NODE_HOME}/node.log"
        tail -50 "${NODE_HOME}/node.log"
        exit 1
    fi
    
    log_info "Waiting for chain to be ready (checking RPC endpoint)..."
    
    # Wait for RPC to be available
    local elapsed=0
    while [ ${elapsed} -lt ${START_TIMEOUT} ]; do
        if curl -s http://localhost:26657/status > /dev/null 2>&1; then
            log_info "RPC endpoint is responding"
            break
        fi
        sleep 1
        elapsed=$((elapsed + 1))
    done
    
    if [ ${elapsed} -ge ${START_TIMEOUT} ]; then
        log_error "RPC endpoint did not become available after ${START_TIMEOUT}s"
        log_error "Check logs at: ${NODE_HOME}/node.log"
        tail -50 "${NODE_HOME}/node.log"
        kill "${DAEMON_PID}" 2>/dev/null || true
        exit 1
    fi
    
    # Wait for first block
    log_info "Waiting for first block to be produced..."
    elapsed=0
    while [ ${elapsed} -lt ${BLOCK_WAIT_TIMEOUT} ]; do
        local latest_block=$(curl -s http://localhost:26657/status 2>/dev/null | grep -o '"latest_block_height":"[0-9]*"' | grep -o '[0-9]*' | head -1 || echo "0")
        if [ -n "${latest_block}" ] && [ "${latest_block}" != "0" ] && [ "${latest_block}" != "null" ]; then
            log_info "Chain is producing blocks! Latest block height: ${latest_block}"
            break
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done
    
    if [ ${elapsed} -ge ${BLOCK_WAIT_TIMEOUT} ]; then
        log_warn "No blocks produced after ${BLOCK_WAIT_TIMEOUT}s, but chain is running"
        log_warn "This might be normal if the chain is still initializing"
    fi
    
    # Display chain status
    log_info "Chain status:"
    curl -s http://localhost:26657/status | grep -E '"latest_block_height"|"chain_id"|"catching_up"' || true
    
    log_info ""
    log_info "Chain is running!"
    log_info "PID: ${DAEMON_PID}"
    log_info "Logs: ${NODE_HOME}/node.log"
    log_info "RPC: http://localhost:26657"
    log_info "API: http://localhost:1317"
    log_info ""
    log_info "To stop the chain, run: kill ${DAEMON_PID}"
    log_info "Or use: pkill -f '${BINARY_NAME} start'"
}

# Stop the chain
stop_chain() {
    log_info "Stopping chain..."
    
    if [ -f "${NODE_HOME}/daemon.pid" ]; then
        local pid=$(cat "${NODE_HOME}/daemon.pid")
        if kill -0 "${pid}" 2>/dev/null; then
            kill "${pid}" 2>/dev/null || true
            log_info "Stopped chain daemon (PID: ${pid})"
        else
            log_warn "Process ${pid} not running"
        fi
        rm -f "${NODE_HOME}/daemon.pid"
    else
        # Try to find and kill by process name
        if pkill -f "${BINARY_NAME} start"; then
            log_info "Stopped chain daemon"
        else
            log_warn "No running chain daemon found"
        fi
    fi
}

# Main execution
main() {
    local command="${1:-start}"
    local non_interactive="${2:-}"
    
    # If CI or non-interactive flag, skip prompts
    if [ -n "${CI:-}" ] || [ "${non_interactive}" = "--non-interactive" ] || [ "${non_interactive}" = "-y" ]; then
        export SKIP_PROMPTS=1
    fi
    
    case "${command}" in
        start)
            check_prerequisites
            build_binary
            init_chain
            start_chain
            ;;
        stop)
            stop_chain
            ;;
        restart)
            stop_chain
            sleep 2
            check_prerequisites
            build_binary
            init_chain
            start_chain
            ;;
        build)
            check_prerequisites
            build_binary
            ;;
        init)
            check_prerequisites
            build_binary
            init_chain
            ;;
        status)
            if curl -s http://localhost:26657/status > /dev/null 2>&1; then
                log_info "Chain is running"
                curl -s http://localhost:26657/status | grep -E '"latest_block_height"|"chain_id"|"catching_up"' || true
            else
                log_error "Chain is not running or RPC is not accessible"
            fi
            ;;
        *)
            echo "Usage: $0 {start|stop|restart|build|init|status}"
            echo ""
            echo "Commands:"
            echo "  start   - Build, initialize (if needed), and start the chain"
            echo "  stop    - Stop the running chain"
            echo "  restart - Stop, rebuild, reinitialize, and start the chain"
            echo "  build   - Build the binary only"
            echo "  init    - Initialize the chain only"
            echo "  status  - Check if chain is running and show status"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"

