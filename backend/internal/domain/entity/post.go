package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorId  uuid.UUID `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NewPost struct {
	AuthorId uuid.UUID `json:"authorId" validate:"required"`
	Content  string    `json:"content" validate:"required"`
	Title    string    `json:"title" validate:"required"`
}

type UpdatePost struct {
	AuthorId uuid.UUID `json:"authorId" validate:"required"`
	Content  string    `json:"content" validate:"required"`
	Id       uuid.UUID `json:"id" validate:"required"`
	Title    string    `json:"title"`
}
