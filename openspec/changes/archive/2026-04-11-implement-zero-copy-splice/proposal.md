# Proposal: Implement High-Performance Zero-copy Forwarding (Splice)

## 1. 问题定位 (Problem)
当前的 `Forward` 逻辑使用 `Read/Write` 循环结合用户态 Buffer进行数据交换。
- **上下文切换**: 每次读写都要在用户态与内核态之间切换。
- **内存开销**: 数据被拷贝到用户内存中，增加了 CPU 缓存压力。
- **在高并发/大带宽场景下**，CPU 会成为系统的首要瓶颈。

## 2. 改进规格 (Specs)
- **内核加速**: 利用 Linux `splice(2)` 系统调用实现零拷贝转发。
- **自动回退**: 如果底层连接不支持 `splice` (如非 TCP 或非 Linux 环境)，自动降级到高性能 `Read/Write`。
- **集成优化**: 保持现有的 `idleTimeout` 超时控制能力。

## 3. 交付价值
- 显著降低单次请求的 CPU 循环指令消耗。
- 极大幅度提升万兆网卡环境下的吞吐极限。
