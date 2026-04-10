# Design: Implement High-Performance Zero-copy Forwarding (Splice)

## 1. 核心原理
在 Go 中，当 `dst` 实现了 `io.ReaderFrom` 接口（如 `*net.TCPConn`），且 `src` 是一个能够提供 FD 的 `*net.TCPConn` 时，`io.Copy` 会底层调用 `splice`。

## 2. 带有超时控制的 Splice
由于 `splice` 是长阻塞操作，我们需要确保 `idleTimeout` 依然生效。
1. 在开始 `io.Copy` 之前设置 `SetReadDeadline`。
2. `io.Copy` 会在超时发生时返回错误。
3. 捕获超时错误并优雅结束。

## 3. 结构重构
我们将重构 `internal/proxy/proxy.go` 中的 `copyWithTimeout`。

### 降级策略
1. 检查 `dst` 是否能断言为 `interface{ ReadFrom(io.Reader) (int64, error) }`。
2. 检查 `src` 是否是 `*net.TCPConn`。
3. 若不成立，使用现有的 `bufferPool` 循环方案。
