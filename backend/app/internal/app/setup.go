package app

import (
	"database/sql"
	"fmt"
	"io"
	dbrepo "leetFalls/internal/adapters/dbRepo"
	"leetFalls/internal/adapters/external"
	"leetFalls/internal/adapters/handlers"
	"leetFalls/internal/adapters/storage"
	"leetFalls/internal/domain"
	"leetFalls/internal/service"
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
	domain.Config.TemplatesPath = "templates" // Default templates path
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

func SetRouter() *http.ServeMux {
	db, storage, external := ConnectAdapters()

	commentRepo := dbrepo.NewCommentRepo(db)
	postRepo := dbrepo.NewPostsRepo(db)
	authRepo := dbrepo.NewAuthRepo(db)

	authServ := service.NewAuthService(*authRepo, *external)
	commentServ := service.NewCommentService(*authRepo, *storage, *commentRepo)
	postServ := service.NewPostService(*authServ, *storage, *postRepo, *commentRepo)

	catalogH := handlers.NewCatalogHandler(*postServ, *authServ, *commentServ)
	archiveH := handlers.NewArchiveHandler(*postServ, *authServ)
	profileH := handlers.NewProfileHandler()

	mux := http.NewServeMux()

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates"))))

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
