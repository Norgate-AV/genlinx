#!/bin/bash
# Package Go binaries into npm platform packages and publish to the registry.
#
# Environment controls:
#   RUN_GORELEASER=false   Skip goreleaser build and reuse an existing dist/
#   MODE=pack              pack (default) | dry-run | publish
#   REGISTRY=<url>         Point at a local Verdaccio instance, e.g. http://localhost:4873
#
# Windows note: run via Git Bash, not WSL bash. WSL bash cannot write to Windows
# paths via Windows executables (go.exe). From PowerShell:
#   & "C:\Program Files\Git\bin\bash.exe" ./scripts/npm-publish.sh
#
# First run cross-compiles the Go stdlib for all targets — this takes a few
# minutes but the results are cached by the Go toolchain for subsequent runs.

set -e

# Scoop-installed tools (npm, go, etc.) refuse to run under WSL 1.
# Detect it early and direct the user to Git Bash instead.
if grep -qi "microsoft" /proc/version 2>/dev/null \
    && ! grep -qi "wsl2" /proc/sys/kernel/osrelease 2>/dev/null; then
    echo "Error: WSL 1 detected. Scoop tools do not work under WSL 1."
    echo "  Run with Git Bash instead:"
    echo "    & \"C:\\Program Files\\Git\\bin\\bash.exe\" ./scripts/test-npm-publish.sh"
    exit 1
fi

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST_DIR="$REPO_ROOT/dist"
TMPL="$REPO_ROOT/npm/package.json.tmpl"
WRAPPER_SRC="$REPO_ROOT/npm/genlinx"

RUN_GORELEASER="${RUN_GORELEASER:-true}"
MODE="${MODE:-pack}"
REGISTRY="${REGISTRY:-}"

# ---------------------------------------------------------------------------
# Tool discovery — handles WSL (where Scoop shims are .exe files that bash
# won't find without the extension) and native Linux/macOS installs.
# ---------------------------------------------------------------------------

# find_cmd NAME — resolves NAME or NAME.exe; returns the usable command string.
find_cmd() {
    if command -v "$1" &>/dev/null;     then echo "$1";     return; fi
    if command -v "${1}.exe" &>/dev/null; then echo "${1}.exe"; return; fi
    echo ""
}

GO_CMD=$(find_cmd go)
GORELEASER_CMD=$(find_cmd goreleaser)
NPM_CMD=$(find_cmd npm)

if [ -z "$GO_CMD" ]; then
    echo "Error: 'go' not found in PATH. Install Go or run from a shell where it is accessible."
    exit 1
fi
if [ -z "$NPM_CMD" ]; then
    echo "Error: 'npm' not found in PATH."
    exit 1
fi

# json_get FILE KEY — extract a top-level JSON string value without external tools.
json_get() {
    grep -o "\"$2\"[[:space:]]*:[[:space:]]*\"[^\"]*\"" "$1" \
        | head -1 \
        | sed 's/.*"[^"]*"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/'
}

# json_stamp_version FILE VERSION — replace all placeholder "0.0.0" version
# strings in the wrapper package.json with the real version.
json_stamp_version() {
    local tmp; tmp=$(mktemp)
    sed "s/\"0\.0\.0\"/\"${2}\"/g" "$1" > "$tmp" && mv "$tmp" "$1"
}

# ---------------------------------------------------------------------------
# 1. Build binaries
# ---------------------------------------------------------------------------

build_with_go() {
    local version="$1"
    echo "==> goreleaser not found — building all platforms with $GO_CMD..."
    cd "$REPO_ROOT"

    # Matches the dist dir layout goreleaser produces so the rest of the script
    # works unchanged regardless of which builder was used.
    declare -A GO_PLATFORMS=(
        ["linux_amd64_v1"]="linux amd64 genlinx"
        ["linux_arm64"]="linux arm64 genlinx"
        ["darwin_amd64_v1"]="darwin amd64 genlinx"
        ["darwin_arm64"]="darwin arm64 genlinx"
        ["windows_amd64_v1"]="windows amd64 genlinx.exe"
        ["windows_arm64"]="windows arm64 genlinx.exe"
    )

    for suffix in "${!GO_PLATFORMS[@]}"; do
        read -r goos goarch binary <<< "${GO_PLATFORMS[$suffix]}"
        abs_out="$DIST_DIR/genlinx_${suffix}/$binary"
        # Use a relative path for -o so Windows go.exe can resolve it via
        # the translated CWD — absolute WSL paths (/mnt/c/...) are opaque to
        # Windows executables.
        rel_out="dist/genlinx_${suffix}/$binary"
        mkdir -p "$(dirname "$abs_out")"
        printf "     %-30s" "$goos/$goarch"
        GOOS=$goos GOARCH=$goarch CGO_ENABLED=0 $GO_CMD build -trimpath -o "$rel_out" . \
            && echo "[ok]" || { echo "[FAILED]"; exit 1; }
    done

    mkdir -p "$DIST_DIR"
    printf '{"version":"%s"}\n' "$version" > "$DIST_DIR/metadata.json"
}

if [ "$RUN_GORELEASER" = "true" ]; then
    if [ -n "$GORELEASER_CMD" ]; then
        echo "==> Building snapshot binaries with goreleaser..."
        cd "$REPO_ROOT"
        $GORELEASER_CMD release --snapshot --clean
    else
        RAW_VERSION=$(git -C "$REPO_ROOT" describe --tags --always --dirty 2>/dev/null || echo "0.0.0")
        build_with_go "${RAW_VERSION#v}-dev"
    fi
fi

[ -d "$DIST_DIR" ] || {
    echo "Error: dist/ not found. Run with RUN_GORELEASER=true or build first."
    exit 1
}
[ -f "$DIST_DIR/metadata.json" ] || {
    echo "Error: dist/metadata.json not found."
    exit 1
}

VERSION=$(json_get "$DIST_DIR/metadata.json" version)

echo ""
echo "  Version  : $VERSION"
echo "  Mode     : $MODE"
[ -n "$REGISTRY" ] && echo "  Registry : $REGISTRY"
echo ""

# ---------------------------------------------------------------------------
# 2. Stage packages into a temp working dir
# ---------------------------------------------------------------------------

WORK_DIR="$(mktemp -d)"
trap 'rm -rf "$WORK_DIR"' EXIT

# GoReleaser default dist dir naming:
#   amd64 → genlinx_{goos}_amd64_v1/{binary}
#   arm64 → genlinx_{goos}_arm64/{binary}
#
# Columns: dist_dir_suffix | node_os | node_arch | binary_name
PLATFORMS=(
    "linux_amd64_v1|linux|x64|genlinx"
    "linux_arm64|linux|arm64|genlinx"
    "darwin_amd64_v1|darwin|x64|genlinx"
    "darwin_arm64|darwin|arm64|genlinx"
    "windows_amd64_v1|windows|x64|genlinx.exe"
    "windows_arm64|windows|arm64|genlinx.exe"
)

ERRORS=0

for entry in "${PLATFORMS[@]}"; do
    IFS='|' read -r dir_suffix node_os node_arch binary <<< "$entry"
    node_pkg="@norgate-av/genlinx-${node_os}-${node_arch}"
    binary_src="$DIST_DIR/genlinx_${dir_suffix}/$binary"

    printf "  %-42s" "$node_pkg"

    if [ ! -f "$binary_src" ]; then
        echo "[MISSING]  $binary_src"
        ERRORS=$((ERRORS + 1))
        continue
    fi

    pkg_dir="$WORK_DIR/$node_pkg"
    mkdir -p "$pkg_dir"
    cp "$binary_src" "$pkg_dir/$binary"

    # Substitute template placeholders using sed (no envsubst dependency)
    sed \
        -e "s|\${node_pkg}|${node_pkg}|g" \
        -e "s|\${version}|${VERSION}|g" \
        -e "s|\${node_os}|${node_os}|g" \
        -e "s|\${node_arch}|${node_arch}|g" \
        -e "s|\${binary}|${binary}|g" \
        "$TMPL" > "$pkg_dir/package.json"

    echo "[ok]"
done

# Main wrapper
wrapper_dir="$WORK_DIR/@norgate-av_genlinx"
cp -r "$WRAPPER_SRC" "$wrapper_dir"
json_stamp_version "$wrapper_dir/package.json" "$VERSION"
printf "  %-42s" "@norgate-av/genlinx"
echo "[ok]"

if [ "$ERRORS" -gt 0 ]; then
    echo ""
    echo "==> $ERRORS platform(s) are missing binaries. Aborting."
    exit 1
fi

# ---------------------------------------------------------------------------
# 3. Run npm for each staged package
# ---------------------------------------------------------------------------

echo ""
echo "==> Running npm $MODE for each package..."
echo ""

npm_args=("--access" "public")
[ -n "$REGISTRY" ] && npm_args+=("--registry" "$REGISTRY")
[ "$MODE" = "dry-run" ] && npm_args+=("--dry-run")
# Pre-release versions must not overwrite the 'latest' dist-tag on npm
[[ "$VERSION" == *"-"* ]] && npm_args+=("--tag" "next")

for entry in "${PLATFORMS[@]}"; do
    IFS='|' read -r dir_suffix node_os node_arch binary <<< "$entry"
    node_pkg="@norgate-av/genlinx-${node_os}-${node_arch}"
    pkg_dir="$WORK_DIR/$node_pkg"

    echo "--- $node_pkg ---"
    cd "$pkg_dir"
    case "$MODE" in
        pack)     $NPM_CMD pack ;;
        dry-run)  $NPM_CMD publish "${npm_args[@]}" ;;
        publish)  $NPM_CMD publish "${npm_args[@]}" ;;
    esac
    echo ""
    cd "$REPO_ROOT"
done

echo "--- @norgate-av/genlinx ---"
cd "$wrapper_dir"
case "$MODE" in
    pack)     $NPM_CMD pack ;;
    dry-run)  $NPM_CMD publish "${npm_args[@]}" ;;
    publish)  $NPM_CMD publish "${npm_args[@]}" ;;
esac
echo ""
cd "$REPO_ROOT"

echo "==> Done."
