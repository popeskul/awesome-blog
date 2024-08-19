package usecase

import "errors"

var (
	ErrCommentNotFound     = errors.New("comment not found")
	ErrEmptyTitleOrContent = errors.New("title and content are required")
	ErrPostNotFound        = errors.New("post not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrUnauthorized        = errors.New("unauthorized to modify this post")
	ErrEmptyCredentials    = errors.New("username and password are required")
	ErrUserExists          = errors.New("username already exists")
	ErrCreateUserFailed    = errors.New("failed to create user")
	ErrInvalidPage         = errors.New("invalid page number")
	ErrInvalidLimit        = errors.New("invalid limit number")
	ErrInvalidComment      = errors.New("invalid comment")
)
