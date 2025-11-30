package fileutil

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/puyokura/revolt-file-uploader/internal/api"
)

type Downloader interface {
	DownloadFile(url string, destPath string) error
}

func DownloadAndJoin(metadataPath string, downloader Downloader) error {
	// Read metadata
	metaFile, err := os.Open(metadataPath)
	if err != nil {
		return err
	}
	defer metaFile.Close()

	var metadata FileMetadata
	if err := json.NewDecoder(metaFile).Decode(&metadata); err != nil {
		return err
	}

	// Sort parts just in case
	sort.Slice(metadata.Parts, func(i, j int) bool {
		return metadata.Parts[i].Index < metadata.Parts[j].Index
	})

	// Create output file
	outPath := metadata.OriginalName
	// Check if exists? Overwrite? Let's assume overwrite or fail?
	// Let's fail if exists to be safe, or just overwrite. Standard CLI usually overwrites or asks.
	// For now, overwrite.
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	tempDir := os.TempDir()

	for _, part := range metadata.Parts {
		// Construct URL if not present
		url := part.URL
		if url == "" {
			// Assume Autumn URL structure if not provided
			// But we don't know the Autumn URL here easily unless we pass it or assume default.
			// The `Downloader` (Client) has the base URL.
			// But `DownloadFile` takes a full URL.
			// We should probably construct the URL here using the default or let the Client handle it.
			// Let's assume we need to construct it.
			// "https://autumn.revolt.chat/attachments/{id}"
			// Wait, the client knows the URL.
			// Maybe we should add a method `DownloadAttachment(id, dest)` to Client?
			// But `Downloader` interface is generic.
			// Let's change `Downloader` to `AttachmentDownloader`?
			// Or just construct the URL here.
			url = api.DefaultAutumnURL + "/attachments/" + part.ID
		}

		partPath := filepath.Join(tempDir, fmt.Sprintf("%s.part%d", metadata.OriginalName, part.Index))
		fmt.Printf("Downloading part %d...\n", part.Index)
		if err := downloader.DownloadFile(url, partPath); err != nil {
			return fmt.Errorf("failed to download part %d: %w", part.Index, err)
		}

		// Append to output file
		partFile, err := os.Open(partPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, partFile)
		partFile.Close()
		os.Remove(partPath) // Clean up temp part
		if err != nil {
			return err
		}
	}

	fmt.Printf("Successfully restored %s\n", outPath)
	return nil
}
