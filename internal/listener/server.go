package listener

import (
	"context"
	"net"
	"syscall"
	"time"

	"github.com/dan/go-sni-proxy/internal/logger"
	"github.com/dan/go-sni-proxy/internal/proxy"
	"github.com/dan/go-sni-proxy/internal/router"
	"github.com/dan/go-sni-proxy/internal/sniffer"
)

// Server 是我们的监听控制中心
type Server struct {
	Addr             string
	NetworkInterface string
	Router           *router.Router
	Pool             *proxy.BackendPool // 新增：连接池
}

// Start 启动监听服务
func (s *Server) Start() error {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var controlErr error
			err := c.Control(func(fd uintptr) {
				if s.NetworkInterface != "" {
					// 仅适用于 Linux 环境的 SO_BINDTODEVICE
					controlErr = syscall.SetsockoptString(int(fd), syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, s.NetworkInterface)
				}
			})
			if err != nil {
				return err
			}
			return controlErr
		},
	}

	l, err := lc.Listen(context.Background(), "tcp", s.Addr)
	if err != nil {
		return err
	}
	return s.Serve(l)
}

// Serve 在指定的 Listener 上启动服务
func (s *Server) Serve(l net.Listener) error {
	defer l.Close()

	logger.Info("Proxy server started on %s", l.Addr().String())

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Accept error: %v", err)
			continue
		}

		go s.HandleConn(conn)
	}
}

func (s *Server) HandleConn(conn net.Conn) {
	start := time.Now()
	defer conn.Close()

	// 1. 探测域名 (SNI)
	sniffStart := time.Now()
	result, err := sniffer.Sniff(conn)
	sniffDuration := time.Since(sniffStart)
	if err != nil {
		logger.Debug("Sniff failed (maybe not TLS) from %v: %v", conn.RemoteAddr(), err)
		return
	}

	// 2. 查找路由
	matchStart := time.Now()
	route, ok := s.Router.Lookup(result.Domain)
	matchDuration := time.Since(matchStart)
	if !ok {
		logger.Warn("No route found for domain: %s (from %v)", result.Domain, conn.RemoteAddr())
		return
	}

	// 3. 连接后端 (使用连接池)
	dialStart := time.Now()
	backendConn, err := s.Pool.Get(route.Addr)
	dialDuration := time.Since(dialStart)
	if err != nil {
		logger.Error("Failed to connect to backend %s: %v", route.Addr, err)
		return
	}
	defer s.Pool.Put(route.Addr, backendConn)

	// 4. 子级流水报告 (Debug 级别)
	logger.Debug("Handled request: domain=%s, sniff=%dms, match=%dns, dial=%dms, total=%dms",
		result.Domain,
		sniffDuration.Milliseconds(),
		matchDuration.Nanoseconds(),
		dialDuration.Milliseconds(),
		time.Since(start).Milliseconds(),
	)

	logger.Info("Forwarding traffic: %s -> %s", result.Domain, route.Addr)
	proxy.Forward(result.Conn, backendConn, time.Duration(route.IdleTimeout)*time.Second)
}
