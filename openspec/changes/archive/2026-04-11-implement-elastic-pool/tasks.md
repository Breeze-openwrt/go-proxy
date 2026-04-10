# Tasks: Implement Elastic Backend Connection Pool

## Phase 1: 结构升级 (Step 1)
- [ ] Task 1: 在 `internal/proxy/pool.go` 中定义 `idleConn` 并重构 `BackendPool`
- [ ] Task 2: 实现 `PreHeat` 对新结构的兼容

## Phase 2: 获取与自愈逻辑 (Step 2)
- [ ] Task 3: 实现 `Get` 逻辑：包含从 Channel 获取、过期检查、Peek 存活检测
- [ ] Task 4: 实现 `Put` 逻辑：包含非阻塞回收与溢出自动关闭

## Phase 3: 集成验证 (Step 3)
- [ ] Task 5: 在 `main.go` / `server.go` 中确保调用正确
- [ ] Task 6: 压力测试验证连接复用率（观察 FD 增长情况）
