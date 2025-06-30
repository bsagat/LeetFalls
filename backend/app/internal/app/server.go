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

func SetRouter() *http.ServeMux {
	db, storage, external := ConnectAdapters()

	commentRepo := dbrepo.NewCommentRepo(db)
	postRepo := dbrepo.NewPostsRepo(db)
	authRepo := dbrepo.NewAuthRepo(db)

	authServ := service.NewAuthService(*authRepo, *external)
	commentServ := service.NewCommentService(*authRepo, *storage, *commentRepo, *postRepo)
	postServ := service.NewPostService(*authServ, *storage, *postRepo, *commentRepo)

	catalogH := handlers.NewCatalogHandler(*postServ, *authServ, *commentServ)
	archiveH := handlers.NewArchiveHandler(*postServ, *authServ)
	// profileH := handlers.NewProfileHandler()

	mux := http.NewServeMux()

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates"))))

	// Catalog API endpoints
	mux.HandleFunc("GET /catalog", catalogH.ServeMainPage)
	mux.HandleFunc("GET /catalog/post/{id}", catalogH.ServeCatalogPost)
	mux.HandleFunc("GET /catalog/post/new", catalogH.ShowCreatePostForm)

	// // Profile API endpoints
	// mux.HandleFunc("GET /profile", profileH.ServeProfilePage)
	// mux.HandleFunc("GET /profile/posts", profileH.GetUserPostsHandler)

	// Archive API endpoints
	mux.HandleFunc("GET /archive", archiveH.ServeArchivePage)
	mux.HandleFunc("GET /archive/post/{id}", archiveH.ServeArchivePost)

	// New subjects API endpoints
	mux.HandleFunc("POST /submit/post", catalogH.CreatePostHandler)
	mux.HandleFunc("POST /submit/comment", catalogH.CreateCommentHandler)

	return mux
}
