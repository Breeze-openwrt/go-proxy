# Tasks: Implement Asynchronous Zero-delay Pool Refill

- [ ] Task 1: 在 `BackendPool` 中新增 `triggerRefill(addr string)` 私有方法
- [ ] Task 2: 在 `Get` 逻辑的所有出口点集成 `triggerRefill` 呼叫
- [ ] Task 3: 优化 `PreHeat` 逻辑，使其与 `triggerRefill` 复用代码
- [ ] Task 4: 连通性与高并发压力测试
