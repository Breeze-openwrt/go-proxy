# Tasks: Go SNI Proxy 实施路线图

## Phase 1: 基础设施 (Step 1-2)
- [x] Step 1: 实现基础 TCP 监听器 (listener)
- [x] Step 2: 集成结构化日志 (logger)
- [x] 验证：使用 telnet 连接，查看后台日志输出。

## Phase 2: SNI 协议识别 (Step 3-4)
- [x] Step 3: 实现 SNI 探测器 (sniffer)
- [x] Step 4: 实现静态路由匹配 (router)
- [x] 验证：使用 `curl -k --resolve` 测试能否正确提取域名。

## Phase 3: 数据搬运与路由 (Step 5-6)
- [x] Step 5: 实现双向数据对换器 (proxy)
- [x] Step 6: 配置文件加载 with router
- [x] 验证：完整闭环测试，通过代理访问指定的 HTTPS 服务。
