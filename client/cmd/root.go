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
	daemon      bool
	httpAddress string
	tcpAddress  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "noah client",
	Short: "Noah client command-line tool",
	Long:  `Noah client is a powerful command-line tool for interacting with Noah server.`,
}

// runCmd represents the run command
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
				HttpAddr: httpAddress,
				TcpAddr:  tcpAddress,
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.client.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolVarP(&daemon, "daemon", "d", false, "Run as a daemon")
	runCmd.Flags().StringVarP(&httpAddress, "http", "", "http://127.0.0.1:8080", "HTTP server address")
	runCmd.Flags().StringVarP(&tcpAddress, "tcp", "t", "127.0.0.1:1234", "TCP server address")
}
