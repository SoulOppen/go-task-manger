/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/SoulOppen/task-manager-go/internal/config"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version",
	Long:  `The actual version of your app.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.AppVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
