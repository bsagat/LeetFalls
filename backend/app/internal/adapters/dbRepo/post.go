package dbrepo

import (
	"database/sql"
	"leetFalls/internal/domain/models"
)

type PostsRepo struct {
	Db *sql.DB
}

func NewPostsRepo(Db *sql.DB) *PostsRepo {
	return &PostsRepo{Db: Db}
}

// Saves post to database and returns his new id
func (repo *PostsRepo) SavePost(post models.Post) (int, error) {
	res, err := repo.Db.Exec(`
	INSERT INTO 
		Posts (Author_id, Content, Title,ImageURL) 
	VALUES
		($1, $2, $3, $4)
	RETURNING ID;`, post.AuthorID, post.Content, post.Title, post.ImageLink)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}
