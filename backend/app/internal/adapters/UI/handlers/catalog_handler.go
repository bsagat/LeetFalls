package handlers

import (
	"leetFalls/internal/service"
	"log/slog"
	"net/http"
)

type CatalogHandler struct {
	postServ service.PostService
}

func NewCatalogHandler() *CatalogHandler {
	return &CatalogHandler{}
}

func (h *CatalogHandler) ServeMainPage(w http.ResponseWriter, r *http.Request) {

}

func (h *CatalogHandler) ServeCatalogPost(w http.ResponseWriter, r *http.Request) {

}

func (h *CatalogHandler) ShowCreatePostForm(w http.ResponseWriter, r *http.Request) {

}

func (h *CatalogHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Cookie parsing
	cookie, err := r.Cookie("session_id")
	if err != nil {
		slog.Warn("Session_id is not exist")
	}

	// MultiPartForm parsing
	//
	// Max size calculation:
	// 10 << 20 = 10 * 2^20 = 10 * 1048576 = 10485760
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		slog.Error("Failed to parse form: ", "error", err.Error())
		return
	}

	// Business logic call
	code, err := h.postServ.CreatePost(cookie)
	if err != nil {
		slog.Error("Failed to create post: ", "error", err.Error())
		ErrorPage(w, err, int(code))
		return
	}

	// MainPage show
}

func (h *CatalogHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {

}
