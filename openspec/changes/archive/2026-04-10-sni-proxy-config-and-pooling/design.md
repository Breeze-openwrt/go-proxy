# Design: SNI Proxy Config & Pool Management

## 1. 配置加载架构 (Config Loader)
- **JSONC 支持**: 采用流水线式加载：`Read File` -> `Filter Comments (//)` -> `JSON Unmarshal`。
- **结构化映射**: 配置对象直接映射到核心模块的初始化参数。

## 2. 深度网卡绑定 (Network Interface Binding)
- **技术栈**: `net.ListenConfig` + `syscall.RawConn`。
- **控制逻辑**: 在 `Control` 钩子中，使用 `syscall.SetsockoptString(fd, SOL_SOCKET, SO_BINDTODEVICE, interfaceName)`。
- **双栈处理**: 默认开启 `ipv6only=0`，确保配置为 `[::]` 时能同时接收 IPv4 和 IPv6。

## 3. 连接预热供应池 (Backend Pre-heating)
- **池管理**: 每个路由条目拥有独立的一个缓冲 Channel。
- **动态补给**: 引入 `PoolGuard` 协程。当 `chan` 中的活跃连接数低于 `jump_start` 时，异步触发新的 `Dial`。
- **清理**: 定时探测池内连接的健康度。

## 4. 闲置超时 (Idle Timeout)
- **追踪器**: 在转发双向流时，封装一个读写追踪器。
- **机制**: 每次读写操作重置 `conn.SetDeadline(time.Now().Add(idle_timeout))`。
