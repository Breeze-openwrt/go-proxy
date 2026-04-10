package sniffer

import (
	"encoding/hex"
	"io"
	"net"
	"testing"
)

// 一个包含 SNI 为 "example.com" 的完整 TLS ClientHello 报文
// 一个包含 SNI 为 "example.com" 的完整 TLS ClientHello 报文
// 一个包含 SNI 为 "example.com" 的完整 TLS ClientHello 报文 (修正了长度字段)
const testClientHello = "16030100c0010000bc0303ec12dd1762a663b0ede9f2f86538545e7078631dd1577c3a072120e7e231000020e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e2e200061301130213030100006d00000010000e00000b6578616d706c652e636f6d000a00080006001d00170018000b00020100000d0012001004030804040105030805050108060601002b0003020304002d00020101003300260024001d0020358072d636ad5ad3028d27949510167230055938855b146a8123bc48f886444e000000000000000000000000"

func TestParseServerName(t *testing.T) {
	data, _ := hex.DecodeString(testClientHello)
	// 跳过 5 字节的 TLS Header
	domain, err := parseServerName(data[5:])
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	expected := "example.com"
	if domain != expected {
		t.Errorf("Expected %s, got %s", expected, domain)
	}
}

// 模拟连接测试
type mockConn struct {
	net.Conn
	data []byte
	off  int
}

func (m *mockConn) Read(b []byte) (int, error) {
	if m.off >= len(m.data) {
		return 0, io.EOF
	}
	n := copy(b, m.data[m.off:])
	m.off += n
	return n, nil
}

func (m *mockConn) Close() error { return nil }

func TestSniff(t *testing.T) {
	data, _ := hex.DecodeString(testClientHello)
	mock := &mockConn{data: data}
	
	result, err := Sniff(mock)
	if err != nil {
		t.Fatalf("Sniff failed: %v", err)
	}

	if result.Domain != "example.com" {
		t.Errorf("Expected example.com, got %s", result.Domain)
	}
}
