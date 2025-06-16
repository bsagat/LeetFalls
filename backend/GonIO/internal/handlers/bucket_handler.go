package handlers

import (
	"GonIO/internal/domain"
	"log"
	"net/http"
)

type BucketHandler struct {
	serv domain.BucketService
}

func NewBucketHandler(serv domain.BucketService) *BucketHandler {
	return &BucketHandler{serv: serv}
}

func (h *BucketHandler) BucketListsHandler(w http.ResponseWriter, r *http.Request) {
	h.serv.BucketList(w)
}

func (h *BucketHandler) CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	if bucketName == "" {
		log.Printf("Failed to get path value: %s", domain.ErrEmptyBucketName.Error())
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	h.serv.CreateBucket(w, bucketName)
}

func (h *BucketHandler) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	if bucketName == "" {
		log.Printf("Failed to get path value: %s", domain.ErrEmptyBucketName.Error())
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	h.serv.DeleteBucket(w, bucketName)
}
