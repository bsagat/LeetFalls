package dbrepo

import (
	"database/sql"
	"leetFalls/internal/domain/models"
)

type CommentRepo struct {
	Db *sql.DB
}

func NewCommentRepo(Db *sql.DB) *CommentRepo {
	return &CommentRepo{Db: Db}
}

// Saves comment to database and returns his new id
func (repo *CommentRepo) SaveComment(comm models.Comment) (int, error) {
	var commentId int
	if err := repo.Db.QueryRow(`
	INSERT INTO 
		Comments (post_id, Reply_to, content, Author_id, ImageURL) 
	VALUES 
		($1, $2, $3, $4, $5)
	RETURNING ID;`,
		comm.PostID, comm.ReplyToID, comm.Content, comm.Author.ID, comm.ImageLink,
	).Scan(&commentId); err != nil {
		return 0, err
	}

	return commentId, nil
}

func (repo *CommentRepo) GetCommentsByPost(postID int) ([]models.Comment, error) {
	rows, err := repo.Db.Query(`
		SELECT c.ID, c.Post_id, COALESCE(c.Reply_to, 0)::INTEGER, c.Content, u.Avatar_URL, u.Name, c.Created_at, COALESCE(c.ImageURL,'')  FROM Comments c
		INNER JOIN Users u On c.Author_id=u.ID
		WHERE c.Post_id = $1 ORDER BY c.Created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.ReplyToID, &comment.Content, &comment.Author.ImageURL, &comment.Author.Name, &comment.CreatedAt, &comment.ImageLink)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
