# Design: Implement Elastic Backend Connection Pool

## 1. 数据结构 (Data Structure)
```go
type idleConn struct {
    conn   net.Conn
    idleAt time.Time
}

type BackendPool struct {
    mu    sync.Mutex
    pools map[string]chan *idleConn
}
```

## 2. 核心算法流程

### Get(addr)
1. 进入 `select` 尝试从 `pools[addr]` 弹出 `item`。
2. **如果拿到 item**：
    - 检查 `time.Since(item.idleAt)` 是否超过阈值。
    - 发起一次非阻塞系统调用（如 `Read` 配合 1 字节 Buffer）检查 Socket 是还活着。
    - 若失效，Close 并循环第 1 步。
3. **如果池子弹空**：
    - 调用 `net.DialTimeout`。

### Put(addr, conn)
1. 检查 `conn` 是否有效（如果 `Forward` 报错了，则不回收）。
2. 构建 `*idleConn`，记录当前时间为 `idleAt`。
3. 执行非阻塞发送到 `pools[addr]`：
    - `select { case ch <- item: default: conn.Close() }`

## 3. 并发安全性
使用 `sync.Mutex` 保护 Map 的创建，使用原生的 `chan` 线程安全性来管理具体的连接交换。
