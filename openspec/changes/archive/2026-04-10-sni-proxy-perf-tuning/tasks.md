# Tasks: SNI Proxy Performance Tuning

## Phase 1: 测量与基准 (Step 1)
- [x] Task 1: 建立基础性能基准 (Time-to-Serve)
- [x] 验证：日志中能精确显示每个连接从“握手”到“转发”的纳秒级耗时。

## Phase 2: 响应式内存管理 (Step 2)
- [x] Task 2: 实现 Buffer 复用 (sync.Pool Integration)
- [x] 验证：程序在高并发下稳健运行，无内存泄漏或 Parnic。

## Phase 3: 路径优化 (Step 3-4)
- [x] Task 3: 优化数据转发路径 (Check Splice Compatibility)
- [x] Task 4: Goroutine 调度模型检查
- [x] 验证：通过日志观察吞吐量提升。

## Phase 4: 系统观测与进阶 (Step 5-6)
- [x] Task 5: 全局日志降级与 Level 控制
- [x] Task 6: [可选] 后端连接复用 (Connection Pooling)
- [x] 验证：性能压测结果对比，展示优化涨幅。
