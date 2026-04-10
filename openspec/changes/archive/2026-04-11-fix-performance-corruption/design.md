# Design: Fix Performance Issues and Data Corruption

## 1. 存活探测重构 (Pool Health Check)
修改 `internal/proxy/pool.go`：
- 禁止调用 `conn.Read()`。
- 使用 `rawConn.Control` 挂钩。
- 执行 `syscall.Recvfrom(fd, buf, syscall.MSG_PEEK | syscall.MSG_DONTWAIT)`。
- **逻辑判断**:
    - `n == 0 && err == nil`: 连接已关闭 (EOF)。
    - `err == EAGAIN`: 连接活跃且无数据（理想状态）。
    - `n > 0`: 连接活跃且有滞留数据（安全）。

## 2. 转发稳定性 (Forwarding Stability)
在 `internal/proxy/proxy.go` 中：
- 确保 `SetReadDeadline` 的设置频率不会过高。
- 增加对 `io.Copy` 返回错误的更细致处理，避免在某些临时错误下重置连接。

## 3. 回退机制
如果 `SyscallConn` 断言失败，则直接跳过健康检查，优先保证数据不被破坏。
