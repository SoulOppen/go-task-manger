package cmd

import (
	"github.com/SoulOppen/task-manager-go/internal/auth"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Inicia sesion de usuario",
	Long:  "Permite iniciar sesion o registrar un usuario con --signup.",
	Run: func(cmd *cobra.Command, args []string) {
		signup, _ := cmd.Flags().GetBool("signup")
		if signup {
			auth.SignUp()
			return
		}
		auth.Login()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().BoolP("signup", "s", false, "registrar usuario nuevo")
}
