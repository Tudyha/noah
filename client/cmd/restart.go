/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart",
	Run: func(cmd *cobra.Command, args []string) {
		stop()
		runAsDaemon()
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
