package domain

import "errors"

var (
	ErrInvalidPortStr   = errors.New("port number is invalid, must be in range between 1100 and 65535")
	ErrEmptyDomain      = errors.New("domain config is empty")
	ErrEmptyBucketName  = errors.New("bucket name is empty")
	ErrEmptyObjectName  = errors.New("object name is empty")
	ErrObjectIsNotExist = errors.New("object is not exist")
	ErrBucketIsNotExist = errors.New("bucket is not exist")
	ErrBucketIsNotEmpty = errors.New("bucket must be empty")
	ErrNameHyphen       = errors.New("name must not begin or end with a hyphen")
	ErrNameLenght       = errors.New("name should be between 3 and 63 characters long")
	ErrNotUniqueName    = errors.New("name must be unique")
	ErrNameLikeIpAdress = errors.New("name must not be formatted as an IP address ")
	ErrNamePeriodOrDash = errors.New("name must not contain two consecutive periods or dashes")
)
