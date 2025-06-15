package service

import (
	"errors"
	"fmt"
	"html/template"
	dbrepo "leetFalls/internal/adapters/dbRepo"
	"leetFalls/internal/adapters/storage"
	"leetFalls/internal/domain"
	"leetFalls/internal/domain/models"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostService struct {
	AuthService AuthService
	storage     storage.GonIO
	commentRepo dbrepo.CommentRepo
	repo        dbrepo.PostsRepo
}

func NewPostService(authService AuthService, storage storage.GonIO, repo dbrepo.PostsRepo, commentRepo dbrepo.CommentRepo) *PostService {
	return &PostService{
		AuthService: authService,
		storage:     storage,
		repo:        repo,
		commentRepo: commentRepo,
	}
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
	userId, err := s.AuthService.Auth(w, cookie)
	if err != nil {
		slog.Error("User authorization check failed: ", "error", err)
		return http.StatusInternalServerError, err
	}

	// User name modification
	if err := s.AuthService.ChangeUserName(userId, userName); err != nil {
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

	slog.Info(fmt.Sprintf("Post with id %d created succesfuly", post.ID))
	return http.StatusCreated, nil
}

func (s *PostService) ShowPost(w http.ResponseWriter, postId string, archive bool) (domain.Code, error) {
	id, err := strconv.Atoi(postId)
	if err != nil {
		return http.StatusBadRequest, domain.ErrInvalidPostId
	}

	post, err := s.repo.GetPost(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if post.ID == 0 {
		return http.StatusNotFound, domain.ErrPostNotFound
	}

	now := time.Now()
	if archive && post.ExpiresAt.After(now) {
		return http.StatusBadRequest, errors.New("post is still active")
	}
	if !archive && post.ExpiresAt.Before(now) {
		return http.StatusBadRequest, errors.New("post is archived")
	}

	// Author information fetching
	author, err := s.AuthService.GetUserById(post.AuthorID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Post comments fetching
	comments, err := s.commentRepo.GetCommentsByPost(id)
	if err != nil {
		slog.Error("Failed to get comments by post: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	temp, err := template.ParseFiles(domain.Config.TemplatesPath + "/post.html")
	if err != nil {
		slog.Error("Failed to serve post page: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	data := struct {
		AuthorImageURL string
		AuthorName     string
		CreatedAt      string
		ID             int
		ImageLink      string
		Title          string
		Content        string
		Comments       []models.Comment
	}{
		AuthorImageURL: author.ImageURL,
		AuthorName:     author.Name,
		CreatedAt:      post.CreatedAt.Format("02 January 2006, 15:04:05 UTC"),
		ID:             post.ID,
		ImageLink:      post.ImageLink,
		Title:          post.Title,
		Content:        post.Content,
		Comments:       comments,
	}

	if err := temp.Execute(w, data); err != nil {
		slog.Error("Failed to execute post page: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *PostService) ShowCatalogWithPosts(w http.ResponseWriter) (domain.Code, error) {
	posts, err := s.repo.ActivePosts()
	if err != nil {
		slog.Error("Failed to get active posts from database: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	temp, err := template.ParseFiles(domain.Config.TemplatesPath + "/main_page.html")
	if err != nil {
		slog.Error("Failed to serve main page: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	if err := temp.Execute(w, posts); err != nil {
		slog.Error("Failed to execute posts on main page: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *PostService) ShowArchiveWithPosts(w http.ResponseWriter) (domain.Code, error) {
	posts, err := s.repo.ArchivePosts()
	if err != nil {
		slog.Error("Failed to get archive posts from database: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	temp, err := template.ParseFiles(domain.Config.TemplatesPath + "/archive.html")
	if err != nil {
		slog.Error("Failed to serve archive page: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	if err := temp.Execute(w, posts); err != nil {
		slog.Error("Failed to execute posts on archive page: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
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
