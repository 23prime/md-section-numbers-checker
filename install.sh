#!/usr/bin/env sh
set -eu

REPO="23prime/md-section-numbers-checker"
BIN_NAME="md-section-numbers-checker"
INSTALL_DIR="${INSTALL_DIR:-${HOME}/.local/bin}"

# ---- detect OS ----
OS="$(uname -s)"
case "${OS}" in
  Linux)  GOOS="linux" ;;
  Darwin) GOOS="darwin" ;;
  *)
    echo "Unsupported OS: ${OS}" >&2
    exit 1
    ;;
esac

# ---- detect arch ----
ARCH="$(uname -m)"
case "${ARCH}" in
  x86_64 | amd64)  GOARCH="amd64" ;;
  arm64 | aarch64) GOARCH="arm64" ;;
  *)
    echo "Unsupported architecture: ${ARCH}" >&2
    exit 1
    ;;
esac

# ---- resolve version ----
if [ -z "${VERSION:-}" ]; then
  echo "Fetching latest release version..."
  VERSION="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')"
  if [ -z "${VERSION}" ]; then
    echo "Failed to fetch latest version." >&2
    exit 1
  fi
fi

echo "Installing ${BIN_NAME} ${VERSION} (${GOOS}/${GOARCH})..."

ARCHIVE="${BIN_NAME}-${GOOS}-${GOARCH}.tar.gz"
BASE_URL="https://github.com/${REPO}/releases/download/${VERSION}"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

# ---- download archive and checksums ----
curl -fsSL "${BASE_URL}/${ARCHIVE}"       -o "${TMP_DIR}/${ARCHIVE}"
curl -fsSL "${BASE_URL}/checksums.txt"    -o "${TMP_DIR}/checksums.txt"

# ---- verify checksum ----
cd "${TMP_DIR}"
if command -v sha256sum > /dev/null 2>&1; then
  grep " ${ARCHIVE}$" checksums.txt | sha256sum -c - || { echo "Checksum mismatch!" >&2; exit 1; }
elif command -v shasum > /dev/null 2>&1; then
  grep " ${ARCHIVE}$" checksums.txt | shasum -a 256 -c - || { echo "Checksum mismatch!" >&2; exit 1; }
else
  echo "Warning: no sha256sum/shasum found, skipping checksum verification." >&2
fi
cd - > /dev/null

# ---- extract ----
tar -xzf "${TMP_DIR}/${ARCHIVE}" -C "${TMP_DIR}"

# ---- install ----
mkdir -p "${INSTALL_DIR}"
if [ -w "${INSTALL_DIR}" ]; then
  mv "${TMP_DIR}/${BIN_NAME}-${GOOS}-${GOARCH}" "${INSTALL_DIR}/${BIN_NAME}"
  chmod +x "${INSTALL_DIR}/${BIN_NAME}"
else
  echo "Root privileges required to install to ${INSTALL_DIR}."
  sudo mv "${TMP_DIR}/${BIN_NAME}-${GOOS}-${GOARCH}" "${INSTALL_DIR}/${BIN_NAME}"
  sudo chmod +x "${INSTALL_DIR}/${BIN_NAME}"
fi

echo "Installed: ${INSTALL_DIR}/${BIN_NAME}"
"${INSTALL_DIR}/${BIN_NAME}" --version
