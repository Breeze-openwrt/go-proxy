# Proposal: Add Daemon Mode

## 1. 目标 (Objectives)
增强 `go-sni-proxy` 的生产运维能力，通过 `-d` 参数支持后台运行，并提供自动化的进程实例冲突解决机制。

## 2. 核心规格 (Specs)
- **命令行参数**: `-d` (Daemon)。
- **PID 文件**:
    - 路径优先级: `/run/sni-proxy.pid` > `/var/run/sni-proxy.pid` > `./sni-proxy.pid`。
    - 内容: 当前运行进程的数字 ID。
- **单例限制**:
    - 启动时若检测到 PID 文件且进程活跃，自动发送信号杀掉旧进程。
    - 成功杀掉或发现文件失效后，写入当前 PID 并继续启动。
- **后台机制**: 后台运行时将 `stdout/stderr` 重定向至日志，并脱离终端控制台。

## 3. 交付价值
- 使程序适配生产环境的长期运行需求。
- 简化实例更新流程（无需手动 `ps | grep` 并 `kill`）。
