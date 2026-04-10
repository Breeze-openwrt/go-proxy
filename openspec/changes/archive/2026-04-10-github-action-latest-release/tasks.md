# Tasks: GitHub Action Latest Release

## Phase 1: 自动化构建流程 (Step 1)
- [x] Task 1: 创建 `.github/workflows/release.yml` 基础脚手架
- [x] Task 2: 实现静态编译与资产打包 (Binaray + Config)
- [x] 验证：本地模拟运行编译指令成功。

## Phase 2: 发布与清理逻辑 (Step 2)
- [x] Task 3: 配置 `softprops/action-gh-release` 自动发布
- [x] Task 4: 实现 **旧版清理脚本** (GH CLI 整合)
- [x] 验证：检查脚本逻辑，确保它不会删除正在处理的当前 Tag。

## Phase 3: 权限与总装
- [x] Task 5: 配置 GITHUB_TOKEN 的写入权限
- [x] 验证：完成 Workflow 全量编写，准备好推送到仓库。
