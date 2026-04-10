# Tasks: Implement TCP Fast Open (TFO)

## Phase 1: 客户端 TFO (Step 1)
- [ ] Task 1: 在 `BackendPool` 的 `Dial` 逻辑中注入 TFO 控制选项
- [ ] Task 2: 在 `PreHeat` 逻辑中同步应用此优化

## Phase 2: 监听端 TFO (Step 2)
- [ ] Task 3: 在 `listener/server.go` 中开启监听端口的 TFO 接收能力

## Phase 3: 验证 (Step 3)
- [ ] Task 4: 编译并在支持 TFO 的环境下进行连通性测试
- [ ] Task 5: 验证在高延迟网络环境下的“起步”提速效果
