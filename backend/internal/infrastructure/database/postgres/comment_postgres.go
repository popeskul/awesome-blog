package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	"github.com/popeskul/awesome-blog/backend/pkg/db"
)

type CommentRepository struct {
	db     *db.PostgresDB
	logger *logrus.Logger
}

func NewCommentRepository(db *db.PostgresDB, logger *logrus.Logger) *CommentRepository {
	return &CommentRepository{
		db:     db,
		logger: logger,
	}
}

func (r *CommentRepository) CreateComment(ctx context.Context, comment *entity.NewComment) (*entity.Comment, error) {
	query := `
        INSERT INTO comments (id, post_id, author_id, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id, post_id, author_id, content, created_at, updated_at
    `

	commentID := uuid.New()
	var createdComment entity.Comment
	err := r.db.QueryRowContext(ctx, query,
		commentID, comment.PostId, comment.AuthorId, comment.Content,
	).Scan(
		&createdComment.Id, &createdComment.PostId, &createdComment.AuthorId,
		&createdComment.Content, &createdComment.CreatedAt, &createdComment.UpdatedAt,
	)

	if err != nil {
		r.logger.WithError(err).Error("Failed to create comment")
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &createdComment, nil
}

func (r *CommentRepository) GetCommentById(ctx context.Context, id uuid.UUID) (*entity.Comment, error) {
	query := `
        SELECT id, post_id, author_id, content, created_at, updated_at
        FROM comments
        WHERE id = $1
    `

	var comment entity.Comment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.Id, &comment.PostId, &comment.AuthorId,
		&comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("comment not found")
		}

		r.logger.WithError(err).Error("Failed to get comment")

		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return &comment, nil
}

func (r *CommentRepository) GetComments(ctx context.Context, postID uuid.UUID, params *entity.Pagination) ([]*entity.Comment, error) {
	query := `SELECT id, post_id, author_id, content, created_at, updated_at FROM comments WHERE post_id = $1`

	r.logger.WithFields(logrus.Fields{
		"postID": postID,
		"params": params,
	}).Info("GetComments called")

	if params.Sort != "" {
		r.logger.WithField("sort", params.Sort).Info("Sorting parameter")

		sortField, sortOrder := "", ""

		switch params.Sort {
		case "created_at_asc":
			sortField = "created_at"
			sortOrder = "asc"
		case "created_at_desc":
			sortField = "created_at"
			sortOrder = "desc"
		default:
			return nil, fmt.Errorf("invalid sort format: %s", params.Sort)
		}

		query += fmt.Sprintf(" ORDER BY %s %s", sortField, sortOrder)
	} else {
		query += " ORDER BY created_at DESC"
	}

	r.logger.WithField("query", query).Info("Final query before adding LIMIT and OFFSET")

	query += " LIMIT $2 OFFSET $3"

	r.logger.WithField("query", query).Info("Final query")

	rows, err := r.db.QueryContext(ctx, query, postID, params.Limit, params.Offset)
	if err != nil {
		r.logger.WithError(err).Error("Failed to get comments")
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		var comment entity.Comment
		err := rows.Scan(
			&comment.Id, &comment.PostId, &comment.AuthorId,
			&comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
		)
		if err != nil {
			r.logger.WithError(err).Error("Failed to scan comment")
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &comment)
	}

	for _, comment := range comments {
		r.logger.WithFields(logrus.Fields{
			"commentID": comment.Id,
			"createdAt": comment.CreatedAt,
			"authorID":  comment.AuthorId,
		}).Info("Comment details")
	}

	return comments, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, comment *entity.UpdateComment) error {
	query := `
        UPDATE comments
        SET content = $1, updated_at = NOW()
        WHERE id = $2
    `

	_, err := r.db.ExecContext(ctx, query, comment.Content, comment.Id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to update comment")
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) DeleteCommentById(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM comments WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete comment")
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) GetTotalCommentsByPostID(ctx context.Context, postID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`

	var total int
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&total)
	if err != nil {
		r.logger.WithError(err).Error("Failed to get total comments")
		return 0, fmt.Errorf("failed to get total comments: %w", err)
	}

	return total, nil
}
