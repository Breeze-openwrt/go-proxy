# Proposal: SNI Proxy Config & Pool Management

## 1. 目标 (Objectives)
实现生产级的配置驱动能力，支持网卡绑定、IPv4/IPv6 混合监听、连接池预热以及闲置连接回收。

## 2. 核心规格 (Specs)
- **配置格式**: `config.jsonc` (支持 // 注释)。
- **网络绑定**: 允许指定 `network_interface` (如 eth0)，并在 Linux 上实现硬件级绑定。
- **高性能连接池**: 
    - `jump_start`: 启动时预先建立的后端连接数。
    - `idle_timeout`: 自动回收长时间不活动的死连接。
- **环境约束**: 仅支持 Linux (amd64)。

## 3. 交付价值
- **零延迟启动**: 通过预热连接池，消除首个请求的握手时间。
- **硬件级隔离**: 支持多网卡环境下的流量定向分发。
- **低资源占用**: 自动清理闲置 TCP 连接。
