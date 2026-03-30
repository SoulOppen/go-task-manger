package cmd

import (
	"database/sql"

	"github.com/SoulOppen/task-manager-go/internal/auth"
	"github.com/SoulOppen/task-manager-go/internal/db"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Iniciar sesion",
	Long:  "Abre sesion o registra con --signup (MySQL activo).",
	RunE: func(cmd *cobra.Command, args []string) error {
		signup, err := cmd.Flags().GetBool("signup")
		if err != nil {
			return err
		}
		return db.WithDB(cmd.Context(), func(d *sql.DB) error {
			store := auth.NewMySQLUserStore(d)
			if signup {
				return auth.RunSignUp(cmd.Context(), store, cmd.InOrStdin(), cmd.OutOrStdout())
			}
			return auth.RunLogin(cmd.Context(), store, cmd.InOrStdin(), cmd.OutOrStdout())
		})
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().BoolP("signup", "s", false, "registrar usuario nuevo")
	loginCmd.SilenceUsage = true
}
