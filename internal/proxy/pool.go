package proxy

import (
	"net"
	"sync"
	"time"
)

// BackendPool 简单的后端连接池
type BackendPool struct {
	mu    sync.Mutex
	pools map[string]chan net.Conn
}

func NewBackendPool() *BackendPool {
	return &BackendPool{
		pools: make(map[string]chan net.Conn),
	}
}

// PreHeat 预热后端连接
func (p *BackendPool) PreHeat(addr string, count int) {
	if count <= 0 {
		return
	}

	p.mu.Lock()
	ch, ok := p.pools[addr]
	if !ok {
		ch = make(chan net.Conn, count+10) // 预留余量
		p.pools[addr] = ch
	}
	p.mu.Unlock()

	// 异步预热
	go func() {
		for i := 0; i < count; i++ {
			conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
			if err != nil {
				continue
			}
			select {
			case ch <- conn:
			default:
				conn.Close()
				return
			}
		}
	}()
}

// Get 获取一个后端连接
func (p *BackendPool) Get(addr string) (net.Conn, error) {
	p.mu.Lock()
	ch, ok := p.pools[addr]
	if !ok {
		ch = make(chan net.Conn, 10)
		p.pools[addr] = ch
	}
	p.mu.Unlock()

	select {
	case conn := <-ch:
		// 这里在实际生产中可以加一次简单的健康检查
		return conn, nil
	default:
		return net.DialTimeout("tcp", addr, 5*time.Second)
}
}

// Put 仅在需要时回缩池，在这里我们暂时不做 L4 的深度复用。
func (p *BackendPool) Put(addr string, conn net.Conn) {
	conn.Close()
}
