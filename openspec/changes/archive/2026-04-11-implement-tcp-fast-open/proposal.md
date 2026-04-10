# Proposal: Implement TCP Fast Open (TFO) for Lower Latency

## 1. 问题定位 (Problem)
虽然我们已经有了连接池，但在“冷启动”或“连接耗尽”需重新 Dial 时，传统的 TCP 三次握手会引入 1 个 RTT 的固定延迟。
- **SNI 探测延迟**: 需要等待握手完成后才能发送 ClientHello。
- **用户体验**: 在长距离跨地域转发中，握手带来的延迟感非常明显。

## 2. 改进规格 (Specs)
- **客户端 TFO**: 在连接后端服务器时，启用 `TCP_FASTOPEN_CONNECT` 选项。
- **服务端 TFO**: 在监听端口时，通过 `SetExtraFiles` 或底层配置开启 TFO 队列，允许客户端通过 TFO 连接到我们。
- **环境适配**: 仅在支持的内核 (Linux 4.11+) 上激活。

## 3. 交付价值
- 消除新连接建立时的 1 个 RTT 等待。
- 极致的单次连接响应速度。
