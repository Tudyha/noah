package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "status",
	Run: func(cmd *cobra.Command, args []string) {
		if p, running := isRunning(); running {
			fmt.Printf("%s is already running with PID: %d\n", programName, p.Pid)
			return
		} else {
			fmt.Printf("%s is not running\n", programName)
		}
	},
}

func isRunning() (*os.Process, bool) {
	process, err := cntxt.Search()
	if err != nil || process == nil {
		return nil, false
	}
	return process, true
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
