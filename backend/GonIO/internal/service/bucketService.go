package service

import (
	"GonIO/internal/domain"
	xmlsender "GonIO/pkg/xmlMsgSender"
	"fmt"
	"log"
	"net/http"
)

type BucketServiceImp struct {
	dal domain.BucketDal
}

var _ domain.BucketService = (*BucketServiceImp)(nil)

func NewBucketService(dal domain.BucketDal) *BucketServiceImp {
	return &BucketServiceImp{dal: dal}
}

// Bucket List retrieve logic
func (serv BucketServiceImp) BucketList(w http.ResponseWriter) {
	list, err := serv.dal.GetBucketList()
	if err != nil {
		log.Printf("Failed to get bucket list: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = xmlsender.SendBucketList(w, list); err != nil {
		log.Printf("Failed to send bucket list: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Bucket create logic (validation)
func (serv BucketServiceImp) CreateBucket(w http.ResponseWriter, bucketName string) {
	unique, err := serv.dal.IsUniqueBucket(bucketName)
	if err != nil {
		log.Printf("Failed to check if bucket name is unique: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !unique {
		log.Printf("Bucket name is not unique")
		http.Error(w, domain.ErrNotUniqueName.Error(), http.StatusConflict)
		return
	}

	if err = Validate(bucketName); err != nil {
		log.Printf("Bucket name validation error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = serv.dal.CreateBucket(bucketName); err != nil {
		log.Printf("Failed to create bucket: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = xmlsender.SendMessage(w, http.StatusCreated, fmt.Sprintf("bucket with name %s created succesfully", bucketName)); err != nil {
		log.Printf("Failed to send xml message: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Bucket Delete logic
func (serv BucketServiceImp) DeleteBucket(w http.ResponseWriter, bucketName string) {
	unique, err := serv.dal.IsUniqueBucket(bucketName)
	if err != nil {
		log.Printf("Failed to check if bucket name is unique: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if unique {
		log.Printf("Bucket is not exist")
		http.Error(w, domain.ErrBucketIsNotExist.Error(), http.StatusNotFound)
		return
	}

	if err = serv.dal.DeleteBucket(bucketName); err != nil {
		log.Printf("Failed to delete bucket: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = xmlsender.SendMessage(w, http.StatusOK, fmt.Sprintf("bucket with name %s deleted succesfully", bucketName)); err != nil {
		log.Printf("Failed to send xml message: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
