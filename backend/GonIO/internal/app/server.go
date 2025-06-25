package app

import (
	"GonIO/internal/dal"
	"GonIO/internal/handlers"
	"GonIO/internal/service"
	"log"
	"log/slog"
	"net/http"
	"os"

	httpswagger "github.com/swaggo/http-swagger"
)

func SetHandler() *http.ServeMux {
	mux := http.NewServeMux()

	// Swagger routes
	SetSwagger(mux)

	// DAL and services
	bucketDal := dal.NewBucketXMLRepo()
	objectDal := dal.NewObjectCSVRepo()

	objectServ := service.NewObjectService(*objectDal, bucketDal)
	bucketServ := service.NewBucketService(*bucketDal)

	// Handlers
	objectHandler := handlers.NewObjectHandler(objectServ)
	bucketHandler := handlers.NewBucketHandler(bucketServ)
	healthHandler := handlers.NewHealthHandler()

	// Healthcheck
	mux.HandleFunc("GET /PING", healthHandler.Ping)

	// Bucket routes
	mux.HandleFunc("GET /buckets", bucketHandler.BucketListsHandler)
	mux.HandleFunc("PUT /buckets/{BucketName}", bucketHandler.CreateBucketHandler)
	mux.HandleFunc("DELETE /buckets/{BucketName}", bucketHandler.DeleteBucketHandler)

	// Object routes
	mux.HandleFunc("GET /objects/{BucketName}", objectHandler.GetObjectList)
	mux.HandleFunc("PUT /objects/{BucketName}/jar", objectHandler.ObjectJarHandler)
	mux.HandleFunc("PUT /objects/{BucketName}/{ObjectKey}", objectHandler.UpdateObject)
	mux.HandleFunc("GET /objects/{BucketName}/{ObjectKey}", objectHandler.RetrieveObject)
	mux.HandleFunc("DELETE /objects/{BucketName}/{ObjectKey}", objectHandler.DeleteObject)

	return mux
}

func SetSwagger(mux *http.ServeMux) {
	swaggerBytes, err := os.ReadFile("docs/swagger.json")
	if err != nil {
		log.Fatal("Failed to read swagger file: ", err)
	}

	// Swagger UI
	mux.HandleFunc("/docs/", httpswagger.Handler(
		httpswagger.URL("/swagger.json"),
	))

	// Swagger JSON
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/openapi+json")
		if _, err := w.Write(swaggerBytes); err != nil {
			slog.Error("Failed to send swagger file", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
