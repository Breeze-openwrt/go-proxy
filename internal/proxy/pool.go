package proxy

import (
	"context"
	"net"
	"sync"
	"syscall"
	"time"
)

// dialHappyEyeballs 实现 Happy Eyeballs v2 (RFC 8305) 赛马竞争算法
func (p *BackendPool) dialHappyEyeballs(ctx context.Context, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return p.dialSingle(ctx, addr)
	}

	// 并发解析 IPv4 和 IPv6
	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		return p.dialSingle(ctx, addr)
	}

	// 准备赛马任务
	type result struct {
		conn net.Conn
		err  error
	}
	resCh := make(chan result, len(ips))
	raceCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 赛马控制器
	for i, ip := range ips {
		go func(targetIP net.IP) {
			targetAddr := net.JoinHostPort(targetIP.String(), port)
			conn, err := p.dialSingle(raceCtx, targetAddr)
			select {
			case resCh <- result{conn: conn, err: err}:
			case <-raceCtx.Done():
				if conn != nil {
					conn.Close()
				}
			}
		}(ip)

		// 250ms 的 Stagger 延迟启动，避免 SYN Flood
		if i < len(ips)-1 {
			select {
			case <-time.After(250 * time.Millisecond):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	// 等待赢家产生
	var lastErr error
	for range ips {
		res := <-resCh
		if res.err == nil {
			cancel() // 赢家已定，通知其他选手退赛
			return res.conn, nil
		}
		lastErr = res.err
	}

	return nil, lastErr
}

// dialSingle 是带 TFO 的单路拨号器，作为赛马的基础单元
func (p *BackendPool) dialSingle(ctx context.Context, addr string) (net.Conn, error) {
	d := net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				// TCP_FASTOPEN_CONNECT
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
			conn, err := p.dialHappyEyeballs(ctx, addr)
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

			// 3. 成功拿到连接，触发异步回填补充水位
			go p.triggerRefill(addr)
			return item.conn, nil
		default:
			// 池子空了，发起同步拨号 (带 TFO)
			// 同时触发异步回填，防止下一次请求继续落入同步拨号
			go p.triggerRefill(addr)
			
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return p.dialHappyEyeballs(ctx, addr)
		}
	}
}

// triggerRefill 尝试为指定地址补充一个空闲连接
func (p *BackendPool) triggerRefill(addr string) {
	p.mu.Lock()
	ch, ok := p.pools[addr]
	p.mu.Unlock()

	if !ok {
		return
	}

	// 提前嗅探：如果通道正忙或已满，直接退出
	if len(ch) >= cap(ch) {
		return
	}

	// 开始异步拨号
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	conn, err := p.dialHappyEyeballs(ctx, addr)
	if err != nil {
		return
	}

	// 尝试存入
	select {
	case ch <- &idleConn{
		conn:   conn,
		idleAt: time.Now(),
	}:
		// 成功存入，它将在池中等待下一个 Get
	default:
		// 存入失败（池子刚被填满），关闭该多余连接
		conn.Close()
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
