package domain

import "io"

// === DAL Interfaces ===

type BucketDal interface {
	GetBucketList() ([]Bucket, error)
	IsUniqueBucket(bucketName string) (bool, error)
	CreateBucket(bucketName string) error
	DeleteBucket(bucketName string) error
}

type ObjectDal interface {
	IsObjectExist(bucketName, objectName string) (bool, error)
	ListObjects(bucketName string) ([]Object, error)
	UploadObject(bucketName, objectName string, data io.ReadCloser, contentLen int64, contentType string) error
	RetrieveObject(bucketName, objectName string) (io.ReadCloser, error)
	DeleteObject(bucketName, objectName string) error
}

// === Service Interfaces ===

type BucketService interface {
	BucketList() ([]Bucket, error)
	CreateBucket(bucketName string) (int, error)
	DeleteBucket(bucketName string) (int, error)
}

type ObjectService interface {
	ObjectList(bucketName string) ([]Object, int, error)
	RetrieveObject(bucketName, objectName string) (io.ReadCloser, int, error)
	UploadObject(bucketName, objectName, contentType string, data io.ReadCloser, contentLen int64) (int, error)
	UploadObjectJar(bucketName string, reqBody io.ReadCloser) (int, error)
	DeleteObject(bucketName, objectName string) (int, error)
}
