package service

import (
	"GonIO/internal/domain"
	"archive/zip"
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
)

type ObjectServiceImp struct {
	dal        domain.ObjectDal
	bucketRepo domain.BucketDal
}

var _ domain.ObjectService = (*ObjectServiceImp)(nil)

func NewObjectService(dal domain.ObjectDal, bucketRepo domain.BucketDal) *ObjectServiceImp {
	return &ObjectServiceImp{dal: dal, bucketRepo: bucketRepo}
}

func (s ObjectServiceImp) ObjectList(bucketname string) ([]domain.Object, int, error) {
	if notExist, err := s.bucketNotExist(bucketname); notExist || err != nil {
		return nil, getStatusFromErr(err), err
	}

	objects, err := s.dal.ListObjects(bucketname)
	if err != nil {
		slog.Error("List_Object failed", "bucket", bucketname, "err", err)
		return nil, http.StatusInternalServerError, err
	}
	return objects, http.StatusOK, nil
}

func (s ObjectServiceImp) RetrieveObject(bucketname, objectname string) (io.ReadCloser, int, error) {
	if notExist, err := s.bucketNotExist(bucketname); notExist || err != nil {
		return nil, getStatusFromErr(err), err
	}

	if err := Validate(objectname); err != nil {
		slog.Error("Invalid object name", "object", objectname, "err", err)
		return nil, http.StatusBadRequest, err
	}

	metaPath := filepath.Join(domain.BucketsPath, bucketname, "objects.csv")
	exist, err := s.dal.IsObjectExist(metaPath, objectname)
	if err != nil {
		slog.Error("Failed to check object existence", "object", objectname, "err", err)
		return nil, http.StatusInternalServerError, err
	}
	if !exist {
		slog.Info("Object does not exist", "object", objectname)
		return nil, http.StatusNotFound, domain.ErrObjectIsNotExist
	}

	reader, err := s.dal.RetrieveObject(bucketname, objectname)
	if err != nil {
		slog.Error("RetrieveObject failed", "object", objectname, "err", err)
		return nil, http.StatusInternalServerError, err
	}
	return reader, http.StatusOK, nil
}

func (s ObjectServiceImp) UploadObject(bucketname, objectname, fileType string, image io.ReadCloser, imageLen int64) (int, error) {
	if notExist, err := s.bucketNotExist(bucketname); notExist || err != nil {
		return getStatusFromErr(err), err
	}

	if err := Validate(objectname); err != nil {
		slog.Error("Invalid object name", "object", objectname, "err", err)
		return http.StatusBadRequest, err
	}

	if err := s.dal.UploadObject(bucketname, objectname, image, imageLen, fileType); err != nil {
		slog.Error("UploadObject failed", "object", objectname, "err", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

func (s ObjectServiceImp) DeleteObject(bucketname, objectname string) (int, error) {
	if notExist, err := s.bucketNotExist(bucketname); notExist || err != nil {
		return getStatusFromErr(err), err
	}

	if err := Validate(objectname); err != nil {
		slog.Error("Invalid object name", "object", objectname, "err", err)
		return http.StatusBadRequest, err
	}

	metaPath := filepath.Join(domain.BucketsPath, bucketname, "objects.csv")
	exist, err := s.dal.IsObjectExist(metaPath, objectname)
	if err != nil {
		slog.Error("Failed to check object existence", "object", objectname, "err", err)
		return http.StatusInternalServerError, err
	}
	if !exist {
		slog.Info("Object does not exist", "object", objectname)
		return http.StatusNotFound, domain.ErrObjectIsNotExist
	}

	if err := s.dal.DeleteObject(bucketname, objectname); err != nil {
		slog.Error("DeleteObject failed", "object", objectname, "err", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (s ObjectServiceImp) UploadObjectJar(bucketName string, reqBody io.ReadCloser) (int, error) {
	if notExist, err := s.bucketNotExist(bucketName); notExist || err != nil {
		slog.Error("Bucket does not exist or error occurred", "bucket", bucketName, "error", err)
		return getStatusFromErr(err), err
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, reqBody); err != nil {
		slog.Error("Failed to read request body", "error", err)
		return http.StatusInternalServerError, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		slog.Error("Failed to read ZIP archive", "error", err)
		return http.StatusBadRequest, err
	}

	for _, zipFile := range zipReader.File {
		// Pass the directory
		if zipFile.FileInfo().IsDir() {
			continue
		}

		fileReader, err := zipFile.Open()
		if err != nil {
			slog.Error("Failed to open file inside ZIP", "file", zipFile.Name, "error", err)
			return http.StatusInternalServerError, err
		}
		defer fileReader.Close()

		var fileBuf bytes.Buffer
		if _, err := io.Copy(&fileBuf, fileReader); err != nil {
			slog.Error("Failed to read file from ZIP", "file", zipFile.Name, "error", err)
			return http.StatusInternalServerError, err
		}

		objectName := zipFile.FileInfo().Name()

		if err := Validate(objectName); err != nil {
			slog.Error("Invalid object name", "object", objectName, "error", err)
			return http.StatusBadRequest, err
		}

		contentType := zipFile.Mode().Type().String()
		contentLen := zipFile.FileInfo().Size()

		err = s.dal.UploadObject(
			bucketName,
			objectName,
			io.NopCloser(bytes.NewReader(fileBuf.Bytes())),
			contentLen,
			contentType,
		)
		if err != nil {
			slog.Error("Failed to upload object from ZIP", "object", objectName, "error", err)
			return http.StatusInternalServerError, err
		}

		slog.Info("Uploaded object from ZIP successfully", "bucket", bucketName, "object", objectName)
	}

	return http.StatusCreated, nil
}

func (s ObjectServiceImp) bucketNotExist(bucketname string) (bool, error) {
	unique, err := s.bucketRepo.IsUniqueBucket(bucketname)
	if err != nil {
		slog.Error("Failed to check if bucket exists", "bucket", bucketname, "err", err)
		return false, err
	}
	if unique {
		slog.Info("Bucket not found", "bucket", bucketname)
		return true, domain.ErrBucketIsNotExist
	}
	return false, nil
}

func getStatusFromErr(err error) int {
	switch err {
	case domain.ErrBucketIsNotExist:
		return http.StatusNotFound
	case domain.ErrObjectIsNotExist:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
