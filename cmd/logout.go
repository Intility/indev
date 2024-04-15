/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command.
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout the current account",
	Long:  `Logout the current account.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("logout called")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
