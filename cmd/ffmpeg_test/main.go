package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

// 启动 ffmpeg 进程
func startFFmpeg() *exec.Cmd {
	// 这里指定你的 ffmpeg 命令和参数
	cmd := exec.Command("ffmpeg", "-i", "input.mp4", "output.mp4")

	// 将输出和错误重定向到标准输出
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatalf("启动 ffmpeg 失败: %v", err)
	}

	fmt.Println("ffmpeg 已启动，PID:", cmd.Process.Pid)
	return cmd
}

// 监控 ffmpeg 进程
func monitorFFmpeg(cmd *exec.Cmd) {
	for {
		// 检查进程是否退出
		err := cmd.Process.Signal(nil)
		if err != nil {
			// 进程已经退出
			fmt.Println("ffmpeg 已停止，尝试重启...")
			cmd = startFFmpeg() // 重新启动 ffmpeg
		}

		// 每隔 5 秒检查一次进程状态
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// 启动 ffmpeg
	cmd := startFFmpeg()

	// 使用 goroutine 监控 ffmpeg
	go monitorFFmpeg(cmd)

	// 主程序继续运行，可以添加其他逻辑
	select {} // 阻塞主线程
}
