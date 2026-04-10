# Tasks: SNI Proxy Testing Suite

## Phase 1: 基础自检 (Step 1)
- [x] Task 1: 编译验证与二进制可用性测试
- [x] 验证：`go build` 正常，且 `./sni-proxy --help` (如果有) 或基础运行不崩溃。

## Phase 2: 核心单元验证 (Step 2-3)
- [x] Task 2: 补全 Router 模块的单元测试 (覆盖各种匹配规则)
- [x] Task 3: 补全 Proxy 模块的并发安全性测试
- [x] 验证：`go test -v ./internal/router/... ./internal/proxy/...` 全部通过。

## Phase 3: 端到端集成测试 (Step 4)
- [x] Task 4: 编写集成测试套件 (Integration Test Suite)
- [x] 验证：运行测试脚本，模拟双后端转发，结果 100% 匹配。

## Phase 5: 异常边界测试
- [x] Task 5: 注入非法报文测试 Sniffer 的容错性
- [x] 验证：非 TLS 流量被正确拒绝，且不影响其他连接。
