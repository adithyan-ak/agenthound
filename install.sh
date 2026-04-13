#!/bin/sh
set -e

REPO="ghcr.io/adithyan-ak/agenthound"
GITHUB_REPO="adithyan-ak/agenthound"
TAG="latest"
NAME="agenthound"
PORT="${AGENTHOUND_PORT:-8080}"
INSTALL_DIR="${AGENTHOUND_INSTALL_DIR:-/usr/local/bin}"
ADMIN_PASSWORD="${AGENTHOUND_ADMIN_PASSWORD:-agenthound}"

echo ""
echo "  AgentHound Installer"
echo "  ===================="
echo ""

# ── Detect container runtime ──
if command -v docker >/dev/null 2>&1; then
  RUNTIME="docker"
elif command -v podman >/dev/null 2>&1; then
  RUNTIME="podman"
else
  echo "Error: Docker or Podman is required."
  echo ""
  echo "Install Docker: https://docs.docker.com/get-docker/"
  exit 1
fi

$RUNTIME info >/dev/null 2>&1 || {
  echo "Error: $RUNTIME daemon is not running."
  echo "Start it and try again."
  exit 1
}

# ── Detect OS and architecture ──
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Warning: Unsupported architecture $ARCH. CLI binary will not be installed." ;;
esac

# ── Check port availability ──
if command -v lsof >/dev/null 2>&1 && lsof -i ":$PORT" >/dev/null 2>&1; then
  echo "Error: Port $PORT is already in use."
  echo ""
  echo "Use a different port:"
  echo "  AGENTHOUND_PORT=9090 sh -c \"\$(curl -sSfL https://raw.githubusercontent.com/${GITHUB_REPO}/main/install.sh)\""
  exit 1
fi

# ── Stop existing container if running ──
if $RUNTIME ps -a --format '{{.Names}}' 2>/dev/null | grep -q "^${NAME}$"; then
  echo "Stopping existing AgentHound container..."
  $RUNTIME stop "$NAME" 2>/dev/null || true
  $RUNTIME rm "$NAME" 2>/dev/null || true
fi

# ── Pull and run server ──
echo "Pulling AgentHound image..."
$RUNTIME pull "$REPO:$TAG"

echo "Starting AgentHound on port $PORT..."
$RUNTIME run -d \
  --name "$NAME" \
  -p "${PORT}:8080" \
  -v agenthound-data:/data \
  --restart unless-stopped \
  "$REPO:$TAG"

# ── Wait for healthy ──
echo ""
echo "Waiting for services to start (this takes ~30s)..."
ATTEMPTS=0
while [ $ATTEMPTS -lt 45 ]; do
  if curl -sf "http://localhost:${PORT}/api/v1/health" >/dev/null 2>&1; then
    break
  fi
  ATTEMPTS=$((ATTEMPTS + 1))
  sleep 2
  printf "."
done

if [ $ATTEMPTS -ge 45 ]; then
  echo ""
  echo ""
  echo "Warning: Health check timed out after 90s."
  echo "Check logs: $RUNTIME logs $NAME"
  exit 1
fi

echo ""
echo "Server is running."

# ── Download native CLI binary ──
CLI_INSTALLED=false
if [ "$ARCH" = "amd64" ] || [ "$ARCH" = "arm64" ]; then
  VERSION=$(curl -sSfL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" 2>/dev/null | grep '"tag_name"' | cut -d'"' -f4 || true)
  if [ -z "$VERSION" ]; then
    echo ""
    echo "Warning: No release found. CLI binary not installed."
    echo "You can still use: $RUNTIME exec $NAME agenthound scan"
  else
    BINARY_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/agenthound-${VERSION#v}-${OS}-${ARCH}.tar.gz"
    echo "Downloading CLI binary ${VERSION} for ${OS}/${ARCH}..."
    if curl -sSfL "$BINARY_URL" | tar -xz -C /tmp 2>/dev/null; then
      if install -m 755 /tmp/agenthound "$INSTALL_DIR/agenthound" 2>/dev/null; then
        CLI_INSTALLED=true
        echo "Installed to ${INSTALL_DIR}/agenthound"
      else
        echo "Warning: Could not install to ${INSTALL_DIR}. Try with sudo or set AGENTHOUND_INSTALL_DIR."
      fi
      rm -f /tmp/agenthound
    else
      echo "Warning: Binary download failed. You can still use: $RUNTIME exec $NAME agenthound scan"
    fi
  fi
fi

# ── Bootstrap CLI auth ──
if [ "$CLI_INSTALLED" = true ]; then
  echo ""
  echo "Setting up CLI authentication..."
  JWT=$(curl -sSf "http://localhost:${PORT}/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"admin\",\"password\":\"${ADMIN_PASSWORD}\"}" 2>/dev/null \
    | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4 || true)

  if [ -n "$JWT" ]; then
    HOSTNAME=$(hostname 2>/dev/null || echo "unknown")
    TOKEN=$(curl -sSf "http://localhost:${PORT}/api/v1/auth/tokens" \
      -H "Authorization: Bearer $JWT" \
      -H "Content-Type: application/json" \
      -d "{\"name\":\"cli-${HOSTNAME}-install\"}" 2>/dev/null \
      | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4 || true)

    if [ -n "$TOKEN" ]; then
      CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/agenthound"
      mkdir -p "$CONFIG_DIR"
      cat > "$CONFIG_DIR/config.yaml" <<YAML
server_url: http://localhost:${PORT}
api_token: ${TOKEN}
YAML
      chmod 600 "$CONFIG_DIR/config.yaml"
      echo "CLI configured."
    else
      echo "Warning: Could not create API token. Run 'agenthound setup' manually."
    fi
  else
    echo "Warning: Could not authenticate. Run 'agenthound setup' manually."
  fi
fi

# ── Print summary ──
echo ""
echo "  AgentHound is running!"
echo ""
echo "  URL:     http://localhost:${PORT}"
echo "  Login:   admin / agenthound"
echo ""
if [ "$CLI_INSTALLED" = true ]; then
  echo "  Scan:    agenthound scan"
  echo "  Query:   agenthound query --findings"
else
  echo "  Scan:    $RUNTIME exec $NAME agenthound scan"
fi
echo ""
echo "  Manage:"
echo "    Stop:    $RUNTIME stop $NAME"
echo "    Start:   $RUNTIME start $NAME"
echo "    Logs:    $RUNTIME logs -f $NAME"
echo "    Remove:  $RUNTIME rm -f $NAME && $RUNTIME volume rm agenthound-data"
echo ""
