package service

import (
	"fmt"
	"io"
	"leetFalls/internal/domain"
	"leetFalls/internal/domain/models"
	"leetFalls/internal/domain/ports"
	"log/slog"
	"net/http"
	"strconv"
)

type CommentService struct {
	storage     ports.Storage
	userRepo    ports.UserRepo
	postRepo    ports.PostsRepo
	commentRepo ports.CommentRepo
}

func NewCommentService(userRepo ports.UserRepo, storage ports.Storage, commentRepo ports.CommentRepo, postRepo ports.PostsRepo) *CommentService {
	return &CommentService{userRepo: userRepo, storage: storage, commentRepo: commentRepo, postRepo: postRepo}
}

func (s *CommentService) CreateComment(authorId int, postId, commentReplyId, content string, file io.Reader) (int, error) {
	var (
		comm models.Comment
		err  error
	)

	// 1) Comment Validation - Post
	if comm.PostID, err = strconv.Atoi(postId); err != nil {
		return http.StatusBadRequest, domain.ErrInvalidPostId
	}

	exist, err := s.postRepo.IsPostExist(comm.PostID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !exist {
		return http.StatusBadRequest, domain.ErrPostNotFound
	}

	// 2) Comment Validation - ReplyID
	if commentReplyId != "" {
		if comm.ReplyToID, err = strconv.Atoi(commentReplyId); err != nil {
			return http.StatusBadRequest, domain.ErrInvalidReplyId
		}

		// Check is Reply comment exist
		exist, err := s.commentRepo.IsCommentExist(comm.PostID, comm.ReplyToID)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if !exist {
			return http.StatusBadRequest, domain.ErrCommentNotFound
		}
	}

	// 3) Comment Validation - other fields check
	comm.Content = content
	if err := ValidateComment(comm); err != nil {
		return http.StatusBadRequest, err
	}

	// 4) Comment author data parsing
	user, err := s.userRepo.GetUserById(authorId)
	if err != nil {
		slog.Error("Failed to get user data by id: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	if user.ID == 0 {
		slog.Error("Failed to get user data by id: ", "error", domain.ErrUserNotExist)
		return http.StatusUnauthorized, domain.ErrUserNotExist
	}
	comm.Author.ID = user.ID

	// 5) Get unique Comment ID
	comm.ID, err = s.commentRepo.NextCommentId()
	if err != nil {
		slog.Error("Failed to get next comment id: ", "error", err)
		return http.StatusInternalServerError, err
	}

	// 6) Comment Image save
	if file != nil {
		if err := s.storage.SaveCommentImage(&comm, file); err != nil {
			slog.Error("Failed to save post image to storage: ", "error", err)
			return http.StatusInternalServerError, err
		}
	}

	// 7) Comment database save
	if err = s.commentRepo.SaveComment(comm); err != nil {
		slog.Error("Failed to save comment data: ", "error", err.Error())
		return http.StatusInternalServerError, err
	}

	// 8) Updating post TTL (time to live)
	if err = s.postRepo.AddExpirationTime(comm.PostID, 15); err != nil {
		slog.Error("Failed to update post expiration time: ", "error", err)
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
