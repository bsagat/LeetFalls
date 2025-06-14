package service

import (
	"fmt"
	"leetFalls/internal/domain"
	"leetFalls/internal/domain/models"
	"log/slog"
	"mime/multipart"
	"net/http"
)

type PostService struct {
	sessionServ MiddlewareService
}

func NewPostService() *PostService {
	return &PostService{}
}

func (s *PostService) CreatePost(w http.ResponseWriter, name, title, content string, file multipart.File, cookie *http.Cookie) (domain.Code, error) {
	post := models.Post{
		Changed_name: name,
		Title:        title,
	}
	PostValidation()
	// нужно написать код который будет сравнивать совпадают ли changed name и name in database
	// если да тогда меняем его на changed
	// его нужно убрать из модельки

	userId, err := s.sessionServ.Auth(w, cookie)
	if err != nil {
		slog.Error("User authorization check failed: ", "error", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func PostValidation(post models.Post) error {
	if post.ID != 0 {
		return fmt.Errorf("post ID must be empty")
	}
	if post.Title == "" {
		return fmt.Errorf("post title is empty")
	}
	if post.Content == "" {
		return fmt.Errorf("post content is empty")
	}
	if post.AuthorID != 0 {
		return fmt.Errorf("author_id field must be empty")
	}
	if !post.CreatedAt.IsZero() {
		return fmt.Errorf("createdAt field must be empty")
	}
	if !post.ExpiresAt.IsZero() {
		return fmt.Errorf("expires_at field must be empty")
	}
	if post.ImageLink != "" {
		return fmt.Errorf("image URL field must be empty")
	}
	return nil
}
