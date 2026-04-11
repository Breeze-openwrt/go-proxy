package proxy

import (
	"context"
	"io"
	"net"
	"sync"
	"time"
)

// Forward 在两个连接之间进行双向数据拷贝，并使用哨兵协程维护闲置超时
func Forward(conn1, conn2 net.Conn, idleTimeout time.Duration) {
	// 使用 context 来协调哨兵协程的生命周期
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if idleTimeout > 0 {
		// 启动哨兵协程：只要转发还在继续，就持续推后截止日期
		go func() {
			ticker := time.NewTicker(idleTimeout / 2)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					// 刷新两端的截止日期，赋予连接“生命值”
					now := time.Now()
					_ = conn1.SetReadDeadline(now.Add(idleTimeout))
					_ = conn2.SetReadDeadline(now.Add(idleTimeout))
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// 开启双向转发 (利用 io.Copy 触发 Linux Splice 零拷贝)
	go func() {
		defer wg.Done()
		_, _ = io.Copy(conn1, conn2)
		cancel() // 任意一端结束，立即通知另一端和哨兵退出
	}()

	go func() {
		defer wg.Done()
		_, _ = io.Copy(conn2, conn1)
		cancel()
	}()

	wg.Wait()
}
