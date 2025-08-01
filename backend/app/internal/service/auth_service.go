package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"leetFalls/internal/domain/models"
	"leetFalls/internal/domain/ports"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type AuthService struct {
	dbrepo   ports.UserRepo
	external ports.ExternalAPI
	stop     chan struct{}
	done     chan struct{}
}

func NewAuthService(dbrepo ports.UserRepo, external ports.ExternalAPI, deleteTimeout int) *AuthService {
	serv := &AuthService{
		dbrepo:   dbrepo,
		external: external,
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}
	serv.StartDeleteExpiredSessions(deleteTimeout)
	return serv
}

func (s *AuthService) StartDeleteExpiredSessions(timeoutMinutes int) {
	t := time.NewTicker(time.Duration(timeoutMinutes) * time.Minute)
	go func() {
		defer func() {
			t.Stop()
			close(s.done)
		}()
		for {
			select {
			case <-s.stop:
				return
			case <-t.C:
				if err := s.dbrepo.DeleteExpiredSessions(); err != nil {
					slog.Warn("Failed to delete expired sessions", "error", err)
				}
			}
		}
	}()
}

func (s *AuthService) Close() {
	close(s.stop)
	<-s.done
}

// Validates the session cookie. If valid, it retrieves and returns the current user's ID.
// If the cookie is missing, invalid, or expired,
// a new session is created (for an unauthenticated/guest user) and its associated user ID is returned.
func (s *AuthService) Auth(w http.ResponseWriter, cookie *http.Cookie) (int, error) {
	if cookie != nil && CheckSessionId(cookie.Value) {
		session_id := cookie.Value
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

	// if session id not exist or invalid, we generate new user
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
func (s *AuthService) CreateNewUser(sessionId string) (int, error) {
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

	slog.Info(fmt.Sprintf("User with id %d created succesfully", user.ID))
	return user.ID, nil
}

func (s *AuthService) ChangeUserName(userId int, userName string) error {
	// User name modification
	if err := s.dbrepo.ChangeUserName(userId, userName); err != nil {
		slog.Error("Failed to change user name: ", "error", err.Error())
		return err
	}
	return nil
}

func (s *AuthService) GetUserById(id int) (models.User, error) {
	user, err := s.dbrepo.GetUserById(id)
	if err != nil {
		slog.Error("Failed to get user by id: ", "error", err.Error())
		return user, err
	}
	return user, nil
}

// Generates session_id randomly
func GenerateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Validates session id
func CheckSessionId(session_id string) bool {
	// session_id must not contain double hyphens (--) or empty
	if strings.Contains(session_id, "--") || len(session_id) == 0 {
		return false
	}

	return true
}
