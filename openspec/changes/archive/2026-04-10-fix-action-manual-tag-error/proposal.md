# Proposal: Fix Action Manual Tag Error

## 1. 问题定位 (Problem)
`softprops/action-gh-release` 在 `workflow_dispatch` (手动触发) 模式下，无法从 `github.ref` 中提取合法的 Git Tag，导致发布中断并报错。

## 2. 改进规格 (Specs)
- **显式标签配置 (Explicit Tag)**: 
    - 如果是 Tag 触发，使用原有 Tag。
    - 如果是手动触发，强制指定为 `latest`。
- **发布语义优化**:
    - 允许覆盖现有 Release (`prerelease: false`)。
    - 确保手动发布的版本始终标记为 `latest` 稳定版。

## 3. 交付价值
- 消除手动触发时的报错。
- 提供一个始终可用的 `latest` 下载链接。
