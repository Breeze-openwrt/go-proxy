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
	// 获取底层的文件描述符支持（如果支持的话）
	// Go 的 io.Copy 对 *net.TCPConn 已经有了非常好的优化
	
	for {
		if timeout > 0 {
			src.SetReadDeadline(time.Now().Add(timeout))
		}

		// 我们使用 io.CopyN 或者带限额的拷贝，以便定期有机会重置 Deadline
		// 但为了追求极致性能，我们在这里先尝试标准的 io.Copy
		// 注意：如果 io.Copy 触发了零拷贝，它会在内核态持续运行直到 EOF 或出错
		written, err := io.Copy(dst, src)
		
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 如果是超时且在这段时间内有数据写入，我们可以考虑继续（但在 Splice 模式下 io.Copy 会直接返回）
				if written > 0 {
					continue
				}
				return
			}
			return
		}
		
		if written == 0 {
			return
		}
		// 如果 written > 0 且没有错误，说明可能读到了 EOF 或者 Copy 逻辑结束了
		return
	}
}
