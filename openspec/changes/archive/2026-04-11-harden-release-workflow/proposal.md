# Proposal: Harden Release Workflow

## 1. 问题定位 (Problem)
- **并发冲突**: 同时推送多个 Tag 时，多个流水线实例并行运行，导致清理逻辑互相干扰，无法实现“只保留一个”的目标。
- **资源丢失**: 二进制文件未成功附加到 Release，可能是由于并发 API 调用失败。

## 2. 改进规格 (Specs)
- **串行化发布**: 引入 `concurrency` 机制。同一时间只有一个 Release Job 在运行，后续任务会自动排队或取消。
- **全局清理**: 在清理 Step 中，不仅通过 `gh release list` 获取，还要通过 `gh api` 确保捕捉到所有悬挂的 Release。
- **资产校验**: 在发布 Step 增加资产存在性校验。

## 3. 交付价值
- 确保仓库 Release 列表永远只有最新的一个正式版本。
- 确保发布的版本 100% 携带可用的 Linux 二进制文件。
