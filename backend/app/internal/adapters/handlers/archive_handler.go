package handlers

import (
	"errors"
	"leetFalls/internal/service"
	"log/slog"
	"net/http"
)

type ArchiveHandler struct {
	postServ service.PostService
	authServ service.AuthService
}

func NewArchiveHandler(postServ service.PostService, authServ service.AuthService) *ArchiveHandler {
	return &ArchiveHandler{postServ: postServ, authServ: authServ}
}

func (h *ArchiveHandler) ServeArchivePage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		slog.Warn("Session_id is not exist")
	}

	if _, err := h.authServ.Auth(w, cookie); err != nil {
		slog.Error("Failed to auth user with cookie: ", "error", err.Error())
		ErrorPage(w, errors.New("user auth error: "+err.Error()), http.StatusInternalServerError)
		return
	}

	code, err := h.postServ.ShowArchiveWithPosts(w)
	if err != nil {
		slog.Error("Failed to show archive page: ", "error", err.Error())
		ErrorPage(w, errors.New("archive showcase error: "+err.Error()), int(code))
		return
	}
}

func (h *ArchiveHandler) ServeArchivePost(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("id")
	isArchive := true
	code, err := h.postServ.ShowPost(w, postId, isArchive)
	if err != nil {
		slog.Error("Failed to show post page: ", "error", err.Error())
		ErrorPage(w, errors.New("post showcase error: "+err.Error()), int(code))
		return
	}
}
