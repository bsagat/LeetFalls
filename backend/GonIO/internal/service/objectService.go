package service

import (
	"GonIO/internal/dal"
	"GonIO/internal/domain"
	xmlsender "GonIO/pkg/xmlMsgSender"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type ObjectServiceImp struct {
	dal domain.ObjectDal
}

var _ domain.ObjectService = (*ObjectServiceImp)(nil)

func NewObjectService(dal domain.ObjectDal) *ObjectServiceImp {
	return &ObjectServiceImp{dal: dal}
}

// Object list retrieve logics
func (serv ObjectServiceImp) ObjectList(w http.ResponseWriter, bucketname string) {
	unique, err := dal.NewBucketXMLRepo().IsUniqueBucket(bucketname)
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

	objectList, err := serv.dal.List_Object(bucketname)
	if err != nil {
		log.Printf("Failed to retrieve object list: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	if err := xml.NewEncoder(w).Encode(objectList); err != nil {
		log.Printf("Failed to encode response: %s", err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Object retrieve logic
func (serv ObjectServiceImp) RetrieveObject(w http.ResponseWriter, bucketname, objectname string) {
	unique, err := dal.NewBucketXMLRepo().IsUniqueBucket(bucketname)
	if err != nil {
		log.Printf("Failed to check if bucket name exists: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if unique {
		log.Printf("Bucket is not exist")
		http.Error(w, domain.ErrBucketIsNotExist.Error(), http.StatusNotFound)
		return
	}

	if err = Validate(objectname); err != nil {
		log.Printf("Object name validation error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exist, err := serv.dal.IsObjectExist(domain.BucketsPath+"/"+bucketname+"/objects.csv", objectname)
	if err != nil {
		log.Printf("Failed to check if object name exists: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exist {
		log.Printf("Object is not exist")
		http.Error(w, domain.ErrObjectIsNotExist.Error(), http.StatusNotFound)
		return
	}

	if err = serv.dal.RetrieveObject(bucketname, objectname, w); err != nil {
		log.Printf("Failed to retrieve object: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Object upload logic
func (serv ObjectServiceImp) UploadObject(w http.ResponseWriter, r *http.Request, bucketname, objectname string) {
	unique, err := dal.NewBucketXMLRepo().IsUniqueBucket(bucketname)
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

	if err = Validate(objectname); err != nil {
		log.Printf("Object name validation error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = serv.dal.UploadObject(bucketname, objectname, r); err != nil {
		log.Printf("Failed to upload object: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = xmlsender.SendMessage(w, http.StatusCreated, fmt.Sprintf("object with name %s created succesfully", objectname)); err != nil {
		log.Printf("Failed to send xml message: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Object delete logic
func (serv ObjectServiceImp) DeleteObject(w http.ResponseWriter, r *http.Request, bucketname, objectname string) {
	unique, err := dal.NewBucketXMLRepo().IsUniqueBucket(bucketname)
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

	if err = Validate(objectname); err != nil {
		log.Printf("Object name validation error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exist, err := serv.dal.IsObjectExist(domain.BucketsPath+"/"+bucketname+"/objects.csv", objectname)
	if err != nil {
		log.Printf("Failed to check if object name exists: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exist {
		log.Printf("Object is not exist")
		http.Error(w, domain.ErrObjectIsNotExist.Error(), http.StatusNotFound)
		return
	}

	if err := serv.dal.DeleteObject(bucketname, objectname); err != nil {
		log.Printf("Failed to delete object : %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = xmlsender.SendMessage(w, http.StatusOK, fmt.Sprintf("object with name %s deleted succesfully", objectname)); err != nil {
		log.Printf("Failed to send xml message: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
