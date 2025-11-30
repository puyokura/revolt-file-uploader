package commands

import (
	"fmt"
	"os"

	"github.com/puyokura/revolt-file-uploader/internal/api"
	"github.com/puyokura/revolt-file-uploader/internal/config"
	"github.com/puyokura/revolt-file-uploader/internal/fileutil"
	"github.com/spf13/cobra"
)

var repairCmd = &cobra.Command{
	Use:   "repair <json-file>",
	Short: "Restore a split file from its metadata JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonPath := args[0]

		// Resolve token (optional for download if public? but usually needed)
		if Token == "" {
			Token = os.Getenv("REVOLT_TOKEN")
		}
		if Token == "" {
			var err error
			Token, err = config.LoadToken()
			if err != nil {
				// Ignore error
			}
		}
		// Note: Autumn might not need token for downloads, but let's pass it if we have it.
		// Our client struct requires it? No, NewClient takes it but we can pass empty.

		client := api.NewClient(Token)

		fmt.Printf("Restoring from %s...\n", jsonPath)
		if err := fileutil.DownloadAndJoin(jsonPath, client); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(repairCmd)
}
