package service

import (
	"fmt"
	"html/template"
	"leetFalls/internal/domain"
	"leetFalls/internal/domain/models"
	"leetFalls/internal/domain/ports"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type PostService struct {
	AuthService ports.AuthService
	storage     ports.Storage
	commentRepo ports.CommentRepo
	postRepo    ports.PostsRepo
}

func NewPostService(authService ports.AuthService, storage ports.Storage, postRepo ports.PostsRepo, commentRepo ports.CommentRepo) *PostService {
	return &PostService{
		AuthService: authService,
		storage:     storage,
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

func (s *PostService) CreatePost(w http.ResponseWriter, userName, title, content string, file multipart.File, cookie *http.Cookie) (int, error) {
	post := models.Post{
		Title:   title,
		Content: content,
	}

	// 1) Form fields validation
	if err := PostValidation(post); err != nil {
		return http.StatusBadRequest, err
	}

	// 2) Cookie session_id validation
	userId, err := s.AuthService.Auth(w, cookie)
	if err != nil {
		slog.Error("User authorization check failed: ", "error", err)
		return http.StatusInternalServerError, err
	}

	// 3) User name modification
	if err := s.AuthService.ChangeUserName(userId, userName); err != nil {
		slog.Error("Failed to change user name: ", "error", err)
		return http.StatusInternalServerError, err
	}
	post.Author.ID = userId

	// 4) Get unique post ID
	post.ID, err = s.postRepo.GetNextPostId()
	if err != nil {
		slog.Error("Failed to get next post id: ", "error", err)
		return http.StatusInternalServerError, err
	}

	// 5) Save image to storage
	if file != nil {
		if err := s.storage.SavePostImage(&post, file); err != nil {
			slog.Error("Failed to save post image to storage: ", "error", err)
			return http.StatusInternalServerError, err
		}
	}

	// 6) Save Post to Database
	err = s.postRepo.SavePost(post)
	if err != nil {
		slog.Error("Failed to save post to database: ", "error", err)
		return http.StatusInternalServerError, err
	}

	slog.Info(fmt.Sprintf("Post with id %d created succesfuly", post.ID))
	return http.StatusCreated, nil
}

func (s *PostService) ShowPost(w http.ResponseWriter, postId string, archive bool) (int, error) {
	// 1) Post Validation - post ID
	id, err := strconv.Atoi(postId)
	if err != nil {
		return http.StatusBadRequest, domain.ErrInvalidPostId
	}

	// 2) Parse post data from database
	post, err := s.postRepo.GetPost(id)
	if err != nil {
		slog.Error("Failed to get post by ID", "postID", id, "error", err)
		return http.StatusInternalServerError, err
	}
	if post.ID == 0 {
		return http.StatusNotFound, domain.ErrPostNotFound
	}

	// 3) Validate post TTL during an archive/unarchive request:
	now := time.Now()
	if archive && post.ExpiresAt.After(now) {
		return http.StatusBadRequest, domain.ErrPostIsActive
	}
	if !archive && post.ExpiresAt.Before(now) {
		return http.StatusBadRequest, domain.ErrPostIsArchived
	}

	// 4) Parse author information
	author, err := s.AuthService.GetUserById(post.Author.ID)
	if err != nil {
		slog.Error("Failed to get user by ID: ", "error", err)
		return http.StatusInternalServerError, err
	}
	if author.ID == 0 {
		slog.Error("Failed to get user data by id: ", "error", domain.ErrUserNotExist)
		return http.StatusNotFound, fmt.Errorf("(post author) %w", domain.ErrUserNotExist)
	}

	// 5) Get Post Comments list
	comments, err := s.commentRepo.CommentsByPost(id)
	if err != nil {
		slog.Error("Failed to get comments by post: ", "error", err)
		return http.StatusInternalServerError, err
	}

	// 6) Show HTML template
	htmlFile := "/post.html"
	if archive {
		htmlFile = "/archive-post.html"
	}
	temp, err := template.ParseFiles(domain.Config.TemplatesPath + htmlFile)
	if err != nil {
		slog.Error("Failed to serve post page: ", "error", err)
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
		slog.Error("Failed to execute post page: ", "error", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *PostService) ShowCatalogWithPosts(w http.ResponseWriter) (int, error) {
	posts, err := s.postRepo.ActivePosts()
	if err != nil {
		slog.Error("Failed to get active posts from database: ", "error", err)
		return http.StatusInternalServerError, err
	}

	temp, err := template.ParseFiles(domain.Config.TemplatesPath + "/main_page.html")
	if err != nil {
		slog.Error("Failed to serve main page: ", "error", err)
		return http.StatusInternalServerError, err
	}

	if err := temp.Execute(w, posts); err != nil {
		slog.Error("Failed to execute posts on main page: ", "error", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *PostService) ShowArchiveWithPosts(w http.ResponseWriter) (int, error) {
	posts, err := s.postRepo.ArchivePosts()
	if err != nil {
		slog.Error("Failed to get archive posts from database: ", "error", err)
		return http.StatusInternalServerError, err
	}

	temp, err := template.ParseFiles(domain.Config.TemplatesPath + "/archive.html")
	if err != nil {
		slog.Error("Failed to serve archive page: ", "error", err)
		return http.StatusInternalServerError, err
	}

	if err := temp.Execute(w, posts); err != nil {
		slog.Error("Failed to execute posts on archive page: ", "error", err)
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
