#!/bin/bash

# error handling
set -euo pipefail

# MZ version to download
MZ_VERSION=1.6

# download the MZ-Automation libiec61850 library and extract it
echo "Downloading libiec61850 version ${MZ_VERSION} from MZ-Automation..."
if [ -d "./libiec61850-repo" ]; then
    echo "Directory ./libiec61850-repo already exists. Skipping download."
else
    git clone -b v${MZ_VERSION} https://github.com/mz-automation/libiec61850.git ./libiec61850-repo

fi

# add the third_party libraries
echo "Downloading third party libraries..."

if [ -d "./libiec61850-repo/third_party/mbedtls/mbedtls-3.6.0" ]; then
    echo "Directory ./libiec61850-repo/third_party/mbedtls/mbedtls-3.6.0 already exists. Skipping download."
else
    echo "Downloading mbedtls version 3.6.0..."
    git clone -b v3.6.0 https://github.com/Mbed-TLS/mbedtls.git ./libiec61850-repo/third_party/mbedtls/mbedtls-3.6.0
fi

# Winpcap
echo "Downloading Winpcap version 4.1.2..."
curl -fL "https://www.winpcap.org/install/bin/WpdPack_4_1_2.zip" -o "WpdPack_4_1_2.zip"

unzip -qo WpdPack_4_1_2.zip

# copy the lib and include directories to the third_party/winpcap directory
cp -r ./WpdPack/Lib ./libiec61850-repo/third_party/winpcap
cp -r ./WpdPack/Include ./libiec61850-repo/third_party/winpcap

# build all the linux and mac libraries using docker
docker compose up --build

# compile for windows locally using the zig CC drop-in compiler in a sub shell
(cd libiec61850-repo || exit && 
make TARGET=WIN64 CC="zig cc -target x86_64-windows-gnu" CPP="zig c++ -target x86_64-windows-gnu" \
AR="zig ar" RANLIB="zig ranlib" WITH_MBEDTLS3=1 \
INSTALL_PREFIX=./build/windows_amd64 install
)

# now copy the built libraries to the libiec61850 directory 
echo "Copying built libraries to libiec61850 directory..."
cp -r ./build/* ./libiec61850/

cp -r ./libiec61850-repo/build/windows_amd64/ ./libiec61850/windows_amd64

# add the go files to each platform's include directory so that they're available for cgo
echo "Copying go files to each platform's include directory..."

# for each directory in the libiec61850 directory, touch a go file and write `package platform` to it, where platform is the name of the directory
for dir in ./libiec61850/*/; do
    platform=$(basename "$dir")
    include_go_file="$dir/include/include.go"
    echo "package $platform" > "$include_go_file"

    lib_go_file="$dir/lib/lib.go"
    echo "package $platform" > "$lib_go_file"
done