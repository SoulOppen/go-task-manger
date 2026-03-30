/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/SoulOppen/task-manager-go/internal/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   config.AppName(),
	Short: config.AppShortDescription(),
	Long:  config.AppLongDescription(),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	rootCmd.Use = config.AppName()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = config.AppVersion()
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}
