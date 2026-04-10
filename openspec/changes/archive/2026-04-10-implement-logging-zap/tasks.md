# Tasks: Implement High-Performance Logging with Zap

## Phase 1: 依赖与底层实现 (Step 1)
- [ ] Task 1: 引入 `go.uber.org/zap` 和 `github.com/natefinch/lumberjack`
- [ ] Task 2: 实现 `internal/logger` 包及其初始化、目录检查逻辑
- [ ] 验证：单元测试检查目录不存在时是否报错。

## Phase 2: 集成与替换 (Step 2)
- [ ] Task 3: 在 `main.go` 中解析 `-v/-vv` 标志并初始化日志
- [ ] Task 4: 将 `main.go` 中的 `log.Printf` 替换为新 Logger
- [ ] Task 5: 将 `internal/` 下所有包的 `log` 调用迁移至 `logger`
- [ ] 验证：手动运行 `./sni-proxy -v` 查看是否带有调试输出。

## Phase 3: 滚动验证
- [ ] Task 6: 验证日志文件超过 20M 时是否发生滚动
