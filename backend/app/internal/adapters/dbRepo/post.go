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
	var postId int
	if err := repo.Db.QueryRow(`
	INSERT INTO 
		Posts (Author_id, Content, Title,ImageURL) 
	VALUES
		($1, $2, $3, $4)
	RETURNING ID;`,
		post.AuthorID, post.Content, post.Title, post.ImageLink,
	).Scan(&postId); err != nil {
		return 0, err
	}

	return postId, nil
}

func (repo *PostsRepo) GetPost(id int) (models.Post, error) {
	var post models.Post
	err := repo.Db.QueryRow(`
		SELECT 
			ID, Title, Content, Author_id, Created_at, Expires_at, coalesce(ImageURL,'empty') 
		FROM 
			Posts
	`, id).Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.ExpiresAt, &post.ImageLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Post{}, nil
		}
		return post, err
	}
	return post, nil
}

func (repo *PostsRepo) ActivePosts() ([]models.Post, error) {
	rows, err := repo.Db.Query(`
	SELECT ID, Title, Content, Author_id, coalesce(ImageURL,'empty') 
	FROM Posts
	WHERE Created_at<Expires_at;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *PostsRepo) ArchivePosts() ([]models.Post, error) {
	rows, err := repo.Db.Query(`
	SELECT ID, Title, Content, Author_id, coalesce(ImageURL,'empty') 
	FROM Posts
	WHERE Created_at>=Expires_at;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
