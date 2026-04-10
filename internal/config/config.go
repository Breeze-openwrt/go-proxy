package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Config 对应 config.jsonc 的结构
type Config struct {
	ListenAddr       string                 `json:"listen_addr"`
	NetworkInterface string                 `json:"network_interface"`
	Log              LogConfig              `json:"log"`
	Routes           map[string]RouteConfig `json:"routes"`
}

type LogConfig struct {
	Level  string `json:"level"`
	Output string `json:"output"`
}

type RouteConfig struct {
	Addr        string `json:"addr"`
	JumpStart   int    `json:"jump_start"`
	IdleTimeout int    `json:"idle_timeout"`
}

// Load 从指定路径加载并解析 JSONC 配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}

	// 剥离注释 (简单的单行注释 // 处理)
	cleanJSON := stripComments(data)

	var cfg Config
	if err := json.Unmarshal(cleanJSON, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config error: %w", err)
	}

	return &cfg, nil
}

func stripComments(data []byte) []byte {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		// 寻找 // 的位置并截断
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}
		buf.WriteString(line + "\n")
	}
	return buf.Bytes()
}
