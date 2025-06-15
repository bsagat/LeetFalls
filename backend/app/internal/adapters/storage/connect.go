package storage

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Like minIO :)
type GonIO struct {
	url            string
	commentDirPath string // comment images directory path
	postsDirPath   string // post images directory path
}

func NewGonIOStrorage(host, port string) *GonIO {
	return &GonIO{url: host + ":" + port, commentDirPath: "comments", postsDirPath: "posts"}
}

func (storage *GonIO) CheckHealth() error {
	resp, err := http.Get(storage.url + "/PING")
	if err != nil {
		return fmt.Errorf("failed to ping storage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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
