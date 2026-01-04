package cmd

import (
	"fmt"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop",
	Run: func(cmd *cobra.Command, args []string) {
		stop()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stop() {
	process, running := isRunning()
	if !running {
		fmt.Printf("%s is not running\n", programName)
		return
	}

	// 尝试优雅关闭
	if err := process.Signal(syscall.SIGTERM); err != nil {
		fmt.Printf("Failed to send SIGTERM to process %d: %v\n", process.Pid, err)
		return
	}

	// 等待进程退出
	fmt.Printf("Sent stop signal to process %d\n", process.Pid)

	// 等待几秒钟检查进程是否已经退出
	for i := 0; i < 20; i++ {
		time.Sleep(500 * time.Millisecond)
		if _, err := cntxt.Search(); err != nil {
			fmt.Println("Process stopped successfully")
			return
		}
	}

	// 如果进程还在运行，强制杀死
	fmt.Println("Process did not stop gracefully, forcing termination...")
	if err := process.Kill(); err != nil {
		fmt.Printf("Failed to kill process: %v\n", err)
	} else {
		fmt.Println("Process terminated forcefully")
	}
}
