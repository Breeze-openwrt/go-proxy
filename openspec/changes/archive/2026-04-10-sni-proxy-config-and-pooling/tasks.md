# Tasks: SNI Proxy Config & Pool Management

## Phase 1: 基础设施与解析 (Step 1)
- [x] Task 1: 实现 JSONC 加载器 (config.jsonc parser)
- [x] 验证：成功读取测试配置文件并打印到日志，跳过单行注释。

## Phase 2: Linux 网络增强 (Step 2)
- [x] Task 2: 实现 SO_BINDTODEVICE 网卡绑定逻辑
- [x] 验证：在 Linux 环境下，通过 `netstat` 或 `ss` 确认套接字是否绑定到特定网卡。

## Phase 3: 连接池进阶 (Step 3-4)
- [x] Task 3: 实现 BackendPool 预热机制 (Jump Start)
- [x] 验证：程序启动瞬间，后端服务能观察到预设数量 of TCP 连接。
- [x] Task 4: 实现数据转发的闲置超时 (Idle Timeout)
- [x] 验证：长连接在停止发送数据指定时间后，双方能否正确触发 `Close`。

## Phase 4: 总装与收尾 (Step 5-6)
- [x] Task 5: 整合全量配置至 main.go
- [x] Task 6: Linux amd64 综合环境验证
- [x] 验证：完整的从配置加载到流量转发的闭环测试。
