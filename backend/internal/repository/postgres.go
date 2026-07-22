package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"cloudpulse/backend/internal/model"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("postgres")

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, email, passwordHash string) (model.User, error) {
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()

	var u model.User
	err := r.db.QueryRowContext(ctx, 
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash, created_at",
		email, passwordHash,
	).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	return u, err
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	ctx, span := tracer.Start(ctx, "GetUserByEmail")
	defer span.End()

	var u model.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, email, password_hash, created_at FROM users WHERE email = $1",
		email,
	).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	return u, err
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id int) (model.User, error) {
	ctx, span := tracer.Start(ctx, "GetUserByID")
	defer span.End()

	var u model.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, email, password_hash, created_at FROM users WHERE id = $1",
		id,
	).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	return u, err
}

func (r *PostgresRepository) CreateTask(ctx context.Context, userID int, title string) (model.Task, error) {
	ctx, span := tracer.Start(ctx, "CreateTask")
	defer span.End()

	var t model.Task
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO tasks (user_id, title) VALUES ($1, $2) RETURNING id, user_id, title, done, attachment_url, attachment_name, attachment_size, mime_type, thumbnail_url, created_at",
		userID, title,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Done, &t.AttachmentURL, &t.AttachmentName, &t.AttachmentSize, &t.MimeType, &t.ThumbnailURL, &t.CreatedAt)
	if err != nil {
		slog.Error("CreateTask Scan Error", slog.String("error", err.Error()))
	}
	return t, err
}

func (r *PostgresRepository) GetTasksByUserID(ctx context.Context, userID int) ([]model.Task, error) {
	ctx, span := tracer.Start(ctx, "GetTasksByUserID")
	defer span.End()

	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, title, done, attachment_url, attachment_name, attachment_size, mime_type, thumbnail_url, created_at FROM tasks WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Done, &t.AttachmentURL, &t.AttachmentName, &t.AttachmentSize, &t.MimeType, &t.ThumbnailURL, &t.CreatedAt); err != nil {
			slog.Error("GetTasksByUserID Scan Error", slog.String("error", err.Error()))
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []model.Task{} // Return empty array instead of null for JSON
	}
	return tasks, rows.Err()
}

func (r *PostgresRepository) GetTaskByID(ctx context.Context, id, userID int) (model.Task, error) {
	ctx, span := tracer.Start(ctx, "GetTaskByID")
	defer span.End()

	var t model.Task
	err := r.db.QueryRowContext(ctx, "SELECT id, user_id, title, done, attachment_url, attachment_name, attachment_size, mime_type, thumbnail_url, created_at FROM tasks WHERE id = $1 AND user_id = $2", id, userID).
		Scan(&t.ID, &t.UserID, &t.Title, &t.Done, &t.AttachmentURL, &t.AttachmentName, &t.AttachmentSize, &t.MimeType, &t.ThumbnailURL, &t.CreatedAt)
	return t, err
}

func (r *PostgresRepository) UpdateTask(ctx context.Context, id, userID int, title *string, done *bool) (model.Task, error) {
	ctx, span := tracer.Start(ctx, "UpdateTask")
	defer span.End()

	// First fetch the existing task
	t, err := r.GetTaskByID(ctx, id, userID)
	if err != nil {
		return t, err
	}

	if title != nil {
		t.Title = *title
	}
	if done != nil {
		t.Done = *done
	}

	err = r.db.QueryRowContext(ctx,
		"UPDATE tasks SET title = $1, done = $2 WHERE id = $3 AND user_id = $4 RETURNING id, user_id, title, done, attachment_url, attachment_name, attachment_size, mime_type, thumbnail_url, created_at",
		t.Title, t.Done, id, userID,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Done, &t.AttachmentURL, &t.AttachmentName, &t.AttachmentSize, &t.MimeType, &t.ThumbnailURL, &t.CreatedAt)

	return t, err
}

func (r *PostgresRepository) UpdateTaskAttachment(ctx context.Context, id, userID int, url, name *string, size *int64, mime *string) (model.Task, error) {
	ctx, span := tracer.Start(ctx, "UpdateTaskAttachment")
	defer span.End()

	var t model.Task
	err := r.db.QueryRowContext(ctx,
		"UPDATE tasks SET attachment_url = $1, attachment_name = $2, attachment_size = $3, mime_type = $4 WHERE id = $5 AND user_id = $6 RETURNING id, user_id, title, done, attachment_url, attachment_name, attachment_size, mime_type, thumbnail_url, created_at",
		url, name, size, mime, id, userID,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Done, &t.AttachmentURL, &t.AttachmentName, &t.AttachmentSize, &t.MimeType, &t.ThumbnailURL, &t.CreatedAt)

	return t, err
}

func (r *PostgresRepository) UpdateTaskThumbnail(ctx context.Context, id, userID int, thumbnailURL *string) error {
	ctx, span := tracer.Start(ctx, "UpdateTaskThumbnail")
	defer span.End()

	_, err := r.db.ExecContext(ctx, "UPDATE tasks SET thumbnail_url = $1 WHERE id = $2 AND user_id = $3", thumbnailURL, id, userID)
	return err
}

func (r *PostgresRepository) DeleteTask(ctx context.Context, id, userID int) error {
	ctx, span := tracer.Start(ctx, "DeleteTask")
	defer span.End()

	res, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("task not found")
	}
	return nil
}
