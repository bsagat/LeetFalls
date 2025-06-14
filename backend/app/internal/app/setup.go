package app

import (
	"fmt"
	"io"
	"leetFalls/internal/adapters/handlers"
	"leetFalls/internal/domain"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func Setup() (*http.Server, func()) {
	logFileCloser := SetLogger()

	mux := SetRouter()

	cleanUp := func() {
		logFileCloser()
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
	domain.Config.TemplatesPath = "LeetFalls/frontend/templates" // Default templates path

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
		Level: slog.LevelInfo, // it can be switched to debug level :)
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Logger setup has been finished...")
	return file.Close
}

func SetRouter() *http.ServeMux {
	catalogH := handlers.NewCatalogHandler()
	profileH := handlers.NewProfileHandler()
	archiveH := handlers.NewArchiveHandler()

	mux := http.NewServeMux()

	// Catalog API endpoints
	mux.HandleFunc("GET /catalog", catalogH.ServeMainPage)
	mux.HandleFunc("GET /catalog/post/{id}", catalogH.ServeCatalogPost)
	mux.HandleFunc("GET /catalog/post/new", catalogH.ShowCreatePostForm)

	// Profile API endpoints
	mux.HandleFunc("GET /profile", profileH.ServeProfilePage)
	mux.HandleFunc("GET /profile/posts", profileH.GetUserPostsHandler)

	// Archive API endpoints
	mux.HandleFunc("GET /archive", archiveH.ServeArchivePage)
	mux.HandleFunc("GET /archive/post/{id}", archiveH.ServeArchivePost)

	// New subjects API endpoints
	mux.HandleFunc("POST /submit/post", catalogH.CreatePostHandler)
	mux.HandleFunc("POST /submit/comment", catalogH.CreateCommentHandler)

	return mux
}
