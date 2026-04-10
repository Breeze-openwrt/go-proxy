# Tasks: Fix Performance Issues and Data Corruption

## Phase 1: 核心 Bug 修复 (Step 1)
- [x] Task 1: 在 `pool.go` 中移除破坏性的 `Read(1)`
- [x] Task 2: 实现基于 `MSG_PEEK` 的 `isAlive` 检查函数
- [x] 验证：确保连接池复用时不再丢失首字节。

## Phase 2: 转发逻辑校准 (Step 2)
- [x] Task 3: 优化 `proxy.go` 中的超时重试逻辑
- [x] Task 4: 移除不必要的 `SetReadDeadline` 操作

## Phase 3: 整体稳定性验证 (Step 3)
- [x] Task 5: 验证大量小文件并发请求下的加载速度
- [x] Task 6: 验证视频流等大带宽场景的稳定性
