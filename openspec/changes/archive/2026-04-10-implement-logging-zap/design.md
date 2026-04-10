# Design: Implement High-Performance Logging with Zap

## 1. 核心架构 (Logger Architecture)
引入 `internal/logger` 包，作为全局日志入口。

### 日志级别映射 (Zap Leveling)
- `Trace`: `zap.DebugLevel - 1` (自定义) 或仅在文本中标记。
- `Debug`: `zap.DebugLevel`
- `Info`: `zap.InfoLevel`

## 2. 目录校验逻辑
在初始化 `lumberjack` 之前：
1. 获取 `logPath` 的父目录。
2. 使用 `os.Stat` 检查。
3. 如果不存在且不是当前目录，返回 `Error`。

## 3. 滚动配置 (Rotation)
```go
lumberjack.Logger{
    Filename:   path,
    MaxSize:    20, // 20 megabytes
    MaxBackups: 3,
    MaxAge:     28, // days
    Compress:   true,
}
```

## 4. 集成点
- `main.go`: 在加载配置后立即调用 `logger.Init`。
- 全局代码替换: 将 `log.Printf` 替换为 `logger.Info/Debug`。
