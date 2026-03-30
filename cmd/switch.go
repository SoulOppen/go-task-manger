package cmd

import (
	"database/sql"

	"github.com/SoulOppen/task-manager-go/internal/auth"
	"github.com/SoulOppen/task-manager-go/internal/db"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Cambia de usuario activo",
	Long:  "Solicita credenciales y cambia la sesion al usuario autenticado (requiere MySQL).",
	RunE: func(cmd *cobra.Command, args []string) error {
		return db.WithDB(cmd.Context(), func(d *sql.DB) error {
			store := auth.NewMySQLUserStore(d)
			return auth.RunLogin(cmd.Context(), store, cmd.InOrStdin(), cmd.OutOrStdout())
		})
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.SilenceUsage = true
}
