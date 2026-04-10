# Proposal: Fix Action Publish Condition

## 1. 现状回顾 (Problem)
目前的 Workflow 设计过于严格，导致通过 "Run workflow" 按钮手动触发时，跳过了 `Create Release` 和 `Cleanup` 步骤。

## 2. 改进规格 (Specs)
- **多维度触发发布**: 
    - 维持 Tag 推送触发 (`refs/tags/`)。
    - **新增**: 支持手动触发 (`workflow_dispatch`) 时的发布。
- **发布语义**:
    - 如果是手动触发，发布版本名将默认使用 `latest` 或基于分支名。
    - 依然保留“只留最新版”的清理逻辑。

## 3. 交付价值
- 允许用户在不打 Tag 的情况下，通过点击按钮直接更新 Release 页面上的二进制文件。
