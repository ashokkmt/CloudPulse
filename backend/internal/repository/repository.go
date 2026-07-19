package repository

import (
	"context"
	"cloudpulse/backend/internal/model"
)

type Repository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	GetUserByID(ctx context.Context, id int) (model.User, error)

	CreateTask(ctx context.Context, userID int, title string) (model.Task, error)
	GetTasksByUserID(ctx context.Context, userID int) ([]model.Task, error)
	UpdateTask(ctx context.Context, id, userID int, title *string, done *bool) (model.Task, error)
	DeleteTask(ctx context.Context, id, userID int) error
}
