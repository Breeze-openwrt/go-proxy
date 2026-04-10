package proxy

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestForward(t *testing.T) {
	// 1. 基础转发测试
	t.Run("Basic Forward", func(t *testing.T) {
		c1, s1 := net.Pipe()
		c2, s2 := net.Pipe()

		msg1 := []byte("hello from client")
		msg2 := []byte("hello from server")

		go Forward(s1, s2, 0)

		// Client writes to c1, should appear at c2
		go c1.Write(msg1)
		buf := make([]byte, len(msg1))
		_, _ = c2.Read(buf)
		if !bytes.Equal(buf, msg1) {
			t.Errorf("Expected %s, got %s", msg1, buf)
		}

		// Server writes to c2, should appear at c1
		go c2.Write(msg2)
		buf = make([]byte, len(msg2))
		_, _ = c1.Read(buf)
		if !bytes.Equal(buf, msg2) {
			t.Errorf("Expected %s, got %s", msg2, buf)
		}

		c1.Close()
		c2.Close()
	})

	// 2. 并发压力测试 (验证 sync.Pool 安全性)
	t.Run("Concurrency Stress", func(t *testing.T) {
		var wg sync.WaitGroup
		numConns := 100
		wg.Add(numConns)

		for i := 0; i < numConns; i++ {
			go func(id int) {
				defer wg.Done()
				c1, s1 := net.Pipe()
				c2, s2 := net.Pipe()
				defer c1.Close()
				defer c2.Close()

				msg := []byte(fmt.Sprintf("message-from-%d", id))
				
				go Forward(s1, s2, 0)

				// Write and Verify
				go c1.Write(msg)
				buf := make([]byte, len(msg))
				_, err := c2.Read(buf)
				if err != nil || !bytes.Equal(buf, msg) {
					t.Errorf("Conn %d: Expected %s, got %s, err: %v", id, msg, buf, err)
				}
			}(i)
		}
		wg.Wait()
	})
}
