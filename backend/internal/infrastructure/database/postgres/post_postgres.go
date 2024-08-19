package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/pkg/db"
)

type PostRepository struct {
	db     *db.PostgresDB
	logger *logrus.Logger
}

func NewPostRepository(db *db.PostgresDB, logger *logrus.Logger) *PostRepository {
	return &PostRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostRepository) CreatePost(ctx context.Context, post *entity.NewPost) (*entity.Post, error) {
	query := `INSERT INTO posts (id, title, content, author_id, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, NOW(), NOW()) 
              RETURNING id, title, content, author_id, created_at, updated_at`

	postID := uuid.New()
	var createdPost entity.Post
	err := r.db.QueryRowContext(ctx, query, postID, post.Title, post.Content, post.AuthorId).Scan(
		&createdPost.Id, &createdPost.Title, &createdPost.Content,
		&createdPost.AuthorId, &createdPost.CreatedAt, &createdPost.UpdatedAt,
	)

	if err != nil {
		r.logger.WithError(err).Error("Failed to create post")
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return &createdPost, nil
}

func (r *PostRepository) GetPostById(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at FROM posts WHERE id = $1`
	var post entity.Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		r.logger.WithError(err).Error("Failed to get post")
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return &post, nil
}

func (r *PostRepository) GetAll(ctx context.Context, params *entity.Pagination) ([]*entity.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at FROM posts`

	r.logger.WithField("params", params).Info("GetAll posts")

	if params.Sort != "" {
		r.logger.Info("params.Sort", params.Sort)

		sortField, sortOrder := "", ""

		switch params.Sort {
		case "created_at_asc":
			sortField = "created_at"
			sortOrder = "asc"
		case "created_at_desc":
			sortField = "created_at"
			sortOrder = "desc"
		case "title_asc":
			sortField = "title"
			sortOrder = "asc"
		case "title_desc":
			sortField = "title"
			sortOrder = "desc"
		default:
			return nil, fmt.Errorf("invalid sort format: %s", params.Sort)
		}

		query += fmt.Sprintf(" ORDER BY %s %s", sortField, sortOrder)
	} else {
		query += " ORDER BY created_at DESC"
	}

	r.logger.WithField("query", query).Info("Final query")

	query += " LIMIT $1 OFFSET $2"

	rows, err := r.db.QueryContext(ctx, query, params.Limit, params.Offset)
	if err != nil {
		r.logger.WithError(err).Error("Failed to get posts")
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		var post entity.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.AuthorId, &post.CreatedAt, &post.UpdatedAt); err != nil {
			r.logger.WithError(err).Error("Failed to scan post")
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, &post)
	}

	for _, post := range posts {
		r.logger.Infof("Post: %v, CreatedAt: %v, Title: %v", post.Id, post.CreatedAt, post.Title)
	}

	return posts, nil
}

func (r *PostRepository) Update(ctx context.Context, post *entity.Post) error {
	query := `UPDATE posts SET title = $1, content = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.Id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to update post")
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

func (r *PostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *PostRepository) GetTotalPosts(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM posts`
	var total int64
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		r.logger.WithError(err).Error("Failed to get total posts")
		return 0, fmt.Errorf("failed to get total posts: %w", err)
	}
	return total, nil
}
