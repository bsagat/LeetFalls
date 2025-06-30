package handlers

import "net/http"

type ProfileHandler struct{}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

// YAGNI, ya ya i know

func (h *ProfileHandler) ServeProfilePage(w http.ResponseWriter, r *http.Request) {
}

func (h *ProfileHandler) GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {

}
