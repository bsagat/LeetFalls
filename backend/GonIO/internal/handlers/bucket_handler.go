package handlers

import (
	"GonIO/internal/domain"
	xmlsender "GonIO/pkg/xmlMsgSender"
	"fmt"
	"log/slog"
	"net/http"
)

type BucketHandler struct {
	serv domain.BucketService
}

func NewBucketHandler(serv domain.BucketService) *BucketHandler {
	return &BucketHandler{serv: serv}
}

func (h *BucketHandler) BucketListsHandler(w http.ResponseWriter, r *http.Request) {
	list, err := h.serv.BucketList()
	if err != nil {
		slog.Error("Failed to get bucket list: ", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = xmlsender.SendBucketList(w, list); err != nil {
		slog.Error("Failed to send bucket list: ", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *BucketHandler) CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	if bucketName == "" {
		slog.Error("Failed to get path value: ", "error", domain.ErrEmptyBucketName)
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	code, err := h.serv.CreateBucket(bucketName)
	if err != nil {
		slog.Error("Failed to create bucket: ", "error", err)
		http.Error(w, err.Error(), code)
		return
	}

	if err = xmlsender.SendMessage(w, code, fmt.Sprintf("bucket with name %s created succesfully", bucketName)); err != nil {
		slog.Error("Failed to send xml message: ", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *BucketHandler) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	if bucketName == "" {
		slog.Error("Failed to get path value: ", "error", domain.ErrEmptyBucketName)
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	code, err := h.serv.DeleteBucket(bucketName)
	if err != nil {
		slog.Error("Failed to delete bucket: ", "error", err)
		http.Error(w, err.Error(), code)
		return
	}

	if err = xmlsender.SendMessage(w, code, fmt.Sprintf("bucket with name %s deleted succesfully", bucketName)); err != nil {
		slog.Error("Failed to send xml message: ", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
