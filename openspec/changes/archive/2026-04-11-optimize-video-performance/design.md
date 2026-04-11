# Design: Optimize Video Performance (Turbo Mode)

## 1. Socket 选项注入计划
在 `dialSingle` 的 `Control` 回调中增加以下系统调用：
- `syscall.SetsockoptInt(fd, SOL_SOCKET, SO_RCVBUF, 4*1024*1024)`: 接收缓冲区 4MB。
- `syscall.SetsockoptInt(fd, SOL_SOCKET, SO_SNDBUF, 4*1024*1024)`: 发送缓冲区 4MB。
- `syscall.SetsockoptInt(fd, IPPROTO_TCP, TCP_NODELAY, 1)`: 禁用 Nagle。
- `syscall.SetsockoptInt(fd, IPPROTO_TCP, TCP_QUICKACK, 1)`: 强制快速确认。

## 2. 赛马竞争算法调优
- `staggerDelay` 从 250ms 下调至 **100ms**。
- 这将使系统在 200ms 内就能完成对 3 个地址的探测。

## 3. 连接池扩容
- 将 `BackendPool` 的默认单地址 Channel 容量从 20 提升至 **50**，以应对视频网站极高的并发 Range 请求。
