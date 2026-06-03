#!/usr/bin/env bash

if [ -z "${BASH_VERSION:-}" ]; then
  exec /usr/bin/env bash
fi

set -e
set -o pipefail

info() {
  printf '%s\n' "$*"
}

die() {
  printf 'Error: %s\n' "$*" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "Requiere '$1' instalado."
}

os_name="$(uname -s 2>/dev/null || true)"
case "$os_name" in
  Linux)
    os="linux"
    ;;
  Darwin)
    os="darwin"
    ;;
  MINGW*|MSYS*|CYGWIN*|Windows_NT)
    die "Windows no está soportado por este instalador. Descargá el binario desde https://github.com/U23205221-S/env-vault/releases"
    ;;
  *)
    die "Sistema operativo no soportado: ${os_name}. Soportados: linux, darwin"
    ;;
esac

arch_name="$(uname -m 2>/dev/null || true)"
case "$arch_name" in
  x86_64)
    arch="amd64"
    ;;
  arm64|aarch64)
    arch="arm64"
    ;;
  *)
    die "Arquitectura no soportada: ${arch_name}. Soportadas: x86_64, arm64, aarch64"
    ;;
esac

require_cmd curl
require_cmd tar

checksum_cmd=""
if command -v sha256sum >/dev/null 2>&1; then
  checksum_cmd="sha256sum"
elif command -v shasum >/dev/null 2>&1; then
  checksum_cmd="shasum"
else
  die "Requiere 'sha256sum' o 'shasum' para verificar checksums."
fi

version_input="${1:-}"
if [ -n "$version_input" ]; then
  version="${version_input#v}"
else
  info "Resolviendo versión latest..."
  # Bajamos el JSON a un archivo temporal porque curl con `-f` reporta
  # "Failure writing output" cuando awk hace `exit` antes de consumir
  # todo el body. Con un archivo destino, ese caso no ocurre.
  tmp_release_json="$(mktemp)"
  if ! curl -sSfL "https://api.github.com/repos/U23205221-S/env-vault/releases/latest" -o "$tmp_release_json"; then
    rm -f "$tmp_release_json"
    die "No se pudo resolver la versión latest desde la API de GitHub."
  fi
  tag_name="$(awk -F'"' '/"tag_name"/ {print $4; exit}' "$tmp_release_json")"
  rm -f "$tmp_release_json"
  [ -n "$tag_name" ] || die "La respuesta de la API no contiene tag_name."
  version="${tag_name#v}"
fi

filename="env-vault_${version}_${os}_${arch}.tar.gz"
release_url="https://github.com/U23205221-S/env-vault/releases/download/v${version}"
archive_url="${release_url}/${filename}"
checksums_url="${release_url}/checksums.txt"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

info "Descargando ${archive_url}"
curl -sSfL "$archive_url" -o "$tmpdir/$filename"
curl -sSfL "$checksums_url" -o "$tmpdir/checksums.txt"

expected_checksum="$(grep " ${filename}$" "$tmpdir/checksums.txt" | awk '{print $1}')"
[ -n "$expected_checksum" ] || die "No se encontró checksum para ${filename}."

if [ "$checksum_cmd" = "sha256sum" ]; then
  actual_checksum="$(sha256sum "$tmpdir/$filename" | awk '{print $1}')"
else
  actual_checksum="$(shasum -a 256 "$tmpdir/$filename" | awk '{print $1}')"
fi
[ -n "$actual_checksum" ] || die "No se pudo calcular el checksum."

if [ "$expected_checksum" != "$actual_checksum" ]; then
  die "Checksum inválido para ${filename}."
fi

tar -xzf "$tmpdir/$filename" -C "$tmpdir"

bin_path="$tmpdir/env-vault"
if [ ! -f "$bin_path" ]; then
  bin_path="$(find "$tmpdir" -maxdepth 2 -type f -name env-vault | awk 'NR==1 {print; exit}')"
fi

[ -n "$bin_path" ] && [ -f "$bin_path" ] || die "No se encontró el binario env-vault en el archivo descargado."

install_dir=""
if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
  install_dir="/usr/local/bin"
else
  if [ -d "$HOME/.local/bin" ] || mkdir -p "$HOME/.local/bin" 2>/dev/null; then
    install_dir="$HOME/.local/bin"
    info "Usando ${install_dir}. Asegurate de tenerlo en tu PATH."
  else
    die "No se pudo usar /usr/local/bin ni crear ~/.local/bin."
  fi
fi

mv "$bin_path" "$install_dir/env-vault"
chmod +x "$install_dir/env-vault"

info "env-vault ${version} instalado en ${install_dir}/env-vault"
info "Run 'env-vault --help' to get started"
