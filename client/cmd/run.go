package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	clientapp "noah/client/app"

	"noah/pkg/app"
	"noah/pkg/config"

	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
)

var (
	programName = "noah-cli"
	pidFile     = ".pid"
	logFile     = ".log"
	cntxt       = &daemon.Context{
		PidFileName: pidFile,
		PidFilePerm: 0644,
		LogFileName: logFile,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	daemonFlag   bool
	configBase64 string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Run: func(cmd *cobra.Command, args []string) {
		if daemonFlag {
			runAsDaemon()
		} else {
			run()
		}
	},
}

func init() {
	runCmd.Flags().BoolVarP(&daemonFlag, "daemon", "d", false, "Run in background (daemon mode)")
	runCmd.Flags().StringVarP(&configBase64, "config", "c", "", "config")
	rootCmd.AddCommand(runCmd)
}

func run() {
	var cfg config.ClientConfig
	b, err := base64.StdEncoding.DecodeString(configBase64)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		panic(err)
	}
	cli := clientapp.NewClient(&cfg)

	a := app.NewApp(cli)
	a.Run()
}

func runAsDaemon() {
	if p, running := isRunning(); running {
		fmt.Printf("%s is already running with PID: %d\n", programName, p.Pid)
		return
	}

	// 启动守护进程
	child, err := cntxt.Reborn()
	if err != nil {
		log.Fatalf("Failed to start daemon: %v", err)
	}

	if child != nil {
		// 父进程，显示启动信息后退出
		fmt.Printf("Daemon started with PID: %d\n", child.Pid)
		fmt.Printf("Logs will be written to: %s\n", logFile)
		return
	}

	// 子进程继续运行
	defer cntxt.Release()
	run()
}
