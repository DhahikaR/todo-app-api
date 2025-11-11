package repository

import (
	"context"
	"todo-app-api/models/domain"

	"gorm.io/gorm"
)

type TodoRepository interface {
	Save(ctx context.Context, tx *gorm.DB, todo domain.Todo) domain.Todo
	Update(ctx context.Context, tx *gorm.DB, todo domain.Todo) domain.Todo
	Delete(ctx context.Context, tx *gorm.DB, todo domain.Todo)
	FindById(ctx context.Context, tx *gorm.DB, todoId int) (domain.Todo, error)
	FindAll(ctx context.Context, tx *gorm.DB) []domain.Todo
}
