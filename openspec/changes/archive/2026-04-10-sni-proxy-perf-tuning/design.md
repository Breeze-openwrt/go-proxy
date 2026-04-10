# Design: SNI Proxy Performance Tuning

## 1. 内存优化：Buffer 复用 (Buffer Reuse)
**痛点**：当前每当有新连接进入，我们都会为转发逻辑分配新的字节数组。在高并发下，这会导致大量的内存碎片和频繁的 GC（垃圾回收）。
**方案**：引入 `sync.Pool`。它像一个“共享仓库”，暂时不用的 Buffer 不会被直接销毁，而是放回池子里，供下一个连接直接复用。

## 2. 转发优化：零拷贝深入 (Zero-Copy)
**方案**：虽然 `io.Copy` 已经很好，但我们要确保 `src` 和 `dst` 都能触发内核态的 `splice`。我们会检查 `bufferedConn` 是否阻碍了这种优化。

## 3. 日志与指标 (Logging & Metrics)
**方案**：
- **计时统计**：在 `handleConn` 中加入延迟计算逻辑。
- **降级**：将“连接进入”等高频日志设为 `Debug` 级别，仅在需要时开启，避免字符串拼接和磁盘 I/O 成为瓶颈。

## 4. 并发模型 (Concurrency)
**方案**：虽然每个连接一个 Goroutine 是 Go 的强项，但我们要确保没有不必要的锁竞争（Lock Contention）。
