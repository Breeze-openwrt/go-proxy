package router

import (
	"testing"
	"github.com/dan/go-sni-proxy/internal/config"
)

func TestRouter_Lookup(t *testing.T) {
	routes := map[string]config.RouteConfig{
		"example.com": {Addr: "1.1.1.1:443", IdleTimeout: 60},
		"google.com":  {Addr: "8.8.8.8:443", IdleTimeout: 30},
	}
	r := NewRouter(routes)

	tests := []struct {
		name     string
		domain   string
		wantAddr string
		wantOk   bool
	}{
		{"Exact match", "example.com", "1.1.1.1:443", true},
		{"Another match", "google.com", "8.8.8.8:443", true},
		{"No match", "baidu.com", "", false},
		{"Empty domain", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, gotOk := r.Lookup(tt.domain)
			if gotCfg.Addr != tt.wantAddr {
				t.Errorf("Lookup() gotAddr = %v, want %v", gotCfg.Addr, tt.wantAddr)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Lookup() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
