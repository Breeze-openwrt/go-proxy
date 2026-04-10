# Tasks: Fix GitHub Action Trigger

## Phase 1: Workflow 增强 (Step 1)
- [ ] Task 1: 在 `release.yml` 中添加 `workflow_dispatch` 与 `push: branches`
- [ ] Task 2: 为 Release 和 Cleanup 步骤添加 `if` 条件语句
- [ ] 验证：确认编译步骤对所有触发方式均有效。

## Phase 2: 同步与测试
- [ ] Task 3: 提交并推送到 GitHub
- [ ] 验证：在浏览器 Actions 页面观察是否出现了 "Run workflow" 按钮。
