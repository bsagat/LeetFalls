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

// Saves post to database
func (repo *PostsRepo) SavePost(post models.Post) error {
	if _, err := repo.Db.Exec(`
	INSERT INTO 
		Posts (ID, Author_id, Content, Title,ImageURL) 
	VALUES
		($1, $2, $3, $4, $5)`,
		post.ID, post.Author.ID, post.Content, post.Title, post.ImageLink,
	); err != nil {
		return err
	}

	return nil
}

// Gets unique post id
func (repo *PostsRepo) GetNextPostId() (int, error) {
	var id int
	if err := repo.Db.QueryRow("SELECT COALESCE(MAX(ID), 0) + 1 FROM Posts").Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *PostsRepo) GetPost(id int) (models.Post, error) {
	var post models.Post
	err := repo.Db.QueryRow(`
		SELECT 
			p.ID, p.Title, p.Content, p.Author_id, p.Created_at, p.Expires_at, coalesce(p.ImageURL,'empty') ,
			u.Name as authorName, u.Avatar_URL as authorAvatar 
		FROM 
			Posts p
		INNER JOIN Users u 
		  ON p.Author_id = u.ID
		WHERE p.ID = $1;
	`, id).Scan(&post.ID, &post.Title, &post.Content, &post.Author.ID, &post.CreatedAt, &post.ExpiresAt, &post.ImageLink, &post.Author.Name, &post.Author.ImageURL)
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
	SELECT 
		p.ID, p.Title, p.Content, p.Author_id, p.Created_at, p.Expires_at, coalesce(p.ImageURL,'empty') ,
		u.Name as authorName, u.Avatar_URL as authorAvatar 
	FROM 
		Posts p
	INNER JOIN Users u 
		ON p.Author_id = u.ID
	WHERE Now() < p.Expires_at;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author.ID, &post.CreatedAt,
			&post.ExpiresAt, &post.ImageLink, &post.Author.Name, &post.Author.ImageURL); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *PostsRepo) ArchivePosts() ([]models.Post, error) {
	rows, err := repo.Db.Query(`
	SELECT 
		p.ID, p.Title, p.Content, p.Author_id, p.Created_at, p.Expires_at, coalesce(p.ImageURL,'empty') ,
		u.Name as authorName, u.Avatar_URL as authorAvatar 
	FROM 
		Posts p
	INNER JOIN Users u 
		ON p.Author_id = u.ID
	WHERE Now() >= p.Expires_at;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author.ID, &post.CreatedAt,
			&post.ExpiresAt, &post.ImageLink, &post.Author.Name, &post.Author.ImageURL); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
