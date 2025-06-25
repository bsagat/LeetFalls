package handlers

import (
	"errors"
	"leetFalls/internal/domain"
	"leetFalls/internal/service"
	"log/slog"
	"net/http"
)

type CatalogHandler struct {
	postServ    service.PostService
	commentServ service.CommentService
	authServ    service.AuthService
}

func NewCatalogHandler(postServ service.PostService, authServ service.AuthService, commentServ service.CommentService) *CatalogHandler {
	return &CatalogHandler{postServ: postServ, authServ: authServ, commentServ: commentServ}
}

// Shows catalog page with posts
func (h *CatalogHandler) ServeMainPage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		slog.Warn("Session_id is not exist")
	}

	if _, err := h.authServ.Auth(w, cookie); err != nil {
		slog.Error("Failed to auth user with cookie: ", "error", err.Error())
		ErrorPage(w, errors.New("user auth error: "+err.Error()), http.StatusInternalServerError)
		return
	}

	code, err := h.postServ.ShowCatalogWithPosts(w)
	if err != nil {
		slog.Error("Failed to show catalog page: ", "error", err.Error())
		ErrorPage(w, errors.New("catalog showcase error: "+err.Error()), int(code))
		return
	}
}

// Serve post from catalog page
func (h *CatalogHandler) ServeCatalogPost(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("id")
	code, err := h.postServ.ShowPost(w, postId, false)
	if err != nil {
		slog.Error("Failed to show post page: ", "error", err.Error())
		ErrorPage(w, errors.New("post showcase error: "+err.Error()), int(code))
		return
	}
}

// Shows post create page
func (h *CatalogHandler) ShowCreatePostForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, domain.Config.TemplatesPath+"/create-post.html")
}

// Handles post create
func (h *CatalogHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Cookie parsing
	cookie, err := r.Cookie("session_id")
	if err != nil {
		slog.Warn("Session_id is not exist")
	}

	// MultiPartForm parsing
	// Max size calculation:
	// 10 << 20 = 10 * 2^20 = 10 * 1048576 = 10485760
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		slog.Error("Failed to parse form: ", "error", err.Error())
		ErrorPage(w, errors.New("form file parse error: "+err.Error()), http.StatusInternalServerError)
		return
	}

	// File Parse
	file, _, err := r.FormFile("File")
	if err != nil && err != http.ErrMissingFile {
		slog.Error("Failed to parse file from form: ", "error", err.Error())
		ErrorPage(w, errors.New("form file reading error: "+err.Error()), http.StatusInternalServerError)
		return
	}
	if file != nil {
		defer file.Close()
	}
	// Business logic call
	code, err := h.postServ.CreatePost(w, r.FormValue("Name"), r.FormValue("Title"), r.FormValue("Content"), file, cookie)
	if err != nil {
		slog.Error("Failed to create post: ", "error", err.Error())
		ErrorPage(w, err, int(code))
		return
	}

	http.Redirect(w, r, "http://localhost:8080/catalog", http.StatusSeeOther)
}

// Handles comment create
func (h *CatalogHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Authorization by cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		slog.Warn("Session_id is not exist")
	}

	authorId, err := h.authServ.Auth(w, cookie)
	if err != nil {
		slog.Error("Failed to auth user with cookie: ", "error", err.Error())
		ErrorPage(w, errors.New("user auth error: "+err.Error()), http.StatusInternalServerError)
		return
	}

	// Comment details parsing from Form
	if err = r.ParseMultipartForm(10 << 20); err != nil {
		slog.Error("Failed to parse form: ", "error", err.Error())
		ErrorPage(w, errors.New("form file parse error: "+err.Error()), http.StatusInternalServerError)
		return
	}

	// File Parse
	file, _, err := r.FormFile("File")
	if err != nil && err != http.ErrMissingFile {
		slog.Error("Failed to parse file from form: ", "error", err.Error())
		ErrorPage(w, errors.New("form file reading error: "+err.Error()), http.StatusInternalServerError)
		return
	}

	if file != nil {
		defer file.Close()
	}

	code, err := h.commentServ.CreateComment(authorId, r.FormValue("postID"), r.FormValue("ReplyTo"), r.FormValue("Content"), file)
	if err != nil {
		slog.Error("Failed to create comment: ", "error", err.Error())
		ErrorPage(w, err, int(code))
		return
	}

	http.Redirect(w, r, "http://localhost:8080/posts/"+r.FormValue("postID"), http.StatusSeeOther)
}
