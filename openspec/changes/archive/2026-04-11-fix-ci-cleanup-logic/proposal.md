# Proposal: Fix CI/CD Cleanup Logic

## 1. 问题定位 (Problem)
当前的 GitHub Actions 清理脚本使用 `awk` 解析 `gh release list` 的纯文本输出。由于非 Latest 版本的列偏移，脚本误将时间戳作为 Tag 名进行删除，导致 `release not found` 错误并中断整个 Workflow。

## 2. 改进规格 (Specs)
- **健壮解析**: 切换到 `--json tagName` 模式，利用内置的 `jq` 过滤器提取 Tag 名。
- **容错处理**: 在执行 `gh release delete` 时增加错误容忍，避免因单个删除失败而导致整个 Job 失败。

## 3. 交付价值
- 恢复 CI/CD 自动发布的稳定性。
- 自动化保持仓库 Release 列表简洁。
