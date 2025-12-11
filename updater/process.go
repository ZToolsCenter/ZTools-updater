package updater

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// WaitForExit 等待应用程序退出
func WaitForExit(appPath string) error {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("等待应用退出超时: %s", appPath)
		case <-ticker.C:
			if !isProcessRunning(appPath) {
				// 双重检查，短暂延迟
				time.Sleep(1 * time.Second)
				if !isProcessRunning(appPath) {
					return nil
				}
			}
		}
	}
}

// isProcessRunning 检查进程是否正在运行
// 在 Windows 上，尝试以独占写入模式打开文件
// 在 macOS/Linux 上，使用 ps 命令检查进程
func isProcessRunning(appPath string) bool {
	if runtime.GOOS == "windows" {
		// Windows：如果无法以写入模式打开，通常意味着已被锁定（正在运行）
		f, err := os.OpenFile(appPath, os.O_RDWR, 0666)
		if err != nil {
			return true // 假设无法打开即为正在运行
		}
		f.Close()
		return false
	} else {
		// Unix 类系统：使用 ps 命令
		cmd := exec.Command("ps", "-ax", "-o", "command")
		out, err := cmd.Output()
		if err != nil {
			fmt.Printf("运行 ps 出错: %v\n", err)
			return false
		}

		output := string(out)
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, appPath) {
				// 我们需要排除 updater 自身包含路径参数的情况
				// 简单的字符串包含可能误判，但通常运行中的应用命令以此路径开头
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, appPath) {
					return true
				}
			}
		}
		return false
	}
}

// StartApp 启动应用程序
func StartApp(appPath string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		// macOS: 如果路径指向 .app 内部的二进制文件，尝试使用 open 命令启动 .app
		// 这样可以确保 Dock 图标和上下文正确
		if strings.Contains(appPath, ".app/Contents/MacOS/") {
			split := strings.Split(appPath, ".app/")
			if len(split) > 0 {
				appBundle := split[0] + ".app"
				cmd = exec.Command("open", appBundle)
			} else {
				cmd = exec.Command(appPath)
			}
		} else {
			cmd = exec.Command(appPath)
		}
	} else {
		// Windows / Linux 直接执行
		cmd = exec.Command(appPath)
	}

	// 分离进程，以便 updater 可以退出
	if err := cmd.Start(); err != nil {
		return err
	}

	// 不等待进程结束
	return nil
}
