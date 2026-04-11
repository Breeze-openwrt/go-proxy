# Design: Implement Watchdog Forwarding Model

## 1. 动态对撞机制 (Dynamic Bumping)
为了既保留 `Splice` 性能，又解决“分段加载”：
- **方案**: 封装一个 `trackedConn`，它在 `Write` 和 `Read` 时更新原子时间戳。
- **难点**: `io.Copy` 触发内核零拷贝时会绕过用户态 `Read/Write`。
- **对策**: 在 `Forward` 时启动一个 **“心脏起搏器” (Heartbeat)** 协程。
    - 该协程每隔 `idleTimeout / 2` 检查一次连接的活跃状态。
    - 如果连接活跃，则异步调用 `src.SetReadDeadline(time.Now().Add(idleTimeout))`。

## 2. 协程治理
- 使用 `context` 确保转发结束时，哨兵协程立即退出。
- 为回填任务增加一个全局信号量（Semaphore），防止突发流量下的协程爆炸。

## 3. 核心接口
```go
type Forwarder struct {
    IdleTimeout time.Duration
    MaxRacing   int // 分流/赛马任务上限
}
```
