package domain

import "errors"

var ErrPortRange = errors.New("port flag range must be in 1025 and 65000")
var ErrInvalidPort = errors.New("invalid port number")
var ErrUnsupportMIME = errors.New("unsupported MIME type")
var ErrUserNameDoubleGyphen = errors.New("user name must not contain double hyphens (--)")
