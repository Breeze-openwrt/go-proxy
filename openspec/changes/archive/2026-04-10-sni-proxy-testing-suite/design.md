# Design: SNI Proxy Testing Suite

## 1. 单元测试策略 (Unit Testing)
- **TDD 模式**: 使用 Table-Driven Tests (表格驱动测试)，涵盖各种 SNI 域名（长域名、短域名、特殊字符）。
- **Mocking**: 对 `net.Conn` 进行 Mock，注入受损的 TLS 包验证防护逻辑。

## 2. 集成测试方案 (Integration Testing)
- **闭环模拟**: 
    1. 在测试中启动两个真实的 HTTPS Server (使用 `httptest`)。
    2. 启动我们的 SNI Proxy。
    3. 编写一个客户端同时向代理发送不同域名的 TLS 请求。
    4. 检查客户端拿回的数据是否对应正确的后端。

## 3. 并发安全性验证
- 专门测试在高并发并发下，`sync.Pool` 是否会出现数据错乱或 Buffer 越界。
