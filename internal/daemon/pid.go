package daemon

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// HandlePIDLock 确保只有一个实例运行，如果发现旧实例则杀掉它
func HandlePIDLock(pidPath string) error {
	if pidPath == "" {
		return nil
	}

	// 1. 检查是否存在旧的 PID 文件
	data, err := os.ReadFile(pidPath)
	if err == nil {
		oldPid, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
		if parseErr == nil {
			// 2. 检查旧进程是否活跃
			process, err := os.FindProcess(oldPid)
			if err == nil {
				// 发送信号 0 检查进程是否存在
				if err := process.Signal(syscall.Signal(0)); err == nil {
					fmt.Printf("Detected existing process (PID %d). Terminating it...\n", oldPid)
					
					// 尝试优雅关闭
					process.Signal(syscall.SIGTERM)
					
					// 等待一段时间
					terminated := false
					for i := 0; i < 10; i++ {
						time.Sleep(200 * time.Millisecond)
						if err := process.Signal(syscall.Signal(0)); err != nil {
							terminated = true
							break
						}
					}

					// 如果还没退，暴力杀掉
					if !terminated {
						fmt.Printf("Process %d still alive. Force killing...\n", oldPid)
						process.Kill()
					}
				}
			}
		}
	}

	// 3. 写入当前进程的 PID
	currentPid := os.Getpid()
	return os.WriteFile(pidPath, []byte(strconv.Itoa(currentPid)), 0644)
}

// CleanPID 移除 PID 文件
func CleanPID(pidPath string) {
	if pidPath != "" {
		os.Remove(pidPath)
	}
}
