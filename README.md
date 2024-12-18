# iec61850

[![License](https://img.shields.io/badge/license-GPL--3.0-green.svg)](https://www.gnu.org/licenses/gpl-3.0.html)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/wendy512/iec61850)](https://pkg.go.dev/mod/github.com/wendy512/iec61850)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.0-61CFDD.svg?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/wendy512/iec61850?style=flat-square)](https://goreportcard.com/report/github.com/wendy512/iec61850)

English | [中文](README_zh_CN.md)

cgo version of IEC 61850 library, reference [libiec61850](https://github.com/mz-automation/libiec61850)

## Overview

iec61850 is an open source (GPL-3.0 license) implementation of the IEC 61850 client and server library that implements the MMS, GOOSE and SV protocols.
It can be used to implement IEC 61850 compliant clients and PCs on embedded systems and PCs running Linux, Windows Server application.
This project relies on and refers to [libiec61850](https://github.com/mz-automation/libiec61850).

## Features

The library support the following IEC 61850 protocol features:

- MMS client/server, GOOSE (IEC 61850-8-1)
- Sampled Values (SV - IEC 61850-9-2)
- Support for buffered and unbuffered reports
- Online report control block configuration
- Data access service (get data, set data)
- Online data model discovery and browsing
- All data set services (get values, set values, browse)
- Dynamic data set services (create and delete)
- Log service
- MMS file services (browse, get file, set file, delete/rename file)
- Setting group handling
- Support for service tracking
- GOOSE and SV control block handling
- TLS support

## How to use

```shell
go get -u github.com/wendy512/iec61850
```

- [Client control operations](test/client_control/client_control_test.go)
- [Client rcb operations](test/client_rcb/client_rcb_test.go)
- [Client read and write](test/client_rw)
- [Client setting groups](test/client_sg/client_sg_test.go)
- [Create tls client](test/tls_client/client_read_test.go)
- [Server handle write access](test/server/complexModel_test.go)
- [Server handle control](test/server/simpleIO_control_test.go)
- [Server handle direct control](test/server/simpleIO_direct_control_goose_test.go)
- [Create tls server](test/tls_server/tls_server_test.go)


## License

iec61850 is based on the [GPL-3.0 license](./LICENSE) agreement, and iec61850 relies on some third-party components whose open source agreement is GPL-3.0 and MIT.

## Contact

- Email：<wendy512@yeah.net>
