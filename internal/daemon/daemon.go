package daemon

import (
	"os"
	"os/exec"
)

const EnvDaemonMarker = "GO_SNI_PROXY_DAEMON"

// Daemonize 如果处于前台且请求了后台模式，则派生子进程并退出父进程
func Daemonize() error {
	// 检查是否已经是子进程（避免无限循环）
	if os.Getenv(EnvDaemonMarker) == "1" {
		return nil
	}

	// 获取当前可执行程序路径
	path, err := os.Executable()
	if err != nil {
		return err
	}

	// 准备参数：去掉其中的 -d 标志，防止子进程再次触发 Daemonize
	var args []string
	for _, arg := range os.Args[1:] {
		if arg == "-d" {
			continue
		}
		args = append(args, arg)
	}

	cmd := exec.Command(path, args...)
	
	// 设置环境变量标记子进程
	cmd.Env = append(os.Environ(), EnvDaemonMarker+"=1")
	
	// 重定向标准输出/错误（子进程将独立运行）
	// 注意：真正工作时的日志重定向会在 main 中通过日志库处理
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	// 启动子进程
	if err := cmd.Start(); err != nil {
		return err
	}

	// 父进程功成身退
	os.Exit(0)
	return nil
}
