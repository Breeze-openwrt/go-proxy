# Tasks: Implement High-Performance Zero-copy Forwarding (Splice)

## Phase 1: 重构转发内核 (Step 1)
- [ ] Task 1: 在 `internal/proxy/proxy.go` 中重构 `copyWithTimeout`
- [ ] Task 2: 优先尝试 `io.Copy` (利用其内部 Splice 优化)
- [ ] Task 3: 保留原有的缓存池方案作为跨平台回退

## Phase 2: 集成与验证 (Step 2)
- [ ] Task 4: 验证超时控制在 Splice 模式下是否依然有效
- [ ] Task 5: 验证在高吞吐下的 CPU 占用率变化
