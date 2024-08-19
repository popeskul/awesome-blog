package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/internal/infrastructure/database/postgres"
	"github.com/popeskul/awesome-blog/backend/pkg/db"
)

var (
	userId1 = uuid.New()
	userId2 = uuid.New()

	postId1 = uuid.New()
	postId2 = uuid.New()

	commentId1 = uuid.New()
	commentId2 = uuid.New()
)

func TestCommentRepository_CreateComment_Success(t *testing.T) {
	tests := []struct {
		name        string
		newComment  *entity.NewComment
		expectedID  uuid.UUID
		expectedErr error
	}{
		// ... (оставьте существующие тестовые случаи без изменений)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewCommentRepository(&db.PostgresDB{DB: mockDB}, logger)

			rows := sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at", "updated_at"}).
				AddRow(tt.expectedID, tt.newComment.PostId, tt.newComment.AuthorId, tt.newComment.Content, time.Now(), time.Now())

			mock.ExpectQuery("INSERT INTO comments").
				WithArgs(sqlmock.AnyArg(), tt.newComment.PostId, tt.newComment.AuthorId, tt.newComment.Content).
				WillReturnRows(rows)

			comment, err := repo.CreateComment(context.Background(), tt.newComment)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				assert.NotNil(t, comment)
				assert.Equal(t, tt.expectedID, comment.Id)
				assert.Equal(t, tt.newComment.PostId, comment.PostId)
				assert.Equal(t, tt.newComment.AuthorId, comment.AuthorId)
				assert.Equal(t, tt.newComment.Content, comment.Content)
				assert.NotZero(t, comment.CreatedAt)
				assert.NotZero(t, comment.UpdatedAt)
			} else {
				assert.Nil(t, comment)
			}
		})
	}
}

func TestCommentRepository_CreateComment_Failed(t *testing.T) {
	tests := []struct {
		name        string
		newComment  *entity.NewComment
		expectedErr error
	}{
		{
			name: "Failed comment creation due to database error",
			newComment: &entity.NewComment{
				PostId:   postId1,
				AuthorId: userId1,
				Content:  "Test comment",
			},
			expectedErr: errors.New("database error"),
		},
		{
			name: "Failed comment creation due to constraint violation",
			newComment: &entity.NewComment{
				PostId:   postId2,
				AuthorId: userId2,
				Content:  "Another test comment",
			},
			expectedErr: errors.New("constraint violation"),
		},
		{
			name: "Failed comment creation with empty content",
			newComment: &entity.NewComment{
				PostId:   postId1,
				AuthorId: userId1,
				Content:  "",
			},
			expectedErr: errors.New("content cannot be empty"),
		},
		{
			name: "Failed comment creation with non-existent post ID",
			newComment: &entity.NewComment{
				PostId:   postId2,
				AuthorId: userId2,
				Content:  "Comment for non-existent post",
			},
			expectedErr: errors.New("foreign key violation"),
		},
		{
			name: "Failed comment creation with non-existent author ID",
			newComment: &entity.NewComment{
				PostId:   postId1,
				AuthorId: userId1,
				Content:  "Comment from non-existent author",
			},
			expectedErr: errors.New("foreign key violation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewCommentRepository(&db.PostgresDB{DB: mockDB}, logger)

			mock.ExpectQuery("INSERT INTO comments").
				WithArgs(sqlmock.AnyArg(), tt.newComment.PostId, tt.newComment.AuthorId, tt.newComment.Content).
				WillReturnError(tt.expectedErr)

			comment, err := repo.CreateComment(context.Background(), tt.newComment)

			assert.Error(t, err)
			assert.Nil(t, comment)
			assert.Contains(t, err.Error(), "failed to create comment")
			assert.Contains(t, err.Error(), tt.expectedErr.Error())
		})
	}
}

func TestCommentRepository_GetCommentById_Success(t *testing.T) {
	tests := []struct {
		name            string
		commentID       uuid.UUID
		expectedComment *entity.Comment
		expectedErr     error
	}{
		{
			name:      "Successful get comment by ID",
			commentID: commentId1,
			expectedComment: &entity.Comment{
				Id:       commentId1,
				PostId:   postId1,
				AuthorId: userId1,
				Content:  "Test comment",
			},
			expectedErr: nil,
		},
		{
			name:      "Successful get comment with long content",
			commentID: commentId2,
			expectedComment: &entity.Comment{
				Id:       commentId2,
				PostId:   postId2,
				AuthorId: userId2,
				Content:  "This is a much longer comment that exceeds the typical length of a short comment.",
			},
			expectedErr: nil,
		},
		{
			name:      "Successful get comment with special characters",
			commentID: commentId1,
			expectedComment: &entity.Comment{
				Id:       commentId1,
				PostId:   postId1,
				AuthorId: userId1,
				Content:  "Comment with special chars: !@#$%^&*()_+-={}[]|\\:;\"'<>,.?/",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewCommentRepository(&db.PostgresDB{DB: mockDB}, logger)

			rows := sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at", "updated_at"}).
				AddRow(tt.expectedComment.Id, tt.expectedComment.PostId, tt.expectedComment.AuthorId, tt.expectedComment.Content, time.Now(), time.Now())

			mock.ExpectQuery("SELECT (.+) FROM comments WHERE id = \\$1").
				WithArgs(tt.commentID).
				WillReturnRows(rows)

			comment, err := repo.GetCommentById(context.Background(), tt.commentID)

			assert.Equal(t, tt.expectedErr, err)
			assert.NotNil(t, comment)
			assert.Equal(t, tt.expectedComment.Id, comment.Id)
			assert.Equal(t, tt.expectedComment.PostId, comment.PostId)
			assert.Equal(t, tt.expectedComment.AuthorId, comment.AuthorId)
			assert.Equal(t, tt.expectedComment.Content, comment.Content)
		})
	}
}

func TestCommentRepository_GetCommentsByPostIdWithPagination_Success(t *testing.T) {
	tests := []struct {
		name        string
		postID      uuid.UUID
		pagination  *entity.Pagination
		expectedLen int
		totalCount  int
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, postID uuid.UUID, pagination *entity.Pagination, expectedLen int, totalCount int)
	}{
		{
			name:   "Successful get comments by post ID with pagination - first page",
			postID: postId1,
			pagination: &entity.Pagination{
				Page:  1,
				Limit: 10,
			},
			expectedLen: 2,
			totalCount:  15,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, postID uuid.UUID, pagination *entity.Pagination, expectedLen int, totalCount int) {
				rows := sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at", "updated_at"})
				for i := 0; i < expectedLen; i++ {
					rows.AddRow(uuid.New(), postID, uuid.New(), fmt.Sprintf("Test comment %d", i+1), time.Now(), time.Now())
				}

				offset := (pagination.Page - 1) * pagination.Limit
				mock.ExpectQuery("SELECT (.+) FROM comments WHERE post_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
					WithArgs(postID, pagination.Limit, offset).
					WillReturnRows(rows)

				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM comments WHERE post_id = \\$1").
					WithArgs(postID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(totalCount))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewCommentRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.setupMock(mock, tt.postID, tt.pagination, tt.expectedLen, tt.totalCount)

			commentList, err := repo.GetComments(context.Background(), tt.postID, tt.pagination)

			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				assert.NotNil(t, commentList)
				assert.Len(t, commentList.Comments, tt.expectedLen)
				assert.Equal(t, int64(tt.totalCount), commentList.Pagination.Total)
				for _, comment := range commentList.Comments {
					assert.NotEqual(t, uuid.Nil, comment.Id)
					assert.Equal(t, tt.postID, comment.PostId)
					assert.NotEqual(t, uuid.Nil, comment.AuthorId)
					assert.NotEmpty(t, comment.Content)
					assert.NotZero(t, comment.CreatedAt)
					assert.NotZero(t, comment.UpdatedAt)
				}
			} else {
				assert.Nil(t, commentList)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCommentRepository_GetCommentsByPostIdWithPagination_Failed(t *testing.T) {
	tests := []struct {
		name        string
		postID      uuid.UUID
		pagination  *entity.Pagination
		expectedErr error
	}{
		{
			name:   "Failed get comments by post ID with pagination - database error",
			postID: postId1,
			pagination: &entity.Pagination{
				Page:  1,
				Limit: 10,
			},
			expectedErr: errors.New("database error"),
		},
		{
			name:   "Failed get comments by post ID with pagination - invalid post ID",
			postID: uuid.Nil,
			pagination: &entity.Pagination{
				Page:  1,
				Limit: 10,
			},
			expectedErr: errors.New("invalid post ID"),
		},
		{
			name:   "Failed get comments by post ID with pagination - count query error",
			postID: postId2,
			pagination: &entity.Pagination{
				Page:  1,
				Limit: 10,
			},
			expectedErr: errors.New("count query error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewCommentRepository(&db.PostgresDB{DB: mockDB}, logger)

			if tt.name == "Failed get comments by post ID with pagination - count query error" {
				mock.ExpectQuery("SELECT (.+) FROM comments WHERE post_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
					WithArgs(tt.postID, tt.pagination.Limit, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at", "updated_at"}))

				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM comments WHERE post_id = \\$1").
					WithArgs(tt.postID).
					WillReturnError(tt.expectedErr)
			} else {
				mock.ExpectQuery("SELECT (.+) FROM comments WHERE post_id = \\$1 ORDER BY created_at DESC LIMIT \\$2 OFFSET \\$3").
					WithArgs(tt.postID, tt.pagination.Limit, 0).
					WillReturnError(tt.expectedErr)
			}

			commentList, err := repo.GetComments(context.Background(), tt.postID, tt.pagination)

			assert.Error(t, err)
			assert.Nil(t, commentList)
			assert.Contains(t, err.Error(), tt.expectedErr.Error())
		})
	}
}

func TestCommentRepository_UpdateComment_Success(t *testing.T) {
	tests := []struct {
		name        string
		comment     *entity.UpdateComment
		expectedErr error
	}{
		{
			name: "Successful comment update - content only",
			comment: &entity.UpdateComment{
				Id:      commentId1,
				Content: "Updated comment content",
			},
			expectedErr: nil,
		},
		{
			name: "Successful comment update - all fields",
			comment: &entity.UpdateComment{
				Id:       commentId2,
				AuthorId: userId2,
				Content:  "Fully updated comment",
			},
			expectedErr: nil,
		},
		{
			name: "Successful comment update - empty content",
			comment: &entity.UpdateComment{
				Id:      commentId1,
				Content: "",
			},
			expectedErr: nil,
		},
		{
			name: "Successful comment update - long content",
			comment: &entity.UpdateComment{
				Id:      commentId2,
				Content: strings.Repeat("Long content ", 100),
			},
			expectedErr: nil,
		},
		{
			name: "Successful comment update - special characters",
			comment: &entity.UpdateComment{
				Id:      commentId1,
				Content: "Special characters: !@#$%^&*()_+|",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewCommentRepository(&db.PostgresDB{DB: mockDB}, logger)

			mock.ExpectExec("UPDATE comments SET content = \\$1, updated_at = NOW\\(\\) WHERE id = \\$2").
				WithArgs(tt.comment.Content, tt.comment.Id).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = repo.UpdateComment(context.Background(), tt.comment)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestCommentRepository_UpdateComment_Failed(t *testing.T) {
	tests := []struct {
		name        string
		comment     *entity.UpdateComment
		expectedErr error
	}{
		{
			name: "Failed comment update - SQL error",
			comment: &entity.UpdateComment{
				Id:      commentId1,
				Content: "Updated comment",
			},
			expectedErr: errors.New("failed to update comment"),
		},
		{
			name: "Failed comment update - missing id",
			comment: &entity.UpdateComment{
				Id:      uuid.Nil,
				Content: "Updated comment",
			},
			expectedErr: errors.New("failed to update comment"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewCommentRepository(&db.PostgresDB{DB: mockDB}, logger)

			mock.ExpectExec("UPDATE comments SET content = \\$1, updated_at = NOW\\(\\) WHERE id = \\$2").
				WithArgs(tt.comment.Content, tt.comment.Id).
				WillReturnError(tt.expectedErr)

			err = repo.UpdateComment(context.Background(), tt.comment)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr.Error())
		})
	}
}

func TestCommentRepository_DeleteCommentById_Success(t *testing.T) {
	tests := []struct {
		name        string
		commentID   uuid.UUID
		expectedErr error
	}{
		{
			name:        "Successful comment deletion",
			commentID:   commentId1,
			expectedErr: nil,
		},
		{
			name:        "Successful comment deletion - non-existing ID",
			commentID:   uuid.Nil,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewCommentRepository(&db.PostgresDB{DB: mockDB}, logger)

			mock.ExpectExec("DELETE FROM comments WHERE id = \\$1").
				WithArgs(tt.commentID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = repo.DeleteCommentById(context.Background(), tt.commentID)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestCommentRepository_DeleteCommentById_Failed(t *testing.T) {
	tests := []struct {
		name        string
		commentID   uuid.UUID
		expectedErr error
	}{
		{
			name:        "Failed comment deletion - SQL error",
			commentID:   commentId1,
			expectedErr: errors.New("failed to delete comment"),
		},
		{
			name:        "Failed comment deletion - non-existent ID",
			commentID:   uuid.Nil,
			expectedErr: errors.New("failed to delete comment"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewCommentRepository(&db.PostgresDB{DB: mockDB}, logger)

			mock.ExpectExec("DELETE FROM comments WHERE id = \\$1").
				WithArgs(tt.commentID).
				WillReturnError(tt.expectedErr)

			err = repo.DeleteCommentById(context.Background(), tt.commentID)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr.Error())
		})
	}
}

func TestCommentRepository_GetTotalCommentsByPostID_Success(t *testing.T) {
	tests := []struct {
		name        string
		postID      uuid.UUID
		totalCount  int
		expectedErr error
	}{
		{
			name:        "Successful get total comments count - 0 comments",
			postID:      postId1,
			totalCount:  0,
			expectedErr: nil,
		},
		{
			name:        "Successful get total comments count - 1 comment",
			postID:      postId2,
			totalCount:  1,
			expectedErr: nil,
		},
		{
			name:        "Successful get total comments count - multiple comments",
			postID:      postId1,
			totalCount:  10,
			expectedErr: nil,
		},
		{
			name:        "Successful get total comments count - large number",
			postID:      postId2,
			totalCount:  1000,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewCommentRepository(&db.PostgresDB{DB: mockDB}, logger)

			mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM comments WHERE post_id = \\$1").
				WithArgs(tt.postID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.totalCount))

			total, err := repo.GetTotalCommentsByPostID(context.Background(), tt.postID)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.totalCount, total)
		})
	}
}

func TestCommentRepository_GetTotalCommentsByPostID_Failed(t *testing.T) {
	tests := []struct {
		name        string
		postID      uuid.UUID
		expectedErr error
	}{
		{
			name:        "Failed get total comments count - SQL error",
			postID:      postId1,
			expectedErr: errors.New("failed to get total comments"),
		},
		{
			name:        "Failed get total comments count - invalid query",
			postID:      postId2,
			expectedErr: errors.New("invalid query"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewCommentRepository(&db.PostgresDB{DB: mockDB}, logger)

			mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM comments WHERE post_id = \\$1").
				WithArgs(tt.postID).
				WillReturnError(tt.expectedErr)

			total, err := repo.GetTotalCommentsByPostID(context.Background(), tt.postID)

			assert.Error(t, err)
			assert.Zero(t, total)
			assert.Contains(t, err.Error(), tt.expectedErr.Error())
		})
	}
}
