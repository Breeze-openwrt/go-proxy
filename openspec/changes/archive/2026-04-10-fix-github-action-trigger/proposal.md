# Proposal: Fix GitHub Action Trigger

## 1. 问题描述 (The Issue)
当前的 GitHub Action 仅在推送 `v*` 格式的 Tag 时触发。这导致：
1. 正常的代码提交 (Push) 不会触发编译检查。
2. GitHub 界面上没有“手动运行”按钮。

## 2. 改进规格 (Specs)
- **支持手动执行**: 引入 `workflow_dispatch` 触发器。
- **支持提交检查**: 引入 `push: branches: [main]` 触发器。
- **发布隔离**: 
    - 确保只有在推送到 Tag 时才会执行 `Create Release` 和 `Cleanup` 步骤。
    - 普通推送仅执行 `Build` 步骤以验证编译。

## 3. 交付价值
- 提供可见的 CI 反馈。
- 允许用户在不推 Tag 的情况下手动测试编译脚本。
