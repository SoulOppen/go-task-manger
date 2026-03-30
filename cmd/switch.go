package cmd

import (
	"github.com/SoulOppen/task-manager-go/internal/auth"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Cambia de usuario activo",
	Long:  "Solicita credenciales y cambia la sesion al usuario autenticado.",
	Run: func(cmd *cobra.Command, args []string) {
		auth.SwitchUser()
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
