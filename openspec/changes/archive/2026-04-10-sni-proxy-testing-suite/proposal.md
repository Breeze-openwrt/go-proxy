# Proposal: SNI Proxy Testing Suite

## 1. 目标 (Objectives)
构建一个全面的测试矩阵，验证 L4 代理在各种网络情况下的健壮性和准确性。

## 2. 核心指标 (Test Specs)
- **编译正确性**: 全量清理后，`go build` 必须无报错通过。
- **协议兼容性**: 正确提取各种主流浏览器生成的 TLS SNI。
- **异常处理**: 遇到非 TLS 连接或格式错误的包，代理应能优雅关闭而不崩溃。
- **路由准确性**: 不同域名必须 100% 被分发到预期的后端。

## 3. 测试覆盖范围 (Coverage)
- `internal/sniffer`: TLS 解析逻辑。
- `internal/router`: 静态映射匹配。
- `internal/proxy`: 双向转发与 `sync.Pool` 复用安全。
- `integration`: 端到端闭环测试。
