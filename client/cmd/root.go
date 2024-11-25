package cmd

import (
	"fmt"
	"noah/client/app"
	"noah/client/app/environment"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	daemon bool
	host   string
	port   int
)

var rootCmd = &cobra.Command{
	Use:   "noah client",
	Short: "Noah client command-line tool",
	Long:  `Noah client is a powerful command-line tool for interacting with Noah server.`,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Noah client",
	Long:  `Run the Noah client with specified server addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		if daemon {
			// 启动后台进程
			a := os.Args[1:]
			a = RemoveString(a, "-d")
			a = RemoveString(a, "--daemon")
			cmd := exec.Command(os.Args[0], a...)
			if err := cmd.Start(); err != nil {
				fmt.Println("Error starting daemon:", err)
				return
			}
			fmt.Printf("Daemon started with PID: %d\n", cmd.Process.Pid)
			return
		}
		env := environment.Environment{
			Server: environment.ServerConfig{
				Host: host,
				Port: port,
			},
		}
		c := app.NewClient(&env)
		c.Start()
	},
}

func RemoveString(slice []string, s string) []string {
	var result []string
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolVarP(&daemon, "daemon", "d", false, "Run as a daemon")
	runCmd.Flags().StringVarP(&host, "host", "", "127.0.0.1", "Server host")
	runCmd.Flags().IntVarP(&port, "tcp", "t", 8080, "Server port")
}
