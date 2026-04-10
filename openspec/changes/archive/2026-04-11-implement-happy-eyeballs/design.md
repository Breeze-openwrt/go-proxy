# Design: Implement Happy Eyeballs v2 (Concurrent Dialing)

## 1. 核心流程
重构 `dialWithTFO` 为 `dialHappyEyeballs`：
1. **DNS 预处理**: 分离 `host` 和 `port`，并发查询 `AAAA` 和 `A` 记录。
2. **连接编排**:
   - 将所有 IP 交错排序（IPv6, IPv4, IPv6, IPv4...）。
   - 启动一个循环，向每个 IP 发起带 `context.WithCancel` 的 `dial`。
   - 每路拨号间隔 250ms。
3. **完成处理**:
   - 第一个成功的 `net.Conn` 通过一个带 Buffer 的 Channel 发回。
   - 立即调用 `cancel()` 取消所有其他协程。

## 2. 与 TFO 的结合
- 每个独立的 Dial 任务依然会设置 `TCP_FASTOPEN_CONNECT`。
- 只有“赢家”才会进入应用层数据写入阶段。

## 3. 错误边界
- 如果所有路都失败，返回最后一个捕获的错误或“全部失败”错误。
