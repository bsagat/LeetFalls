package ports

import "leetFalls/internal/domain/models"

type ExternalAPI interface {
	AvatarCount() (int, error)
	SetUser(user *models.User) error
}
