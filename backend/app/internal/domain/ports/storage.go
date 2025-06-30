package ports

import (
	"io"
	"leetFalls/internal/domain/models"
)

type Storage interface {
	SaveCommentImage(comment *models.Comment, img io.Reader) error
	SavePostImage(post *models.Post, img io.Reader) error
	CheckHealth() error
	InitBuckets() error
}
