#!/bin/bash
set -e

BINARY_NAME="aha-mcp-server"
CLAUDE_DESKTOP_CONFIG="$HOME/Library/Application Support/Claude/claude_desktop_config.json"
CLAUDE_CODE_CONFIG="$HOME/.claude.json"

echo "Aha! MCP Server Installer"
echo "========================="
echo ""

# --- 1. Check for Go ---

if ! command -v go &>/dev/null; then
  echo "Go is not installed. Installing via Homebrew..."
  if ! command -v brew &>/dev/null; then
    echo "Error: Homebrew is also not installed."
    echo "Install Homebrew first: https://brew.sh"
    exit 1
  fi
  brew install go
fi

echo "Go found: $(go version)"
echo ""

# --- 2. Build and install ---

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
echo "Building $BINARY_NAME from source..."
cd "$SCRIPT_DIR"
go install ./cmd/aha-mcp-server
echo "Done."
echo ""

INSTALL_PATH="$(go env GOPATH)/bin/$BINARY_NAME"

# --- 3. Get API token ---

echo "You need an Aha! API token."
echo "Get one at: https://jamf.aha.io/settings/api_keys"
echo ""
read -rsp "Paste your Aha! API token: " AHA_TOKEN
echo ""

if [ -z "$AHA_TOKEN" ]; then
  echo "Error: API token cannot be empty."
  exit 1
fi
echo ""

# --- 4. Configure Claude ---

configure_mcp() {
  local config_path="$1"
  local label="$2"

  python3 - <<EOF
import json, os

config_path = os.path.expanduser("$config_path")
os.makedirs(os.path.dirname(config_path), exist_ok=True)

if os.path.exists(config_path):
    with open(config_path) as f:
        try:
            config = json.load(f)
        except json.JSONDecodeError:
            config = {}
else:
    config = {}

if "mcpServers" not in config:
    config["mcpServers"] = {}

config["mcpServers"]["aha"] = {
    "command": "$INSTALL_PATH",
    "env": {
        "AHA_API_TOKEN": "$AHA_TOKEN",
        "AHA_DOMAIN": "jamf",
        "GODEBUG": "netdns=go"
    }
}

with open(config_path, "w") as f:
    json.dump(config, f, indent=2)
EOF
  echo "  Configured $label"
}

echo "Configuring Claude..."

if [ -d "$HOME/Library/Application Support/Claude" ] || [ ! -f "$CLAUDE_DESKTOP_CONFIG" ]; then
  configure_mcp "$CLAUDE_DESKTOP_CONFIG" "Claude Desktop"
fi

if [ -f "$CLAUDE_CODE_CONFIG" ]; then
  configure_mcp "$CLAUDE_CODE_CONFIG" "Claude Code"
fi

echo ""
echo "All done!"
echo ""
echo "Next steps:"
echo "  - Restart Claude Desktop if it's running"
echo "  - Restart Claude Code if it's running"
echo "  - You'll have 16 Aha! tools available in both apps"
