# Design: Fix GitHub Action Trigger

## 1. 触发器重构 (Trigger Refactoring)
```yaml
on:
  push:
    branches: [main] # 提交代码到 main 时触发
    tags: ['v*']     # 推送版本号时触发
  workflow_dispatch:  # 允许手动手动执行
```

## 2. 逻辑分流 (Conditional Execution)
为了防止普通提交也去创建 Release，我们需要在发布步骤加上条件判断：
- **Create Release**: `if: startsWith(github.ref, 'refs/tags/')`
- **Cleanup**: `if: startsWith(github.ref, 'refs/tags/')`

## 3. 实现细节
- 修改 `.github/workflows/release.yml`。
- 保持原有的静态编译逻辑不变。
