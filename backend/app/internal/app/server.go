package app

import (
	"context"
	dbrepo "leetFalls/internal/adapters/dbRepo"
	"leetFalls/internal/adapters/handlers"
	"leetFalls/internal/service"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

func StartServer(srv *http.Server) {
	go func() {
		slog.Info("Server has been started on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %s", err.Error())
		}
	}()
}

// Listens for system signals (e.g., SIGINT, SIGTERM) to ensure a graceful shutdown of the HTTP server.
func WaitForShutDown(srv *http.Server) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	slog.Info("ShutDown signal received!!!")
	slog.Info("Shutting down HTTP server")

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(context); err != nil {
		slog.Error("Server shutdown failed: ", "error", err.Error())
		return
	}

	slog.Info("HTTP server gracefully stopped")
}

func SetRouter() (*http.ServeMux, func()) {
	db, storage, external := ConnectAdapters()

	commentRepo := dbrepo.NewCommentRepo(db)
	postRepo := dbrepo.NewPostsRepo(db)
	authRepo := dbrepo.NewAuthRepo(db)

	authServ := service.NewAuthService(authRepo, external, 5)
	commentServ := service.NewCommentService(authRepo, storage, commentRepo, postRepo)
	postServ := service.NewPostService(authServ, storage, postRepo, commentRepo)

	catalogH := handlers.NewCatalogHandler(postServ, authServ, commentServ)
	archiveH := handlers.NewArchiveHandler(postServ, authServ)

	mux := http.NewServeMux()
	SetSwagger(mux)

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web"))))

	// Catalog API endpoints
	mux.HandleFunc("GET /catalog", catalogH.ServeMainPage)
	mux.HandleFunc("GET /catalog/post/{id}", catalogH.ServeCatalogPost)
	mux.HandleFunc("GET /catalog/post/new", catalogH.ShowCreatePostForm)

	// Archive API endpoints
	mux.HandleFunc("GET /archive", archiveH.ServeArchivePage)
	mux.HandleFunc("GET /archive/post/{id}", archiveH.ServeArchivePost)

	// New subjects API endpoints
	mux.HandleFunc("POST /submit/post", catalogH.CreatePostHandler)
	mux.HandleFunc("POST /submit/comment", catalogH.CreateCommentHandler)

	cleanup := authServ.Close
	return mux, cleanup
}

func SetSwagger(mux *http.ServeMux) {
	// Swagger File
	swaggerBytes, err := os.ReadFile("swagger.json")
	if err != nil {
		log.Fatal("Failed to read swagger file: ", err)
	}

	// Swagger UI
	mux.HandleFunc("GET /docs/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	// Swagger JSON API
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/openapi+json")
		if _, err := w.Write(swaggerBytes); err != nil {
			slog.Error("Failed to send swagger file", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
