package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/internal/infrastructure/database/postgres"
	"github.com/popeskul/awesome-blog/backend/pkg/db"
)

var (
	authorId1 = uuid.New()
	authorId2 = uuid.New()
)

func TestPostRepository_CreatePost_Success(t *testing.T) {
	tests := []struct {
		name         string
		post         *entity.NewPost
		expectedPost *entity.Post
		setupMocks   func(mock sqlmock.Sqlmock, post *entity.NewPost)
	}{
		{
			name: "Successful post creation",
			post: &entity.NewPost{
				Title:    "Test Title",
				Content:  "Test Content",
				AuthorId: authorId1,
			},
			expectedPost: &entity.Post{
				Id:        postId1,
				Title:     "Test Title",
				Content:   "Test Content",
				AuthorId:  authorId1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setupMocks: func(mock sqlmock.Sqlmock, post *entity.NewPost) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
					AddRow(postId1, post.Title, post.Content, post.AuthorId, time.Now(), time.Now())

				mock.ExpectQuery(`INSERT INTO posts \(id, title, content, author_id, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, NOW\(\), NOW\(\)\) RETURNING id, title, content, author_id, created_at, updated_at`).
					WithArgs(sqlmock.AnyArg(), post.Title, post.Content, post.AuthorId).
					WillReturnRows(rows)
			},
		},
		{
			name: "Successful post creation with missing fields",
			post: &entity.NewPost{
				Title:    "Incomplete Post",
				Content:  "",
				AuthorId: authorId2,
			},
			expectedPost: &entity.Post{
				Id:        postId2,
				Title:     "Incomplete Post",
				Content:   "",
				AuthorId:  authorId2,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setupMocks: func(mock sqlmock.Sqlmock, post *entity.NewPost) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
					AddRow(postId2, post.Title, post.Content, post.AuthorId, time.Now(), time.Now())

				mock.ExpectQuery(`INSERT INTO posts \(id, title, content, author_id, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, NOW\(\), NOW\(\)\) RETURNING id, title, content, author_id, created_at, updated_at`).
					WithArgs(sqlmock.AnyArg(), post.Title, post.Content, post.AuthorId).
					WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewPostRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.setupMocks(mock, tt.post)

			post, err := repo.CreatePost(context.Background(), tt.post)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedPost.Title, post.Title)
			assert.Equal(t, tt.expectedPost.Content, post.Content)
			assert.Equal(t, tt.expectedPost.AuthorId, post.AuthorId)
			assert.NotEqual(t, uuid.Nil, post.Id)
		})
	}
}

func TestPostRepository_CreatePost_Failed(t *testing.T) {
	tests := []struct {
		name         string
		post         *entity.NewPost
		expectedErr  error
		mockBehavior func(mock sqlmock.Sqlmock, post *entity.NewPost)
	}{
		{
			name: "Failed to create post - SQL error",
			post: &entity.NewPost{
				Title:    "Test Title",
				Content:  "Test Content",
				AuthorId: authorId1,
			},
			expectedErr: errors.New("failed to create post"),
			mockBehavior: func(mock sqlmock.Sqlmock, post *entity.NewPost) {
				mock.ExpectQuery(`INSERT INTO posts \(id, title, content, author_id, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, NOW\(\), NOW\(\)\) RETURNING id, title, content, author_id, created_at, updated_at`).
					WithArgs(sqlmock.AnyArg(), post.Title, post.Content, post.AuthorId).
					WillReturnError(errors.New("failed to create post"))
			},
		},
		{
			name: "Failed to create post - Unique constraint violation",
			post: &entity.NewPost{
				Title:    "Unique Title",
				Content:  "Content with unique constraint",
				AuthorId: authorId2,
			},
			expectedErr: errors.New("failed to create post"),
			mockBehavior: func(mock sqlmock.Sqlmock, post *entity.NewPost) {
				mock.ExpectQuery(`INSERT INTO posts \(id, title, content, author_id, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, NOW\(\), NOW\(\)\) RETURNING id, title, content, author_id, created_at, updated_at`).
					WithArgs(sqlmock.AnyArg(), post.Title, post.Content, post.AuthorId).
					WillReturnError(errors.New("unique constraint violation"))
			},
		},
		{
			name: "Failed to create post - Data type mismatch",
			post: &entity.NewPost{
				Title:    "Test Title",
				Content:  "Content with type mismatch",
				AuthorId: authorId1,
			},
			expectedErr: errors.New("failed to create post"),
			mockBehavior: func(mock sqlmock.Sqlmock, post *entity.NewPost) {
				mock.ExpectQuery(`INSERT INTO posts \(id, title, content, author_id, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, NOW\(\), NOW\(\)\) RETURNING id, title, content, author_id, created_at, updated_at`).
					WithArgs(sqlmock.AnyArg(), post.Title, post.Content, post.AuthorId).
					WillReturnError(errors.New("type mismatch"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewPostRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockBehavior(mock, tt.post)

			post, err := repo.CreatePost(context.Background(), tt.post)
			assert.Error(t, err)
			assert.Nil(t, post)
			assert.Contains(t, err.Error(), tt.expectedErr.Error())
		})
	}
}

func TestPostRepository_GetPostById_Success(t *testing.T) {
	tests := []struct {
		name         string
		id           uuid.UUID
		expectedPost *entity.Post
		setupMocks   func(mock sqlmock.Sqlmock, id uuid.UUID, post *entity.Post)
	}{
		{
			name: "Successful get post by ID",
			id:   postId1,
			expectedPost: &entity.Post{
				Id:        postId1,
				Title:     "Test Title",
				Content:   "Test Content",
				AuthorId:  authorId1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setupMocks: func(mock sqlmock.Sqlmock, id uuid.UUID, post *entity.Post) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
					AddRow(post.Id, post.Title, post.Content, post.AuthorId, post.CreatedAt, post.UpdatedAt)

				mock.ExpectQuery(`SELECT id, title, content, author_id, created_at, updated_at FROM posts WHERE id = \$1`).
					WithArgs(id).
					WillReturnRows(rows)
			},
		},
		{
			name: "Get post with different ID",
			id:   postId2,
			expectedPost: &entity.Post{
				Id:        postId2,
				Title:     "Another Title",
				Content:   "Another Content",
				AuthorId:  authorId2,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setupMocks: func(mock sqlmock.Sqlmock, id uuid.UUID, post *entity.Post) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
					AddRow(post.Id, post.Title, post.Content, post.AuthorId, post.CreatedAt, post.UpdatedAt)

				mock.ExpectQuery(`SELECT id, title, content, author_id, created_at, updated_at FROM posts WHERE id = \$1`).
					WithArgs(id).
					WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewPostRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.setupMocks(mock, tt.id, tt.expectedPost)

			post, err := repo.GetPostById(context.Background(), tt.id)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedPost.Title, post.Title)
			assert.Equal(t, tt.expectedPost.Content, post.Content)
			assert.Equal(t, tt.expectedPost.AuthorId, post.AuthorId)
		})
	}
}

func TestPostRepository_GetPostById_Failed(t *testing.T) {
	tests := []struct {
		name        string
		id          uuid.UUID
		setupMocks  func(mock sqlmock.Sqlmock, id uuid.UUID, err error)
		expectedErr string
	}{
		{
			name: "Failed to get post by ID - not found",
			id:   postId1,
			setupMocks: func(mock sqlmock.Sqlmock, id uuid.UUID, err error) {
				mock.ExpectQuery(`SELECT id, title, content, author_id, created_at, updated_at FROM posts WHERE id = \$1`).
					WithArgs(id).
					WillReturnError(sql.ErrNoRows)
			},
			expectedErr: "post not found",
		},
		{
			name: "Failed to get post by ID - SQL error",
			id:   postId2,
			setupMocks: func(mock sqlmock.Sqlmock, id uuid.UUID, err error) {
				mock.ExpectQuery(`SELECT id, title, content, author_id, created_at, updated_at FROM posts WHERE id = \$1`).
					WithArgs(id).
					WillReturnError(err)
			},
			expectedErr: "failed to get post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewPostRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.setupMocks(mock, tt.id, errors.New(tt.expectedErr))

			post, err := repo.GetPostById(context.Background(), tt.id)
			assert.Error(t, err)
			assert.Nil(t, post)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestPostRepository_GetAllByParams_Success(t *testing.T) {
	tests := []struct {
		name          string
		params        *entity.Pagination
		expectedPosts []*entity.Post
	}{
		{
			name: "Successful get all posts",
			params: &entity.Pagination{
				Limit:  10,
				Offset: 0,
			},
			expectedPosts: []*entity.Post{
				{
					Id:        postId1,
					Title:     "Post 1",
					Content:   "Content 1",
					AuthorId:  authorId1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					Id:        postId2,
					Title:     "Post 2",
					Content:   "Content 2",
					AuthorId:  authorId2,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		},
		{
			name: "Successful get posts with pagination",
			params: &entity.Pagination{
				Limit:  5,
				Offset: 5,
			},
			expectedPosts: []*entity.Post{
				{
					Id:        postId1,
					Title:     "Post 3",
					Content:   "Content 3",
					AuthorId:  authorId1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewPostRepository(&db.PostgresDB{DB: mockDB}, logger)

			rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"})
			for _, post := range tt.expectedPosts {
				rows.AddRow(post.Id, post.Title, post.Content, post.AuthorId, post.CreatedAt, post.UpdatedAt)
			}

			mock.ExpectQuery(`SELECT id, title, content, author_id, created_at, updated_at FROM posts ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
				WithArgs(tt.params.Limit, tt.params.Offset).
				WillReturnRows(rows)

			posts, err := repo.GetAll(context.Background(), tt.params)
			assert.NoError(t, err)
			assert.Len(t, posts, len(tt.expectedPosts))
			for i, post := range tt.expectedPosts {
				assert.Equal(t, post.Title, posts[i].Title)
				assert.Equal(t, post.Content, posts[i].Content)
			}
		})
	}
}

func TestPostRepository_GetAllByParams_Failed(t *testing.T) {
	tests := []struct {
		name        string
		params      *entity.Pagination
		expectedErr string
	}{
		{
			name: "Failed to get all posts - SQL error",
			params: &entity.Pagination{
				Limit:  10,
				Offset: 0,
			},
			expectedErr: "failed to get posts",
		},
		{
			name: "Failed to get all posts - No rows error",
			params: &entity.Pagination{
				Limit:  5,
				Offset: 10,
			},
			expectedErr: "no rows in result set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewPostRepository(&db.PostgresDB{DB: mockDB}, logger)

			if tt.expectedErr == "no rows in result set" {
				mock.ExpectQuery(`SELECT id, title, content, author_id, created_at, updated_at FROM posts ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
					WithArgs(tt.params.Limit, tt.params.Offset).
					WillReturnError(sql.ErrNoRows)
			} else {
				mock.ExpectQuery(`SELECT id, title, content, author_id, created_at, updated_at FROM posts ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
					WithArgs(tt.params.Limit, tt.params.Offset).
					WillReturnError(errors.New(tt.expectedErr))
			}

			posts, err := repo.GetAll(context.Background(), tt.params)
			assert.Error(t, err)
			assert.Nil(t, posts)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestPostRepository_Update_Success(t *testing.T) {
	tests := []struct {
		name string
		post *entity.Post
	}{
		{
			name: "Successful update post",
			post: &entity.Post{
				Id:        postId1,
				Title:     "Updated Title",
				Content:   "Updated Content",
				UpdatedAt: time.Now(),
			},
		},
		{
			name: "Successful update with different post data",
			post: &entity.Post{
				Id:        postId2,
				Title:     "Another Updated Title",
				Content:   "Another Updated Content",
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewPostRepository(&db.PostgresDB{DB: mockDB}, logger)

			mock.ExpectExec(`UPDATE posts SET title = \$1, content = \$2, updated_at = NOW\(\) WHERE id = \$3`).
				WithArgs(tt.post.Title, tt.post.Content, tt.post.Id).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = repo.Update(context.Background(), tt.post)
			assert.NoError(t, err)
		})
	}
}

func TestPostRepository_Update_Failed(t *testing.T) {
	tests := []struct {
		name        string
		post        *entity.Post
		expectedErr string
	}{
		{
			name: "Failed to update post - SQL error",
			post: &entity.Post{
				Id:      postId1,
				Title:   "Test Title",
				Content: "Test Content",
			},
			expectedErr: "failed to update post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewPostRepository(&db.PostgresDB{DB: mockDB}, logger)

			mock.ExpectExec(`UPDATE posts SET title = \$1, content = \$2, updated_at = NOW\(\) WHERE id = \$3`).
				WithArgs(tt.post.Title, tt.post.Content, tt.post.Id).
				WillReturnError(errors.New(tt.expectedErr))

			err = repo.Update(context.Background(), tt.post)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestPostRepository_Delete_Success(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		mockSetup func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Successful delete post with ID 1",
			id:   postId1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM posts WHERE id = \$1`).
					WithArgs(postId1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Successful delete post with ID 2",
			id:   postId2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM posts WHERE id = \$1`).
					WithArgs(postId2).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewPostRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			err = repo.Delete(context.Background(), tt.id)

			assert.NoError(t, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestPostRepository_Delete_Failed(t *testing.T) {
	tests := []struct {
		name        string
		id          uuid.UUID
		mockSetup   func(mock sqlmock.Sqlmock)
		expectedErr string
	}{
		{
			name: "Failed to delete post - SQL error",
			id:   postId1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM posts WHERE id = \$1`).
					WithArgs(postId1).
					WillReturnError(errors.New("failed to delete post"))
			},
			expectedErr: "failed to delete post",
		},
		{
			name: "Failed to delete post - No rows affected",
			id:   postId2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM posts WHERE id = \$1`).
					WithArgs(postId2).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: "no rows affected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewPostRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			err = repo.Delete(context.Background(), tt.id)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestPostRepository_GetTotalPosts_Success(t *testing.T) {
	tests := []struct {
		name          string
		expectedTotal int64
		mockSetup     func(mock sqlmock.Sqlmock)
	}{
		{
			name:          "Successful get total posts",
			expectedTotal: 10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(10)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM posts`).WillReturnRows(rows)
			},
		},
		{
			name:          "Successful get total posts with zero count",
			expectedTotal: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM posts`).WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewPostRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			total, err := repo.GetTotalPosts(context.Background())

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedTotal, total)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostRepository_GetTotalPosts_Failed(t *testing.T) {
	tests := []struct {
		name        string
		mockError   error
		expectedErr string
		mockSetup   func(mock sqlmock.Sqlmock)
	}{
		{
			name:        "Failed to get total posts - SQL error",
			mockError:   errors.New("failed to get total posts"),
			expectedErr: "failed to get total posts",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM posts`).WillReturnError(errors.New("failed to get total posts"))
			},
		},
		{
			name:        "Failed to get total posts - No rows affected",
			mockError:   sql.ErrNoRows,
			expectedErr: "sql: no rows in result set",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM posts`).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:        "Failed to get total posts - Connection issue",
			mockError:   errors.New("database connection error"),
			expectedErr: "database connection error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM posts`).WillReturnError(errors.New("database connection error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			logger := logrus.New()
			repo := postgres.NewPostRepository(&db.PostgresDB{DB: mockDB}, logger)

			tt.mockSetup(mock)

			total, err := repo.GetTotalPosts(context.Background())

			if tt.mockError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Zero(t, total)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
