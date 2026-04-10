# Design: Fix Action Manual Tag Error

## 1. 动态标签逻辑 (Dynamic Tagging)
在 `action-gh-release` 步骤中，利用表达式来决定标签：

```yaml
tag_name: ${{ github.ref_type == 'tag' && github.ref_name || 'latest' }}
```

## 2. 发布参数修补 (Parameter Tuning)
- `make_latest: true`: 确保该 Release 在页面上始终置顶。
- `name`: 
    - 如果是 Tag: 使用 `Release ${{ github.ref_name }}`。
    - 如果是手动: 使用 `Latest Build`。

## 3. 实现细节
- 更新 `.github/workflows/release.yml`。
