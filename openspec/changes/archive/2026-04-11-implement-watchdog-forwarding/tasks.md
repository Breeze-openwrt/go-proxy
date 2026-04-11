# Tasks: Implement Watchdog Forwarding Model

- [ ] Task 1: 在 `proxy.go` 中实现 `Forward` 哨兵逻辑
- [ ] Task 2: 移除 `copyWithTimeout` 中的过时 `SetReadDeadline`
- [ ] Task 3: 在 `BackendPool` 的 `dialHappyEyeballs` 中引入并发限制（信号量）
- [ ] Task 4: 验证在慢速下载下的图片完整性
