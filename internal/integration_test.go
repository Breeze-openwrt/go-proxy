package internal

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/dan/go-sni-proxy/internal/config"
	"github.com/dan/go-sni-proxy/internal/listener"
	"github.com/dan/go-sni-proxy/internal/proxy"
	"github.com/dan/go-sni-proxy/internal/router"
)

func TestIntegration(t *testing.T) {
	// 1. 启动两个模拟后端
	backend1Addr := startMockBackend(t, "RESPONSE_FROM_BACKEND_1")
	backend2Addr := startMockBackend(t, "RESPONSE_FROM_BACKEND_2")

	// 2. 启动代理服务器
	proxyAddr, l := startProxy(t, map[string]config.RouteConfig{
		"app1.com": {Addr: backend1Addr, IdleTimeout: 5},
		"app2.com": {Addr: backend2Addr, IdleTimeout: 5},
	})
	defer l.Close() 

	time.Sleep(200 * time.Millisecond)

	// 3. 构造 TLS ClientHello 报文 (SNI: app1.com)
	app1Hello := []byte{
		0x16, 0x03, 0x01, 0x00, 0x42, // Record Header
		0x01, 0x00, 0x00, 0x3e,       // Handshake Header
		0x03, 0x03,                   // Version
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
		0x00, 
		0x00, 0x02, 0x00, 0x2f, 
		0x01, 0x00,
		0x00, 0x13,
		0x00, 0x00, 0x00, 0x0f, 0x00, 0x0d, 0x00, 0x00, 0x08,
		0x61, 0x70, 0x70, 0x31, 0x2e, 0x63, 0x6f, 0x6d,
		0x00, 0x0a, 0x00, 0x00,
	}

	t.Run("Route to Backend 1", func(t *testing.T) {
		resp := sendAndReceive(t, proxyAddr, app1Hello)
		if resp != "RESPONSE_FROM_BACKEND_1" {
			t.Errorf("Expected response from backend 1, got: %q", resp)
		}
	})
}

func startMockBackend(t *testing.T, response string) string {
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.Write([]byte(response))
			conn.Close()
		}
	}()
	return l.Addr().String()
}

func startProxy(t *testing.T, routes map[string]config.RouteConfig) (string, net.Listener) {
	r := router.NewRouter(routes)
	p := proxy.NewBackendPool()
	srv := &listener.Server{
		Addr:   "127.0.0.1:0",
		Router: r,
		Pool:   p,
	}
	l, _ := net.Listen("tcp", srv.Addr)
	actualAddr := l.Addr().String()

	go func() {
		_ = srv.Serve(l)
	}()
	return actualAddr, l
}

func sendAndReceive(t *testing.T, addr string, data []byte) string {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return fmt.Sprintf("DIAL_ERROR: %v", err)
	}
	defer conn.Close()
	
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Sprintf("WRITE_ERROR: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Sprintf("READ_ERROR: %v", err)
	}
	return string(buf[:n])
}
