# Design: Fix CI/CD Cleanup Logic

## 1. 提取逻辑重构
原来的方案：
```bash
OLD_TAGS=$(gh release list | awk '{print $3}')
```
**新方案**：
```bash
OLD_TAGS=$(gh release list --limit 100 --json tagName --jq '.[].tagName' | grep -v "^$CURRENT_TAG$")
```

## 2. 循环删除强化
- 增加 `continue` 逻辑或 `|| true`。
- 只有在 `OLD_TAGS` 非空时才进入循环。
