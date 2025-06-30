package ports

import (
	"io"
	"leetFalls/internal/domain/models"
	"mime/multipart"
	"net/http"
)

type PostService interface {
	CreatePost(w http.ResponseWriter, userName string, title string, content string, file multipart.File, cookie *http.Cookie) (int, error)
	ShowPost(w http.ResponseWriter, postId string, archive bool) (int, error)
	ShowArchiveWithPosts(w http.ResponseWriter) (int, error)
	ShowCatalogWithPosts(w http.ResponseWriter) (int, error)
}

type CommentService interface {
	CreateComment(authorId int, postId string, commentReplyId string, content string, file io.Reader) (int, error)
}

type AuthService interface {
	Auth(w http.ResponseWriter, cookie *http.Cookie) (int, error)
	ChangeUserName(userId int, userName string) error
	CreateNewUser(sessionId string) (int, error)
	GetUserById(id int) (models.User, error)
}
