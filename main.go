package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dan/go-sni-proxy/internal/config"
	"github.com/dan/go-sni-proxy/internal/daemon"
	"github.com/dan/go-sni-proxy/internal/listener"
	"github.com/dan/go-sni-proxy/internal/proxy"
	"github.com/dan/go-sni-proxy/internal/router"
)

func main() {
	// 定义命令行参数
	configPath := flag.String("c", "config.jsonc", "Path to config file")
	daemonMode := flag.Bool("d", false, "Run in background (daemon mode)")
	flag.Parse()

	// 1. 处理后台化
	if *daemonMode {
		if err := daemon.Daemonize(); err != nil {
			log.Fatalf("Failed to daemonize: %v", err)
		}
	}

	// 2. 尝试获取 PID 锁 (优先标准位置 /run/sni-proxy.pid)
	pidPath := "/run/sni-proxy.pid"
	if err := daemon.HandlePIDLock(pidPath); err != nil {
		// 如果权限不足（非 root），使用当前目录
		pidPath = "./sni-proxy.pid"
		if err := daemon.HandlePIDLock(pidPath); err != nil {
			fmt.Printf("Warning: could not write PID file: %v\n", err)
			pidPath = ""
		}
	}

	// 3. 注册信号处理用于清理 PID
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, exiting...", sig)
		daemon.CleanPID(pidPath)
		os.Exit(0)
	}()

	// 4. 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化路由
	r := router.NewRouter(cfg.Routes)

	// 3. 初始化连接池并预热
	p := proxy.NewBackendPool()
	for _, routeCfg := range cfg.Routes {
		if routeCfg.JumpStart > 0 {
			log.Printf("Pre-heating pool for %s (count: %d)", routeCfg.Addr, routeCfg.JumpStart)
			p.PreHeat(routeCfg.Addr, routeCfg.JumpStart)
		}
	}

	// 4. 初始化并启动服务器
	srv := &listener.Server{
		Addr:             cfg.ListenAddr,
		NetworkInterface: cfg.NetworkInterface,
		Router:           r,
		Pool:             p,
	}

	log.Printf("Starting SNI Proxy on %s (Interface: %s)...", cfg.ListenAddr, cfg.NetworkInterface)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
