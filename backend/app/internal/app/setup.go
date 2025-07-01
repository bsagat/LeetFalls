package app

import (
	"database/sql"
	"fmt"
	"io"
	dbrepo "leetFalls/internal/adapters/dbRepo"
	"leetFalls/internal/adapters/external"
	"leetFalls/internal/adapters/storage"
	"leetFalls/internal/domain"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func Setup() (*http.Server, func()) {
	logFileCloser := SetLogger()

	mux, servCleanup := SetRouter()

	cleanUp := func() {
		logFileCloser()
		servCleanup()
	}

	srv := &http.Server{
		Addr:    ":" + *domain.Config.Port,
		Handler: mux,
	}

	return srv, cleanUp
}

func SetConfigs() {
	CheckFlags()

	domain.Config.DatabaseDsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	domain.Config.TemplatesPath = "web" // Default templates path
	domain.Config.StorageHost = os.Getenv("S3_HOST")
	domain.Config.StoragePort = os.Getenv("S3_PORT")
	domain.Config.GravityFallsHost = os.Getenv("GRAVITYFALLS_HOST")
	domain.Config.GravityFallsPort = os.Getenv("GRAVITYFALLS_PORT")
}

func SetLogger() func() error {
	logFilePath := "logs/app.log"
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err.Error())
	}

	writer := io.MultiWriter(os.Stdout, file)
	// It can be changed:
	// Why? Because it duplicates all log output to multiple destinations,
	// but it does not allow you to configure different log levels (INFO,DEBUG,WARN...)

	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelInfo, // it can be switched to debug level
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Logger setup has been finished...")
	return file.Close
}

func ConnectAdapters() (*sql.DB, *storage.GonIO, *external.GravityFallsAPI) {
	storage, err := storage.InitStorage(domain.Config.StorageHost, domain.Config.StoragePort)
	if err != nil {
		log.Fatal("Failed to connect s3 storage: ", err)
	}

	db, err := dbrepo.Connect()
	if err != nil {
		log.Fatal("Failed to connect Database: ", err)
	}

	gravityFallsURL := domain.Config.GravityFallsHost + ":" + domain.Config.GravityFallsPort
	external := external.NewGravityFallsAPI(gravityFallsURL, domain.Config.StoragePort)

	return db, storage, external
}
