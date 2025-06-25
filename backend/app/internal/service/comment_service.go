package service

import (
	"fmt"
	"io"
	dbrepo "leetFalls/internal/adapters/dbRepo"
	"leetFalls/internal/adapters/storage"
	"leetFalls/internal/domain"
	"leetFalls/internal/domain/models"
	"log/slog"
	"net/http"
	"strconv"
)

type CommentService struct {
	storage     storage.GonIO
	authRepo    dbrepo.AuthRepo
	commentRepo dbrepo.CommentRepo
}

func NewCommentService(authRepo dbrepo.AuthRepo, storage storage.GonIO, commentRepo dbrepo.CommentRepo) *CommentService {
	return &CommentService{authRepo: authRepo, storage: storage, commentRepo: commentRepo}
}

func (s *CommentService) CreateComment(authorId int, postId, commentReplyId, content string, file io.Reader) (domain.Code, error) {
	var (
		comm models.Comment
		err  error
	)

	// Comment Validation
	if comm.PostID, err = strconv.Atoi(postId); err != nil {
		return http.StatusBadRequest, domain.ErrInvalidPostId
	}

	if commentReplyId != "" {
		if comm.ReplyToID, err = strconv.Atoi(commentReplyId); err != nil {
			return http.StatusBadRequest, domain.ErrInvalidReplyId
		}
	}

	comm.Content = content
	if err := ValidateComment(comm); err != nil {
		return http.StatusBadRequest, err
	}

	// User data parsing
	user, err := s.authRepo.GetUserById(authorId)
	if err != nil {
		slog.Error("Failed to get user data by id: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	if user.ID == 0 {
		slog.Error("Failed to get user data by id: ", "error", domain.ErrUserNotExist)
		return http.StatusUnauthorized, domain.ErrUserNotExist
	}
	comm.Author.ID = user.ID

	// Get next Comment ID
	comm.ID, err = s.commentRepo.GetNextCommentId()
	if err != nil {
		slog.Error("Failed to get next comment id: ", "error", err)
		return http.StatusInternalServerError, err
	}

	// Comment Image save
	if file != nil {
		if err := s.storage.SaveCommentImage(&comm, file); err != nil {
			slog.Error("Failed to save post image to storage: ", "error", err)
			return http.StatusInternalServerError, err
		}
	}

	// Comment database save
	err = s.commentRepo.SaveComment(comm)
	if err != nil {
		slog.Error("Failed to save comment data: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	slog.Info(fmt.Sprintf("Comment %d on post %d created succesfully", comm.ID, comm.PostID))
	return http.StatusCreated, nil
}

func ValidateComment(comm models.Comment) error {
	if comm.PostID <= 0 {
		return domain.ErrLessPostId
	}
	if comm.ReplyToID < 0 {
		return domain.ErrLessReplyId
	}

	if len(comm.Content) == 0 {
		return domain.ErrEmptyContent
	}

	return nil
}
