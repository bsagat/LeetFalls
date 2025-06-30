package ports

import "leetFalls/internal/domain/models"

/* ----------Comment repository ports---------- */

type CommentRepo interface {
	CommentGetter
	CommentChecker
	CommentSaver
}

type CommentGetter interface {
	CommentsByPost(postID int) ([]models.Comment, error)
	NextCommentId() (int, error)
}

type CommentChecker interface {
	IsCommentExist(postID int, commentID int) (bool, error)
}

type CommentSaver interface {
	SaveComment(comm models.Comment) error
}

/* ----------Post repository ports---------- */

type PostsRepo interface {
	PostGetter
	PostChecker
	PostSaver
	PostModificator
}

type PostGetter interface {
	ActivePosts() ([]models.Post, error)
	ArchivePosts() ([]models.Post, error)
	GetPost(id int) (models.Post, error)
	GetNextPostId() (int, error)
}

type PostSaver interface {
	SavePost(post models.Post) error
}

type PostChecker interface {
	IsPostExist(postID int) (bool, error)
}

type PostModificator interface {
	AddExpirationTime(postId int, minutes int) error
}

/* ----------User repository ports---------- */

type UserRepo interface {
	UserModificator
	UserGetter
	UserSaver
}

type UserGetter interface {
	GetUserIDBySession(sessionID string) (int, error)
	GetUserById(id int) (models.User, error)
	GetNextUserId() (int, error)
}

type UserSaver interface {
	SaveUser(user models.User) error
}

type UserModificator interface {
	ChangeUserName(id int, changedName string) error
}
