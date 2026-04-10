# Design: Initialize Remote Repository

## 1. 本地预处理 (Local Pre-processing)
- **.gitignore**: 创建并配置，忽略二进制文件 (`sni-proxy`)、日志文件 (`*.log`) 以及 `.gemini` 等临时目录。
- **本地归约**: 执行 `git init` 并创建第一个正式提交 `initial commit: high performance sni proxy with config & pooling`。

## 2. 远程仓库创建 (Remote Creation)
- **指令**: `gh repo create Breeze-openwrt/go-proxy --public --source=. --remote=origin --push`。
- **权限管理**: 依赖已完成的 `gh auth login` 认证。

## 3. 分支策略
- 默认分支设为 `main`。
