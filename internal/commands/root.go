package commands

import (
	"fmt"
	"os"

	"github.com/puyokura/revolt-file-uploader/internal/config"
	"github.com/spf13/cobra"
)

var (
	Token string
)

var rootCmd = &cobra.Command{
	Use:   "rev-up",
	Short: "Revolt File Uploader",
	Long:  `A CLI tool to upload large files to Revolt by splitting them.`,
	Run: func(cmd *cobra.Command, args []string) {
		if Token != "" {
			if err := config.SaveToken(Token); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to save token: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Token saved successfully!")
		} else {
			cmd.Help()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&Token, "token", "t", "", "Revolt Bot/Session Token (or set REVOLT_TOKEN env var)")
}
