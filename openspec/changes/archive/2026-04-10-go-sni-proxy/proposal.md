# Proposal: Go SNI Proxy (高性能 L4 透明代理)

## 1. 背景与目标
实现一个轻量级、高性能的 Go 语言 TCP 代理，支持基于 TLS SNI 字段的转发。

## 2. 核心功能
- L4 转发（不终止 TLS）
- SNI 探测（提取域名）
- 路由表映射
- 零拷贝优化

## 3. 技术栈
- 语言：Go 1.21+
- 核心库：`net`
- 并发模型：Goroutine + Channel
