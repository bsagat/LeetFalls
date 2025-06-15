package storage

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Like minIO :)
type GonIO struct {
	url            string
	commentDirPath string // comment images directory path
	postsDirPath   string // post images directory path
}

func InitStorage(host, port string) (*GonIO, error) {
	storage := &GonIO{url: host + ":" + port, commentDirPath: "comments", postsDirPath: "posts"}
	if err := storage.CheckHealth(); err != nil {
		slog.Error("Failed to send PING message: ", "error", err.Error())
		return nil, err
	}

	if err := storage.InitBuckets(); err != nil {
		slog.Error("Failed to init buckets: ", "error", err.Error())
		return nil, err
	}

	return storage, nil
}

// Initializes comments and posts buckets
func (storage *GonIO) InitBuckets() error {
	dirPaths := []string{storage.postsDirPath, storage.commentDirPath}

	for _, path := range dirPaths {
		url := fmt.Sprintf("%s/%s", storage.url, path)
		req, err := http.NewRequest("PUT", url, nil)
		if err != nil {
			return fmt.Errorf("failed to create new request: %w", err)
		}
		defer req.Body.Close()

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send response: %w", err)
		}
		defer resp.Body.Close()

		// error can be raised if bucket is already exist...
		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("unexpected status code: %d, expected: %d", resp.StatusCode, http.StatusCreated)
		}
	}
	return nil
}

func (storage *GonIO) CheckHealth() error {
	resp, err := http.Get(storage.url + "/PING")
	if err != nil {
		return fmt.Errorf("failed to ping storage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, expected: %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if string(body) != "PONG" {
		return errors.New("invalid health check response: expected 'PONG'")
	}

	return nil
}
