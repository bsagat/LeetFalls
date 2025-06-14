package handlers

import "net/http"

type ArchiveHandler struct {
}

func NewArchiveHandler() *ArchiveHandler {
	return &ArchiveHandler{}
}

func (h *ArchiveHandler) ServeArchivePage(w http.ResponseWriter, r *http.Request) {

}

func (h *ArchiveHandler) ServeArchivePost(w http.ResponseWriter, r *http.Request) {

}
