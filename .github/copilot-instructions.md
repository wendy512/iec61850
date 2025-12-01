# IEC 61850 Go Library - AI Coding Instructions

## Project Overview
CGo-based Go library implementing IEC 61850 protocols (MMS, GOOSE, SV) for industrial automation. Wraps [libiec61850](https://github.com/mz-automation/libiec61850) C library with platform-specific precompiled binaries.

## Architecture

### Core Components
- **Client/Server** (`client.go`, `server.go`): MMS protocol implementation for reading/writing data objects
- **GOOSE** (`goose_publisher.go`, `goose_subscriber.go`): Generic Object-Oriented Substation Event publisher/subscriber
- **SV** (`sv_publisher.go`, `sv_subscriber.go`): Sampled Values for electrical measurements
- **SCL Parser** (`scl/`): XML-based System Configuration Language parser for `.icd`/`.cid` files
- **Model** (`model.go`, `data_model.go`): IEC 61850 data model (IED → LogicalDevice → LogicalNode → DataObject)

### Platform-Specific CGo Configuration
Build tags control platform-specific C library linking (`config_*.go`):
- `config_linux64.go`, `config_armv7l.go`, `config_armv8.go`, `config_darwinarmv8.go`, `config_win64.go`
- Each imports corresponding precompiled library from `libiec61850/lib/<platform>/`
- CGo directives specify CFLAGS/LDFLAGS for header paths and static linking

## Critical Patterns

### CGo Memory Management
**Always** use `defer C.free(unsafe.Pointer(cString))` immediately after `C.CString()`:
```go
cObjectRef := C.CString(objectRef)
defer C.free(unsafe.Pointer(cObjectRef))
C.IedConnection_readBooleanValue(c.conn, &clientError, cObjectRef, C.FunctionalConstraint(fc))
```

### Model Configuration Files
Two approaches to create IED models:
1. **From SCL/ICD files** (XML): Parse with `scl.NewParser()`, generate C code with `scl.NewStaticModelGenerator()`
2. **From .cfg files**: Load directly via `CreateModelFromConfigFileEx("model.cfg")`

Example from `test/server/complexModel_test.go`:
```go
model, err := iec61850.CreateModelFromConfigFileEx("complexModel.cfg")
server := iec61850.NewServerWithConfig(iec61850.NewServerConfig(), model)
```

### Error Handling
Use sentinel errors from `errors.go` (e.g., `NotConnected`, `Timeout`, `ObjectDoesNotExist`):
```go
value, err := client.ReadInt32(objectRef, fc)
if errors.Is(err, iec61850.Timeout) {
    // handle timeout
}
```

### Directory Discovery & Type Inspection
Modern API functions in `client_directory.go` for clean model discovery:

**Recommended approach** - Use `GetServerDirectory()` instead of deprecated `GetLogicalDeviceList()`:
```go
// Get logical devices (no verbose output)
devices, err := client.GetLogicalDeviceNames()

// Get data attributes WITH FC annotations (solves enumeration problem!)
entries, err := client.GetDataDirectoryWithFC("myLD/myLN.myDO")
// Returns: [{Name:"mag", FC:MX}, {Name:"stVal", FC:ST}]

// Get type info BEFORE reading (eliminates guesswork)
spec, err := client.GetVariableSpecification("myLD/myLN.myDO.mag.f", MX)
if spec.Type == Float {
    value, _ := client.ReadFloat32(objectRef, MX)
}

// Filter by specific FC
names, err := client.GetDataDirectoryByFC("myLD/myLN.myDO", MX)
```

### Concurrency for Servers
Server callbacks (`SetHandleWriteAccess`) execute in C thread context:
- **Lock data model** before updates: `server.LockDataModel()` / `server.UnlockDataModel()`
- **Thread-safe updates** via `UpdateFloatAttributeValue()`, `UpdateInt32AttributeValue()`, etc.

## Development Workflow

### Building
```bash
go build ./...                    # Builds with platform-specific CGo config
go test -v ./test/...            # Run integration tests
cd cmd/scltool && go build       # Build SCL tool CLI
```

### Testing Conventions
- Tests require running servers: `server.Start(102)` blocks until signal
- Use `.cfg` files in `test/*/` directories for model initialization
- Client/server tests typically manual (require `syscall.SIGINT` to stop)

### SCL Tool Usage
Generate static C model code from ICD files:
```bash
./scltool genmodel <file.icd> <output_dir> --ied <ied_name> --out static_model
```
Produces `.c`/`.h` files defining IedModel structure.

## Project-Specific Conventions

### Naming Patterns
- **No stutter**: `iec61850.Client` not `iec61850.IEDClient` (package already namespaced)
- **Object references**: Use IEC 61850 dotted notation: `"ied1Inverter/ZINV1.OutVarSet.setMag.f"`
- **Functional Constraints (FC)**: Enum type `FC` (e.g., `FC_MX`, `FC_ST`, `FC_CF`)

### Configuration
- `Settings` struct for client connections (host, port, timeouts)
- `ServerConfig` for server tuning (buffer sizes, max connections)
- `TLSConfig` for mutual TLS authentication (client/server certs, CA chain)

### Platform Build Tags
Limited build tags (e.g., `//go:build linux && amd64` in `goose_subscriber.go`):
- GOOSE/SV features may be platform-restricted
- Check build constraints before using subscriber features

## Common Pitfalls

1. **Memory leaks**: Forgetting `defer C.free()` on C strings
2. **Model lifecycle**: Call `model.Destroy()` and `server.Destroy()` before exit
3. **Blocked tests**: Server tests wait for OS signals; design for manual termination
4. **Platform mismatches**: Ensure test environment matches `config_*.go` build tags

## Key Files Reference
- `types.go`: MmsValue, FC enums, quality flags
- `client_directory.go`: Modern directory/discovery API (GetServerDirectory, GetDataDirectoryWithFC, GetVariableSpecification)
- `client_ld.go`: Legacy GetLogicalDeviceList (deprecated - use GetServerDirectory instead)
- `scl/model.go`: SCL XML structure (IED, LogicalNode, DataSet, ReportControl)
- `scl/static_model_generator.go`: C code generation from SCL (1363 lines)
- `test/icd_file/*.cid`: Sample ICD configuration files
