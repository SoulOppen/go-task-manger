/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/SoulOppen/task-manager-go/internal/auth"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		signup, _ := cmd.Flags().GetBool("signup")
		if signup {
			auth.SignUp()
		} else {
			auth.Login()
		}

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().BoolP("signup", "s", false, "sign for the first time")
}
