# Proposal: Implement High-Performance Logging with Zap

## 1. 问题定位 (Problem)
目前项目使用 Go 标准库 `log` 包，缺乏：
- 结构化日志。
- 日志级别过滤 (Info/Debug/Trace)。
- 日志文件滚动管理。
- 命令行优先级的配置覆盖。

## 2. 改进规格 (Specs)
- **技术栈**: 使用 `zap` 构建核心，`lumberjack` 负责滚动。
- **配置优先级**:
    - `-vv`: TRACE (Zap Debug + Trace 标记)
    - `-v`: DEBUG
    - 无参数: 使用 `config.jsonc` 中的 `log.level` (默认 Info)。
- **安全检查**: 如果指定的日志目录不存在，程序必须报错并立即退出。
- **滚动策略**: 单个日志文件上限 20MiB。

## 3. 交付价值
- 提供工业级的可观测性。
- 保护磁盘空间，防止日志无限增长。
- 提升调试效率。
