package main

import (
	"log"
	"os"

	"github.com/dan/go-sni-proxy/internal/config"
	"github.com/dan/go-sni-proxy/internal/listener"
	"github.com/dan/go-sni-proxy/internal/proxy"
	"github.com/dan/go-sni-proxy/internal/router"
)

func main() {
	// 1. 加载配置 (默认查找 config.jsonc)
	cfgPath := "config.jsonc"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
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
