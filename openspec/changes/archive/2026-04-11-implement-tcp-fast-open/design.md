# Design: Implement TCP Fast Open (TFO)

## 1. 客户端 TFO 实现 (Outgoing Connections)
在 `internal/proxy/pool.go` 的 `Dial` 逻辑中：
- 使用 `net.Dialer`。
- 设置 `Control` 函数。
- 使用 `syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, 30, 1)`。 (30 即 `TCP_FASTOPEN_CONNECT`)。

## 2. 服务端 TFO 实现 (Incoming Listener)
在 `internal/listener/server.go` 的 `lc.Listen` 阶段：
- 在 `Control` 函数中设置 `syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, 23, 256)`。 (23 即 `TCP_FASTOPEN`)。
- 这允许来自客户端的 TFO 连接。

## 3. 安全性与依赖
- **Fallback**: 若内核不支持，`Setsockopt` 会报错，我们将捕获此错误并静默忽略，回退至普通 TCP。
- **配置项**: 该功能目前默认开启。
