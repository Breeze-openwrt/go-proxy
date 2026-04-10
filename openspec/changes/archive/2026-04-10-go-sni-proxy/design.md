# Design: Go SNI Proxy

## 1. 架构概览 (DDD 模式)

我们将项目划分为以下领域（Domain）：

### 1.1 Listener (监听层)
- **职责**：监听本地端口，接收 `net.Conn`。
- **关键技术**：`net.ListenTCP`。

### 1.2 Sniffer (探测层)
- **职责**：分析连接前几个字节，识别 TLS 握手，提取 SNI。
- **关键技术**：Bufio Peek（不破坏后续读取）或手动解析 TLS ClientHello。

### 1.3 Router (路由层)
- **职责**：根据 SNI 字符串查找后端地址。
- **数据结构**：`map[string]string`。

### 1.4 Proxy (转发层)
- **职责**：建立后端连接，同步搬运原始流量。
- **关键技术**：`io.Copy` (利用 `sendfile` 或 `splice` 的潜在优化)。

## 2. 交互流程
1. 连接进入 -> Listener 接收。
2. Listener 将连接交给 Sniffer。
3. Sniffer 预读数据，提取 SNI 后归还“原始”连接流。
4. Router 根据 SNI 找到 Target Address。
5. Proxy 建立 Dial 到 Target，并开始 `io.Copy`。
