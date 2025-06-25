package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"leetFalls/internal/domain/models"
	"log/slog"
	"net/http"
)

type GravityFallsAPI struct {
	url    string
	s3port string
}

func NewGravityFallsAPI(url string, s3port string) *GravityFallsAPI {
	return &GravityFallsAPI{url: url, s3port: s3port}
}

// Associates a specific character with a user via an external API
func (ext *GravityFallsAPI) SetUser(user *models.User) error {
	count, err := ext.AvatarCount()
	if err != nil {
		return fmt.Errorf("failed to get avatar count: %w", err)
	}

	characterID := user.ID
	if user.ID > count {
		characterID = user.ID % count
	}

	url := fmt.Sprintf("%s/characters/%d", ext.url, characterID)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var character models.Character
	if err := json.Unmarshal(data, &character); err != nil {
		return fmt.Errorf("failed to unmarshal character data: %w", err)
	}

	user.Name = character.Name
	user.ImageURL = "http://127.0.0.1:" + ext.s3port + "/objects/characters/" + character.ImageURL

	slog.Info("Parsed character from external API", "character", character)
	return nil
}

func (ext *GravityFallsAPI) AvatarCount() (int, error) {
	url := fmt.Sprintf("%s/characters", ext.url)
	res, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to send response: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Count int `json:"Count"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if result.Count == 0 {
		return 0, errors.New("characters count is 0")
	}

	return result.Count, nil
}
