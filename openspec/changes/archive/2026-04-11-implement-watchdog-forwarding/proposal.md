# Proposal: Implement Watchdog Forwarding Model

## 1. 问题定位 (Problem)
目前使用 `SetReadDeadline` 来实现闲置超时，但这在 Go 的零拷贝 (`Splice`) 场景下会被解析为“总传输时长限制”。如果图片过大或传输较慢，下载会因触及绝对死线而被强行中断，导致“分段加载”。

## 2. 改进规格 (Specs)
- **解耦超时与转发**: 移除 `io.Copy` 路径上的 `SetReadDeadline`。
- **引入观测哨兵 (Watchdog)**:
    - 转发逻辑保持 `Native Splice` 性能。
    - 使用一个原子时间戳记录“最后活动时间”。
    - 使用后台协程扫描该时间戳，实现真正的“空闲超时”关闭。

## 3. 交付价值
- **完美加载**: 彻底解决大文件或慢网下的分段加载问题。
- **性能飞跃**: 移除频繁的 `SetDeadline` 系统调用，降低 CPU 负载。
