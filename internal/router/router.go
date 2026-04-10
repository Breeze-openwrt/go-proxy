package router

import (
	"github.com/dan/go-sni-proxy/internal/config"
)

// Router 负责根据域名查找对应的后端配置
type Router struct {
	routes map[string]config.RouteConfig
}

// NewRouter 创建一个新的路由引擎
func NewRouter(routes map[string]config.RouteConfig) *Router {
	return &Router{
		routes: routes,
	}
}

// Lookup 返回域名对应的后端配置
func (r *Router) Lookup(domain string) (config.RouteConfig, bool) {
	cfg, ok := r.routes[domain]
	return cfg, ok
}
