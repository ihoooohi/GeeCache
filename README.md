# GeeCache

[中文文档](./README.zh-CN.md)

GeeCache is a lightweight distributed caching framework implemented in Go.
It demonstrates the core building blocks of a distributed cache system: local LRU storage, consistent hashing for peer selection, HTTP + Protobuf peer communication, and singleflight request coalescing.

## Features

- Local in-memory cache with LRU eviction
- Namespace-based cache management via `Group`
- Consistent hash ring with virtual nodes
- HTTP-based peer communication with Protobuf payloads
- `singleflight` protection to avoid cache breakdown
- Thread-safe cache operations

## Project Structure

```text
.
├── byteview.go               # immutable value wrapper
├── cache.go                  # thread-safe cache wrapper
├── geecache.go               # Group and load path (local/remote)
├── http.go                   # HTTPPool and peer getter
├── peers.go                  # PeerPicker / PeerGetter interfaces
├── lru/                      # LRU implementation and tests
├── consistenthash/           # hash ring implementation and tests
├── singleflight/             # duplicate suppression for concurrent loads
└── geecachepb/               # protobuf definitions and generated code
```

## Requirements

- Go `1.25.4` (as declared in `go.mod`)
- `google.golang.org/protobuf`

## Quick Start

### 1. Clone

```bash
git clone git@github.com:ihoooohi/GeeCache.git
cd GeeCache
```

### 2. Run tests

```bash
go test ./...
```

## How It Works (Read Path)

1. Call `Group.Get(key)`.
2. Try local cache first.
3. On miss, enter `singleflight` to merge concurrent requests for the same key.
4. If peers are configured, pick a node by consistent hashing and fetch remotely.
5. If remote fetch fails (or no peer is available), fall back to local `Getter`.
6. Store result in local cache and return.

## Minimal Example

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

## License

This project is for learning and experimentation.
