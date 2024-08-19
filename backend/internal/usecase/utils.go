package usecase

import (
	"errors"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
)

func ValidatePost(post interface{}) error {
	switch p := post.(type) {
	case *entity.NewPost:
		if p.Title == "" || p.Content == "" {
			return ErrEmptyTitleOrContent
		}
	case *entity.Post:
		if p.Title == "" || p.Content == "" {
			return ErrEmptyTitleOrContent
		}
	case *entity.UpdateUser:
		if p.Username == "" || p.Password == "" {
			return ErrEmptyCredentials
		}
	default:
		return errors.New("invalid post type")
	}
	return nil
}
