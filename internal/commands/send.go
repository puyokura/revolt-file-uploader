package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/puyokura/revolt-file-uploader/internal/api"
	"github.com/puyokura/revolt-file-uploader/internal/config"
	"github.com/puyokura/revolt-file-uploader/internal/fileutil"
	"github.com/spf13/cobra"
)

var (
	ServerID  string
	ChannelID string
)

var sendCmd = &cobra.Command{
	Use:   "send <file>",
	Short: "Upload a file to Revolt",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		// Resolve token
		if Token == "" {
			Token = os.Getenv("REVOLT_TOKEN")
		}
		if Token == "" {
			var err error
			Token, err = config.LoadToken()
			if err != nil {
				// Ignore error, maybe config doesn't exist
			}
		}
		if Token == "" {
			return fmt.Errorf("token is required. Use --token, set REVOLT_TOKEN env var, or save it using 'rev-up --token <token>'")
		}

		// Validate flags
		if ChannelID == "" {
			return fmt.Errorf("channel ID is required. Use --channel")
		}

		client := api.NewClient(Token)

		// Check file size
		info, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if info.Size() > fileutil.SplitSize {
			fmt.Println("File is larger than 15MB, splitting and uploading...")
			metadata, err := fileutil.SplitAndUpload(filePath, client)
			if err != nil {
				return err
			}

			// Save metadata to JSON file
			jsonPath := filePath + ".json"
			jsonFile, err := os.Create(jsonPath)
			if err != nil {
				return err
			}
			defer jsonFile.Close()

			encoder := json.NewEncoder(jsonFile)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(metadata); err != nil {
				return err
			}

			// Upload JSON file
			fmt.Println("Uploading metadata JSON...")
			jsonFile.Seek(0, 0) // Reset read pointer
			// Actually we need to open it again for reading or just use the struct?
			// UploadFile takes io.Reader.
			// Let's re-open to be safe or use bytes buffer.

			// Wait, we need to upload the JSON file to Revolt too?
			// The README says "parts + parts analysis json are sent".
			// So yes, we upload the JSON file.

			// Re-open file for reading
			f, _ := os.Open(jsonPath)
			defer f.Close()

			jsonID, err := client.UploadFile(f, filepath.Base(jsonPath))
			if err != nil {
				return fmt.Errorf("failed to upload metadata JSON: %w", err)
			}

			// Send message with JSON attachment
			msg := fmt.Sprintf("Uploaded split file: %s", filepath.Base(filePath))
			if err := client.SendMessage(ChannelID, msg, jsonID); err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}

			fmt.Println("Upload complete!")
			fmt.Printf("Metadata JSON ID: %s\n", jsonID)

		} else {
			fmt.Println("Uploading file...")
			f, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer f.Close()

			id, err := client.UploadFile(f, filepath.Base(filePath))
			if err != nil {
				return err
			}

			if err := client.SendMessage(ChannelID, "", id); err != nil {
				return err
			}
			fmt.Println("Upload complete!")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringVarP(&ServerID, "server", "s", "", "Server ID (optional, currently unused by API but good for future)")
	sendCmd.Flags().StringVarP(&ChannelID, "channel", "c", "", "Channel ID (required)")
}
