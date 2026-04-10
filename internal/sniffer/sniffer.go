package sniffer

import (
	"bufio"
	"errors"
	"io"
	"net"
)

// SniffResult 包含提取出的域名和重新包装后的连接
type SniffResult struct {
	Domain string
	Conn   net.Conn
}

// Sniff 解析 TLS ClientHello 并提取 SNI
func Sniff(conn net.Conn) (*SniffResult, error) {
	// 1. 创建一个带缓存的读取器，这样我们可以“窥视”内容
	br := bufio.NewReader(conn)

	// 2. 尝试读取 SNI
	// 我们只需要解析 ClientHello，逻辑比较复杂，这里实现一个轻量级解析
	domain, err := peekSNI(br)
	if err != nil {
		return nil, err
	}

	// 3. 将已经读取的缓存部分和原始连接重新包装，确保后续模块能读到完整数据
	return &SniffResult{
		Domain: domain,
		Conn:   &bufferedConn{Conn: conn, br: br},
	}, nil
}

// peekSNI 从缓冲读取器中解析 SNI
func peekSNI(br *bufio.Reader) (string, error) {
	// TLS Record Header (5 bytes)
	header, err := br.Peek(5)
	if err != nil {
		return "", err
	}

	if header[0] != 0x16 {
		return "", errors.New("not a TLS handshake")
	}

	length := int(header[3])<<8 | int(header[4])
	
	// 这里 Peek 出完整的 Record
	data, err := br.Peek(5 + length)
	if err != nil {
		return "", err
	}

	// 确保传给 parseServerName 的是 Handshake Data (跳过 5 字节 Header)
	return parseServerName(data[5:])
}

// bufferedConn 重新包装了 net.Conn，使得 Read 操作先从 bufio 读取
type bufferedConn struct {
	net.Conn
	br *bufio.Reader
}

func (c *bufferedConn) Read(b []byte) (int, error) {
	return c.br.Read(b)
}

// WriteTo 实现了 io.WriterTo 接口
func (c *bufferedConn) WriteTo(w io.Writer) (int64, error) {
	// 先倒出 bufio 缓冲区里的数据
	return c.br.WriteTo(w)
}

// parseServerName 是一个简化的 TLS 扩展解析逻辑
// 注意：实际生产中建议使用更完备的库，这里为了演示纯手写关键逻辑
func parseServerName(data []byte) (string, error) {
	if len(data) < 38 {
		return "", errors.New("handshake too short")
	}

	// 1. Skip Handshake Type(1), Length(3), Version(2), Random(32)
	pos := 1 + 3 + 2 + 32

	// 2. Skip Session ID
	if pos+1 > len(data) { return "", errors.New("invalid session id len") }
	sessionIDLen := int(data[pos])
	pos += 1 + sessionIDLen

	// 3. Skip Cipher Suites
	if pos+2 > len(data) { return "", errors.New("invalid cipher suites len") }
	cipherSuiteLen := int(data[pos])<<8 | int(data[pos+1])
	pos += 2 + cipherSuiteLen

	// 4. Skip Compression Methods
	if pos+1 > len(data) { return "", errors.New("invalid compression len") }
	compressionLen := int(data[pos])
	pos += 1 + compressionLen

	// 5. Extensions
	if pos+2 > len(data) { return "", errors.New("no extensions") }
	extensionsLen := int(data[pos])<<8 | int(data[pos+1])
	pos += 2
	
	if pos+extensionsLen > len(data) {
		return "", errors.New("extensions out of bounds")
	}
	end := pos + extensionsLen

	for pos + 4 <= end {
		extType := int(data[pos])<<8 | int(data[pos+1])
		extLen := int(data[pos+2])<<8 | int(data[pos+3])
		pos += 4

		// 0x0000 是 Server Name 扩展
		if extType == 0x0000 {
			if pos + 2 > end { break }
			// Server Name List Length (2)
			pos += 2
			if pos + 3 > end { break }
			// Name Type(1), Name Length(2)
			nameType := data[pos]
			nameLen := int(data[pos+1])<<8 | int(data[pos+2])
			pos += 3
			if nameType == 0 && pos + nameLen <= end {
				return string(data[pos : pos+nameLen]), nil
			}
		}
		pos += extLen
	}

	return "", errors.New("SNI not found")
}
