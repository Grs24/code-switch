#!/usr/bin/env bash
set -euo pipefail

# Configuration
R2_ENDPOINT="https://7098a4e60134735b01cf4c3ef96f0f3f.r2.cloudflarestorage.com"
R2_BUCKET="bewildcard"
R2_PREFIX="code-switch"
PUBLIC_URL_BASE="https://static.bewildcard.com/${R2_PREFIX}"

if ! command -v aws >/dev/null 2>&1; then
  echo "aws CLI is required. Install via 'brew install awscli'" >&2
  exit 1
fi

if [ $# -lt 1 ]; then
  echo "Usage: scripts/publish_r2.sh <tag> [notes-file]" >&2
  exit 1
fi

TAG="$1"
NOTES="${2:-RELEASE_NOTES.md}"

if [ ! -f "$NOTES" ]; then
  echo "Release notes file '$NOTES' not found" >&2
  exit 1
fi

# Ensure tag starts with v
if [[ ! "$TAG" =~ ^v ]]; then
  TAG="v$TAG"
fi

VERSION="${TAG#v}" # Remove v prefix for some uses if needed

echo "==> Preparing release $TAG"

# 1. Update version in code
perl -0pi -e "s/const\\s+AppVersion\\s*=\\s*\"[^\"]*\"/const AppVersion = \"$TAG\"/" version_service.go

# 2. Build Assets
echo "==> Building frontend assets..."
wails3 task common:update:build-assets

# 3. Build & Package macOS
MAC_APP_PRIMARY="bin/CodeSwitch.app"
MAC_ARCHS=("arm64" "amd64")
MAC_ZIPS=()

package_macos_arch() {
  local arch="$1"
  local staging_dir="bin/package-${arch}"
  local staging_app="${staging_dir}/CodeSwitch.app"
  local zip_name="CodeSwitch-macos-${arch}-${TAG}.zip"
  local zip_path="bin/${zip_name}"

  echo "==> Building macOS ${arch}..."
  env ARCH="$arch" wails3 task package ${BUILD_OPTS:-}

  if [ ! -d "$MAC_APP_PRIMARY" ]; then
    echo "Missing asset: $MAC_APP_PRIMARY" >&2
    exit 1
  fi

  rm -rf "$staging_dir"
  mkdir -p "$staging_dir"
  cp -R "$MAC_APP_PRIMARY" "$staging_app"

  echo "==> Archiving macOS app bundle (${arch})..."
  rm -f "$zip_path"
  ditto -c -k --sequesterRsrc --keepParent "$staging_app" "$zip_path"
  rm -rf "$staging_dir"

  MAC_ZIPS+=("$zip_path")
}

for arch in "${MAC_ARCHS[@]}"; do
  package_macos_arch "$arch"
done

# 4. Build Windows
echo "==> Building Windows amd64..."
env ARCH=amd64 wails3 task windows:package ${BUILD_OPTS:-}
WIN_INSTALLER="bin/codeswitch-amd64-installer.exe"
WIN_EXE="bin/codeswitch.exe"
# Rename for upload consistency
WIN_INSTALLER_UPLOAD="CodeSwitch-windows-setup-${TAG}.exe"
cp "$WIN_INSTALLER" "bin/${WIN_INSTALLER_UPLOAD}"

# 5. Upload to R2
echo "==> Uploading artifacts to R2..."

upload_file() {
  local file="$1"
  local name="$2"
  echo "  -> Uploading $name"
  aws s3 cp "$file" "s3://${R2_BUCKET}/${R2_PREFIX}/${TAG}/${name}" \
    --endpoint-url "$R2_ENDPOINT" \
    --acl public-read
}

# Upload macOS zips
for zip in "${MAC_ZIPS[@]}"; do
  filename=$(basename "$zip")
  upload_file "$zip" "$filename"
done

# Upload Windows installer
upload_file "bin/${WIN_INSTALLER_UPLOAD}" "$WIN_INSTALLER_UPLOAD"

# 6. Generate and Upload latest.json
echo "==> Updating latest.json..."

# Read release notes
BODY=$(cat "$NOTES")
# Escape json string (basic)
BODY_JSON=$(echo "$BODY" | jq -R -s '.')

# Construct JSON
# Note: We provide direct download links.
# For simplicity, we assume:
# - macOS arm64 -> CodeSwitch-macos-arm64-${TAG}.zip
# - macOS amd64 -> CodeSwitch-macos-amd64-${TAG}.zip
# - Windows -> CodeSwitch-windows-setup-${TAG}.exe

cat > bin/latest.json <<EOF
{
  "version": "${TAG}",
  "published_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "body": ${BODY_JSON},
  "downloads": {
    "darwin_arm64": "${PUBLIC_URL_BASE}/${TAG}/CodeSwitch-macos-arm64-${TAG}.zip",
    "darwin_amd64": "${PUBLIC_URL_BASE}/${TAG}/CodeSwitch-macos-amd64-${TAG}.zip",
    "windows_amd64": "${PUBLIC_URL_BASE}/${TAG}/${WIN_INSTALLER_UPLOAD}"
  }
}
EOF

aws s3 cp "bin/latest.json" "s3://${R2_BUCKET}/${R2_PREFIX}/latest.json" \
  --endpoint-url "$R2_ENDPOINT" \
  --acl public-read

echo "==> Release $TAG published successfully!"
echo "    Manifest: ${PUBLIC_URL_BASE}/latest.json"
