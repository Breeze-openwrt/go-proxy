# Proposal: GitHub Action Latest Release

## 1. 目标 (Objectives)
构建一个自动化的 CI/CD 流程，确保 GitHub 仓库始终提供且仅提供最新的编译二进制文件。

## 2. 核心规格 (Specs)
- **触发机制**: 
    - 推送特定前缀的 Tag (如 `v*`) 时触发发布。
    - 推送代码到 `main` 分支时触发测试与编译验证。
- **构建目标**: Linux amd64 (静态编译以确保最大兼容性)。
- **发布策略 (Latest-Only)**: 
    - 每次成功发布新 Release 后，通过 GitHub CLI 识别并删除旧的所有 Release。
    - 同时清理本地和远程多余的旧 Tag。
- **资产附件**: `sni-proxy` 二进制文件及 `config.jsonc` 模板。

## 3. 交付价值
- 保持 Release 页面的极致整洁。
- 确保下载链接始终指向最新的稳定版本。
