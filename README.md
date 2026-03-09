# GeeCache | 分布式缓存（Go）

[English](#english) | [中文](#中文)

---

## 中文

GeeCache 是一个使用 Go 实现的轻量级分布式缓存项目，核心目标是提供：

- 本地高性能 LRU 缓存
- 基于一致性哈希的节点选择
- 节点间 HTTP + Protobuf 通信
- `singleflight` 并发控制，避免缓存击穿

该仓库适合作为分布式缓存的学习与实验项目。

### 功能特性

- `lru/`: 手写 LRU 缓存淘汰策略，支持容量限制与淘汰回调
- `cache.go`: 为缓存访问增加互斥锁，保证并发安全
- `geecache.go`: `Group` 作为命名空间，统一管理缓存读取流程
- `consistenthash/`: 一致性哈希环 + 虚拟节点
- `http.go`: `HTTPPool` 负责节点间请求转发与服务
- `geecachepb/`: Protobuf 协议定义与生成代码
- `singleflight/`: 同 key 请求合并，减少对后端数据源的重复访问

### 项目结构

```text
.
├── byteview.go               # 只读值封装（防止外部修改）
├── cache.go                  # 并发安全缓存包装
├── geecache.go               # Group、加载流程、本地/远程回源
├── http.go                   # HTTPPool 与远程节点访问
├── peers.go                  # PeerPicker / PeerGetter 接口
├── lru/                      # LRU 实现与测试
├── consistenthash/           # 一致性哈希实现与测试
├── singleflight/             # 请求合并实现
└── geecachepb/               # protobuf 定义与生成文件
```

### 快速开始

#### 1) 克隆并进入项目

```bash
git clone git@github.com:ihoooohi/GeeCache.git
cd GeeCache
```

#### 2) 运行测试

```bash
go test ./...
```

### 核心流程（读缓存）

1. 调用 `Group.Get(key)`
2. 先查本地 LRU，命中直接返回
3. 未命中时进入 `singleflight`，合并并发请求
4. 若配置了远程节点，使用一致性哈希选择 peer 获取
5. 远程失败或无可用 peer，则调用本地 `Getter` 回源
6. 回源结果写入本地缓存并返回

### 使用示例（本地 Getter）

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
        }))

    v, err := group.Get("Tom")
    if err != nil {
        panic(err)
    }
    fmt.Println(v.String()) // 630
}
```

### 依赖

- Go `1.25.4`（以 `go.mod` 为准）
- `google.golang.org/protobuf`

---

## English

GeeCache is a lightweight distributed cache written in Go. It focuses on:

- High-performance local LRU cache
- Peer selection via consistent hashing
- Inter-node communication through HTTP + Protobuf
- Cache-breakdown protection with `singleflight`

This repository is suitable for learning and experimenting with distributed caching internals.

### Features

- `lru/`: custom LRU with capacity limit and eviction callback
- `cache.go`: mutex-protected cache access for concurrency safety
- `geecache.go`: `Group` namespace and unified load path
- `consistenthash/`: hash ring with virtual nodes
- `http.go`: peer HTTP server and client (`HTTPPool`)
- `geecachepb/`: protobuf schema and generated code
- `singleflight/`: de-duplicate concurrent loads for the same key

### Project Layout

```text
.
├── byteview.go               # immutable value wrapper
├── cache.go                  # thread-safe cache wrapper
├── geecache.go               # Group, local/remote load path
├── http.go                   # HTTPPool and remote fetch
├── peers.go                  # PeerPicker / PeerGetter interfaces
├── lru/                      # LRU implementation and tests
├── consistenthash/           # consistent-hash implementation and tests
├── singleflight/             # request coalescing
└── geecachepb/               # protobuf definitions and generated files
```

### Quick Start

#### 1) Clone and enter the repo

```bash
git clone git@github.com:ihoooohi/GeeCache.git
cd GeeCache
```

#### 2) Run tests

```bash
go test ./...
```

### Read Path Overview

1. Call `Group.Get(key)`
2. Check local LRU first
3. On miss, enter `singleflight` to merge concurrent loads
4. If peers are registered, pick one via consistent hashing
5. If remote fetch fails (or no peer), fallback to local `Getter`
6. Populate local cache and return

### Minimal Example

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
        }))

    v, err := group.Get("Tom")
    if err != nil {
        panic(err)
    }
    fmt.Println(v.String()) // 630
}
```

### Dependencies

- Go `1.25.4` (as declared in `go.mod`)
- `google.golang.org/protobuf`
