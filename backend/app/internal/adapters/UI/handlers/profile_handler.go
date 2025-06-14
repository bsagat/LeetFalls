package handlers

import "net/http"

type ProfileHandler struct{}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

func (h *ProfileHandler) ServeProfilePage(w http.ResponseWriter, r *http.Request) {
}

func (h *ProfileHandler) GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {

}
