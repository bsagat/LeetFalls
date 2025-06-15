package service

import (
	"crypto/rand"
	"encoding/hex"
	dbrepo "leetFalls/internal/adapters/dbRepo"
	"leetFalls/internal/adapters/external"
	"leetFalls/internal/domain/models"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type MiddlewareService struct {
	dbrepo   dbrepo.SessionRepo
	external external.GravityFallsAPI
}

// Validates the session cookie. If valid, it retrieves and returns the current user's ID.
// If the cookie is missing, invalid, or expired,
// a new session is created (for an unauthenticated/guest user) and its associated user ID is returned.
func (s *MiddlewareService) Auth(w http.ResponseWriter, cookie *http.Cookie) (int, error) {
	if cookie != nil {
		session_id, err := CheckSessionId(cookie.Value)
		if err != nil {
			slog.Error("Session id is invalid: ", "error", err.Error())
			return 0, err
		}
		userId, err := s.dbrepo.GetUserIDBySession(session_id)
		if err != nil {
			slog.Error("Failed to check session existence: ", "error", err.Error())
			return 0, err
		}

		// return id, if it exist
		if userId != 0 {
			return userId, nil
		}
	}

	// if session id not exist, we generate new user
	session_id, err := GenerateSessionID()
	if err != nil {
		slog.Error("Failed to generate session id: ", "error", err.Error())
		return 0, err
	}

	userId, err := s.CreateNewUser(session_id)
	if err != nil {
		slog.Error("Failed to create new session: ", "error", err.Error())
		return 0, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session_id,
		Path:     "/",
		Expires:  time.Now().Add(24 * 7 * time.Hour), // Expire time: 7 weeks
		HttpOnly: true,
	})
	return userId, nil
}

// Creates a new user and returns their ID.
func (s *MiddlewareService) CreateNewUser(sessionId string) (int, error) {
	var (
		user models.User
		err  error
	)

	user.Token_ID = sessionId
	user.ID, err = s.dbrepo.GetNextUserId()
	if err != nil {
		slog.Error("Failed to get next user id: ", "error", err.Error())
		return 0, err
	}

	if err := s.external.SetUser(&user); err != nil {
		slog.Error("Failed to set user information with external data: ", "error", err.Error())
		return 0, err
	}

	if err := s.dbrepo.SaveUser(user); err != nil {
		slog.Error("Failed to save user information: ", "error", err.Error())
		return 0, err
	}

	return user.ID, nil
}

// Generates session_id randomly
func GenerateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Validates session id, if it invalid, generates new one
func CheckSessionId(session_id string) (string, error) {
	// session_id must not contain double hyphens (--)
	if strings.Contains(session_id, "--") {
		session_id = ""
	}

	if len(session_id) == 0 {
		generated, err := GenerateSessionID()
		if err != nil {
			return "", err
		}
		return generated, nil
	}

	return session_id, nil
}
