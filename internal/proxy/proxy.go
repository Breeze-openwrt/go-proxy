package proxy

import (
	"net"
	"sync"
	"time"
)

// bufferPool 用于复用固定大小的字节切片
var bufferPool = sync.Pool{
	New: func() any {
		// 每个连接分配 32KB 的缓冲区
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
		_ = conn1.Close()
	}()

	go func() {
		defer wg.Done()
		copyWithTimeout(conn2, conn1, idleTimeout)
		_ = conn2.Close()
	}()

	wg.Wait()
}

func copyWithTimeout(dst net.Conn, src net.Conn, timeout time.Duration) {
	// 从池子里借一个碗 (Buffer)
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)

	for {
		if timeout > 0 {
			src.SetReadDeadline(time.Now().Add(timeout))
		}
		nr, err := src.Read(buf)
		if nr > 0 {
			_, ew := dst.Write(buf[0:nr])
			if ew != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
}
