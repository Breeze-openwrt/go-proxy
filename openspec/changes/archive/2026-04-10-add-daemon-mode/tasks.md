# Tasks: Add Daemon Mode

## Phase 1: 基础设施与 PID 管理 (Step 1)
- [ ] Task 1: 实现 PID 文件工具类 (Check, Read, Write, Delete)
- [ ] Task 2: 实现 PID 冲突解决逻辑 (杀掉并替换旧进程)
- [ ] 验证：手动运行两个实例，确认第二个能成功杀掉第一个。

## Phase 2: 后台化与 Fork (Step 2)
- [ ] Task 3: 实现基于 Re-exec 的后台 Fork 逻辑
- [ ] Task 4: 在 `main.go` 中整合 `-d` 参数解析
- [ ] 验证：执行 `./sni-proxy -d` 后，控制台立即返回，但 `ps aux` 显示进程在后台。

## Phase 3: 信号与收尾
- [ ] Task 5: 优化信号处理，确保关闭时清理 PID 文件
- [ ] 验证：`kill <pid>` 后 PID 文件被自动移除。
