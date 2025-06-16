package domain

import (
	"net/http"
)

type BucketDal interface {
	GetBucketList() ([]Bucket, error)
	IsUniqueBucket(bucketName string) (bool, error)
	CreateBucket(bucketname string) error
	DeleteBucket(bucketName string) error
}

type ObjectDal interface {
	IsObjectExist(path, name string) (bool, error)
	List_Object(bucketname string) (ObjectsList, error)
	UploadObject(bucketname, objectname string, r *http.Request) error
	RetrieveObject(bucketname, objectname string, w http.ResponseWriter) error
	DeleteObject(bucketname, objectname string) error
}

type BucketService interface {
	BucketList(w http.ResponseWriter)
	CreateBucket(w http.ResponseWriter, bucketName string)
	DeleteBucket(w http.ResponseWriter, bucketName string)
}

type ObjectService interface {
	ObjectList(w http.ResponseWriter, bucketname string)
	RetrieveObject(w http.ResponseWriter, bucketname, objectname string)
	UploadObject(w http.ResponseWriter, r *http.Request, bucketname, objectname string)
	DeleteObject(w http.ResponseWriter, r *http.Request, bucketname, objectname string)
}
