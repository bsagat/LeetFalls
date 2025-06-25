package handlers

import (
	"GonIO/internal/domain"
	xmlsender "GonIO/pkg/xmlMsgSender"
	"fmt"
	"io"
	"log/slog"
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
		slog.Error("Missing bucket name in path")
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	objectList, code, err := h.serv.ObjectList(bucketName)
	if err != nil {
		slog.Error("Failed to retrieve object list", "bucket", bucketName, "error", err)
		http.Error(w, err.Error(), code)
		return
	}

	if err := xmlsender.SendObjectList(w, objectList); err != nil {
		slog.Error("Failed to send object list XML response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Object list retrieved successfully", "bucket", bucketName)
}

func (h *ObjectHandler) RetrieveObject(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	if bucketName == "" {
		slog.Error("Missing bucket name in path")
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	objectName := r.PathValue("ObjectKey")
	if objectName == "" {
		slog.Error("Missing object name in path")
		http.Error(w, domain.ErrEmptyObjectName.Error(), http.StatusBadRequest)
		return
	}

	objectReader, code, err := h.serv.RetrieveObject(bucketName, objectName)
	if err != nil {
		slog.Error("Failed to retrieve object", "bucket", bucketName, "object", objectName, "error", err)
		http.Error(w, err.Error(), code)
		return
	}
	defer objectReader.Close()

	if _, err := io.Copy(w, objectReader); err != nil {
		slog.Error("Failed to write object to response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Object retrieved successfully", "bucket", bucketName, "object", objectName)
}

func (h *ObjectHandler) UpdateObject(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	if bucketName == "" {
		slog.Error("Missing bucket name in path")
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	objectName := r.PathValue("ObjectKey")
	if objectName == "" {
		slog.Error("Missing object name in path")
		http.Error(w, domain.ErrEmptyObjectName.Error(), http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		slog.Error("Request body is empty")
		http.Error(w, domain.ErrEmptyReqBody.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	contentType := r.Header.Get("Content-Type")
	code, err := h.serv.UploadObject(bucketName, objectName, contentType, r.Body, r.ContentLength)
	if err != nil {
		slog.Error("Failed to upload object", "bucket", bucketName, "object", objectName, "error", err)
		http.Error(w, err.Error(), code)
		return
	}

	if err := xmlsender.SendMessage(w, code, fmt.Sprintf("Object '%s' uploaded successfully", objectName)); err != nil {
		slog.Error("Failed to send XML response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Object uploaded successfully", "bucket", bucketName, "object", objectName)
}

func (h *ObjectHandler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	if bucketName == "" {
		slog.Error("Missing bucket name in path")
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	objectName := r.PathValue("ObjectKey")
	if objectName == "" {
		slog.Error("Missing object name in path")
		http.Error(w, domain.ErrEmptyObjectName.Error(), http.StatusBadRequest)
		return
	}

	code, err := h.serv.DeleteObject(bucketName, objectName)
	if err != nil {
		slog.Error("Failed to delete object", "bucket", bucketName, "object", objectName, "error", err)
		http.Error(w, err.Error(), code)
		return
	}

	if err := xmlsender.SendMessage(w, code, fmt.Sprintf("Object '%s' deleted successfully", objectName)); err != nil {
		slog.Error("Failed to send XML response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Object deleted successfully", "bucket", bucketName, "object", objectName)
}

func (h *ObjectHandler) ObjectJarHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.PathValue("BucketName")
	if bucketName == "" {
		slog.Error("Missing bucket name in path")
		http.Error(w, domain.ErrEmptyBucketName.Error(), http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		slog.Error("Request body is empty")
		http.Error(w, domain.ErrEmptyReqBody.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	code, err := h.serv.UploadObjectJar(bucketName, r.Body)
	if err != nil {
		slog.Error("Failed to upload object jar: ", "error", err)
		http.Error(w, err.Error(), code)
		return
	}

	slog.Info("Object jar uploaded successfully", "bucket", bucketName)
}
