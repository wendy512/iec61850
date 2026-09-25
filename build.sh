#!/bin/bash
#
# Linux variants are built in Docker (see docker-compose.yaml).
# macOS variants are built natively, because Mach-O cannot be produced from a Linux container without a macOS SDK / cctools cross toolchain.
# Windows is built locally via the zig drop-in C compiler.
#
# After every build, the resulting archive is validated to ensure (a) it is in the correct binary format for its target OS and (b) that the TLS symbols we depend on at the Go layer are actually present.
# The script exits non-zero on any mismatch so we don't silently ship broken libraries again.

set -euo pipefail

# Versions
MZ_VERSION=1.6.2.1
# Must match the mbedtls-3.6.0 path hardcoded in upstream's Makefile.
MBEDTLS_VERSION=3.6.0
WINPCAP_VERSION=4.1.2

REPO_DIR="./libiec61850-repo"
MBEDTLS_DIR="${REPO_DIR}/third_party/mbedtls/mbedtls-${MBEDTLS_VERSION}"
WINPCAP_ZIP="WpdPack_${WINPCAP_VERSION//./_}.zip"

# patch_makefile fixes two things in upstream's plain Makefile (WITH_MBEDTLS3=1 itself ships with it since v1.6.2.1):
#   1. Its WITH_MBEDTLS3 block enables R-GOOSE/R-SMV but, unlike CMake, doesn't compile src/r_session, leaving RSession_* undefined at link time.
#   2. `ar rcs` instead of `ar r` + ranlib. With mbedtls 3.6 the object list
#      exceeds the threshold above which plain `ar r` produces an archive that
#      ranlib then rejects as "malformed".
patch_makefile() {
    local makefile="$1/Makefile"

    awk '
        /^LIB_SOURCE_DIRS \+= hal\/tls\/mbedtls3$/ {
            print
            print "LIB_SOURCE_DIRS += src/r_session"
            next
        }
        /^\t\$\(AR\) r \$\(LIB_NAME\) \$\(LIB_OBJS\)$/ {
            print "\trm -f $(LIB_NAME)"
            print "\t$(AR) rcs $(LIB_NAME) $(LIB_OBJS)"
            next
        }
        { print }
    ' "${makefile}" > "${makefile}.tmp" && mv "${makefile}.tmp" "${makefile}"
}

# Download sources
echo "Downloading libiec61850 version ${MZ_VERSION} from MZ-Automation..."
if [ -d "${REPO_DIR}" ]; then
    echo "Directory ${REPO_DIR} already exists. Skipping download."
else
    git clone --depth=1 -b "v${MZ_VERSION}" https://github.com/mz-automation/libiec61850.git "${REPO_DIR}"
    patch_makefile "${REPO_DIR}"
    # Upstream enables IED server debug output (printf to stdout) by default.
    sed -i.bak 's/^#define DEBUG_IED_SERVER 1$/#define DEBUG_IED_SERVER 0/' "${REPO_DIR}/config/stack_config.h"
    rm "${REPO_DIR}/config/stack_config.h.bak"
fi

echo "Downloading mbedtls version ${MBEDTLS_VERSION}..."
if [ -d "${MBEDTLS_DIR}" ]; then
    echo "Directory ${MBEDTLS_DIR} already exists. Skipping download."
else
    git clone --depth=1 -b "v${MBEDTLS_VERSION}" https://github.com/Mbed-TLS/mbedtls.git "${MBEDTLS_DIR}"
fi

echo "Downloading Winpcap version ${WINPCAP_VERSION}..."
curl -fL "https://www.winpcap.org/install/bin/${WINPCAP_ZIP}" -o "${WINPCAP_ZIP}"
unzip -qo "${WINPCAP_ZIP}"
cp -r ./WpdPack/Lib "${REPO_DIR}/third_party/winpcap"
cp -r ./WpdPack/Include "${REPO_DIR}/third_party/winpcap"

# verify_archive <archive_path> <expected_format>
#   expected_format: "macho" or "elf"
#
# Fails the script if the archive is in the wrong binary format or if it does not export TLSConfiguration_create (i.e. TLS support was not compiled in).
verify_archive() {
    local archive="$1"
    local expected="$2"

    if [ ! -f "${archive}" ]; then
        echo "ERROR: expected archive ${archive} was not produced" >&2
        exit 1
    fi

    # Resolve to an absolute path so the subshell can find it after `cd`.
    local archive_abs
    archive_abs=$(cd "$(dirname "${archive}")" && pwd)/$(basename "${archive}")

    # Pick an arbitrary object out of the archive and inspect it.
    # bsdtar reads both BSD- and GNU-style archives; macOS ar(1) silently extracts nothing from GNU archives yet exits 0, so it is only the fallback.
    local tmp
    tmp=$(mktemp -d)
    if ! (cd "${tmp}" && bsdtar -xf "${archive_abs}" 2>/dev/null); then
        (cd "${tmp}" && ar -x "${archive_abs}")
    fi
    local sample
    sample=$(find "${tmp}" -name '*.o' | head -n1)
    if [ -z "${sample}" ]; then
        echo "ERROR: ${archive} contains no .o members" >&2
        rm -rf "${tmp}"
        exit 1
    fi

    local info
    info=$(file "${sample}")
    case "${expected}" in
        macho)
            if ! echo "${info}" | grep -q "Mach-O"; then
                echo "ERROR: ${archive} member ${sample##*/} is not Mach-O: ${info}" >&2
                rm -rf "${tmp}"
                exit 1
            fi
            ;;
        elf)
            if ! echo "${info}" | grep -q "ELF"; then
                echo "ERROR: ${archive} member ${sample##*/} is not ELF: ${info}" >&2
                rm -rf "${tmp}"
                exit 1
            fi
            ;;
        *)
            echo "ERROR: unknown expected format '${expected}'" >&2
            rm -rf "${tmp}"
            exit 1
            ;;
    esac
    rm -rf "${tmp}"

    # The Go bindings unconditionally reference these TLS symbols via cgo, so an archive without them will fail to link in any downstream project.
    if ! nm "${archive}" 2>/dev/null | grep -E " T _?TLSConfiguration_create$" >/dev/null; then
        echo "ERROR: ${archive} is missing TLSConfiguration_create — was the library built with WITH_MBEDTLS3=1?" >&2
        exit 1
    fi

    echo "OK: ${archive} (${expected}, TLS symbols present)"
}

# Build Linux variants in Docker
echo "Building Linux variants via docker compose..."
docker compose up --build

# Build macOS variants natively
build_darwin_native() {
    local target_dir="$1"   # e.g. darwin_armv8
    local arch_flag="$2"    # e.g. -arch arm64
    # Install into ./build (not the repo's build dir) so the staging step below picks it up.
    local prefix
    prefix="$(pwd)/build/${target_dir}"

    echo "Building ${target_dir} natively on $(uname -s)/$(uname -m)..."
    rm -rf "${prefix}"
    (
        cd "${REPO_DIR}"
        make clean >/dev/null
        # Via the environment, not as make arguments: those would replace the Makefile's own CFLAGS (defines, -std).
        CFLAGS="${arch_flag} -O2 -g" LDFLAGS="${arch_flag}" \
            make WITH_MBEDTLS3=1 INSTALL_PREFIX="${prefix}" install
    )
}

if [ "$(uname -s)" = "Darwin" ]; then
    case "$(uname -m)" in
        arm64)
            build_darwin_native darwin_armv8 "-arch arm64"
            # Apple clang on Apple Silicon targets x86_64 out of the box.
            build_darwin_native darwin_amd64 "-arch x86_64"
            ;;
        x86_64)
            # If someone runs this on an Intel Mac, build the amd64 variant natively.
            # arm64 cross-compilation from Intel macOS would require a recent Xcode + macOSX.sdk and is intentionally not attempted here.
            build_darwin_native darwin_amd64 "-arch x86_64"
            ;;
        *)
            echo "WARNING: unknown macOS architecture $(uname -m); skipping darwin build" >&2
            ;;
    esac
else
    echo "WARNING: macOS targets cannot be built on $(uname -s); the darwin_*"\
         "archives currently in libiec61850/ will not be refreshed."\
         "Run build.sh on a macOS host (or in CI on a macos-* runner) to"\
         "rebuild them." >&2
fi

# Build Windows locally via zig
(cd "${REPO_DIR}" &&
    make TARGET=WIN64 \
         CC="zig cc -target x86_64-windows-gnu -fno-sanitize=undefined" \
         CPP="zig c++ -target x86_64-windows-gnu" \
         AR="zig ar" RANLIB="zig ranlib" \
         WITH_MBEDTLS3=1 \
         INSTALL_PREFIX=./build/windows_amd64 install
)

# Stage produced libraries into ./libiec61850/<platform>
echo "Copying built libraries to libiec61850 directory..."
mkdir -p ./libiec61850
cp -r ./build/* ./libiec61850/
cp -r "${REPO_DIR}/build/windows_amd64/" ./libiec61850/windows_amd64
# x64 import library for wpcap.dll, renamed so -lwpcap finds it.
cp ./WpdPack/Lib/x64/wpcap.lib ./libiec61850/windows_amd64/lib/libwpcap.a

# Stub Go files so each platform directory is a valid package
echo "Writing Go package stubs for each platform..."
for dir in ./libiec61850/*/; do
    platform=$(basename "${dir}")
    echo "package ${platform}" > "${dir}/include/include.go"
    echo "package ${platform}" > "${dir}/lib/lib.go"
done

# Validate every produced archive before declaring success
echo "Validating produced archives..."
for dir in ./libiec61850/*/; do
    platform=$(basename "${dir}")
    archive="${dir}lib/libiec61850.a"
    case "${platform}" in
        darwin_*)   verify_archive "${archive}" macho ;;
        linux_*)    verify_archive "${archive}" elf  ;;
        win64|windows_*)
            # Windows COFF archives use a different validation path; for now just sanity-check the file exists and exposes the TLS symbol.
            if [ ! -f "${archive}" ]; then
                echo "ERROR: missing ${archive}" >&2
                exit 1
            fi
            if ! nm "${archive}" 2>/dev/null | grep -E " T _?TLSConfiguration_create$" >/dev/null; then
                echo "ERROR: ${archive} is missing TLSConfiguration_create" >&2
                exit 1
            fi
            echo "OK: ${archive} (windows, TLS symbols present)"
            ;;
        *)
            echo "WARNING: no validation rule for platform ${platform}" >&2
            ;;
    esac
done

echo "All archives built and validated successfully."
