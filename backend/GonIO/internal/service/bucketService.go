package service

import (
	"GonIO/internal/domain"
	"log/slog"
	"net/http"
)

type BucketServiceImp struct {
	dal domain.BucketDal
}

var _ domain.BucketService = (*BucketServiceImp)(nil)

func NewBucketService(dal domain.BucketDal) *BucketServiceImp {
	return &BucketServiceImp{dal: dal}
}

func (s BucketServiceImp) BucketList() ([]domain.Bucket, error) {
	buckets, err := s.dal.GetBucketList()
	if err != nil {
		slog.Error("Could not retrieve bucket list", "error", err)
		return nil, err
	}
	return buckets, nil
}

func (s BucketServiceImp) CreateBucket(bucketName string) (int, error) {
	isUnique, err := s.dal.IsUniqueBucket(bucketName)
	if err != nil {
		slog.Error("Failed to check bucket name uniqueness", "error", err)
		return http.StatusInternalServerError, err
	}
	if !isUnique {
		slog.Info("Bucket name already exists", "bucket", bucketName)
		return http.StatusConflict, domain.ErrNotUniqueName
	}

	if err := Validate(bucketName); err != nil {
		slog.Error("Invalid bucket name", "bucket", bucketName, "error", err)
		return http.StatusBadRequest, err
	}

	if err := s.dal.CreateBucket(bucketName); err != nil {
		slog.Error("Failed to create bucket", "bucket", bucketName, "error", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (s BucketServiceImp) DeleteBucket(bucketName string) (int, error) {
	isUnique, err := s.dal.IsUniqueBucket(bucketName)
	if err != nil {
		slog.Error("Failed to check bucket existence", "bucket", bucketName, "error", err)
		return http.StatusInternalServerError, err
	}
	if isUnique {
		slog.Info("Bucket does not exist", "bucket", bucketName)
		return http.StatusNotFound, domain.ErrBucketIsNotExist
	}

	if err := s.dal.DeleteBucket(bucketName); err != nil {
		slog.Error("Failed to delete bucket", "bucket", bucketName, "error", err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
