package external

import (
	"encoding/json"
	"fmt"
	"io"
	"leetFalls/internal/domain/models"
	"log/slog"
	"net/http"
)

type GravityFallsAPI struct {
	url string
}

func NewGravityFallsAPI(url string) *GravityFallsAPI {
	return &GravityFallsAPI{url: url}
}

// Associates a specific character with a user via an external API
func (ext *GravityFallsAPI) SetUser(user *models.User) error {
	path := fmt.Sprintf("%s/characters/%d", ext.url, user.ID)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer req.Body.Close()

	data, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read request: %w", err)
	}

	// Temporary variable for parsing
	var character models.Character
	err = json.Unmarshal(data, &character)
	if err != nil {
		return fmt.Errorf("failed to unmarshall character data: %w", err)
	}

	// Writing data in user
	user.Name = character.Name
	user.ImageURL = character.ImageURL

	slog.Info("Parsed response from external API: ", "message", character)
	return nil
}
