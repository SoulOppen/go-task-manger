package cmd

import (
	"github.com/SoulOppen/task-manager-go/internal/auth"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Cierra la sesion actual",
	Run: func(cmd *cobra.Command, args []string) {
		auth.Logout()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
