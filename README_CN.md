<div align="center">

# GeeCache

[![License](https://img.shields.io/badge/License-Unspecified-lightgrey.svg)]()
[![Language](https://img.shields.io/badge/language-Go-00ADD8.svg)](https://go.dev/)
[![Go Version](https://img.shields.io/badge/go-1.25.4-00ADD8.svg)](https://go.dev/doc/devel/release)
[![Protobuf](https://img.shields.io/badge/protobuf-v1.36.11-blue.svg)](https://pkg.go.dev/google.golang.org/protobuf)
[![Architecture](https://img.shields.io/badge/architecture-LRU%20%7C%20Consistent%20Hash%20%7C%20Singleflight-informational.svg)]()

[**English**](./README.md) | [**中文**](./README_CN.md)

</div>

---

GeeCache 是一个用 Go 实现的轻量级分布式缓存框架。
它演示了分布式缓存系统中的关键能力：本地 LRU 缓存、一致性哈希节点选择、HTTP + Protobuf 节点通信，以及 singleflight 并发请求合并。

## 功能特性

- 基于内存的本地 LRU 淘汰缓存
- 通过 `Group` 做命名空间级缓存管理
- 支持虚拟节点的一致性哈希环
- 基于 HTTP + Protobuf 的节点间通信
- 使用 `singleflight` 防止缓存击穿
- 并发安全的缓存读写

## 项目结构

```text
.
├── byteview.go               # 不可变值封装
├── cache.go                  # 并发安全缓存包装
├── geecache.go               # Group 与加载流程（本地/远程）
├── http.go                   # HTTPPool 与远程节点获取
├── peers.go                  # PeerPicker / PeerGetter 接口
├── lru/                      # LRU 实现与测试
├── consistenthash/           # 一致性哈希实现与测试
├── singleflight/             # 并发请求合并
└── geecachepb/               # protobuf 定义与生成代码
```

## 环境要求

- Go `1.25.4`（以 `go.mod` 为准）
- `google.golang.org/protobuf`

## 快速开始

### 1. 克隆仓库

```bash
git clone git@github.com:ihoooohi/GeeCache.git
cd GeeCache
```

### 2. 运行测试

```bash
go test ./...
```

## 工作流程（读路径）

1. 调用 `Group.Get(key)`。
2. 优先查询本地缓存。
3. 未命中时进入 `singleflight`，合并同 key 并发请求。
4. 如果配置了节点，基于一致性哈希选择远程节点获取。
5. 若远程失败（或无可用节点），回退到本地 `Getter`。
6. 将结果写入本地缓存并返回。

## 最小示例

```go
package main

import (
    "fmt"
    "log"

    "geecache"
)

func main() {
    db := map[string]string{"Tom": "630", "Jack": "589"}

    group := geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
        func(key string) ([]byte, error) {
            log.Println("[SlowDB] search key", key)
            if v, ok := db[key]; ok {
                return []byte(v), nil
            }
            return nil, fmt.Errorf("%s not found", key)
        }),
    )

    v, err := group.Get("Tom")
    if err != nil {
        panic(err)
    }
    fmt.Println(v.String()) // 630
}
```

## 许可说明

本项目用于学习和实验。
