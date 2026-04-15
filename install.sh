#!/bin/sh
# BitBadges Chain — One-liner installer
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/BitBadges/bitbadgeschain/master/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/BitBadges/bitbadgeschain/master/install.sh | sh -s -- --testnet
#   curl -fsSL https://raw.githubusercontent.com/BitBadges/bitbadgeschain/master/install.sh | sh -s -- --version v28
#
# Options:
#   --testnet           Install testnet binary (default: mainnet)
#   --version <tag>     Install a specific release (default: latest)
#   --install-dir <dir> Install to a custom directory (default: /usr/local/bin)
#   --no-sudo           Don't use sudo even if install dir requires it

set -e

REPO="BitBadges/bitbadgeschain"
NETWORK="mainnet"
VERSION=""
INSTALL_DIR="/usr/local/bin"
USE_SUDO="auto"
BINARY_NAME="bitbadgeschaind"

# Parse arguments
while [ $# -gt 0 ]; do
  case "$1" in
    --testnet)
      NETWORK="testnet"
      shift
      ;;
    --version)
      VERSION="$2"
      shift 2
      ;;
    --install-dir)
      INSTALL_DIR="$2"
      shift 2
      ;;
    --no-sudo)
      USE_SUDO="no"
      shift
      ;;
    --help)
      echo "BitBadges Chain Installer"
      echo ""
      echo "Usage: curl -fsSL https://raw.githubusercontent.com/BitBadges/bitbadgeschain/master/install.sh | sh"
      echo ""
      echo "Options:"
      echo "  --testnet           Install testnet binary (default: mainnet)"
      echo "  --version <tag>     Install a specific release (default: latest)"
      echo "  --install-dir <dir> Install directory (default: /usr/local/bin)"
      echo "  --no-sudo           Don't use sudo"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Detect OS
detect_os() {
  case "$(uname -s)" in
    Linux*)   echo "linux" ;;
    Darwin*)  echo "darwin" ;;
    MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
    *)
      echo "Unsupported OS: $(uname -s)" >&2
      exit 1
      ;;
  esac
}

# Detect architecture
detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64)   echo "amd64" ;;
    aarch64|arm64)   echo "arm64" ;;
    *)
      echo "Unsupported architecture: $(uname -m)" >&2
      exit 1
      ;;
  esac
}

# Get latest release tag from GitHub API
get_latest_version() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p'
  elif command -v wget >/dev/null 2>&1; then
    wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p'
  else
    echo "Error: curl or wget required" >&2
    exit 1
  fi
}

# Download a file
download() {
  local url="$1"
  local output="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL -o "$output" "$url"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$output" "$url"
  fi
}

# Main
main() {
  local os arch asset_name download_url tmp_dir

  os="$(detect_os)"
  arch="$(detect_arch)"

  echo "Detected platform: ${os}/${arch}"

  # Windows only supports amd64
  if [ "$os" = "windows" ] && [ "$arch" != "amd64" ]; then
    echo "Error: Windows builds are only available for amd64" >&2
    exit 1
  fi

  # Darwin only supports amd64 and arm64
  if [ "$os" = "darwin" ] && [ "$arch" != "amd64" ] && [ "$arch" != "arm64" ]; then
    echo "Error: macOS builds are available for amd64 and arm64 only" >&2
    exit 1
  fi

  # Resolve version
  if [ -z "$VERSION" ]; then
    echo "Fetching latest release..."
    VERSION="$(get_latest_version)"
    if [ -z "$VERSION" ]; then
      echo "Error: could not determine latest release version" >&2
      exit 1
    fi
  fi
  echo "Version: ${VERSION}"

  # Build asset name
  # Mainnet: bitbadgeschain-{os}-{arch}
  # Testnet: bitbadgeschain-testnet-{os}-{arch}
  if [ "$NETWORK" = "testnet" ]; then
    asset_name="bitbadgeschain-testnet-${os}-${arch}"
  else
    asset_name="bitbadgeschain-${os}-${arch}"
  fi

  # Windows binaries have .exe extension
  if [ "$os" = "windows" ]; then
    asset_name="${asset_name}.exe"
    BINARY_NAME="bitbadgeschaind.exe"
  fi

  download_url="https://github.com/${REPO}/releases/download/${VERSION}/${asset_name}"
  echo "Downloading ${asset_name}..."

  # Create temp dir
  tmp_dir="$(mktemp -d)"
  trap 'rm -rf "$tmp_dir"' EXIT

  download "$download_url" "${tmp_dir}/${asset_name}"

  if [ ! -f "${tmp_dir}/${asset_name}" ]; then
    echo "Error: download failed" >&2
    echo ""
    echo "The binary for your platform (${os}/${arch}) may not be available yet for ${VERSION}."
    echo "Check available assets at: https://github.com/${REPO}/releases/tag/${VERSION}"
    exit 1
  fi

  # Make executable (not needed on Windows but doesn't hurt)
  chmod +x "${tmp_dir}/${asset_name}"

  # Install
  echo "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
  mkdir -p "$INSTALL_DIR" 2>/dev/null || true

  if [ -w "$INSTALL_DIR" ] || [ "$USE_SUDO" = "no" ]; then
    mv "${tmp_dir}/${asset_name}" "${INSTALL_DIR}/${BINARY_NAME}"
  elif [ "$USE_SUDO" = "auto" ] && command -v sudo >/dev/null 2>&1; then
    sudo mv "${tmp_dir}/${asset_name}" "${INSTALL_DIR}/${BINARY_NAME}"
  else
    echo "Error: ${INSTALL_DIR} is not writable. Run with sudo or use --install-dir." >&2
    exit 1
  fi

  echo ""
  echo "Successfully installed bitbadgeschaind ${VERSION} (${NETWORK}) to ${INSTALL_DIR}/${BINARY_NAME}"
  echo ""

  # Verify
  if command -v "${BINARY_NAME}" >/dev/null 2>&1; then
    echo "Verify: $(${BINARY_NAME} version 2>/dev/null || echo 'installed')"
  else
    echo "Note: ${INSTALL_DIR} may not be in your PATH. Add it with:"
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
  fi

  # Install SDK CLI
  echo ""
  if command -v bun >/dev/null 2>&1; then
    echo "Installing BitBadges SDK CLI (bitbadges-cli) via bun..."
    bun install -g bitbadges 2>&1 | tail -1
  elif command -v npm >/dev/null 2>&1; then
    echo "Installing BitBadges SDK CLI (bitbadges-cli) via npm..."
    npm install -g bitbadges 2>&1 | tail -1
  else
    echo "npm/bun not found — skipping SDK CLI install. To install later:"
    echo "  npm install -g bitbadges"
  fi

  if command -v bitbadges-cli >/dev/null 2>&1; then
    echo "Successfully installed bitbadges-cli"
  elif command -v bun >/dev/null 2>&1 || command -v npm >/dev/null 2>&1; then
    echo "Note: bitbadges-cli installed but may not be in PATH. Check your global bin directory."
  fi

  echo ""
  echo "Quick start:"
  echo "  bitbadgeschaind init my-node --chain-id bitbadges-1"
  echo "  bitbadgeschaind start"
  echo "  bitbadges-cli --help"
}

main
