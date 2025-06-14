package service

import (
	"leetFalls/internal/domain"
	"net/http"
)

type PostService struct {
}

func NewPostService() *PostService {
	return &PostService{}
}

func (s *PostService) CreatePost(cookie *http.Cookie) (domain.Code, error) {
	return http.StatusOK, nil
}
