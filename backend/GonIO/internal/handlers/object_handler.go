package handlers

import (
	"GonIO/internal/domain"
	"log"
	"net/http"
)

type ObjectHandler struct {
	serv domain.ObjectService
}

func NewObjectHandler(serv domain.ObjectService) *ObjectHandler {
	return &ObjectHandler{serv: serv}
}

func (h *ObjectHandler) GetObjectList(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	if bucketName == "" {
		log.Printf("Failed to get path value: %s", domain.ErrEmptyBucketName.Error())
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	h.serv.ObjectList(w, bucketName)
}

func (h *ObjectHandler) RetrieveObject(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	if bucketName == "" {
		log.Printf("Failed to get path value: %s", domain.ErrEmptyBucketName.Error())
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	objectName := r.PathValue("ObjectKey")
	if objectName == "" {
		log.Printf("Failed to get path value: %s", domain.ErrEmptyObjectName.Error())
		http.Error(w, domain.ErrEmptyObjectName.Error(), http.StatusBadRequest)
		return
	}

	h.serv.RetrieveObject(w, bucketName, objectName)
}

func (h *ObjectHandler) UpdateObject(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	if bucketName == "" {
		log.Printf("Failed to get path value: %s", domain.ErrEmptyBucketName.Error())
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	objectName := r.PathValue("ObjectKey")
	if objectName == "" {
		log.Printf("Failed to get path value: %s", domain.ErrEmptyObjectName.Error())
		http.Error(w, domain.ErrEmptyObjectName.Error(), http.StatusBadRequest)
		return
	}

	h.serv.UploadObject(w, r, bucketName, objectName)
}

func (h *ObjectHandler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	if bucketName == "" {
		log.Printf("Failed to get path value: %s", domain.ErrEmptyBucketName.Error())
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	objectName := r.PathValue("ObjectKey")
	if objectName == "" {
		log.Printf("Failed to get path value: %s", domain.ErrEmptyObjectName.Error())
		http.Error(w, domain.ErrEmptyObjectName.Error(), http.StatusBadRequest)
		return
	}

	h.serv.DeleteObject(w, r, bucketName, objectName)
}
