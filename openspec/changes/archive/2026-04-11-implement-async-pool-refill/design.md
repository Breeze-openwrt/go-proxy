# Design: Implement Asynchronous Zero-delay Pool Refill

## 1. 核心钩子 (Hook)
在 `Get(addr)` 方法成功返回一个连接（不论是来自缓存还是 Dial）后，增加一行：
```go
go p.triggerRefill(addr)
```

## 2. 回填逻辑 (triggerRefill)
1. **获取 Channel**: 根据 `addr` 找到对应的池。
2. **非阻塞检查**: 检查 Channel 此时是否能塞入（不需要严谨，`select default` 即可）。
3. **Dial**: 调用 `dialWithTFO`。
4. **入池**:
   ```go
   select {
   case ch <- &idleConn{...}:
   default:
       conn.Close() // 满了，关掉
   }
```

## 3. 并发治理
回填任务应该是“短小精悍”的。为了防止恶意堆积，每个地址可以限制同时进行的回填任务数，或简单依靠协程自愈。
