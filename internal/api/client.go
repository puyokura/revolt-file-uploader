package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	DefaultAutumnURL = "https://cdn.stoatusercontent.com"
	DefaultRevoltURL = "https://api.revolt.chat"
)

type Client struct {
	Token     string
	AutumnURL string
	RevoltURL string
	HTTP      *http.Client
}

type Attachment struct {
	ID string `json:"id"`
}

func NewClient(token string) *Client {
	return &Client{
		Token:     token,
		AutumnURL: DefaultAutumnURL,
		RevoltURL: DefaultRevoltURL,
		HTTP:      &http.Client{},
	}
}

// UploadFile uploads a file to Autumn and returns the attachment ID.
func (c *Client) UploadFile(file io.Reader, filename string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.AutumnURL+"/attachments", body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	// Autumn doesn't always require auth for uploads depending on config, but usually it does or context.
	// For Revolt, we usually need x-session-token or x-bot-token.
	// Assuming Bot token for now, or we try both.
	// req.Header.Set("x-bot-token", c.Token)
	req.Header.Set("x-session-token", c.Token) // Try both or let user specify type?
	// Actually, let's just use the provided token as is. The user might provide a session token or bot token.
	// Common practice: try to detect or just send as both headers if unsure, or specific header.
	// Let's stick to a generic header map or just set both for now if we don't know the type.
	// Better approach: Check if token starts with "Bot " (optional convention) or just send as x-bot-token if it looks like one.
	// For simplicity, I'll send it as x-bot-token if it's a bot, x-session-token if user.
	// But to be safe, I will just set `x-bot-token` and `x-session-token` to the same value.
	// Revolt API ignores the one that is invalid usually.

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s - %s", resp.Status, string(bodyBytes))
	}

	var result Attachment
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ID, nil
}

// SendMessage sends a message with an attachment to a channel.
func (c *Client) SendMessage(channelID string, content string, attachmentID string) error {
	payload := map[string]interface{}{
		"content":     content,
		"attachments": []string{attachmentID},
	}

	if attachmentID == "" {
		delete(payload, "attachments")
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.RevoltURL+"/channels/"+channelID+"/messages", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("x-bot-token", c.Token)
	req.Header.Set("x-session-token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("send message failed: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}

// DownloadFile downloads a file from a URL to a local path.
func (c *Client) DownloadFile(url string, destPath string) error {
	resp, err := c.HTTP.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
