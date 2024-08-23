# iec61850

[![License](https://img.shields.io/badge/license-Apache--2.0-green.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

中文 | [English](README.md)

cgo 版本的 IEC61850 库，参考[libiec61850](https://github.com/mz-automation/libiec61850)

## 概述

iec61850 是实现 MMS、GOOSE 和 SV 协议的 IEC 61850 客户端和服务器库的开源 (Apache-2.0 license) 实现，它可用于在运行 Linux、Windows 的嵌入式系统和 PC 上实施符合 IEC 61850 的客户端和服务器应用程序。本项目依赖并参考了[libiec61850](https://github.com/mz-automation/libiec61850)。

## 功能特性

该库支持以下 IEC 61850 协议功能：

- MMS 客户端/服务器、GOOSE (IEC 61850-8-1)
- 采样值 (SV - IEC 61850-9-2)
- 支持缓冲和非缓冲报告
- 在线报告控制块配置
- 数据访问服务（获取数据、设置数据）
- 在线数据模型发现和浏览
- 所有数据集服务（获取值、设置值、浏览）
- 动态数据集服务（创建和删除）
- 日志服务
- MMS 文件服务（浏览、获取文件、设置文件、删除/重命名文件）
- 设置组处理
- 支持服务跟踪
- GOOSE 和 SV 控制块处理
- TLS 支持

## 如何使用

```shell
go get -u github.com/wendy512/iec61850
```

- [客户端读和写入值](test/client_test.go)
- [客户端控制](test/client_control_test.go)
- [客户端读取 icd 文件](test/scl_test.go)
- [创建服务端](test/server_test.go)

## 开源许可

iec61850 基于 [GPL-3.0 license](./LICENSE) 协议，iec61850 依赖了一些第三方组件，它们的开源协议也为 GPL-3.0 和 MIT。

## 联系方式

- 邮箱：<wendy512@yeah.net>
