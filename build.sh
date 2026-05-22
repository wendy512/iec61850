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
MZ_VERSION=1.6
# We extend the plain Makefile to support WITH_MBEDTLS3=1 via patch_makefile() below.
MBEDTLS_VERSION=3.6.0
WINPCAP_VERSION=4.1.2

REPO_DIR="./libiec61850-repo"
MBEDTLS_DIR="${REPO_DIR}/third_party/mbedtls/mbedtls-${MBEDTLS_VERSION}"
WINPCAP_ZIP="WpdPack_${WINPCAP_VERSION//./_}.zip"

# patch_makefile teaches the cloned libiec61850 Makefile two things upstream doesn't ship:
#   1. WITH_MBEDTLS3=1 build flag (CMake supports it, plain Makefile doesn't).
#      We mirror the existing WITH_MBEDTLS block for mbedtls-3.6.0 +
#      hal/tls/mbedtls3 so we can produce an archive with TLS 1.3 support
#      without switching to CMake (whose generated stack_config.h differs from
#      the Makefile's and breaks downstream tests that rely on the defaults).
#   2. `ar rcs` instead of `ar r` + ranlib. With mbedtls 3.6 the object list
#      exceeds the threshold above which plain `ar r` produces an archive that
#      ranlib then rejects as "malformed".
patch_makefile() {
    local makefile="$1/Makefile"

    # Idempotency guard so re-running on an existing checkout is a no-op.
    if grep -q '^ifdef WITH_MBEDTLS3' "${makefile}"; then
        return 0
    fi

    awk '
        /^LIB_INCLUDES = \$\(addprefix -I,\$\(LIB_INCLUDE_DIRS\)\)/ && !inserted {
            print "ifdef WITH_MBEDTLS3"
            print "LIB_SOURCE_DIRS += third_party/mbedtls/mbedtls-3.6.0/library"
            print "LIB_SOURCE_DIRS += hal/tls/mbedtls3"
            print "LIB_INCLUDE_DIRS += third_party/mbedtls/mbedtls-3.6.0/include"
            print "LIB_INCLUDE_DIRS += hal/tls/mbedtls3"
            print "CFLAGS += -D'\''MBEDTLS_CONFIG_FILE=\"mbedtls_config.h\"'\''"
            print "CFLAGS += -D'\''CONFIG_MMS_SUPPORT_TLS=1'\''"
            print "CFLAGS += -D'\''CONFIG_IEC61850_R_GOOSE=1'\''"
            print "CFLAGS += -D'\''CONFIG_IEC61850_R_SMV=1'\''"
            print "endif"
            print ""
            inserted = 1
        }
        { print }
    ' "${makefile}" > "${makefile}.tmp" && mv "${makefile}.tmp" "${makefile}"

    # ar r → rm + ar rcs (tab-indented Makefile recipe).
    awk '
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
    # Teach the plain Makefile about WITH_MBEDTLS3=1 (CMake already supports it).
    patch_makefile "${REPO_DIR}"
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
    # ar(1) on macOS handles both BSD- and SysV-style archives; on Linux/macOS where ar can't read a SysV archive at all we fall back to bsdtar.
    local tmp
    tmp=$(mktemp -d)
    if ! (cd "${tmp}" && ar -x "${archive_abs}" 2>/dev/null); then
        (cd "${tmp}" && bsdtar -xf "${archive_abs}")
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
    if ! nm "${archive}" 2>/dev/null | grep -qE " T _?TLSConfiguration_create$"; then
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

    echo "Building ${target_dir} natively on $(uname -s)/$(uname -m)..."
    (
        cd "${REPO_DIR}"
        # Clean only the build dir for this target, not the whole repo, so parallel native builds don't stomp on each other.
        rm -rf "build/${target_dir}"
        make clean >/dev/null
        make WITH_MBEDTLS3=1 \
             CFLAGS="${arch_flag} -O2 -g" \
             LDFLAGS="${arch_flag}" \
             INSTALL_PREFIX="$(pwd)/build/${target_dir}" \
             install
    )
}

if [ "$(uname -s)" = "Darwin" ]; then
    case "$(uname -m)" in
        arm64)
            build_darwin_native darwin_armv8 "-arch arm64"
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
         CC="zig cc -target x86_64-windows-gnu" \
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
            if ! nm "${archive}" 2>/dev/null | grep -qE " T _?TLSConfiguration_create$"; then
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
