package storage

import (
	"fmt"
	"io"
	"leetFalls/internal/domain"
	"log/slog"
	"net/http"
	"time"
)

// Post image save logic
// for request use:
// localhost:9090/posts/{post_id}.{MIME type}
func (storage *GonIO) SavePostImage(post_id int, img io.Reader) error {
	mimeType, err := DetectMIME(img)
	if err != nil {
		return fmt.Errorf("failed to detect mime type: %w", err)
	}

	responseUrl := fmt.Sprintf("%s/%s/%d.%s", storage.url, storage.postsDirPath, post_id, mimeType)
	req, err := http.NewRequest("PUT", responseUrl, img)
	if err != nil {
		return fmt.Errorf("failed to create response: %w", err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	slog.Info("Response from s3 storage server: ", "message", string(body))

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// Comment image save logic
// for request use:
// localhost:9090/comments/{post_id}_{comment_id}{MIME type}
// Example: localhost:9090/posts/1_21.png
func (storage *GonIO) SaveCommentImage(post_id int, comment_id int, img io.Reader) error {
	mimeType, err := DetectMIME(img)
	if err != nil {
		return fmt.Errorf("failed to detect mime type: %w", err)
	}

	responseUrl := fmt.Sprintf("%s/%s/%d_%d%s", storage.url, storage.postsDirPath, post_id, comment_id, mimeType)
	req, err := http.NewRequest("PUT", responseUrl, img)
	if err != nil {
		return fmt.Errorf("failed to create response: %w", err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	slog.Info("Response from s3 storage server: ", "message", string(body))

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func DetectMIME(file io.Reader) (string, error) {
	var extension string
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	mime := http.DetectContentType(data)
	switch mime {
	case "image/jpeg":
		extension = ".jpeg"
	case "image/png":
		extension = ".png"
	case "image/bmp":
		extension = ".bmp"
	default:
		return extension, domain.ErrUnsupportMIME
	}
	return extension, nil
}
