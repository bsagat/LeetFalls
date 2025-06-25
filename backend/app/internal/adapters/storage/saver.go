package storage

import (
	"bytes"
	"fmt"
	"io"
	"leetFalls/internal/domain"
	"leetFalls/internal/domain/models"
	"log/slog"
	"net/http"
	"time"
)

// SavePostImage saves post image to s3 storage
// URL request: http://127.0.0.1:9090/objects/posts/{post_id}.{MIME type}
func (storage *GonIO) SavePostImage(post *models.Post, img io.Reader) error {
	data, err := io.ReadAll(img)
	if err != nil {
		return fmt.Errorf("failed to read image: %w", err)
	}

	mimeType, err := DetectMIME(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to detect mime type: %w", err)
	}

	path := fmt.Sprintf(":%s/objects/%s/%d%s", storage.port, storage.postsDirPath, post.ID, mimeType)
	post.ImageLink = "http://127.0.0.1" + path

	req, err := http.NewRequest("PUT", storage.host+path, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
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

// SaveCommentImage saves comment image to local storage
// URL request: http://127.0.0.1:9090/objects/comments/{post_id}_{comment_id}.{MIME type}
func (storage *GonIO) SaveCommentImage(comment *models.Comment, img io.Reader) error {
	data, err := io.ReadAll(img)
	if err != nil {
		return fmt.Errorf("failed to read image: %w", err)
	}

	mimeType, err := DetectMIME(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to detect mime type: %w", err)
	}

	path := fmt.Sprintf(":%s/objects/%s/%d_%d%s", storage.port, storage.postsDirPath, comment.PostID, comment.ID, mimeType)
	comment.ImageLink = "http://127.0.0.1" + path

	req, err := http.NewRequest("PUT", storage.host+path, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
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

// DetectMIME returns MIME extension of image
func DetectMIME(file io.Reader) (string, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	mime := http.DetectContentType(data)
	switch mime {
	case "image/jpeg":
		return ".jpeg", nil
	case "image/png":
		return ".png", nil
	case "image/bmp":
		return ".bmp", nil
	default:
		return "", domain.ErrUnsupportMIME
	}
}
