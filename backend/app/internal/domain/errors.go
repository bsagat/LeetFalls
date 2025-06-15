package domain

import "errors"

var (
	ErrPortRange            = errors.New("port flag range must be in 1025 and 65000")
	ErrInvalidPort          = errors.New("invalid port number")
	ErrUnsupportMIME        = errors.New("unsupported MIME type")
	ErrUserNameDoubleGyphen = errors.New("user name must not contain double hyphens (--)")
	ErrUserNotExist         = errors.New("user is not found from database")
	ErrInvalidPostId        = errors.New("post id is invalid: must be integer")
	ErrInvalidReplyId       = errors.New("comment reply id is invalid: must be integer")
	ErrEmptyContent         = errors.New("comment content is empty")
	ErrLessReplyId          = errors.New("reply id must be more or equal to 0")
	ErrLessPostId           = errors.New("post id must be more than 0")
	ErrPostNotFound         = errors.New("post is not found")
)
