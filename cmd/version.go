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
	Short: "Mostrar version",
	Long:  "Imprime el numero de version (internal/config).",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), config.AppVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
