package app

import (
	"GonIO/internal/dal"
	"GonIO/internal/handlers"
	"GonIO/internal/service"
	"net/http"
)

func SetHandler() *http.ServeMux {
	mux := http.NewServeMux()

	objectDal := dal.NewObjectCSVRepo()
	objectServ := service.NewObjectService(*objectDal)
	objectHandler := handlers.NewObjectHandler(objectServ)

	bucketDal := dal.NewBucketXMLRepo()
	bucketServ := service.NewBucketService(*bucketDal)
	bucketHandler := handlers.NewBucketHandler(bucketServ)

	healthHandler := handlers.NewHealthHandler()

	mux.HandleFunc("GET /PING", healthHandler.Ping) // Healthcheck

	mux.HandleFunc("GET /", bucketHandler.BucketListsHandler)                 // Bucket list
	mux.HandleFunc("PUT /{BucketName}", bucketHandler.CreateBucketHandler)    // Create bucket
	mux.HandleFunc("DELETE /{BucketName}", bucketHandler.DeleteBucketHandler) // Delete bucket

	mux.HandleFunc("GET /{BucketName}", objectHandler.GetObjectList)               // Get object list in bucket
	mux.HandleFunc("GET /{BucketName}/{ObjectKey}", objectHandler.RetrieveObject)  // Retrieve an object
	mux.HandleFunc("PUT /{BucketName}/{ObjectKey}", objectHandler.UpdateObject)    // Upload an object
	mux.HandleFunc("DELETE /{BucketName}/{ObjectKey}", objectHandler.DeleteObject) // Delete an object

	return mux
}
