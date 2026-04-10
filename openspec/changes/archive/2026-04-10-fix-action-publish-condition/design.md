# Design: Fix Action Publish Condition

## 1. 逻辑修正 (Conditional Merge)
修改 `if` 判断语句，将其从单纯的 Tag 检查扩展为“Tag 检查 或 手动触发”。

```yaml
if: startsWith(github.ref, 'refs/tags/') || github.event_name == 'workflow_dispatch'
```

## 2. 处理 Tag 变量
在 `Cleanup` 脚本中，如果 `github.ref_name` 指向的是分支（手动触发时），我们需要一个兜底逻辑，或者让脚本只清理除自己以外的真正 Tag。

## 3. 实现细节
- 更新 `.github/workflows/release.yml`。
