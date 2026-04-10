package proxy

import (
	"io"
	"net"
	"sync"
	"time"
)

// bufferPool 用于复用固定大小的字节切片 (Fallback 模式使用)
var bufferPool = sync.Pool{
	New: func() any {
		return make([]byte, 32*1024)
	},
}

// Forward 在两个连接之间进行双向数据拷贝，支持闲置超时
func Forward(conn1, conn2 net.Conn, idleTimeout time.Duration) {
	var wg sync.WaitGroup
	wg.Add(2)

	// 开启两个方向的转发
	go func() {
		defer wg.Done()
		copyWithTimeout(conn1, conn2, idleTimeout)
	}()

	go func() {
		defer wg.Done()
		copyWithTimeout(conn2, conn1, idleTimeout)
	}()

	wg.Wait()
}

func copyWithTimeout(dst net.Conn, src net.Conn, timeout time.Duration) {
	if timeout > 0 {
		src.SetReadDeadline(time.Now().Add(timeout))
	}

	// 直接执行拷贝 (在 Linux TCP 间会自动触发 Splice 零拷贝)
	// Go 的 io.Copy 会在内部循环 Read/Write 直到 EOF 或错误
	_, _ = io.Copy(dst, src)
}
