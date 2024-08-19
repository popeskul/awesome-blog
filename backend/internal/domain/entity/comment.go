package entity

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	AuthorId  uuid.UUID `json:"authorId"`
	PostId    uuid.UUID `json:"postId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NewComment struct {
	AuthorId uuid.UUID `json:"authorId" validate:"required"`
	Content  string    `json:"content" validate:"required,max=1000"`
	PostId   uuid.UUID `json:"postId" validate:"required"`
}

type UpdateComment struct {
	Id       uuid.UUID `json:"id" validate:"required"`
	AuthorId uuid.UUID `json:"authorId" validate:"required"`
	Content  string    `json:"content" validate:"required,max=1000"`
}
