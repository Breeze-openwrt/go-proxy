package proxy

import (
	"context"
	"net"
	"sync"
	"syscall"
	"time"
)

// dialWithTFO 是一个支持 TCP Fast Open 的拨号器
func dialWithTFO(ctx context.Context, addr string) (net.Conn, error) {
	d := net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				// 30 是 Linux 上的 TCP_FASTOPEN_CONNECT 选项
				// 设置为 1 开启客户端端的 TFO
				_ = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, 30, 1)
			})
		},
	}
	return d.DialContext(ctx, "tcp", addr)
}

// idleConn 包装了一个空闲连接及其进入空闲状态的时间
type idleConn struct {
	conn   net.Conn
	idleAt time.Time
}

// BackendPool 具有自愈能力的弹性后端连接池
type BackendPool struct {
	mu    sync.Mutex
	pools map[string]chan *idleConn
}

func NewBackendPool() *BackendPool {
	return &BackendPool{
		pools: make(map[string]chan *idleConn),
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
		ch = make(chan *idleConn, count+20) // 预留余量
		p.pools[addr] = ch
	}
	p.mu.Unlock()

	// 异步预热
	go func() {
		for i := 0; i < count; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			conn, err := dialWithTFO(ctx, addr)
			cancel()
			if err != nil {
				continue
			}
			select {
			case ch <- &idleConn{conn: conn, idleAt: time.Now()}:
			default:
				conn.Close()
				return
			}
		}
	}()
}

// Get 获取一个后端连接，具备过期检查和存活状态探测能力
func (p *BackendPool) Get(addr string) (net.Conn, error) {
	p.mu.Lock()
	ch, ok := p.pools[addr]
	if !ok {
		ch = make(chan *idleConn, 20)
		p.pools[addr] = ch
	}
	p.mu.Unlock()

	for {
		select {
		case item := <-ch:
			// 1. 检查是否过期 (默认 60s)
			if time.Since(item.idleAt) > 60*time.Second {
				item.conn.Close()
				continue
			}

			// 2. 检查连接是否还存活 (非破坏性探测)
			if !isAlive(item.conn) {
				item.conn.Close()
				continue
			}

			return item.conn, nil
		default:
			// 池子空了，发起新的握手 (带 TFO)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return dialWithTFO(ctx, addr)
		}
	}
}

// isAlive 使用 MSG_PEEK 检查连接是否可用，不消耗任何数据
func isAlive(conn net.Conn) bool {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return true // 非 TCP 连接默认跳过检查
	}

	raw, err := tcpConn.SyscallConn()
	if err != nil {
		return false
	}

	var closed bool
	_ = raw.Control(func(fd uintptr) {
		var b [1]byte
		// MSG_PEEK: 只窥视不读取
		// MSG_DONTWAIT: 立即返回不阻塞
		n, _, err := syscall.Recvfrom(int(fd), b[:], syscall.MSG_PEEK|syscall.MSG_DONTWAIT)
		
		// 如果读到 0 且没有错误，表示 EOF (对端已关闭)
		if n == 0 && err == nil {
			closed = true
		}
		// 如果报错且不是 EAGAIN/EWOULDBLOCK，也视为失效
		if err != nil && err != syscall.EAGAIN && err != syscall.EWOULDBLOCK {
			closed = true
		}
	})

	return !closed
}

// Put 将健康的连接归还至池中
func (p *BackendPool) Put(addr string, conn net.Conn) {
	if conn == nil {
		return
	}

	p.mu.Lock()
	ch, ok := p.pools[addr]
	if !ok {
		ch = make(chan *idleConn, 20)
		p.pools[addr] = ch
	}
	p.mu.Unlock()

	item := &idleConn{
		conn:   conn,
		idleAt: time.Now(),
	}

	// 尝试非阻塞塞回池子
	select {
	case ch <- item:
		// 归还成功
	default:
		// 池子满了，物理关闭
		conn.Close()
	}
}
