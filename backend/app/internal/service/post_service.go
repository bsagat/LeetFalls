package service

import (
	"fmt"
	dbrepo "leetFalls/internal/adapters/dbRepo"
	"leetFalls/internal/adapters/storage"
	"leetFalls/internal/domain"
	"leetFalls/internal/domain/models"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"
)

type PostService struct {
	sessionServ MiddlewareService
	storage     storage.GonIO
	repo        dbrepo.PostsRepo
}

func NewPostService(sessionServ MiddlewareService, storage storage.GonIO) *PostService {
	return &PostService{sessionServ: sessionServ, storage: storage}
}

func (s *PostService) CreatePost(w http.ResponseWriter, userName, title, content string, file multipart.File, cookie *http.Cookie) (domain.Code, error) {
	post := models.Post{
		Title:   title,
		Content: content,
	}

	if err := PostValidation(post); err != nil {
		return http.StatusBadRequest, err
	}

	// Sql injection check
	if strings.Contains(userName, "--") {
		return http.StatusBadRequest, domain.ErrUserNameDoubleGyphen
	}

	// Cookie session_id validation
	userId, err := s.sessionServ.Auth(w, cookie)
	if err != nil {
		slog.Error("User authorization check failed: ", "error", err)
		return http.StatusInternalServerError, err
	}

	// User name modification
	err = s.sessionServ.dbrepo.ChangeUserName(userId, userName)
	if err != nil {
		slog.Error("Failed to change user name: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}
	post.AuthorID = userId

	// Post save
	post.ID, err = s.repo.SavePost(post)
	if err != nil {
		slog.Error("Failed to save post to database: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	// Post Image save
	if file != nil {
		if err := s.storage.SavePostImage(post.ID, file); err != nil {
			slog.Error("Failed to save post image to storage: ", "error", err)
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusCreated, nil
}

func PostValidation(post models.Post) error {
	if len(post.Title) == 0 {
		return fmt.Errorf("post title is empty")
	}
	if len(post.Content) == 0 {
		return fmt.Errorf("post content is empty")
	}
	return nil
}
